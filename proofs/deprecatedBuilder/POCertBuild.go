package deprecatedBuilder

import (
	"github.com/dileepaj/tracified-gateway/proofs/executer/deprecatedStellarExecuter"
	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"
)

// type InsertPOA struct{}
type POCertInsertInterface interface {
	InsertPOCertHash() model.InsertDataResponse
}

type AbstractPOCertInsert struct {
	InsertPOCertStruct apiModel.InsertPOCertStruct
}

func (AP *AbstractPOCertInsert) POCertInsert() model.InsertDataResponse {

	object := deprecatedStellarExecuter.ConcreteInsertPOCert{InsertPOCertStruct: AP.InsertPOCertStruct}

	result := object.InsertPOCertHash()

	return result
}
