package airtable

import (
	"cefp/database"
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

type Record struct {
	ID          string                 `json:"id"`
	CreatedTime string                 `json:"createdTime"`
	Fields      map[string]interface{} `json:"fields"`
}

type Table struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	PrimaryFieldId  string                 `json:"primaryFieldId"`
	Fields          map[string]interface{} `json:"fields"`
	PermissionLevel string                 `json:"permissionLevel"`
}

type TablesResponse struct {
	Tables []Table `json:"tables"`
}

type Framework struct {
	ID          string                 `json:"id"`
	CreatedTime string                 `json:"createdTime"`
	Fields      map[string]interface{} `json:"fields"`
}

type FrameworksTable struct {
}

type ViewsRoot struct {
	Views []View `json:"views"`
}

type View struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	PersonalForUserId string `json:"personalForUserId,omitempty"`
	Type              string `json:"type"`
}

type FrameworksResponse struct {
	Records []Framework `json:"records"`
	Offset  string      `json:"offset,omitempty"`
}

//func GetFrameworksTables(apiKey string) error {
//	reqURL := frameworks_TablesURL //fmt.Sprintf("%s?view=%s&Rand=%s", frameworks_BaseURL, frameworksViewName, generateRandomString())
//
//	done := false
//
//	for {
//		response, err := makeHTTPRequest(reqURL, apiKey)
//		if err != nil {
//			log.Fatalf("Error making request: %v", err)
//			return err
//		}
//		fmt.Println(response)
//		// strResponses = strResponses + response
//		var airtableResp TablesResponse
//		err = json.Unmarshal([]byte(response), &airtableResp)
//		if err != nil {
//			log.Fatalf("Error parsing JSON: %v", err)
//			return err
//		}
//
//		for _, table := range airtableResp.Tables {
//			fmt.Printf("%s", table)
//		}
//
//	}
//
//	return nil
//}

func GetAirtableTablesAndViews(apiKey string) (string, error) {
	baseID := "app5fTueYfRM65SzX"
	reqURL := fmt.Sprintf("https://api.airtable.com/v0/meta/bases/%s/tables", baseID)

	response, err := makeHTTPRequest(reqURL, apiKey)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	return response, nil
}

// GetFrameworksLookup function to read the Frameworks Build table on Airtable
func GetFrameworksLookup(apiKey string) ([]Framework, error) {
	reqURL := fmt.Sprintf("%s?view=%s&Rand=%s", frameworksBaseURL, frameworksViewName, GenerateRandomString())
	done := false

	var allRecords []Framework

	for !done {
		response, err := makeHTTPRequest(reqURL, apiKey)
		if err != nil {
			log.Fatalf("Error making request: %v", err)
			return allRecords, err
		}

		var airtableFrameworksResp FrameworksResponse
		err = json.Unmarshal([]byte(response), &airtableFrameworksResp)
		if err != nil {
			log.Fatalf("Error parsing JSON: %v", err)
			return allRecords, err
		}

		// Append the records to the slice of all records
		allRecords = append(allRecords, airtableFrameworksResp.Records...)
		// for _, record := range airtableResp.Records {
		// 	fmt.Printf("%s", record)
		// }

		if airtableFrameworksResp.Offset == "" {
			done = true
		} else {
			reqURL = fmt.Sprintf("%s?offset=%s&view=%s&Rand=%s", frameworksBaseURL, airtableFrameworksResp.Offset, frameworksViewName, GenerateRandomString())
		}

	}
	return allRecords, nil
}

func GetFrameworkData(db *sql.DB, apiKey, tableName, tableID, tableView string) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}
	if apiKey == "" {
		return fmt.Errorf("no apiKey provided")
	}

	reqURL := fmt.Sprintf("https://api.airtable.com/v0/%s/%s?view=%s&Rand=%s", devMasterBase, tableID, tableView, GenerateRandomString())

	done := false

	delQry := fmt.Sprintf("DELETE FROM Framework WHERE Framework='%s';", tableName)
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

		//log.Print(response)

		var airtableFrameworksResp FrameworksResponse
		err = json.Unmarshal([]byte(response), &airtableFrameworksResp)

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
			//log.Println(record)
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
			frameworkName := tableName
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

			//message := fmt.Sprintf("Processing EvidenceID: %d, Evidence: %s", int(evidenceID), evidenceTitle)
			//runtime.EventsEmit(ctx, "progress", message)
			//log.Printf("identifier: %s, parentID: %s, displayName: %s, testType: %s", identifier, parentID, displayName, testType)
			// Insert records
			err := database.InsertFrameworkRecord(db, sortID, promptID, controlNarrative, frameworkName, identifier, parentID, displayName, description, guidance, tags, testType)
			// sortID int, frameworkName, identifier, parentID, displayName, description, guidance, tags, testType, policyID, controlNarrative
			if err != nil {
				log.Printf("skipping Framework %s Identifier %s due to error: %v", identifier, frameworkName, err)
				//runtime.EventsEmit(ctx, "progress", fmt.Sprintf("Skipping EvidenceID %d due to error: %v", int(evidenceID), err))
				continue
			}

			if airtableFrameworksResp.Offset == "" {
				done = true
			} else {
				reqURL = fmt.Sprintf("https://api.airtable.com/v0/%s/%s?offset=%s&view=%s&Rand=%s", devMasterBase, tableID, airtableFrameworksResp.Offset, tableView, GenerateRandomString())
			}
		}
	}
	return nil
}
