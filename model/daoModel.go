package model

import (
	"github.com/stellar/go/txnbuild"
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
type TransactionCollectionBody struct {
	Identifier           string
	RealIdentifier       string
	TdpId                string
	SequenceNo           int64
	ProfileID            string
	TxnHash              string
	PreviousTxnHash      string
	FromIdentifier1      string
	FromIdentifier2      string
	ToIdentifier         string
	MapIdentifier1       string
	MapIdentifier2       string
	MapIdentifier        string
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
	XDR txnbuild.Transaction
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
	Status          string
	Blockchain      string
	Txnhash         string
	TxnType         string
	SequenceNo      int64
	Url             string
	LabUrl          string
	From            string
	To              string
	SourceAccount   string
	Identifier      string
	TdpId           string
	Timestamp       string
	Ledger          string
	FeePaid         string
	AvailableProof  []string
	DataHash        string
	ProductName     string
	Itemcount       string
	AssetCode       string
	FromIdentifier1 string
	FromIdentifier2 string
	ToIdentifier    string
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
	Accepted
)

func (O Option) String() string {
	return [...]string{"undefined", "approved", "rejected", "expired", "pending", "accepted"}[O]
}

type TransactionCollectionBodyWithCount struct {
	Count        int64
	Transactions []TransactionCollectionBody
}

type Services struct {
	ServiceName       string
	ServiceURL        string
	DocumentationLink string
	// UIAutomationSteps	UIAutomationSteps
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
	// Segmants				Segmants
	Steps []Steps
}

type TransactionHashWithIdentifier struct {
	Status          string
	Txnhash         string
	TxnType         string
	Identifier      string
	FromIdentifier1 string
	FromIdentifier2 string
	ToIdentifier    string
	AvailableProof  []string
	ProductName     string
	ProductID       string
	Timestamp       string
}

type NFTSolana struct {
	OwnerPK       string
	Asset_code    string
	NFTURL        string
	Description   string
	Collection    string
	NFTBlockChain string
	Tags          string
	Categories    string
	Copies        string
	NFTLinks      string
	ArtistName    string
	ArtistLink    string
}

type NFTWithTransactionSolana struct {
	Identifier                       string
	Categories                       string
	Collection                       string
	ImageBase64                      string `json:"imagebase64"`
	NftTransactionExistingBlockchain string
	NftIssuingBlockchain             string
	NFTTXNhash                       string
	Timestamp                        string
	NftURL                           string
	NftContentName                   string
	NftContent                       string
	MinterPK                         string
	OwnerPK                          string
	NFTArtistName                    string
	NFTArtistURL                     string
	Description                      string
	Copies                           string
	InitialDistributorPK             string
}

type MarketPlaceNFT struct {
	Identifier                       string
	Categories                       string
	Collection                       string
	ImageBase64                      string
	NftTransactionExistingBlockchain string
	NftIssuingBlockchain             string
	NFTTXNhash                       string
	Timestamp                        string
	NftURL                           string
	NftContentName                   string
	NftContent                       string
	NFTArtistName                    string
	NFTArtistURL                     string
	TrustLineCreatedAt               string
	Description                      string
	Copies                           string
	OriginPK                         string
	SellingStatus                    string
	Amount                           string
	Price                            string
	InitialDistributorPK             string
	InitialIssuerPK                  string
	MainAccountPK                    string
	PreviousOwnerNFTPK               string
	CurrentOwnerNFTPK                string
}

type NFTContracts struct {
	NFTContract         string
	MarketplaceContract string
	MintNFTTxn          string
	OwnerPK             string
	Asset_code          string
	NFTURL              string
	Description         string
	Collection          string
	NFTBlockChain       string
	Tags                string
	Categories          string
	Copies              string
	NFTLinks            string
	ArtistName          string
	ArtistLink          string
	Identifier          string
}

type NFTWithTransactionContracts struct {
	Identifier                       string
	Categories                       string
	Collection                       string
	ImageBase64                      string `json:"imagebase64"`
	NftTransactionExistingBlockchain string
	NftIssuingBlockchain             string
	NFTTXNhash                       string
	Timestamp                        string
	NftURL                           string
	NftContentName                   string
	NftContent                       string
	NFTContract                      string
	MarketplaceContract              string
	OwnerPK                          string
	NFTArtistName                    string
	NFTArtistURL                     string
	Description                      string
	Copies                           string
}

type NFTCreactedResponse struct {
	NFTTxnHash         string
	TDPTxnHash         string
	NFTName            string
	NFTIssuerPublicKey string
}

type Minter struct {
	NFTIssuerPK   string `json:"NFTIssuerPK"`
	NFTTxnHash    string `json:"NFTTxnHash"`
	NFTIdentifier string `json:"NFTIdentifier"`
	CreatorUserID string `json:"CreatorUserID"`
}

type NFTWithTransaction struct {
	Identifier                       string
	Categories                       string
	Collection                       string
	ImageBase64                      string
	NftTransactionExistingBlockchain string
	NftIssuingBlockchain             string
	NFTTXNhash                       string
	Timestamp                        string
	NftURL                           string
	NftContentName                   string
	NftContent                       string
	CuurentIssuerPK                  string
	MainAccountPK                    string
	InitialDistributorPublickKey     string
	InitialIssuerPK                  string
	NFTArtistName                    string
	NFTArtistURL                     string
	TrustLineCreatedAt               string
	Description                      string
	Copies                           string
}

type TrustLineResponseNFT struct {
	DistributorPublickKey string
	IssuerPublicKey       string
	Asset_code            string
	NFTURL                string
	Description           string
	Collection            string
	NFTBlockChain         string
	Tags                  string
	Categories            string
	Copies                string
	NFTLinks              string
	ArtistName            string
	ArtistLink            string
	Successfull           bool
	TrustLineCreatedAt    string
}

type MarketPlaceNFTTrasactionWithCount struct {
	Count               int64
	MarketPlaceNFTItems []MarketPlaceNFT
}

type NFTKeys struct {
	PublicKey string
	SecretKey []byte
}

type NFTIssuerAccount struct {
	NFTIssuerPK string
}

type StellarMintTXN struct {
	NFTTxnHash string `json:"NFTTxnHash"`
}

type PublicKey struct {
	PublicKey string `json:"PublicKey"`
}

type XDRRuri struct {
	XDR string
}
