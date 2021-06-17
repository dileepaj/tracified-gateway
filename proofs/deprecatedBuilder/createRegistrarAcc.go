package deprecatedBuilder

import (
	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/deprecatedStellarExecuter"

)

type createRegistrarAccInterface interface {
	RegistrarAccount() model.InsertDataResponse
}

type AbstractCreateRegistrar struct {
	RegistrarAccount apiModel.RegistrarAccount
}

func (AP *AbstractCreateRegistrar) CreateRegistrarAcc() string {

	object := deprecatedStellarExecuter.ConcreteRegistrarAcc{RegistrarAccount: AP.RegistrarAccount}

	result := object.SetupAccount()

	return result
}
