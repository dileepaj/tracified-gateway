package interpreter

import (
	"fmt"
	"main/model"
	"main/proofs/retriever/stellarRetriever"
	"net/http"
)

// type AbstractPOC struct {
// 	Txn       string
// 	ProfileID string
// 	DBTree    []model.Current
// 	BCTree    []model.Current
// }

func (AP *AbstractPOC) InterpretFullPOC() model.POC {
	var pocObj model.POC

	object := stellarRetriever.ConcretePOC{
		Txn:       AP.Txn,
		ProfileID: AP.ProfileID,
		DBTree:    AP.DBTree}

	pocObj.RetrievePOC = object.RetrieveFullPOC()

	fmt.Println(pocObj.RetrievePOC.BCHash)

	if pocObj.RetrievePOC.BCHash == nil {
		return pocObj
	} else {
		pocObj.RetrievePOC.Error = fullCompare(AP.DBTree, pocObj.RetrievePOC.BCHash)
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
						if db[i].Identifier != bc[i].Identifier {
							Rerr.Code = http.StatusOK
							Rerr.Message = "Error! BC Tree & DB Tree Genesis Identifiers din't match."
							return Rerr
						} else {
							Rerr.Code = http.StatusOK
							Rerr.Message = "Success! BC Tree & DB Tree Genesis Identifiers matched."
							return Rerr
						}
					case "1":
						if db[i].Identifier != bc[i].Identifier && db[i].PreviousProfileID != bc[i].PreviousProfileID {
							Rerr.Code = http.StatusOK
							Rerr.Message = "Error! BC Tree & DB Tree profile Identifiers & previous profileID din't match."
							return Rerr
						} else {
							Rerr.Code = http.StatusOK
							Rerr.Message = "Success! BC Tree & DB Tree profile Identifiers & previous profileID matched."
							return Rerr
						}

					case "2":
						fmt.Println(db[i].DataHash + " = " + bc[i].DataHash)
						fmt.Println(db[i].ProfileID + " = " + bc[i].ProfileID)
						if db[i].DataHash != bc[i].DataHash && db[i].ProfileID != bc[i].ProfileID {
							Rerr.Code = http.StatusOK
							Rerr.Message = "Error! BC Tree & DB Tree TDP DataHash & profileID din't match."
							return Rerr
						} else {

						}
					case "5":

					case "6":
						return fullCompare(db[i].MergedChain, bc[i].MergedChain)
					default:
						Rerr.Code = http.StatusOK
						Rerr.Message = "Error! Invalid Txn Type."
						return Rerr
					}
				} else {
					Rerr.Code = http.StatusOK
					Rerr.Message = "Error! BC Tree & DB Tree TxnID & Type din't match."
					return Rerr
				}
			}
		} else {
			Rerr.Code = http.StatusOK
			Rerr.Message = "Error! BC Tree & DB Tree length din't match."
			return Rerr
		}

	}

	Rerr.Code = http.StatusOK
	Rerr.Message = "Error! BC Tree & DB Tree din't match."
	return Rerr
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
