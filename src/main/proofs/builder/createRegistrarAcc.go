package builder

import (
	"main/api/apiModel"
	"main/model"
	"main/proofs/executer/stellarExecuter"
)

type createRegistrarAccInterface interface {
	RegistrarAccount() model.InsertDataResponse
}

type AbstractCreateRegistrar struct {
	RegistrarAccount apiModel.RegistrarAccount
}

func (AP *AbstractCreateRegistrar) CreateRegistrarAcc() string {

	object := stellarExecuter.ConcreteRegistrarAcc{RegistrarAccount: AP.RegistrarAccount}

	result := object.SetupAccount()

	return result
}
