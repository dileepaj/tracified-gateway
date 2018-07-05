package model

//proof of POC

// type Root struct {
// 	Previous []Node
// }

type Node struct {
	Previous []Node  `json:"Previous"`
	Current  Current `json:"Current"`
}

type Current struct {
	TDPID    string `json:"TDPID"`
	DataHash string `json:"DataHash"`
	TxNHash  string `json:"TxNHash"`
	Sequence string `json:"previous"`
}

//poc Success response
type PocSuccess struct {
	RootHash string
	Error    Error1
}

//proof of POE
type POE struct {
	TxNHash  string
	RootHash string
	Error    Error1
}
