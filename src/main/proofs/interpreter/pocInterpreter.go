package interpreter

import (
	"fmt"
	"main/model"
	"main/proofs/retriever/stellarRetriever"
)

//all splits and normal scenarios
func InterpretPOC(rootHash string, treeObj string) model.PocSuccess {
	//
	blockchainTree := stellarRetriever.RetrievePOC(rootHash)

	dBTree := dBTree()

	// isMapped := compare(dBTree, blockchainTree, "")
	isMapped := testCompare(dBTree, blockchainTree, "1000504")
	fmt.Println(isMapped)
	var returnRes model.PocSuccess

	if isMapped == true {
		returnRes = model.PocSuccess{rootHash, model.Error1{0, ""}}
	}

	return returnRes
}

// Used to check for the unmatched transaction ids and datahashes
//in normal and split scenarios
func InterpretPOCFakeTree(rootHash string, treeObj string) model.PocSuccess {
	//
	blockchainTree := stellarRetriever.RetrievePOC(rootHash)

	dBTree := dBTreeFalse()

	// isMapped := compare(dBTree, blockchainTree, "")
	isMapped := testCompare(dBTree, blockchainTree, "1000504")
	fmt.Println(isMapped)
	var returnRes model.PocSuccess

	if isMapped == true {
		returnRes = model.PocSuccess{rootHash, model.Error1{0, ""}}
	}

	return returnRes
}

//used to check for error in Sequence Order of each transmission
//in normal and split scenarios
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

//merge scenario
func InterpretPOC2(rootHash string, treeObj string) model.PocSuccess {
	//
	blockchainTree := stellarRetriever.RetrievePOC2(rootHash)
	// blockchainTree := dBTree2()

	dBTree := dBTree2()

	// isMapped := compare(dBTree, blockchainTree, "")
	isMapped := testCompare(dBTree, blockchainTree, "1000504")
	fmt.Println(isMapped)
	var returnRes model.PocSuccess

	if isMapped == true {
		returnRes = model.PocSuccess{rootHash, model.Error1{0, ""}}
	}

	return returnRes
}

//Used to check for the unmatched transaction ids and datahashes
//in merge scenario
func InterpretPOC2FakeTree(rootHash string, treeObj string) model.PocSuccess {
	//
	blockchainTree := stellarRetriever.RetrievePOC2(rootHash)
	// blockchainTree := dBTree2()

	dBTree := dBTree2False()

	// isMapped := compare(dBTree, blockchainTree, "")
	isMapped := testCompare(dBTree, blockchainTree, "1000504")
	fmt.Println(isMapped)
	var returnRes model.PocSuccess

	if isMapped == true {
		returnRes = model.PocSuccess{rootHash, model.Error1{0, ""}}
	}

	return returnRes
}

//used to check for error in Sequence Order of each transmission
//in merge scenario
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
	// if seq == "" {
	// 	seq = bc.Sequence
	// }
	fmt.Println(seq, ">", db.Sequence)

	//If current and previous are both present
	//will be incase of all normal transactions and splits.
	if bc.Current != (model.Current{}) && db.Previous != nil && bc.Previous != nil {

		if db.Current.TDPID == bc.Current.TDPID && db.Current.DataHash == bc.Current.DataHash {
			fmt.Println("Sequence ", bc.Sequence, ": TDPID and DataHash Matched")
			isMatch = true
			if db.Sequence == bc.Sequence && bc.Sequence < seq {

				matchArray := make([]bool, len(bc.Previous))

				for i := 0; i < len(db.Previous); i++ {

					matchArray[i] = testCompare(db.Previous[i], bc.Previous[i], bc.Sequence)
				}
				isMatch = checkBoolArray(matchArray)

			} else {
				fmt.Println("Sequence Mismatch")
				isMatch = false
			}

			return isMatch
		}
		fmt.Println("Sequence ", bc.Sequence, ": TDPID and DataHash UnMatched")
		isMatch = false

	}
	//if Only current is present...
	//like at:- Genesis Block,
	if bc.Current != (model.Current{}) && db.Previous == nil && bc.Previous == nil {
		if db.Current.TDPID == bc.Current.TDPID && db.Current.DataHash == bc.Current.DataHash {
			fmt.Println("Sequence ", bc.Sequence, ": TDPID and DataHash Matched")
			isMatch = true
			return isMatch
		}
		fmt.Println("Sequence ", bc.Sequence, ": TDPID and DataHash UnMatched")
		isMatch = false

	}

	//if no Current but multiple previous are present..
	//in case of merge scenarios..
	if bc.Current == (model.Current{}) && db.Previous != nil && bc.Previous != nil {
		fmt.Println("Sequence ", bc.Sequence, ": No Current")

		if db.Sequence == bc.Sequence && db.Sequence < seq {
			matchArray := make([]bool, len(bc.Previous))

			for i := 0; i < len(bc.Previous); i++ {

				matchArray[i] = testCompare(db.Previous[i], bc.Previous[i], bc.Sequence)
			}
			isMatch = checkBoolArray(matchArray)

		} else {
			fmt.Println("Sequence Mismatch")
			isMatch = false
		}
	}

	return isMatch
}

//checks the multiple boolean indexes in an array and returns the combined result.
func checkBoolArray(array []bool) bool {
	isMatch := true
	for i := 0; i < len(array); i++ {
		if array[i] == false {
			isMatch = false
			return isMatch
		}
	}
	return isMatch
}
