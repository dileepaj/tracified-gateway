package services

import (
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"

	// "time"
	// "fmt"
	// "github.com/stellar/go/xdr"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"

	// "fmt"
	"github.com/dileepaj/tracified-gateway/adminDAO"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
)

// CheckTempOrphan ...
func CheckTempOrphan() {
	if (commons.GoDotEnvVariable("LOGSTYPE")=="DEBUG"){
		log.Debug("=================== CheckTempOrphan ==================")
	}
	adminDBConnectionObj := adminDAO.Connection{}
	clientList := adminDBConnectionObj.GetPublicKeysOfFO()
	// log.Info("PK count : " + strconv.Itoa(len(clientList)))
	object := dao.Connection{}
	// loop through clients
	for _, address := range clientList {
		kp, _ := keypair.Parse(address)
		client := commons.GetHorizonClient()
		ar := horizonclient.AccountRequest{AccountID: kp.Address()}
		sourceAccount, err := client.AccountDetail(ar)
		if err != nil {
			// log.Error("Error while loading account from horizon " + err.Error())
		} else {
			seq, err := strconv.Atoi(fmt.Sprint(sourceAccount.Sequence))
			if err != nil {
				log.Error("Error while convert string to int " + err.Error())
			}
			stop := false // for infinite loop
			// loop through sequence incrementally and see match
			for i := seq + 1; ; i++ {
				if commons.GoDotEnvVariable("LOGSTYPE")=="DEBUG"{
				log.Info("Find tempOrphan by ", kp.Address(), "    -   ", i)
				}
				data, errorAsync := object.GetSpecialForPkAndSeq(kp.Address(), int64(i)).Then(func(data interface{}) interface{} {
					return data
				}).Await()
				if errorAsync != nil {
					stop = true // to break loop
				} else if data == nil {
					stop = true
				} else {
					result := data.(model.TransactionCollectionBody)
					var UserTxnHash string

					///HARDCODED CREDENTIALS
					publicKey := constants.PublicKey
					secretKey := constants.SecretKey
					tracifiedAccount, err := keypair.ParseFull(secretKey)
					if err != nil {
						log.Error("Account loading issue ", err)
					}
					client := commons.GetHorizonClient()
					pubaccountRequest := horizonclient.AccountRequest{AccountID: publicKey}
					pubaccount, err := client.AccountDetail(pubaccountRequest)
					if commons.GoDotEnvVariable("LOGSTYPE") == "DEBUG" {
						log.Info("clientList PublicKey ", address, " Sequence number ", seq)
						log.Info("PublicKey key of XDR ", ar.AccountID)
						log.Info("Sequence number ", i)
						log.Info("Type of XDR ", result.TxnType)
					}
					switch result.TxnType {
					case "0":

						display := stellarExecuter.ConcreteSubmitXDR{XDR: result.XDR}
						response := display.SubmitXDR(result.TxnType)
						UserTxnHash = response.TXNID
						if commons.GoDotEnvVariable("LOGSTYPE")=="DEBUG"{
							log.Info("type 0 submission ", response)
						}
						if response.Error.Code == 400 {
							log.Println("response.Error.Code 400 for SubmitXDR")
							break
						}
						// var PreviousTXNBuilder txnbuild.ManageData
						PreviousTXNBuilder := txnbuild.ManageData{
							Name:  "PreviousTXN",
							Value: []byte(""),
						}
						TypeTxnBuilder := txnbuild.ManageData{
							Name:  "Type",
							Value: []byte("G" + result.TxnType),
						}
						CurrentTXNBuilder := txnbuild.ManageData{
							Name:  "CurrentTXN",
							Value: []byte(UserTxnHash),
						}
						// BUILD THE GATEWAY XDR
						tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
							SourceAccount:        &pubaccount,
							IncrementSequenceNum: true,
							Operations:           []txnbuild.Operation{&PreviousTXNBuilder, &TypeTxnBuilder, &CurrentTXNBuilder},
							BaseFee:              constants.MinBaseFee,
							Memo:                 nil,
							Preconditions:        txnbuild.Preconditions{TimeBounds: constants.TransactionTimeOut},
						})
						if err != nil {
							log.Println("Error while building XDR " + err.Error())
							break
						}
						// SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
						GatewayTXE, err := tx.Sign(commons.GetStellarNetwork(), tracifiedAccount)
						if err != nil {
							log.Println("Error while getting GatewayTXE by secretKey " + err.Error())
							break
						}
						// CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
						txeB64, err := GatewayTXE.Base64()
						if err != nil {
							log.Println("Error while converting GatewayTXE to base64 " + err.Error())
							break
						}
						// SUBMIT THE GATEWAY'S SIGNED XDR
						display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
						response1 := display1.SubmitXDR("G" + result.TxnType)
						if response1.Error.Code == 400 {
							log.Println("Error code 400 for SubmitXDR")
							break
						}
						result.TxnHash = response1.TXNID
						result.Status = "done"
						///INSERT INTO TRANSACTION COLLECTION
						err2 := object.InsertTransaction(result)
						if err2 != nil {
							log.Println("Error while InsertTransaction " + err2.Error())
							break
						} else {
							err := object.RemoveFromTempOrphanList(result.PublicKey, result.SequenceNo)
							if err != nil {
								log.Println("Error while RemoveFromTempOrphanList " + err.Error())
								break
							}
						}
					case "2":
						data, errorLastTXN := object.GetLastTransactionbyIdentifier(result.Identifier).Then(func(data interface{}) interface{} {
							return data
						}).Await()
						var PreviousTXNBuilder txnbuild.ManageData
						if errorLastTXN != nil || data == nil {
							PreviousTXNBuilder = txnbuild.ManageData{
								Name:  "PreviousTXN",
								Value: []byte(""),
							}
						} else {
							///ASSIGN PREVIOUS MANAGE DATA BUILDER
							res := data.(model.TransactionCollectionBody)
							PreviousTXNBuilder = txnbuild.ManageData{
								Name:  "PreviousTXN",
								Value: []byte(res.TxnHash),
							}
							result.PreviousTxnHash = res.TxnHash
						}
						TypeTxnBuilder := txnbuild.ManageData{
							Name:  "Type",
							Value: []byte("G" + result.TxnType),
						}
						display := stellarExecuter.ConcreteSubmitXDR{XDR: result.XDR}
						response := display.SubmitXDR(result.TxnType)
						if commons.GoDotEnvVariable("LOGSTYPE")=="DEBUG"{
						log.Info("type 1 submission ", response)
						}
						UserTxnHash = response.TXNID

						CurrentTXNBuilder := txnbuild.ManageData{
							Name:  "CurrentTXN",
							Value: []byte(UserTxnHash),
						}
						if response.Error.Code == 400 {
							log.Println("Response code 400 for SubmitXDR")
							break
						}
						// BUILD THE GATEWAY XDR
						tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
							SourceAccount:        &pubaccount,
							IncrementSequenceNum: true,
							Operations:           []txnbuild.Operation{&PreviousTXNBuilder, &TypeTxnBuilder, &CurrentTXNBuilder},
							BaseFee:              constants.MinBaseFee,
							Memo:                 nil,
							Preconditions:        txnbuild.Preconditions{TimeBounds: constants.TransactionTimeOut},
						})
						if err != nil {
							log.Println("Error while buliding XDR " + err.Error())
							break
						}
						// SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
						GatewayTXE, err := tx.Sign(commons.GetStellarNetwork(), tracifiedAccount)
						if err != nil {
							log.Println("Error while getting GatewayTXE by secretKey " + err.Error())
							break
						}

						// CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
						txeB64, err := GatewayTXE.Base64()
						if err != nil {
							log.Println("Error while converting GatewayTXE to base64 " + err.Error())
							break
						}

						// SUBMIT THE GATEWAY'S SIGNED XDR
						display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
						response1 := display1.SubmitXDR("G" + result.TxnType)

						if response1.Error.Code == 400 {
							log.Println("Error response code 400 while SubmitXDR")
							break
						}
						result.TxnHash = response1.TXNID
						result.Status = "done"
						///INSERT INTO TRANSACTION COLLECTION
						err2 := object.InsertTransaction(result)
						if err2 != nil {
							log.Println("Error while InsertTransaction " + err2.Error())
							break
						} else {
							err := object.RemoveFromTempOrphanList(result.PublicKey, result.SequenceNo)
							if err != nil {
								log.Println("Error while RemoveFromTempOrphanList " + err.Error())
								break
							}
						}
					case "9":
						var PreviousTXNBuilder txnbuild.ManageData
						// var PreviousTxn string
						data, errorLastTXN := object.GetLastTransactionbyIdentifier(result.Identifier).Then(func(data interface{}) interface{} {
							return data
						}).Await()

						if errorLastTXN != nil || data == nil {
							PreviousTXNBuilder = txnbuild.ManageData{
								Name:  "PreviousTXN",
								Value: []byte(""),
							}
						} else {
							///ASSIGN PREVIOUS MANAGE DATA BUILDER
							res := data.(model.TransactionCollectionBody)
							PreviousTXNBuilder = txnbuild.ManageData{
								Name:  "PreviousTXN",
								Value: []byte(res.TxnHash),
							}
							result.PreviousTxnHash = res.TxnHash
						}
						TypeTxnBuilder := txnbuild.ManageData{
							Name:  "Type",
							Value: []byte("G" + result.TxnType),
						}
						CurrentTXNBuilder := txnbuild.ManageData{
							Name:  "CurrentTXN",
							Value: []byte(UserTxnHash),
						}
						display := stellarExecuter.ConcreteSubmitXDR{XDR: result.XDR}
						response := display.SubmitXDR(result.TxnType)
						UserTxnHash = response.TXNID
						if response.Error.Code == 400 {
							log.Println("400 SubmitXDR")
							break
						}
						// BUILD THE GATEWAY XDR
						tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
							SourceAccount:        &pubaccount,
							IncrementSequenceNum: true,
							Operations:           []txnbuild.Operation{&PreviousTXNBuilder, &TypeTxnBuilder, &CurrentTXNBuilder},
							BaseFee:              constants.MinBaseFee,
							Memo:                 nil,
							Preconditions:        txnbuild.Preconditions{TimeBounds: constants.TransactionTimeOut},
						})
						if err != nil {
							log.Println("Error while buliding XDR " + err.Error())
							break
						}
						// SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
						GatewayTXE, err := tx.Sign(commons.GetStellarNetwork(), tracifiedAccount)
						if err != nil {
							log.Println("Error while getting GatewayTXE " + err.Error())
							break
						}

						// CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
						txeB64, err := GatewayTXE.Base64()
						if err != nil {
							log.Println("Error while converting to base64 " + err.Error())
							break
						}
						// SUBMIT THE GATEWAY'S SIGNED XDR
						display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
						response1 := display1.SubmitXDR("G" + result.TxnType)
						if commons.GoDotEnvVariable("LOGSTYPE")=="DEBUG"{
							log.Info("type 9 submission ", response1)
						}
						if response1.Error.Code == 400 {
							log.Println("400 from SubmitXDR")
							break
						}
						result.TxnHash = response1.TXNID
						result.Status = "done"
						///INSERT INTO TRANSACTION COLLECTION
						err2 := object.InsertTransaction(result)
						if err2 != nil {
							log.Println("Error while InsertTransaction " + err2.Error())
							break
						} else {
							err := object.RemoveFromTempOrphanList(result.PublicKey, result.SequenceNo)
							if err != nil {
								log.Println("Error while RemoveFromTempOrphanList " + err.Error())
								break
							}
						}

					}
				}

				if stop {
					break
				}
			}
		}
	}
}
