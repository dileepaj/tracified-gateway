package model

//proof of POC

// type Root struct {
// 	Previous []Node
// }

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
