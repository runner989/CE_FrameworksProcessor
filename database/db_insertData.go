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

func UpdateFrameworkLookupTable(db *sql.DB, missingFrameworkName, name, uatStage, stageNumber, prodNumber string) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	fmt.Printf("ELMN: %s, CEName: %s, UAT: %s, stage: %s, prod: %s\n", missingFrameworkName, name, uatStage, stageNumber, prodNumber)

	query := `INSERT INTO Framework_Lookup (EvidenceLibraryMappedName, CEFramework, FrameworkId_UAT, FrameworkId_Staging, FrameworkId_Prod)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(EvidenceLibraryMappedName) DO UPDATE SET
		  CEFramework = excluded.CEFramework,
		  FrameworkId_UAT = excluded.FrameworkId_UAT,
		  FrameworkId_Staging = excluded.FrameworkId_Staging,
		  FrameworkId_Prod = excluded.FrameworkId_Prod;
		`
	//UPDATE Framework_Lookup SET CEFramework = ?, FrameworkId_UAT = ?, FrameworkId_Staging = ?, FrameworkId_Prod = ? WHERE EvidenceLibraryMappedName = ?`

	statement, err := db.Prepare(query)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(missingFrameworkName, name, uatStage, stageNumber, prodNumber)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return err
}
