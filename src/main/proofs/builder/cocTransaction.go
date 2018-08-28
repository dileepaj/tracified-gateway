package builder

import (
	"main/model"
	"main/proofs/executer/stellarExecuter"
)

type cocTransactionInterface interface {
	CreateTrustline() model.InsertDataResponse
}

type AbstractCoCTransaction struct {
	Code       string
	Amount     string
	IssuerKey  string
	Reciverkey string
	Sender     string
}

func (AP *AbstractCoCTransaction) CoCTransaction() string {

	object := stellarExecuter.ConcreteChangeOfCustody{Code: AP.Code, Amount: AP.Amount, IssuerKey: AP.IssuerKey, Reciverkey: AP.Reciverkey, Sender: AP.Sender}

	result := object.ChangeOfCustody()

	return result
}
