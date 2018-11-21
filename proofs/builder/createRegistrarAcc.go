package builder

import (
	"github.com/tracified-gateway/api/apiModel"
	"github.com/tracified-gateway/model"
	"github.com/tracified-gateway/proofs/executer/stellarExecuter"
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
