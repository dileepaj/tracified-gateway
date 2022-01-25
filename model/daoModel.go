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

type TotalTransaction struct {
	TotalTransactionCount int64
}

type NFTWithTransaction struct {
	Identifier                       string
	TDPTxnHash                       string
	TDPID                            string
	TxnType                          string
	DataHash                         string
	NftTransactionExistingBlockchain string
	NftIssuingBlockchain             string
	NFTTXNhash                       string
	Timestamp                        string
	NftAssetName                     string
	NftContentName                   string
	NftContent                       string
	CuurentIssuerPK                  string
	MainAccountPK					 string
	InitialDistributorPublickKey     string
	InitialIssuerPK                  string
	ProductName                      string
	TrustLineCreatedAt				 string
}

type NFTKeys struct {
	PublicKey string
	SecretKey string
}

// type CurrentIssuer struct {
// 	PublicKey string
// }

type NFTKeyItems struct {
	NFTKeyItem []NFTKeys
}

type MarketPlaceNFT struct {
	Identifier                       string
	TDPTxnHash                       string
	TDPID                            string
	TxnType                          string
	DataHash                         string
	NftTransactionExistingBlockchain string
	NftIssuingBlockchain             string
	NFTTXNhash                       string
	Timestamp                        string
	NftAssetName                     string
	NftContentName                   string
	NftContent                       string
	TrustLineCreatedAt               string
	ProductName                      string
	OriginPK                         string
	SellingStatus                    string
	Amount                           string
	Price                            string
	InitialDistributorPK           	 string
	InitialIssuerPK                  string
	MainAccountPK                    string
	PreviousOwnerNFTPK               string
	CurrentOwnerNFTPK                string
}

type TrustLineResponseNFT struct {
	DistributorPublickKey string
	IssuerPublicKey       string
	Asset_code            string
	TDPtxnhash            string
	TDPID                 string
	NFTBlockChain         string
	Successfull           bool
	TrustLineCreatedAt    string
	Identifier            string
	ProductName           string
}

type NFTCreactedResponse struct {
	NFTTxnHash         string
	TDPTxnHash         string
	NFTName            string
	NFTIssuerPublicKey string
}

type NFTIssuerAccount struct {
	NFTIssuerPK string
}

type TransactionCollectionBody struct {
	Identifier           string
	TdpId                string
	SequenceNo           int64
	ProfileID            string
	TxnHash              string
	PreviousTxnHash      string
	FromIdentifier1      string
	FromIdentifier2      string
	ToIdentifier         string
	ItemCode             string
	ItemAmount           string
	PublicKey            string
	TxnType              string
	XDR                  string
	Status               string
	MergeID              string
	Orphan               bool
	PreviousStage        string
	CurrentStage         string
	AppAccount           string
	DataHash             string
	ProductName          string
	ProductID            string
	PreviousSplitProfile string
	CurrentTxnHash       string
	PreviousTxnHash2     string
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
	Blockchain     string
	Txnhash        string
	TxnType        string
	SequenceNo     int64
	Url            string
	LabUrl         string
	From           string
	To             string
	SourceAccount  string
	Identifier     string
	TdpId          string
	Timestamp      string
	Ledger         string
	FeePaid        string
	AvailableProof []string
	DataHash       string
	ProductName    string
	Itemcount      string
	AssetCode      string
}

type POCOCResponse struct {
	Status         string
	Txnhash        string
	TxnType        string
	SequenceNo     int64
	Url            string
	LabUrl         string
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
	Ledger         string
	FeePaid        string
}

type POGResponse struct {
	Status         string
	Txnhash        string
	TxnType        string
	SequenceNo     int64
	Url            string
	LabUrl         string
	SourceAccount  string
	Identifier     string
	BlockchainName string
	Timestamp      string
	Ledger         string
	FeePaid        string
	ProductName    string
	ProductId      string
}

type POEResponse struct {
	Status         string
	Txnhash        string
	TxnType        string
	SequenceNo     int64
	Url            string
	LabUrl         string
	SourceAccount  string
	Identifier     string
	BlockchainName string
	Timestamp      string
	Ledger         string
	FeePaid        string
	DbHash         string
	BcHash         string
}

type POCResponse struct {
	Status         string
	Txnhash        string
	TxnType        string
	SequenceNo     int64
	Identifier     string
	DataHash       string
	BlockchainName string
	Timestamp      string
	Ledger         string
	FeePaid        string
	Url            string
	SourceAccount  string
	AvailableProof []string
	COCStatus      string
	Quantity       string
	AssetCode      string
	From           string
	To             string
	FromSigned     bool
	ToSigned       bool
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

type ArtifactTransaction struct {
	TenantId       string `json:"tenantId"`
	ArtifactId     string `json:"artifactId"`
	ArtifactDataId string `json:"artifactDataId"`
	CellId         string `json:"cellId"`
	PublicKey      string `json:"publicKey"`
	XDR            string `json:"XDR"`
	Txnhash        string `json:"Txnhash"`
}
type TestimonialOrganization struct {
	Name           string
	Description    string
	Logo           string
	Email          string
	Phone          string
	PhoneSecondary string
	AcceptTxn      string
	AcceptXDR      string
	RejectTxn      string
	RejectXDR      string
	TxnHash        string
	Author         string
	SubAccount     string
	SequenceNo     string
	Status         string
	ApprovedBy     string
	ApprovedOn     string
}

type TestimonialOrganizationResponse struct {
	AcceptTxn  string
	AcceptXDR  string
	RejectTxn  string
	RejectXDR  string
	TxnHash    string
	SequenceNo string
	Status     string
}

type RawData struct {
	Name        string
	Title       string
	Description string
	Image       string
}

type Testimonial struct {
	Sender      string
	Reciever    string
	AcceptTxn   string
	AcceptXDR   string
	RejectTxn   string
	RejectXDR   string
	TxnHash     string
	Subaccount  string
	SequenceNo  string
	Status      string
	Testimonial RawData
}

type TestimonialResponse struct {
	TxnHash     string
	SequenceNo  string
	Status      string
	Testimonial RawData
}

type Option int

const (
	Undefined Option = iota
	Approved
	Rejected
	Expired
	Pending
)

func (O Option) String() string {
	return [...]string{"undefined", "approved", "rejected", "expired", "pending"}[O]
}

type TransactionCollectionBodyWithCount struct {
	Count        int64
	Transactions []TransactionCollectionBody
}
type MarketPlaceNFTTrasactionWithCount struct {
	Count               int64
	MarketPlaceNFTItems []MarketPlaceNFT
}
type Services struct {
	ServiceName       string
	ServiceURL        string
	DocumentationLink string
	//UIAutomationSteps	UIAutomationSteps
}

type ActionParams struct {
	InputVariables  []string
	OutputVariables []string
	Services        Services
}

type Action struct {
	ActionType       string
	ActionParameters ActionParams
}

type Steps struct {
	StepNo      int32
	SegmentNo   string
	Discription string
	Action      Action
}

type ProofProtocol struct {
	ProofName            string
	ProofDescriptiveName string
	NumberofSteps        string
	//Segmants				Segmants
	Steps []Steps
}
