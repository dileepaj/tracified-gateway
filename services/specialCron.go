package services

import (
	"github.com/dileepaj/tracified-gateway/commons"
	log "github.com/sirupsen/logrus"
	"strconv"
	// "time"
	// "fmt"
	// "github.com/stellar/go/xdr"


	"github.com/stellar/go/build"
	// "fmt"
	"github.com/dileepaj/tracified-gateway/adminDAO"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
)

//CheckTempOrphan ...
func CheckTempOrphan() {
	log.Debug("=================== CheckTempOrphan ==================")
	netClient := commons.GetHorizonClient()
	//clientList := commons.CallAdminBE()
	adminDBConnectionObj := adminDAO.Connection{}
	clientList := adminDBConnectionObj.GetPublicKeysOfFO()
	log.Info("PK count : " + strconv.Itoa(len(clientList)))
	object := dao.Connection{}
	//loop through clients
	for _, address := range clientList {
		//load horizon account
		account, err := netClient.LoadAccount(address)
		if err != nil {
			log.Error("Error while loading account from horizon "+err.Error())
		}else {
			//log.Println("Current Sequence for address:", address)
			//log.Println(account.Sequence)
			seq, err := strconv.Atoi(account.Sequence)
			if err != nil {
				log.Error("Error while convert string to int "+err.Error())
			}
			stop := false //for infinite loop
			//loop through sequence incrementally and see match
			for i := seq + 1; ; i++ {
				p := object.GetSpecialForPkAndSeq(address, int64(i))
				p.Then(func(data interface{}) interface{} {
					result := data.(model.TransactionCollectionBody)
					var UserTxnHash string
					///HARDCODED CREDENTIALS
					publicKey := constants.PublicKey
					secretKey := constants.SecretKey
					switch result.TxnType {
					case "0":
						display := stellarExecuter.ConcreteSubmitXDR{XDR: result.XDR}
						response := display.SubmitXDR(result.TxnType)
						UserTxnHash = response.TXNID
						if response.Error.Code == 400 {
							log.Println("response.Error.Code 400 for SubmitXDR")
							break
						}

						var PreviousTXNBuilder build.ManageDataBuilder
						PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(""))

						//BUILD THE GATEWAY XDR
						tx, err := build.Transaction(
							build.PublicNetwork,
							build.SourceAccount{publicKey},
							build.AutoSequence{netClient},
							build.SetData("Type", []byte("G"+result.TxnType)),
							PreviousTXNBuilder,
							build.SetData("CurrentTXN", []byte(UserTxnHash)),
						)

						//SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
						GatewayTXE, err := tx.Sign(secretKey)
						if err != nil {
							log.Println("Error while getting GatewayTXE by secretKey "+err.Error())
							break
						}

						//CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
						txeB64, err := GatewayTXE.Base64()
						if err != nil {
							log.Println("Error while converting GatewayTXE to base64 "+err.Error())
							break
						}

						//SUBMIT THE GATEWAY'S SIGNED XDR
						display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
						response1 := display1.SubmitXDR("G"+result.TxnType)

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
						}

					case "2":

						var PreviousTXNBuilder build.ManageDataBuilder

						// var PreviousTxn string
						p := object.GetLastTransactionbyIdentifier(result.Identifier)
						p.Then(func(data interface{}) interface{} {
							///ASSIGN PREVIOUS MANAGE DATA BUILDER
							res := data.(model.TransactionCollectionBody)
							PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(res.TxnHash))

							result.PreviousTxnHash = res.TxnHash
							return nil
						}).Catch(func(error error) error {
							PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(""))

							return error
						})
						p.Await()

						display := stellarExecuter.ConcreteSubmitXDR{XDR: result.XDR}
						response := display.SubmitXDR(result.TxnType)
						UserTxnHash = response.TXNID
						if response.Error.Code == 400 {
							log.Println("Response code 400 for SubmitXDR")
							break
						}
						//BUILD THE GATEWAY XDR
						tx, err := build.Transaction(
							build.PublicNetwork,
							build.SourceAccount{publicKey},
							build.AutoSequence{netClient},
							build.SetData("Type", []byte("G"+result.TxnType)),
							PreviousTXNBuilder,
							build.SetData("CurrentTXN", []byte(UserTxnHash)),
						)

						//SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
						GatewayTXE, err := tx.Sign(secretKey)
						if err != nil {
							log.Println("Error while getting GatewayTXE by secretKey " + err.Error())
							break
						}

						//CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
						txeB64, err := GatewayTXE.Base64()
						if err != nil {
							log.Println("Error while converting GatewayTXE to base64 " + err.Error())
							break
						}

						//SUBMIT THE GATEWAY'S SIGNED XDR
						display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
						response1 := display1.SubmitXDR("G"+result.TxnType)

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
						}
					case "9":

						var PreviousTXNBuilder build.ManageDataBuilder

						// var PreviousTxn string
						p := object.GetLastTransactionbyIdentifier(result.Identifier)
						p.Then(func(data interface{}) interface{} {
							///ASSIGN PREVIOUS MANAGE DATA BUILDER
							res := data.(model.TransactionCollectionBody)
							PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(res.TxnHash))
							result.PreviousTxnHash = res.TxnHash

							return nil
						}).Catch(func(error error) error {
							PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(""))

							return error
						})
						p.Await()

						display := stellarExecuter.ConcreteSubmitXDR{XDR: result.XDR}
						response := display.SubmitXDR(result.TxnType)
						UserTxnHash = response.TXNID
						if response.Error.Code == 400 {
							log.Println("400 SubmitXDR")
							break
						}
						//BUILD THE GATEWAY XDR
						tx, err := build.Transaction(
							build.PublicNetwork,
							build.SourceAccount{publicKey},
							build.AutoSequence{netClient},
							build.SetData("Type", []byte("G"+result.TxnType)),
							PreviousTXNBuilder,
							build.SetData("CurrentTXN", []byte(UserTxnHash)),
						)

						//SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
						GatewayTXE, err := tx.Sign(secretKey)
						if err != nil {
							log.Println("Error while getting GatewayTXE " + err.Error())
							break
						}

						//CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
						txeB64, err := GatewayTXE.Base64()
						if err != nil {
							log.Println("Error while converting to base64 " + err.Error())
							break
						}

						//SUBMIT THE GATEWAY'S SIGNED XDR
						display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
						response1 := display1.SubmitXDR("G"+result.TxnType)

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
						}

					}
					return nil
				}).Catch(func(error error) error {
					log.Error("Error while GetSpecialForPkAndSeq "+error.Error())
					// return error
					//log.Println("No transactions in the scheduler")
					stop = true //to break loop
					return nil
				})
				p.Await()
				if stop {
					break
				}
			}
		}
	}
}
