package builder

import (
	"encoding/json"
	"fmt"
	"github.com/dileepaj/tracified-gateway/commons"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"

	"github.com/dileepaj/tracified-gateway/api/apiModel"

	"strings"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"

	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/stellar/go/build"
	"github.com/stellar/go/xdr"
)

/*SubmitMerge - WORKING MODEL
@author - Azeem Ashraf
@desc - Builds the TXN Type 6 for the gateway where it receives the user XDR 
and decodes it's contents and submit's to stellar and further maps the received TXN 
to Gateway Signed TXN's to maintain the profile, also records the activity in the gateway datastore.
@note - Should implement a validation layer to validate the contents of the XDR per builder before submission.
@params - ResponseWriter,Request
*/
func (AP *AbstractXDRSubmiter) SubmitMerge(w http.ResponseWriter, r *http.Request) {
	log.Debug("============================== SubmitMerge ==============================")
	var Done []bool
	Done = append(Done, true)
	object := dao.Connection{}

	var UserMergeTxnHashes []string
	var PreviousTxn string
	var MergeID string

	///HARDCODED CREDENTIALS
	publicKey := constants.PublicKey
	secretKey := constants.SecretKey
	// var result model.SubmitXDRResponse

	for i, TxnBody := range AP.TxnBody {

		var txe xdr.Transaction

		//decode the XDR
		err := xdr.SafeUnmarshalBase64(TxnBody.XDR, &txe)
		if err != nil {
			log.Error("Error while SafeUnmarshalBase64 @SubmitMerge " + err.Error())
		}

		//GET THE TYPE, IDENTIFIER, FROM IDENTIFERS, ITEM CODE AND ITEM AMOUNT FROM THE XDR
		AP.TxnBody[i].PublicKey = txe.SourceAccount.Address()
		AP.TxnBody[i].TxnType = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].Identifier = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[1].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].FromIdentifier1 = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[2].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].FromIdentifier2 = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[3].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].ItemCode = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[4].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].ItemAmount = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[5].Body.ManageDataOp.DataValue), "&")

		log.Debug(AP.TxnBody)
		//FOR THE MERGE FIRST BLOCK RETRIEVE THE PREVIOUS TXN FROM GATEWAY DB
		if i == 0 {
			p := object.GetLastTransactionbyIdentifier(AP.TxnBody[i].FromIdentifier1)
			p.Then(func(data interface{}) interface{} {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER
				result := data.(model.TransactionCollectionBody)
				PreviousTxn = result.TxnHash
				AP.TxnBody[i].PreviousTxnHash = result.TxnHash

				log.Debug(AP.TxnBody[i].PreviousTxnHash)
				return nil
			}).Catch(func(error error) error {
				log.Error("Error while GetLastTransactionbyIdentifier @SubmitMerge "+error.Error())
				///ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
				//DUE TO THE CHILD HAVING A NEW IDENTIFIER
				AP.TxnBody[i].PreviousTxnHash = ""
				return error
			})
			p.Await()
		}

		//SUBMIT THE FIRST XDR SIGNED BY THE USER
		display := stellarExecuter.ConcreteSubmitXDR{XDR: AP.TxnBody[i].XDR}
		result := display.SubmitXDR(AP.TxnBody[i].TxnType)
		UserMergeTxnHashes = append(UserMergeTxnHashes, result.TXNID)

		if result.Error.Code == 400 {
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
			var PreviousTXNBuilder build.ManageDataBuilder

			////GET THE PREVIOUS TRANSACTION FOR THE IDENTIFIER
			//INCASE OF FIRST MERGE BLOCK THE PREVIOUS IS TAKEN FROM IDENTIFIER
			//&
			//INCASE OF GREATER THAN ONE THE PREVIOUS TXN IS THE PREVIOUS MERGE
			if i == 0 {
				PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(PreviousTxn))
				AP.TxnBody[i].PreviousTxnHash = PreviousTxn
			} else {
				PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(PreviousTxn))
				AP.TxnBody[i].PreviousTxnHash = PreviousTxn
			}

			if i == 0 {
				p := object.GetLastTransactionbyIdentifier(TxnBody.FromIdentifier2)
				p.Then(func(data interface{}) interface{} {
					///ASSIGN PREVIOUS MANAGE DATA BUILDER
					result := data.(model.TransactionCollectionBody)
					MergeID = result.TxnHash
					AP.TxnBody[i].MergeID = result.TxnHash

					fmt.Println(AP.TxnBody[i].MergeID)
					return nil
				}).Catch(func(error error) error {
					log.Error("Error while GetLastTransactionbyIdentifier @SubmitMerge "+ error.Error())
					///ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
					//DUE TO THE CHILD HAVING A NEW IDENTIFIER
					AP.TxnBody[i].MergeID = ""
					return error
				})
				p.Await()
			}

			//BUILD THE GATEWAY XDR
			tx, err := build.Transaction(
				build.PublicNetwork,
				build.SourceAccount{publicKey},
				build.AutoSequence{commons.GetHorizonClient()},
				build.SetData("Type", []byte("G"+TxnBody.TxnType)),
				PreviousTXNBuilder,
				build.SetData("CurrentTXN", []byte(UserMergeTxnHashes[i])),
				build.SetData("MergeID", []byte(AP.TxnBody[i].MergeID)),
			)

			//SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
			GatewayTXE, err := tx.Sign(secretKey)
			if err != nil {
				log.Error("Error while build Transaction @SubmitMerge "+err.Error())
				AP.TxnBody[i].TxnHash = UserMergeTxnHashes[i]
				AP.TxnBody[i].Status = "Pending"

				///INSERT INTO TRANSACTION COLLECTION
				err2 := object.InsertTransaction(AP.TxnBody[i])
				if err2 != nil {
					log.Error("Error while InsertTransaction @SubmitMerge "+err2.Error())
				}
			}
			//CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
			txeB64, err := GatewayTXE.Base64()
			if err != nil {
				log.Error("Error while convert GatewayTXE to base64 @SubmitMerge "+err.Error())
				AP.TxnBody[i].TxnHash = UserMergeTxnHashes[i]
				AP.TxnBody[i].Status = "Pending"

				///INSERT INTO TRANSACTION COLLECTION
				err2 := object.InsertTransaction(AP.TxnBody[i])
				if err2 != nil {
					log.Error("Error while InsertTransaction @SubmitMerge "+err2.Error())
				}
			}

			//SUBMIT THE GATEWAY'S SIGNED XDR
			display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
			response1 := display1.SubmitXDR("G"+AP.TxnBody[i].TxnType)

			if response1.Error.Code == 400 {
				log.Error("Error got 400 from ConcreteSubmitXDR @SubmitMerge ")
				AP.TxnBody[i].TxnHash = UserMergeTxnHashes[i]
				AP.TxnBody[i].Status = "Pending"

				///INSERT INTO TRANSACTION COLLECTION
				err2 := object.InsertTransaction(AP.TxnBody[i])
				if err2 != nil {
					log.Error("Error while InsertTransaction @SubmitMerge "+err2.Error())
				}
			} else {
				//UPDATE THE TRANSACTION COLLECTION WITH TXN HASH
				AP.TxnBody[i].TxnHash = response1.TXNID
				if i == 0 {
					PreviousTxn = response1.TXNID
				}
				///INSERT INTO TRANSACTION COLLECTION
				err = object.InsertTransaction(AP.TxnBody[i])
				if err != nil {
					log.Error("Error while InsertTransaction @SubmitMerge "+err.Error())
				} else if i == 0 {

					var PreviousProfile string
					p := object.GetProfilebyIdentifier(AP.TxnBody[i].FromIdentifier1)
					p.Then(func(data interface{}) interface{} {

						result := data.(model.ProfileCollectionBody)
						PreviousProfile = result.ProfileTxn
						return nil
					}).Catch(func(error error) error {
						log.Error("Error while GetProfilebyIdentifier @SubmitMerge "+error.Error())
						PreviousProfile = ""
						return error
					})
					p.Await()

					Profile := model.ProfileCollectionBody{
						ProfileTxn:         response1.TXNID,
						ProfileID:          AP.TxnBody[i].ProfileID,
						Identifier:         AP.TxnBody[i].Identifier,
						PreviousProfileTxn: PreviousProfile,
						TriggerTxn:         UserMergeTxnHashes[i],
						TxnType:            AP.TxnBody[i].TxnType,
					}
					err3 := object.InsertProfile(Profile)
					if err3 != nil {
						log.Error("Error while InsertProfile @SubmitMerge "+err3.Error())
					}

				}

				// Done = true
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
