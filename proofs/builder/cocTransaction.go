package builder

import (
	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
)

type cocTransactionInterface interface {
	ChangeOfCustody() model.COCResponse
}

type AbstractCoCTransaction struct {
	// Code       string
	// Amount     string
	// IssuerKey  string
	// Reciverkey string
	// Sender     string
	ChangeOfCustody apiModel.ChangeOfCustody
}

func (AP *AbstractCoCTransaction) CoCTransaction() model.COCResponse {
	// object2 := stellarExecuter.ConcreteProfile{Identifiers: AP.ChangeOfCustody.Identifier, InsertType: "1", PreviousTXNID: AP.ChangeOfCustody.PreviousTXNID, PreviousProfileID: AP.ChangeOfCustody.PreviousProfileID}

	// result2 := object2.InsertProfile()

	object := stellarExecuter.ConcreteChangeOfCustody{COC: AP.ChangeOfCustody}

	result := object.ChangeOfCustody()

	return result
}
