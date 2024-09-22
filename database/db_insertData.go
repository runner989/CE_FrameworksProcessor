package database

import (
	"database/sql"
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
	if err != nil && err != sql.ErrNoRows {
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
