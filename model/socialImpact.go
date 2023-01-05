package model

import (
	"time"

	"github.com/stellar/go/txnbuild"
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
	Value             any
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
	MetricExpertFormula ExpertFormula `json:"MetricExpertFormula" bson:"metricExpertFormula"`
	User                User          `json:"User" bson:"user"`
}

type ExpertFormula struct {
	Blockchain      string               `json:"Blockchain" bson:"blockchain"`
	ID              string               `json:"ID" bson:"id" validate:"required"`
	Name            string               `json:"Name" bson:"name" validate:"required"`
	Formula         []FormulaItemRequest `json:"Formula" bson:"formula" validate:"required"`
	FormulaAsString string               `json:"FormulaAsString" bson:"formulaAsString" validate:"required"`
	FormulaAsQuery  string               `json:"FormulaAsQuery" bson:"formulaAsQuery" validate:"required"`
	CreatedAt       string               `json:"CreatedAt" bson:"createdAt" validate:"required"`
	UpdatedAt       string               `json:"UpdatedAt" bson:"updatedAt" validate:"required"`
	CiperText       string               `json:"CiperText" bson:"ciperText"`
}

type User struct {
	ID        string `json:"ID" bson:"id" validate:"required"`
	TenantID  string `json:"TenantID" bson:"tenantID" validate:"required"`
	Publickey string `json:"PublicKey" bson:"publickey" validate:"required"`
}

type MetricReference struct {
	ID              string
	Name            string
	MeasurementUnit string
	Description     string
	Reference       string
}

type MetricItem struct {
	ID          string `json:"ID" bson:"id" validate:"required"`
	Name        string `json:"Name" bson:"name" validate:"required"`
	Description string `json:"Description" bson:"description" validate:"required"`
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
	ManageDataOrder int
	ManageDataName  string
	FormulaMapID    uint64
	ManageDataKey   string
	ManageDataValue []byte
}

type AuthorIdentity struct {
	ManageDataOrder int
	ManageDataName  string
	PublicKeyPart   string
	Expert          ExpertIDMap
	ManageDataKey   string
	ManageDataValue []byte
}

type ValueDefinition struct {
	ManageDataOrder   int
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
	AuthorIdentity    []AuthorIdentity
	ValueDefinitions  []ValueDefinition
	ExecutionTemplate []Template
}
type FormulaTransaction struct {
	TransactionHash   string
	TransactionStatus string
	Memo              []byte
	Manifest          string
	FormulaMapID      uint64
	NoOfVariables     int
	ManageData        FormulaManageData
	TransactionTime   string
	Cost              string
}

type FormulaStore struct {
	MetricExpertFormula ExpertFormula
	User                User
	FormulaID           string
	FormulaMapID        uint64
	VariableCount       int
	ExecutionTemplate   ExecutionTemplate
	TotalNoOfManageData int
	NoOfManageDataInTxn int
	Memo                []byte
	TxnHash             string
	TxnSenderPK         string
	XDR                 string
	SequenceNo          int64
	Status              string
	Timestamp           string
	TransactionTime     string
	TransactionCost     string
	ErrorMessage        string
	TxnUUID             string
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
	BindingType  uint
	ResourceID   string
	FormulaID    uint64
}

type ResourceIdMap struct {
	ResourceType string
	ResourceID   string
	MapID        uint64
}

type MetricDataBindingRequest struct {
	Metric MetricReq
	User   User
}

type MetricBindingStore struct {
	MetricId            string
	MetricMapID         uint64
	Metric              MetricReq
	User                User
	TotalNoOfManageData int
	NoOfManageDataInTxn int
	Memo                []byte
	TxnHash             string
	TxnSenderPK         string
	XDR                 string
	SequenceNo          int64
	Status              string
	Timestamp           string
	ErrorMessage        string
	TxnUUID             string
	TransactionTime     string
	TransactionCost     string
}

type MetricReq struct {
	ID             string                          `json:"ID" bson:"id" validate:"required"`
	Blockchain     string                          `json:"Blockchain" bson:"blockchain"`
	Name           string                          `json:"Name" bson:"name" validate:"required"`
	Description    string                          `json:"Description" bson:"description" validate:"required"`
	BenchmarkRef   string                          `json:"BenchmarkRef" bson:"benchmarkRef"`
	BenchmarkValue string                          `json:"BenchmarkValue" bson:"benchmarkValue"`
	BenchmarkUnit  string                          `json:"BenchmarkUnit" bson:"benchmarkUnit"`
	Status         string                          `json:"Status" bson:"status"`
	TenantId       string                          `json:"TenantId" bson:"tenantId" validate:"required"`
	CreatedAt      string                          `json:"CreatedAt" bson:"createdAt" validate:"required"`
	UpdatedAt      string                          `json:"UpdatedAt" bson:"updatedAt" validate:"required"`
	Activities     []MetricDataBindActivityRequest `json:"Activities" bson:"activities" validate:"required"`
}
type MetricDataBindActivityRequest struct {
	ID                                  string                              `json:"ID" bson:"id" validate:"required"`
	Name                                string                              `json:"Name" bson:"name" validate:"required"`
	StageID                             string                              `json:"StageID" bson:"stageid" validate:"required"`
	Stage                               StageReq                            `json:"Stage" bson:"stage" validate:"required"`
	MetricID                            string                              `json:"MetricID" bson:"metricid" validate:"required"`
	MetricFormula                       MetricFormulaReq                    `json:"MetricFormula" bson:"metricformula" validate:"required"`
	WorkflowID                          string                              `json:"WorkflowID" bson:"workflowid" validate:"required"`
	Revision                            int                                 `json:"Revision" bson:"revision" validate:"required"`
	TenantID                            string                              `json:"TenantID" bson:"tenantid" validate:"required"`
	CreatedAt                           string                              `json:"CreatedAt" bson:"createdat" validate:"required"`
	UpdatedAt                           string                              `json:"UpdatedAt" bson:"updatedat" validate:"required"`
	ActivityFormulaDefinitionManageData ActivityFormulaDefinitionManageData `json:"ActivityFormulaDefinitionManageData" bson:"activityformuladefinitionmanagedata"`
	ActivityNameMangeData               ManageDataActivityName
}

type MetricFormulaReq struct {
	ID                  string `json:"ID" bson:"id" validate:"required"`
	Formula             []FormulaDetails
	MetricExpertFormula MetricExpertFormula `json:"MetricExpertFormula" bson:"metricexpertformula" validate:"required"`
	TenantID            string              `json:"TenantID" bson:"tenantid" validate:"required"`
	PivotFields          []PivotField
	Active              bool `json:"Active" bson:"active"`
}

type StageReq struct {
	StageID string `json:"StageID" bson:"stageid" validate:"required"`
	Name    string `json:"Name" bson:"name" validate:"required"`
}
type TransacionDetailsMetricBinding struct {
	MetricID                    string
	MtericMapId                 uint64
	TenantMapId                 uint32
	MetricName                  string
	MetricNameMangeData         ManageDataMetricName
	NoOfActivityFormula         int
	PublisherIdentityManageData PublisherIdentity
	TotalNumberOfManageData     int
	MaxNumOfManageDatePerHash   int
	TXNHashes                   []TransactionHash
}

type ManageDataMetricName struct {
	ManageDataOrder int
	ManageDataType  string
	MetricName      string
	Name            string
	Value           []byte
}

type ManageDataStageName struct {
	ManageDataOrder int
	ManageDataType  string
	StageName       string
	Name            string
	Value           []byte
}

type ManageDataMasterDetails struct {
	ManageDataOrder      int
	ManageDataType       string
	ArtifactID           string
	TracabilityDataType  uint
	Name                 string
	Value                []byte
	ArtifactTemplateName ManageDataArtifactTemplateName
	ArtifactFieldKey     ManageDataFieldKey
	ArtifactFieldName    ManageDataFieldName
}
type ManageDataKeyName struct {
	ManageDataOrder int
	ManageDataType  string
	KeyName         string
	Name            string
	Value           []byte
}

type ManageDataArtifactTemplateName struct {
	ManageDataOrder      int
	ManageDataType       string
	ArtifactTemplateName string
	Name                 string
	Value                []byte
}

type ManageDataFieldKey struct {
	ManageDataOrder int
	ManageDataType  string
	FieldKey        string
	Name            string
	Value           []byte
}

type ManageDataFieldName struct {
	ManageDataOrder int
	ManageDataType  string
	FieldName       string
	Name            string
	Value           []byte
}

type ManageDataActivityName struct {
	ManageDataOrder int
	ManageDataType  string
	ActivityName    string
	Name            string
	Value           []byte
}

type TransactionHash struct {
	Order int
	Memo  []byte
	Hash  string
}
type MetricDataBindArtifactRequest struct {
	ID                                  string        `json:"ID" bson:"id" validate:"required"`
	Name                                string        `json:"Name" bson:"name" validate:"required"`
	StageID                             string        `json:"StageID" bson:"stageId" validate:"required"`
	MetricID                            string        `json:"MetricID" bson:"metricId" validate:"required"`
	MetricFormula                       MetricFormula `json:"MetricFormula" bson:"metricformula" validate:"required"`
	Revision                            int           `json:"Revision" bson:"revision" validate:"required"`
	TenantID                            string        `json:"TenantID" bson:"tenantId" validate:"required"`
	CreatedAt                           string        `json:"CreatedAt" bson:"createdAt" validate:"required"`
	UpdatedAt                           string        `json:"UpdatedAt" bson:"updatedAt" validate:"required"`
	ActivityFormulaDefinitionManageData ActivityFormulaDefinitionManageData
}

type MetricFormula struct {
	ID                  string           `json:"ID" bson:"id" validate:"required"`
	Formula             []FormulaDetails `json:"Formula" bson:"formula" validate:"required"`
	MetricExpertFormula MetricExpertFormula
	TenantID            string
	PivotField          []PivotField
	Active              bool
}

type FormulaDetails struct {
	ID                 string `json:"ID" bson:"id" validate:"required"`
	Key                string `json:"Key" bson:"key" validate:"required"`
	Field              string
	Type               int
	ArtifactTemplateID string
	ArtifactTemplate   ArtifactTemplate
	Description        string `json:"Description" bson:"description" validate:"required"`
	BindManageData     BindManageData
}

type ArtifactTemplate struct {
	ID        string
	Name      string
	FieldName string
}
type MetricExpertFormula struct {
	ID              string        `json:"ID" bson:"id" validate:"required"`
	Name            string        `json:"Name" bson:"name" validate:"required"`
	Formula         []FullFormula `json:"Formula" bson:"formula" validate:"required"`
	FormulaAsString string        `json:"FormulaAsString" bson:"formulaAsString" validate:"required"`
	FormulaAsQuery  string        `json:"FormulaAsQuery" bson:"formulaAsQuery" validate:"required"`
	CreatedAt       string        `json:"CreatedAt" bson:"createdAt" validate:"required"`
	UpdatedAt       string        `json:"UpdatedAt" bson:"updatedAt" validate:"required"`
}

type FullFormula struct {
	Type  string `json:"Type" bson:"type" validate:"required"`
	Value any
	Name  string `json:"Name" bson:"name"`
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

type ActivityMapDetails struct {
	ActivityID string
	Name       string
	MetricID   string
	StageID    string
	MapID      uint64
}
type ActivityFormulaDefinitionManageData struct {
	ManageDataOrder     int
	ManageDataType      string
	FormulaMapID        uint64
	FormulaID           string
	ActivityMapID       uint64
	ActivityID          string
	NoOfDynamicVariable uint32
	Key                 string
	Value               []byte
}
type BindManageData struct {
	ManageDataOrder int
	ManageDataType  string
	BindData        ValueBuilder
	Key             string
	Value           []byte
	Stage           ManageDataStageName
	KeyName         ManageDataKeyName
	Master          ManageDataMasterDetails
}

type PublisherIdentity struct {
	ManageDataOrder int
	ManageDataType  string
	UserID          string
	PublicKeyHash   string
	ManageDataKey   string
	ManageDataValue []byte
}

type SuccessResponseMetricBinding struct {
	Code     int
	ID       string
	MetricID string
	Message  string
}

type SuccessResponseExpertFormula struct {
	Code      int
	FormulaID string
	Message   string
}

type ValueBuilder struct {
	ValueUUID           string
	WorkflowID          string
	StageID             string
	TracabilityDataType uint
	BindingType         uint
}

type WorkflowMap struct {
	WorkflowID string
	MapID      uint64
}

type ArtifactIDMap struct {
	ArtifactID string
	MapID      uint64
}
type PivotField struct {
	Name               string
	Condition          string
	Value              string
	Field              string
	ArtifactDataId     string
	ArtifactTemplateId string
}

type Activity struct {
	ID              string `json:"ID" bson:"id" validate:"required"`
	Name            string `json:"Name" bson:"name" validate:"required"`
	ProductName     string
	TracifiedItemId string `json:"TracifiedItemId" bson:"tracifiedItemId" validate:"required"`
	StageId         string `json:"StageId" bson:"stageId" validate:"required"`
}

type ExecutionTemplate struct {
	Lst_Commands      []Command
	P_Entity          P_Entity
	S_CodeLine        string
	S_StartVarName    string
	Ul_SpecialCommand uint32
	Ul_type           uint64
	Error             string
}

type P_Entity struct {
	Ul_type int64
	Value   any
}

type Command struct {
	P_Arg                ExecutionTemplate
	S_AdditionalFuncName string
	Ul_CommandType       uint32
}

type SendToQueue struct {
	MetricBinding         MetricBindingStore
	ExpertFormula         FormulaStore
	EthereumExpertFormula EthereumExpertFormula
	EthereumMetricBind    EthereumMetricBind
	Type                  string
	User                  User
	Memo                  []byte
	Status                string
	Operations            []txnbuild.ManageData
	ErrorMessage          string
}

type ContractGeneral struct {
	License                   string
	PragmaLine                string
	ImportCalculationsSol     string
	ContractStart             string
	ResultVariable            string
	MetaDataStructure         string
	ValueDataStructure        string
	VariableStructure         string
	SemanticConstantStructure string
	ReferredConstant          string
	MetadataDeclaration       string
	ResultDeclaration         string
	CalculationObject         string
	ContractEnd               string
}

type EthereumExpertFormula struct {
	FormulaID           string
	FormulaName         string
	ExecutionTemplate   ExecutionTemplate
	MetricExpertFormula ExpertFormula
	VariableCount       int32
	ContractName        string
	TemplateString      string
	BINstring           string
	ABIstring           string
	GOstring            string
	SetterNames		 	[]string
	ContractAddress     string
	Status              string
	Timestamp           string
	TransactionHash     string
	TransactionCost     string
	TransactionTime     string
	TransactionUUID     string
	TransactionSender   string
	User                User
	ErrorMessage        string
}

type EthFormulaIDMap struct {
	FormulaID string
	MapID     uint64
}
type EthereumMetricBind struct {
	MetricID          string
	MetricName        string
	Metric            MetricReq
	ContractName      string
	TemplateString    string
	BINstring         string
	ABIstring         string
	Timestamp         string
	ContractAddress   string
	TransactionHash   string
	TransactionCost   string
	TransactionTime   string
	TransactionUUID   string
	TransactionSender string
	User              User
	ErrorMessage      string
	Status            string
	FormulaIDs		  []string
	ValueIDs		  []string
}

type MetricContractGeneral struct {
	License                   string
	PragmaLine                string
	ContractStart             string
	MetaDataStructure         string
	FormulaStructure          string
	ValueDataStructure        string
	MetadataDeclaration       string
	ValueMap                  string
	FormulaMap                string
	AddValueFunction          string
	AddFormulaFunction        string
	GetFormulaDetailsFunction string
	GetValueDetailsFunction   string
	ContractEnd               string
}
