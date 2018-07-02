package interpreter

import (
	"fmt"
	. "main/model"

	stellarRetriever "main/proofs/retriever/stellarRetriever"
)

func InterpretPOC(rootHash string, isValid bool) Node {
	result := stellarRetriever.RetrievePOC(rootHash, isValid)

	if result.Previous == "" && result.Current == "" {
		fmt.Println("Error in stellarRetriever!")
	}

	return result
}
