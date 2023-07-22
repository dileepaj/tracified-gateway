package model

type TransactionData struct {
	FOUser        string `json:"foUser"`
	XDR           string `json:"xdr"`
	AccountIssuer string `json:"accountIssuer"`
}

type TransactionPayload struct {
	FOUser string `json:"foUser"`
	XDR    string `json:"xdr"`
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
