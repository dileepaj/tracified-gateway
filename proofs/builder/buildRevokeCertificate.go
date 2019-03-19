package builder

import (
	"encoding/json"
	"fmt"
	"net/http"

	// "strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/xdr"
)

//I AM THE GENESIS BUILDER
func (AP *AbstractCertificateSubmiter) SubmitRevokeCertificate(w http.ResponseWriter, r *http.Request) {
	var Done []bool
	Done = append(Done, true)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	object := dao.Connection{}
	var UserTxnHashes []string

	valid:=false
	
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
		AP.TxnBody[i].CertificateType = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[1].Body.ManageDataOp.DataValue), "&")
		AP.TxnBody[i].CertificateID = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[2].Body.ManageDataOp.DataValue), "&")
		// AP.TxnBody[i].ValidityPeriod = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[3].Body.ManageDataOp.DataValue), "&")

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
			result := display.SubmitXDR()
			UserTxnHashes = append(UserTxnHashes, result.TXNID)

			if result.Error.Code == 400 {
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

			var PreviousTXNBuilder build.ManageDataBuilder

			//GET THE PREVIOUS CERTIFICATE FOR THE PUBLIC KEY
			p := object.GetLastCertificatebyPublicKey(AP.TxnBody[i].PublicKey)
			p.Then(func(data interface{}) interface{} {
				result := data.(model.CertificateCollectionBody)
				PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(result.CertificateID))
				AP.TxnBody[i].PreviousCertificate = result.CertificateID
				return nil
			}).Catch(func(error error) error {
				PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(""))
				return error
			})
			p.Await()

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
				AP.TxnBody[i].CertificateID = UserTxnHashes[i]
				AP.TxnBody[i].Status = "Pending"

				///INSERT INTO TRANSACTION COLLECTION
				err2 := object.InsertCertificate(AP.TxnBody[i])
				if err2 != nil {

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

				}
			}
			//SUBMIT THE GATEWAY'S SIGNED XDR
			display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
			response1 := display1.SubmitXDR()

			if response1.Error.Code == 400 {
				AP.TxnBody[i].CertificateID = UserTxnHashes[i]
				AP.TxnBody[i].Status = "Pending"

				///INSERT INTO TRANSACTION COLLECTION
				err2 := object.InsertCertificate(AP.TxnBody[i])
				if err2 != nil {

				}
			} else {
				//UPDATE THE TRANSACTION COLLECTION WITH TXN HASH
				AP.TxnBody[i].CertificateID = response1.TXNID
				AP.TxnBody[i].Status = "Done"

				///INSERT INTO TRANSACTION COLLECTION
				err2 := object.InsertCertificate(AP.TxnBody[i])
				if err2 != nil {

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
