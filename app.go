package main

import (
	"cefp/airtable"
	"cefp/database"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

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
func (a *App) GetFrameworkLookup() ([]airtable.Framework, error) {
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

	//selectedFrameworkDetails, ok := data["selectedFrameworkDetails"].(map[string]interface{})
	//if !ok {
	//	return fmt.Errorf("invalid missing framework details")
	//}

	// Extract details
	cename, _ := data["cename"].(string)
	uatStage, _ := data["uatStage"].(string)
	prodNumber, _ := data["prodNumber"].(string)
	stageNumber, _ := data["stageNumber"].(string)
	tableID, _ := data["tableID"].(string)
	tableName, _ := data["tableName"].(string)
	tableView, _ := data["tableView"].(string)

	err := database.UpdateFrameworkLookupTable(a.db, missingFrameworkName, cename, uatStage, stageNumber, prodNumber, tableID, tableName, tableView)
	if err != nil {
		return fmt.Errorf("failed to update framework lookup: %v", err)
	}
	return nil
}

func (a *App) UpdateBuildFrameworkLookupTable(records []map[string]interface{}) error {
	if a.db == nil {
		return fmt.Errorf("database connection is nil")
	}

	for _, fields := range records {
		ceFramework, _ := fields["Name"].(string)
		frameworkIdUAT, _ := fields["UAT_Stage"].(string)
		frameworkIdStaging, _ := fields["Stage Framework Number"].(string)
		frameworkIdProd, _ := fields["Production Framework Number"].(string)

		if ceFramework == "" {
			continue // Skip records without a name
		}

		err := database.UpdateBuildFramework_LookupTable(a.db, ceFramework, frameworkIdUAT, frameworkIdStaging, frameworkIdProd)
		if err != nil {
			return fmt.Errorf("failed to update framework lookup: %v", err)
		}
	}
	return nil
}

func (a *App) GetAirtableBaseTables() (map[string]interface{}, error) {
	if a.apiKey == "" {
		log.Fatal("API Key is missing")
	}
	tables, err := airtable.GetAirtableTablesAndViews(a.apiKey)
	if err != nil {
		log.Printf("Error fetching Airtable tables: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal([]byte(tables), &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing tables JSON: %v", err)
	}

	tablesArray, ok := result["tables"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("error extracting tables array")
	}

	sort.SliceStable(tablesArray, func(i, j int) bool {
		tableI := tablesArray[i].(map[string]interface{})
		tableJ := tablesArray[j].(map[string]interface{})

		nameI, okI := tableI["name"].(string)
		nameJ, okJ := tableJ["name"].(string)

		if !okI || !okJ {
			return false
		}

		return nameI < nameJ
	})

	result["tables"] = tablesArray

	return result, err
}

func (a *App) GetMappedFrameworks() ([]string, error) {
	if a.db == nil {
		log.Fatal("database connection is missing")
	}

	uniqueFrameworks, err := database.GetMappedFrameworkRecords(a.db)
	if err != nil {
		log.Printf("Error fetching mapped frameworks: %v", err)
	}
	return uniqueFrameworks, err
}

func (a *App) GetUniqueFrameworks() ([]string, error) {
	if a.db == nil {
		log.Fatal("database connection is missing")
	}

	uniqueFrameworks, err := database.GetFrameworkLookupFrameworks(a.db)
	if err != nil {
		return nil, fmt.Errorf("error fetching unique frameworks: %v", err)
	}
	return uniqueFrameworks, err
}

func (a *App) GetFrameworkDetails(framework string) (map[string]interface{}, error) {
	if a.db == nil {
		log.Fatal("Database is missing")
	}
	//log.Printf("Getting framework details for %s", framework)
	frameworkInfo, err := database.GetFrameworkInfoBackend(a.db, framework)
	if err != nil {
		return nil, fmt.Errorf("error fetching framework details: %v", err)
	}

	return frameworkInfo, err
}

func (a *App) GetFrameworkRecords(data map[string]interface{}) error {
	tableName, _ := data["tableName"].(string)
	tableID, _ := data["tableID"].(string)
	tableView, _ := data["tableView"].(string)
	tableView = strings.ReplaceAll(tableView, " ", "%20")

	err := airtable.GetFrameworkData(a.db, a.apiKey, tableName, tableID, tableView)
	if err != nil {
		return fmt.Errorf("error fetching framework data: %v", err)
	}
	return err
}
