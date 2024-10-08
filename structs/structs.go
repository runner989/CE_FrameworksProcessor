package structs

import "database/sql"

type AirtableBases struct {
	Bases  []Base `json:"bases"`
	Offset string `json:"offset,omitempty"`
}

type Base struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	PermissionLevel string `json:"permissionLevel"`
}

type AirtableResponse struct {
	Records []Evidence `json:"records"`
	Offset  string     `json:"offset,omitempty"`
}

type Record struct {
	ID          string                 `json:"id"`
	CreatedTime string                 `json:"createdTime"`
	Fields      map[string]interface{} `json:"fields"`
}

type Table struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	PrimaryFieldId  string                 `json:"primaryFieldId"`
	Fields          map[string]interface{} `json:"fields"`
	PermissionLevel string                 `json:"permissionLevel"`
}

type TablesResponse struct {
	Tables []Table `json:"tables"`
}

type Evidence struct {
	ID     string                 `json:"id"`
	Fields map[string]interface{} `json:"fields"`
}

type EvidenceRecord struct {
	EvidenceID           int            `json:"EvidenceID"`
	EvidenceTitle        sql.NullString `json:"Evidence Title"`
	Description          sql.NullString `json:"Description_FromEvidence"`
	Requirement          sql.NullString `json:"Requirement"`
	AnecdotesEvidenceIds sql.NullString `json:"AnecdotesEvidenceIds"`
	Priority             sql.NullString `json:"Priority"`
	EvidenceType         sql.NullString `json:"Evidence Type"`
}

type EvidenceMapRecord struct {
	EvidenceID      int // Assuming EvidenceID should be an int
	Framework       string
	FrameworkID     int // Assuming FrameworkID should be an int
	Requirement     sql.NullString
	Description     sql.NullString
	Guidance        sql.NullString
	RequirementType sql.NullString
	Delete          sql.NullString
}

type FrameworkMappedCount struct {
	Framework string
	Count     int
}

type FrameworkTable struct {
	TableName string `json:"tableName"`
	TableID   string `json:"tableID"`
	TableView string `json:"tableView"`
}

type Framework struct {
	ID          string                 `json:"id"`
	CreatedTime string                 `json:"createdTime"`
	Fields      map[string]interface{} `json:"fields"`
}

type FrameworkRecord struct {
	SortID           int    `json:"sortId"`
	PromptID         int    `json:"promptId"`
	ControlNarrative int    `json:"controlNarrative"`
	FrameworkName    string `json:"frameworkName"`
	Identifier       string `json:"identifier"`
	ParentID         string `json:"parentId"`
	DisplayName      string `json:"displayName"`
	Description      string `json:"description"`
	Guidance         string `json:"guidance"`
	Tags             string `json:"tags"`
	TestType         string `json:"testType"`
}

type ViewsRoot struct {
	Views []View `json:"views"`
}

type View struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	PersonalForUserId string `json:"personalForUserId,omitempty"`
	Type              string `json:"type"`
}

type FrameworksResponse struct {
	Records []Framework `json:"records"`
	Offset  string      `json:"offset,omitempty"`
}

type FrameworkLookup struct {
	RowID       sql.NullFloat64 `json:"rowId"`
	MappedName  sql.NullString  `json:"mappedName"`
	CeName      sql.NullString  `json:"ceName"`
	UatStage    sql.NullFloat64 `json:"uatStage"`
	StageNumber sql.NullFloat64 `json:"stageNumber"`
	ProdNumber  sql.NullFloat64 `json:"prodNumber"`
	TableBase   sql.NullString  `json:"tableBase"`
	TableID     sql.NullString  `json:"tableID"`
	TableName   sql.NullString  `json:"tableName"`
	TableView   sql.NullString  `json:"tableView"`
	Version     sql.NullString  `json:"version"`
	Description sql.NullString  `json:"description"`
	Comments    sql.NullString  `json:"comments"`
}

type FileMetadata struct {
	Filename string `json:"filename"`
	Content  []byte `json:"content"`
}

type CEMappingRecord struct {
	EvidenceID  int64  `json:"EvidenceID"`
	Framework   string `json:"framework"`
	Requirement string `json:"requirement"`
}

type NonCEFrameworkRecord struct {
	CE_EvidenceID     int64  `json:"ce_evidence_id"`      // From CE Framework Mapping
	MappingEvidenceID int64  `json:"mapping_evidence_id"` // From tblMapping
	Framework         string `json:"framework"`
	Requirement       string `json:"requirement"`
}
