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
