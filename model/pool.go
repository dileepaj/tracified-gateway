package model

type Coin struct {
	Id            string
	FieldName     string
	CoinName      string
	GeneratedName string
	Amount        string
}

type BuildPathPayment struct {
	SendingCoin        Coin
	ReceivingCoin      Coin
	BatchAccountPK     string
	BatchAccountSK     string
	CoinIssuerAccontPK string
	PoolId             string
	Hash               string
}

type BuildPathPaymentJSon struct {
	CoinConertions []BuildPathPayment
	ProductId      string
	ProductIdName  string
	EquationId     string
	TenantId       string
}

type BuildPool struct {
	Coin1               string
	DepositeAmountCoin1 string
	Coin2               string
	DepositeAmountCoin2 string
	Ratio               string
	poolDepositorPK     string
	PoolId              string
	EquationId          string
	PoolDepositeHash    string
}

type BuildPoolResponse struct {
	ProductId   string
	ProductName string
	MetricId    string
	EquationId  string
	TenantId    string
	CoinMap     []CoinMap
	BuildPools  []BuildPool
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

type TransactionHashWithIdentifier struct {
	Status          string
	Txnhash         string
	TxnType         string
	Identifier      string
	FromIdentifier1 string
	FromIdentifier2 string
	ToIdentifier    string
	AvailableProof  []string
	ProductName     string
	ProductID       string
}

type TrustlineHistory struct {
	CoinIssuer   string
	CoinReceiver string
	Asset        string
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
	Id            string
	CoinName      string
	GeneratedName string
	FieldName     string
	Value         string
}
type BatchCoinConvert struct {
	EquationID  string      `json:"EquationId" bson:"EquationId" validate:"required"`
	TenantID    string      `json:"TenantId" bson:"TenantId" validate:"required"`
	ProductName string      `json:"ProductName" bson:"ProductName" validate:"required"`
	ProductID   string      `json:"ProductID" bson:"ProductID" validate:"required"`
	BatchID     string      `json:"BatchId" bson:"BatchId" validate:"required"`
	BatchName   string      `json:"BatchName" bson:"BatchName" validate:"required"`
	StageId     string      `json:"stageId" bson:"stageId" validate:"required"`
	MetricCoin  MetricCoin  `json:"metricCoin" bson:"metricCoin" validate:"required"`
	UserInputs  []UserInput `json:"userInputs" bson:"userInputs"`
}

type BatchAccount struct {
	TenantID       string
	BatchID        string
	EquationID     string
	BatchName      string
	ProductID      string
	ProductName    string
	MetricCoin     MetricCoin
	StageID        string
	BatchAccountPK string
	BatchAccountSK string
}

type MetricFormulas struct {
	ActivityID      string
	MetricFormulaID string
}

type MetricCoin struct {
	Id            string
	CoinName      string
	GeneratedName string
	FieldName     string
	Description   string
}

type FieldAndCoin struct {
	ID            string
	CoinName      string
	GeneratedName string
	FieldName     string
	Description   string
	UserInputType string
	VariableType  string
	Value         string
}

type EquationSubPortion struct {
	Subportion   string
	FieldAndCoin []FieldAndCoin
}

type CreatePool struct {
	EquationID           string               `json:"equationId" bson:"equationId" validate:"required"`
	ProductName          string               `json:"productName" bson:"productname" validate:"required"`
	ProductID            string               `json:"productID" bson:"productid" validate:"required"`
	Timestamp            string               `json:"timestamp" bson:"timestamp" validate:"required"`
	Description          string               `json:"description" bson:"description" validate:"required"`
	TenantID             string               `json:"tenantId" bson:"tenantId" validate:"required"`
	MetricID             string               `json:"metricId" bson:"metricId" validate:"required"`
	MetricName           string               `json:"metricName" bson:"metricName" validate:"required"`
	MetrixType           string               `json:"metrixType" bson:"metrixType" validate:"required"`
	ActivityId           string               `json:"activityId" bson:"activityId" validate:"required"`
	MetricFormulaId      string               `json:"metricFormulaId" bson:"metricFormulaId" validate:"required"`
	EquationType         string               `json:"equationType" bson:"equationType" validate:"required"`
	EquatinStringFormate string               `json:"equatinStringFormate" bson:"equatinStringFormate" validate:"required"`
	SimpleifedEquation   string               `json:"simpleifedEquation" bson:"simpleifedEquation" validate:"required"`
	MetricCoin           MetricCoin           `json:"metricCoin" bson:"metricCoin"`
	EquationSubPortion   []EquationSubPortion `json:"equationSubPortion" bson:"equationSubPortion"`
}

type CalculateEquationForBatch struct {
	TenantID    string
	ProductName string
	ProductID   string
	BatchID     string
	BatchName   string
	StageId     string
	EquationID  string
	MetrixType  string
}
type EquationResultForBatch struct {
	TenantID       string
	ProductName    string
	ProductID      string
	BatchID        string
	BatchName      string
	StageId        string
	EquationID     string
	MetrixType     string
	MetricCoin     MetricCoin
	BatchAccount   string
	EquationResult string
}

type CoinMap struct {
	CoinName      string
	GeneratedName string
}

type SendToQueue struct {
	Type              string
	EqationJson       CreatePool
	CoinMap           []CoinMap
	PoolCreationArray []BuildPool
	CoinConvert       BatchCoinConvert
}

type CoinName struct {
	TenantID          string
	EquationID        string
	Type              string
	CoinName          string
	GeneratedCoinName string
	FieldName         string
	Description       string
	Count             string
}

type Pool struct {
	EquationId          string
	ProductId           string
	TenantId            string
	FormulatType        string
	Coin1               string
	DepositeAmountCoin1 string
	Coin2               string
	DepositeAmountCoin2 string
	Ratio               string
}
