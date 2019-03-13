package builder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/xdr"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
)


type AbstractCertificateSubmiter struct {
	TxnBody []model.CertificateCollectionBody
}


//I AM THE GENESIS BUILDER
func (AP *AbstractCertificateSubmiter) SubmitInsertCertificate(w http.ResponseWriter, r *http.Request) {
	var Done []bool
	Done = append(Done, true)

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
		errx := xdr.SafeUnmarshalBase64(TxnBody.XDR, &txe)
		if errx != nil {
		}

		//GET THE TYPE AND IDENTIFIER FROM THE XDR
		AP.TxnBody[i].PublicKey = txe.SourceAccount.Address()
		AP.TxnBody[i].TxnType = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].PreviousCertificate = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[1].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].CertificateType = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[2].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].Data= strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[3].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].ValidityPeriod= strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[4].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].Asset= strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[5].Body.ManageDataOp.DataValue), "&")

		//SUBMIT THE FIRST XDR SIGNED BY THE USER
		display := stellarExecuter.ConcreteSubmitXDR{XDR: AP.TxnBody[i].XDR}
		result := display.SubmitXDR()
		UserTxnHashes = append(UserTxnHashes, result.TXNID)

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
			PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(AP.TxnBody[i].TxnType))

			//BUILD THE GATEWAY XDR
			tx, err := build.Transaction(
				build.TestNetwork,
				build.SourceAccount{publicKey},
				build.AutoSequence{horizon.DefaultTestNetClient},
				build.SetData("Type", []byte("G"+TxnBody.TxnType)),
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
			response1 := display1.SubmitXDR()

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

				} else {
					var PreviousProfile string
					p := object.GetProfilebyIdentifier(AP.TxnBody[i].Identifier)
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
						TriggerTxn:         UserTxnHashes[i],
						TxnType:            AP.TxnBody[i].TxnType,
					}
					err3 := object.InsertProfile(Profile)
					if err3 != nil {

					}

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
					fmt.Println(err.Error())
				}
				return nil
			}).Catch(func(error error) error {
				return error
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
