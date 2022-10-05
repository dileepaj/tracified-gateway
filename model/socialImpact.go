package model

import (
	"time"
)

type BuildFormula struct {
	FormulaID       string
	FormulaType     string
	Name            string
	FieldCount      uint32
	Query           string
	FormulaAsString string
	TenantId        string
	Metric          FormulaItem
	Activity        Activity
	Expert          Expert
	Formula         []FormulaItem
	Active          bool
	CreatedAt       string
	UpdatedAt       string
}
type FormulaItem struct {
	ID           string
	Name         string
	Description  string
	value        float32
	VariableType string
	DataType     string
	URL          string
	Unit         string
}

type FormulaItemRequest struct {
	ID                string `json:"ID" bson:"id" validate:"required"`
	Value             float64
	Name              string
	Description       string
	Key               string
	Type              string `json:"Type" bson:"type" validate:"required"`
	MeasurementUnit   string
	MetricReferenceId string
	Precision         uint16
	MetricReference   MetricReference
}

type Expert struct {
	ExpertID string `json:"ExpertID" bson:"expertID" validate:"required"`
	ExpertPK string `json:"ExpertPK" bson:"expertPK" validate:"required"`
}

type FormulaIDMap struct {
	FormulaID string
	MapID     uint64
}

type ExpertIDMap struct {
	ExpertID  string
	ExpertPK  string
	MapID     uint64
	FormulaID string
}

type Counters struct {
	ID            string
	SequenceValue uint64
}

type ValueIDMap struct {
	ValueId   string
	ValueType string
	ValueName string
	Key       string
	FormulaID string
	MapID     uint64
}

type UnitIDMap struct {
	Unit  string
	MapID uint16
}

type FormulaBuildingRequest struct {
	Blockchain      string               `json:"Blockchain" bson:"blockchain" validate:"required"`
	ID              string               `json:"ID" bson:"id" validate:"required"`
	Name            string               `json:"Name" bson:"name" validate:"required"`
	Formula         []FormulaItemRequest `json:"Formula" bson:"formula" validate:"required"`
	FormulaAsString string               `json:"FormulaAsString" bson:"formulaAsString" validate:"required"`
	FormulaAsQuery  string               `json:"FormulaAsQuery" bson:"formulaAsQuery" validate:"required"`
	Expert          Expert               `json:"Expert" bson:"expert" validate:"required"`
	CreatedAt       string               `json:"CreatedAt" bson:"createdAt" validate:"required"`
	UpdatedAt       string               `json:"UpdatedAt" bson:"updatedAt" validate:"required"`
}

type MetricReference struct {
	Name            string
	MeasurementUnit string
	Url             string
}

type MetricItem struct {
	ID          string `json:"ID" bson:"id" validate:"required"`
	Name        string `json:"Name" bson:"name" validate:"required"`
	Description string `json:"Description" bson:"description" validate:"required"`
}

type MetricBindingRequest struct {
	Blockchain   string                    `json:"Blockchain" bson:"blockchain" validate:"required"`
	TenantID     string                    `json:"TenantID" bson:"tenantID" validate:"required"`
	TenantPK     string                    `json:"TenantPK" bson:"tenantPK" validate:"required"`
	Activity     ActivityForMetricBinding  `json:"Activity" bson:"activity" validate:"required"`
	Metric       MetricForMetricBinding    `json:"Metric" bson:"metric" validate:"required"`
	NoOfFormulas int                       `json:"NoOfFormulas" bson:"noOfFormulas" validate:"required"`
	Formula      []FormulaForMetricBinding `json:"Formula" bson:"formula" validate:"required"`
	CreatedAt    string                    `json:"CreatedAt" bson:"createdAt" validate:"required"`
}

type ActivityForMetricBinding struct {
	ID              string `json:"ID" bson:"id" validate:"required"`
	Name            string `json:"Name" bson:"name" validate:"required"`
	StageId         string `json:"StageId" bson:"stageId" validate:"required"`
	TracifiedItemId string `json:"TracifiedItemId" bson:"tracifiedItemId" validate:"required"`
}

type MetricForMetricBinding struct {
	ID   string `json:"ID" bson:"id" validate:"required"`
	Name string `json:"Name" bson:"name" validate:"required"`
}

type FormulaForMetricBinding struct {
	FormulaId     string
	NoOfVariables int
	ExpertId      string
	Variable      []VariableStructure
}

// Value(float) and the Binding type(int) cannot be 0 in this library
type VariableStructure struct {
	ID              string `json:"ID" bson:"id" validate:"required"`
	Value           float64
	Name            string `json:"Name" bson:"name" validate:"required"`
	Description     string `json:"Description" bson:"description" validate:"required"`
	Key             string `json:"Key" bson:"key" validate:"required"`
	Type            string `json:"Type" bson:"type" validate:"required"`
	MeasurementUnit string `json:"MeasurementUnit" bson:"measurementUnit" validate:"required"`
	Precision       uint16 `json:"Precision" bson:"precision" validate:"required"`
	BindingType     int
	BindData        BindData
}

type BindData struct {
	Stage  Stage
	Master Master
}

type Stage struct {
	WorkflowId    string
	StageId       string
	StageName     string
	FieldId       string
	FieldName     string
	FieldDataType string
}

type Master struct {
	MetaDataName    string
	PrimaryKeyName  string
	KeyValue        string
	KeyDataType     string
	ValueColumnName string
	ValueDataType   string
}

type ThrottlerRecord struct {
	RequestEntityType string
	RequestEntity     string
	Timestamp         time.Time
	FormulaID         string
}

type API_ThrottlerRequest struct {
	RequestEntityType string
	RequestEntity     string
	FormulaID         string
	AllowedAmount     int
	FromTime          time.Time
	ToTime            time.Time
}

type RequestCount struct {
	TotalRequestCount int64
}

// structs for storing expert formula in the DB -> FormulaIdentity, AuthorIdentity, ValueDefinition, ExpertFormulaManageData, ExpertFormulaTransaction, FormulaStore
type FormulaIdentity struct {
	ManageDataName  string
	FormulaMapID    uint64
	ManageDataKey   string
	ManageDataValue []byte
}

type AuthorIdentity struct {
	ManageDataName  string
	AuthorMapID     uint64
	ManageDataKey   string
	ManageDataValue []byte
}

type ValueDefinition struct {
	ValueType         string
	ValueMapID        uint64
	UnitMapID         uint16
	Precision         uint16
	Value             float64
	MetricReferenceID string
	ManageDataKey     string
	ManageDataValue   []byte
}

type Template struct {
	manageDataKey   string
	manageDataValue string
}

type FormulaManageData struct {
	FormulaIdentity   FormulaIdentity
	AuthorIdentity    AuthorIdentity
	ValueDefinitions  []ValueDefinition
	ExecutionTemplate []Template
}

type FormulaTransaction struct {
	TransactionHash   string
	TransactionStatus string
	Memo              []byte
	ManageData        FormulaManageData
	TransactionTime   string
	Cost              string
}

type FormulaStore struct {
	Blockchain             string
	FormulaID              string
	ExpertPK               string
	VariableCount          int
	FormulaJsonRequestBody FormulaBuildingRequest
	Transactions           []FormulaTransaction
	OverflowAmount         int
	Status                 string
	CreatedAt              string
}

type ValueDefOutParmas struct {
	ValueMapID uint64
	UnitMapID  uint16
}

type GeneralValueDefBuildRequest struct {
	ResourceType string
	ResourceName string
	Key          string
	VariableUUID string
	VariableName string
	BindingType  string
	ResourceID   string
	FormulaID    string
}

type ResourceIdMap struct {
	ResourceType string
	ResourceID   string
	MapID        uint64
}

type MetricDataBindingRequest struct {
	ID             string `json:"ID" bson:"id" validate:"required"`
	Blockchain     string `json:"Blockchain" bson:"blockchain" validate:"required"`
	UserID         string `json:"UserID" bson:"userid"`
	UserPK         string `json:"UserPK" bson:"userpk"`
	Name           string `json:"Name" bson:"name" validate:"required"`
	Description    string `json:"Description" bson:"description" validate:"required"`
	BenchmarkRef   string `json:"BenchmarkRef" bson:"benchmarkRef" validate:"required"`
	BenchmarkValue string `json:"BenchmarkValue" bson:"benchmarkValue" validate:"required"`
	BenchmarkUnit  string `json:"BenchmarkUnit" bson:"benchmarkUnit" validate:"required"`
	Status         string `json:"Status" bson:"status" validate:"required"`
	TenantId       string `json:"TenantId" bson:"tenantId" validate:"required"`
	CreatedAt      string `json:"CreatedAt" bson:"createdAt" validate:"required"`
	UpdatedAt      string `json:"UpdatedAt" bson:"updatedAt" validate:"required"`
	Activity       []MetricDataBindArtifactRequest
}

type MetricDataBindArtifactRequest struct {
	ID            string        `json:"ID" bson:"id" validate:"required"`
	Name          string        `json:"Name" bson:"name" validate:"required"`
	StageID       string        `json:"StageID" bson:"stageId" validate:"required"`
	MetricID      string        `json:"MetricID" bson:"metricId" validate:"required"`
	MetricFormula MetricFormula `json:"MetricFormula" bson:"metricFormula" validate:"required"`
	Revision      int           `json:"Revision" bson:"revision" validate:"required"`
	TenantID      string        `json:"TenantID" bson:"tenantId" validate:"required"`
	CreatedAt     string        `json:"CreatedAt" bson:"createdAt" validate:"required"`
	UpdatedAt     string        `json:"UpdatedAt" bson:"updatedAt" validate:"required"`
}

type MetricFormula struct {
	ID                  string `json:"ID" bson:"id" validate:"required"`
	Formula             []FormulaDetails
	MetricExpertFormula []MetricExpertFormula
	TenantID            string `json:"TenantID" bson:"tenantId" validate:"required"`
	PivotFields         []PivotField
	Active              bool `json:"Active" bson:"active" validate:"required"`
}

type FormulaDetails struct {
	ID                 string `json:"ID" bson:"id" validate:"required"`
	Key                string `json:"Key" bson:"key" validate:"required"`
	Field              string
	ArtifactTemplateID string
	Description        string `json:"Description" bson:"description" validate:"required"`
}

type MetricExpertFormula struct {
	ID              string `json:"ID" bson:"description" validate:"required"`
	Name            string `json:"Name" bson:"description" validate:"required"`
	Formula         []FullFormula
	FormulaAsString string `json:"FormulaAsString" bson:"formulaAsString" validate:"required"`
	FormulaAsQuery  string `json:"FormulaAsQuery" bson:"formulaAsQuery" validate:"required"`
	CreatedAt       string `json:"CreatedAt" bson:"createdAt" validate:"required"`
	UpdatedAt       string `json:"UpdatedAt" bson:"updatedAt" validate:"required"`
}

type FullFormula struct {
	Type  string `json:"Type" bson:"type" validate:"required"`
	Value string
	Name  string `json:"Name" bson:"name" validate:"required"`
	ID    string `json:"ID" bson:"id" validate:"required"`
	Key   string
}

type MetricMapDetails struct {
	MetricID   string
	MetricName string
	MapID      uint64
}


type TenantID struct {
	MetricID   string
	MetricName string
	MapID      uint64
}

type TenentMapDetails struct {
	TenentID string
	MapID    uint64
}
