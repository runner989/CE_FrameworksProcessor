package airtable

import (
	"cefp/database"
	"cefp/structs"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

const (
	//frameworks_TablesURL = "https://api.airtable.com/v0/meta/bases/appspojzJxIM9tUaC/tables" // this is the Framework Build table
	frameworksBaseURL  = "https://api.airtable.com/v0/appspojzJxIM9tUaC/tblRjgSEfrpsd4Llp"
	frameworksViewName = "All%20tasks%20grid"

	devMasterBase = "app5fTueYfRM65SzX"
	//tableViewsMetaURL = "https://api.airtable.com/v0/meta/bases/{baseId}/views"
	//tableViews        = "/views"
	//
	//getRecordURL = "https://api.airtable.com/v0/{baseId}/{tableIdOrName}/{recordId}"
)

func GetAirtableTablesAndViews(apiKey, baseID string) (string, error) {
	reqURL := fmt.Sprintf("https://api.airtable.com/v0/meta/bases/%s/tables", baseID)

	response, err := makeHTTPRequest(reqURL, apiKey)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	return response, nil
}

// GetFrameworksLookup function to read the Frameworks Build table on Airtable
func GetFrameworksLookup(apiKey string) ([]structs.Framework, error) {
	reqURL := fmt.Sprintf("%s?view=%s&Rand=%s", frameworksBaseURL, frameworksViewName, GenerateRandomString())
	done := false

	var allRecords []structs.Framework

	for !done {
		response, err := makeHTTPRequest(reqURL, apiKey)
		if err != nil {
			log.Fatalf("Error making request: %v", err)
			return allRecords, err
		}

		var airtableFrameworksResp structs.FrameworksResponse
		err = json.Unmarshal([]byte(response), &airtableFrameworksResp)
		if err != nil {
			log.Fatalf("Error parsing JSON: %v", err)
			return allRecords, err
		}
		//log.Printf("airtableFrameworksResp: %v", airtableFrameworksResp)

		// Append the records to the slice of all records
		allRecords = append(allRecords, airtableFrameworksResp.Records...)

		if airtableFrameworksResp.Offset == "" {
			done = true
		} else {
			reqURL = fmt.Sprintf("%s?offset=%s&view=%s&Rand=%s", frameworksBaseURL, airtableFrameworksResp.Offset, frameworksViewName, GenerateRandomString())
		}

	}
	return allRecords, nil
}

func GetFrameworkData(db *sql.DB, apiKey string, lr structs.FrameworkLookup) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}
	if apiKey == "" {
		return fmt.Errorf("no apiKey provided")
	}

	tableView := strings.ReplaceAll(lr.TableView.String, " ", "+")
	reqURL := fmt.Sprintf("https://api.airtable.com/v0/%s/%s?view=%s&Rand=%s", devMasterBase, lr.TableID.String, tableView, GenerateRandomString())
	//log.Printf("Getting framework data for %s", reqURL)

	done := false

	delQry := fmt.Sprintf("DELETE FROM Framework WHERE Framework='%s';", lr.TableName.String)
	_, err := db.Exec(delQry)
	if err != nil {
		//runtime.EventsEmit(ctx, "progress", fmt.Sprintf("Error deleting from Framework: %v", err))
		log.Printf("error deleting from Framework: %v", err)
		return fmt.Errorf("error deleting from Framework: %v", err)
	}

	for !done {
		response, err := makeHTTPRequest(reqURL, apiKey)
		if err != nil {
			log.Fatalf("Error making request: %v", err)
			return err
		}

		var airtableFrameworksResp structs.FrameworksResponse
		err = json.Unmarshal([]byte(response), &airtableFrameworksResp)
		if err != nil {
			return fmt.Errorf("unable to parse the framework data")
		}

		if strings.Contains(response, `"error":{`) {
			errorType := airtableFrameworksResp.Records[0].ID
			if errorType != "" {
				//runtime.EventsEmit(ctx, "progress", fmt.Sprintf("There is an error: %v", err))
				return fmt.Errorf("there is an error: %s", errorType)
			}
			if strings.Contains(response, "NOT_FOUND") {
				//runtime.EventsEmit(ctx, "progress", fmt.Sprintf("The framework was not found.: %v", err))
				return fmt.Errorf("the framework was not found. please check the name and try again")
			}
		}
		sortID := 0
		for _, record := range airtableFrameworksResp.Records {
			var testType string
			identifier, ok := record.Fields["Identifier"].(string)
			if !ok {
				//runtime.EventsEmit(ctx, "progress", fmt.Sprintf("Skipping record due to missing or invalid EvidenceID: %v", err))
				log.Printf("skipping record due to missing or invalid Identifier")
			}
			parentID, _ := record.Fields["ParentIdentifier"].(string)
			displayName, _ := record.Fields["DisplayName"].(string)
			description, _ := record.Fields["Description"].(string)
			guidance, _ := record.Fields["Guidance"].(string)
			tags, _ := record.Fields["Tags"].(string)
			promptID, _ := record.Fields["Prompt ID"].(int)
			controlNarrative, _ := record.Fields["ControlNarrativeAIPromptTemplateId"].(int)
			frameworkName := lr.TableName.String
			sortID += 1
			switch v := record.Fields["TestType"].(type) {
			case []interface{}:
				if len(v) > 0 {
					testType, _ = v[0].(string)
				}
			case string:
				testType = v
			case nil:
				testType = ""
			default:
				log.Printf("unknown type for TestType: %T", v)
			}

			frameworkRecord := structs.FrameworkRecord{
				SortID:           sortID,
				PromptID:         promptID,
				ControlNarrative: controlNarrative,
				FrameworkName:    frameworkName,
				Identifier:       identifier,
				ParentID:         parentID,
				DisplayName:      displayName,
				Description:      description,
				Guidance:         guidance,
				Tags:             tags,
				TestType:         testType,
			}

			// Insert records
			err := database.InsertFrameworkRecord(db, frameworkRecord)
			if err != nil {
				log.Printf("skipping Framework %s Identifier %s due to error: %v", identifier, frameworkName, err)
				//runtime.EventsEmit(ctx, "progress", fmt.Sprintf("Skipping EvidenceID %d due to error: %v", int(evidenceID), err))
				continue
			}

			if airtableFrameworksResp.Offset == "" {
				done = true
			} else {
				reqURL = fmt.Sprintf("https://api.airtable.com/v0/%s/%s?offset=%s&view=%s&Rand=%s", devMasterBase, lr.TableID.String, airtableFrameworksResp.Offset, tableView, GenerateRandomString())
			}
		}
	}
	return nil
}
