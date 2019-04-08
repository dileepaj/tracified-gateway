package deprecatedBuilder

import (
	// "github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
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

	object := stellarExecuter.ConcreteSubmitXDR{XDR: AP.XDR}
	// object := stellarExecuter.ConcreteInsertData{Hash: AP.Hash, InsertType: AP.InsertType, PreviousTXNID: AP.PreviousTXNID, ProfileId: AP.ProfileId}

	result := object.SubmitXDR()

	return result
}
