package builder

import (
	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
)

// type InsertPOA struct{}
type POCertInsertInterface interface {
	InsertPOCertHash() model.InsertDataResponse
}

type AbstractPOCertInsert struct {
	InsertPOCertStruct apiModel.InsertPOCertStruct
}

func (AP *AbstractPOCertInsert) POCertInsert() model.InsertDataResponse {

	object := stellarExecuter.ConcreteInsertPOCert{InsertPOCertStruct: AP.InsertPOCertStruct}

	result := object.InsertPOCertHash()

	return result
}
