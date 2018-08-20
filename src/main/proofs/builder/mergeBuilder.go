package builder

import (
	"main/proofs/executer/stellarExecuter"
	"main/model"
)

// type struct{}
type ProflieMergeInterface interface {
	InsertMerge() model.MergeProfileResponse
	InsertProfile() model.MergeProfileResponse
}

type AbstractMergeProfile struct {
	MergingTXNs        []string
	PreviousTXNID      string
	PreviousProfileID  string
	Identifiers        string
	MergingIdentifiers []string
	InsertType         string
	ProfileID          string
	Assets             string
	Code               string
}

func (AP *AbstractMergeProfile) ProfileMerge() model.MergeProfileResponse {
	// result2 := ProflieMergeInterface.InsertProfile()
	// if result2.Txn == "" {
	// 	return result2
	// }

	// result := ProflieMergeInterface.InsertMerge()
	// if result.Txn == "" {
	// 	return result
	// }
	object := stellarExecuter.ConcreteProfile{Identifiers: AP.Identifiers, InsertType: "1", PreviousTXNID: AP.PreviousTXNID, PreviousProfileID: AP.PreviousProfileID}

	result := object.InsertProfile()

	object1 := stellarExecuter.ConcreteMerge{Identifiers: AP.Identifiers, InsertType: AP.InsertType, PreviousTXNID: result.Txn,
		PreviousProfileID: AP.ProfileID, MergingTXNs: AP.MergingTXNs, ProfileID: result.Txn, MergingIdentifiers: AP.MergingIdentifiers}
	result1:=object1.InsertMerge()
	return result1
}
