package builder

import (
	// "main/api/apiModel"
	"main/model"
	"main/proofs/executer/stellarExecuter"
)

// type InsertData struct{}
type TDPInsertInterface interface {
	InsertDataHash() model.InsertDataResponse
}

type AbstractTDPInsert struct {
	XDR string
	// InsertTDP apiModel.TestTDP
	// Hash          string
	// InsertType    string
	// PreviousTXNID string
	// ProfileId     string
}

func (AP *AbstractTDPInsert) TDPInsert() model.SubmitXDRResponse {

	object := stellarExecuter.ConcreteInsertData{XDR: AP.XDR}
	// object := stellarExecuter.ConcreteInsertData{Hash: AP.Hash, InsertType: AP.InsertType, PreviousTXNID: AP.PreviousTXNID, ProfileId: AP.ProfileId}

	result := object.InsertDataHash()

	return result
}
