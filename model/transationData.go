package model

type TransactionData struct {
	FOUser          string `json:"foUser"`
	AccountIssuer   string `json:"accountIssuer"`
	Type            string `json:"Type"`
	Identifier      string `json:"Identifier"`
	ProductName     string `json:"productName"`
	ProductID       string `json:"productId"`
	AppAccount      string `json:"appAccount"`
	CurrentStage    string `json:"CurrentStage"`
	TypeName        string `json:"TypeName"`
	TimeStamp       string `json:"Timestamp"`
	DataHash        string `json:"dataHash"`
	FromIdentifier  string `json:"FromIdentifier"`
	ToIdentifiers   string `json:"ToIdentifiers"`
	FromIdentifier1 string `json:"FromIdentifier1"`
	FromIdentifier2 string `json:"FromIdentifier2"`
	PreviousStage   string `json:"PreviousStage"`
}

type TransactionDataKeys struct {
	FOUser          string `json:"foUser"`
	AccountIssuerPK string `json:"accountIssuerpk"`
	AccountIssuerSK []byte `json:"accountIssuersk"`
}

type BalanceResponse struct {
	Balances []struct {
		Balance string `json:"balance"`
	} `json:"balances"`
}

type AccountResponse struct {
	AccountID string `json:"account_id"`
}
