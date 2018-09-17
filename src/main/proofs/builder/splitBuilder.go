package builder

import (
	"main/api/apiModel"
	"main/model"
	"main/proofs/executer/stellarExecuter"
)

// type struct{}
type ProflieSplitInterface interface {
	InsertSplit() model.SplitProfileResponse
	InsertProfile() model.SplitProfileResponse
}

type AbstractSplitProfile struct {
	SplitProfileStruct apiModel.SplitProfileStruct
	// PreviousTXNID     string
	// Identifiers       string
	// SplitIdentifiers  []string
	// InsertType        string
	// ProfileID         string
	// Assets            []string
	// Code              string
	// PreviousProfileID string
}

func (AP *AbstractSplitProfile) ProfileSplit() model.SplitProfileResponse {

	var splitProfileID []string
	var splitTXN []string
	var result2 model.SplitProfileResponse

	if len(AP.SplitProfileStruct.SplitIdentifiers) >= 1 {

		for i := 0; i < len(AP.SplitProfileStruct.SplitIdentifiers); i++ {
			temp:=apiModel.InsertProfileStruct{
				Type:"1",
				PreviousProfileID:AP.SplitProfileStruct.PreviousProfileID,
				PreviousTXNID:AP.SplitProfileStruct.PreviousTXNID,
				Identifier:AP.SplitProfileStruct.SplitIdentifiers[i]}

			object := stellarExecuter.ConcreteProfile{
				InsertProfileStruct:temp}

			result := object.InsertProfile()

			splitProfileID = append(splitProfileID, result.ProfileTxn)

			object1 := stellarExecuter.ConcreteSplit{
				SplitProfileStruct: AP.SplitProfileStruct, 
				CurAssets: AP.SplitProfileStruct.Assets[i]}
		

			result2 = object1.InsertSplit()

			splitTXN = append(splitTXN, result2.Txn)

		}

	}

	result2.SplitProfiles = splitProfileID
	result2.SplitTXN = splitTXN

	return result2
}
