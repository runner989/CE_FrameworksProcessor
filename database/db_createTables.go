package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CreateFrameworkTable(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	var tableName string
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='Framework';`
	err := db.QueryRow(query).Scan(&tableName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error checking for table existence: %v", err)
	}

	if tableName != "" {
		log.Println("Framework table already exists, skipping creation")
		return nil
	}

	createTableSQL := `CREATE TABLE Framework (
		"sortID" INTEGER,
		"Framework" TEXT,
		"FrameworkID" TEXT,
		"Identifier" TEXT,
		"ParentIdentifier" TEXT,
		"DisplayName" TEXT,
		"Description" TEXT,
		"Guidance" TEXT,
		"Observations" TEXT,
		"Recommendations" TEXT,
		"Notes" TEXT, 
		"Tags" TEXT,
		"TestType" TEXT,
		"RequirementType" TEXT,
		"PolicyAndProcedureAIPromptTemplateId" integer,
		"ControlNarrativeAIPromptTemplateId" integer
	);`

	log.Println("Create Framework table...")
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		return fmt.Errorf("error creating Framework table: %v", err)
	}
	log.Println("Framework table created")
	return err
}

func CreateFrameworkLookupTable(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	var tableName string
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='Framework_Lookup';`
	err := db.QueryRow(query).Scan(&tableName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error checking for table existence: %v", err)
	}

	if tableName != "" {
		log.Println("Framework_Lookup table already exists, skipping creation")
		return nil
	}

	createTableSQL := `CREATE TABLE Framework_Lookup (
		"AirtableBase" TEXT,
		"AirtableTableID" TEXT,
		"AirtableFramework" TEXT,
		"AirtableView" TEXT,
		"EvidenceLibraryMappedName" TEXT,
		"CEFramework" TEXT UNIQUE,
		"FrameworkId_UAT" INTEGER,
		"FrameworkId_Staging" INTEGER,
		"FrameworkId_Prod" INTEGER,
		"Version" INTEGER,
		"Description" TEXT,
		"Comments" TEXT
	);`

	log.Println("Create Framework_Lookup table...")
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		return fmt.Errorf("error creating Framework_Lookup table: %v", err)
	}
	log.Println("Framework_Lookup table created")
	return err
}

func CreateAirTableBaseTable(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	var tableName string
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='Airtable_Base';`
	err := db.QueryRow(query).Scan(&tableName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error checking for table existence: %v", err)
	}

	if tableName != "" {
		log.Println("Airtable_Base table already exists, skipping creation")
		return nil
	}

	createTableSQL := `CREATE TABLE Airtable_Base (
		"ID" INTEGER PRIMARY KEY,
		"BaseName" TEXT,
		"BaseID" TEXT
	);`

	log.Println("Create Airtable_Base table...")
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		return fmt.Errorf("error creating Airtable_Base table: %v", err)
	}
	log.Println("Airtable_Base table created")
	return err
}

func CreateCEMappingProdTable(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	var tableName string
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='CEMapping_Prod';`
	err := db.QueryRow(query).Scan(&tableName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error checking for table existence: %v", err)
	}

	if tableName != "" {
		log.Println("CEMapping-Prod table already exists, skipping creation")
		return nil
	}

	createTableSQL := `CREATE TABLE CEMapping_Prod (
		"EvidenceID" INTEGER,
		"Framework" TEXT,
		"FrameworkId" INTEGER,
		"Requirement" TEXT,
		"Description" TEXT,
		"Guidance" TEXT,
		"RequirementType" TEXT,
		"Delete" TEXT
	);`

	log.Println("Create CEMapping_Prod table...")
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		return fmt.Errorf("error creating CEMapping_Prod table: %v", err)
	}
	log.Println("CEMapping_Prod table created")
	return err
}

func CreateCEMappingStagingTable(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	var tableName string
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='CEMapping_Staging';`
	err := db.QueryRow(query).Scan(&tableName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error checking for table existence: %v", err)
	}

	if tableName != "" {
		log.Println("CEMapping_Staging table already exists, skipping creation")
		return nil
	}

	createTableSQL := `CREATE TABLE CEMapping_Staging (
		"EvidenceID" INTEGER,
		"Framework" TEXT,
		"FrameworkId" INTEGER,
		"Requirement" TEXT,
		"Description" TEXT,
		"Guidance" TEXT,
		"RequirementType" TEXT,
		"Delete" TEXT
	);`

	log.Println("Create CEMapping_Staging table...")
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		return fmt.Errorf("error creating CEMapping_Staging table: %v", err)
	}
	log.Println("CEMapping_Staging table created")
	return err
}

func CreateEvidenceTable(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	var tableName string
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='Evidence';`
	err := db.QueryRow(query).Scan(&tableName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error checking for table existence: %v", err)
	}

	if tableName != "" {
		log.Println("Evidence table already exists, skipping creation")
		return nil
	}

	createTableSQL := `CREATE TABLE Evidence (
		"EvidenceID" INTEGER NOT NULL PRIMARY KEY,
		"Evidence" TEXT,
		"Description" TEXT,
		"AnecdotesEvidenceIds" TEXT,
		"Priority" TEXT,
		"EvidenceType" TEXT
	);`

	log.Println("Create Evidence table...")
	statement, createErr := db.Prepare(createTableSQL)
	if createErr != nil {
		log.Fatal(createErr.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		return fmt.Errorf("error creating Evidence table: %v", err)
	}
	log.Println("Evidence table created")
	return nil
}

func CreatePlaceholderMappingsTable(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	var tableName string
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='Placeholder_Mappings';`
	err := db.QueryRow(query).Scan(&tableName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error checking for table existence: %v", err)
	}

	if tableName != "" {
		log.Println("Placeholder_Mappings table already exists, skipping creation")
		return nil
	}

	createTableSQL := `CREATE TABLE Placeholder_Mappings (
		"FromIdentifier" TEXT,
		"FromDescription" TEXT,
		"ToFrameworkName" TEXT,
		"ToFrameworkID" TEXT,
		"ToFrameworkVersion" TEXT,
		"ToIdentifier" TEXT,
		"ToIdentifierType" TEXT,
		"ToDescription" TEXT
	);`

	log.Println("Create Placeholder_Mappings table...")
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		return fmt.Errorf("error creating Placeholder_Mappings table: %v", err)
	}
	log.Println("Placeholder_Mappings table created")
	return err
}

func CreateMappingTable(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	var tableName string
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='Mapping';`
	err := db.QueryRow(query).Scan(&tableName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error checking for table existence: %v", err)
	}

	if tableName != "" {
		log.Println("Mapping table already exists, skipping creation")
		return nil
	}

	createTableSQL := `CREATE TABLE Mapping (
		"EvidenceID" INTEGER,
		"Framework" TEXT,
		"FrameworkId" TEXT,
		"Requirement" TEXT,
		"Description" TEXT,
		"Guidance" TEXT,
		"RequirementType" TEXT,
		"Delete" TEXT
	);`

	log.Println("Create Mapping table...")
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		return fmt.Errorf("error creating Mapping table: %v", err)
	}
	log.Println("Mapping table created")
	return err
}

func CreateTestProceduresTable(db *sql.DB) error {
	return nil
}

func CreateTestProceduresLookupTable(db *sql.DB) error {
	return nil
}
