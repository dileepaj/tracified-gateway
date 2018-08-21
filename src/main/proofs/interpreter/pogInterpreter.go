package interpreter

import (
	"fmt"
	"main/model"
	"main/proofs/retriever/stellarRetriever"
)

type POGInterface interface {
	RetrievePOG() model.RetrievePOG
}

type AbstractPOG struct {
	LastTxn    string
	POGTxn     string
	Identifier string
}

func (AP *AbstractPOG) InterpretPOG() model.POG {
	var pogObj model.POG

	object := stellarRetriever.ConcretePOG{LastTxn: AP.LastTxn}

	pogObj.RetrievePOG = object.RetrievePOG()

	fmt.Println("POGInterpreter")
	fmt.Println(pogObj.RetrievePOG)

	if pogObj.RetrievePOG.Error.Message == "Txn Hash retrieved from the blockchain." {
		if pogObj.RetrievePOG.CurTxn != AP.POGTxn {
			pogObj.RetrievePOG.Error.Message = "Proof of Genesis Failed, Genesis Txn hashes didn't match!"
			return pogObj
		} else if pogObj.RetrievePOG.Identifier != AP.Identifier {
			pogObj.RetrievePOG.Error.Message = "Proof of Genesis Failed, Genesis Identifier hash didn't match!"
			return pogObj
		} else if pogObj.RetrievePOG.PreTxn != "" {
			pogObj.RetrievePOG.Error.Message = "Proof of Genesis Failed, Genesis previousTxn ID is not empty!"
			return pogObj
		} else {
			pogObj.RetrievePOG.Error.Message = "Proof of Genesis Success!"
			return pogObj
		}
	} else {
		return pogObj
	}
}
