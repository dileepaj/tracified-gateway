package builder

import (
	"main/model"
	"main/proofs/executer/stellarExecuter"
)

// type InsertData struct{}

func TDPInsert(hash string, insertType string, previousTDPID string, profileId string) model.InsertDataResponse {
	result := stellarExecuter.InsertDataHash(hash, insertType, previousTDPID, profileId)

	return result
}
