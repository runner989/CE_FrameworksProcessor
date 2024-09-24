package main

import (
	"cefp/airtable"
	"cefp/database"
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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
	runtime.WindowSetPosition(ctx, 1, 5) // Keeps window from opening outside the desktop bounds
}

// ReadAPIEvidenceTable get all records from the Evidence table in Airtable
func (a *App) ReadAPIEvidenceTable() (string, error) {
	err := airtable.ReadAPI_EvidenceTable(a.ctx, a.db, a.apiKey)
	if err != nil {
		log.Printf("Error updating evidence table: %v", err)
		return "", fmt.Errorf("failed to read/update evidence table")
	}
	message := "Updated Evidence and Mapping tables"
	return message, nil
}

// GetMissingFramework find frameworks in Mapping not in Framework_Lookup
func (a *App) GetMissingFramework() ([]string, error) {
	records, err := database.GetMissingFrameworks(a.db)
	if err != nil {
		log.Printf("Error fetching missing frameworks: %v", err)
		return nil, fmt.Errorf("failed to retrieve missing frameworks")
	}
	return records, nil
}

// GetFrameworkLookup Expose to the frontend
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

func (a *App) UpdateFrameworkLookup(data map[string]interface{}) error {
	missingFrameworkName, ok := data["missingFrameworkName"].(string)
	if !ok {
		return fmt.Errorf("invalid missing framework name")
	}

	selectedFrameworkDetails, ok := data["selectedFrameworkDetails"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid missing framework details")
	}

	// Extract details
	name, _ := selectedFrameworkDetails["name"].(string)
	uatStage, _ := selectedFrameworkDetails["uatStage"].(string)
	prodNumber, _ := selectedFrameworkDetails["prodNumber"].(string)
	stageNumber, _ := selectedFrameworkDetails["stageNumber"].(string)

	err := database.UpdateFrameworkLookupTable(a.db, missingFrameworkName, name, uatStage, stageNumber, prodNumber)
	if err != nil {
		return fmt.Errorf("failed to update framework lookup: %v", err)
	}
	return nil
}

func (a *App) UpdateBuildFrameworkLookupTable(records []map[string]interface{}) error {
	if a.db == nil {
		return fmt.Errorf("database connection is nil")
	}

	// tx, err := a.db.Begin()
	// if err != nil {
	// 	return fmt.Errorf("failed to begin transaction: %v", err)
	// }

	// defer func() {
	// 	if err != nil {
	// 		_ = tx.Rollback()
	// 	} else {
	// 		_ = tx.Commit()
	// 	}
	// }()

	// query := `
	// 	INSERT OR REPLACE INTO Framework_Lookup (CEFramework, FrameworkId_UAT, FrameworkId_Staging, FrameworkId_Prod)
	// 	VALUES (?, ?, ?, ?);
	// `

	// Let's plan to move this to the database package and out of app.go
	query := `
		INSERT INTO Framework_Lookup (CEFramework, FrameworkId_UAT, FrameworkId_Staging, FrameworkId_Prod)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(CEFramework) DO UPDATE SET
		  FrameworkId_UAT = excluded.FrameworkId_UAT,
		  FrameworkId_Staging = excluded.FrameworkId_Staging,
		  FrameworkId_Prod = excluded.FrameworkId_Prod;
		`
	// log.Println("running upsert for Framework_Lookup table")

	stmt, err := a.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Printf("error closing statement: %v", err)
		}
	}(stmt)

	for _, fields := range records {
		ceFramework, _ := fields["Name"].(string)
		frameworkIdUAT, _ := fields["UAT_Stage"].(string)
		frameworkIdStaging, _ := fields["Stage Framework Number"].(string)
		frameworkIdProd, _ := fields["Production Framework Number"].(string)

		if ceFramework == "" {
			continue // Skip records without a name
		}

		// log.Printf("Insert or replace %s", ceFramework)

		_, err = stmt.Exec(ceFramework, frameworkIdUAT, frameworkIdStaging, frameworkIdProd)
		if err != nil {
			return fmt.Errorf("failed to update framework lookup table: %v", err)
		}
	}
	return nil
}
