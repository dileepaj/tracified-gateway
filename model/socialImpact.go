package model

import (
	"time"
)

type BuildFormula struct {
	FormulaID       string
	FormulaType     string
	Name            string
	FieldCount      int32
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
	Precision         int64
	MetricReference   MetricReference
}

type Expert struct {
	ExpertID string `json:"ExpertID" bson:"expertID" validate:"required"`
	ExpertPK string `json:"ExpertPK" bson:"expertPK" validate:"required"`
}

type FormulaIDMap struct {
	FormulaID string
	MapID     int64
}

type ExpertIDMap struct {
	ExpertID  string
	ExpertPK  string
	MapID     int64
	FormulaID string
}

type Counters struct {
	ID            string
	SequenceValue int64
}

type ValueIDMap struct {
	ValueId   string
	ValueType string
	ValueName string
	MapID     int64
}

type UnitIDMap struct {
	Unit  string
	MapID int64
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
	Precision       int64  `json:"Precision" bson:"precision" validate:"required"`
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
	FormulaMapID    string
	ManageDataKey   string
	ManageDataValue string
}

type AuthorIdentity struct {
	AuthorMapID     string
	ManageDataKey   string
	ManageDataValue string
}

type ValueDefinition struct {
	ValueType         string
	ValueMapID        string
	UnitMapID         string
	Precision         int64
	Value             float64
	MetricReferenceID string
	ManageDataKey     string
	ManageDataValue   string
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
	Memo              string
	VariableCount     int
	ManageData        FormulaManageData
	TransactionTime   string
}

type FormulaStore struct {
	FormulaID              string
	ExpertPK               string
	FormulaJsonRequestBody FormulaBuildingRequest
	Transactions           []FormulaTransaction
	OverflowAmount         int
	Status                 string
	CreatedAt              string
}
