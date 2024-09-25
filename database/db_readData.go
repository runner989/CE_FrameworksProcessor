package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func GetMissingFrameworks(db *sql.DB) ([]string, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	uniqueFrameworks := []string{}

	rows, err := db.Query("SELECT DISTINCT Framework FROM Mapping WHERE Framework NOT IN (SELECT DISTINCT EvidenceLibraryMappedName FROM Framework_Lookup WHERE EvidenceLibraryMappedName IS NOT NULL) ORDER BY Framework;")
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

func GetFrameworkLookupFrameworks(db *sql.DB) ([]string, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	uniqueFrameworks := []string{}

	rows, err := db.Query("SELECT DISTINCT CEFramework FROM Framework_Lookup WHERE CEFramework IS NOT NULL ORDER BY CEFramework;")
	if err != nil {
		return nil, fmt.Errorf("error getting frameworks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var framework string
		if err := rows.Scan(&framework); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		uniqueFrameworks = append(uniqueFrameworks, framework)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return uniqueFrameworks, nil
}

func GetFrameworkInfoBackend(db *sql.DB, framework string) (map[string]interface{}, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	query := "SELECT EvidenceLibraryMappedName, AirtableFramework, AirtableView, FrameworkId_UAT, FrameworkId_Staging, FrameworkId_Staging FROM Framework_Lookup WHERE CEFramework = ?"
	row := db.QueryRow(query, framework)

	var details map[string]interface{}
	var evidenceLibraryName, airtableFramework, airtableView, frameworkId_UAT, frameworkId_Stage, frameworkId_Prod sql.NullString

	err := row.Scan(&evidenceLibraryName, &airtableFramework, &airtableView, &frameworkId_UAT, &frameworkId_Stage, &frameworkId_Prod)
	if err != nil {
		return nil, fmt.Errorf("error querying framework details: %w", err)
	}

	details = map[string]interface{}{
		"CEName":                    framework,
		"EvidenceLibraryMappedName": nullStringToString(evidenceLibraryName),
		"AirtableFramework":         nullStringToString(airtableFramework),
		"AirtableView":              nullStringToString(airtableView),
		"FrameworkId_UAT":           nullStringToString(frameworkId_UAT),
		"FrameworkId_Staging":       nullStringToString(frameworkId_Stage),
		"FrameworkId_Prod":          nullStringToString(frameworkId_Prod),
	}

	return details, nil
}

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func GetMappedFrameworkRecords(db *sql.DB) ([]string, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	query := "SELECT DISTINCT Framework FROM Mapping ORDER BY Framework;"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying mapped frameworks: %w", err)
	}
	defer rows.Close()
	var frameworks []string
	for rows.Next() {
		var framework string
		if err := rows.Scan(&framework); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		frameworks = append(frameworks, framework)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return frameworks, nil
}
