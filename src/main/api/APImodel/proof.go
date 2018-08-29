package apiModel

import (
	"main/model"
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
	ProfileTxn        string
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
	ProfileTxn  string
	GenesisTxn  string
	Identifiers string
	Type        string
}

type SplitSuccess struct {
	Message          string
	TxnHash          string
	PreviousTXNID    string
	Identifier       string
	SplitProfiles    []string
	SplitTXN         []string
	SplitIdentifiers []string
	Type             string
}

type MergeSuccess struct {
	Message            string
	TxnHash            string
	PreviousTXNID      string
	ProfileID          string
	Identifier         string
	MergeTXNs          []string
	MergingIdentifiers []string
	Type               string
}

type TransactionStruct struct {
	TType             string   `json:"TType"`
	ProfileID         []string `json:"ProfileID"`
	PreviousTXNID     []string `json:"PreviousTXNID"`
	Data              string   `json:"Data"`
	Identifiers       []string `json:"Identifiers"`
	Identifier        string
	MergingTXNs       []string
	PreviousProfileID []string `json:"PreviousProfileID"`
}

type CreateTrustLine struct {
	Code      string `json:"Code"`
	Limit     string `json:"Limit"`
	Issuerkey string `json:"Issuerkey"`
	Signerkey string `json:"Signerkey"`
}
type SendAssest struct {
	Code       string `json:"Code"`
	Amount     string `json:"Amount"`
	Issuerkey  string `json:"Issuerkey"`
	Reciverkey string `json:"Reciverkey"`
	Signer     string `json:"Signer"`
}

type AppointRegistrar struct {
	Registrar  string `json:"Registrar"`
	AccountKey string `json:"AccountKey"`
	Weight     string `json:"Weight"`
	Low        string `json:"Low"`
	Medium     string `json:"Medium"`
	High       string `json:"High"`
}

type RegistrarAccount struct {
	SignerKeys []model.SignerKey
	SignerKey  string `json:"SignerKey"`
	Low        string `json:"Low"`
	Medium     string `json:"Medium"`
	High       string `json:"High"`
}

type ChangeOfCustody struct {
	Code       string `json:"Code"`
	Amount     string `json:"Amount"`
	IssuerKey  string `json:"IssuerKey"`
	Reciverkey string `json:"Reciverkey"`
	Sender     string `json:"Sender"`
}

// type AssestTransfer struct {
// 	Code     string `json:"Code"`
// 	Limit    string `json:"Amount"`
// 	Reciver1 string `json:"Reciver1"`
// 	Reciver2 string `json:"Reciver2"`
// }

type AssestTransfer struct {
	Asset             []model.Asset
	Issuer            string `json:"Issuer"`
	Sender            string `json:"Sender"`
	Reciver           string `json:"Reciver"`
	PreviousProfileID string `json:"PreviousProfileID"`
	PreviousTxnHash   string `json:"PreviousTxnHash"`
}

type POCBody struct {
	Chain []model.Current
}

type RegSuccess struct {
	Message string
	Xdr     string
}
