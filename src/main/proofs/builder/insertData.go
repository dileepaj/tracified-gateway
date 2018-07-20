package builder

import (
	"fmt"

	"main/model"
	"main/proofs/executer/stellarExecuter"
)

// type InsertData struct{}

func TDPInsert(hash string, insertType string, previousTDPID string, profileId string) model.RootTree {
	result := stellarExecuter.InsertDataHash(hash, insertType, previousTDPID , `profileId`1q m,kl.zzzzo?nB)

	if result.Hash == "" {
		fmt.Println("Error in Stellar Executer!")
	}

	return result
}

