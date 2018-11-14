package model

type RootTree struct {
	Hash  string
	Error Error
}

type InsertDataResponse struct {
	TDPID     string
	ProfileID string
	TxnType   string
	Error     Error
}

type SubmitXDRResponse struct {
	TDPID     string
	Identifier string
	PublicKey   string
	TXNID string
	TxnType string
	Error     Error
}

type InsertProfileResponse struct {
	ProfileTxn        string
	PreviousProfileID string
	PreviousTXNID     string
	Identifiers       string
	TxnType           string
	Error             Error
}

type InsertGenesisResponse struct {
	ProfileTxn  string
	GenesisTxn  string
	Identifiers string
	TxnType     string
	Error       Error
}

type SplitProfileResponse struct {
	Txn               string
	PreviousProfileID string
	PreviousTXNID     string
	Identifiers       string
	SplitProfiles     []string
	SplitTXN          []string
	TxnType           string
	Error             Error
}

type MergeProfileResponse struct {
	Txn                 string
	PreviousProfileID   string
	PreviousTXNID       string
	Identifiers         string
	PreviousIdentifiers []string
	MergeTXNs           []string
	TxnType             string
	Error               Error
	ProfileID           string
}

type SendAssetResponse struct {
	Txn               string
	To                string
	From              string
	Code              string
	Amount            string
	PreviousTXNID     string
	PreviousProfileID string
	Error             Error
}

type COCResponse struct {
	TxnXDR            string
	To                string
	From              string
	Code              string
	Amount            string
	PreviousTXNID     string
	PreviousProfileID string
	Error             Error
}
