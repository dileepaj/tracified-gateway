package interpreter

import (
	"main/model"
	"main/proofs/retriever/stellarRetriever"
	"net/http"
)

type POEInterface interface {
	RetrievePOE() model.RetrievePOE
}

type AbstractPOE struct {
	Txn       string
	ProfileID string
	Hash      string
}

func (AP *AbstractPOE) InterpretPOE() model.POE {
	var poeObj model.POE

	object := stellarRetriever.ConcretePOE{Txn: AP.Txn, ProfileID: AP.ProfileID, Hash: AP.Hash}

	poeObj.RetrievePOE = object.RetrievePOE()

	if poeObj.RetrievePOE.BCHash == "" {
		return poeObj
	} else {
		poeObj.RetrievePOE.Error = MatchingHash(poeObj.RetrievePOE.BCHash, poeObj.RetrievePOE.DBHash)
		return poeObj
	}
}

func MatchingHash(bcHash string, dbHash string) model.Error {
	var Rerr model.Error

	if bcHash == dbHash {
		Rerr.Code = http.StatusOK
		Rerr.Message = "BC Hash & DB Hash match."
		return Rerr
	} else {
		Rerr.Code = http.StatusOK
		Rerr.Message = "Error! BC Hash & DB Hash din't match."
		return Rerr
	}
}
