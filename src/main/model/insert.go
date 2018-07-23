package model

type RootTree struct {
	Hash  string
	Error Error
}

type InsertDataResponse struct {
	Txn       string
	ProfileID string
	TxnType   string
	Error     Error
}

type InsertProfileResponse struct {
	Txn               string
	PreviousProfileID string
	PreviousTDPID     string
	Identifiers       string
	TxnType           string
	Error             Error
}

type InsertGenesisResponse struct {
	Txn         string
	GenesisTxn  string
	Identifiers string
	TxnType     string
	Error       Error
}

type SplitProfileResponse struct {
	Txn   string
	Error Error
}
