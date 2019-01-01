package builder

import (
	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
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
