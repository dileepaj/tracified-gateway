package interpreter

import (
	"fmt"

	"main/model"
	"main/proofs/retriever/stellarRetriever"
)

func InterpretPOC(rootHash string, isValid bool) model.Node {
	result := stellarRetriever.RetrievePOC(rootHash, isValid)

	if result.Previous == "" && result.Current == "" {
		fmt.Println("Error in stellarRetriever!")
	}

	return result
}
