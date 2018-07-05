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

	// isMapped := compare(dBTree, blockchainTree, "")
	isMapped := testCompare(dBTree, blockchainTree, "1000503")
	fmt.Println(isMapped)
	var returnRes model.PocSuccess

	if isMapped == true {
		returnRes = model.PocSuccess{rootHash, model.Error1{0, ""}}
	}

	return returnRes
}

func InterpretPOCFakeTree(rootHash string, treeObj string) model.PocSuccess {
	//
	blockchainTree := stellarRetriever.RetrievePOC(rootHash)

	dBTree := dBTreeFalse()

	// isMapped := compare(dBTree, blockchainTree, "")
	isMapped := testCompare(dBTree, blockchainTree, "1000503")
	fmt.Println(isMapped)
	var returnRes model.PocSuccess

	if isMapped == true {
		returnRes = model.PocSuccess{rootHash, model.Error1{0, ""}}
	}

	return returnRes
}
func InterpretPOCError(rootHash string, treeObj string) model.PocSuccess {
	//
	blockchainTree := stellarRetriever.RetrievePOC(rootHash)

	dBTree := dBTree()

	// isMapped := compare(dBTree, blockchainTree, "")
	isMapped := testCompare(dBTree, blockchainTree, "1000500")
	fmt.Println(isMapped)
	var returnRes model.PocSuccess

	if isMapped == true {
		returnRes = model.PocSuccess{rootHash, model.Error1{0, ""}}
	}

	return returnRes
}

func InterpretPOC2(rootHash string, treeObj string) model.PocSuccess {
	//
	blockchainTree := stellarRetriever.RetrievePOC2(rootHash)
	// blockchainTree := dBTree2()

	dBTree := dBTree2()

	// isMapped := compare(dBTree, blockchainTree, "")
	isMapped := testCompare(dBTree, blockchainTree, "1000503")
	fmt.Println(isMapped)
	var returnRes model.PocSuccess

	if isMapped == true {
		returnRes = model.PocSuccess{rootHash, model.Error1{0, ""}}
	}

	return returnRes
}

func InterpretPOC2FakeTree(rootHash string, treeObj string) model.PocSuccess {
	//
	blockchainTree := stellarRetriever.RetrievePOC2(rootHash)
	// blockchainTree := dBTree2()

	dBTree := dBTree2False()

	// isMapped := compare(dBTree, blockchainTree, "")
	isMapped := testCompare(dBTree, blockchainTree, "1000503")
	fmt.Println(isMapped)
	var returnRes model.PocSuccess

	if isMapped == true {
		returnRes = model.PocSuccess{rootHash, model.Error1{0, ""}}
	}

	return returnRes
}

func InterpretPOC2Error(rootHash string, treeObj string) model.PocSuccess {
	//
	blockchainTree := stellarRetriever.RetrievePOC2(rootHash)
	// blockchainTree := dBTree2()

	dBTree := dBTree2()

	// isMapped := compare(dBTree, blockchainTree, "")
	isMapped := testCompare(dBTree, blockchainTree, "1000500")
	fmt.Println(isMapped)
	var returnRes model.PocSuccess

	if isMapped == true {
		returnRes = model.PocSuccess{rootHash, model.Error1{0, ""}}
	}

	return returnRes
}


// func compare(db model.Node, bc model.Node, seq string) bool {

// 	if db.Previous != nil && bc.Previous != nil {
// 		if db.Current.TDPID == bc.Current.TDPID && db.Current.DataHash == bc.Current.DataHash {
// 			//check whether there a re multiple pre?
// 			//traverse through the length.array
// 			// compare(db.Previous[0], bc.Previous[0])
// 			if seq != "" && db.Sequence == db.Sequence {
// 				return true
// 			} else {
// 				return false
// 			}
// 		} else {
// 			return false
// 		}

// 	} else {
// 		return false
// 	}

// 	return true
// }

func testCompare(db model.Node, bc model.Node, seq string) bool {
	var isMatch bool
	fmt.Println(seq, ">", db.Sequence)
	if db.Current != (model.Current{}) && db.Previous != nil && bc.Previous != nil {
		// fmt.Println(len(db.Previous), len(bc.Previous))
		// fmt.Println(db.Sequence)
		if db.Current.TDPID == bc.Current.TDPID && db.Current.DataHash == bc.Current.DataHash {
			fmt.Println("Sequence ", db.Sequence, ": TDPID and DataHash Matched")
			isMatch = true
			if db.Sequence == bc.Sequence && db.Sequence < seq {

				matchArray := make([]bool, len(bc.Previous))

				for i := 0; i < len(db.Previous); i++ {

					matchArray[i] = testCompare(db.Previous[i], bc.Previous[i], db.Sequence)
				}
				isMatch = checkBoolArray(matchArray)

			} else {
				fmt.Println("Sequence Mismatch")
				isMatch = false
			}

			return isMatch
		}
			isMatch=false
		
		
	}

	if db.Current != (model.Current{}) && db.Previous == nil && bc.Previous == nil {
		// fmt.Println(db.Sequence)
		if db.Current.TDPID == bc.Current.TDPID && db.Current.DataHash == bc.Current.DataHash {
			fmt.Println("Sequence ", db.Sequence, ": TDPID and DataHash Matched")
			isMatch = true
			return isMatch
		}
			isMatch=false
		
	}

	if db.Current == (model.Current{}) && db.Previous != nil && bc.Previous != nil {
		// fmt.Println(db.Sequence)
		if db.Sequence == bc.Sequence && db.Sequence < seq {
			matchArray := make([]bool, len(bc.Previous))

			for i := 0; i < len(db.Previous); i++ {

				matchArray[i] = testCompare(db.Previous[i], bc.Previous[i], db.Sequence)
			}
			isMatch = checkBoolArray(matchArray)

		} else {
			fmt.Println("Sequence Mismatch")
			isMatch = false
		}
	}

	return isMatch
}

func checkBoolArray(array []bool) bool {
	isMatch:=true
	for i := 0; i < len(array); i++ {
		if array[i]==false {
			isMatch=false
			return isMatch
		}
	}
	return isMatch
}
