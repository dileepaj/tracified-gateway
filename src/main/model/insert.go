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
