package builder

import (
	"main/model"
)

// type struct{}
type ProflieSplitInterface interface {
	InsertSplit() model.SplitProfileResponse
}

type AbstractSplitProfile struct {
}

func (AP *AbstractSplitProfile) ProfileSplit(ProflieSplitInterface ProflieSplitInterface) model.SplitProfileResponse {
	result := ProflieSplitInterface.InsertSplit()

	return result
}
