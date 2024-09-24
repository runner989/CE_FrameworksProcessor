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

	// fmt.Printf("ELMN: %s, CEName: %s, UAT: %s, stage: %s, prod: %s\n", missingFrameworkName, name, uatStage, stageNumber, prodNumber)

	query := `INSERT INTO Framework_Lookup (EvidenceLibraryMappedName, CEFramework, FrameworkId_UAT, FrameworkId_Staging, FrameworkId_Prod)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(CEFramework) DO UPDATE SET
		  EvidenceLibraryMappedName = excluded.EvidenceLibraryMappedName,
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

func UpdateBuildFrameworkLookupTable(db *sql.DB, records []map[string]interface{}) error {
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
	// 	INSERT OR REPLACE INTO Framework_Lookup (CEFramework, FrameworkId_UAT, FrameworkId_Staging, FrameworkId_Prod)
	// 	VALUES (?, ?, ?, ?);
	// `
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

	for _, fields := range records {
		ceFramework, _ := fields["Name"].(string)
		frameworkIdUAT, _ := fields["UAT_Stage"].(string)
		frameworkIdStaging, _ := fields["Stage Framework Number"].(string)
		frameworkIdProd, _ := fields["Production Framework Number"].(string)

		if ceFramework == "" {
			continue // Skip records without a name
		}

		_, err = stmt.Exec(ceFramework, frameworkIdUAT, frameworkIdStaging, frameworkIdProd)
		if err != nil {
			return fmt.Errorf("failed to update framework lookup table: %v", err)
		}
	}
	return nil
}
