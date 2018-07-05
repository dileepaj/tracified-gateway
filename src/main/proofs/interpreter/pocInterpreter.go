package interpreter

import (
	"fmt"
	"main/model"
	"main/proofs/retriever/stellarRetriever"
)

func InterpretPOC(rootHash string, treeObj string) model.PocSuccess {
	//
	blockchainTree := stellarRetriever.RetrievePOC(rootHash)

	dBTree := dBTree()

	isMapped := CompareTree(dBTree, blockchainTree, blockchainTree.Current.Sequence)
	fmt.Println("isMapped")
	fmt.Println(isMapped)
	fmt.Println(treeObj)

	returnRes := model.PocSuccess{rootHash, model.Error1{0, "error in InterpretPOC"}}

	if isMapped == true {
		returnRes = model.PocSuccess{rootHash, model.Error1{0, "Mapped"}}
	} else {
		returnRes = model.PocSuccess{rootHash, model.Error1{0, "Not Mapped"}}
	}

	return returnRes
}

func dBTree() model.Node {
	cur1 := model.Current{"001", "cur1datahash", "cur1txn", "1000"}
	cur20 := model.Current{"0020", "cur20datahash", "cur20txn", "2000"}
	cur21 := model.Current{"0021", "cur21datahash", "cur21txn", "21000"}
	cur4 := model.Current{"004", "cur4datahash", "cur4txn", "4000"}
	nullCurr := model.Current{}

	preArr1 := []model.Node{}
	pre1 := model.Node{preArr1, cur1}

	preArr20 := []model.Node{pre1}
	pre20 := model.Node{preArr20, cur20}

	preArr21 := []model.Node{}
	pre21 := model.Node{preArr21, cur21}

	preArr3 := []model.Node{pre20, pre21}
	pre3 := model.Node{preArr3, nullCurr}

	preArr4 := []model.Node{pre3}
	pre4 := model.Node{preArr4, cur4}

	return pre4
}

// func compare(db model.Node, bc model.Node, seq string) bool {

// 	isSame := false
// 	fmt.Println(db)
// 	fmt.Println(bc)

// 	if len(db.Previous) < 2 {

// 		if db.Current.TDPID == bc.Current.TDPID && db.Current.DataHash == bc.Current.DataHash {
// 			fmt.Println("seq")
// 			fmt.Println(seq)
// 			fmt.Println("db.Sequence")
// 			fmt.Println(db.Current.Sequence)
// 			fmt.Println("bc.Sequence")
// 			fmt.Println(bc.Current.Sequence)
// 			// fmt.Println(bc.Current.DataHash)
// 			// if seq != "" && db.Sequence == bc.Sequence {

// 			// if seq > db.Sequence && seq > bc.Sequence {
// 			// 	// fmt.Println(bc.Sequence)
// 			// 	isSame = true
// 			// } else {
// 			// 	isSame = false
// 			// 	return isSame
// 			// 	fmt.Println("check1")
// 			// }

// 			// } else {
// 			// 	isSame = false
// 			// 	// return isSame
// 			// 	fmt.Println("check2")
// 			// }

// 			// isSame = true
// 			// fmt.Println("check seq")
// 			// fmt.Println(bc.Sequence)

// 		} else {
// 			isSame = false
// 			return isSame
// 			fmt.Println("check3")
// 		}
// 	}

// 	for i := 0; i < len(db.Previous); i++ {
// 		fmt.Println("loop")
// 		fmt.Println(i + 1)
// 		compare(db.Previous[i], bc.Previous[i], db.Current.Sequence)

// 	}

// 	return isSame

// }

func CompareTree(db model.Node, bc model.Node, seq string) bool {

	isSame := false
	nullCurr := model.Current{}
	fmt.Println(nullCurr)
	fmt.Println(bc.Current)

	if bc.Current != nullCurr {
		if len(bc.Previous) == 0 {
			isSame = CheckCurrent(bc.Current.TDPID, bc.Current.TDPID, db.Current.TDPID, db.Current.TDPID)
			if isSame == false {
				return isSame
			} else {
				isSame = true
			}
			if len(db.Previous) == 1 {
				isSame = CompareTree(db.Previous[0], bc.Previous[0], bc.Current.Sequence)
				if isSame == false {
					return isSame
				} else {
					isSame = true
				}
			}
		} 
	} else {
		if len(db.Previous) > 1 {
			isTrue := true
			for i := 0; i < len(db.Previous); i++ {
				fmt.Println("loop")
				fmt.Println(i + 1)
				isTrue = CompareTree(db.Previous[i], bc.Previous[i], bc.Current.Sequence)
				isSame = isSame && isTrue
			}
			if isSame == false {
				return isSame
			} else {
				isSame = true
			}
	}else{ 
	isSame = false
		fmt.Println("chk")
		return isSame
		//tree is breaked coz cannot be without a current
	}
	return isSame
}

func CheckCurrent(bcTDPID string, bcHash string, dbTDPID string, dbHash string) bool {

	if bcTDPID == dbTDPID && bcHash == dbHash {
		return true
	} else {
		return false
	}

}
