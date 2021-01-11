package builder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"

	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/xdr"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
)

/*SubmitXDR - WORKING MODEL
@author - Azeem Ashraf
@desc - Builds any with generic TXN Type for the gateway where it receives the user XDR
and decodes it's contents and submit's to stellar and further maps the received TXN
to Gateway Signed TXN's to maintain the profile, also records the activity in the gateway datastore.
@note - Should implement a validation layer to validate the contents of the XDR per builder before submission.
@params - ResponseWriter,Request
*/
func (AP *AbstractXDRSubmiter) SubmitXDR(w http.ResponseWriter, r *http.Request, NotOrphan bool) {
	var Done []bool
	Done = append(Done, NotOrphan)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	object := dao.Connection{}
	var UserTxnHashes []string

	///HARDCODED CREDENTIALS
	publicKey := constants.PublicKey
	secretKey := constants.SecretKey
	// var result model.SubmitXDRResponse
	var OrphanBoolArray []bool
	var OrphanBool bool

	for i, TxnBody := range AP.TxnBody {
		var txe xdr.Transaction
		//decode the XDR
		errx := xdr.SafeUnmarshalBase64(TxnBody.XDR, &txe)
		if errx != nil {

		}
		//GET THE TYPE AND IDENTIFIER FROM THE XDR
		AP.TxnBody[i].Identifier = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[1].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].PublicKey = txe.SourceAccount.Address()
		AP.TxnBody[i].TxnType = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].Status = "pending"

		fmt.Println(AP.TxnBody[i].Identifier)
		p := object.GetLastTransactionbyIdentifier(AP.TxnBody[i].Identifier)
		p.Then(func(data interface{}) interface{} {
			///ASSIGN PREVIOUS MANAGE DATA BUILDER
			result := data.(model.TransactionCollectionBody)
			AP.TxnBody[i].PreviousTxnHash = result.TxnHash
			AP.TxnBody[i].Orphan = false
			OrphanBoolArray = append(OrphanBoolArray, false)
			fmt.Println("Adopting from Orphanage! or Just Submitting")
			return nil
		}).Catch(func(error error) error {
			fmt.Println("Sending to Orphanage!")
			OrphanBoolArray = append(OrphanBoolArray, true)
			AP.TxnBody[i].Orphan = true
			return error
		})
		p.Await()
	}

	OrphanBool = checkBoolArray(OrphanBoolArray)

	if OrphanBool {
		for i, _ := range AP.TxnBody {
			//INSERT THE TXN INTO THE BUFFER
			err1 := object.InsertToOrphan(AP.TxnBody[i])
			if err1 != nil {
				Done = append(Done, false)
				w.WriteHeader(400)
				response := apiModel.SubmitXDRSuccess{
					Status: "Index[" + strconv.Itoa(i) + "] TXN: Orphanage Admission Revoked",
				}
				json.NewEncoder(w).Encode(response)
				return
			}
		}
	} else {
		for i, _ := range AP.TxnBody {
			//SUBMIT THE FIRST XDR SIGNED BY THE USER
			display := stellarExecuter.ConcreteSubmitXDR{XDR: AP.TxnBody[i].XDR}
			result1 := display.SubmitXDR(false,AP.TxnBody[i].TxnType)
			UserTxnHashes = append(UserTxnHashes, result1.TXNID)

			if result1.Error.Code == 400 {
				Done = append(Done, false)
				w.WriteHeader(result1.Error.Code)
				response := apiModel.SubmitXDRSuccess{
					Status: "Index[" + strconv.Itoa(i) + "] TXN: Blockchain Transaction Failed!",
				}
				json.NewEncoder(w).Encode(response)
				return
			}
		}
	}
	// for i, _ := range AP.TxnBody {

	// 	if AP.TxnBody[i].Orphan {

	// 		//INSERT THE TXN INTO THE BUFFER
	// 		err1 := object.InsertToOrphan(AP.TxnBody[i])
	// 		if err1 != nil {
	// 			Done = append(Done, false)
	// 			w.WriteHeader(400)
	// 			response := apiModel.SubmitXDRSuccess{
	// 				Status: "Index[" + strconv.Itoa(i) + "] TXN: Orphanage Admission Revoked",
	// 			}
	// 			json.NewEncoder(w).Encode(response)
	// 			return
	// 		}
	// 	} else {
	// 		//SUBMIT THE FIRST XDR SIGNED BY THE USER
	// 		display := stellarExecuter.ConcreteSubmitXDR{XDR: AP.TxnBody[i].XDR}
	// 		result1 := display.SubmitXDR()
	// 		UserTxnHashes = append(UserTxnHashes, result1.TXNID)

	// 		if result1.Error.Code == 400 {
	// 			Done = append(Done, false)
	// 			w.WriteHeader(result1.Error.Code)
	// 			response := apiModel.SubmitXDRSuccess{
	// 				Status: "Index[" + strconv.Itoa(i) + "] TXN: Blockchain Transaction Failed!",
	// 			}
	// 			json.NewEncoder(w).Encode(response)
	// 			return
	// 		}
	// 	}
	// }

	go func() {
		for i, TxnBody := range AP.TxnBody {
			if !TxnBody.Orphan {

				var PreviousTXNBuilder build.ManageDataBuilder
				PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(AP.TxnBody[i].PreviousTxnHash))

				//BUILD THE GATEWAY XDR
				tx, err := build.Transaction(
					build.TestNetwork,
					build.SourceAccount{publicKey},
					build.AutoSequence{horizon.DefaultTestNetClient},
					build.SetData("Type", []byte("G"+AP.TxnBody[i].TxnType)),
					PreviousTXNBuilder,
					build.SetData("CurrentTXN", []byte(UserTxnHashes[i])),
				)

				//SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
				GatewayTXE, err := tx.Sign(secretKey)
				if err != nil {
					AP.TxnBody[i].TxnHash = UserTxnHashes[i]
					AP.TxnBody[i].Status = "Pending"

					///INSERT INTO TRANSACTION COLLECTION
					err2 := object.InsertTransaction(AP.TxnBody[i])
					if err2 != nil {

					}
				}
				//CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
				txeB64, err := GatewayTXE.Base64()
				if err != nil {
					AP.TxnBody[i].TxnHash = UserTxnHashes[i]
					AP.TxnBody[i].Status = "Pending"

					///INSERT INTO TRANSACTION COLLECTION
					err2 := object.InsertTransaction(AP.TxnBody[i])
					if err2 != nil {

					}
				}

				//SUBMIT THE GATEWAY'S SIGNED XDR
				display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
				response1 := display1.SubmitXDR(false,"G"+AP.TxnBody[i].TxnType)

				if response1.Error.Code == 400 {
					AP.TxnBody[i].TxnHash = UserTxnHashes[i]
					AP.TxnBody[i].Status = "Pending"

					///INSERT INTO TRANSACTION COLLECTION
					err2 := object.InsertTransaction(AP.TxnBody[i])
					if err2 != nil {

					}
				} else {
					//UPDATE THE TRANSACTION COLLECTION WITH TXN HASH
					AP.TxnBody[i].TxnHash = response1.TXNID
					AP.TxnBody[i].Status = "done"

					///INSERT INTO TRANSACTION COLLECTION
					err2 := object.InsertTransaction(AP.TxnBody[i])
					if err2 != nil {

					}
				}
			}
		}
	}()

	if checkBoolArray(Done) {
		w.WriteHeader(http.StatusOK)
		result := apiModel.SubmitXDRSuccess{
			Status: "Success",
		}
		json.NewEncoder(w).Encode(result)
	}

	return
}
