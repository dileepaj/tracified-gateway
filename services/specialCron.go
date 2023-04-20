package services

import (
	"fmt"
	"strconv"

	"github.com/dileepaj/tracified-gateway/adminDAO"
	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
	log "github.com/sirupsen/logrus"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
)

// CheckTempOrphan ...
func CheckTempOrphan() {
	if commons.GoDotEnvVariable("LOGSTYPE") == "DEBUG" {
		log.Debug("=================== CheckTempOrphan ==================")
	}
	adminDBConnectionObj := adminDAO.Connection{}
	clientList := adminDBConnectionObj.GetPublicKeysOfFO()
	object := dao.Connection{}
	// loop through clients
	for _, address := range clientList {
		kp, _ := keypair.Parse(address)
		client := commons.GetHorizonClient()
		ar := horizonclient.AccountRequest{AccountID: kp.Address()}
		sourceAccount, err := client.AccountDetail(ar)
		if err != nil {
			if commons.GoDotEnvVariable("LOGSTYPE") == "DEBUG" {
			log.Error("Error while loading account from horizon " + err.Error())
			}
		} else {
			seq, err := strconv.Atoi(fmt.Sprint(sourceAccount.Sequence))
			if err != nil {
				log.Debug("Error while convert string to int " + err.Error())
			}
			stop := false // for infinite loop
			// loop through sequence incrementally and see match
			for i := seq + 1; ; i++ {
				if commons.GoDotEnvVariable("LOGSTYPE") == "DEBUG" {
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
					publicKey := constants.PublicKey
					secretKey := constants.SecretKey
					tracifiedAccount, err := keypair.ParseFull(secretKey)
					if err != nil {
						log.Error(err)
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
						if commons.GoDotEnvVariable("LOGSTYPE") == "DEBUG" {
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
							log.Println("Error code 400 for SubmitXDR")
							break
						}
						result.TxnHash = response1.TXNID
						result.Status = "done"
						// INSERT INTO TRANSACTION COLLECTION
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
						data, errorLastTXN := object.GetLastTransactionbyIdentifierAndTenantId(result.Identifier, result.TenantID).Then(func(data interface{}) interface{} {
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
						if commons.GoDotEnvVariable("LOGSTYPE") == "DEBUG" {
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
						data, errorLastTXN := object.GetLastTransactionbyIdentifierAndTenantId(result.Identifier, result.TenantID).Then(func(data interface{}) interface{} {
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
						if commons.GoDotEnvVariable("LOGSTYPE") == "DEBUG" {
							log.Info("type 9 submission ", response1)
						}
						if response1.Error.Code != 200 {
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
					case "5":
						var UserSplitTxnHashes string
						var PreviousTxn string
						// ParentIdentifier = Identifier
						pData, errAsnc := object.GetLastTransactionbyIdentifierAndTenantId(result.Identifier, result.TenantID).Then(func(data interface{}) interface{} {
							return data
						}).Await()
						if pData == nil || errAsnc != nil {
							log.Error("Error @GetLastTransactionbyIdentifier @SubmitSplit ")
							// ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
							// DUE TO THE CHILD HAVING A NEW IDENTIFIER
							PreviousTxn = ""
							result.PreviousTxnHash = ""
						} else {
							// ASSIGN PREVIOUS MANAGE DATA BUILDER
							result := pData.(model.TransactionCollectionBody)
							PreviousTxn = result.TxnHash
							result.PreviousTxnHash = result.TxnHash
						}
						// SUBMIT THE FIRST XDR SIGNED BY THE USER
						display := stellarExecuter.ConcreteSubmitXDR{XDR: result.XDR}
						result1 := display.SubmitXDR(result.TxnType)
						UserSplitTxnHashes = result1.TXNID
						if result1.Error.Code == 400 {
							log.Error("Index[" + strconv.Itoa(i) + "] TXN: Blockchain Transaction Failed!")
						} else {
							log.Info((i + 1), " Submitted")
							log.Info((i + 1), " PreviousTxn ", PreviousTxn)
							var SplitParentProfile string
							var PreviousSplitProfile string
							/*
								When constructing a backlink transaction(put from gateway) for a split, it is important to exclude the split-parent transaction as its previous transaction.
								Instead, you should obtain the most recent transaction that is specific to the identifier and disregard the split-parent transaction.
							*/
							previousTXNBuilder := txnbuild.ManageData{Name: "PreviousTXN", Value: []byte(PreviousTxn)}
							typeTXNBuilder := txnbuild.ManageData{Name: "Type", Value: []byte("G" + result.TxnType)}
							currentTXNBuilder := txnbuild.ManageData{Name: "CurrentTXN", Value: []byte(UserSplitTxnHashes)}
							identifierTXNBuilder := txnbuild.ManageData{Name: "Identifier", Value: []byte(result.Identifier)}
							profileIDTXNBuilder := txnbuild.ManageData{Name: "ProfileID", Value: []byte(result.ProfileID)}
							PreviousProfileTXNBuilder := txnbuild.ManageData{Name: "PreviousProfile", Value: []byte(PreviousSplitProfile)}
							result.PreviousTxnHash = PreviousTxn
							// ASSIGN THE PREVIOUS PROFILE ID USING THE PARENT FOR THE CHILDREN AND A DB CALL FOR PARENT
							if result.TxnType == "5" {
								PreviousSplitProfile = ""
								SplitParentProfile = result.ProfileID
							} else {
								PreviousSplitProfile = SplitParentProfile
							}
							// BUILD THE GATEWAY XDR
							tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
								SourceAccount:        &pubaccount,
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
								result.TxnHash = UserSplitTxnHashes
								result.Status = "Pending"

								///INSERT INTO TRANSACTION COLLECTION
								err2 := object.InsertTransaction(result)
								if err2 != nil {
									log.Error("Error @InsertTransaction @SubmitSplit " + err2.Error())
								}
							}
							// CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
							txeB64, err := GatewayTXE.Base64()
							if err != nil {
								log.Error("Error @GatewayTXE.Base64 @SubmitSplit " + err.Error())
								result.TxnHash = UserSplitTxnHashes
								result.Status = "Pending"

								///INSERT INTO TRANSACTION COLLECTION
								err2 := object.InsertTransaction(result)
								if err2 != nil {
									log.Error("Error @InsertTransaction @SubmitSplit " + err2.Error())
								}
							}
							// SUBMIT THE GATEWAY'S SIGNED XDR
							display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
							response1 := display1.SubmitXDR("G" + result.TxnType)
							if response1.Error.Code != 200 {
								result.TxnHash = UserSplitTxnHashes
								result.Status = "Pending"
								///INSERT INTO TRANSACTION COLLECTION
								err2 := object.InsertTransaction(result)
								if err2 != nil {
									log.Error("Error @InsertTransaction @SubmitSplit " + err2.Error())
								}
								break
							} else {
								// UPDATE THE TRANSACTION COLLECTION WITH TXN HASH
								result.TxnHash = response1.TXNID
								if result.TxnType == "5" {
									PreviousTxn = response1.TXNID
								}
								///INSERT INTO TRANSACTION COLLECTION
								err1 := object.InsertTransaction(result)
								if err1 != nil {
									log.Error("Error @InsertTransaction @SubmitSplit " + err1.Error())
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
					case "6":
						var UserSplitTxnHashes string
						var PreviousTxn string
						// SUBMIT THE FIRST XDR SIGNED BY THE USER
						display := stellarExecuter.ConcreteSubmitXDR{XDR: result.XDR}
						result1 := display.SubmitXDR(result.TxnType)
						UserSplitTxnHashes = result1.TXNID

						if result1.Error.Code == 400 {
							log.Error("Index[" + strconv.Itoa(i) + "] TXN: Blockchain Transaction Failed!")
						} else {
							var id apiModel.IdentifierModel
							id.MapValue = result.Identifier
							id.Identifier = result.MapIdentifier
							err3 := object.InsertIdentifier(id)
							if err3 != nil {
								log.Error("identifier map failed" + err3.Error())
							}
							log.Info((i + 1), " Submitted")
							log.Info((i + 1), " PreviousTxn ", PreviousTxn)
							var SplitParentProfile string
							var PreviousSplitProfile string
							/*
								When constructing a backlink transaction(put from gateway) for a split, it is important to exclude the split-parent transaction as its previous transaction.
								Instead, you should obtain the most recent transaction that is specific to the identifier and disregard the split-parent transaction.
							*/
							backlinkData, errAsnc := object.GetLastTransactionbyIdentifierNotSplitParent(result.FromIdentifier1, result.TenantID).Then(func(data interface{}) interface{} {
								return data
							}).Await()
							if backlinkData == nil || errAsnc != nil {
								log.Info("Can not find transaction form database ", "build Split")
							} else {
								result := backlinkData.(model.TransactionCollectionBody)
								PreviousTxn = result.TxnHash
								result.PreviousTxnHash = result.TxnHash
							}

							previousTXNBuilder := txnbuild.ManageData{Name: "PreviousTXN", Value: []byte(PreviousTxn)}
							typeTXNBuilder := txnbuild.ManageData{Name: "Type", Value: []byte("G" + result.TxnType)}
							currentTXNBuilder := txnbuild.ManageData{Name: "CurrentTXN", Value: []byte(UserSplitTxnHashes)}
							identifierTXNBuilder := txnbuild.ManageData{Name: "Identifier", Value: []byte(result.Identifier)}
							profileIDTXNBuilder := txnbuild.ManageData{Name: "ProfileID", Value: []byte(result.ProfileID)}
							PreviousProfileTXNBuilder := txnbuild.ManageData{Name: "PreviousProfile", Value: []byte(PreviousSplitProfile)}
							result.PreviousTxnHash = PreviousTxn
							// ASSIGN THE PREVIOUS PROFILE ID USING THE PARENT FOR THE CHILDREN AND A DB CALL FOR PARENT
							if result.TxnType == "5" {
								PreviousSplitProfile = ""
								SplitParentProfile = result.ProfileID
							} else {
								PreviousSplitProfile = SplitParentProfile
							}
							// BUILD THE GATEWAY XDR
							tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
								SourceAccount:        &pubaccount,
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
								result.TxnHash = UserSplitTxnHashes
								result.Status = "Pending"

								///INSERT INTO TRANSACTION COLLECTION
								err2 := object.InsertTransaction(result)
								if err2 != nil {
									log.Error("Error @InsertTransaction @SubmitSplit " + err2.Error())
								}
							}
							// CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
							txeB64, err := GatewayTXE.Base64()
							if err != nil {
								log.Error("Error @GatewayTXE.Base64 @SubmitSplit " + err.Error())
								result.TxnHash = UserSplitTxnHashes
								result.Status = "Pending"

								///INSERT INTO TRANSACTION COLLECTION
								err2 := object.InsertTransaction(result)
								if err2 != nil {
									log.Error("Error @InsertTransaction @SubmitSplit " + err2.Error())
								}
							}
							// SUBMIT THE GATEWAY'S SIGNED XDR
							display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
							response1 := display1.SubmitXDR("G" + result.TxnType)
							if response1.Error.Code != 200 {
								result.TxnHash = UserSplitTxnHashes
								result.Status = "Pending"
								///INSERT INTO TRANSACTION COLLECTION
								err2 := object.InsertTransaction(result)
								if err2 != nil {
									log.Error("Error @InsertTransaction @SubmitSplit " + err2.Error())
								}
								break
							} else {
								// UPDATE THE TRANSACTION COLLECTION WITH TXN HASH
								result.TxnHash = response1.TXNID
								if result.TxnType == "5" {
									PreviousTxn = response1.TXNID
								}
								///INSERT INTO TRANSACTION COLLECTION
								err1 := object.InsertTransaction(result)
								if err1 != nil {
									log.Error("Error @InsertTransaction @SubmitSplit " + err1.Error())
									break
								} else {
									err := object.RemoveFromTempOrphanList(result.PublicKey, result.SequenceNo)
									if err != nil {
										log.Println("Error while RemoveFromTempOrphanList " + err.Error())
										break
									}
								}
								var PreviousProfile string
								pData1, errorAsync1 := object.GetProfilebyIdentifier(result.FromIdentifier1).Then(func(data interface{}) interface{} {
									return data
								}).Await()
								if pData1 == nil || errorAsync1 != nil {
									log.Error("Error @GetProfilebyIdentifier @SubmitSplit ")
									PreviousProfile = ""
									break
								} else {
									result := pData1.(model.ProfileCollectionBody)
									PreviousProfile = result.ProfileTxn
								}
								Profile := model.ProfileCollectionBody{
									ProfileTxn:         response1.TXNID,
									ProfileID:          result.ProfileID,
									Identifier:         result.Identifier,
									PreviousProfileTxn: PreviousProfile,
									TriggerTxn:         UserSplitTxnHashes,
									TxnType:            result.TxnType,
								}
								err2 := object.InsertProfile(Profile)
								if err2 != nil {
									log.Error("Error @InsertProfile @SubmitSplit " + err2.Error())
								}
								break
							}
						}
					case "7":
						var UserMergeTxnHashes string
						// FOR THE MERGE FIRST BLOCK RETRIEVE THE PREVIOUS TXN FROM GATEWAY DB
						if result.Identifier != result.FromIdentifier1 {
							pData, errorAsync := object.GetLastTransactionbyIdentifierAndTenantId(result.FromIdentifier1, result.TenantID).Then(func(data interface{}) interface{} {
								return data
							}).Await()
							if errorAsync != nil || pData == nil {
								log.Error("Error while GetLastTransactionbyIdentifier @SubmitMerge " + errorAsync.Error())
								// ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
								// DUE TO THE CHILD HAVING A NEW IDENTIFIER
								result.PreviousTxnHash = ""
							} else {
								///ASSIGN PREVIOUS MANAGE DATA BUILDER
								result1 := pData.(model.TransactionCollectionBody)
								result.PreviousTxnHash = result1.TxnHash
								log.Debug(result.PreviousTxnHash)
							}

							pData2, errorAsync2 := object.GetLastTransactionbyIdentifierAndTenantId(result.FromIdentifier2, result.TenantID).Then(func(data interface{}) interface{} {
								return data
							}).Await()

							if errorAsync2 != nil || pData2 == nil {
								log.Error("Error while GetLastTransactionbyIdentifier @SubmitMerge " + errorAsync.Error())
								result.PreviousTxnHash2 = ""
							} else {
								///ASSIGN PREVIOUS MANAGE DATA BUILDER
								result2 := pData2.(model.TransactionCollectionBody)
								result.PreviousTxnHash2 = result2.TxnHash
								log.Debug(result.PreviousTxnHash)
							}
						} else {
							pData3, errorAsync3 := object.GetLastTransactionbyIdentifierAndTenantId(result.FromIdentifier2, result.TenantID).Then(func(data interface{}) interface{} {
								return data
							}).Await()

							if errorAsync3 != nil || pData3 == nil {
								log.Error("Error while GetLastTransactionbyIdentifier @SubmitMerge " + errorAsync3.Error())
								result.PreviousTxnHash2 = ""
							} else {
								///ASSIGN PREVIOUS MANAGE DATA BUILDER
								result3 := pData3.(model.TransactionCollectionBody)
								result.PreviousTxnHash2 = result3.TxnHash
								log.Debug(result.PreviousTxnHash2)
							}
						}
						// SUBMIT THE FIRST XDR SIGNED BY THE USER
						display := stellarExecuter.ConcreteSubmitXDR{XDR: result.XDR}
						result1 := display.SubmitXDR(result.TxnType)
						UserMergeTxnHashes = result1.TXNID

						if result1.Error.Code == 400 {
							log.Error("Index[" + strconv.Itoa(i) + "] TXN: Blockchain Transaction Failed!")
						} else {
							// var PreviousTxn string
							var id apiModel.IdentifierModel
							id.MapValue = result.Identifier
							id.Identifier = result.MapIdentifier
							err3 := object.InsertIdentifier(id)
							if err3 != nil {
								log.Error("identifier map failed" + err3.Error())
							}
							// var PreviousTxn string
							var TypeTXNBuilder txnbuild.ManageData
							var PreviousTXNBuilder txnbuild.ManageData
							var MergeIDBuilder txnbuild.ManageData
							// FIRST MERGE BLOCK
							if result.MergeBlock == 0 {
								TypeTXNBuilder = txnbuild.ManageData{
									Name:  "Type",
									Value: []byte("G8"),
								}
								PreviousTXNBuilder = txnbuild.ManageData{
									Name:  "PreviousTXN",
									Value: []byte(result.PreviousTxnHash),
								}
								MergeIDBuilder = txnbuild.ManageData{
									Name:  "MergeID",
									Value: []byte(result.PreviousTxnHash2),
								}
								result.MergeID = result.PreviousTxnHash2
							} else {
								previousTXNHash := ""
								previousTxn, err := object.GetLastMergeTransactionbyIdentifierAndOrder(result.Identifier, result.TenantID, result.MergeBlock-1).Then(func(data interface{}) interface{} {
									return data
								}).Await()
								if err != nil {
									log.Error("Error while GetLastTransactionbyIdentifier @@SubmitMerge Identifier ", result.Identifier, " mergeBlock: ", result.MergeBlock-1)
								} else if previousTxn == nil {
									log.Error("Can not find GetLastTransactionbyIdentifier @SubmitMerge Identifier ", result.Identifier, " mergeBlock: ", result.MergeBlock-1)
								} else {
									previousTxnData := previousTxn.(model.TransactionCollectionBody)
									previousTXNHash = previousTxnData.TxnHash
								}
								TypeTXNBuilder = txnbuild.ManageData{
									Name:  "Type",
									Value: []byte("G7"),
								}
								PreviousTXNBuilder = txnbuild.ManageData{
									Name:  "PreviousTXN",
									Value: []byte(previousTXNHash),
								}
								MergeIDBuilder = txnbuild.ManageData{
									Name:  "MergeID",
									Value: []byte(result.PreviousTxnHash2),
								}
								result.MergeID = result.PreviousTxnHash2
							}
							if result.MergeBlock == 0 {
								pData, errorAsync := object.GetLastTransactionbyIdentifierAndTenantId(result.FromIdentifier2, result.TenantID).Then(func(data interface{}) interface{} {
									return data
								}).Await()

								if errorAsync != nil || pData == nil {
									log.Error("Error while GetLastTransactionbyIdentifier @SubmitMerge " + errorAsync.Error())
									result.MergeID = ""
								} else {
									///ASSIGN PREVIOUS MANAGE DATA BUILDER
									result4 := pData.(model.TransactionCollectionBody)
									// MergeID = result.TxnHash
									result.MergeID = result4.TxnHash
									log.Error(result.MergeID)
								}
							}
							CurrentTXN := txnbuild.ManageData{
								Name:  "CurrentTXN",
								Value: []byte(UserMergeTxnHashes),
							}
							tx, err := txnbuild.NewTransaction(
								txnbuild.TransactionParams{
									SourceAccount:        &pubaccount,
									IncrementSequenceNum: true,
									Operations:           []txnbuild.Operation{&TypeTXNBuilder, &PreviousTXNBuilder, &MergeIDBuilder, &CurrentTXN},
									BaseFee:              constants.MinBaseFee,
									Preconditions:        txnbuild.Preconditions{TimeBounds: constants.TransactionTimeOut},
								},
							)

							// SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
							GatewayTXE, err := tx.Sign(commons.GetStellarNetwork(), tracifiedAccount)
							if err != nil {
								log.Error("Error while build Transaction @SubmitMerge " + err.Error())
								result.TxnHash = UserMergeTxnHashes
								result.Status = "Pending"
								///INSERT INTO TRANSACTION COLLECTION
								err2 := object.InsertTransaction(result)
								if err2 != nil {
									log.Error("Error while InsertTransaction @SubmitMerge " + err2.Error())
								}
							}
							// CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
							txeB64, err := GatewayTXE.Base64()
							if err != nil {
								log.Error("Error while convert GatewayTXE to base64 @SubmitMerge " + err.Error())
								result.TxnHash = UserMergeTxnHashes
								result.Status = "Pending"

								///INSERT INTO TRANSACTION COLLECTION
								err2 := object.InsertTransaction(result)
								if err2 != nil {
									log.Error("Error while InsertTransaction @SubmitMerge " + err2.Error())
								}
							}

							// SUBMIT THE GATEWAY'S SIGNED XDR
							display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
							response1 := display1.SubmitXDR("G" + result.TxnType)
							if response1.Error.Code != 200 {
								result.TxnHash = UserMergeTxnHashes
								result.Status = "Pending"

								///INSERT INTO TRANSACTION COLLECTION
								err2 := object.InsertTransaction(result)
								if err2 != nil {
									log.Error("Error @InsertTransaction @SubmitSplit " + err2.Error())
								}
								break
							} else {
								// UPDATE THE TRANSACTION COLLECTION WITH TXN HASH
								result.TxnHash = response1.TXNID
								// PreviousTxn = response1.TXNID
								err1 := object.InsertTransaction(result)
								if err1 != nil {
									log.Error("Error @InsertTransaction @SubmitSplit " + err1.Error())
									break
								} else {
									err := object.RemoveFromTempOrphanList(result.PublicKey, result.SequenceNo)
									if err != nil {
										log.Println("Error while RemoveFromTempOrphanList " + err.Error())
										break
									}
								}
								var PreviousProfile string
								pData, errorAsync := object.GetProfilebyIdentifier(result.FromIdentifier1).Then(func(data interface{}) interface{} {
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
									ProfileID:          result.ProfileID,
									Identifier:         result.Identifier,
									PreviousProfileTxn: PreviousProfile,
									TriggerTxn:         UserMergeTxnHashes,
									TxnType:            result.TxnType,
								}
								err3 := object.InsertProfile(Profile)
								if err3 != nil {
									log.Error("Error while InsertProfile @SubmitMerge " + err3.Error())
								}
								break
							}
						}
					case "8":
						var UserMergeTxnHashes string
						// var PreviousTxn string
						// FOR THE MERGE FIRST BLOCK RETRIEVE THE PREVIOUS TXN FROM GATEWAY DB
						if result.Identifier != result.FromIdentifier1 {
							pData, errorAsync := object.GetLastTransactionbyIdentifierAndTenantId(result.FromIdentifier1, result.TenantID).Then(func(data interface{}) interface{} {
								return data
							}).Await()

							if errorAsync != nil || pData == nil {
								log.Error("Error while GetLastTransactionbyIdentifier @SubmitMerge " + errorAsync.Error())
								///ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
								//DUE TO THE CHILD HAVING A NEW IDENTIFIER
								result.PreviousTxnHash = ""
							} else {
								///ASSIGN PREVIOUS MANAGE DATA BUILDER
								result1 := pData.(model.TransactionCollectionBody)
								result.PreviousTxnHash = result1.TxnHash
								log.Debug(result.PreviousTxnHash)
							}

							pData2, errorAsync2 := object.GetLastTransactionbyIdentifierAndTenantId(result.FromIdentifier2, result.TenantID).Then(func(data interface{}) interface{} {
								return data
							}).Await()

							if errorAsync2 != nil || pData2 == nil {
								log.Error("Error while GetLastTransactionbyIdentifier @SubmitMerge " + errorAsync.Error())
								///ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
								//DUE TO THE CHILD HAVING A NEW IDENTIFIER
								result.PreviousTxnHash2 = ""
							} else {
								///ASSIGN PREVIOUS MANAGE DATA BUILDER
								result2 := pData2.(model.TransactionCollectionBody)
								result.PreviousTxnHash2 = result2.TxnHash
								log.Debug(result.PreviousTxnHash)
							}
						} else {
							pData3, errorAsync3 := object.GetLastTransactionbyIdentifierAndTenantId(result.FromIdentifier2, result.TenantID).Then(func(data interface{}) interface{} {
								return data
							}).Await()

							if errorAsync3 != nil || pData3 == nil {
								log.Error("Error while GetLastTransactionbyIdentifier @SubmitMerge " + errorAsync3.Error())
								///ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
								//DUE TO THE CHILD HAVING A NEW IDENTIFIER
								result.PreviousTxnHash2 = ""
							} else {
								///ASSIGN PREVIOUS MANAGE DATA BUILDER
								result3 := pData3.(model.TransactionCollectionBody)
								result.PreviousTxnHash2 = result3.TxnHash
								log.Debug(result.PreviousTxnHash2)
							}
						}

						// SUBMIT THE FIRST XDR SIGNED BY THE USER
						display := stellarExecuter.ConcreteSubmitXDR{XDR: result.XDR}
						result1 := display.SubmitXDR(result.TxnType)
						UserMergeTxnHashes = result1.TXNID

						if result1.Error.Code == 400 {
							log.Error("Index[" + strconv.Itoa(i) + "] TXN: Blockchain Transaction Failed!")
						} else {
							// var PreviousTxn string
							var id apiModel.IdentifierModel
							id.MapValue = result.Identifier
							id.Identifier = result.MapIdentifier
							err3 := object.InsertIdentifier(id)
							if err3 != nil {
								log.Error("identifier map failed" + err3.Error())
							}
							// var PreviousTxn string
							var TypeTXNBuilder txnbuild.ManageData
							var PreviousTXNBuilder txnbuild.ManageData
							var MergeIDBuilder txnbuild.ManageData
							// FIRST MERGE BLOCK
							if result.MergeBlock == 0 {
								// TypeTXNBuilder = build.SetData("Type", []byte("G8"))
								// PreviousTXNBuilder = build.SetData("PreviousTXN", []byte(AP.TxnBody[i].PreviousTxnHash))
								// MergeIDBuilder = build.SetData("MergeID", []byte(AP.TxnBody[i].PreviousTxnHash2))
								TypeTXNBuilder = txnbuild.ManageData{
									Name:  "Type",
									Value: []byte("G8"),
								}
								PreviousTXNBuilder = txnbuild.ManageData{
									Name:  "PreviousTXN",
									Value: []byte(result.PreviousTxnHash),
								}
								MergeIDBuilder = txnbuild.ManageData{
									Name:  "MergeID",
									Value: []byte(result.PreviousTxnHash2),
								}
								result.MergeID = result.PreviousTxnHash2
							} else {
								previousTXNHash := ""
								previousTxn, err := object.GetLastMergeTransactionbyIdentifierAndOrder(result.Identifier, result.TenantID, result.MergeBlock-1).Then(func(data interface{}) interface{} {
									return data
								}).Await()
								if err != nil {
									log.Error("Error while GetLastTransactionbyIdentifier @@SubmitMerge Identifier ", result.Identifier, " mergeBlock: ", result.MergeBlock-1)
								} else if previousTxn == nil {
									log.Error("Can not find GetLastTransactionbyIdentifier @SubmitMerge Identifier ", result.Identifier, " mergeBlock: ", result.MergeBlock-1)
								} else {
									previousTxnData := previousTxn.(model.TransactionCollectionBody)
									previousTXNHash = previousTxnData.TxnHash
								}

								TypeTXNBuilder = txnbuild.ManageData{
									Name:  "Type",
									Value: []byte("G7"),
								}
								PreviousTXNBuilder = txnbuild.ManageData{
									Name:  "PreviousTXN",
									Value: []byte(previousTXNHash),
								}
								MergeIDBuilder = txnbuild.ManageData{
									Name:  "MergeID",
									Value: []byte(result.PreviousTxnHash2),
								}
								result.MergeID = result.PreviousTxnHash2
							}
							if result.MergeBlock == 0 {
								pData, errorAsync := object.GetLastTransactionbyIdentifierAndTenantId(result.FromIdentifier2, result.TenantID).Then(func(data interface{}) interface{} {
									return data
								}).Await()

								if errorAsync != nil || pData == nil {
									log.Error("Error while GetLastTransactionbyIdentifier @SubmitMerge " + errorAsync.Error())
									///ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
									//DUE TO THE CHILD HAVING A NEW IDENTIFIER
									result.MergeID = ""
								} else {
									///ASSIGN PREVIOUS MANAGE DATA BUILDER
									result4 := pData.(model.TransactionCollectionBody)
									// MergeID = result.TxnHash
									result.MergeID = result4.TxnHash
									log.Error(result.MergeID)
								}
							}
							CurrentTXN := txnbuild.ManageData{
								Name:  "CurrentTXN",
								Value: []byte(UserMergeTxnHashes),
							}
							tx, err := txnbuild.NewTransaction(
								txnbuild.TransactionParams{
									SourceAccount:        &pubaccount,
									IncrementSequenceNum: true,
									Operations:           []txnbuild.Operation{&TypeTXNBuilder, &PreviousTXNBuilder, &MergeIDBuilder, &CurrentTXN},
									BaseFee:              constants.MinBaseFee,
									Preconditions:        txnbuild.Preconditions{TimeBounds: constants.TransactionTimeOut},
								},
							)
							// SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
							GatewayTXE, err := tx.Sign(commons.GetStellarNetwork(), tracifiedAccount)
							if err != nil {
								log.Error("Error while build Transaction @SubmitMerge " + err.Error())
								result.TxnHash = UserMergeTxnHashes
								result.Status = "Pending"

								///INSERT INTO TRANSACTION COLLECTION
								err2 := object.InsertTransaction(result)
								if err2 != nil {
									log.Error("Error while InsertTransaction @SubmitMerge " + err2.Error())
								}
							}
							// CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
							txeB64, err := GatewayTXE.Base64()
							if err != nil {
								log.Error("Error while convert GatewayTXE to base64 @SubmitMerge " + err.Error())
								result.TxnHash = UserMergeTxnHashes
								result.Status = "Pending"

								///INSERT INTO TRANSACTION COLLECTION
								err2 := object.InsertTransaction(result)
								if err2 != nil {
									log.Error("Error while InsertTransaction @SubmitMerge " + err2.Error())
								}
							}
							// SUBMIT THE GATEWAY'S SIGNED XDR
							display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
							response1 := display1.SubmitXDR("G" + result.TxnType)
							if response1.Error.Code != 200 {
								result.TxnHash = UserMergeTxnHashes
								result.Status = "Pending"
								///INSERT INTO TRANSACTION COLLECTION
								err2 := object.InsertTransaction(result)
								if err2 != nil {
									log.Error("Error @InsertTransaction @SubmitSplit " + err2.Error())
								}
								break
							} else {
								// UPDATE THE TRANSACTION COLLECTION WITH TXN HASH
								result.TxnHash = response1.TXNID
								// PreviousTxn = response1.TXNID
								err1 := object.InsertTransaction(result)
								if err1 != nil {
									log.Error("Error @InsertTransaction @SubmitSplit " + err1.Error())
									break
								} else {
									err := object.RemoveFromTempOrphanList(result.PublicKey, result.SequenceNo)
									if err != nil {
										log.Println("Error while RemoveFromTempOrphanList " + err.Error())
										break
									}
								}
								var PreviousProfile string
								pData, errorAsync := object.GetProfilebyIdentifier(result.FromIdentifier1).Then(func(data interface{}) interface{} {
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
									ProfileID:          result.ProfileID,
									Identifier:         result.Identifier,
									PreviousProfileTxn: PreviousProfile,
									TriggerTxn:         UserMergeTxnHashes,
									TxnType:            result.TxnType,
								}
								err3 := object.InsertProfile(Profile)
								if err3 != nil {
									log.Error("Error while InsertProfile @SubmitMerge " + err3.Error())
								}
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
