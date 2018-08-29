package builder

import (
	"main/api/apiModel"
	"main/model"
	"main/proofs/executer/stellarExecuter"
)

type cocTransactionInterface interface {
	CreateTrustline() model.InsertDataResponse
}

type AbstractCoCTransaction struct {
	// Code       string
	// Amount     string
	// IssuerKey  string
	// Reciverkey string
	// Sender     string
	ChangeOfCustody apiModel.ChangeOfCustody
}

func (AP *AbstractCoCTransaction) CoCTransaction() string {
	object2 := stellarExecuter.ConcreteProfile{Identifiers: AP.ChangeOfCustody.Identifier, InsertType: "1", PreviousTXNID: AP.ChangeOfCustody.PreviousTXNID, PreviousProfileID: AP.ChangeOfCustody.PreviousProfileID}

	result2 := object2.InsertProfile()

	object := stellarExecuter.ConcreteChangeOfCustody{COC: AP.ChangeOfCustody, ProfileId: result2.ProfileTxn}

	result := object.ChangeOfCustody()

	return result
}
