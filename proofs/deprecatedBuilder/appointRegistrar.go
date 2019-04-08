package deprecatedBuilder

import (
	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/proofs/executer/deprecatedStellarExecuter"

	"github.com/dileepaj/tracified-gateway/model"
)

type AppointRegistrarInterface interface {
	CreateTrustline() model.InsertDataResponse
}

type AbstractAppointRegistrar struct {
	// Publickey string
	// SignerKey string
	// Weight    uint32
	AppointRegistrar apiModel.AppointRegistrar
}

func (AP *AbstractAppointRegistrar) AppointReg() string {

	object := deprecatedStellarExecuter.ConcreteAppointReg{AppointRegistrar: AP.AppointRegistrar}

	result := object.RegistrarRequest()

	return result
}
