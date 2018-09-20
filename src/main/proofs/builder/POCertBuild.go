package builder

import (
	"main/api/apiModel"
	"main/model"
	"main/proofs/executer/stellarExecuter"
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
