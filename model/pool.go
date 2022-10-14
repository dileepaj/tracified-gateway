package model

import "go.mongodb.org/mongo-driver/bson/primitive"

//! CoinName = first 4 chacter form coin name
//! fullCoinName = user inserted coin name
//! generatedName = gateway generated coin name
type Coin struct {
	ID              string
	FullCoinName    string
	CoinName        string
	GeneratedName   string
	Description     string
	Amount          string
	RescaledAmmount string
}

type BuildPathPayment struct {
	SendingCoin        Coin
	ReceivingCoin      Coin
	BatchAccountPK     string
	BatchAccountSK     string
	CoinIssuerAccontPK string
	PoolId             string
	CoinConversionHash string
	PoolTradeFees      string
	PoolTradeFeesHash  string
}

type BuildPathPaymentResponse struct {
	SendingCoin        Coin
	ReceivingCoin      Coin
	BatchAccountPK     string
	CoinIssuerAccontPK string
	PoolId             string
	Hash               string
}

type BuildPathPaymentJSonResponse struct {
	CoinConertions []BuildPathPaymentResponse
	ProductId      string
	ProductIdName  string
	EquationId     string
	TenantId       string
	FormulaType    string
}

type BuildPathPaymentJSon struct {
	Metric         Metric
	Inputs         []Input
	RealAnswer     float64
	ActualAnswer   float64
	ErrorRate      float64
	CoinConertions []BuildPathPayment
	CoinAccount    CoinAccount
	Event          CoinConvertBody
	AccountPK      string
	CreatedAt      string
}

type BuildPool struct {
	Coin1               string
	Coin1Name           string
	Coin1FullName       string
	DepositeAmountCoin1 string
	Coin2Name           string
	Coin2FullName       string
	Coin2               string
	DepositeAmountCoin2 string
	Ratio               string
	poolDepositorPK     string
	PoolId              string
	EquationId          string
	TenantId            string
	PoolDepositeHash    string
	FormulatType        string
	Activity            []Activity
	MetricCoin          MetricCoin
}

type BuildPoolResponse struct {
	EquationId           string
	Products             []Product
	TenantId             string
	Activity             []Activity
	MetricCoin           MetricCoin
	EquationStringFormat string
	SimplifiedEquation   string
	FormulaType          string
	CoinMap              []CoinMap
	BuildPools           []BuildPool
	EquationSubPortion   []EquationSubPortion
	CreatedAt            string
	UpdatedAt            string
	Active               bool
}

type CoinPath struct {
	Type     string
	CoinName string
	Issuer   string
}

type DestinationCoin struct {
	Destination      Coin
	IntermediateCoin []CoinPath
}

type TrustlineHistory struct {
	CoinIssuer   string
	CoinReceiver string
	Asset        string
	Hash         string
}

type MetrixCoin struct {
	Id          string
	CointName   string
	FieldName   string
	Description string
}

type SendingCoin struct {
	Id       string
	CoinName string
	Amount   string
}

type IntermediateCoin struct {
	Id       string
	CoinName string
	Amount   string
}

type ReceivingCoin struct {
	Id       string
	CoinName string
	Amount   string
}

type PathPayment struct {
	SendingCoin      []SendingCoin
	IntermediateCoin []IntermediateCoin
	ReceivingCoin    []ReceivingCoin
}

type UserInput struct {
	ID            string
	CoinName      string
	GeneratedName string
	FullCoinName  string
	Value         string
}
type BatchCoinConvert struct {
	EquationID       string      `json:"EquationId" bson:"equationId" validate:"required"`
	TenantID         string      `json:"TenantId" bson:"TenantId" validate:"required"`
	ProductName      string      `json:"ProductName" bson:"ProductName"`
	ProductID        string      `json:"ProductID" bson:"ProductID" validate:"required"`
	FormulaType      string      `json:"FormulaType" bson:"FormulaType" validate:"required"`
	FormulaTypeID    string      `json:"FormulaTypeID" bson:"FormulaTypeID" validate:"required"`
	FormulaTypeName  string      `json:"FormulaTypeName" bson:"FormulaTypeName" validate:"required"`
	StageId          string      `json:"stageId" bson:"stageId" validate:"required"`
	MetricCoin       MetricCoin  `json:"metricCoin" bson:"metricCoin" validate:"required"`
	UserInputs       []UserInput `json:"userInputs" bson:"userInputs" validate:"required"`
	Event            string      `json:"Event" bson:"event"`
	EventModel       string      `json:"EventModel" bson:"eventModel"`
	EventModelId     string      `json:"EventModelId" bson:"eventModelId"`
	MetricFormulaId  string      `json:"MetricFormulaId" bson:"metricFormulaId"`
	MetricActivityId string      `json:"MetricActivityId" bson:"metricActivityId"`
	CreatedAt        string      `json:"CreatedAt" bson:"createdAt"`
}

type ArtifactCoinConvert struct {
	EquationID      string      `json:"EquationId" bson:"equationId" validate:"required"`
	TenantID        string      `json:"TenantId" bson:"TenantId" validate:"required"`
	ProductName     string      `json:"ProductName" bson:"ProductName"`
	ProductID       string      `json:"ProductID" bson:"ProductID"`
	FormulaType     string      `json:"FormulaType" bson:"FormulaType" validate:"required"`
	FormulaTypeID   string      `json:"FormulaTypeID" bson:"FormulaTypeID" validate:"required"`
	FormulaTypeName string      `json:"FormulaTypeName" bson:"FormulaTypeName" validate:"required"`
	StageId         string      `json:"stageId" bson:"stageId"`
	MetricCoin      MetricCoin  `json:"metricCoin" bson:"metricCoin" validate:"required"`
	UserInputs      []UserInput `json:"userInputs" bson:"userInputs"`
}

type Formula struct {
	StageID           string
	TracifiedItemId   string
	MetricFormulaId   string
	MetricActivivtyId string
}
type CoinAccount struct {
	Event         Event
	Type          string
	TenantID      string
	CreatedAt     string
	CoinAccountPK string
	CoinAccountSK []byte
}

type MetricFormulas struct {
	ActivityID      string
	MetricFormulaID string
}

type MetricCoin struct {
	ID            string
	Name          string
	CoinName      string
	GeneratedName string
	FullCoinName  string
	Description   string
}

// in coin conversion request name eqaualto full coin name
type Metric struct {
	ID            string
	Name          string
	GeneratedName string
	FullCoinName  string
	CoinName      string
	Description   string
}

type FieldAndCoin struct {
	ID            string
	CoinName      string
	GeneratedName string
	FullCoinName  string
	Description   string
	VariableType  string
	Value         string
}

type EquationSubPortion struct {
	Subportion   string
	FieldAndCoin []FieldAndCoin
}

type Subsets struct {
	SubsetAsString string
	SubSet         []Subset
}

type Subset struct {
	ID                 string
	Type               string
	Value              string
	CoinName           string
	Description        string
	ArtifactTemplateId string
}

type Product struct {
	ProductID   string
	ProductName string
}
type CreatePool struct {
	EquationID           string
	Products             []Product
	TenantID             string
	Activity             []Activity
	FormulaType          string
	EquatinStringFormate string
	SimpleifedEquation   string
	MetricCoin           MetricCoin
	EquationSubPortion   []EquationSubPortion
	CreatedAt            string
	UpdatedAt            string
	Active               bool
}

type CreatePoolBody struct {
	ID              string `json:"ID" bson:"Id" validate:"required"`
	ProductName     string `json:"ProductName" bson:"productname"`
	ProductID       string `json:"ProductID" bson:"productid"`
	Metric          Metric `json:"Metric" bson:"metric"`
	Activity        Activity
	TenantID        string `json:"tenantId" bson:"tenantId" validate:"required"`
	FormulaAsString string `json:"FormulaAsString" bson:"formulaasstring" validate:"required"`
	Subsets         []Subsets
	Active          bool
	CreatedAt       string `json:"CreatedAt" bson:"createdat" validate:"required"`
	UpdatedAt       string `json:"UpdatedAt" bson:"updatedat" validate:"required"`
	PivotFields     []PivotField
}

type Activity struct {
	ID              string `json:"ID" bson:"id" validate:"required"`
	Name            string `json:"Name" bson:"name" validate:"required"`
	ProductName     string
	TracifiedItemId string `json:"TracifiedItemId" bson:"tracifiedItemId" validate:"required"`
	StageId         string `json:"StageId" bson:"stageId" validate:"required"`
}

// type CreatePoolForArtifact struct {
// 	EquationID           string               `json:"equationId" bson:"equationId" validate:"required"`
// 	ProductName          string               `json:"productName" bson:"productname"`
// 	ProductID            string               `json:"productID" bson:"productid"`
// 	TenantID             string               `json:"tenantId" bson:"tenantId" validate:"required"`
// 	ActivityId           string               `json:"activityId" bson:"activityId" validate:"required"`
// 	FormulaType          string               `json:"FormulaType" bson:"FormulaType" validate:"required"`
// 	EquatinStringFormate string               `json:"equatinStringFormate" bson:"equatinStringFormate" validate:"required"`
// 	SimpleifedEquation   string               `json:"simpleifedEquation" bson:"simpleifedEquation" validate:"required"`
// 	MetricCoin           MetricCoin           `json:"metricCoin" bson:"metricCoin"`
// 	EquationSubPortion   []EquationSubPortion `json:"equationSubPortion" bson:"equationSubPortion"`
// }

type CalculateEquation struct {
	MetricFormulaId  string `json:"MetricFormulaId" bson:"metricformulaid" validate:"required"`
	TenantID         string `json:"TenantID" bson:"tenantid" validate:"required"`
	MetricActivityId string `json:"MetricActivityId" bson:"metricactivityid"`
	Type             string `json:"Type" bson:"type" validate:"required"`
	BatchID          string `json:"BatchID" bson:"batchid"`
	ArtifactID       string `json:"artifactID" bson:"artifactid"`
	TracifiedItemId  string `json:"TracifiedItemId" bson:"tracifieditemid"`
}
type EquationResultForBatch struct {
	MetricFormulaId  string
	TenantID         string
	MetricActivityId string
	Type             string
	BatchID          string
	ArtifactID       string
	StageID          string
	TracifiedItemId  string
	Metric           Metric
	Account          string
	EquationResult   string
}

type CoinMap struct {
	ID            string
	CoinName      string
	FullCoinName  string
	GeneratedName string
	Description   string
}

type SendToQueue struct {
	Type              string
	EqationJson       CreatePool
	CoinMap           []CoinMap
	PoolCreationArray []BuildPool
	CoinConvert       CoinConvertBody
}

// TODO add FORMULA TYPE
type CoinName struct {
	TenantID          string
	EquationID        string
	Type              string
	CoinID            string
	CoinName          string
	GeneratedCoinName string
	FullCoinName      string
	Description       string
	Count             string
	MetricID          string
	Timestamp         primitive.DateTime
}

type Pool struct {
	EquationId          string
	Products            []Product
	ProductName         string
	TenantId            string
	FormulatType        string
	Coin1               string
	Coin1Name           string
	Coin1FullName       string
	DepositeAmountCoin1 string
	Coin2               string
	Coin2Name           string
	Coin2FullName       string
	DepositeAmountCoin2 string
	Ratio               string
}

type Sendingcoin struct {
	Id            string
	Fieldname     string
	Coinname      string
	GeneratedName string
	Amount        string
}

type Receivingcoin struct {
	Id            string
	Fieldname     string
	Coinname      string
	GeneratedName string
	Amount        string
}
type Coinconertions struct {
	Sendingcoin        Sendingcoin
	ReceivingCoin      Receivingcoin
	BatchAccountPK     string
	BatchAccountSK     string
	CoinIssuerAccontPK string
	PoolID             string
	Hash               string
}
type CoinConversionDetails struct {
	Coinconertions []Coinconertions
	ProductID      string
	EquationID     string
	TenantID       string
	FormulaType    string
}

type AddStageAndProduct struct {
	EquationID  string
	TenantID    string
	Type        string
	ProductID   string
	ProductName string
	Activity    Activity
}

type Input struct {
	ID            string
	Key           string
	Input         float32
	CoinName      string
	GeneratedName string
	Description   string
}

type Identifier struct {
	Identifier string
	Type       string
}

type Details struct {
	StageID           string
	BatchID           string
	ArtifactID        string
	Identifier        Identifier
	TracifiedItemId   string
	TracifiedItemName string
}
type Event struct {
	Type    string
	ID      string
	Model   string
	Details Details
}
type CoinConvertBody struct {
	ID               string
	Value            float64
	Metric           Metric
	Inputs           []Input
	Event            Event
	Type             string
	MetricFormulaId  string
	MetricActivityId string
	TenantID         string
	CreatedAt        string
}

type PivotField struct {
	Name               string
	Condition          string
	Value              string
	Field              string
	ArtifactDataId     string
	ArtifactTemplateId string
}
