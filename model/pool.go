package model

type Coin struct {
	CoinName string
	Amount   string
}

type BuildPathPayment struct {
	SendingCoin        Coin
	IntermediateCoins  []Coin
	ReceivingCoin      Coin
	BatchAccountPK     string
	BatchAccountSK     string
	CoinIssuerAccontPK string
	PoolId             string
	ProductId          string
	EquationId         string
	TenantId           string
}
type BuildPool struct {
	Coin1               string
	DepositeAmountCoin1 string
	Coin2               string
	DepositeAmountCoin2 string
	Ratio               string
	poolDepositorPK     string
	PoolId              string
	ProductId           string
	EquationId          string
	TenantId            string
}