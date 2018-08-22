package model

type Node struct {
	Previous []Node
	Current  Current
}

type Current struct {
	TType string
	TXNID    string
	DataHash	[]string
}

//proof of POE
type RetrievePOE struct {
	Txn    string `json:"txn"`
	BCHash string `json:"bcHash"`
	DBHash string `json:"dbHash"`
	Error  Error  `json:"error"`
}

type POE struct {
	RetrievePOE RetrievePOE
}

//proof of POG
type RetrievePOG struct {
	CurTxn     string `json:"CurTxn"`
	PreTxn     string `json:"PreTxn"`
	Identifier string `json:"Identifier"`
	Error      Error  `json:"error"`
}

type POG struct {
	RetrievePOG RetrievePOG
}

//proof of POC
type RetrievePOC struct {
	Txn    string `json:"txn"`
	BCHash []Current `json:"bcHash"`
	DBHash []Current `json:"dbHash"`
	Error  Error  `json:"error"`
}

type POC struct {
	RetrievePOC RetrievePOC
}
