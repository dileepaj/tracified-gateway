package apiModel


import (
	"main/model"
)

//poc Success response
type PocSuccess struct {
	Message string
	Chain []model.Current
}
//poc Failure response
type PocFailure struct {
	RootHash string
	Error    string
}

//poe Success response
// type PoeSuccess struct {
// 	DataHash string
// 	TxNHash  string
// 	Message  string
// }

type PoeSuccess struct {
	Message string
	TxNHash string
}

//poe Failure response
type PoeFailure struct {
	RootHash string
	Error    string
}

type InsertSuccess struct {
	Message   string
	TxNHash   string
	ProfileID string
	Type      string
}

type TransactionStruct struct {
	TType             string   `json:"TType"`
	ProfileID         []string `json:"ProfileID"`
	PreviousTDPID     string   `json:"PreviousTDPID"`
	Data              []string `json:"Data"`
	Identifiers       []string `json:"Identifiers"`
	PreviousProfileID string   `json:"PreviousProfileID"`
}
