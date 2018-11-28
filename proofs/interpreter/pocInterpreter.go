package interpreter

import (
	"fmt"
	"github.com/tracified-gateway/api/apiModel"
	"github.com/tracified-gateway/model"
	"github.com/tracified-gateway/proofs/retriever/stellarRetriever"
	"net/http"
)

type AbstractPOC struct {
	POCStruct apiModel.POCStruct
	// Txn       string
	// ProfileID string
	// DBTree    []model.Current
	// BCTree    []model.Current
}

func (AP *AbstractPOC) InterpretPOC() model.POC {
	var pocObj model.POC

	object := stellarRetriever.ConcretePOC{POCStruct: AP.POCStruct}

	pocObj.RetrievePOC = object.RetrievePOC()

	fmt.Println(AP.POCStruct.DBTree)
	fmt.Println(pocObj.RetrievePOC.BCHash)

	if pocObj.RetrievePOC.BCHash == nil {
		return pocObj
	} else {
		pocObj.RetrievePOC.Error = testCompare(AP.POCStruct.DBTree, pocObj.RetrievePOC.BCHash)
		return pocObj
	}

	// return pocObj
}

func testCompare(db []model.Current, bc []model.Current) model.Error {
	var Rerr model.Error
	if db != nil && bc != nil {
		if len(db) == len(bc) {
			for i := 0; i < len(db); i++ {
				if db[i].TXNID == bc[i].TXNID && db[i].TType == bc[i].TType && db[i].DataHash == bc[i].DataHash {
					Rerr.Code = http.StatusOK
					Rerr.Message = "Success! The Tree exists in the Blockchain"

				} else {
					Rerr.Code = http.StatusOK
					Rerr.Message = "Error! TXN:" + db[i].TXNID + " is invalid."
					return Rerr
				}
			}
		} else {
			Rerr.Code = http.StatusOK
			Rerr.Message = "Error! BC Tree & DB Tree are of different length."
			return Rerr
		}

	} else {
		Rerr.Code = http.StatusOK
		Rerr.Message = "Error! BC Tree & DB Tree are Non-existant."
		return Rerr
	}

	return Rerr
}

//checks the multiple boolean indexes in an array and returns the combined result.
// func checkBoolArray(array []bool) bool {
// 	isMatch := true
// 	for i := 0; i < len(array); i++ {
// 		if array[i] == false {
// 			isMatch = false
// 			return isMatch
// 		}
// 	}
// 	return isMatch
// }
