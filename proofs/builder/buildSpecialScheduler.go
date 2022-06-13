package builder

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	// "github.com/dileepaj/tracified-gateway/model"
	// "github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"

	// "github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"

	// "github.com/stellar/go/build"
	// "github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/xdr"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
)

/*SubmitSpecial - EXPERIMENTAL
@author - Azeem Ashraf
@desc -
//get genesis and tdp transactions and push to temp orphan
//also add then sequence numebr to the database
@note -
@params - ResponseWriter,Request
*/
type Identifier struct {
	Id   string
	Type string
}

func (AP *AbstractXDRSubmiter) SubmitSpecial(w http.ResponseWriter, r *http.Request) {

	// log.Debug("------------------------- SubmitSpecial --------------------------")
	fmt.Println("-------------Submit Special-------------------------")
	var Done []bool           //array to decide whether the actions are done
	Done = append(Done, true) //starting with a true for bipass
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var id apiModel.IdentifierModel
	object := dao.Connection{}
	// var UserTxnHashes []string

	// ///HARDCODED CREDENTIALS
	// publicKey := constants.PublicKey
	// secretKey := constants.SecretKey

	for i, TxnBody := range AP.TxnBody {
		var txe xdr.Transaction
		log.Debug("index")
		log.Debug(i)
		log.Debug("TxnBody.XDR")
		log.Debug(TxnBody.XDR)
		//decode the XDR
		err := xdr.SafeUnmarshalBase64(TxnBody.XDR, &txe)
		if err != nil {
			log.Error("Error @ SafeUnmarshalBase64 @SubmitSpecial " + err.Error())
		}
		AP.TxnBody[i].PublicKey = txe.SourceAccount.Address()
		AP.TxnBody[i].SequenceNo = int64(txe.SeqNum)
		//GET THE TYPE AND IDENTIFIER FROM THE XDR
		// AP.TxnBody[i].Identifier = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[3].Body.ManageDataOp.DataValue), "&")
		// AP.TxnBody[i].ProductName = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[1].Body.ManageDataOp.DataValue), "&")
		// AP.TxnBody[i].TxnType = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
		// AP.TxnBody[i].DataHash = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[4].Body.ManageDataOp.DataValue), "&")

		fmt.Println("Source Account ", AP.TxnBody[i].PublicKey)
		//fmt.Println(" =======================AP TXN body ", &AP.TxnBody[i])
		
		//fmt.Println("************MAPXDR operations",stellarRetriever.MapXDROperations(&AP.TxnBody[i], txe.Operations))

		AP.TxnBody[i].Status = "pending"
		// AP.TxnBody[i].DataHash = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[4].Body.ManageDataOp.DataValue), "&")

		log.Debug(AP.TxnBody[i].Identifier)
		err2 := object.InsertSpecialToTempOrphan(AP.TxnBody[i])
		if err2 != nil {
			fmt.Println("Error @ InsertSpecialToTempOrphan @SubmitSpecial " + err2.Error())
			Done = append(Done, false)
			w.WriteHeader(http.StatusBadRequest)
			response := apiModel.SubmitXDRSuccess{
				Status: "Index[" + strconv.Itoa(i) + "] TXN: Scheduling Failed!",
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		if AP.TxnBody[i].TxnType == "0" {
			rawDecodedText, err := base64.StdEncoding.DecodeString(TxnBody.Identifier)
			if err != nil {
				panic(err)
			}

			var jsonID Identifier
			json.Unmarshal([]byte(rawDecodedText), &jsonID)
			id.MapValue = AP.TxnBody[i].Identifier
			id.Identifier = jsonID.Id
			id.Type = jsonID.Type

			err3 := object.InsertIdentifier(id)
			if err3 != nil {
				fmt.Println("identifier map failed" + err3.Error())
			}
		}
	}

	if checkBoolArray(Done) {
		w.WriteHeader(http.StatusOK)
		result := apiModel.SubmitXDRSuccess{
			Status: "Success",
		}
		json.NewEncoder(w).Encode(result)
		return
	}

}

//SubmitSpecialData - EXPERIMENTAL

func (AP *AbstractXDRSubmiter) SubmitSpecialData(w http.ResponseWriter, r *http.Request) {

	log.Debug("------------------------------------ SubmitSpecialData ----------------------------------")

	var Done []bool           //array to decide whether the actions are done
	Done = append(Done, true) //starting with a true for bipass
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	object := dao.Connection{}
	// var UserTxnHashes []string

	// ///HARDCODED CREDENTIALS
	// publicKey := constants.PublicKey
	// secretKey := constants.SecretKey

	for i, TxnBody := range AP.TxnBody {
		var txe xdr.Transaction
		//decode the XDR
		err := xdr.SafeUnmarshalBase64(TxnBody.XDR, &txe)
		if err != nil {
			log.Error("Error @ SafeUnmarshalBase64 @SubmitSpecialData " + err.Error())
		}
		//GET THE TYPE AND IDENTIFIER FROM THE XDR
		AP.TxnBody[i].Identifier = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[1].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].PublicKey = txe.SourceAccount.Address()
		AP.TxnBody[i].SequenceNo = int64(txe.SeqNum)
		AP.TxnBody[i].TxnType = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].DataHash = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[2].Body.ManageDataOp.DataValue), "&")

		AP.TxnBody[i].Status = "pending"

		fmt.Println(AP.TxnBody[i].Identifier)
		err = object.InsertSpecialToTempOrphan(AP.TxnBody[i])
		if err != nil {
			log.Error("Error @ InsertSpecialToTempOrphan @SubmitSpecialData " + err.Error())
			Done = append(Done, false)
			w.WriteHeader(http.StatusBadRequest)
			response := apiModel.SubmitXDRSuccess{
				Status: "Index[" + strconv.Itoa(i) + "] TXN: Scheduling Failed!",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	if checkBoolArray(Done) {
		w.WriteHeader(http.StatusOK)
		result := apiModel.SubmitXDRSuccess{
			Status: "Success",
		}
		json.NewEncoder(w).Encode(result)
		return
	}

}

//SubmitSpecialTransfer - EXPERIMENTAL

func (AP *AbstractXDRSubmiter) SubmitSpecialTransfer(w http.ResponseWriter, r *http.Request) {
	log.Debug("-------------------------------- @SubmitSpecialTransfer ----------------------------------")
	var Done []bool           //array to decide whether the actions are done
	Done = append(Done, true) //starting with a true for bipass
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	object := dao.Connection{}
	// var UserTxnHashes []string

	// ///HARDCODED CREDENTIALS
	// publicKey := constants.PublicKey
	// secretKey := constants.SecretKey

	for i, TxnBody := range AP.TxnBody {
		var txe xdr.Transaction
		//decode the XDR
		err := xdr.SafeUnmarshalBase64(TxnBody.XDR, &txe)
		if err != nil {
			log.Error("Error @ SafeUnmarshalBase64 @SubmitSpecialTransfer " + err.Error())
		}
		//GET THE TYPE AND IDENTIFIER FROM THE XDR
		AP.TxnBody[i].Identifier = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[1].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].PublicKey = txe.SourceAccount.Address()
		AP.TxnBody[i].SequenceNo = int64(txe.SeqNum)
		AP.TxnBody[i].TxnType = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].PreviousStage = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[2].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].CurrentStage = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[3].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].AppAccount = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[4].Body.ManageDataOp.DataValue), "&")

		AP.TxnBody[i].Status = "pending"

		log.Debug(AP.TxnBody[i].Identifier)
		err = object.InsertSpecialToTempOrphan(AP.TxnBody[i])
		if err != nil {
			log.Error("Error @ InsertSpecialToTempOrphan @SubmitSpecialTransfer " + err.Error())
			Done = append(Done, false)
			w.WriteHeader(http.StatusBadRequest)
			response := apiModel.SubmitXDRSuccess{
				Status: "Index[" + strconv.Itoa(i) + "] TXN: Scheduling Failed!",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	if checkBoolArray(Done) {
		w.WriteHeader(http.StatusOK)
		result := apiModel.SubmitXDRSuccess{
			Status: "Success",
		}
		json.NewEncoder(w).Encode(result)
		return
	}

}
