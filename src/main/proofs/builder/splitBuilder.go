package builder

import (
	"main/model"
	"main/proofs/executer/stellarExecuter"
)

// type struct{}
type ProflieSplitInterface interface {
	InsertSplit() model.SplitProfileResponse
	InsertProfile() model.SplitProfileResponse
}

type AbstractSplitProfile struct {
	PreviousTXNID     string
	Identifiers       string
	SplitIdentifiers  []string
	InsertType        string
	ProfileID         string
	Assets            string
	Code              string
	PreviousProfileID string
}

func (AP *AbstractSplitProfile) ProfileSplit() model.SplitProfileResponse {

	var splitProfileID []string
	var result2 model.SplitProfileResponse
	if len(AP.SplitIdentifiers) >= 1 {
		for i := 0; i < len(AP.SplitIdentifiers); i++ {
			object := stellarExecuter.ConcreteProfile{
				Identifiers:       AP.SplitIdentifiers[i],
				InsertType:        "1",
				PreviousTXNID:     AP.PreviousTXNID,
				PreviousProfileID: AP.PreviousProfileID}
			result := object.InsertProfile()

			splitProfileID = append(splitProfileID, result.Txn)

			object1 := stellarExecuter.ConcreteSplit{
				PreviousTXNID: result.Txn,
				InsertType:    AP.InsertType,
				ProfileID:     result.Txn,
				Assets:        AP.Assets,
				Code:          AP.Code}

			result2 = object1.InsertSplit()

			AP.PreviousTXNID = result2.Txn
		}

	}

	result2.SplitProfiles = splitProfileID

	return result2
}
