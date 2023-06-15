package builder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

/*SubmitTransfer - WORKING MODEL
@author - Azeem Ashraf
@desc - Builds the TXN Type 9 for the gateway where it receives the user XDR
and decodes it's contents and submit's to stellar and further maps the received TXN
to Gateway Signed TXN's to maintain the profile, also records the activity in the gateway datastore.
@note - Should implement a validation layer to validate the contents of the XDR per builder before submission.
@params - ResponseWriter,Request
*/
func (AP *AbstractXDRSubmiter) SubmitTransfer(w http.ResponseWriter, r *http.Request, NotOrphan bool) {
	
	log.Debug("---------------------------------- SubmitTransfer -------------------------------")
	var Done []bool
	Done = append(Done, NotOrphan)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	object := dao.Connection{}
	var UserTxnHashes []string

	///HARDCODED CREDENTIALS
	publicKey := constants.PublicKey
	secretKey := constants.SecretKey
	tracifiedAccount, err := keypair.ParseFull(secretKey)
	if err != nil {
		log.Error(err)
	}
	// netClient := commons.GetHorizonClient()
	// accountRequest := horizonclient.AccountRequest{AccountID: publicKey}
	// account, err := netClient.AccountDetail(accountRequest)

	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: publicKey}
	account, err := client.AccountDetail(ar)
	if err != nil {
		log.Fatal(err)
	}
	
	// var result model.SubmitXDRResponse
	var OrphanBoolArray []bool
	var OrphanBool bool

	for i, TxnBody := range AP.TxnBody {
		var txe xdr.Transaction
		//decode the XDR
		err := xdr.SafeUnmarshalBase64(TxnBody.XDR, &txe)
		if err != nil {
			log.Error("Error @SafeUnmarshalBase64 @SubmitTransfer "+err.Error())
		}
		//GET THE TYPE,IDENTIFIER,PREVIOUSSTAGE, CURRENTSTAGE, APPACCOUNT FROM THE XDR
		AP.TxnBody[i].Identifier = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[1].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].PublicKey = txe.SourceAccount.Address()
		AP.TxnBody[i].SequenceNo = int64(txe.SeqNum)
		AP.TxnBody[i].TxnType = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].Status = "pending"

		log.Debug(AP.TxnBody[i].Identifier)
		p := object.GetLastTransactionbyIdentifierAndTenantId(AP.TxnBody[i].Identifier, AP.TxnBody[i].TenantID, AP.TxnBody[i].ProductID)
		p.Then(func(data interface{}) interface{} {
			///ASSIGN PREVIOUS MANAGE DATA BUILDER
			result := data.(model.TransactionCollectionBody)
			AP.TxnBody[i].PreviousTxnHash = result.TxnHash
			AP.TxnBody[i].Orphan = false
			OrphanBoolArray = append(OrphanBoolArray, false)
			log.Info("Adopting from Orphanage! or Just Submitting")
			return nil
		}).Catch(func(error error) error {
			log.Error("Error @GetLastTransactionbyIdentifier @SubmitTransfer "+error.Error())
			log.Info("Sending to Orphanage!")
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
			err := object.InsertToOrphan(AP.TxnBody[i])
			if err != nil {
				log.Error("Error @InsertToOrphan @SubmitTransfer " + err.Error())
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
			result1 := display.SubmitXDR(AP.TxnBody[i].TxnType)
			UserTxnHashes = append(UserTxnHashes, result1.TXNID)

			if result1.Error.Code == 400 {
				log.Error("Error got 400 for SubmitXDR @SubmitTransfer ")
				Done = append(Done, false)
				if NotOrphan {
					w.WriteHeader(result1.Error.Code)
					response := apiModel.SubmitXDRSuccess{
						Status: "Index[" + strconv.Itoa(i) + "] TXN: Blockchain Transaction Failed!",
					}
					json.NewEncoder(w).Encode(response)
					return
				}

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

				//var PreviousTXNBuilder txnbuild.ManageDataBuilder
				//PreviousTXNBuilder = txnbuild.SetData("PreviousTXN", []byte(AP.TxnBody[i].PreviousTxnHash))
				PreviousTXNBuilder := txnbuild.ManageData{
					Name:  "PreviousTXN",
					Value:  []byte(AP.TxnBody[i].PreviousTxnHash),
				}
				TypeTxnBuilder := txnbuild.ManageData{
					Name:  "Type",
					Value:  []byte("G"+AP.TxnBody[i].TxnType),
				}

				CurrentTXNBuilder := txnbuild.ManageData{
					Name:  "CurrentTXN",
					Value:[]byte(UserTxnHashes[i]),
				}

				// BUILD THE GATEWAY XDR
				tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
					SourceAccount:        &account,
					IncrementSequenceNum: true,
					Operations:           []txnbuild.Operation{&PreviousTXNBuilder, &TypeTxnBuilder, &CurrentTXNBuilder},
					BaseFee:              constants.MinBaseFee,
					Memo:                 nil,
					Preconditions:       txnbuild.Preconditions{TimeBounds: constants.TransactionTimeOut},
				})

				if err != nil{
					log.Error("Error @ build.Transaction @SubmitTransfer " +err.Error())
				}

				//SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
				GatewayTXE, err := tx.Sign(commons.GetStellarNetwork(), tracifiedAccount)
				if err != nil {
					log.Error("Error @sign"+err.Error())
					AP.TxnBody[i].TxnHash = UserTxnHashes[i]
					AP.TxnBody[i].Status = "Pending"

					///INSERT INTO TRANSACTION COLLECTION
					err2 := object.InsertTransaction(AP.TxnBody[i])
					if err2 != nil {
						log.Error("Error @InsertTransaction @SubmitTransfer " +err2.Error())
					}
				}
				//CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
				txeB64, err := GatewayTXE.Base64()
				if err != nil {
					log.Error("Error @GatewayTXE.Base64 @SubmitTransfer "+err.Error())
					AP.TxnBody[i].TxnHash = UserTxnHashes[i]
					AP.TxnBody[i].Status = "Pending"

					///INSERT INTO TRANSACTION COLLECTION
					err2 := object.InsertTransaction(AP.TxnBody[i])
					if err2 != nil {
						log.Error("Error @InsertTransaction @SubmitTransfer " +err2.Error())
					}
				}

				//SUBMIT THE GATEWAY'S SIGNED XDR
				display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
				response1 := display1.SubmitXDR("G"+AP.TxnBody[i].TxnType)

				if response1.Error.Code == 400 {
					log.Error("Error got 400 for SubmitXDR @SubmitTransfer ")
					AP.TxnBody[i].TxnHash = UserTxnHashes[i]
					AP.TxnBody[i].Status = "Pending"

					///INSERT INTO TRANSACTION COLLECTION
					err2 := object.InsertTransaction(AP.TxnBody[i])
					if err2 != nil {
						log.Error("Error @InsertTransaction @SubmitTransfer " +err2.Error())
					}
				} else {
					//UPDATE THE TRANSACTION COLLECTION WITH TXN HASH
					AP.TxnBody[i].TxnHash = response1.TXNID
					AP.TxnBody[i].Status = "done"

					///INSERT INTO TRANSACTION COLLECTION
					err2 := object.InsertTransaction(AP.TxnBody[i])
					if err2 != nil {
						log.Error("Error @InsertTransaction @SubmitTransfer " +err2.Error())
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
		return
	}

}
