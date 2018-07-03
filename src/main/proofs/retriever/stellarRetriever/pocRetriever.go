package stellarRetriever

import (
	"fmt"

	"main/model"
)

func RetrievePOC(rootHash string) model.Node {
	fmt.Println("retieve poc!")
	poc := blockchainTree()
	return poc

}

func blockchainTree() model.Node {
	cur1 := model.Current{"001", "E3FC18CB4776193F8AD15A947406DBYE", "cda5c9845b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee13311etxn"}

	cur2 := model.Current{"002", "E3FC18CB4776193F8AD15A947402000", "cda5c9845b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee132000txn"}

	preArr1 := []model.Node{}
	preArr2 := []model.Node{}

	pre1 := model.Node{preArr1, cur1, "1000500"}
	preArr2[0] = pre1

	pre2 := model.Node{preArr2, cur2, "1000501"}

	return pre2
}
