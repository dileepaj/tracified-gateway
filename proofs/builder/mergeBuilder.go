package builder

import (
	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
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

	temp := apiModel.InsertProfileStruct{Type: "1",
		PreviousProfileID: AP.MergeProfileStruct.PreviousProfileID,
		PreviousTXNID:     AP.MergeProfileStruct.PreviousTXNID,
		Identifier:        AP.MergeProfileStruct.Identifier}

	object := stellarExecuter.ConcreteProfile{
		InsertProfileStruct: temp}

	result := object.InsertProfile()

	AP.MergeProfileStruct.PreviousTXNID = result.ProfileTxn
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
