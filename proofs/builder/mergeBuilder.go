package builder

import (
	"github.com/Tracified-Gateway/model"
)

// type struct{}
type ProflieMergeInterface interface {
	InsertMerge() model.MergeProfileResponse
	InsertProfile() model.MergeProfileResponse
}

type AbstractMergeProfile struct {
}

func (AP *AbstractMergeProfile) ProfileMerge(ProflieMergeInterface ProflieMergeInterface) model.MergeProfileResponse {
	result2 := ProflieMergeInterface.InsertProfile()
	if result2.Txn == "" {
		return result2
	}

	result := ProflieMergeInterface.InsertMerge()
	if result.Txn == "" {
		return result
	}

	return result
}
