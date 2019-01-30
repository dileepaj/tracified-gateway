package model

//proof of POE
type RetrievePOE struct {
	TdpId       string `json:"txnId"`
	BCHash    string `json:"bcHash"`
	DBHash    string `json:"dbHash"`
	Identifier string `json:"identifier"`
	Error     Error  `json:"error"`
}

type POE struct {
	RetrievePOE RetrievePOE
}

//proof of POG
type RetrievePOG struct {
	CurTxn     string `json:"genesisTxn"`
	PreTxn     string `json:"previousTxn"`
	Identifier string `json:"identifier"`
	Message      Error  `json:"status"`
}

type POG struct {
	RetrievePOG RetrievePOG
}

//proof of POC

type Current struct {
	TType             string   `json:"TType"`
	TXNID             string   `json:"TXNID"`
	DataHash          string   `json:"DataHash"`
	MergedID          string   `json:"MergedID"`
	Identifier        string   `json:"Identifier"`
	PreviousProfileID string   `json:"PreviousProfileID"`
	ProfileID         string   `json:"ProfileID"`
	Assets            []string `json:"Assets"`
	Time              string   `json:"Time"`
	MergedChain       []Current
}

type RetrievePOC struct {
	Txn    string    `json:"txn"`
	BCHash []Current `json:"bcHash"`
	DBHash []Current `json:"dbHash"`
	Error  Error     `json:"error"`
}

type RetrievePrevious struct {
	// Txn    string    `json:"txn"`
	// BCHash []Current `json:"bcHash"`
	HashList []Current `json:"dbHash"`
	// Error  Error     `json:"error"`
}

type POC struct {
	RetrievePOC RetrievePOC
}
