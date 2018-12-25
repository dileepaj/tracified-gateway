package builder

import (
	"github.com/dileepaj/tracified-gateway/api/apiModel"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
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

	object := stellarExecuter.ConcreteAppointReg{AppointRegistrar: AP.AppointRegistrar}

	result := object.RegistrarRequest()

	return result
}
