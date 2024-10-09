package database

import (
	"cefp/structs"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func GetMissingFrameworks(db *sql.DB) ([]string, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	uniqueFrameworks := []string{}

	query := `SELECT DISTINCT Mapping.Framework
		FROM Mapping LEFT JOIN Framework_Lookup ON Mapping.Framework = Framework_Lookup.EvidenceLibraryMappedName
		WHERE (((Framework_Lookup.EvidenceLibraryMappedName) Is Null)) ORDER BY Framework;
		`
	//`
	//SELECT DISTINCT Mapping.Framework FROM Mapping
	//LEFT JOIN Framework_Lookup ON Mapping.Framework = Framework_Lookup.EvidenceLibraryMappedName
	//WHERE Framework_Lookup.EvidenceLibraryMappedName IS NULL ORDER BY Mapping.Framework;
	//`
	rows, err := db.Query(query)
	//SELECT DISTINCT Framework FROM Mapping WHERE Framework NOT IN (SELECT DISTINCT EvidenceLibraryMappedName FROM Framework_Lookup WHERE EvidenceLibraryMappedName IS NOT NULL) ORDER BY Framework;
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var framework string
		if err := rows.Scan(&framework); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		uniqueFrameworks = append(uniqueFrameworks, framework)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return uniqueFrameworks, nil
}

func GetFrameworkLookupFrameworks(db *sql.DB) ([]string, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	uniqueFrameworks := []string{}

	rows, err := db.Query("SELECT DISTINCT CEFramework FROM Framework_Lookup WHERE CEFramework IS NOT NULL ORDER BY CEFramework;")
	if err != nil {
		return nil, fmt.Errorf("error getting frameworks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var framework string
		if err := rows.Scan(&framework); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		uniqueFrameworks = append(uniqueFrameworks, framework)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return uniqueFrameworks, nil
}

func GetFrameworkInfoBackend(db *sql.DB, framework string) (map[string]interface{}, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	query := "SELECT EvidenceLibraryMappedName, AirtableTableID, AirtableFramework, AirtableView, FrameworkId_UAT, FrameworkId_Staging, FrameworkId_Staging FROM Framework_Lookup WHERE CEFramework = ?"
	row := db.QueryRow(query, framework)

	var details map[string]interface{}
	var evidenceLibraryName, airtableID, airtableFramework, airtableView, frameworkId_UAT, frameworkId_Stage, frameworkId_Prod sql.NullString

	err := row.Scan(&evidenceLibraryName, &airtableID, &airtableFramework, &airtableView, &frameworkId_UAT, &frameworkId_Stage, &frameworkId_Prod)
	if err != nil {
		return nil, fmt.Errorf("error querying framework details: %w", err)
	}

	details = map[string]interface{}{
		"CEName":                    framework,
		"EvidenceLibraryMappedName": nullStringToString(evidenceLibraryName),
		"AirtableTableID":           nullStringToString(airtableID),
		"AirtableFramework":         nullStringToString(airtableFramework),
		"AirtableView":              nullStringToString(airtableView),
		"FrameworkId_UAT":           nullStringToString(frameworkId_UAT),
		"FrameworkId_Staging":       nullStringToString(frameworkId_Stage),
		"FrameworkId_Prod":          nullStringToString(frameworkId_Prod),
	}

	return details, nil
}

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func GetMappedFrameworkRecords(db *sql.DB) ([]string, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	query := "SELECT DISTINCT Framework FROM Mapping ORDER BY Framework;"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying mapped frameworks: %w", err)
	}
	defer rows.Close()
	var frameworks []string
	for rows.Next() {
		var framework string
		if err := rows.Scan(&framework); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		frameworks = append(frameworks, framework)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return frameworks, nil
}

func GetReadyFrameworks(db *sql.DB) ([]structs.FrameworkLookup, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	query := "SELECT DISTINCT AirtableBase, AirtableTableID, AirtableFramework, AirtableView, CEFramework FROM Framework_Lookup WHERE CEFramework IS NOT NULL AND AirtableFramework IS NOT NULL AND AirtableBase IS NOT NULL AND AirtableView IS NOT NULL AND AirtableTableID IS NOT NULL"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying framework_lookup: %w", err)
	}
	defer rows.Close()

	var frameworks []structs.FrameworkLookup

	for rows.Next() {
		var framework structs.FrameworkLookup

		if err := rows.Scan(
			&framework.TableBase,
			&framework.TableID,
			&framework.TableName,
			&framework.TableView,
			&framework.CeName,
		); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		frameworks = append(frameworks, framework)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return frameworks, nil
}

func GetDistinctFrameworks(db *sql.DB) ([]string, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	query := ` 
		SELECT DISTINCT Framework.Framework 
		FROM Framework_Lookup 
		INNER JOIN Framework ON Framework_Lookup.EvidenceLibraryMappedName = Framework.Framework
		WHERE Framework IS NOT NULL AND EvidenceLibraryMappedName IS NOT NULL ORDER BY Framework.Framework;
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying frameworks: %w", err)
	}
	defer rows.Close()

	var frameworks []string
	for rows.Next() {
		var framework string
		err := rows.Scan(&framework)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		frameworks = append(frameworks, framework)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return frameworks, nil
}

func CheckForMissing(db *sql.DB, table string) ([]int, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	query := fmt.Sprintf("SELECT DISTINCT [CEMapping_%s].EvidenceID FROM [CEMapping_%s] LEFT JOIN Evidence ON [CEMapping_%s].EvidenceID = Evidence.EvidenceID WHERE (((Evidence.EvidenceID) Is Null));", table, table, table)

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed checking for missing evidence: %v", err)
	}
	defer rows.Close()

	evidenceIDs := []int{}
	for rows.Next() {
		var evidenceID int
		if err := rows.Scan(&evidenceID); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		evidenceIDs = append(evidenceIDs, evidenceID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return evidenceIDs, nil
}

func getEvidenceSheet(db *sql.DB) ([]structs.EvidenceRecord, error) {
	// This query is to get only Evidence that is mapped to a framework
	query := `
		WITH CEFrameworkMapping AS (
		-- First Query: Select CE Framework mappings
		SELECT Mapping.EvidenceID, Mapping.Framework, Mapping.Requirement
		FROM Mapping
		WHERE Mapping.Framework = 'CE Framework'
		),
		NonCEFrameworkMapping AS (
			-- Second Query: Select Non-CE Framework mappings
			SELECT DISTINCT CEFrameworkMapping.EvidenceID AS CEEvidenceID, Mapping.EvidenceID, Mapping.Framework, CEFrameworkMapping.Requirement
			FROM CEFrameworkMapping
			LEFT JOIN Mapping ON CEFrameworkMapping.EvidenceID = Mapping.EvidenceID
			WHERE Mapping.Framework <> 'CE Framework' OR CEFrameworkMapping.Requirement = ''
		)
		-- Final Query: Join CE and Non-CE mappings with the Evidence table
		SELECT DISTINCT Evidence.*
		FROM CEFrameworkMapping
		INNER JOIN NonCEFrameworkMapping ON CEFrameworkMapping.EvidenceID = NonCEFrameworkMapping.CEEvidenceID
		INNER JOIN Evidence ON CEFrameworkMapping.EvidenceID = Evidence.EvidenceID;
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying CE Framework Mapping: %v", err)
	}
	defer rows.Close()
	var evidenceList []structs.EvidenceRecord
	for rows.Next() {
		var evidence structs.EvidenceRecord
		err := rows.Scan(&evidence.EvidenceID, &evidence.EvidenceTitle, &evidence.Description, &evidence.AnecdotesEvidenceIds, &evidence.Priority, &evidence.EvidenceType)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		evidenceList = append(evidenceList, evidence)
	}

	return evidenceList, nil
}

func GetDeletions(db *sql.DB, table string) ([]structs.EvidenceMapRecord, error) {
	var designatedTable string
	if table == "UAT" {
		designatedTable = "Staging"
	} else {
		designatedTable = table
	}
	delQuery := fmt.Sprintf(`WITH Mapping_%s_FWID AS (
		SELECT Mapping.EvidenceID, Mapping.Framework, Framework_Lookup.FrameworkId_%s, Mapping.Requirement, Mapping.Description, Mapping.Guidance, Mapping.RequirementType, Mapping."Delete", Framework_Lookup.CEFramework
		FROM Mapping LEFT JOIN Framework_Lookup ON Mapping.Framework = Framework_Lookup.EvidenceLibraryMappedName
		)
		SELECT [CEMapping_%s].EvidenceID, [CEMapping_%s].Framework, [CEMapping_%s].FrameworkId, TRIM([CEMapping_%s].Requirement), [CEMapping_%s].Description, [CEMapping_%s].Guidance, [CEMapping_%s].RequirementType, "X" AS "Delete"
		FROM [CEMapping_%s] LEFT JOIN Mapping_%s_FWID ON (TRIM([CEMapping_%s].Requirement) = TRIM(Mapping_%s_FWID.Requirement)) AND ([CEMapping_%s].FrameworkId = Mapping_%s_FWID.FrameworkId_%s) AND ([CEMapping_%s].EvidenceID = Mapping_%s_FWID.EvidenceID)
		WHERE ((Mapping_%s_FWID.EvidenceID Is Null));
		`, designatedTable, designatedTable, designatedTable, designatedTable, designatedTable, designatedTable, designatedTable, designatedTable, designatedTable, designatedTable, designatedTable, designatedTable, designatedTable, designatedTable, designatedTable, designatedTable, designatedTable, designatedTable, designatedTable)

	delRows, err := db.Query(delQuery)
	if err != nil {
		return nil, fmt.Errorf("error getting deletions list: %v", err)
	}
	defer delRows.Close()
	var deleteList []structs.EvidenceMapRecord
	for delRows.Next() {
		var deleted structs.EvidenceMapRecord
		err := delRows.Scan(&deleted.EvidenceID, &deleted.Framework, &deleted.FrameworkID, &deleted.Requirement, &deleted.Description, &deleted.Guidance, &deleted.RequirementType, &deleted.Delete)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		deleteList = append(deleteList, deleted)
	}
	if err := delRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return deleteList, nil
}

func getEvidenceMapping(db *sql.DB, table string) ([]structs.EvidenceMapRecord, error) {
	var designatedTable string
	if table == "UAT" {
		designatedTable = "Staging"
	} else {
		designatedTable = table
	}
	query := fmt.Sprintf(`
		WITH CEFrameworkMapping AS (
			-- First Query: Select CE Framework mappings
			SELECT Mapping.EvidenceID, Mapping.Framework, Mapping.Requirement
			FROM Mapping
			WHERE Mapping.Framework = 'CE Framework'
		),
		NonCEFrameworkMapping AS (
			-- Second Query: Select Non-CE Framework mappings
			SELECT DISTINCT CEFrameworkMapping.EvidenceID AS CEEvidenceID, Mapping.EvidenceID, Mapping.Framework, CEFrameworkMapping.Requirement
			FROM CEFrameworkMapping
			LEFT JOIN Mapping ON CEFrameworkMapping.EvidenceID = Mapping.EvidenceID
			WHERE Mapping.Framework <> 'CE Framework' OR CEFrameworkMapping.Requirement = ''
		),
		EvidenceExport AS (
			-- Third Query: Select relevant evidence records
			SELECT DISTINCT Evidence.*
			FROM CEFrameworkMapping
			INNER JOIN NonCEFrameworkMapping ON CEFrameworkMapping.EvidenceID = NonCEFrameworkMapping.CEEvidenceID
			INNER JOIN Evidence ON CEFrameworkMapping.EvidenceID = Evidence.EvidenceID
		)
		
		-- Final Query: Export relevant data
		SELECT DISTINCT 
			Mapping.EvidenceID, 
			Framework_Lookup.CEFramework AS Framework, 
			Framework_Lookup.FrameworkId_%s AS FrameworkId, 
			Mapping.Requirement, 
			Mapping.Description, 
			Mapping.Guidance, 
			'Requirement' AS RequirementType, 
			Mapping."Delete"
		FROM EvidenceExport
		INNER JOIN Mapping ON EvidenceExport.EvidenceID = Mapping.EvidenceID
		INNER JOIN Framework_Lookup ON Mapping.Framework = Framework_Lookup.EvidenceLibraryMappedName
		WHERE Framework_Lookup.FrameworkId_%s <> 0 
		  AND Mapping.Requirement <> ''
	`, designatedTable, designatedTable)

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying CE Framework Mapping: %v", err)
	}
	defer rows.Close()
	var mappingList []structs.EvidenceMapRecord
	for rows.Next() {
		var mapping structs.EvidenceMapRecord
		err := rows.Scan(&mapping.EvidenceID, &mapping.Framework, &mapping.FrameworkID, &mapping.Requirement, &mapping.Description, &mapping.Guidance, &mapping.RequirementType, &mapping.Delete)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		mappingList = append(mappingList, mapping)
	}
	return mappingList, nil
}

func GetEvidenceMappingCounts(db *sql.DB) ([]structs.FrameworkMappedCount, error) {
	qry := `
		WITH Unique_EvID AS (
		SELECT DISTINCT Mapping.Framework, Mapping.EvidenceID
		FROM Mapping
		GROUP BY Mapping.Framework, Mapping.EvidenceID
		)
		SELECT Unique_EvID.Framework, Count(Unique_EvID.EvidenceID) AS CountOfEvidenceID
		FROM Unique_EvID
		GROUP BY Unique_EvID.Framework;
		`

	rows, err := db.Query(qry)
	if err != nil {
		return nil, fmt.Errorf("error querying CE Framework Mapping: %v", err)
	}
	defer rows.Close()

	var results []structs.FrameworkMappedCount

	for rows.Next() {
		var count structs.FrameworkMappedCount
		err := rows.Scan(&count.Framework, &count.Count)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		results = append(results, count)
	}
	return results, nil
}
