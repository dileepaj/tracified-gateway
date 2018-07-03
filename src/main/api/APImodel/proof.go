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
type PoeSuccess struct {
	DataHash string
	TxNHash  string
	Message  string
}

//poe Failure response
type PoeFailure struct {
	RootHash string
	Error    string
}
