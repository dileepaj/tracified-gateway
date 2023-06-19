package builder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

/*SubmitSplit - WORKING MODEL
@author - Azeem Ashraf
@desc - Builds the TXN Type 5 for the gateway where it receives the user XDR
and decodes it's contents and submit's to stellar and further maps the received TXN
to Gateway Signed TXN's to maintain the profile, also records the activity in the gateway datastore.
@note - Should implement a validation layer to validate the contents of the XDR per builder before submission.
@params - ResponseWriter,Request
*/
func (AP *AbstractXDRSubmiter) SubmitSplit(w http.ResponseWriter, r *http.Request) {
	log.Debug("--------------------------------- SubmitSplit -------------------------------------")

	var Done []bool
	Done = append(Done, true)
	var id apiModel.IdentifierModel
	object := dao.Connection{}

	var UserSplitTxnHashes []string
	// var ParentIdentifier string
	// var ParentTxn string
	var PreviousTxn string

	///HARDCODED CREDENTIALS
	publicKey := constants.PublicKey
	secretKey := constants.SecretKey
	tracifiedAccount, err := keypair.ParseFull(secretKey)
	if err != nil {
		log.Error(err)
	}
	// var result model.SubmitXDRResponse
	client := commons.GetHorizonClient()
	ar := horizonclient.AccountRequest{AccountID: publicKey}
	account, err := client.AccountDetail(ar)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	time.Sleep(4 * time.Second)

	for i, TxnBody := range AP.TxnBody {

		var txe xdr.Transaction

		// decode the XDR
		err := xdr.SafeUnmarshalBase64(TxnBody.XDR, &txe)
		if err != nil {
			log.Error("Error @SafeUnmarshalBase64 @SubmitSplit " + err.Error())
		}

		// GET THE TYPE AND IDENTIFIER FROM THE XDR
		AP.TxnBody[i].PublicKey = txe.SourceAccount.Address()
		AP.TxnBody[i].SequenceNo = int64(txe.SeqNum)
		AP.TxnBody[i].TxnType = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")

		if AP.TxnBody[i].TxnType == "5" {
			AP.TxnBody[i].Identifier = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[1].Body.ManageDataOp.DataValue), "&")
			AP.TxnBody[i].ToIdentifier = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[2].Body.ManageDataOp.DataValue), "&")
			AP.TxnBody[i].ItemCode = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[3].Body.ManageDataOp.DataValue), "&")
			AP.TxnBody[i].ProductName = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[4].Body.ManageDataOp.DataValue), "&")
			AP.TxnBody[i].AppAccount = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[5].Body.ManageDataOp.DataValue), "&")
		} else {
			AP.TxnBody[i].Identifier = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[1].Body.ManageDataOp.DataValue), "&")
			AP.TxnBody[i].FromIdentifier1 = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[2].Body.ManageDataOp.DataValue), "&")
			AP.TxnBody[i].ItemCode = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[3].Body.ManageDataOp.DataValue), "&")
			AP.TxnBody[i].ItemAmount = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[4].Body.ManageDataOp.DataValue), "&")
			AP.TxnBody[i].ProductName = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[5].Body.ManageDataOp.DataValue), "&")
			AP.TxnBody[i].AppAccount = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[6].Body.ManageDataOp.DataValue), "&")
		}

		// FOR THE SPLIT PARENT RETRIEVE THE PREVIOUS TXN FROM GATEWAY DB
		if AP.TxnBody[i].TxnType == "5" {
			// ParentIdentifier = Identifier
			pData, errAsnc := object.GetLastTransactionbyIdentifierAndTenantId(AP.TxnBody[i].Identifier, AP.TxnBody[i].TenantID, AP.TxnBody[i].ProductID).Then(func(data interface{}) interface{} {
				return data
			}).Await()

			if pData == nil || errAsnc != nil {
				log.Error("Error @GetLastTransactionbyIdentifier @SubmitSplit ")
				///ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
				//DUE TO THE CHILD HAVING A NEW IDENTIFIER
				PreviousTxn = ""
				AP.TxnBody[i].PreviousTxnHash = ""
			} else {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER
				result := pData.(model.TransactionCollectionBody)
				PreviousTxn = result.TxnHash
				AP.TxnBody[i].PreviousTxnHash = result.TxnHash
			}
		}

		// SUBMIT THE FIRST XDR SIGNED BY THE USER
		display := stellarExecuter.ConcreteSubmitXDR{XDR: TxnBody.XDR}
		result := display.SubmitXDR(AP.TxnBody[i].TxnType)
		UserSplitTxnHashes = append(UserSplitTxnHashes, result.TXNID)

		if result.Error.Code == 400 {
			Done = append(Done, false)
			w.WriteHeader(result.Error.Code)
			response := apiModel.SubmitXDRSuccess{
				Status: "Index[" + strconv.Itoa(i) + "] TXN: Blockchain Transaction Failed!",
			}
			json.NewEncoder(w).Encode(response)
			return
		} else {
			log.Info((i + 1), " Submitted")
		}

		if AP.TxnBody[i].TxnType == "6" {
    
			id.MapValue = AP.TxnBody[i].Identifier
			id.Identifier = AP.TxnBody[i].MapIdentifier
      
			err3 := object.InsertIdentifier(id)
			if err3 != nil {
				fmt.Println("identifier map failed" + err3.Error())
			}
		}

	}
	go func() {
		var SplitParentProfile string
		var PreviousSplitProfile string
		for i, TxnBody := range AP.TxnBody {
			/*
			When constructing a backlink transaction(put from gateway) for a split, it is important to exclude the split-parent transaction as its previous transaction.
			Instead, you should obtain the most recent transaction that is specific to the identifier and disregard the split-parent transaction.
			*/
			if TxnBody.TxnType == "6" {
			backlinkData, errAsnc := object.GetLastTransactionbyIdentifierNotSplitParent(AP.TxnBody[i].FromIdentifier1, AP.TxnBody[i].TenantID).Then(func(data interface{}) interface{} {
				return data
			}).Await()
			if backlinkData == nil || errAsnc != nil {
				log.Info("Can not find transaction form database ","build Split")
			} else {
				result := backlinkData.(model.TransactionCollectionBody)
				PreviousTxn = result.TxnHash
				AP.TxnBody[i].PreviousTxnHash = result.TxnHash
			}
			}
			previousTXNBuilder := txnbuild.ManageData{Name: "PreviousTXN", Value: []byte(PreviousTxn)}
			typeTXNBuilder := txnbuild.ManageData{Name: "Type", Value: []byte("G" + TxnBody.TxnType)}
			currentTXNBuilder := txnbuild.ManageData{Name: "CurrentTXN", Value: []byte(UserSplitTxnHashes[i])}
			identifierTXNBuilder := txnbuild.ManageData{Name: "Identifier", Value: []byte(AP.TxnBody[i].Identifier)}
			profileIDTXNBuilder := txnbuild.ManageData{Name: "ProfileID", Value: []byte(AP.TxnBody[i].ProfileID)}
			PreviousProfileTXNBuilder := txnbuild.ManageData{Name: "PreviousProfile", Value: []byte(PreviousSplitProfile)}

			AP.TxnBody[i].PreviousTxnHash = PreviousTxn

			// ASSIGN THE PREVIOUS PROFILE ID USING THE PARENT FOR THE CHILDREN AND A DB CALL FOR PARENT
			if AP.TxnBody[i].TxnType == "5" {
				PreviousSplitProfile = ""
				SplitParentProfile = AP.TxnBody[i].ProfileID
			} else {
				PreviousSplitProfile = SplitParentProfile
			}

			// BUILD THE GATEWAY XDR
			tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
				SourceAccount:        &account,
				IncrementSequenceNum: true,
				Operations:           []txnbuild.Operation{&previousTXNBuilder, &typeTXNBuilder, &currentTXNBuilder, &identifierTXNBuilder, &profileIDTXNBuilder, &PreviousProfileTXNBuilder},
				BaseFee:              constants.MinBaseFee,
				Memo:                 nil,
				Preconditions:        txnbuild.Preconditions{TimeBounds: constants.TransactionTimeOut},
			})

			// SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
			GatewayTXE, err := tx.Sign(commons.GetStellarNetwork(), tracifiedAccount)
			if err != nil {
				log.Error("Error @tx.Sign @SubmitSplit " + err.Error())
				AP.TxnBody[i].TxnHash = UserSplitTxnHashes[i]
				AP.TxnBody[i].Status = "Pending"

				///INSERT INTO TRANSACTION COLLECTION
				err2 := object.InsertTransaction(AP.TxnBody[i])
				if err2 != nil {
					log.Error("Error @InsertTransaction @SubmitSplit " + err2.Error())
				}
			}
			// CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
			txeB64, err := GatewayTXE.Base64()
			if err != nil {
				log.Error("Error @GatewayTXE.Base64 @SubmitSplit " + err.Error())
				AP.TxnBody[i].TxnHash = UserSplitTxnHashes[i]
				AP.TxnBody[i].Status = "Pending"

				///INSERT INTO TRANSACTION COLLECTION
				err2 := object.InsertTransaction(AP.TxnBody[i])
				if err2 != nil {
					log.Error("Error @InsertTransaction @SubmitSplit " + err2.Error())
				}
			}

			// SUBMIT THE GATEWAY'S SIGNED XDR
			display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
			response1 := display1.SubmitXDR("G" + AP.TxnBody[i].TxnType)

			if response1.Error.Code == 400 {
				AP.TxnBody[i].TxnHash = UserSplitTxnHashes[i]
				AP.TxnBody[i].Status = "Pending"

				///INSERT INTO TRANSACTION COLLECTION
				err2 := object.InsertTransaction(AP.TxnBody[i])
				if err2 != nil {
					log.Error("Error @InsertTransaction @SubmitSplit " + err2.Error())
				}
			} else {
				// UPDATE THE TRANSACTION COLLECTION WITH TXN HASH
				AP.TxnBody[i].TxnHash = response1.TXNID
				if AP.TxnBody[i].TxnType == "5" {
					PreviousTxn = response1.TXNID
				}

				///INSERT INTO TRANSACTION COLLECTION
				err1 := object.InsertTransaction(AP.TxnBody[i])
				if err1 != nil {
					log.Error("Error @InsertTransaction @SubmitSplit " + err1.Error())
				} else if i > 0 {
					var PreviousProfile string
					pData1, errorAsync1 := object.GetProfilebyIdentifier(AP.TxnBody[i].FromIdentifier1).Then(func(data interface{}) interface{} {
						return data
					}).Await()
					if pData1 == nil || errorAsync1 != nil {
						log.Error("Error @GetProfilebyIdentifier @SubmitSplit ")
						PreviousProfile = ""
					} else {
						result := pData1.(model.ProfileCollectionBody)
						PreviousProfile = result.ProfileTxn
					}
					Profile := model.ProfileCollectionBody{
						ProfileTxn:         response1.TXNID,
						ProfileID:          AP.TxnBody[i].ProfileID,
						Identifier:         AP.TxnBody[i].Identifier,
						PreviousProfileTxn: PreviousProfile,
						TriggerTxn:         UserSplitTxnHashes[i],
						TxnType:            AP.TxnBody[i].TxnType,
					}
					err2 := object.InsertProfile(Profile)
					if err2 != nil {
						log.Error("Error @InsertProfile @SubmitSplit " + err2.Error())
					}
				}
			}
		}
	}()
	// }
	if checkBoolArray(Done) {
		w.WriteHeader(http.StatusOK)
		result := apiModel.SubmitXDRSuccess{
			Status: "Success",
		}
		json.NewEncoder(w).Encode(result)
	}
	return
}
