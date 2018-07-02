package proofBuilder

import (
	"fmt"
	. "main/model"
	stellarexecuter "main/proofs/proofExecuter/stellarExecuter"
)

func InsertTDP(hash string, secret string, profileId string, rootHash string) RootTree {
	result := stellarexecuter.InsertDataHash(hash, secret, profileId, rootHash)

	if result.Hash == "" {
		fmt.Println("Error in Stellar Executer!")
	}

	return result
}
