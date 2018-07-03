package stellarRetriever

import (
	"fmt"

	"main/model"
)

func RetrievePOE(hash string, TDPId string, rootHash string) model.POE {
	fmt.Println("retieve poc!")
	poe := model.POE{"TxNHash", "", model.Error1{0, ""}}
	return poe

}
