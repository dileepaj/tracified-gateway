package builder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
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
	var Done []bool
	Done = append(Done, true)

	object := dao.Connection{}

	var UserSplitTxnHashes []string
	// var ParentIdentifier string
	// var ParentTxn string
	var PreviousTxn string

	///HARDCODED CREDENTIALS
	publicKey := constants.PublicKey
	secretKey := constants.SecretKey
	// var result model.SubmitXDRResponse

	time.Sleep(4 * time.Second)

	for i, TxnBody := range AP.TxnBody {

		var txe xdr.Transaction

		//decode the XDR
		err := xdr.SafeUnmarshalBase64(TxnBody.XDR, &txe)
		if err != nil {
		}

		//GET THE TYPE AND IDENTIFIER FROM THE XDR
		AP.TxnBody[i].PublicKey = txe.SourceAccount.Address()

		AP.TxnBody[i].TxnType = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].Identifier = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[1].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].ItemCode = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[3].Body.ManageDataOp.DataValue), "&")
		if i == 0 {
			AP.TxnBody[i].ToIdentifier = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[2].Body.ManageDataOp.DataValue), "&")
		} else {
			AP.TxnBody[i].FromIdentifier1 = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[2].Body.ManageDataOp.DataValue), "&")
			AP.TxnBody[i].ItemAmount = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[4].Body.ManageDataOp.DataValue), "&")
		}

		//FOR THE SPLIT PARENT RETRIEVE THE PREVIOUS TXN FROM GATEWAY DB
		if i == 0 {
			// ParentIdentifier = Identifier
			p := object.GetLastTransactionbyIdentifier(AP.TxnBody[i].Identifier)
			p.Then(func(data interface{}) interface{} {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER
				result := data.(model.TransactionCollectionBody)
				PreviousTxn = result.TxnHash
				AP.TxnBody[i].PreviousTxnHash = result.TxnHash
				return nil
			}).Catch(func(error error) error {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
				//DUE TO THE CHILD HAVING A NEW IDENTIFIER
				
				AP.TxnBody[i].PreviousTxnHash = ""
				return error
			})
			p.Await()
		}

		//SUBMIT THE FIRST XDR SIGNED BY THE USER
		display := stellarExecuter.ConcreteSubmitXDR{XDR: TxnBody.XDR}
		result := display.SubmitXDR(false,AP.TxnBody[i].TxnType)
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
			fmt.Println((i + 1), " Submitted")
		}
	}
	go func() {

		var SplitParentProfile string
		var PreviousSplitProfile string
		for i, TxnBody := range AP.TxnBody {
			var PreviousTXNBuilder build.ManageDataBuilder

			PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(PreviousTxn))
			AP.TxnBody[i].PreviousTxnHash = PreviousTxn

			//ASSIGN THE PREVIOUS PROFILE ID USING THE PARENT FOR THE CHILDREN AND A DB CALL FOR PARENT
			if i == 0 {
				PreviousSplitProfile = ""
				SplitParentProfile = AP.TxnBody[i].ProfileID
			} else {
				PreviousSplitProfile = SplitParentProfile
			}
			//BUILD THE GATEWAY XDR
			tx, err := build.Transaction(
				build.PublicNetwork,
				build.SourceAccount{publicKey},
				build.AutoSequence{horizon.DefaultPublicNetClient},
				build.SetData("Type", []byte("G"+TxnBody.TxnType)),
				PreviousTXNBuilder,
				build.SetData("CurrentTXN", []byte(UserSplitTxnHashes[i])),
				build.SetData("Identifier", []byte(AP.TxnBody[i].Identifier)),
				build.SetData("ProfileID", []byte(AP.TxnBody[i].ProfileID)),
				build.SetData("PreviousProfile", []byte(PreviousSplitProfile)),
			)

			//SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
			GatewayTXE, err := tx.Sign(secretKey)
			if err != nil {
				AP.TxnBody[i].TxnHash = UserSplitTxnHashes[i]
				AP.TxnBody[i].Status = "Pending"

				///INSERT INTO TRANSACTION COLLECTION
				err2 := object.InsertTransaction(AP.TxnBody[i])
				if err2 != nil {
				}
			}
			//CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
			txeB64, err := GatewayTXE.Base64()
			if err != nil {
				AP.TxnBody[i].TxnHash = UserSplitTxnHashes[i]
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
				AP.TxnBody[i].TxnHash = UserSplitTxnHashes[i]
				AP.TxnBody[i].Status = "Pending"

				///INSERT INTO TRANSACTION COLLECTION
				err2 := object.InsertTransaction(AP.TxnBody[i])
				if err2 != nil {
				}
			} else {
				//UPDATE THE TRANSACTION COLLECTION WITH TXN HASH
				AP.TxnBody[i].TxnHash = response1.TXNID
				if i == 0 {
					PreviousTxn = response1.TXNID
				}

				///INSERT INTO TRANSACTION COLLECTION
				err1 := object.InsertTransaction(AP.TxnBody[i])
				if err1 != nil {
				} else if i > 0 {

					var PreviousProfile string
					p := object.GetProfilebyIdentifier(AP.TxnBody[i].FromIdentifier1)
					p.Then(func(data interface{}) interface{} {

						result := data.(model.ProfileCollectionBody)
						PreviousProfile = result.ProfileTxn
						return nil
					}).Catch(func(error error) error {
						PreviousProfile = ""
						return error
					})
					p.Await()

					Profile := model.ProfileCollectionBody{
						ProfileTxn:         response1.TXNID,
						ProfileID:          AP.TxnBody[i].ProfileID,
						Identifier:         AP.TxnBody[i].Identifier,
						PreviousProfileTxn: PreviousProfile,
						TriggerTxn:         UserSplitTxnHashes[i],
						TxnType:            AP.TxnBody[i].TxnType,
					}
					err3 := object.InsertProfile(Profile)
					if err3 != nil {

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
