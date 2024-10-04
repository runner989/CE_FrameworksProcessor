package database

import (
	"database/sql"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func ExportEvidenceMapReportToExcel(db *sql.DB, table string) error {
	evidenceList, err := getEvidenceSheet(db)
	if err != nil {
		return fmt.Errorf("error getting mapped evidence list: %v", err)
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Println("excel close err:", err)
		}
	}()

	f.NewSheet("Evidence")
	f.DeleteSheet("Sheet1")

	// Set headers for the Frameworks sheet
	f.SetCellValue("Evidence", "A1", "EvidenceID")
	f.SetCellValue("Evidence", "B1", "Evidence")
	f.SetCellValue("Evidence", "C1", "Description")
	f.SetCellValue("Evidence", "D1", "AnecdotesEvidenceIds")
	f.SetCellValue("Evidence", "E1", "Priority")
	f.SetCellValue("Evidence", "F1", "EvidenceType")

	rowIndex := 2
	for _, evidence := range evidenceList {
		f.SetCellValue("Evidence", "A"+strconv.Itoa(rowIndex), evidence.EvidenceID)
		f.SetCellValue("Evidence", "B"+strconv.Itoa(rowIndex), safeString(evidence.EvidenceTitle))
		f.SetCellValue("Evidence", "C"+strconv.Itoa(rowIndex), safeString(evidence.Description))
		f.SetCellValue("Evidence", "D"+strconv.Itoa(rowIndex), safeString(evidence.AnecdotesEvidenceIds))
		f.SetCellValue("Evidence", "E"+strconv.Itoa(rowIndex), safeString(evidence.Priority))
		f.SetCellValue("Evidence", "F"+strconv.Itoa(rowIndex), safeString(evidence.EvidenceType))
		rowIndex++
	}

	mappingList, err := getEvidenceMapping(db)
	if err != nil {
		return fmt.Errorf("error getting mapped evidence list: %v", err)
	}

	f.NewSheet("Mapping")

	f.SetCellValue("Mapping", "A1", "EvidenceID")
	f.SetCellValue("Mapping", "B1", "Framework")
	f.SetCellValue("Mapping", "C1", "FrameworkId")
	f.SetCellValue("Mapping", "D1", "Requirement")
	f.SetCellValue("Mapping", "E1", "Description")
	f.SetCellValue("Mapping", "F1", "Guidance")
	f.SetCellValue("Mapping", "G1", "RequirementType")
	f.SetCellValue("Mapping", "H1", "Delete")

	rowIndex = 2
	for _, mapping := range mappingList {
		f.SetCellValue("Mapping", "A"+strconv.Itoa(rowIndex), mapping.EvidenceID)
		f.SetCellValue("Mapping", "B"+strconv.Itoa(rowIndex), mapping.Framework)
		f.SetCellValue("Mapping", "C"+strconv.Itoa(rowIndex), mapping.FrameworkID)
		f.SetCellValue("Mapping", "D"+strconv.Itoa(rowIndex), safeString(mapping.Requirement))
		f.SetCellValue("Mapping", "E"+strconv.Itoa(rowIndex), safeString(mapping.Description))
		f.SetCellValue("Mapping", "F"+strconv.Itoa(rowIndex), safeString(mapping.Guidance))
		f.SetCellValue("Mapping", "G"+strconv.Itoa(rowIndex), safeString(mapping.RequirementType))
		f.SetCellValue("Mapping", "H"+strconv.Itoa(rowIndex), safeString(mapping.Delete))
		rowIndex++
	}

	// Define the "Mappings" folder path
	folderPath := "Mappings"

	// Check if the "Mappings" folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// Create the "Mappings" folder if it doesn't exist
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create Mappings folder: %v", err)
		}
	}
	// Save the file to disk
	// Get the current date in MMDDYYYY format
	currentDate := time.Now().Format("01022006")
	fileName := fmt.Sprintf("Evidence_Mapping_%s_%s.xlsx", table, currentDate)
	filePath := filepath.Join(folderPath, fileName)
	err = f.SaveAs(filePath)
	if err != nil {
		return fmt.Errorf("failed to save mappings: %v", err)
	}
	return nil
}

func ExportFrameworkToExcel(db *sql.DB, selectedFramework string) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()

	// Sheet 1: Frameworks
	f.NewSheet("Frameworks")
	f.DeleteSheet("Sheet1")
	frameworkQuery := `SELECT CEFramework, Version, Description, Comments FROM Framework_Lookup WHERE EvidenceLibraryMappedName = ?`
	rows, err := db.Query(frameworkQuery, selectedFramework)
	if err != nil {
		return fmt.Errorf("error querying Framework_Lookup for export: %v", err)
	}
	defer rows.Close()

	// Set headers for the Frameworks sheet
	f.SetCellValue("Frameworks", "A1", "Name")
	f.SetCellValue("Frameworks", "B1", "Version")
	f.SetCellValue("Frameworks", "C1", "Description")
	f.SetCellValue("Frameworks", "D1", "Comments")

	rowIndex := 2
	for rows.Next() {
		var version sql.NullInt64
		var name, description, comments sql.NullString
		err := rows.Scan(&name, &version, &description, &comments)
		if err != nil {
			return fmt.Errorf("error scanning Framework_Lookup row: %v", err)
		}
		f.SetCellValue("Frameworks", fmt.Sprintf("A%d", rowIndex), safeString(name))
		if version.Valid {
			if safeInt(version) == 0 {
				f.SetCellValue("Frameworks", fmt.Sprintf("B%d", rowIndex), 1)
			} else {
				f.SetCellValue("Frameworks", fmt.Sprintf("B%d", rowIndex), safeInt(version))
			}
		} else {
			f.SetCellValue("Frameworks", fmt.Sprintf("B%d", rowIndex), 1)
		}
		f.SetCellValue("Frameworks", fmt.Sprintf("C%d", rowIndex), safeString(description))
		f.SetCellValue("Frameworks", fmt.Sprintf("D%d", rowIndex), safeString(comments))
		rowIndex++
	}

	// Sheet 2: Requirements
	f.NewSheet("Requirements")
	requirementsQuery := `SELECT Identifier, ParentIdentifier, DisplayName, Description, Guidance, Recommendations, Observations, Notes, Tags, TestType, PolicyAndProcedureAIPromptTemplateId, ControlNarrativeAIPromptTemplateId FROM Framework WHERE Framework = ?`
	reqRows, err := db.Query(requirementsQuery, selectedFramework)
	if err != nil {
		return fmt.Errorf("error querying Framework table for requirements: %v", err)
	}
	defer reqRows.Close()

	// Set headers for the Requirements sheet
	f.SetCellValue("Requirements", "A1", "Identifier")
	f.SetCellValue("Requirements", "B1", "ParentIdentifier")
	f.SetCellValue("Requirements", "C1", "DisplayName")
	f.SetCellValue("Requirements", "D1", "Description")
	f.SetCellValue("Requirements", "E1", "Guidance")
	f.SetCellValue("Requirements", "F1", "Recommendations")
	f.SetCellValue("Requirements", "G1", "Observations")
	f.SetCellValue("Requirements", "H1", "Notes")
	f.SetCellValue("Requirements", "I1", "Tags")
	f.SetCellValue("Requirements", "J1", "TestType")
	f.SetCellValue("Requirements", "K1", "PolicyAndProcedureTemplateId")
	f.SetCellValue("Requirements", "L1", "ControlNarrativeAIPromptTemplateId")

	rowIndex = 2
	for reqRows.Next() {
		var identifier, displayName, description sql.NullString
		var policyAndProcedureTemplateId, controlNarrativeAIPromptTemplateId int
		var parentIdentifier, guidance, recommendations, observations, notes, tags, testType sql.NullString
		err := reqRows.Scan(&identifier, &parentIdentifier, &displayName, &description, &guidance, &recommendations, &observations, &notes, &tags, &testType, &policyAndProcedureTemplateId, &controlNarrativeAIPromptTemplateId)
		if err != nil {
			return fmt.Errorf("error scanning Framework row for requirements: %v", err)
		}
		//log.Printf("displayName: %v, Valid: %v", displayName.String, displayName.Valid)
		//log.Printf("description: %v, Valid: %v", description.String, description.Valid)

		f.SetCellValue("Requirements", fmt.Sprintf("A%d", rowIndex), safeString(identifier))
		f.SetCellValue("Requirements", fmt.Sprintf("B%d", rowIndex), safeString(parentIdentifier))
		f.SetCellValue("Requirements", fmt.Sprintf("C%d", rowIndex), safeString(displayName))
		f.SetCellValue("Requirements", fmt.Sprintf("D%d", rowIndex), safeString(description))
		f.SetCellValue("Requirements", fmt.Sprintf("E%d", rowIndex), safeString(guidance))
		f.SetCellValue("Requirements", fmt.Sprintf("F%d", rowIndex), safeString(recommendations))
		f.SetCellValue("Requirements", fmt.Sprintf("G%d", rowIndex), safeString(observations))
		f.SetCellValue("Requirements", fmt.Sprintf("H%d", rowIndex), safeString(notes))
		f.SetCellValue("Requirements", fmt.Sprintf("I%d", rowIndex), safeString(tags))
		f.SetCellValue("Requirements", fmt.Sprintf("J%d", rowIndex), safeString(testType))
		if policyAndProcedureTemplateId == 0 {
			f.SetCellValue("Requirements", fmt.Sprintf("K%d", rowIndex), "")
		} else {
			f.SetCellValue("Requirements", fmt.Sprintf("K%d", rowIndex), policyAndProcedureTemplateId)
		}
		if controlNarrativeAIPromptTemplateId == 0 {
			f.SetCellValue("Requirements", fmt.Sprintf("L%d", rowIndex), "")
		} else {
			f.SetCellValue("Requirements", fmt.Sprintf("L%d", rowIndex), controlNarrativeAIPromptTemplateId)
		}
		rowIndex++
	}

	// Sheet 3: Tests (Heading only)
	f.NewSheet("Tests")
	f.SetCellValue("Tests", "A1", "TestId")
	f.SetCellValue("Tests", "B1", "ParentTestId")
	f.SetCellValue("Tests", "C1", "RequirementId")
	f.SetCellValue("Tests", "D1", "TestType")
	f.SetCellValue("Tests", "E1", "Description")
	f.SetCellValue("Tests", "F1", "Guidance")
	f.SetCellValue("Tests", "G1", "Recommendations")
	f.SetCellValue("Tests", "H1", "Observations")
	f.SetCellValue("Tests", "I1", "Tags")

	// Sheet 4: Mappings (Heading only)
	f.NewSheet("Mappings")
	f.SetCellValue("Mappings", "A1", "FromIdentifier")
	f.SetCellValue("Mappings", "B1", "FromDescription")
	f.SetCellValue("Mappings", "C1", "ToFrameworkName")
	f.SetCellValue("Mappings", "D1", "ToFrameworkID")
	f.SetCellValue("Mappings", "E1", "ToFrameworkVersion")
	f.SetCellValue("Mappings", "F1", "ToIdentifier")
	f.SetCellValue("Mappings", "G1", "ToIdentifierType")
	f.SetCellValue("Mappings", "H1", "ToDescription")

	// Define the "Frameworks" folder path
	folderPath := "Frameworks"

	// Check if the "Frameworks" folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// Create the "Frameworks" folder if it doesn't exist
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create Frameworks folder: %v", err)
		}
	}
	// Save the file to disk
	// Get the current date in MMDDYYYY format
	currentDate := time.Now().Format("01022006")
	fixedFramework := strings.ReplaceAll(selectedFramework, ":", "-")
	fileName := fmt.Sprintf("%s-Framework-%s.xlsx", fixedFramework, currentDate)
	// Full path to the file
	filePath := filepath.Join(folderPath, fileName)

	err = f.SaveAs(filePath)
	if err != nil {
		return fmt.Errorf("failed to save Excel file: %v", err)
	}

	return nil
}

func safeString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func safeInt(ns sql.NullInt64) int64 {
	if ns.Valid {
		return ns.Int64
	}
	return 0
}
