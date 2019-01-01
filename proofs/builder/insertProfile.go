package builder

import (
	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
)

// type InsertProfile struct{}
type ProflieInsertInterface interface {
	InsertProfile() model.InsertProfileResponse
}

type AbstractProfileInsert struct {
	InsertProfileStruct apiModel.InsertProfileStruct
	// Identifiers       string
	// InsertType        string
	// PreviousTXNID     string
	// PreviousProfileID string
}

func (AP *AbstractProfileInsert) ProfileInsert() model.InsertProfileResponse {

	object := stellarExecuter.ConcreteProfile{InsertProfileStruct: AP.InsertProfileStruct}

	result := object.InsertProfile()

	return result
}
