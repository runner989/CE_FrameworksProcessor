package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// func DisplayRecords(db *sql.DB) error {
// 	if db == nil {
// 		return fmt.Errorf("database connection is nil")
// 	}
// 	row, err := db.Query("SELECT * FROM framework ORDER BY sortID")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer row.Close()

// 	for row.Next() {
// 		var sortid int
// 		var identifier string
// 		var parentident string
// 		var description string
// 		var displayName string
// 		var guidance string
// 		var recommendations string
// 		var requirementType string
// 		var pandpPromtId int
// 		var controlNarrativeId int
// 		var testType string
// 		var framework string
// 		var frameworkId int
// 		row.Scan(&sortid, &identifier, &parentident, &description, &displayName, &guidance, &recommendations, &requirementType, &pandpPromtId,
// 			&controlNarrativeId, &testType, &framework, &frameworkId)
// 		log.Println(sortid, identifier, parentident, description, displayName, guidance, recommendations, requirementType, pandpPromtId, controlNarrativeId, testType, framework, frameworkId)
// 	}
// 	return err
// }

func GetMissingFrameworks(db *sql.DB) ([]string, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	var uniqueFrameworks []string

	rows, err := db.Query("SELECT DISTINCT Framework FROM Framework WHERE Framework NOT IN (SELECT AirtableFramework FROM Framework_Lookup);")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var framework string
		if err := rows.Scan(&framework); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		uniqueFrameworks = append(uniqueFrameworks, framework)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return uniqueFrameworks, nil
}
