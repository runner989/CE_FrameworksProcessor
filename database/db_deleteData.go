package database

import (
	"cefp/structs"
	"database/sql"
	"fmt"
	"log"
)

func DeleteFromFrameworkLookup(db *sql.DB, framework structs.FrameworkLookup) error {
	qry := "DELETE FROM Framework_Lookup WHERE ROWID=?"
	stmt, err := db.Prepare(qry)
	if err != nil {
		return fmt.Errorf("unable to prepare the delete from framework_lookup query: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(framework.RowID)
	if err != nil {
		return fmt.Errorf("unable to delete from framework_lookup query: %v", err)
	}
	return nil
}

func DeleteFromFramework(db *sql.DB, framework structs.FrameworkLookup) error {
	qry := "DELETE FROM Framework WHERE Framework=?"
	stmt, err := db.Prepare(qry)
	if err != nil {
		return fmt.Errorf("unable to prepare the delete from framework query: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(framework.TableName.String)
	if err != nil {
		log.Printf("failed to delete from framework query: %v", err)
		return fmt.Errorf("unable to delete from framework query: %v", err)
	}
	return nil
}
