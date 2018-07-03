package interpreter

import (
	"main/model"
)

func InterpretPOC(rootHash string, treeObj string) model.PocSuccess {
	//
	// blockchainTree := stellarRetriever.RetrievePOC(rootHash)

	// dBTree := dBTree()

	// isMapped := compare(dBTree, blockchainTree, "")

	// if isMapped == true {
	returnRes := model.PocSuccess{rootHash, model.Error1{0, ""}}
	// }

	return returnRes
}

func dBTree() model.Node {
	cur1 := model.Current{"001", "E3FC18CB4776193F8AD15A947406DBYE", "cda5c9845b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee13311etxn"}

	cur2 := model.Current{"002", "E3FC18CB4776193F8AD15A947402000", "cda5c9845b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee132000txn"}

	preArr1 := []model.Node{}
	preArr2 := []model.Node{}

	pre1 := model.Node{preArr1, cur1, "1000500"}
	preArr2[0] = pre1

	pre2 := model.Node{preArr2, cur2, "1000501"}

	return pre2
}

func compare(db model.Node, bc model.Node, seq string) bool {

	if db.Previous != nil && bc.Previous != nil {
		if db.Current.TDPID == bc.Current.TDPID && db.Current.DataHash == bc.Current.DataHash {
			//check whether there a re multiple pre?
			//traverse through the length.array
			// compare(db.Previous[0], bc.Previous[0])
			if seq != "" && db.Sequence == db.Sequence {
				return true
			} else {
				return false
			}
		} else {
			return false
		}

	} else {
		return false
	}

	return true
}
