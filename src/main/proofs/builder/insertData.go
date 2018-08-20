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
	PreviousTXNID string
	ProfileId     string
}

func (AP *AbstractTDPInsert) TDPInsert() model.InsertDataResponse {

	object := stellarExecuter.ConcreteInsertData{Hash: AP.Hash, InsertType: AP.InsertType, PreviousTXNID: AP.PreviousTXNID, ProfileId: AP.ProfileId}

	result := object.InsertDataHash()

	return result
}
