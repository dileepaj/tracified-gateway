package builder

import (
	"main/model"
	"main/proofs/executer/stellarExecuter"
)

// type InsertData struct{}
type TDPInsertInterface interface {
	InsertDataHash() model.InsertDataResponse
}

type AbstractTDPInsert struct {
	Hash          string
	InsertType    string
	PreviousTDPID string
	ProfileId     string
}

func (AP *AbstractTDPInsert) TDPInsert() model.InsertDataResponse {

	object := stellarExecuter.ConcreteInsertData{Hash: AP.Hash, InsertType: AP.InsertType, PreviousTDPID: AP.PreviousTDPID, ProfileId: AP.ProfileId}

	result := object.InsertDataHash()

	return result
}
