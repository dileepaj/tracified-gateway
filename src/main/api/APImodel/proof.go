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

type InsertSuccess struct {
	Message   string
	TxNHash   string
	ProfileID string
	Type      string
}
