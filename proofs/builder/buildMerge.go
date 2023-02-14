package builder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/commons"

	"strings"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"

	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
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
	var id apiModel.IdentifierModel
	var UserMergeTxnHashes []string
	//netClient := commons.GetHorizonClient()
	// var MergeID string

	///HARDCODED CREDENTIALS
	publicKey := constants.PublicKey
	secretKey := constants.SecretKey
	tracifiedAccount, err := keypair.ParseFull(secretKey)
	if err != nil {
		log.Error(err)
	}
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
		AP.TxnBody[i].SequenceNo = int64(txe.SeqNum)
		AP.TxnBody[i].TxnType = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].Identifier = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[1].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].FromIdentifier1 = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[2].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].FromIdentifier2 = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[3].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].ItemCode = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[4].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].ItemAmount = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[5].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].AppAccount = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[6].Body.ManageDataOp.DataValue), "&")
		if len(txe.Operations)==8{
			AP.TxnBody[i].ProductName = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[7].Body.ManageDataOp.DataValue), "&")
		}

		log.Debug(AP.TxnBody)
		//FOR THE MERGE FIRST BLOCK RETRIEVE THE PREVIOUS TXN FROM GATEWAY DB
		if AP.TxnBody[i].Identifier != AP.TxnBody[i].FromIdentifier1 {
			pData, errorAsync := object.GetLastTransactionbyIdentifierNotSplitParent(AP.TxnBody[i].FromIdentifier1).Then(func(data interface{}) interface{} {
				return data
			}).Await()

			if errorAsync != nil || pData == nil {
				log.Error("Error while GetLastTransactionbyIdentifier @SubmitMerge " + errorAsync.Error())
				///ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
				//DUE TO THE CHILD HAVING A NEW IDENTIFIER
				AP.TxnBody[i].PreviousTxnHash = ""
			} else {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER
				result := pData.(model.TransactionCollectionBody)
				AP.TxnBody[i].PreviousTxnHash = result.TxnHash
				log.Debug(AP.TxnBody[i].PreviousTxnHash)
			}

			pData2, errorAsync2 := object.GetLastTransactionbyIdentifierNotSplitParent(AP.TxnBody[i].FromIdentifier2).Then(func(data interface{}) interface{} {
				return data
			}).Await()

			if errorAsync2 != nil || pData2 == nil {
				log.Error("Error while GetLastTransactionbyIdentifier @SubmitMerge " + errorAsync.Error())
				///ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
				//DUE TO THE CHILD HAVING A NEW IDENTIFIER
				AP.TxnBody[i].PreviousTxnHash2 = ""
			} else {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER
				result2 := pData2.(model.TransactionCollectionBody)
				AP.TxnBody[i].PreviousTxnHash2 = result2.TxnHash
				log.Debug(AP.TxnBody[i].PreviousTxnHash)
			}
		} else {
			pData3, errorAsync3 := object.GetLastTransactionbyIdentifierNotSplitParent(AP.TxnBody[i].FromIdentifier2).Then(func(data interface{}) interface{} {
				return data
			}).Await()

			if errorAsync3 != nil || pData3 == nil {
				log.Error("Error while GetLastTransactionbyIdentifier @SubmitMerge " + errorAsync3.Error())
				///ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
				//DUE TO THE CHILD HAVING A NEW IDENTIFIER
				AP.TxnBody[i].PreviousTxnHash2 = ""
			} else {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER
				result := pData3.(model.TransactionCollectionBody)
				AP.TxnBody[i].PreviousTxnHash2 = result.TxnHash
				log.Debug(AP.TxnBody[i].PreviousTxnHash2)
			}
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

		if AP.TxnBody[i].TxnType == "7" {

			id.MapValue = AP.TxnBody[i].Identifier
			id.Identifier = AP.TxnBody[i].MapIdentifier

			err3 := object.InsertIdentifier(id)
			if err3 != nil {
				fmt.Println("identifier map failed" + err3.Error())
			}
		}

	}
	go func() {

		var PreviousTxn string

		for i, TxnBody := range AP.TxnBody {
			var TypeTXNBuilder txnbuild.ManageData
			var PreviousTXNBuilder txnbuild.ManageData
			var MergeIDBuilder txnbuild.ManageData

			////GET THE PREVIOUS TRANSACTION FOR THE IDENTIFIER
			//INCASE OF FIRST MERGE BLOCK THE PREVIOUS IS TAKEN FROM IDENTIFIER
			//&
			//INCASE OF GREATER THAN ONE THE PREVIOUS TXN IS THE PREVIOUS MERGE
			if i == 0 {
				// TypeTXNBuilder = build.SetData("Type", []byte("G8"))
				// PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(AP.TxnBody[i].PreviousTxnHash))
				// MergeIDBuilder = build.SetData("MergeID", []byte(AP.TxnBody[i].PreviousTxnHash2))
				TypeTXNBuilder = txnbuild.ManageData{
					Name: "Type",
					Value: []byte("G8"),
				}
				PreviousTXNBuilder = txnbuild.ManageData{
					Name: "PreviousTXN",
					Value: []byte(AP.TxnBody[i].PreviousTxnHash),
				}
				MergeIDBuilder = txnbuild.ManageData{
					Name: "MergeID",
					Value: []byte(AP.TxnBody[i].PreviousTxnHash2),
				}
				AP.TxnBody[i].MergeID = AP.TxnBody[i].PreviousTxnHash2
			} else {
				// TypeTXNBuilder = build.SetData("Type", []byte("G7"))
				// PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(PreviousTxn))
				// MergeIDBuilder = build.SetData("MergeID", []byte(AP.TxnBody[i].PreviousTxnHash2))
				TypeTXNBuilder = txnbuild.ManageData{
					Name: "Type",
					Value: []byte("G7"),
				}
				PreviousTXNBuilder = txnbuild.ManageData{
					Name: "PreviousTXN",
					Value: []byte(PreviousTxn),
				}
				MergeIDBuilder = txnbuild.ManageData{
					Name: "MergeID",
					Value: []byte(AP.TxnBody[i].PreviousTxnHash2),
				}
				AP.TxnBody[i].MergeID = AP.TxnBody[i].PreviousTxnHash2
			}

			if i == 0 {
				pData, errorAsync := object.GetLastTransactionbyIdentifierNotSplitParent(TxnBody.FromIdentifier2).Then(func(data interface{}) interface{} {
					return data
				}).Await()

				if errorAsync != nil || pData == nil {
					log.Error("Error while GetLastTransactionbyIdentifier @SubmitMerge " + errorAsync.Error())
					///ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
					//DUE TO THE CHILD HAVING A NEW IDENTIFIER
					AP.TxnBody[i].MergeID = ""
				} else {
					///ASSIGN PREVIOUS MANAGE DATA BUILDER
					result := pData.(model.TransactionCollectionBody)
					// MergeID = result.TxnHash
					AP.TxnBody[i].MergeID = result.TxnHash
					fmt.Println(AP.TxnBody[i].MergeID)
				}
			}

			// pubaccountRequest := horizonclient.AccountRequest{AccountID: publicKey}
			// pubaccount, err := netClient.AccountDetail(pubaccountRequest)

			kp,_ := keypair.Parse(publicKey)
			client := commons.GetHorizonNetwork()
			ar := horizonclient.AccountRequest{AccountID: kp.Address()}
			pubaccount, err := client.AccountDetail(ar)
			
			if err != nil{
				log.Fatal(err)
			}

			//BUILD THE GATEWAY XDR
			// tx, err := build.Transaction(
			// 	commons.GetHorizonNetwork(),
			// 	build.SourceAccount{publicKey},
			// 	build.AutoSequence{commons.GetHorizonClient()},
			// 	TypeTXNBuilder,
			// 	PreviousTXNBuilder,
			// 	build.SetData("CurrentTXN", []byte(UserMergeTxnHashes[i])),
			// 	MergeIDBuilder,
			// )
			CurrentTXN := txnbuild.ManageData{
				Name: "CurrentTXN",
				Value: []byte(UserMergeTxnHashes[i]),
			}
			tx, err := txnbuild.NewTransaction(
				txnbuild.TransactionParams{
					SourceAccount: &pubaccount,
					IncrementSequenceNum: true,
					Operations: []txnbuild.Operation{&TypeTXNBuilder, &PreviousTXNBuilder, &MergeIDBuilder, &CurrentTXN},
					BaseFee: constants.MinBaseFee,
					Preconditions: txnbuild.Preconditions{TimeBounds: constants.TransactionTimeOut},
				},
			)

			//SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
			GatewayTXE, err := tx.Sign(commons.GetStellarNetwork(), tracifiedAccount)
			if err != nil {
				log.Error("Error while build Transaction @SubmitMerge " + err.Error())
				AP.TxnBody[i].TxnHash = UserMergeTxnHashes[i]
				AP.TxnBody[i].Status = "Pending"

				///INSERT INTO TRANSACTION COLLECTION
				err2 := object.InsertTransaction(AP.TxnBody[i])
				if err2 != nil {
					log.Error("Error while InsertTransaction @SubmitMerge " + err2.Error())
				}
			}
			//CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
			txeB64, err := GatewayTXE.Base64()
			if err != nil {
				log.Error("Error while convert GatewayTXE to base64 @SubmitMerge " + err.Error())
				AP.TxnBody[i].TxnHash = UserMergeTxnHashes[i]
				AP.TxnBody[i].Status = "Pending"

				///INSERT INTO TRANSACTION COLLECTION
				err2 := object.InsertTransaction(AP.TxnBody[i])
				if err2 != nil {
					log.Error("Error while InsertTransaction @SubmitMerge " + err2.Error())
				}
			}

			//SUBMIT THE GATEWAY'S SIGNED XDR
			display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
			response1 := display1.SubmitXDR("G" + AP.TxnBody[i].TxnType)

			if response1.Error.Code == 400 {
				log.Error("Error got 400 from ConcreteSubmitXDR @SubmitMerge ")
				AP.TxnBody[i].TxnHash = UserMergeTxnHashes[i]
				AP.TxnBody[i].Status = "Pending"

				///INSERT INTO TRANSACTION COLLECTION
				err2 := object.InsertTransaction(AP.TxnBody[i])
				if err2 != nil {
					log.Error("Error while InsertTransaction @SubmitMerge " + err2.Error())
				}
			} else {
				//UPDATE THE TRANSACTION COLLECTION WITH TXN HASH
				AP.TxnBody[i].TxnHash = response1.TXNID
				PreviousTxn = response1.TXNID

				///INSERT INTO TRANSACTION COLLECTION
				err = object.InsertTransaction(AP.TxnBody[i])
				if err != nil {
					log.Error("Error while InsertTransaction @SubmitMerge " + err.Error())
				} else if i == 0 {
					var PreviousProfile string
					pData, errorAsync := object.GetProfilebyIdentifier(AP.TxnBody[i].FromIdentifier1).Then(func(data interface{}) interface{} {
						return data
					}).Await()

					if errorAsync != nil || pData == nil {
						log.Error("Error while GetProfilebyIdentifier @SubmitMerge" + errorAsync.Error())
						PreviousProfile = ""
					} else {
						result := pData.(model.ProfileCollectionBody)
						PreviousProfile = result.ProfileTxn
					}

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
						log.Error("Error while InsertProfile @SubmitMerge " + err3.Error())
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
