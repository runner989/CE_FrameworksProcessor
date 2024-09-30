package airtable

import (
	"cefp/structs"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/mattn/go-sqlite3"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Config struct {
	AdditionalSkipFields []string `yaml:"additional_skip_fields"`
}

const (
	baseURL = "https://api.airtable.com/v0/applEdk0gS7gMZ9o7/tbl6gMhn2VNnl4cOA"
	//baseURL = "https://api.airtable.com/v0/app5fTueYfRM65SzX/tblYaPqsXmknYtwIx"
	evidenceViewName = "Active+Break+Out+View"
	//evidenceViewName = "View%20For%20Export"
)

func loadAdditionalSkipFields() (map[string]struct{}, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config.yaml: %v", err)
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config.yaml: %v", err)
	}

	additionalSkipFields := make(map[string]struct{})
	for _, field := range config.AdditionalSkipFields {
		additionalSkipFields[field] = struct{}{}
	}
	return additionalSkipFields, nil
}

func safeString(value interface{}) sql.NullString {
	if value == nil {
		return sql.NullString{String: "", Valid: false}
	}

	str, ok := value.(string)
	if !ok {
		return sql.NullString{String: "", Valid: false}
	}

	return sql.NullString{String: str, Valid: true}
}

func ReadAPIEvidenceTable(ctx context.Context, db *sql.DB, apiKey string) error {
	skipFields := map[string]struct{}{
		"EvidenceID":                             {},
		"Evidence Title":                         {},
		"Requirement":                            {},
		"Description_FromEvidence":               {},
		"AnecdotesEvidenceIds":                   {},
		"Control Families CCM (from Card Title)": {},
		"Card Title":                             {},
		"Sync Source":                            {},
		"Evidence Type":                          {},
		"Priority":                               {},
	}

	additionalSkipFields, err := loadAdditionalSkipFields()
	if err != nil {
		// log.Printf("proceeding with hardcoded skipped fields only.", err)
	}

	for field := range additionalSkipFields {
		skipFields[field] = struct{}{}
	}

	reqURL := fmt.Sprintf("%s?view=%s&Rand=%s", baseURL, evidenceViewName, GenerateRandomString())
	//log.Printf("Fetching Evidence Table from %s", reqURL)

	done := false

	// allResponses = append(allResponses, airtableResp)
	_, err = db.Exec("DELETE FROM Mapping")
	if err != nil {
		runtime.EventsEmit(ctx, "progress", fmt.Sprintf("Error deleting from Mapping: %v", err))
		return fmt.Errorf("error deleting from Mapping: %v", err)
	}
	_, err = db.Exec("DELETE FROM Evidence")
	if err != nil {
		runtime.EventsEmit(ctx, "progress", fmt.Sprintf("Error deleting from Evidence: %v", err))
		return fmt.Errorf("error deleting from Evidence: %v", err)
	}

	runtime.EventsEmit(ctx, "progress", "Cleared existing evidence and mapping data.")

	for !done {
		response, err := makeHTTPRequest(reqURL, apiKey)
		if err != nil {
			runtime.EventsEmit(ctx, "progress", fmt.Sprintf("Error making request: %v", err))
			log.Fatalf("Error making request: %v", err)
			return err
		}

		//log.Printf("response: %s", response)
		// strResponses = strResponses + response
		var airtableResp structs.AirtableResponse
		err = json.Unmarshal([]byte(response), &airtableResp)
		if err != nil {
			runtime.EventsEmit(ctx, "progress", fmt.Sprintf("Error parsing JSON: %v", err))
			log.Fatalf("Error parsing JSON: %v", err)
			return err
		}

		if strings.Contains(response, `"error":{`) {
			errorType := airtableResp.Records[0].ID
			if errorType != "" {
				runtime.EventsEmit(ctx, "progress", fmt.Sprintf("There is an error: %v", err))
				return fmt.Errorf("there is an error: %s", errorType)
			}
			if strings.Contains(response, "NOT_FOUND") {
				runtime.EventsEmit(ctx, "progress", fmt.Sprintf("The framework was not found.: %v", err))
				return fmt.Errorf("the framework was not found. please check the name and try again")
			}
		}

		for _, record := range airtableResp.Records {
			evidenceID, ok := record.Fields["EvidenceID"].(float64)
			if !ok {
				runtime.EventsEmit(ctx, "progress", fmt.Sprintf("Skipping record due to missing or invalid EvidenceID: %v", err))
				log.Printf("skipping record due to missing or invalid EvidenceID")
				continue
			}

			var anecdotesEvidenceIds string
			switch v := record.Fields["AnecdotesEvidenceIds"].(type) {
			case []interface{}:
				if len(v) > 0 {
					anecdotesEvidenceIds, _ = v[0].(string)
				}
			case string:
				anecdotesEvidenceIds = v
			case nil:
				anecdotesEvidenceIds = ""
			default:
				log.Printf("unknown type for AnecdotesEvidenceIds: %T", v)
			}

			evidenceRecord := structs.EvidenceRecord{
				EvidenceID:           int(evidenceID),
				EvidenceTitle:        safeString(record.Fields["Evidence Title"]),
				Description:          safeString(record.Fields["Description_FromEvidence"]),
				Requirement:          safeString(record.Fields["Requirement"]),
				AnecdotesEvidenceIds: safeString(anecdotesEvidenceIds),
				Priority:             safeString(record.Fields["Priority"]),
				EvidenceType:         safeString(record.Fields["Evidence Type"]),
			}

			message := fmt.Sprintf("Processing EvidenceID: %d, Evidence: %s", evidenceRecord.EvidenceID, evidenceRecord.EvidenceTitle.String)
			runtime.EventsEmit(ctx, "progress", message)

			// Insert CE Framework mapping
			err = insertMappingRecord(db, int(evidenceID), "CE Framework", strings.Trim(evidenceRecord.Requirement.String, " "))
			if err != nil {
				log.Printf("skipping CE Framework mapping due to error: %v", err)
				runtime.EventsEmit(ctx, "progress", fmt.Sprintf("Skipping CE Framework mapping Evidence ID: %d due to error: %v", int(evidenceID), err))
			}

			// Insert records
			err := insertEvidenceRecord(db, evidenceRecord) // int(evidenceID), evidenceTitle, description, anecdotesIds, priority, evidenceType)
			if err != nil {
				log.Printf("skipping EvidenceID %d due to error: %v", int(evidenceID), err)
				runtime.EventsEmit(ctx, "progress", fmt.Sprintf("Skipping EvidenceID %d due to error: %v", int(evidenceID), err))
				continue
			}

			for key, value := range record.Fields {
				if _, skip := skipFields[key]; skip {
					continue
				}

				// split value and iterate through each with insertMappingRecord:
				values := strings.Split(fmt.Sprintf("%v", value), ", ")
				for i := 0; i < len(values); i++ {
					err = insertMappingRecord(db, int(evidenceID), key, fmt.Sprintf("%v", values[i]))
					if err != nil {
						log.Printf("skipping dynamic field %s for EvidenceID %d due to error: %v", key, int(evidenceID), err)
						runtime.EventsEmit(ctx, "progress", fmt.Sprintf("Skipping dynamic field %s for EvidenceID %d due to error: %v", key, int(evidenceID), err))
					}
				}
			}

		}

		if airtableResp.Offset == "" {
			done = true
		} else {
			reqURL = fmt.Sprintf("%s?offset=%s&view=%s&Rand=%s", baseURL, airtableResp.Offset, evidenceViewName, GenerateRandomString())
		}
	}

	runtime.EventsEmit(ctx, "progress", "Done updating Evidence and Mapping tables!")
	return nil
}

func insertEvidenceRecord(db *sql.DB, er structs.EvidenceRecord) error { //evidenceID int, evidenceTitle, description, anecdotesIds, priority, evidenceType string) error {
	_, err := db.Exec("INSERT INTO Evidence (EvidenceID, Evidence, Description, AnecdotesEvidenceIds, Priority, EvidenceType) VALUES (?, ?, ?, ?, ?, ?)",
		er.EvidenceID, er.EvidenceTitle.String, er.Description.String, er.AnecdotesEvidenceIds.String, er.Priority.String, er.EvidenceType.String)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			// Handle UNIQUE constraint violation (duplicate EvidenceID)
			return fmt.Errorf("duplicate EvidenceID: %d", er.EvidenceID)
		}
		return fmt.Errorf("error inserting Evidence: %v", err)
	}
	return nil
}

func insertMappingRecord(db *sql.DB, evidenceID int, framework, requirement string) error {
	_, err := db.Exec("INSERT INTO Mapping (EvidenceID, Framework, Requirement) VALUES (?, ?, ?)",
		evidenceID, framework, requirement)
	if err != nil {
		return fmt.Errorf("error inserting into Mapping: %v", err)
	}
	return nil
}

func makeHTTPRequest(reqURL, apiKey string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "error creating HTTP request:", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")

	resp, err := client.Do(req)
	if err != nil {
		return "error performing HTTP request", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "error reading response body", err
	}

	return string(body), nil
}

func GenerateRandomString() string {
	return time.Now().Format("20060102150405") + fmt.Sprintf("%.0f", float64(time.Now().UnixNano())/1e6)
}
