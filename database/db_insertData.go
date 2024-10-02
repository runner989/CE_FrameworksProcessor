package database

import (
	"cefp/structs"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"log"
	"strconv"

	"github.com/mattn/go-sqlite3"
	"github.com/xuri/excelize/v2"
)

func UpdateFrameworkLookupTable(db *sql.DB, lookupRecord structs.FrameworkLookup) error { //missingFrameworkName, cename, uatStage, stageNumber, prodNumber, tableID, tableName, viewName string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	//log.Printf("missing framework %s, cename: %s, uatStage: %s, stageNumber: %s, prodNumber: %s, tableName: %s, viewName: %s", missingFrameworkName, cename, uatStage, stageNumber, prodNumber, tableName, viewName)
	query := `
		INSERT INTO Framework_Lookup (EvidenceLibraryMappedName, CEFramework, FrameworkId_UAT, FrameworkId_Staging, FrameworkId_Prod, AirtableTableID ,AirtableFramework, AirtableView, AirtableBase)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(CEFramework) DO UPDATE SET
			EvidenceLibraryMappedName = excluded.EvidenceLibraryMappedName,
			FrameworkId_UAT = excluded.FrameworkId_UAT,
			FrameworkId_Staging = excluded.FrameworkId_Staging,
			FrameworkId_Prod = excluded.FrameworkId_Prod,
			AirtableTableID = excluded.AirtableTableID,
			AirtableFramework = excluded.AirtableFramework,
			AirtableView = excluded.AirtableView,
			AirtableBase = excluded.AirtableBase;
		`

	statement, err := db.Prepare(query)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(lookupRecord.MappedName, lookupRecord.CeName, lookupRecord.UatStage, lookupRecord.StageNumber, lookupRecord.ProdNumber, lookupRecord.TableID, lookupRecord.TableName, lookupRecord.TableView, lookupRecord.TableBase)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return err
}

func UpdateBuildFramework_LookupTable(db *sql.DB, lr structs.FrameworkLookup) error { // CEFramework, FrameworkidUat, FrameworkidStaging, FrameworkidProd string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	query := `
		INSERT INTO Framework_Lookup (CEFramework, FrameworkId_UAT, FrameworkId_Staging, FrameworkId_Prod)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(CEFramework) DO UPDATE SET
		  FrameworkId_UAT = excluded.FrameworkId_UAT,
		  FrameworkId_Staging = excluded.FrameworkId_Staging,
		  FrameworkId_Prod = excluded.FrameworkId_Prod;
		`

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Printf("error closing statement: %v", err)
		}
	}(stmt)

	_, err = stmt.Exec(lr.CeName, lr.UatStage, lr.StageNumber, lr.ProdNumber) //FrameworkidUat, FrameworkidStaging, FrameworkidProd)
	if err != nil {

		return fmt.Errorf("failed to update framework lookup table: %v", err)
	}
	return nil
}

func InsertFrameworkRecord(db *sql.DB, fr structs.FrameworkRecord) error {
	_, err := db.Exec("INSERT INTO Framework (Framework, sortID, Identifier, ParentIdentifier, DisplayName, Description, Guidance, TestType, Tags, PolicyAndProcedureAIPromptTemplateId, ControlNarrativeAIPromptTemplateId) VALUES (?, ?, ?, ?, ?, ?, ?, ?,?,?,?)",
		fr.FrameworkName, fr.SortID, fr.Identifier, fr.ParentID, fr.DisplayName, fr.Description, fr.Guidance, fr.TestType, fr.Tags, fr.PromptID, fr.ControlNarrative)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			// Handle UNIQUE constraint violation (e.g. duplicate record)
			return fmt.Errorf("duplicate record: %s", fr.Identifier)
		}
		return fmt.Errorf("error inserting Framework: %v", err)
	}
	return nil

}

func ReadExcelAndSaveToDB(ctx context.Context, db *sql.DB, file io.Reader, table string) error {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return fmt.Errorf("error opening Excel file: %v", err)
	}
	defer f.Close()

	runtime.EventsEmit(ctx, "mappingprogress", fmt.Sprintf("Opening CE Mapping %s file: %v", table, err))

	qry := fmt.Sprintf("DELETE FROM CEMapping_%s", table)

	_, err = db.Exec(qry)
	if err != nil {
		runtime.EventsEmit(ctx, "mappingprogress", fmt.Sprintf("Error deleting from CEMapping_%s: %v", table, err))
		return fmt.Errorf("error deleting from CEMapping_%s: %v", table, err)
	}

	rows, err := f.GetRows("Mapping")
	if err != nil {
		return fmt.Errorf("error reading Mapping sheet: %v", err)
	}

	for i, row := range rows {
		if i == 0 && row[0] == "EvidenceID" {
			continue
		}

		evidenceMapRecord := structs.EvidenceMapRecord{
			EvidenceID:      getIntOrZero(getStringOrEmpty(row, 0)),
			Framework:       getStringOrEmpty(row, 1),
			FrameworkID:     getIntOrZero(getStringOrEmpty(row, 2)),
			Requirement:     getStringOrEmpty(row, 3),
			Description:     getStringOrEmpty(row, 4),
			Guidance:        getStringOrEmpty(row, 5),
			RequirementType: getStringOrEmpty(row, 6),
			Delete:          getStringOrEmpty(row, 7),
		}

		message := fmt.Sprintf("Processing EvidenceID: %d, Evidence: %s", evidenceMapRecord.EvidenceID, evidenceMapRecord.Framework)
		runtime.EventsEmit(ctx, "mappingprogress", message)

		err = saveEvidenceRecordToDB(db, evidenceMapRecord, table)
		if err != nil {
			log.Printf("error saving evidence record: %v", err)
		}
	}

	message := fmt.Sprintf("Done updating CEMapping_%s table!", table)
	runtime.EventsEmit(ctx, "mappingprogress", message)
	return nil
}

// Function to save the record to the database
func saveEvidenceRecordToDB(db *sql.DB, record structs.EvidenceMapRecord, table string) error {

	query := fmt.Sprintf("INSERT INTO CEMapping_%s ("+
		"EvidenceID, Framework, FrameworkID, Requirement, Description, Guidance, RequirementType, \"Delete\" "+
		") VALUES (?, ?, ?, ?, ?, ?, ?, ?)", table)

	_, err := db.Exec(query,
		record.EvidenceID, record.Framework, record.FrameworkID,
		record.Requirement, record.Description, record.Guidance,
		record.RequirementType, record.Delete)
	return err
}

// Helper function to get int from string or return 0 if conversion fails
func getIntOrZero(value string) int {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0 // or handle it as needed
	}
	return intValue
}
func getStringOrEmpty(row []string, idx int) string {
	if len(row) > idx {
		return row[idx]
	}
	return ""
}
