package airtable

import (
	"cefp/secret"
	"testing"
)

func TestGetFrameworksLookup(t *testing.T) {
	// Retrieve the API key from an environment variable
	apiKey := secret.APIKey
	if apiKey == "" {
		t.Fatal("AIRTABLE_API_KEY environment variable is not set")
	}

	// Call the GetFrameworksLookup function
	records, err := GetFrameworksLookup(apiKey)
	if err != nil {
		t.Fatalf("Error fetching frameworks lookup: %v", err)
	}

	// Check that records were returned
	if len(records) == 0 {
		t.Fatal("No records returned from GetFrameworksLookup")
	}

	// Optionally, print the records for debugging
	t.Logf("Fetched %d records", len(records))
	// t.Logf("Records: %+v", records)
}
