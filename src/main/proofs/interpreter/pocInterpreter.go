package interpreter

import (
	"main/model"
	// "main/api/apiModel"
)

type POCInterface interface {
	RetrievePOC() model.RetrievePOC
}

type AbstractPOC struct {
}

func (AP *AbstractPOC) InterpretPOC(POCInterface POCInterface) model.RetrievePOC {
	var pocObj model.POC

	pocObj.RetrievePOC = POCInterface.RetrievePOC()

	isMapped := testCompare(pocObj.RetrievePOC.DBHash, pocObj.RetrievePOC.BCHash)

	var returnRes model.RetrievePOC
	if isMapped == true {
		pocObj.RetrievePOC.Error.Message = "Chain Exists"
		returnRes = pocObj.RetrievePOC
	}

	return returnRes
}

func testCompare(db model.Node, bc model.Node) bool {
	var isMatch bool
	// if seq == "" {
	// 	seq = bc.Sequence
	// }

	//If current and previous are both present
	//will be incase of all normal transactions and splits.
	if bc.Current != (model.Current{}) && db.Previous != nil && bc.Previous != nil {

		if db.Current.TDPID == bc.Current.TDPID {
			isMatch = true
			matchArray := make([]bool, len(bc.Previous))

			for i := 0; i < len(db.Previous); i++ {

				matchArray[i] = testCompare(db.Previous[i], bc.Previous[i])
			}
			isMatch = checkBoolArray(matchArray)
			return isMatch
		}

	}
	//if Only current is present...
	//like at:- Genesis Block,
	if bc.Current != (model.Current{}) && db.Previous == nil && bc.Previous == nil {
		if db.Current.TDPID == bc.Current.TDPID {
			isMatch = true
			return isMatch
		}
	}

	//if no Current but multiple previous are present..
	//in case of merge scenarios..
	if bc.Current == (model.Current{}) && db.Previous != nil && bc.Previous != nil {

		matchArray := make([]bool, len(bc.Previous))

		for i := 0; i < len(bc.Previous); i++ {

			matchArray[i] = testCompare(db.Previous[i], bc.Previous[i])
		}
		isMatch = checkBoolArray(matchArray)

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
