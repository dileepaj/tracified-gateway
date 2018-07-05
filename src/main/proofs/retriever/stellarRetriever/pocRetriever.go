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

func RetrievePOC2(rootHash string) model.Node {
	fmt.Println("retieve poc!")
	poc := blockchainTree2()
	return poc

}

func blockchainTree() model.Node {
	cur1 := model.Current{"001", "E3FC18CB4776193F8AD15A947406DBYE", "cda5c9845b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee13311etxn"}

	cur2 := model.Current{"002", "E3FC18CB4776193F8AD15A947402000", "cda5c9845b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee132000txn"}

	cur3 := model.Current{"003", "E3FC18CB4776193F8AD15A947402000", "cda5c9845b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee132000txn"}

	preArr1 := []model.Node{}
	preArr2 := []model.Node{}
	preArr3 := []model.Node{}

	pre1 := model.Node{preArr1, cur1, "1000500"}
	preArr2 =[]model.Node{pre1}

	pre2 := model.Node{preArr2, cur2, "1000501"}
	preArr3 =[]model.Node{pre2}

	pre3 := model.Node{preArr3, cur3, "1000502"}

	return pre3
}
func blockchainTree2() model.Node {
	cur1 := model.Current{"001", "E3FC18CB4776193F8AD15A947406DBYE", "cda5c9845b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee13311etxn"}

	cur2 := model.Current{"002", "E3FC18CB4776193F8AD15A947402000", "cda5c9845b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee132000txn"}
	// cur2 := model.Current{"00dvvdfv2", "E3FC18CB4776193F8AdffvfdD15A947402000", "cda5c9845b218fdd8e8f0dfvdfvdfvfd4d2db3a82db6b59d1d78785a214cbce3ee132000txn"}

	cur3 := model.Current{}

	preArr1 := []model.Node{}
	preArr2 := []model.Node{}
	preArr3 := []model.Node{}

	pre1 := model.Node{preArr1, cur1, "1000500"}
	preArr2 = []model.Node{pre1}
	pre2 := model.Node{preArr2, cur2, "1000501"}

	preArr3 = []model.Node{pre2, pre1}

	pre3 := model.Node{preArr3, cur3, "1000502"}

	return pre3
}


