package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InsertFrameworkRecord(
	db *sql.DB,
	id int,
	identifier string,
	parentident string,
	description string,
	displayName string,
	guidance string,
	recommendations string,
	requirementType string,
	pandpPromptId int,
	controlNarrativeId int,
	testType string,
	framework string,
	frameworkId int,
) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	var existingID int
	query := `SELECT sortID FROM framework WHERE sortID = ?;`
	err := db.QueryRow(query, id).Scan(&existingID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error checking for existing sortID: %v", err)
	} else if err == nil {
		log.Printf("Record with sortID %d already exists, skipping insertion", id)
		return nil
	}

	log.Println("Inserting framework record")
	insertFrameworkRecord := `INSERT INTO framework(sortID, Identifier, ParentIdentifier, Description, Guidance, 
	DisplayName, Recommendations, RequirementType, PolicyAndProcedureAIPromptTemplateId, ControlNarrativeAllPromptTemplateId,
	TestType, Framework, FrameworkID) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertFrameworkRecord)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(id, identifier, parentident, description, displayName, guidance, recommendations, requirementType, pandpPromptId, controlNarrativeId, testType, framework, frameworkId)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return err
}

func UpdateFrameworkLookupTable(db *sql.DB, missingFrameworkName, cename, uatStage, stageNumber, prodNumber, tableName, viewName string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	log.Printf("missing framework %s, cename: %s, uatStage: %s, stageNumber: %s, prodNumber: %s, tableName: %s, viewName: %s", missingFrameworkName, cename, uatStage, stageNumber, prodNumber, tableName, viewName)
	query := `INSERT INTO Framework_Lookup (EvidenceLibraryMappedName, CEFramework, FrameworkId_UAT, FrameworkId_Staging, FrameworkId_Prod, AirtableFramework, AirtableView)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(CEFramework) DO UPDATE SET
		  EvidenceLibraryMappedName = excluded.EvidenceLibraryMappedName,
		  FrameworkId_UAT = excluded.FrameworkId_UAT,
		  FrameworkId_Staging = excluded.FrameworkId_Staging,
		  FrameworkId_Prod = excluded.FrameworkId_Prod,
		  AirtableFramework = excluded.AirtableFramework,
		  AirtableView = excluded.AirtableView;
		`
	//UPDATE Framework_Lookup SET CEFramework = ?, FrameworkId_UAT = ?, FrameworkId_Staging = ?, FrameworkId_Prod = ? WHERE EvidenceLibraryMappedName = ?`

	statement, err := db.Prepare(query)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(missingFrameworkName, cename, uatStage, stageNumber, prodNumber, tableName, viewName)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return err
}

func UpdateBuildFramework_LookupTable(db *sql.DB, CEFramework, FrameworkId_UAT, FrameworkId_Staging, FrameworkId_Prod string) error {
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

	_, err = stmt.Exec(CEFramework, FrameworkId_UAT, FrameworkId_Staging, FrameworkId_Prod)
	if err != nil {

		return fmt.Errorf("failed to update framework lookup table: %v", err)
	}
	return nil
}
