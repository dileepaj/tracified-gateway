package builder

import (
	"main/model"
)

// type struct{}
type ProflieSplitInterface interface {
	InsertSplit() model.SplitProfileResponse
	InsertProfile() model.SplitProfileResponse
}

type AbstractSplitProfile struct {
}

func (AP *AbstractSplitProfile) ProfileSplit(ProflieSplitInterface ProflieSplitInterface) model.SplitProfileResponse {
	result2 := ProflieSplitInterface.InsertProfile()
	if result2.Txn == "" {
		return result2
	}

	result := ProflieSplitInterface.InsertSplit()
	if result.Txn == "" {
		return result
	}

	return result
}
