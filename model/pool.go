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
	FormulaType    string
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
	FormulatType        string
}

type BuildPoolResponse struct {
	EquationId           string
	ProductId            string
	ProductName          string
	TenantId             string
	ActivityId           string
	EquationStringFormat string
	SimplifiedEquation   string
	FormulaType          string
	CoinMap              []CoinMap
	BuildPools           []BuildPool
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
	EquationID      string      `json:"EquationId" bson:"EquationId" validate:"required"`
	TenantID        string      `json:"TenantId" bson:"TenantId" validate:"required"`
	ProductName     string      `json:"ProductName" bson:"ProductName" validate:"required"`
	ProductID       string      `json:"ProductID" bson:"ProductID" validate:"required"`
	FormulaType     string      `json:"FormulaType" bson:"FormulaType" validate:"required"`
	FormulaTypeID   string      `json:"FormulaTypeID" bson:"FormulaTypeID" validate:"required"`
	FormulaTypeName string      `json:"FormulaTypeName" bson:"FormulaTypeName" validate:"required"`
	StageId         string      `json:"stageId" bson:"stageId" validate:"required"`
	MetricCoin      MetricCoin  `json:"metricCoin" bson:"metricCoin" validate:"required"`
	UserInputs      []UserInput `json:"userInputs" bson:"userInputs"`
}

type ArtifactCoinConvert struct {
	EquationID      string      `json:"EquationId" bson:"EquationId" validate:"required"`
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

type CoinAccount struct {
	TenantID        string
	FormulaType     string
	FormulaTypeID   string
	FormulaTypeName string
	EquationID      string
	ProductID       string
	ProductName     string
	MetricCoin      MetricCoin
	StageID         string
	CoinAccountPK   string
	CoinAccountSK   []byte
}

type MetricFormulas struct {
	ActivityID      string
	MetricFormulaID string
}

type MetricCoin struct {
	CoinName      string
	GeneratedName string
	FieldName     string
	Description   string
}

type FieldAndCoin struct {
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
	TenantID             string               `json:"tenantId" bson:"tenantId" validate:"required"`
	ActivityId           string               `json:"activityId" bson:"activityId" validate:"required"`
	FormulaType          string               `json:"FormulaType" bson:"FormulaType" validate:"required"`
	EquatinStringFormate string               `json:"equatinStringFormate" bson:"equatinStringFormate" validate:"required"`
	SimpleifedEquation   string               `json:"simpleifedEquation" bson:"simpleifedEquation" validate:"required"`
	MetricCoin           MetricCoin           `json:"metricCoin" bson:"metricCoin"`
	EquationSubPortion   []EquationSubPortion `json:"equationSubPortion" bson:"equationSubPortion"`
}

type CreatePoolForArtifact struct {
	EquationID           string               `json:"equationId" bson:"equationId" validate:"required"`
	ProductName          string               `json:"productName" bson:"productname"`
	ProductID            string               `json:"productID" bson:"productid"`
	TenantID             string               `json:"tenantId" bson:"tenantId" validate:"required"`
	ActivityId           string               `json:"activityId" bson:"activityId" validate:"required"`
	FormulaType          string               `json:"FormulaType" bson:"FormulaType" validate:"required"`
	EquatinStringFormate string               `json:"equatinStringFormate" bson:"equatinStringFormate" validate:"required"`
	SimpleifedEquation   string               `json:"simpleifedEquation" bson:"simpleifedEquation" validate:"required"`
	MetricCoin           MetricCoin           `json:"metricCoin" bson:"metricCoin"`
	EquationSubPortion   []EquationSubPortion `json:"equationSubPortion" bson:"equationSubPortion"`
}

type CalculateEquationForBatch struct {
	TenantID        string
	ProductName     string
	FormulaType     string
	FormulaTypeID   string
	FormulaTypeName string
	StageId         string
	EquationID      string
	MetrixType      string
}
type EquationResultForBatch struct {
	TenantID        string
	ProductName     string
	FormulaType     string
	FormulaTypeName string
	StageId         string
	EquationID      string
	MetricCoin      MetricCoin
	Account         string
	EquationResult  string
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

// TODO add FORMULA TYPE
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
	ProductName         string
	TenantId            string
	FormulatType        string
	Coin1               string
	DepositeAmountCoin1 string
	Coin2               string
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
