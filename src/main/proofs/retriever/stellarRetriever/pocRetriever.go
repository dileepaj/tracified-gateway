package stellarRetriever

import (
	"fmt"

	"main/model"
)

func RetrievePOC(rootHash string, isValid bool) model.Node {
	fmt.Println("retieve poc!")
	poc := model.Node{"PRE", "CUR"}
	return poc

}
