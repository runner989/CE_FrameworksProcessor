package airtable

import (
	"encoding/json"
	"fmt"
	"log"
)

const (
	frameworks_TablesURL = "https://api.airtable.com/v0/meta/bases/appspojzJxIM9tUaC/tables"
	frameworks_BaseURL   = "https://api.airtable.com/v0/appspojzJxIM9tUaC/tblRjgSEfrpsd4Llp"

	frameworksViewName = "All%20tasks%20grid"
)

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

type AirtableFrameworks struct {
	ID          string                 `json:"id"`
	CreatedTime string                 `json:"createdTime"`
	Fields      map[string]interface{} `json:"fields"`
}

type AirtableFrameworksResponse struct {
	Records []AirtableFrameworks `json:"records"`
	Offset  string               `json:"offset,omitempty"`
}

func GetFrameworksTables(apiKey string) error {
	reqURL := frameworks_TablesURL //fmt.Sprintf("%s?view=%s&Rand=%s", frameworks_BaseURL, frameworksViewName, generateRandomString())

	done := false

	for !done {
		response, err := makeHTTPRequest(reqURL, apiKey)
		if err != nil {
			log.Fatalf("Error making request: %v", err)
			return err
		}
		fmt.Println(response)
		// strResponses = strResponses + response
		var airtableResp TablesResponse
		err = json.Unmarshal([]byte(response), &airtableResp)
		if err != nil {
			log.Fatalf("Error parsing JSON: %v", err)
			return err
		}

		for _, table := range airtableResp.Tables {
			fmt.Printf("%s", table)
		}

	}

	return nil
}

func GetFrameworksLookup(apiKey string) ([]AirtableFrameworks, error) {
	reqURL := fmt.Sprintf("%s?view=%s&Rand=%s", frameworks_BaseURL, frameworksViewName, GenerateRandomString())
	// reqURL := frameworks_BaseURL
	done := false

	var allRecords []AirtableFrameworks

	for !done {
		response, err := makeHTTPRequest(reqURL, apiKey)
		if err != nil {
			log.Fatalf("Error making request: %v", err)
			return allRecords, err
		}

		var airtableFrameworksResp AirtableFrameworksResponse
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
			reqURL = fmt.Sprintf("%s?offset=%s&view=%s&Rand=%s", frameworks_BaseURL, airtableFrameworksResp.Offset, frameworksViewName, GenerateRandomString())
		}

	}
	return allRecords, nil
}