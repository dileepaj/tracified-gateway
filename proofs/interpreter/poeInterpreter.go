package interpreter

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/retriever/stellarRetriever"
	"github.com/stellar/go/support/log"
)

type POEInterface interface {
	RetrievePOE() model.RetrievePOE
}

type AbstractPOE struct {
	POEStruct apiModel.POEStruct
	// Txn       string
	// ProfileID string
	// Hash      string
}

/*InterpretPOE - Working Model
@author - Azeem Ashraf
@desc - Interprets All the required fields necessary to perform POE
@params - ResponseWriter,Request
*/
func (AP *AbstractPOE) InterpretPOE(tdpid string) model.POE {
	var poeObj model.POE
	fmt.Println("Inside POE verfication with tdp", tdpid)
	object := stellarRetriever.ConcretePOE{POEStruct: AP.POEStruct}

	poeObj.RetrievePOE = object.RetrievePOE()

	if poeObj.RetrievePOE.BCHash == "" {
		return poeObj
	} else {
		poeObj.RetrievePOE.Error = matchingHash(
			poeObj.RetrievePOE.BCHash,
			poeObj.RetrievePOE.DBHash,
			AP.POEStruct.ProfileID,
			poeObj.RetrievePOE.Identifier,
			tdpid,
		)
		return poeObj
	}
}

func (AP *AbstractPOE) InterpretPOAC(tdpid string) model.POE {
	var poeObj model.POE

	object := stellarRetriever.ConcretePOE{POEStruct: AP.POEStruct}

	poeObj.RetrievePOE = object.RetrievePOE()

	if poeObj.RetrievePOE.BCHash == "" {
		return poeObj
	} else {
		poeObj.RetrievePOE.Error = matchingHashToCheckPOAC(
			poeObj.RetrievePOE.BCHash,
			poeObj.RetrievePOE.DBHash,
			AP.POEStruct.ProfileID,
			poeObj.RetrievePOE.Identifier,
			tdpid,
		)
		return poeObj
	}
}

func matchingHashToCheckPOAC(bcHash string, dbHash string, bcProfile string, dbProfile string, tdpid string) model.Error {
	var Rerr model.Error
	fmt.Println("bcHash", "dbHash inside POAC matching hash")
	fmt.Println(bcHash, dbHash)
	if strings.ToUpper(bcHash) == strings.ToUpper(dbHash) {
		Rerr.Code = http.StatusOK
		Rerr.Message = "BC Hash & DB Hash match."

		// if bcProfile == dbProfile {
		// 	Rerr.Code = http.StatusOK
		// 	Rerr.Message = "POE Success! DB & BC Hash and Profile match."
		// } else {
		// 	Rerr.Code = http.StatusOK
		// 	Rerr.Message = "BC Profile & DB Profile didn't match."
		// }
		return Rerr

	} else {
		Rerr.Code = http.StatusOK
		Rerr.Message = "Error! BC Hash & DB Hash din't match."
		return Rerr
	}
}

func matchingHash(bcHash string, dbHash string, bcProfile string, dbProfile string, tdpid string) model.Error {
	var Rerr model.Error
	fmt.Println("bcHash", "dbHash")
	fmt.Println(bcHash, dbHash)
	if strings.ToUpper(bcHash) == strings.ToUpper(dbHash) {
		Rerr.Code = http.StatusOK
		Rerr.Message = "BC Hash & DB Hash match."

		// if bcProfile == dbProfile {
		// 	Rerr.Code = http.StatusOK
		// 	Rerr.Message = "POE Success! DB & BC Hash and Profile match."
		// } else {
		// 	Rerr.Code = http.StatusOK
		// 	Rerr.Message = "BC Profile & DB Profile didn't match."
		// }
		return Rerr

	} else {
		if tdpid != "" {
			var result model.TransactionCollectionBody
			object := dao.Connection{}
			p := object.GetTransactionForTdpId(tdpid)
			p.Then(func(data interface{}) interface{} {
				result = data.(model.TransactionCollectionBody)
				return nil
			}).Catch(func(error error) error {
				log.Error("Error while GetTransactionForTdpId " + error.Error())
				response := model.Error{Message: "TDPID NOT FOUND IN DATASTORE"}
				fmt.Println(response)
				return error
			}).Await()
			poeStructObj := apiModel.POEStruct{Txn: result.TxnHash, Hash: dbHash}
			display := AbstractPOE{POEStruct: poeStructObj}
			display.InterpretPOAC("")
			Rerr.Code = http.StatusOK
			Rerr.Message = "POE failed but POAC is present"
			return Rerr

		} else {
			Rerr.Code = http.StatusOK
			Rerr.Message = "Error! BC Hash & DB Hash din't match."
			return Rerr
		}
	}
}
