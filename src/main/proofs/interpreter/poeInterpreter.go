package interpreter

import (
	"fmt"
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
		poeObj.RetrievePOE.Error = MatchingHash(
			poeObj.RetrievePOE.BCHash,
			poeObj.RetrievePOE.DBHash,
			AP.ProfileID,
			poeObj.RetrievePOE.BCProfile)
		return poeObj
	}
}

func MatchingHash(bcHash string, dbHash string, bcProfile string, dbProfile string) model.Error {
	var Rerr model.Error
	fmt.Println("bcHash", "dbHash")
	fmt.Println(bcHash, dbHash)
	if bcHash == dbHash {
		Rerr.Code = http.StatusOK
		Rerr.Message = "BC Hash & DB Hash match."
		if bcProfile == dbProfile {
			Rerr.Code = http.StatusOK
			Rerr.Message = "POE Success! DB & BC Hash and Profile match."
		} else {
			Rerr.Code = http.StatusOK
			Rerr.Message = "BC Profile & DB Profile didn't match."
		}
		return Rerr

	} else {
		Rerr.Code = http.StatusOK
		Rerr.Message = "Error! BC Hash & DB Hash din't match."
		return Rerr
	}
}
