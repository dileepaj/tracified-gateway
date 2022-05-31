package builder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/commons"
	log "github.com/sirupsen/logrus"

	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
)

/*SubmitGenesis - WORKING MODEL
@author - Azeem Ashraf
@desc - Builds the TXN Type 0 for the gateway where it receives the user XDR
and decodes it's contents and submit's to stellar and further maps the received TXN
to Gateway Signed TXN's to maintain the profile, also records the activity in the gateway datastore.
@note - Should implement a validation layer to validate the contents of the XDR per builder before submission.
@params - ResponseWriter,Request
*/
func (AP *AbstractXDRSubmiter) SubmitGenesis(w http.ResponseWriter, r *http.Request) {
	log.Debug("=========================== buildGenesis.go - SubmitGenesis =============================")
	var Done []bool
	Done = append(Done, true)
	netClient := commons.GetHorizonClient()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	object := dao.Connection{}
	var UserTxnHashes []string

	///HARDCODED CREDENTIALS
	publicKey := constants.PublicKey
	secretKey := constants.SecretKey
	// var result model.SubmitXDRResponse

	for i, TxnBody := range AP.TxnBody {
		var txe xdr.Transaction
		//decode the XDR
		err := xdr.SafeUnmarshalBase64(TxnBody.XDR, &txe)
		if err != nil {
			log.Error("Error while SafeUnmarshalBase64 @SubmitGenesis " + err.Error())
		}
		//GET THE TYPE AND IDENTIFIER FROM THE XDR
		AP.TxnBody[i].Identifier = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[1].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].PublicKey = txe.SourceAccount.Address()
		AP.TxnBody[i].SequenceNo = int64(txe.SeqNum)
		AP.TxnBody[i].TxnType = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].Status = "pending"

		//SUBMIT THE FIRST XDR SIGNED BY THE USER
		display := stellarExecuter.ConcreteSubmitXDR{XDR: AP.TxnBody[i].XDR}
		result := display.SubmitXDR(AP.TxnBody[i].TxnType)
		UserTxnHashes = append(UserTxnHashes, result.TXNID)

		if result.Error.Code == 400 {
			log.Error("Error while ConcreteSubmitXDR in SubmitGenesis @SubmitGenesis ")
			Done = append(Done, false)
			w.WriteHeader(result.Error.Code)
			response := apiModel.SubmitXDRSuccess{
				Status: "Index[" + strconv.Itoa(i) + "] TXN: Blockchain Transaction Failed!",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	go func() {

		for i, TxnBody := range AP.TxnBody {

			pubaccountRequest := horizonclient.AccountRequest{AccountID: publicKey}
			pubaccount, err := netClient.AccountDetail(pubaccountRequest)			

			//var PreviousTXNBuilder txnbuild.ManageData
			//PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(AP.TxnBody[i].PreviousTxnHash))
			PreviousTXNBuilder := txnbuild.ManageData{
				Name: "PreviousTXN",
				Value: []byte(AP.TxnBody[i].PreviousTxnHash),
			}

			TypeTxnBuilder := txnbuild.ManageData{
				Name: "Type",
				Value: []byte("G"+TxnBody.TxnType),
			}

			CurrentTXNBuilder := txnbuild.ManageData{
				Name: "CurrentTXN",
				Value: []byte(UserTxnHashes[i]),
			}

			//BUILD THE GATEWAY XDR
			// tx, err := build.Transaction(
			// 	commons.GetHorizonNetwork(),
			// 	build.SourceAccount{publicKey},
			// 	build.AutoSequence{commons.GetHorizonClient()},
			// 	build.SetData("Type", []byte("G"+TxnBody.TxnType)),
			// 	PreviousTXNBuilder,
			// 	build.SetData("CurrentTXN", []byte(UserTxnHashes[i])),
			// )

			tx, err := txnbuild.NewTransaction(
				txnbuild.TransactionParams{
					SourceAccount: &pubaccount,
					IncrementSequenceNum: true,
					Operations: []txnbuild.Operation{&PreviousTXNBuilder, &TypeTxnBuilder, &CurrentTXNBuilder},
					BaseFee: txnbuild.MinBaseFee,
					Preconditions: txnbuild.Preconditions{},
				},
			)

			if err != nil{
				log.Error("Error while build Transaction @SubmitGenesis " + err.Error())
			}

			//SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
			GatewayTXE, err := tx.Sign(secretKey)
			if err != nil {
				log.Error("Error while Sign @SubmitGenesis " + err.Error())
				AP.TxnBody[i].TxnHash = UserTxnHashes[i]
				AP.TxnBody[i].Status = "Pending"
				///INSERT INTO TRANSACTION COLLECTION
				err = object.InsertTransaction(AP.TxnBody[i])
				if err != nil {
					log.Error("Error while InsertTransaction @SubmitGenesis " + err.Error())
				}
			}

			//CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
			txeB64, err := GatewayTXE.Base64()
			if err != nil {
				log.Error("Error while converting gatewayTXE to base64 @SubmitGenesis " + err.Error())
				AP.TxnBody[i].TxnHash = UserTxnHashes[i]
				AP.TxnBody[i].Status = "Pending"

				///INSERT INTO TRANSACTION COLLECTION
				err = object.InsertTransaction(AP.TxnBody[i])
				if err != nil {
					log.Error("Error while InsertTransaction @SubmitGenesis " + err.Error())
				}
			}

			//SUBMIT THE GATEWAY'S SIGNED XDR
			display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
			response1 := display1.SubmitXDR("G"+AP.TxnBody[i].TxnType)

			if response1.Error.Code == 400 {
				log.Error("Error while SubmitXDR in ConcreteSubmitXDR @SubmitGenesis ")
				AP.TxnBody[i].TxnHash = UserTxnHashes[i]
				AP.TxnBody[i].Status = "Pending"

				///INSERT INTO TRANSACTION COLLECTION
				err = object.InsertTransaction(AP.TxnBody[i])
				if err != nil {
					log.Error("Error while InsertTransaction @SubmitGenesis " +err.Error())
				}
			} else {
				//UPDATE THE TRANSACTION COLLECTION WITH TXN HASH
				AP.TxnBody[i].TxnHash = response1.TXNID
				AP.TxnBody[i].Status = "done"

				///INSERT INTO TRANSACTION COLLECTION
				err = object.InsertTransaction(AP.TxnBody[i])
				if err != nil {
					log.Error("Error while InsertTransaction @SubmitGenesis " +err.Error())
				} else {
					// var PreviousProfile string
					// p := object.GetProfilebyIdentifier(AP.TxnBody[i].Identifier)
					// p.Then(func(data interface{}) interface{} {

					// 	result := data.(model.ProfileCollectionBody)
					// 	PreviousProfile = result.ProfileTxn
					// 	return nil
					// }).Catch(func(error error) error {
					// 	PreviousProfile = ""
					// 	return nil
					// })
					// p.Await()

					// Profile := model.ProfileCollectionBody{
					// 	ProfileTxn:         response1.TXNID,
					// 	ProfileID:          AP.TxnBody[i].ProfileID,
					// 	Identifier:         AP.TxnBody[i].Identifier,
					// 	PreviousProfileTxn: PreviousProfile,
					// 	TriggerTxn:         UserTxnHashes[i],
					// 	TxnType:            AP.TxnBody[i].TxnType,
					// }
					// err3 := object.InsertProfile(Profile)
					// if err3 != nil {

					// }

				}
			}
		}

		//ORPHAN TXNS TO BE COLLECTED HERE TO BE CALLED IN AGAIN
		var Orphans []model.TransactionCollectionBody
		for _, TxnBody := range AP.TxnBody {
			p := object.GetOrphanbyIdentifier(TxnBody.Identifier)
			p.Then(func(data interface{}) interface{} {

				result := data.(model.TransactionCollectionBody)
				Orphans = append(Orphans, result)
				err := object.RemoveFromOrphanage(TxnBody.Identifier)
				if err != nil {
					log.Error("Error while RemoveFromOrphanage @SubmitGenesis " +err.Error())
				}
				return nil
			}).Catch(func(error error) error {
				log.Error("Error while GetOrphanbyIdentifier @SubmitGenesis ")
				// return error
				return nil
			})
			p.Await()
		}

		display := AbstractXDRSubmiter{TxnBody: Orphans}
		display.SubmitData(w, r,false)
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
