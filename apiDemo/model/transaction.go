package model

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
	MergeBlock           int
	TenantID             string
	StageID              string
}