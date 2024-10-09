package airtable

import (
	"cefp/database"
	"cefp/structs"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
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

func GetFrameworkData(db, memDB *sql.DB, apiKey string, lr structs.FrameworkLookup) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}
	if apiKey == "" {
		return fmt.Errorf("no apiKey provided")
	}

	tableView := strings.ReplaceAll(lr.TableView.String, " ", "+")
	reqURL := fmt.Sprintf("https://api.airtable.com/v0/%s/%s?view=%s&Rand=%s", devMasterBase, lr.TableID.String, tableView, GenerateRandomString())

	done := false
	delQry := fmt.Sprintf("DELETE FROM Framework WHERE Framework='%s';", lr.MappedName.String)
	_, err := db.Exec(delQry)
	_, err = memDB.Exec(delQry)
	if err != nil {
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
			var parentID string
			switch v := record.Fields["ParentIdentifier"].(type) {
			case []interface{}:
				if len(v) > 0 {
					parentID, _ = v[0].(string)
				}
			case string:
				parentID = v
			case nil:
				parentID = ""
			default:
				log.Printf("unknown type for ParentIdentifier: %T", v)
				parentID = ""
			}
			//parentID, _ := record.Fields["ParentIdentifier"].(string)
			var displayName string
			switch v := record.Fields["DisplayName"].(type) {
			case []interface{}:
				if len(v) > 0 {
					displayName, _ = v[0].(string)
				}
			case string:
				displayName = v
			case nil:
				displayName = ""
			default:
				log.Printf("unknown type for DisplayName: %T", v)
				displayName = ""
			}
			var description string
			switch v := record.Fields["Description"].(type) {
			case []interface{}:
				if len(v) > 0 {
					description, _ = v[0].(string)
				}
			case string:
				description = v
			case nil:
				description = ""
			default:
				log.Printf("unknown type for Description: %T", v)
				description = ""
			}
			//guidance, _ := record.Fields["Guidance"].(string)
			var guidance string
			switch v := record.Fields["Guidance"].(type) {
			case []interface{}:
				if len(v) > 0 {
					guidance, _ = v[0].(string)
				}
			case string:
				guidance = v
			case nil:
				guidance = ""
			default:
				log.Printf("unknown type for Guidance: %T", v)
				guidance = ""
			}
			tags, _ := record.Fields["Tags"].(string)
			//promptID, _ := record.Fields["Prompt ID"].(int)
			var promptIDStr string
			switch v := record.Fields["Prompt ID"].(type) {
			case []interface{}:
				if len(v) > 0 {
					promptIDStr = fmt.Sprintf("%v", v[0])
				} else {
					promptIDStr = ""
				}
			case int:
				promptIDStr = fmt.Sprintf("%v", v)
			case string:
				promptIDStr = v
			case nil:
				promptIDStr = ""
			default:
				log.Printf("unknown type for PromptID: %T", v)
				promptIDStr = ""
			}
			promptID, err := strconv.Atoi(promptIDStr)
			if err != nil {
				promptID = 0
			}

			controlNarrative, _ := record.Fields["ControlNarrativeAIPromptTemplateId"].(int)
			frameworkName := lr.MappedName.String
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
				testType = ""
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
			err = database.InsertFrameworkRecord(memDB, frameworkRecord)
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
