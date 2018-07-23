package apiModel

//poc Success response
type PocSuccess struct {
	RootHash string
	Message  string
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
	TxNHash           string
	PreviousTDPID     string
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
	TxnHash     string
	GenesisTxn  string
	Identifiers string
	Type        string
}

type SplitSuccess struct {
	Message       string
	TxnHash       string
	PreviousTDPID string
	ProfileID     string
	Identifiers   []string
	Type          string
}

type TransactionStruct struct {
	TType             string   `json:"TType"`
	ProfileID         []string `json:"ProfileID"`
	PreviousTDPID     []string `json:"PreviousTDPID"`
	Data              []string `json:"Data"`
	Identifiers       []string `json:"Identifiers"`
	PreviousProfileID []string `json:"PreviousProfileID"`
}
