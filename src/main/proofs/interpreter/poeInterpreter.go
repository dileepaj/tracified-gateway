package interpreter

import (
	"fmt"

	"main/model"
	"main/proofs/retriever/stellarRetriever"
)

func InterpretPOE(hash string, TDPId string, rootHash string) model.POE {
	result := stellarRetriever.RetrievePOE(hash, TDPId, rootHash)

	if result.TxNHash == "" && result.RootHash == "" {
		fmt.Println("Error in Stellar!")
	}

	return result
}
