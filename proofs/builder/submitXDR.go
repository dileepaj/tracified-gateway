package builder

import (
	// "encoding/json"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	// "net/http"
	// "strconv"

	// "github.com/dileepaj/tracified-gateway/api/apiModel"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"

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

/*SubmitTransfer - Deprecated
@author - Azeem Ashraf
*/

// func (AP *AbstractXDRSubmiter) SubmitTransfer() bool {
// 	var Done bool
// 	object := dao.Connection{}
// 	var copy []model.TransactionCollectionBody

// 	var UserTxnHashes []string
// 	// var PreviousTxn []string
// 	///HARDCODED CREDENTIALS
// 	publicKey := constants.PublicKey
// 	secretKey := constants.SecretKey
// 	// var result model.SubmitXDRResponse

// 	for i, TxnBody := range AP.TxnBody {
// 		var TDP model.TransactionCollectionBody
// 		var txe xdr.Transaction

// 		//decode the XDR
// 		err := xdr.SafeUnmarshalBase64(TxnBody.XDR, &txe)
// 		if err != nil {
// 			fmt.Println(err)
// 		}

// 		//GET THE TYPE AND IDENTIFIER FROM THE XDR
// 		TxnBody.PublicKey = txe.SourceAccount.Address()
// 		// TDP.PreviousTxnHash=
// 		TxnType := strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
// 		Identifier := strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[1].Body.ManageDataOp.DataValue), "&")
// 		TxnBody.Identifier = Identifier
// 		TxnBody.TxnType = TxnType
// 		TxnBody.Status = "pending"
// 		// TxnBody.TdpId=
// 		AP.TxnBody[i].Identifier = Identifier
// 		AP.TxnBody[i].TxnType = TxnType

// 		fmt.Println(Identifier)
// 		p := object.GetLastTransactionbyIdentifier(Identifier)
// 		p.Then(func(data interface{}) interface{} {
// 			///ASSIGN PREVIOU y S MANAGE DATA BUILDER

// 			result := data.(model.TransactionCollectionBody)
// 			if result.TxnHash == "" {
// 				fmt.Println("Sending to Orphanage!")
// 				AP.TxnBody[i].Orphan = true
// 				// TxnBody.Orphan = true

// 				//INSERT THE TXN INTO THE BUFFER
// 				err1 := object.InsertToOrphan(TxnBody)
// 				if err1 != nil {
// 					TDP.Status = "failed"
// 				} else {
// 					Done = true

// 				}

// 			} else {
// 				TxnBody.PreviousTxnHash = result.TxnHash
// 				fmt.Println("Previous TXN: " + result.TxnHash)

// 				copy = append(copy, TxnBody)
// 				///INSERT INTO TRANSACTION COLLECTION
// 				err1 := object.InsertTransaction(TxnBody)
// 				if err1 != nil {
// 					TDP.Status = "failed"
// 				}
// 				//SUBMIT THE FIRST XDR SIGNED BY THE USER
// 				display := stellarExecuter.ConcreteSubmitXDR{XDR: TxnBody.XDR}
// 				result1 := display.SubmitXDR(false,AP.TxnBody[i].TxnType)
// 				UserTxnHashes = append(UserTxnHashes, result1.TXNID)

// 				if result1.Error.Code != 404 {
// 					Done = true
// 					// return Done
// 				}

// 				var PreviousTXNBuilder build.ManageDataBuilder

// 				PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(copy[i].PreviousTxnHash))

// 				//BUILD THE GATEWAY XDR
// 				tx, err := build.Transaction(
// 					commons.GetHorizonNetwork(),
// 					build.SourceAccount{publicKey},
// 					build.AutoSequence{commons.GetHorizonClient()},
// 					build.SetData("Type", []byte("G"+TxnBody.TxnType)),
// 					PreviousTXNBuilder,
// 					build.SetData("CurrentTXN", []byte(UserTxnHashes[i])),
// 				)

// 				//SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
// 				GatewayTXE, err := tx.Sign(secretKey)
// 				if err != nil {
// 					fmt.Println(err)
// 				}
// 				//CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
// 				txeB64, err := GatewayTXE.Base64()
// 				if err != nil {
// 					fmt.Println(err)
// 				}

// 				//SUBMIT THE GATEWAY'S SIGNED XDR
// 				display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
// 				response1 := display1.SubmitXDR(false,"G"+AP.TxnBody[i].TxnType)

// 				if response1.Error.Code == 404 {
// 					TxnBody.Status = "pending"
// 					Done = false
// 				} else {
// 					//UPDATE THE TRANSACTION COLLECTION WITH TXN HASH
// 					TxnBody.TxnHash = response1.TXNID

// 					upd := model.TransactionCollectionBody{
// 						TxnHash: response1.TXNID,
// 						Status:  "done",
// 					}
// 					err2 := object.UpdateTransaction(copy[i], upd)
// 					if err2 != nil {
// 						TxnBody.Status = "pending"
// 					} else {
// 						TxnBody.Status = "done"
// 					}
// 					Done = true
// 				}
// 			}

// 			return nil
// 		}).Catch(func(error error) error {
// 			///ASSIGN PREVIOUS MANAGE DATA BUILDER - LEAVE IT EMPTY
// 			fmt.Println("Sending to Orphanage!")
// 			AP.TxnBody[i].Orphan = true
// 			// TxnBody.Orphan = true

// 			//INSERT THE TXN INTO THE BUFFER
// 			err1 := object.InsertToOrphan(TxnBody)
// 			if err1 != nil {
// 				TDP.Status = "failed"
// 			} else {
// 				Done = true

// 			}

// 			return error
// 		})
// 		p.Await()

// 	}

// 	return Done
// }

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
	// Get information about the account we just created
	//accountRequest := horizonclient.AccountRequest{AccountID: publicKey}
	//account, err := netClient.AccountDetail(accountRequest)

	kp,_ := keypair.Parse(publicKey)
	client := horizonclient.DefaultTestNetClient
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

		TDP[i].PublicKey = txe.SourceAccount.Address()
		TxnType := strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
		Identifier := strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[1].Body.ManageDataOp.DataValue), "&")
		TDP[i].Identifier = Identifier
		TDP[i].TxnType = TxnType
		TDP[i].Status = "success"

		display := stellarExecuter.ConcreteSubmitXDR{XDR: TDP[i].XDR}

		response := display.SubmitXDR(TDP[i].TxnType)
		ret = response
		fmt.Println(ret.Error.Code)
		if response.Error.Code == 400 {
			log.Error("Error got 400 for SubmitXDR @XDRSubmitter ")
			TDP[i].Status = "pending"
			status = append(status, false)
		} else {
			TDP[i].TxnHash = response.TXNID
			UserTxnHashes = append(UserTxnHashes, response.TXNID)
			status = append(status, true)
			// TDP[i].Status = "done"
			// err1 := object.InsertTransaction(TDP[i])
			// if err1 != nil {
			// 	status = append(status, false)
			// }
		}
	}

	if checkBoolArray(status) {
		go func() {
			for i := 0; i < len(TDP); i++ {
				p := object.GetLastTransactionbyIdentifier(TDP[i].Identifier)
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
					BaseFee:              txnbuild.MinBaseFee,
					Memo:                 nil,
					Preconditions:        txnbuild.Preconditions{},
				})
				if err != nil {
					log.Error("Error @ builder @XDRSubmitter " + err.Error())
					return
				}

				// SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
				GatewayTXE, err := tx.Sign(secretKey)
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
