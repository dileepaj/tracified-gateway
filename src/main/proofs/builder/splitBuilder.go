package builder

import (
	"main/proofs/executer/stellarExecuter"
	"main/model"
)

// type struct{}
type ProflieSplitInterface interface {
	InsertSplit() model.SplitProfileResponse
	InsertProfile() model.SplitProfileResponse
}

type AbstractSplitProfile struct {
	SplitProfiles string
	PreviousTXNID string
	Identifiers   string
	InsertType    string
	ProfileID     string
	Assets        string
	Code          string
	PreviousProfileID string
}

func (AP *AbstractSplitProfile) ProfileSplit() model.SplitProfileResponse {

	object := stellarExecuter.ConcreteProfile{
		Identifiers: AP.Identifiers, 
		InsertType: AP.InsertType, 
		PreviousTXNID: AP.PreviousTXNID, 
		PreviousProfileID: AP.PreviousProfileID}

	result := object.InsertProfile()

	object1:=stellarExecuter.ConcreteSplit{SplitProfiles:AP.SplitProfiles,
		PreviousTXNID :AP.PreviousTXNID,
		Identifiers :AP.Identifiers,
		InsertType:AP.InsertType,
		ProfileID:result.Txn ,
		Assets:AP.Assets,
		Code:AP.Code}

	result2:=object1.InsertSplit()	

	return result2
}
