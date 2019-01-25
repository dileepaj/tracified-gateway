package builder

import (
	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
)

// type InsertPOA struct{}
type POAInsertInterface interface {
	InsertPOAHash() model.InsertDataResponse
}

type AbstractPOAInsert struct {
	InsertPOAStruct apiModel.InsertPOAStruct
}

func (AP *AbstractPOAInsert) POAInsert() model.InsertDataResponse {

	object := stellarExecuter.ConcreteInsertPOA{InsertPOAStruct: AP.InsertPOAStruct}

	result := object.InsertPOAHash()

	return result
}
