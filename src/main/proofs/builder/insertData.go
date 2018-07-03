package builder

import (
	"fmt"

	"main/model"
	"main/proofs/executer/stellarExecuter"
)

func InsertTDP(hash string, secret string, profileId string, rootHash string) model.RootTree {
	result := stellarExecuter.InsertDataHash(hash, secret, profileId, rootHash)

	if result.Hash == "" {
		fmt.Println("Error in Stellar Executer!")
	}

	return result
}
