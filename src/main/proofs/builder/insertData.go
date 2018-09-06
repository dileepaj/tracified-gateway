package builder

import (
	"main/api/apiModel"
	"main/model"
	"main/proofs/executer/stellarExecuter"
)

// type InsertData struct{}
type TDPInsertInterface interface {
	InsertDataHash() model.InsertDataResponse
}

type AbstractTDPInsert struct {
	InsertTDP apiModel.InsertTDP
	// Hash          string
	// InsertType    string
	// PreviousTXNID string
	// ProfileId     string
}

func (AP *AbstractTDPInsert) TDPInsert() model.InsertDataResponse {

	object := stellarExecuter.ConcreteInsertData{InsertTDP: AP.InsertTDP}
	// object := stellarExecuter.ConcreteInsertData{Hash: AP.Hash, InsertType: AP.InsertType, PreviousTXNID: AP.PreviousTXNID, ProfileId: AP.ProfileId}

	result := object.InsertDataHash()

	return result
}
