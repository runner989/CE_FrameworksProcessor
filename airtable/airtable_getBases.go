package airtable

import (
	"cefp/structs"
	"encoding/json"
	"fmt"
	"log"
)

const airtableBasesURL = "https://api.airtable.com/v0/meta/bases"

func GetAirtableBases(apiKey string) ([]structs.Base, error) {
	reqURL := airtableBasesURL
	done := false

	var allBases []structs.Base

	for !done {
		response, err := makeHTTPRequest(reqURL, apiKey)
		if err != nil {
			log.Fatalf("Error making request: %v", err)
			return allBases, err
		}

		var airtableBasesResp structs.AirtableBases
		err = json.Unmarshal([]byte(response), &airtableBasesResp)
		if err != nil {
			log.Fatalf("Error parsing JSON: %v", err)
			return allBases, err
		}

		// Append the records to the slice of all records
		allBases = append(allBases, airtableBasesResp.Bases...)

		if airtableBasesResp.Offset == "" {
			done = true
		} else {
			reqURL = fmt.Sprintf("%s?offset=%s&Rand=%s", frameworksBaseURL, airtableBasesResp.Offset, GenerateRandomString())
		}

	}
	return allBases, nil
}
