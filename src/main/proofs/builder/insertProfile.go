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
	PreviousTDPID     string
	PreviousProfileID string
}

func (AP *AbstractProfileInsert) ProfileInsert() model.InsertProfileResponse {

	object := stellarExecuter.ConcreteProfile{Identifiers: AP.Identifiers, InsertType: AP.InsertType, PreviousTDPID: AP.PreviousTDPID, PreviousProfileID: AP.PreviousProfileID}

	result := object.InsertProfile()

	return result
}
