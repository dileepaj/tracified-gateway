package interpreter

import (
	"github.com/Tracified-Gateway/model"
)

type POCInterface interface {
	RetrievePOC() model.RetrievePOC
}

type AbstractPOC struct {
}

func (AP *AbstractPOC) InterpretPOC(POCInterface POCInterface) model.POC {
	var pocObj model.POC

	pocObj.RetrievePOC = POCInterface.RetrievePOC()

	if pocObj.RetrievePOC.BCHash == "" {
		return pocObj
	} else {
		pocObj.RetrievePOC.Error = MatchingHash(pocObj.RetrievePOC.BCHash, pocObj.RetrievePOC.DBHash)
		return pocObj
	}
}
