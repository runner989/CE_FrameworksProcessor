package database

import (
	"database/sql"
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
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking for table existence: %v", err)
	}

	if tableName != "" {
		log.Println("Framework table already exists, skipping creation")
		return nil
	}

	createTableSQL := `CREATE TABLE Framework (
		"sortID" integer NOT NULL PRIMARY KEY,
		"Identifier" TEXT,
		"ParentIdentifier" TEXT,
		"Description" TEXT,
		"Guidance" TEXT,
		"DisplayName" TEXT,
		"Recommendations" TEXT,
		"RequirementType" TEXT,
		"PolicyAndProcedureAIPromptTemplateId" integer,
		"ControlNarrativeAllPromptTemplateId" integer,
		"TestType" TEXT,
		"Framework" TEXT,
		"FrameworkID" TEXT
	);`

	log.Println("Create Framework table...")
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
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
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking for table existence: %v", err)
	}

	if tableName != "" {
		log.Println("Framework_Lookup table already exists, skipping creation")
		return nil
	}

	createTableSQL := `CREATE TABLE Framework_Lookup (
		"AirtableBase" TEXT,
		"AirtableFramework" TEXT,
		"AirtableView" TEXT,
		"EvidenceLibraryMappedName" TEXT,
		"CEFramework" TEXT,
		"FrameworkId_Staging" INTEGER,
		"FrameworkId_Prod" INTEGER,
		"FrameworkId_UAT" INTEGER,
		"Version" INTEGER,
		"Description" TEXT,
		"Comments" TEXT
	);`

	log.Println("Create Framework_Lookup table...")
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
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
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking for table existence: %v", err)
	}

	if tableName != "" {
		log.Println("Airtable_Base table already exists, skipping creation")
		return nil
	}

	createTableSQL := `CREATE TABLE Airtable_Base (
		"ID" INTEGER,
		"BaseName" TEXT,
		"BaseID" TEXT
	);`

	log.Println("Create Airtable_Base table...")
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
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
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking for table existence: %v", err)
	}

	if tableName != "" {
		log.Println("CEMapping-Prod table already exists, skipping creation")
		return nil
	}

	createTableSQL := `CREATE TABLE CEMapping_Prod (
		"EvidenceID" integer NOT NULL PRIMARY KEY,
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
	statement.Exec()
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
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking for table existence: %v", err)
	}

	if tableName != "" {
		log.Println("CEMapping_Staging table already exists, skipping creation")
		return nil
	}

	createTableSQL := `CREATE TABLE CEMapping_Staging (
		"EvidenceID" integer NOT NULL PRIMARY KEY,
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
	statement.Exec()
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
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking for table existence: %v", err)
	}

	if tableName != "" {
		log.Println("Evidence table already exists, skipping creation")
		return nil
	}

	createTableSQL := `CREATE TABLE Evidence (
		"EvidenceID" integer NOT NULL PRIMARY KEY,
		"Evidence" TEXT,
		"Description" TEXT,
		"AnecdotesEvidenceIds" TEXT,
		"Priority" TEXT,
		"EvidenceType" TEXT
	);`

	log.Println("Create Evidence table...")
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("Evidence table created")
	return err
}

func CreatePlaceholderMappingsTable(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	var tableName string
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='Placeholder_Mappings';`
	err := db.QueryRow(query).Scan(&tableName)
	if err != nil && err != sql.ErrNoRows {
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
	statement.Exec()
	log.Println("Placeholder_Mappings table created")
	return err
}

func CreatetblMappingTable(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	var tableName string
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='tblMapping';`
	err := db.QueryRow(query).Scan(&tableName)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking for table existence: %v", err)
	}

	if tableName != "" {
		log.Println("tblMapping table already exists, skipping creation")
		return nil
	}

	createTableSQL := `CREATE TABLE tblMapping (
		"EvidenceID" INTEGER,
		"Framework" TEXT,
		"FrameworkId" TEXT,
		"Requirement" TEXT,
		"Description" TEXT,
		"Guidance" TEXT,
		"RequirementType" TEXT,
		"Delete" TEXT
	);`

	log.Println("Create Mappings table...")
	statement, err := db.Prepare(createTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("tblMapping table created")
	return err
}

func CreateTestProceduresTable(db *sql.DB) error {
	return nil
}

func CreateTestProceduresLookupTable(db *sql.DB) error {
	return nil
}
