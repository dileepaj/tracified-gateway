package model

import (
	"github.com/stellar/go/build"
)

type COCCollectionBody struct {
	TxnHash    string
	Sender     string
	SubAccount string
	SequenceNo string
	Receiver   string
	AcceptXdr  string
	RejectXdr  string
	AcceptTxn  string
	RejectTxn  string
	Identifier string
	Status     string
}
type TransactionCollectionBody struct {
	Identifier      string
	TdpId           string
	SequenceNo      int64
	ProfileID       string
	TxnHash         string
	PreviousTxnHash string
	FromIdentifier1 string
	FromIdentifier2 string
	ToIdentifier    string
	ItemCode        string
	ItemAmount      string
	PublicKey       string
	TxnType         string
	XDR             string
	Status          string
	MergeID         string
	Orphan          bool
	PreviousStage   string
	CurrentStage    string
	AppAccount      string
	DataHash        string
}

type ProfileCollectionBody struct {
	ProfileTxn         string
	ProfileID          string
	TxnType            string
	PreviousProfileTxn string
	Identifier         string
	TriggerTxn         string
}

type CertificateCollectionBody struct {
	TxnType             string
	PreviousCertificate string
	CertificateType     string
	Data                string
	ValidityPeriod      string
	Asset               string
	PublicKey           string
	XDR                 string
	CertificateID       string
	Status              string
}
type XDR struct {
	XDR build.TransactionMutator
}

type LastTxnResponse struct {
	LastTxn string
}

type TransactionId struct {
	Txnhash string
	Url     string
}

type TransactionIds struct {
	Status     string
	Txnhash    string
	Url        string
	Identifier string
	TdpId      string
}

type PrevTxnResponse struct {
	Status         string
	Txnhash        string
	TxnType        string
	SequenceNo     int64
	Url            string
	From           string
	To             string
	SourceAccount  string
	Identifier     string
	TdpId          string
	Timestamp      string
	Ledger         string
	FeePaid        string
	AvailableProof string
	DataHash       string
}

type POCOCResponse struct {
	Status         string
	Txnhash        string
	Url            string
	Identifier     string
	Quantity       string
	AssetCode      string
	From           string
	To             string
	FromSigned     bool
	ToSigned       bool
	BlockchainName string
	COCStatus      string
	Timestamp      string
}

type POGResponse struct {
	Status         string
	Txnhash        string
	Url            string
	SourceAccount  string
	Identifier     string
	BlockchainName string
	Timestamp      string
}

type COCCollectionList struct {
	List []COCCollectionBody
}
type TransactionCollectionList struct {
	List []TransactionCollectionBody
}
type TransactionUpdate struct {
	Selector TransactionCollectionBody
	Update   TransactionCollectionBody
}
