package airtable

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
)

const (
	baseURL          = "https://api.airtable.com/v0/applEdk0gS7gMZ9o7/tbl6gMhn2VNnl4cOA"
	evidenceViewName = "Active+Break+Out+View"
)

type Evidence struct {
	ID     string                 `json:"id"`
	Fields map[string]interface{} `json:"fields"`
}

//struct {
// 	EvidenceID           int    `json:"EvidenceID"`
// 	EvidenceTitle        string `json:"Evidence Title"`
// 	Description          string `json:"Description_FromEvidence"`
// 	AnecdotesEvidenceIds string `json:"AnecdotesEvidenceIds"`
// 	Priority             string `json:"Priority"`
// 	EvidenceType         string `json:"Evidence Type"`
// }

// CardTitle       string `json:"Card Title"`
// FrameworkID     string `json:"FrameworkdId"`
// Requirement     string `json:"Requirement"`
// RequirementType string `json:"RequirementType"`

// "EvidenceID" integer NOT NULL PRIMARY KEY,
// "Evidence" TEXT,
// "Description" TEXT,
// "AnecdotesEvidenceIds" TEXT,
// "Priority" TEXT,
// "EvidenceType" TEXT

type AirtableResponse struct {
	Records []Evidence `json:"records"`
	Offset  string     `json:"offset,omitempty"`
}

func ReadAPI_EvidenceTable(db *sql.DB, apiKey string) error {
	reqURL := fmt.Sprintf("%s?view=%s&Rand=%s", baseURL, evidenceViewName, GenerateRandomString())

	done := false

	// allResponses = append(allResponses, airtableResp)
	_, err := db.Exec("DELETE FROM tblMapping")
	if err != nil {
		return fmt.Errorf("error deleting from tblMapping: %v", err)
	}
	_, err = db.Exec("DELETE FROM Evidence")
	if err != nil {
		return fmt.Errorf("error deleting from Evidence: %v", err)
	}

	for !done {
		response, err := makeHTTPRequest(reqURL, apiKey)
		if err != nil {
			log.Fatalf("Error making request: %v", err)
			return err
		}
		// strResponses = strResponses + response
		var airtableResp AirtableResponse
		err = json.Unmarshal([]byte(response), &airtableResp)
		if err != nil {
			log.Fatalf("Error parsing JSON: %v", err)
			return err
		}

		if strings.Contains(response, `"error":{`) {
			// if strings.HasPrefix(response, `{"error"`) {
			errorType := airtableResp.Records[0].ID
			if errorType != "" {
				return fmt.Errorf("there is an error: %s", errorType)
			}
			if strings.Contains(response, "NOT_FOUND") {
				return fmt.Errorf("the framework was not found. please check the name and try again")
			}
		}

		for _, record := range airtableResp.Records {
			evidenceID, ok := record.Fields["EvidenceID"].(float64)
			if !ok {
				log.Printf("skipping record due to missing or invalid EvidenceID")
			}
			evidenceTitle, _ := record.Fields["Evidence Title"].(string)
			description, _ := record.Fields["Description_FromEvidence"].(string)
			anecdotesIds, _ := record.Fields["AnecdotesEvidenceIds"].(string)
			priority, _ := record.Fields["Priority"].(string)
			evidenceType, _ := record.Fields["Evidence Type"].(string)

			fmt.Printf("processing EvidenceID: %d, Evidence: %s \n", int(evidenceID), evidenceTitle)

			// Insert records
			err := insertEvidenceRecord(db, int(evidenceID), evidenceTitle, description, anecdotesIds, priority, evidenceType)
			if err != nil {
				log.Printf("skipping EvidenceID %d due to error: %v", int(evidenceID), err)
				continue
			}

			// This is wrong!  I need to fix this
			// I need to split the value (requirements) and iterate through the list and save each to the db
			for key, value := range record.Fields {
				if key == "EvidenceID" || key == "Evidence Title" || key == "Requirement" || key == "Description_FromEvidence" || key == "AnecdotesEvidenceIds" ||
					key == "Control Families CCM (from Card Title)" || key == "Card Title" || key == "Sync Source" || key == "Evidence Type" || key == "Priority" {
					continue
				}

				// split value and iterate through each with insertMappingRecord:
				values := strings.Split(fmt.Sprintf("%v", value), ", ")
				for i := 0; i < len(values); i++ {
					err = insertMappingRecord(db, int(evidenceID), key, fmt.Sprintf("%v", values[i]))
					if err != nil {
						log.Printf("skipping dynamic field %s for EvidenceID %d due to error: %v", key, int(evidenceID), err)
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

	// err := saveStringResponsesToFile(strResponses)
	// if err != nil {
	// 	return fmt.Errorf("error saving raw responses to file: %v", err)
	// }

	// err = saveResponsesToFile(allResponses)
	// if err != nil {
	// 	return fmt.Errorf("error saving responses to file: %v", err)
	// }

	fmt.Println("Done!")
	return nil
}

func insertEvidenceRecord(db *sql.DB, evidenceID int, evidenceTitle, description, anecdotesIds, priority, evidenceType string) error {
	_, err := db.Exec("INSERT INTO Evidence (EvidenceID, Evidence, Description, AnecdotesEvidenceIds, Priority, EvidenceType) VALUES (?, ?, ?, ?, ?, ?)",
		evidenceID, evidenceTitle, description, anecdotesIds, priority, evidenceType)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			// Handle UNIQUE constraint violation (duplicate EvidenceID)
			return fmt.Errorf("duplicate EvidenceID: %d", evidenceID)
		}
		return fmt.Errorf("error inserting Evidence: %v", err)
	}
	return nil
}

func insertMappingRecord(db *sql.DB, evidenceID int, framework, requirement string) error {
	_, err := db.Exec("INSERT INTO tblMapping (EvidenceID, Framework, Requirement) VALUES (?, ?, ?)",
		evidenceID, framework, requirement)
	if err != nil {
		return fmt.Errorf("error inserting into tblMapping: %v", err)
	}
	return nil
}

// func saveStringResponsesToFile(responses string) error {
// 	file, err := os.Create("raw_responses.json")
// 	if err != nil {
// 		return fmt.Errorf("error creating file: %v", err)
// 	}
// 	defer file.Close()
// 	_, err = file.Write([]byte(responses))
// 	if err != nil {
// 		return fmt.Errorf("error writing text to file: %v", err)
// 	}
// 	return nil
// }

// func saveResponsesToFile(responses []AirtableResponse) error {
// 	file, err := os.Create("responses.json")
// 	if err != nil {
// 		return fmt.Errorf("error creating file: %v", err)
// 	}
// 	defer file.Close()
// 	data, err := json.MarshalIndent(responses, "", "  ")
// 	if err != nil {
// 		return fmt.Errorf("error marshalling responses to JSON: %v", err)
// 	}

// 	_, err = file.Write(data)
// 	if err != nil {
// 		return fmt.Errorf("error writing data to file: %v", err)
// 	}
// 	return nil
// }

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
