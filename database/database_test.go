package database

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestGetMissingFrameworks(t *testing.T) {
	db, err := sql.Open("sqlite3", "../cefp.db")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	records, err := GetMissingFrameworks(db)
	if err != nil {
		t.Fatalf("GetMissingFrameworks returned an error: %v", err)
	}

	// Check that records were returned
	if len(records) == 0 {
		t.Fatal("No records returned from GetFrameworksLookup")
	}

	// Optionally, print the records for debugging
	t.Logf("Fetched %d records", len(records))
	// t.Logf("Records: %+v", records)
}
