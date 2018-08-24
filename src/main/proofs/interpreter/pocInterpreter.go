package interpreter

import (
	"fmt"
	"main/proofs/retriever/stellarRetriever"
	"main/model"
)

// type POCInterface interface {
// 	RetrievePOC() model.RetrievePOC
// }

type AbstractPOC struct {
	Txn       string
	ProfileID string
	DBTree    []model.Current
	BCTree    []model.Current
}

func (AP *AbstractPOC) InterpretPOC() model.POC {
	var pocObj model.POC

	 object:= stellarRetriever.ConcretePOC{
		 Txn: AP.Txn, 
		 ProfileID: AP.ProfileID, 
		 DBTree: AP.DBTree}
		 
	 pocObj.RetrievePOC=object.RetrievePOC()
	
	 fmt.Println(pocObj.RetrievePOC.DBHash)
	 fmt.Println(pocObj.RetrievePOC.BCHash)

	isMapped := testCompare(pocObj.RetrievePOC.DBHash, pocObj.RetrievePOC.BCHash)
	
	if isMapped == true {
		pocObj.RetrievePOC.Error.Message = "Chain Exists in the Blockchain"
		pocObj.RetrievePOC.Error.Code = 200
	} else {
		pocObj.RetrievePOC.Error.Message = "Chain Doesn't Exist in the Blockchain"
		pocObj.RetrievePOC.Error.Code = 200
	}

	return pocObj
}

func testCompare(db []model.Current, bc []model.Current) bool {
	isMatch := []bool{}
	if db != nil && bc != nil {
		if len(db) == len(bc) {
			for i := 0; i < len(db); i++ {
				if db[i].TXNID == bc[i].TXNID && db[i].TType == bc[i].TType {
					datamatch:=[]bool{}

					if len(db[i].DataHash)== len(bc[i].DataHash)  {
						for j:=0;j<len(db[i].DataHash);j++{
							if db[i].DataHash[j]==bc[i].DataHash[j]{
								datamatch=append(datamatch,true)
							}else{
								datamatch=append(datamatch,false)
							}
						}
					}
					// else{
					// 	datamatch=append(datamatch,false)
					// }
					isMatch = append(isMatch, checkBoolArray(datamatch))
				}else{
					isMatch = append(isMatch, false)
				}

			}
		}else{
			isMatch = append(isMatch, false)
		}

	}

	return checkBoolArray(isMatch)
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
