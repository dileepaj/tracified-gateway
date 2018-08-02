package interpreter

import (
	// "fmt"
	// "github.com/stellar/go/support/http"
	// "fmt"
	"github.com/tracified-gateway/model"
	// "main/api/apiModel"
)

type POCInterface interface {
	RetrievePOC() model.RetrievePOC
}

type AbstractPOC struct {
}

func (AP *AbstractPOC) InterpretPOC(POCInterface POCInterface) model.POC {
	var pocObj model.POC

	pocObj.RetrievePOC = POCInterface.RetrievePOC()
	// fmt.Println(pocObj.RetrievePOC.BCHash)
	// fmt.Println(pocObj.RetrievePOC.BCHash)
	// fmt.Println(pocObj.RetrievePOC.BCHash.Previous)
	isMapped := testCompare(pocObj.RetrievePOC.DBHash, pocObj.RetrievePOC.BCHash)
	// isMapped := true

	// var returnRes model.POC
	if isMapped == true {
		pocObj.RetrievePOC.Error.Message = "Chain Exists in the Blockchain"
		// pocObj.RetrievePOC.Txn= "Chain Exists"
		pocObj.RetrievePOC.Error.Code = 200
		// returnRes := pocObj.RetrievePOC
	} else {
		pocObj.RetrievePOC.Error.Message = "Chain Doesn't Exist in the Blockchain"
		// pocObj.RetrievePOC.Current
		pocObj.RetrievePOC.Error.Code = 200
	}

	return pocObj
}

func testCompare(db []model.Current, bc []model.Current) bool {
	isMatch := []bool{}
	if db != nil && bc != nil {
		if len(db) == len(bc) {
			for i := 0; i < len(db); i++ {
				if db[i].TDPID == bc[i].TDPID && db[i].Hash == bc[i].Hash {
					isMatch = append(isMatch, true)
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
