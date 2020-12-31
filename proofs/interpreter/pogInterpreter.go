package interpreter

import (
	"fmt"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/retriever/stellarRetriever"
)

type POGInterface interface {
	RetrievePOG() model.RetrievePOG
}

type AbstractPOG struct {
	POGStruct apiModel.POGStruct
	// LastTxn    string
	// POGTxn     string
	// Identifier string
}

/*InterpretPOG - Working Model
@author - Azeem Ashraf
@desc - Interprets All the required fields necessary to perform POG
@params - ResponseWriter,Request
*/
func (AP *AbstractPOG) InterpretPOG() model.POG {
	var pogObj model.POG

	object := stellarRetriever.ConcretePOG{POGStruct: AP.POGStruct}

	pogObj.RetrievePOG = object.RetrievePOG()

	fmt.Println("POGInterpreter")
	fmt.Println(pogObj.RetrievePOG)

	if pogObj.RetrievePOG.Message.Message == "Txn Hash retrieved from the blockchain." {
		if pogObj.RetrievePOG.CurTxn != AP.POGStruct.POGTxn {
			pogObj.RetrievePOG.Message.Message = "Proof of Genesis Failed, Genesis Txn hashes didn't match!"
			return pogObj
		} else if pogObj.RetrievePOG.Identifier != AP.POGStruct.Identifier {
			pogObj.RetrievePOG.Message.Message = "Proof of Genesis Failed, Genesis Identifier hash didn't match!"
			return pogObj
		} else if pogObj.RetrievePOG.PreTxn != "" {
			if pogObj.RetrievePOG.PreTxn != "0" {
				pogObj.RetrievePOG.Message.Message = "Proof of Genesis Failed, Genesis previousTxn ID is not empty!"
				return pogObj
			} else {
				pogObj.RetrievePOG.Message.Message = "Proof of Genesis Success!"
				return pogObj
			}

		} else {
			pogObj.RetrievePOG.Message.Message = "Proof of Genesis Success!"
			return pogObj
		}
	} else {
		return pogObj
	}
}
