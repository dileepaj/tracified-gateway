package builder

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	// "strconv"
	"strings"

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

/*SubmitRenewCertificate - WORKING MODEL
@author - Azeem Ashraf
@desc - Builds the TXN Type C2 for the gateway where it receives the user XDR
and decodes it's contents and submit's to stellar and further maps the received TXN
to Gateway Signed TXN's to maintain the profile, also records the activity in the gateway datastore
@note - Should implement a validation layer to validate the contents of the XDR per builder before submission.
@params - ResponseWriter,Request
*/
func (AP *AbstractCertificateSubmiter) SubmitRenewCertificate(w http.ResponseWriter, r *http.Request) {
	log.Debug("======================= SubmitRenewCertificate ========================")
	var Done []bool
	Done = append(Done, true)
	//netClient := commons.GetHorizonClient()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	object := dao.Connection{}
	var UserTxnHashes []string

	valid:=false
	
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
		errx := xdr.SafeUnmarshalBase64(TxnBody.XDR, &txe)
		if errx != nil {
			log.Error("Error while SafeUnmarshalBase64 @SubmitRenewCertificate "+errx.Error())
		}

		//GET THE TYPE AND IDENTIFIER FROM THE XDR
		AP.TxnBody[i].PublicKey = txe.SourceAccount.Address()
		AP.TxnBody[i].TxnType = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].CertificateType = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[1].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].CertificateID = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[2].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].ValidityPeriod = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[3].Body.ManageDataOp.DataValue), "&")

		p := object.GetLastCertificatebyCertificateID(AP.TxnBody[i].CertificateID)
		p.Then(func(data interface{}) interface{} {
			result := data.(model.CertificateCollectionBody)
			fmt.Println(result.PublicKey+" "+AP.TxnBody[i].PublicKey)
			if result.PublicKey == AP.TxnBody[i].PublicKey {
				valid = true
			} else {
				valid = false
				Done = append(Done, false)
				w.WriteHeader(400)
				response := apiModel.SubmitXDRSuccess{
					Status: "Certificate Blockchain Transaction Failed! Wrong Issuer Key",
				}
				json.NewEncoder(w).Encode(response)
			}
			return nil
		}).Catch(func(error error) error {
			log.Error("Error while GetLastCertificatebyCertificateID @SubmitRenewCertificate "+error.Error())
			valid = false
			Done = append(Done, false)
			w.WriteHeader(400)
			response := apiModel.SubmitXDRSuccess{
				Status: "Certificate Blockchain Transaction Failed! No Certificate Found ",
			}
			json.NewEncoder(w).Encode(response)
			return error
		})
		p.Await()

		if valid {
			//SUBMIT THE FIRST XDR SIGNED BY THE USER
			display := stellarExecuter.ConcreteSubmitXDR{XDR: AP.TxnBody[i].XDR}
			result := display.SubmitXDR(AP.TxnBody[i].TxnType)
			UserTxnHashes = append(UserTxnHashes, result.TXNID)

			if result.Error.Code == 400 {
				log.Error("Error got 400 for ConcreteSubmitXDR @SubmitRenewCertificate ")
				Done = append(Done, false)
				w.WriteHeader(result.Error.Code)
				response := apiModel.SubmitXDRSuccess{
					Status: "Certificate Blockchain Transaction Failed!",
				}
				json.NewEncoder(w).Encode(response)
				return
			}
		} else {
			return
		}
	}

	go func() {

		for i, TxnBody := range AP.TxnBody {

			var PreviousTXNBuilder txnbuild.ManageData

			//GET THE PREVIOUS CERTIFICATE FOR THE PUBLIC KEY
			p := object.GetLastCertificatebyPublicKey(AP.TxnBody[i].PublicKey)
			p.Then(func(data interface{}) interface{} {
				result := data.(model.CertificateCollectionBody)
				//PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(result.CertificateID))
				PreviousTXNBuilder = txnbuild.ManageData{
					Name:"PreviousTXN",
					Value: []byte(result.CertificateID),
				}
				AP.TxnBody[i].PreviousCertificate = result.CertificateID
				return nil
			}).Catch(func(error error) error {
				log.Error("Error while GetLastCertificatebyPublicKey @SubmitRenewCertificate "+error.Error())
				//PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(""))
				PreviousTXNBuilder = txnbuild.ManageData{
					Name: "PreviousTXN",
					Value: []byte(""),
				}
				return error
			})
			p.Await()

			// pubaccountRequest := horizonclient.AccountRequest{AccountID: publicKey}
			// pubaccount, err := netClient.AccountDetail(pubaccountRequest)
			kp,_ := keypair.Parse(publicKey)
			client := commons.GetHorizonClient()
			ar := horizonclient.AccountRequest{AccountID: kp.Address()}
			pubaccount, err := client.AccountDetail(ar)
			if err != nil {
				log.Fatal(err)
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

			TypeTxnBuilder := txnbuild.ManageData{
				Name: "Type",
				Value: []byte("G"+TxnBody.TxnType),
			}

			CurrentTXNBuilder := txnbuild.ManageData{
				Name: "CurrentTXN",
				Value: []byte(UserTxnHashes[i]),
			}

			tx, err := txnbuild.NewTransaction(
					txnbuild.TransactionParams{
					SourceAccount: &pubaccount,
					IncrementSequenceNum: true,
					Operations: []txnbuild.Operation{&PreviousTXNBuilder, &TypeTxnBuilder, &CurrentTXNBuilder},
					BaseFee: constants.MinBaseFee,
					Preconditions:txnbuild.Preconditions{TimeBounds: constants.TransactionTimeOut},
				},
			)

			//SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
			GatewayTXE, err := tx.Sign(commons.GetStellarNetwork(), tracifiedAccount)
			if err != nil {
				log.Error("Error while sign"+err.Error())
				AP.TxnBody[i].CertificateID = UserTxnHashes[i]
				AP.TxnBody[i].Status = "Pending"

				///INSERT INTO TRANSACTION COLLECTION
				err2 := object.InsertCertificate(AP.TxnBody[i])
				if err2 != nil {
					log.Error("Error while InsertCertificate @SubmitRenewCertificate "+err2.Error())
				}
			}

			//CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
			txeB64, err := GatewayTXE.Base64()
			if err != nil {
				AP.TxnBody[i].CertificateID = UserTxnHashes[i]
				AP.TxnBody[i].Status = "Pending"

				///INSERT INTO TRANSACTION COLLECTION
				err2 := object.InsertCertificate(AP.TxnBody[i])
				if err2 != nil {
					log.Error("Error while InsertCertificate @SubmitRenewCertificate "+err2.Error())
				}
			}
			//SUBMIT THE GATEWAY'S SIGNED XDR
			display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
			response1 := display1.SubmitXDR("G"+AP.TxnBody[i].TxnType)

			if response1.Error.Code == 400 {
				AP.TxnBody[i].CertificateID = UserTxnHashes[i]
				AP.TxnBody[i].Status = "Pending"

				///INSERT INTO TRANSACTION COLLECTION
				err2 := object.InsertCertificate(AP.TxnBody[i])
				if err2 != nil {
					log.Error("Error while InsertCertificate @SubmitRenewCertificate "+err2.Error())
				}
			} else {
				//UPDATE THE TRANSACTION COLLECTION WITH TXN HASH
				AP.TxnBody[i].CertificateID = response1.TXNID
				AP.TxnBody[i].Status = "Done"

				///INSERT INTO TRANSACTION COLLECTION
				err2 := object.InsertCertificate(AP.TxnBody[i])
				if err2 != nil {
					log.Error("Error while InsertCertificate @SubmitRenewCertificate "+err2.Error())
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
