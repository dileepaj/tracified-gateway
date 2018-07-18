package model

type Node struct {
	Previous []Node
	Current  Current
	Sequence string
}

type Current struct {
	TDPID    string
	DataHash string
	TxNHash  string
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

//proof of POC
type RetrievePOC struct {
	Txn    string `json:"txn"`
	BCHash string `json:"bcHash"`
	DBHash string `json:"dbHash"`
	Error  Error  `json:"error"`
}

type POC struct {
	RetrievePOC RetrievePOC
}
