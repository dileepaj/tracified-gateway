package apiModel

import (
	"github.com/Tracified-Gateway/model"
)

//poc Success response
type PocSuccess struct {
	Message string
	Chain   []model.Current
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

type ProfileSuccess struct {
	Message           string
	TxNHash           string
	PreviousTXNID     string
	PreviousProfileID string
	Identifiers       string
	Type              string
}

type InsertSuccess struct {
	Message   string
	TxNHash   string
	ProfileID string
	Type      string
}

type GenesisSuccess struct {
	Message     string
	TxnHash     string
	GenesisTxn  string
	Identifiers string
	Type        string
}

type SplitSuccess struct {
	Message       string
	TxnHash       string
	PreviousTXNID string
	ProfileID     string
	SplitProfiles []string
	Identifiers   []string
	Type          string
}

type MergeSuccess struct {
	Message       string
	TxnHash       string
	PreviousTXNID string
	ProfileID     string
	Identifier   string
	MergingTXNs []string
	MergingIdentifiers []string
	Type          string
}

type TransactionStruct struct {
	TType             string   `json:"TType"`
	ProfileID         []string `json:"ProfileID"`
	PreviousTXNID     []string `json:"PreviousTXNID"`
	Data              []string `json:"Data"`
	Identifiers       []string `json:"Identifiers"`
	Identifier string
	MergingTXNs []string
	PreviousProfileID []string `json:"PreviousProfileID"`
}

type GenesisStruct struct{
	TType string   `json:"TType"`

}

type ProfileStruct struct{
	TType             string   `json:"TType"`

}

type MergeStruct struct{
	TType             string   `json:"TType"`
	ProfileID         []string `json:"ProfileID"`
	PreviousTXNID     []string `json:"PreviousTXNID"`
	Identifiers       []string `json:"Identifiers"`
	Identifier string
	MergingTXNs []string
	PreviousProfileID []string `json:"PreviousProfileID"`

}

type SplitStruct struct{
	TType             string   `json:"TType"`

}

type InsertStruct struct{
	TType             string   `json:"TType"`

}


