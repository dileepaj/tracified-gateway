package builder

import (
	"main/model"
	"main/proofs/executer/stellarExecuter"
)

// type InsertProfile struct{}
type ProflieInsertInterface interface {
	InsertProfile() model.InsertProfileResponse
}

type AbstractProfileInsert struct {
	Identifiers       string
	InsertType        string
	PreviousTXNID     string
	PreviousProfileID string
}

func (AP *AbstractProfileInsert) ProfileInsert() model.InsertProfileResponse {

	object := stellarExecuter.ConcreteProfile{Identifiers: AP.Identifiers, InsertType: AP.InsertType, PreviousTXNID: AP.PreviousTXNID, PreviousProfileID: AP.PreviousProfileID}

	result := object.InsertProfile()

	return result
}
