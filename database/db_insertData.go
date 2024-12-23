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
	"path/filepath"
	"strconv"

	"github.com/mattn/go-sqlite3"
	"github.com/xuri/excelize/v2"
)

func BackupMemoryToFile(memDB, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Fetch data from in-memory Evidence table
	rows, err := memDB.Query("SELECT * FROM Evidence")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to query Evidence from in-memory DB: %v", err)
	}
	defer rows.Close()

	// Insert data into the file-based Evidence table
	stmt, err := tx.Prepare("INSERT INTO Evidence (EvidenceID, Evidence, Description, AnecdotesEvidenceIds, Priority, EvidenceType) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare insert statement for Evidence: %v", err)
	}
	defer stmt.Close()

	for rows.Next() {
		var EvidenceID int
		var Evidence string
		var Description string
		var AnecdotesEvidenceIds sql.NullString
		var Priority sql.NullString
		var EvidenceType sql.NullString

		err = rows.Scan(&EvidenceID, &Evidence, &Description, &AnecdotesEvidenceIds, &Priority, &EvidenceType)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to scan row from in-memory Evidence: %v", err)
		}

		_, err = stmt.Exec(EvidenceID, Evidence, Description, AnecdotesEvidenceIds, Priority, EvidenceType)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert row into file-based Evidence: %v", err)
		}
	}

	// Fetch data from in-memory Mapping table
	rows, err = memDB.Query("SELECT EvidenceID, Framework, FrameworkId, Requirement, Description, Guidance, RequirementType FROM Mapping")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to query Mapping from in-memory DB: %v", err)
	}
	defer rows.Close()

	// Insert data into the file-based Mapping table
	stmt, err = tx.Prepare("INSERT INTO Mapping (EvidenceID, Framework, FrameworkId, Requirement, Description, Guidance, RequirementType) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare insert statement for Mapping: %v", err)
	}
	defer stmt.Close()

	for rows.Next() {
		var EvidenceID int
		var Framework string
		var FrameworkId sql.NullInt32
		var Requirement sql.NullString
		var Description sql.NullString
		var Guidance sql.NullString
		var RequirementType sql.NullString

		err = rows.Scan(&EvidenceID, &Framework, &FrameworkId, &Requirement, &Description, &Guidance, &RequirementType)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to scan row from in-memory Mapping: %v", err)
		}

		_, err = stmt.Exec(EvidenceID, Framework, FrameworkId, Requirement, Description, Guidance, RequirementType)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert row into file-based Mapping: %v", err)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func UpdateFrameworkLookupTable(db *sql.DB, lookupRecord structs.FrameworkLookup) error { //missingFrameworkName, cename, uatStage, stageNumber, prodNumber, tableID, tableName, viewName string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

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
			fmt.Println("duplicate record.")
			return fmt.Errorf("duplicate record: %s", fr.Identifier)
		}
		fmt.Println("error inserting Framework.")
		return fmt.Errorf("error inserting Framework: %v", err)
	}
	return nil
}

func AddPlaceholders(db *sql.DB, id int) error {
	evidenceTitle := "Place holder"
	evidenceDesc := "Deleted mappings..."
	query := fmt.Sprintf("INSERT INTO Evidence (EvidenceID, Evidence, Description) VALUES (?, ?,?)")
	_, err := db.Exec(query, id, evidenceTitle, evidenceDesc)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	return nil
}

//func SafeString(ns sql.NullString) string {
//	if ns.Valid {
//		return ns.String
//	}
//	return ""
//}

func insertSafeString(value interface{}) sql.NullString {
	if value == nil {
		return sql.NullString{String: "", Valid: false}
	}

	str, ok := value.(string)
	if !ok {
		return sql.NullString{String: "", Valid: false}
	}

	return sql.NullString{String: str, Valid: true}
}

func ReadExcelAndSaveToDB(ctx context.Context, memDB, db *sql.DB, file io.Reader, filePath, table string) error {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return fmt.Errorf("error opening Excel file: %v", err)
	}
	defer f.Close()

	fileName := filepath.Base(filePath)
	runtime.EventsEmit(ctx, "mappingprogress", fmt.Sprintf("Opening CE Mapping %s file: %s", table, fileName))

	qry := fmt.Sprintf("DELETE FROM CEMapping_%s", table)

	_, err = db.Exec(qry)
	_, err = memDB.Exec(qry)
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
			EvidenceID:      getIntOrZero(row[0]),
			Framework:       getStringOrEmpty(row, 1),
			FrameworkID:     getIntOrZero(getStringOrEmpty(row, 2)),
			Requirement:     insertSafeString(row[3]),
			Description:     insertSafeString(row[4]),
			Guidance:        insertSafeString(row[5]),
			RequirementType: insertSafeString(row[6]),
			//Delete:          insertSafeString(row[7]),
		}

		message := fmt.Sprintf("Processing EvidenceID: %d, Evidence: %s", evidenceMapRecord.EvidenceID, evidenceMapRecord.Framework)
		runtime.EventsEmit(ctx, "mappingprogress", message)

		err = saveMappingRecordToMemDB(memDB, evidenceMapRecord, table)
		if err != nil {
			log.Printf("error saving evidence record: %v", err)
		}
	}
	err = saveMappingRecordsToDB(memDB, db, table)

	message := fmt.Sprintf("Done updating CEMapping_%s table!", table)
	runtime.EventsEmit(ctx, "mappingprogress", message)
	return nil
}

func MoveFrameworkMemDBToFile(db, memDB *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	if memDB == nil {
		return fmt.Errorf("memory database connection is nil")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	rows, err := memDB.Query("SELECT SORTID, FRAMEWORK, IDENTIFIER, PARENTIDENTIFIER, DISPLAYNAME, DESCRIPTION, GUIDANCE, TESTTYPE, TAGS, POLICYANDPROCEDUREAIPROMPTTEMPLATEID, CONTROLNARRATIVEAIPROMPTTEMPLATEID FROM Framework")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error reading Framework: %v", err)
	}
	defer rows.Close()

	insertQuery := `INSERT INTO Framework (sortID, Framework, Identifier, ParentIdentifier, DisplayName, Description, Guidance, TestType, Tags, PolicyAndProcedureAIPromptTemplateId, ControlNarrativeAIPromptTemplateId) 
    	VALUES (?,?,?,?,?,?,?,?,?,?,?)
    `
	stmt, err := tx.Prepare(insertQuery)
	if err != nil {
		tx.Rollback()
		log.Printf("failed to prepare insert statement for Framework: %v", err)
		return fmt.Errorf("failed to prepare insert statement for Framework: %v", err)
	}
	defer stmt.Close()

	for rows.Next() {
		var sortID sql.NullString
		var frameworkName sql.NullString
		var identifier sql.NullString
		var parentID sql.NullString
		var displayName sql.NullString
		var description sql.NullString
		var guidance sql.NullString
		var tags sql.NullString
		var testType sql.NullString
		var promptID sql.NullInt32
		var controlNarrative sql.NullInt32

		err = rows.Scan(&sortID, &frameworkName, &identifier, &parentID, &displayName, &description, &guidance, &tags, &testType, &promptID, &controlNarrative)
		if err != nil {
			tx.Rollback()
			log.Printf("error scanning row: %v", err)
			return fmt.Errorf("error scanning row: %v", err)
		}
		_, err = stmt.Exec(sortID, frameworkName, identifier, parentID, displayName, description, guidance, testType, tags, promptID, controlNarrative)
		if err != nil {
			tx.Rollback()
			log.Printf("error executing insert statement for Framework: %v", err)
			return fmt.Errorf("failed to insert row into file-based Framework: %v", err)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("failed to commit transaction: %v", err)
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func saveMappingRecordsToDB(memDB, db *sql.DB, table string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	if memDB == nil {
		return fmt.Errorf("memory database connection is nil")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Fetch data from in-memory Mapping table
	query := fmt.Sprintf("SELECT EvidenceID, Framework, FrameworkId, Requirement, Description, Guidance, RequirementType FROM CEMapping_%s", table)
	log.Printf("memDB select: %s", query)
	rows, err := memDB.Query(query)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to query Mapping from in-memory DB: %v", err)
	}
	defer rows.Close()

	// Insert data into the file-based Mapping table
	insertQuery := fmt.Sprintf("INSERT INTO CEMapping_%s (EvidenceID, Framework, FrameworkId, Requirement, Description, Guidance, RequirementType) VALUES (?, ?, ?, ?, ?, ?, ?)", table)

	stmt, err := tx.Prepare(insertQuery)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare insert statement for Mapping: %v", err)
	}
	defer stmt.Close()

	for rows.Next() {
		var EvidenceID int
		var Framework string
		var FrameworkId sql.NullInt32
		var Requirement sql.NullString
		var Description sql.NullString
		var Guidance sql.NullString
		var RequirementType sql.NullString

		err = rows.Scan(&EvidenceID, &Framework, &FrameworkId, &Requirement, &Description, &Guidance, &RequirementType)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to scan row from in-memory Mapping: %v", err)
		}

		_, err = stmt.Exec(EvidenceID, Framework, FrameworkId, Requirement, Description, Guidance, RequirementType)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert row into file-based Mapping: %v", err)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// Function to save the record to the database
func saveMappingRecordToMemDB(memDB *sql.DB, record structs.EvidenceMapRecord, table string) error {

	query := fmt.Sprintf("INSERT INTO CEMapping_%s "+
		"(EvidenceID, Framework, FrameworkID, Requirement, Description, Guidance, RequirementType"+
		") VALUES (?, ?, ?, ?, ?, ?, ?)", table)

	_, err := memDB.Exec(query,
		record.EvidenceID, record.Framework, record.FrameworkID,
		record.Requirement, record.Description, record.Guidance,
		record.RequirementType)
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
