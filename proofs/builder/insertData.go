package builder

import (
	"fmt"

	"github.com/Tracified-Gateway/model"
	"github.com/Tracified-Gateway/proofs/executer/stellarExecuter"
)

type InsertData struct{}

func InsertTDP(hash string, secret string, profileId string, rootHash string) model.RootTree {
	result := stellarExecuter.InsertDataHash(hash, secret, profileId, rootHash)

	if result.Hash == "" {
		fmt.Println("Error in Stellar Executer!")
	}

	return result
}

func (I *InsertData) TDPInsert(hash string, secret string, profileId string, rootHash string) model.RootTree {
	result := stellarExecuter.InsertDataHash(hash, secret, profileId, rootHash)

	if result.Hash == "" {
		fmt.Println("Error in Stellar Executer!")
	}

	return result
}
