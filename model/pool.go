package model

type Coin struct {
	Id        string
	FieldName string
	CoinName  string
	Amount    string
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
	BuildPathPayment []BuildPathPayment
	ProductId        string
	ProductIdName    string
	EquationId       string
	TenantId         string
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
	ProductId    string
	ProducutName string
	MetricId     string
	EquationId   string
	TenantId     string
	BuildPools   []BuildPool
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
	Id        string
	CoinName  string
	FieldName string
	Value     string
}
type BatchCoinConvert struct {
	EquationID  string
	TenantId    string
	ProductName string
	ProductID   string
	BatchID     string
	BatchName   string
	StageId     string
	MetrixCoin  MetrixCoin
	UserInputs  []UserInput
}

type BatchAccount struct {
	TenentID       string
	BatchID        string
	EquationID     string
	BatchName      string
	ProductName    string
	StageID        string
	BatchAccountPK string
	BatchAccountSK string
}

type MetricFormulas struct {
	ActivityID      string
	MetricFormulaID string
}

type MetricCoin struct {
	Id          string
	CoinName    string
	FieldName   string
	Description string
}

type FieldAndCoin struct {
	ID            string
	CoinName      string
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
	EquationID           string
	Timestamp            string
	Description          string
	TenantID             string
	MetricID             string
	MetricName           string
	MetrixType           string
	MetricFormulas       []MetricFormulas
	EquatinStringFormate string
	SimpleifedEquation   string
	MetricCoin           MetricCoin
	EquationSubPortion   []EquationSubPortion
}
