package main

import (
	"cefp/airtable"
	"cefp/database"
	"context"
	"database/sql"
	"fmt"
	"log"
)

// App struct
type App struct {
	ctx    context.Context
	db     *sql.DB
	apiKey string
}

// NewApp creates a new App application struct
func NewApp(apiKey string, db *sql.DB) *App {
	return &App{
		apiKey: apiKey,
		db:     db,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ReadAPIEvidenceTable() (string, error) {
	err := airtable.ReadAPI_EvidenceTable(a.ctx, a.db, a.apiKey)
	if err != nil {
		log.Printf("Error updating evidence table: %v", err)
		return "", fmt.Errorf("failed to read/update evidence table")
	}
	message := "Updated Evidence and Mapping tables"
	return message, nil
}

func (a *App) GetMissingFramework() ([]string, error) {
	records, err := database.GetMissingFrameworks(a.db)
	if err != nil {
		log.Printf("Error fetching missing frameworks: %v", err)
		return nil, fmt.Errorf("failed to retrieve missing frameworks")
	}
	return records, nil
}

// Expose GetFrameworkLookup to the frontend
func (a *App) GetFrameworkLookup() ([]airtable.AirtableFrameworks, error) {
	if a.apiKey == "" {
		log.Fatal("API Key is missing")
	}
	records, err := airtable.GetFrameworksLookup(a.apiKey)
	if err != nil {
		log.Printf("Error fetching frameworks lookup: %v", err)
		return nil, fmt.Errorf("failed to retrieve frameworks lookup")
	}
	return records, nil
}
