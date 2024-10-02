package main

import (
	"cefp/airtable"
	"cefp/database"
	"cefp/structs"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
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
	err := airtable.ReadAPIEvidenceTable(a.ctx, a.db, a.apiKey)
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
func (a *App) GetFrameworkLookup() ([]structs.Framework, error) {
	if a.apiKey == "" {
		log.Fatal("API Key is missing")
	}
	records, err := airtable.GetFrameworksLookup(a.apiKey)
	if err != nil {
		log.Printf("Error fetching frameworks lookup: %v", err)
		return nil, fmt.Errorf("failed to retrieve frameworks lookup")
	}

	sort.SliceStable(records, func(i, j int) bool {
		nameI, okI := records[i].Fields["Name"].(string)
		nameJ, okJ := records[j].Fields["Name"].(string)

		if !okI && !okJ {
			return false
		}
		if !okI {
			return false
		}
		if !okJ {
			return true
		}
		return nameI < nameJ
	})

	return records, nil
}

// Helper function to convert an interface{} to float64
func toFloat64(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			//fmt.Printf("Error converting string to float64: %v\n", err)
			return 0 // Fallback value if conversion fails
		}
		return f
	default:
		//fmt.Printf("Unexpected type: %T\n", v)
		return 0 // Fallback value for unexpected types
	}
}

func (a *App) UpdateFrameworkLookup(data map[string]interface{}) error {
	missingFrameworkName, ok := data["missingFrameworkName"].(string)
	if !ok {
		return fmt.Errorf("invalid missing framework name")
	}

	lookupRecord := structs.FrameworkLookup{
		MappedName:  sql.NullString{String: missingFrameworkName, Valid: true},
		CeName:      sql.NullString{String: data["ceName"].(string), Valid: true},
		UatStage:    sql.NullFloat64{Float64: toFloat64(data["uatStage"]), Valid: true},
		ProdNumber:  sql.NullFloat64{Float64: toFloat64(data["prodNumber"]), Valid: true},
		StageNumber: sql.NullFloat64{Float64: toFloat64(data["stageNumber"]), Valid: true},
		TableBase:   sql.NullString{String: data["baseID"].(string), Valid: true},
		TableID:     sql.NullString{String: data["tableID"].(string), Valid: true},
		TableName:   sql.NullString{String: data["tableName"].(string), Valid: true},
		TableView:   sql.NullString{String: data["tableView"].(string), Valid: true},
	}

	err := database.UpdateFrameworkLookupTable(a.db, lookupRecord) //missingFrameworkName, cename, uatStage, stageNumber, prodNumber, tableID, tableName, tableView)
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
		lookupRecord := structs.FrameworkLookup{
			CeName:      SafeString(fields["Name"]),
			UatStage:    safeFloat(toFloat64(fields["UAT_Stage"])),
			StageNumber: safeFloat(toFloat64(fields["Stage Framework Number"])),
			ProdNumber:  safeFloat(toFloat64(fields["Production Framework Number"])),
		}

		if !lookupRecord.CeName.Valid || lookupRecord.CeName.String == "" {
			continue // Skip records without a name
		}
		err := database.UpdateBuildFramework_LookupTable(a.db, lookupRecord)
		if err != nil {
			return fmt.Errorf("failed to update framework lookup: %v", err)
		}
	}
	return nil
}

func SafeString(value interface{}) sql.NullString {
	if value == nil {
		return sql.NullString{String: "", Valid: false}
	}

	str, ok := value.(string)
	if !ok {
		return sql.NullString{String: "", Valid: false}
	}

	return sql.NullString{String: str, Valid: true}
}

func safeFloat(value interface{}) sql.NullFloat64 {
	if value == nil {
		return sql.NullFloat64{Float64: 0, Valid: false}
	}
	valFloat, ok := value.(float64)
	if !ok {
		return sql.NullFloat64{Float64: valFloat, Valid: true}
	}
	return sql.NullFloat64{Float64: valFloat, Valid: true}
}

func (a *App) GetAvailableAirtableBases() ([]structs.Base, error) {
	if a.apiKey == "" {
		log.Fatal("API Key is missing")
	}
	bases, err := airtable.GetAirtableBases(a.apiKey)
	if err != nil {
		log.Printf("Error fetching available airtable bases: %v", err)
	}
	return bases, err
}

func (a *App) UpdateAirtableBasesTable(records []map[string]interface{}) error {
	_, err := a.db.Exec("DELETE FROM Airtable_Base")
	if err != nil {
		return fmt.Errorf("failed to clear airtable bases table: %v", err)
	}

	query := `INSERT INTO Airtable_Base (BaseID, BaseName) VALUES (?, ?)`

	for _, record := range records {
		_, err = a.db.Exec(
			query,
			record["id"].(string),
			record["name"].(string),
		)
		if err != nil {
			log.Printf("Error updating airtable bases: %v", err)
			return fmt.Errorf("failed to update airtable bases: %v", err)
		}
	}

	return nil
}

func (a *App) GetAirtableBaseTables() (map[string]interface{}, error) {
	if a.apiKey == "" {
		log.Fatal("API Key is missing")
	}
	baseID := "app5fTueYfRM65SzX"
	tables, err := airtable.GetAirtableTablesAndViews(a.apiKey, baseID)
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

	frameworkInfo, err := database.GetFrameworkInfoBackend(a.db, framework)
	if err != nil {
		return nil, fmt.Errorf("error fetching framework details: %v", err)
	}

	return frameworkInfo, err
}

func (a *App) GetFrameworkRecords(data map[string]interface{}) error {
	tableView, _ := data["tableView"].(string)
	tableView = strings.ReplaceAll(tableView, " ", "%20")

	lookupRecord := structs.FrameworkLookup{
		CeName:    sql.NullString{String: data["ceName"].(string), Valid: true},
		TableName: sql.NullString{String: data["tableName"].(string), Valid: true},
		TableID:   sql.NullString{String: data["tableID"].(string), Valid: true},
		TableView: sql.NullString{String: tableView, Valid: true},
	}
	err := airtable.GetFrameworkData(a.db, a.apiKey, lookupRecord)
	if err != nil {
		return fmt.Errorf("error fetching framework data: %v", err)
	}
	return err
}

func (a *App) GetFrameworkLookupTable() ([]map[string]interface{}, error) {
	query := "SELECT ROWID, * FROM Framework_Lookup ORDER BY CEFramework"
	rows, err := a.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying framework lookup table: %v", err)
	}

	defer rows.Close()

	var records []map[string]interface{}
	for rows.Next() {
		// Use sql.NullString and sql.NullInt64 for nullable columns
		var airtableBase, airtableTableID, airtableFramework, airtableView, evidenceLibraryMappedName, ceFramework, description, comments sql.NullString
		var frameworkidUat, frameworkidStaging, frameworkidProd, version sql.NullString
		var rowID int64

		// Scan the values into the appropriate variables
		err := rows.Scan(
			&rowID,
			&airtableBase,
			&airtableTableID,
			&airtableFramework,
			&airtableView,
			&evidenceLibraryMappedName,
			&ceFramework,
			&frameworkidUat,
			&frameworkidStaging,
			&frameworkidProd,
			&version,
			&description,
			&comments,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Convert sql.NullString and sql.NullInt64 to normal types or handle NULLs
		record := map[string]interface{}{
			"RowID":                     rowID,
			"AirtableBase":              nullStringToString(airtableBase),
			"AirtableTableID":           nullStringToString(airtableTableID),
			"AirtableFramework":         nullStringToString(airtableFramework),
			"AirtableView":              nullStringToString(airtableView),
			"EvidenceLibraryMappedName": nullStringToString(evidenceLibraryMappedName),
			"CEFramework":               nullStringToString(ceFramework),
			"FrameworkId_UAT":           stringToInt(frameworkidUat.String),
			"FrameworkId_Staging":       stringToInt(frameworkidStaging.String),
			"FrameworkId_Prod":          stringToInt(frameworkidProd.String),
			"Version":                   stringToInt(version.String),
			"Description":               nullStringToString(description),
			"Comments":                  nullStringToString(comments),
		}

		records = append(records, record)
	}
	return records, nil
}

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return "" // Return empty string if NULL
}

// Helper function to convert string to int safely
func stringToInt(s string) int {
	if s == "" {
		return 0 // Return 0 if the string is empty
	}

	// Convert the string to an integer
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0 // Handle conversion error (return 0 if conversion fails)
	}
	return i
}

func (a *App) UpdateFrameworkLookupRecord(updatedRecord map[string]interface{}) error {

	query := `
        UPDATE Framework_Lookup SET
            AirtableBase = ?,
            AirtableTableID = ?,
            AirtableFramework = ?,
            AirtableView = ?,
            EvidenceLibraryMappedName = ?,
            CEFramework = ?,
            FrameworkId_UAT = ?,
            FrameworkId_Staging = ?,
            FrameworkId_Prod = ?,
            Version = ?,
            Description = ?,
            Comments = ?
        WHERE ROWID = ?
    `
	_, err := a.db.Exec(
		query,
		updatedRecord["airtableBase"],
		updatedRecord["airtableTableID"],
		updatedRecord["airtableFramework"],
		updatedRecord["airtableView"],
		updatedRecord["evidenceLibraryMappedName"],
		updatedRecord["ceFramework"],
		updatedRecord["frameworkId_UAT"],
		updatedRecord["frameworkId_Staging"],
		updatedRecord["frameworkId_Prod"],
		updatedRecord["version"],
		updatedRecord["description"],
		updatedRecord["comments"],
		updatedRecord["rowID"],
	)
	if err != nil {
		log.Printf("Error updating framework lookup record: %v", err)
		return fmt.Errorf("error updating framework record: %v", err)
	}

	return nil
}

func (a *App) DeleteSelectedFramework(selectedRecord map[string]interface{}) error {
	framework := structs.FrameworkLookup{
		RowID:       safeFloat(selectedRecord["rowID"]),
		MappedName:  SafeString(selectedRecord["mappedName"]),
		CeName:      SafeString(selectedRecord["ceFramework"]),
		UatStage:    safeFloat(selectedRecord["frameworkId_UAT"]),
		StageNumber: safeFloat(selectedRecord["frameworkId_Staging"]),
		ProdNumber:  safeFloat(selectedRecord["frameworkId_Prod"]),
		TableBase:   SafeString(selectedRecord["airtableBase"]),
		TableID:     SafeString(selectedRecord["airtableTableID"]),
		TableName:   SafeString(selectedRecord["airtableFramework"]),
		TableView:   SafeString(selectedRecord["airtableView"]),
		Version:     SafeString(selectedRecord["version"]),
		Description: SafeString(selectedRecord["description"]),
		Comments:    SafeString(selectedRecord["comments"]),
	}

	// Delete framework from Framework_Lookup
	err := database.DeleteFromFrameworkLookup(a.db, framework)
	if err != nil {
		log.Printf("Error deleting selected framework lookup record: %v", err)
		return fmt.Errorf("error deleting selected framework lookup record: %v", err)
	}

	// Delete framework from Framework
	err = database.DeleteFromFramework(a.db, framework)
	if err != nil {
		log.Printf("Error deleting selected framework record: %v", err)
	}
	return nil

}

func (a *App) ProcessEvidenceStagingFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	err = database.ReadExcelAndSaveToDB(a.ctx, a.db, file, "Staging")
	if err != nil {
		return fmt.Errorf("error processing evidence staging file: %v", err)
	}

	ids, err := database.CheckForMissing(a.db, "Staging")
	if err != nil {
		return fmt.Errorf("error checking for missing evidence IDs in staging: %v", err)
	}
	if len(ids) > 0 {
		for id := range ids {
			err = database.AddPlaceholders(a.db, id)
			if err != nil {
				return fmt.Errorf("error adding placeholders: %v", err)
			}
		}
	}

	return nil
}

func (a *App) ProcessEvidenceProdFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	err = database.ReadExcelAndSaveToDB(a.ctx, a.db, file, "Prod")
	if err != nil {
		return fmt.Errorf("error processing evidence prod file: %v", err)
	}

	ids, err := database.CheckForMissing(a.db, "Prod")
	if err != nil {
		return fmt.Errorf("error checking for missing evidence IDs in prod: %v", err)
	}
	if len(ids) > 0 {
		for id := range ids {
			err = database.AddPlaceholders(a.db, id)
			if err != nil {
				return fmt.Errorf("error adding placeholders: %v", err)
			}
		}
	}

	return nil
}

func (a *App) OpenFileDialog() (string, error) {
	// Open the file dialog and allow the user to select a file
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Excel File",
		Filters: []runtime.FileFilter{
			{DisplayName: "Excel Files", Pattern: "*.xlsx"},
		},
	})

	if err != nil {
		return "", fmt.Errorf("error opening file dialog: %v", err)
	}

	return filePath, nil
}

func (a *App) GetAirtableTables(baseID string) (map[string]interface{}, error) {
	if a.apiKey == "" {
		log.Fatal("API Key is missing")
	}

	tables, err := airtable.GetAirtableTablesAndViews(a.apiKey, baseID)
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

func (a *App) UpdateAllFrameworks() (string, error) {
	frameworks, err := database.GetReadyFrameworks(a.db)
	if err != nil {
		return "", fmt.Errorf("error fetching ready frameworks: %v", err)
	}

	for _, framework := range frameworks {
		runtime.EventsEmit(a.ctx, "progress", fmt.Sprintf("Updating framework: %v", framework.CeName.String))
		err := airtable.GetFrameworkData(a.db, a.apiKey, framework)
		if err != nil {
			return "", fmt.Errorf("error getting framework from airtable: %v", err)
		}
	}

	return "All frameworks updated successfully!", nil
}

func (a *App) GetAllFrameworks() ([]string, error) {
	frameworks, err := database.GetDistinctFrameworks(a.db)
	if err != nil {
		return nil, fmt.Errorf("error fetching ready frameworks: %v", err)
	}
	return frameworks, nil
}

func (a *App) ExportAFramework(framework string) error {
	err := database.ExportFrameworkToExcel(a.db, framework)
	if err != nil {
		return fmt.Errorf("error exporting framework: %v", err)
	}
	return nil
}

func (a *App) ExportAllFrameworks() error {
	frameworks, err := database.GetDistinctFrameworks(a.db)
	if err != nil {
		return fmt.Errorf("error fetching ready frameworks: %v", err)
	}
	for _, framework := range frameworks {
		err := database.ExportFrameworkToExcel(a.db, framework)
		if err != nil {
			return fmt.Errorf("error exporting framework: %v", err)
		}
	}
	return nil
}

func (a *App) ExportEvidenceMapReport(table string) error {
	err := database.ExportEvidenceMapReportToExcel(a.db, table)
	if err != nil {
		return fmt.Errorf("error exporting evidence map report: %v", err)
	}
	return nil
}

//func (a *App) UpdateFrameworkName(data map[string]string) error {
//	oldFramework := data["oldFramework"]
//	newFramework := data["newFramework"]
//
//	query := `UPDATE Framework SET Framework = ? WHERE Framework = ?`
//	_, err := a.db.Exec(query, newFramework, oldFramework)
//	if err != nil {
//		return fmt.Errorf("error updating framework name: %v", err)
//	}
//	return nil
//}
