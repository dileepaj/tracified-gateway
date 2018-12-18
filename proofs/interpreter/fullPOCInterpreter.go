package interpreter

import (
	"fmt"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/retriever/stellarRetriever"
	"net/http"
)

func (AP *AbstractPOC) InterpretFullPOC() model.POC {
	var pocObj model.POC

	object := stellarRetriever.ConcretePOC{POCStruct: AP.POCStruct}

	pocObj.RetrievePOC = object.RetrieveFullPOC()

	fmt.Println(pocObj.RetrievePOC.BCHash)
	fmt.Println(AP.POCStruct.DBTree)

	if pocObj.RetrievePOC.BCHash == nil {
		return pocObj
	} else {
		pocObj.RetrievePOC.Error = fullCompare(AP.POCStruct.DBTree, pocObj.RetrievePOC.BCHash)
		return pocObj
	}

	// return pocObj
}

func fullCompare(db []model.Current, bc []model.Current) model.Error {
	var Rerr model.Error
	// isMatch := []bool{}
	if db != nil && bc != nil {
		if len(db) == len(bc) {
			for i := 0; i < len(db); i++ {
				if db[i].TXNID == bc[i].TXNID && db[i].TType == bc[i].TType {
					switch db[i].TType {
					case "0":
						if db[i].Identifier == bc[i].Identifier {
							Rerr.Code = http.StatusOK
							Rerr.Message = "Success! BC Tree & DB Tree matched."
						} else {
							Rerr.Code = http.StatusOK
							Rerr.Message = "Error! TXN: " + db[i].TXNID + " is invalid."
							return Rerr
						}
					case "1":
						if db[i].Identifier == bc[i].Identifier && db[i].PreviousProfileID == bc[i].PreviousProfileID {
							Rerr.Code = http.StatusOK
							Rerr.Message = "Success! BC Tree & DB Tree matched."
						} else {
							Rerr.Code = http.StatusOK
							Rerr.Message = "Error! TXN: " + db[i].TXNID + " is invalid."
							return Rerr
						}

					case "2":
						fmt.Println(db[i].DataHash + " = " + bc[i].DataHash)
						fmt.Println(db[i].ProfileID + " = " + bc[i].ProfileID)
						if db[i].Identifier == bc[i].Identifier && db[i].DataHash == bc[i].DataHash && db[i].ProfileID == bc[i].ProfileID {
							Rerr.Code = http.StatusOK
							Rerr.Message = "Success! BC Tree & DB Tree matched."
						} else {

							Rerr.Code = http.StatusOK
							Rerr.Message = "Error! TXN: " + db[i].TXNID + " is invalid."
							return Rerr
						}
					case "5":

					case "6":
						if db[i].Identifier == bc[i].Identifier && db[i].ProfileID == bc[i].ProfileID {
							Rerr.Code = http.StatusOK
							Rerr.Message = "Success! BC Tree & DB Tree matched."

							return fullCompare(db[i].MergedChain, bc[i].MergedChain)

						} else {

							Rerr.Code = http.StatusOK
							Rerr.Message = "Error! TXN: " + db[i].TXNID + " is invalid."
							return Rerr
						}
						
					default:
						Rerr.Code = http.StatusOK
						Rerr.Message = "Error! Invalid Txn Type in TXN: " + db[i].TXNID + "."
						return Rerr
					}
				} else {
					Rerr.Code = http.StatusOK
					Rerr.Message = "Error! TXN: " + db[i].TXNID + " is invalid."
					return Rerr
				}
			}
		} else {
			Rerr.Code = http.StatusOK
			Rerr.Message = "Error! BC Tree & DB Tree length din't match "
			return Rerr
		}

		return Rerr
	}

	Rerr.Code = http.StatusOK
	Rerr.Message = "Error! BC Tree & DB Tree are non-existent."
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
