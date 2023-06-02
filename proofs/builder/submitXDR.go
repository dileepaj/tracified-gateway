package builder

import (
	// "encoding/json"

	log "github.com/sirupsen/logrus"

	// "net/http"
	// "strconv"

	// "github.com/dileepaj/tracified-gateway/api/apiModel"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
)

/*AbstractXDRSubmiter - WORKING MODEL
@author - Azeem Ashraf
@desc - Abstract Struct that hold's the TransactionModel
*/
type AbstractXDRSubmiter struct {
	TxnBody []model.TransactionCollectionBody
}

/*XDRSubmitter - Deprecated
@author - Azeem Ashraf
*/
func XDRSubmitter(TDP []model.TransactionCollectionBody) (bool, model.SubmitXDRResponse) {
	log.Debug("---------------------------------- XDRSubmitter ----------------------------------")
	var status []bool
	object := dao.Connection{}
	var ret model.SubmitXDRResponse
	var UserTxnHashes []string
	status = append(status, true)

	///HARDCODED CREDENTIALS
	publicKey := constants.PublicKey
	secretKey := constants.SecretKey
	tracifiedAccount, err := keypair.ParseFull(secretKey)
	if err != nil {
		log.Error(err)
	}
	// Get information about the account we just created
	kp,_ := keypair.Parse(publicKey)
	client := commons.GetHorizonClient()
	ar := horizonclient.AccountRequest{AccountID: kp.Address()}
	account, err := client.AccountDetail(ar)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(TDP); i++ {
		TDP[i].Status = "Pending"
		var txe xdr.Transaction
		err := xdr.SafeUnmarshalBase64(TDP[i].XDR, &txe)
		if err != nil {
			log.Error("Error @SafeUnmarshalBase64 @XDRSubmitter " + err.Error())
		}
		TDP[i].MapIdentifier = TDP[i].Identifier
		TDP[i].TxnType = "10"
		TDP[i].Status = "success"
		display := stellarExecuter.ConcreteSubmitXDR{XDR: TDP[i].XDR}
		response := display.SubmitXDR(TDP[i].TxnType)
		ret = response
		if response.Error.Code == 400 {
			log.Error("Error got 400 for SubmitXDR @XDRSubmitter ")
			TDP[i].Status = "pending"
			status = append(status, false)
		} else {
			TDP[i].TxnHash = response.TXNID
			UserTxnHashes = append(UserTxnHashes, response.TXNID)
			status = append(status, true)
		}
	}

	if checkBoolArray(status) {
		go func() {
			for i := 0; i < len(TDP); i++ {
				p := object.GetLastTransactionbyIdentifierAndTenantIdAndTxnType(TDP[i].Identifier, TDP[i].TenantID, TDP[i].TxnType)
				p.Then(func(data interface{}) interface{} {
					///ASSIGN PREVIOUS MANAGE DATA BUILDER
					result := data.(model.TransactionCollectionBody)
					TDP[i].PreviousTxnHash = result.TxnHash
					return nil
				}).Catch(func(error error) error {
					log.Error("Error @ GetLastTransactionbyIdentifier @XDRSubmitter " + error.Error())
					TDP[i].PreviousTxnHash = ""
					return error
				})
				p.Await()
				PreviousTXNBuilder := txnbuild.ManageData{
					Name:  "PreviousTXN",
					Value: []byte(TDP[i].PreviousTxnHash),
				}
				TypeTxnBuilder := txnbuild.ManageData{
					Name:  "Type",
					Value: []byte("G" + TDP[i].TxnType),
				}

				CurrentTXNBuilder := txnbuild.ManageData{
					Name:  "CurrentTXN",
					Value: []byte(UserTxnHashes[i]),
				}

				// BUILD THE GATEWAY XDR
				tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
					SourceAccount:        &account,
					IncrementSequenceNum: true,
					Operations:           []txnbuild.Operation{&PreviousTXNBuilder, &TypeTxnBuilder, &CurrentTXNBuilder},
					BaseFee:              constants.MinBaseFee,
					Memo:                 nil,
					Preconditions:        txnbuild.Preconditions{TimeBounds: constants.TransactionTimeOut},
				})
				if err != nil {
					log.Error("Error @ builder @XDRSubmitter " + err.Error())
					return
				}

				// SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
				GatewayTXE, err := tx.Sign(commons.GetStellarNetwork(), tracifiedAccount)
				if err != nil {
					log.Error("Error @ Sign @XDRSubmitter " + err.Error())

					TDP[i].TxnHash = UserTxnHashes[i]
					TDP[i].Status = "Pending"

					///INSERT INTO TRANSACTION COLLECTION
					err2 := object.InsertTransaction(TDP[i])
					if err2 != nil {
						log.Error("Error @ InsertTransaction @XDRSubmitter " + err2.Error())
					}
				}
				// CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
				txeB64, err := GatewayTXE.Base64()
				if err != nil {
					log.Error("Error while convert GatewayTXE to base64 @XDRSubmitter " + err.Error())
					TDP[i].TxnHash = UserTxnHashes[i]
					TDP[i].Status = "Pending"

					///INSERT INTO TRANSACTION COLLECTION
					err2 := object.InsertTransaction(TDP[i])
					if err2 != nil {
						log.Error("Error @ InsertTransaction @XDRSubmitter " + err2.Error())
					}
				}

				// SUBMIT THE GATEWAY'S SIGNED XDR
				display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
				response1 := display1.SubmitXDR("G" + TDP[i].TxnType)

				if response1.Error.Code == 400 {
					log.Error("Error got 400 for ConcreteSubmitXDR @XDRSubmitter")
					TDP[i].TxnHash = UserTxnHashes[i]
					TDP[i].Status = "Pending"

					///INSERT INTO TRANSACTION COLLECTION
					err2 := object.InsertTransaction(TDP[i])
					if err2 != nil {
						log.Error("Error @ InsertTransaction @XDRSubmitter " + err2.Error())
					}
				} else {
					// UPDATE THE TRANSACTION COLLECTION WITH TXN HASH
					TDP[i].TxnHash = response1.TXNID
					TDP[i].Status = "done"
					TDP[i].XDR = txeB64
					TDP[i].SequenceNo = int64(GatewayTXE.SequenceNumber())

					///INSERT INTO TRANSACTION COLLECTION
					err2 := object.InsertTransaction(TDP[i])
					if err2 != nil {
						log.Error("Error @ InsertTransaction @XDRSubmitter " + err2.Error())
					}
				}
			}
		}()
	}

	return checkBoolArray(status), ret
}

/*
@author - Azeem Ashraf
@desc - checks the multiple boolean indexes in an array and returns the combined result.
*/
func checkBoolArray(array []bool) bool {
	isMatch := true
	for i := 0; i < len(array); i++ {
		if array[i] == false {
			isMatch = false
			return isMatch
		}
	}
	return isMatch
}
