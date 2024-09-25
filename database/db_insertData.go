package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func UpdateFrameworkLookupTable(db *sql.DB, missingFrameworkName, cename, uatStage, stageNumber, prodNumber, tableID, tableName, viewName string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	//log.Printf("missing framework %s, cename: %s, uatStage: %s, stageNumber: %s, prodNumber: %s, tableName: %s, viewName: %s", missingFrameworkName, cename, uatStage, stageNumber, prodNumber, tableName, viewName)
	query := `INSERT INTO Framework_Lookup (EvidenceLibraryMappedName, CEFramework, FrameworkId_UAT, FrameworkId_Staging, FrameworkId_Prod, AirtableTableID ,AirtableFramework, AirtableView)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(CEFramework) DO UPDATE SET
		  EvidenceLibraryMappedName = excluded.EvidenceLibraryMappedName,
		  FrameworkId_UAT = excluded.FrameworkId_UAT,
		  FrameworkId_Staging = excluded.FrameworkId_Staging,
		  FrameworkId_Prod = excluded.FrameworkId_Prod,
		  AirtableTableID = excluded.AirtableTableID,
		  AirtableFramework = excluded.AirtableFramework,
		  AirtableView = excluded.AirtableView;
		`
	//UPDATE Framework_Lookup SET CEFramework = ?, FrameworkId_UAT = ?, FrameworkId_Staging = ?, FrameworkId_Prod = ? WHERE EvidenceLibraryMappedName = ?`

	statement, err := db.Prepare(query)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(missingFrameworkName, cename, uatStage, stageNumber, prodNumber, tableID, tableName, viewName)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return err
}

func UpdateBuildFramework_LookupTable(db *sql.DB, CEFramework, FrameworkidUat, FrameworkidStaging, FrameworkidProd string) error {
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

	_, err = stmt.Exec(CEFramework, FrameworkidUat, FrameworkidStaging, FrameworkidProd)
	if err != nil {

		return fmt.Errorf("failed to update framework lookup table: %v", err)
	}
	return nil
}

func InsertFrameworkRecord(db *sql.DB, sortID, policyID, controlNarrative int, frameworkName, identifier, parentID, displayName, description, guidance, tags, testType string) error {
	_, err := db.Exec("INSERT INTO Framework (Framework, sortID, Identifier, ParentIdentifier, DisplayName, Description, Guidance, TestType, Tags, PolicyAndProcedureAIPromptTemplateId, ControlNarrativeAllPromptTemplateId) VALUES (?, ?, ?, ?, ?, ?, ?, ?,?,?,?)",
		frameworkName, sortID, identifier, parentID, displayName, description, guidance, testType, tags, policyID, controlNarrative)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			// Handle UNIQUE constraint violation (duplicate EvidenceID)
			return fmt.Errorf("duplicate record: %s", identifier)
		}
		return fmt.Errorf("error inserting Framework: %v", err)
	}
	return nil

}
