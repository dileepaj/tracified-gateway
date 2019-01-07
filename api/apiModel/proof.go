package apiModel

import (
	"github.com/dileepaj/tracified-gateway/model"
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
/*todelete*/
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

// type SubmitXDRSuccess struct {
// 	Message   string
// 	TxNHash   string
// 	TdpId	string
// 	Identifier string
// 	Type      string
// 	PublicKey string
// 	Status string

// }
type SubmitXDRSuccess struct {
	Status string 
}

// type SubmitXDRSuccess struct {
// 	Message   string
// 	Txns []model.TransactionCollectionBody

// }

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
	Assets            []string
	Code              string
}

type CreateTrustLine struct {
	Code      string `json:"Code"`
	Limit     string `json:"Limit"`
	Issuerkey string `json:"Issuerkey"`
	Signerkey string `json:"Signerkey"`
}
type SendAssest struct {
	Code          string `json:"Code"`
	Amount        string `json:"Amount"`
	Issuerkey     string `json:"Issuerkey"`
	Reciverkey    string `json:"Reciverkey"`
	Signer        string `json:"Signer"`
	PreviousTXNID string `json:"PreviousTXNID"`
	ProfileID     string `json:"ProfileID"`
	Type          string `json:"Type"`
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
	Code              string `json:"Code"`
	Amount            string `json:"Amount"`
	IssuerKey         string `json:"IssuerKey"`
	Reciverkey        string `json:"Reciverkey"`
	Sender            string `json:"Sender"`
	PreviousTXNID     string `json:"PreviousTXNID"`
	PreviousProfileID string `json:"PreviousProfileID"`
	Identifier        string `json:"Identifier"`
	Type              string `json:"Type"`
}

type ChangeOfCustodyLink struct {
	COCTxn            string `json:"COCTxn"`
	SignerKey         string `json:"SignerKey"`
	Type              string `json:"Type"`
	PreviousTXNID     string `json:"PreviousTXNID"`
	PreviousProfileID string `json:"PreviousProfileID"`
	Identifier        string `json:"Identifier"`
}

type AssetTransfer struct {
	Asset             []model.Asset
	Issuer            string `json:"Issuer"`
	Sender            string `json:"Sender"`
	Reciver           string `json:"Reciver"`
	Type              string `json:"Type"`
	PreviousTXNID     string `json:"PreviousTXNID"`
	PreviousProfileID string `json:"PreviousProfileID"`
	Identifier        string `json:"Identifier"`
}

type POCBody struct {
	Chain []model.Current
}

type RegSuccess struct {
	Message string
	Xdr     string
}

type SendAssetRes struct {
	Txn               string
	To                string
	From              string
	Code              string
	Amount            string
	PreviousTXNID     string
	PreviousProfileID string
	Message           string
}

type COCRes struct {
	TxnXDR            string
	To                string
	From              string
	Code              string
	Amount            string
	PreviousTXNID     string
	PreviousProfileID string
	Message           string
}

type InsertTDP struct {
	Type          string
	PreviousTXNID string
	ProfileID     string
	Identifier    string `json:"Identifier"`
	DataHash      string
}

type TestTDP struct {
	XDR string
	// RawTDP string
}

type TestXDRSubmit struct {
	XDR        string
	Identifier string
	TdpId      string
	PublicKey  string
}

type InsertGenesisStruct struct {
	Type       string `json:"Type"`
	Identifier string `json:"Identifier"`
}

type InsertProfileStruct struct {
	Type              string `json:"Type"`
	PreviousTXNID     string `json:"PreviousTXNID"`
	PreviousProfileID string `json:"PreviousProfileID"`
	Identifier        string `json:"Identifier"`
}

type InsertPOAStruct struct {
	Type          string   `json:"Type"`
	PreviousTXNID string   `json:"PreviousTXNID"`
	ProfileID     string   `json:"ProfileID"`
	Identifier    []string `json:"Identifier"`
}

type InsertPOCertStruct struct {
	Type     string `json:"Type"`
	CertType string `json:"CertType"`
	CertBody string `json:"CertBody"`
	Validity string `json:"Validity"`
	Issued   string `json:"Issued"`
	Expired  string `json:"Expired"`
}

type SplitProfileStruct struct {
	Type              string `json:"Type"`
	PreviousTXNID     string `json:"PreviousTXNID"`
	PreviousProfileID string `json:"PreviousProfileID"`
	Identifier        string `json:"Identifier"`
	SplitIdentifiers  []string
	ProfileID         string
	Assets            []string
	Code              string
	// InsertProfileStruct InsertProfileStruct
}

type MergeProfileStruct struct {
	Type               string `json:"Type"`
	PreviousTXNID      string `json:"PreviousTXNID"`
	PreviousProfileID  string `json:"PreviousProfileID"`
	Identifier         string `json:"Identifier"`
	MergingTXNs        []string
	MergingIdentifiers []string
	ProfileID          string
	Assets             string
	Code               string
	// InsertProfileStruct InsertProfileStruct
}

type TrustlineStruct struct {
	Code      string
	Limit     string
	Issuerkey string
	Signerkey string
}

type POEStruct struct {
	Txn       string
	ProfileID string
	Hash      string
}

type POCStruct struct {
	Txn       string
	ProfileID string
	DBTree    []model.Current
	BCTree    []model.Current
}

type POGStruct struct {
	LastTxn    string
	POGTxn     string
	Identifier string
}

type POCOBJ struct {
	Chain []model.Current
}
