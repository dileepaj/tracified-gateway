package builder

import (
	"main/api/apiModel"
	"main/model"
	"main/proofs/executer/stellarExecuter"
)

// type struct{}
type ProflieMergeInterface interface {
	InsertMerge() model.MergeProfileResponse
	InsertProfile() model.MergeProfileResponse
}

type AbstractMergeProfile struct {
	MergeProfileStruct apiModel.MergeProfileStruct
	// MergingTXNs        []string
	// PreviousTXNID      string
	// PreviousProfileID  string
	// Identifiers        string
	// MergingIdentifiers []string
	// InsertType         string
	// ProfileID          string
	// Assets             string
	// Code               string
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

	object := stellarExecuter.ConcreteProfile{InsertProfileStruct: AP.MergeProfileStruct.InsertProfileStruct}

	result := object.InsertProfile()

	AP.MergeProfileStruct.InsertProfileStruct.PreviousTXNID = result.ProfileTxn
	AP.MergeProfileStruct.ProfileID = result.ProfileTxn

	object1 := stellarExecuter.ConcreteMerge{MergeProfileStruct: AP.MergeProfileStruct}
	// object1 := stellarExecuter.ConcreteMerge{Identifiers: AP.MergeProfileStruct.InsertProfileStruct.Identifier,
	// 	InsertType:         AP.MergeProfileStruct.InsertProfileStruct.Type,
	// 	PreviousTXNID:      result.ProfileTxn,
	// 	PreviousProfileID:  AP.MergeProfileStruct.ProfileID,
	// 	MergingTXNs:        AP.MergeProfileStruct.MergingTXNs,
	// 	ProfileID:          result.ProfileTxn,
	// 	MergingIdentifiers: AP.MergeProfileStruct.MergingIdentifiers}
	result1 := object1.InsertMerge()

	return result1
}
