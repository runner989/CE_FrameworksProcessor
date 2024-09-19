package database

import (
	"database/sql"
	"log"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

// create new database connection and create the file if it doen't exit
func NewDB(path string) (*sql.DB, error) {
	dbase := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	sqlDB, err := dbase.ensureDB()
	return sqlDB, err
}

// create the database file if it does not exist
func (db *DB) ensureDB() (*sql.DB, error) {
	_, err := os.ReadFile(db.path)
	if os.IsNotExist(err) {
		file, err := os.Create(db.path)
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()

		db, err := sql.Open(
			"sqlite3",
			db.path,
		)
		if err != nil {
			return nil, err
		}
		err = CreateFrameworkTable(db)
		if err != nil {
			log.Fatalf("unable to create framework table: %v", err)
		}
		err = CreateAirTableBaseTable(db)
		if err != nil {
			log.Fatalf("unable to create Airtable Base table: %v", err)
		}
		err = CreateCEMappingProdTable(db)
		if err != nil {
			log.Fatalf("unable to create CE Mapping Prod table: %v", err)
		}
		err = CreateCEMappingStagingTable(db)
		if err != nil {
			log.Fatalf("unable to create CE Mapping Staging table: %v", err)
		}
		err = CreateEvidenceTable(db)
		if err != nil {
			log.Fatalf("unable to create Evidence table: %v", err)
		}
		err = CreateFrameworkLookupTable(db)
		if err != nil {
			log.Fatalf("unable to create Framework Lookup table: %v", err)
		}
		err = CreatePlaceholderMappingsTable(db)
		if err != nil {
			log.Fatalf("unable to create Placeholder Mappings table: %v", err)
		}
		err = CreatetblMappingTable(db)
		if err != nil {
			log.Fatalf("unable to create tblMapping table: %v", err)
		}
		err = CreateTestProceduresTable(db)
		if err != nil {
			log.Fatalf("unable to create Test Procedures table: %v", err)
		}
		err = CreateTestProceduresLookupTable(db)
		if err != nil {
			log.Fatalf("unable to create Test Procedures Lookup table: %v", err)
		}
		defer db.Close()
		log.Println("cefp.db created")
	}
	sqliteDatabase, err := sql.Open(
		"sqlite3",
		"cefp.db",
	)
	if err != nil {
		return nil, err
	}
	//defer sqliteDatabase.Close()
	return sqliteDatabase, err
}
