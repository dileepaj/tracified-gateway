package services

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/configs"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
	"github.com/dileepaj/tracified-gateway/proofs/retriever/stellarRetriever"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

func SubmitUserDataToStellar(deliver amqp091.Delivery) {
	// object := dao.Connection{}
	log.Printf("@SubmitUserDataToStellar Reciever a message to  : %s", deliver.Body)
	// Convert the JSON string to the txnBody struct
	var txnBody model.TransactionCollectionBody
	err := json.Unmarshal(deliver.Body, &txnBody)
	if err != nil {
		logrus.Error("Error Unmarshal @SubmitData " + err.Error())
		deliver.Nack(false, true)
		return
	}
	var txe xdr.TransactionEnvelope
	// decode the XDR
	err1 := xdr.SafeUnmarshalBase64(txnBody.XDR, &txe)
	if err1 != nil {
		logrus.Error("Error SafeUnmarshalBase64 @SubmitData " + err.Error())
		deliver.Nack(false, true)
		return
	}
	utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST, configs.BenchmarkLogsAction.CONSUMING_TRANSACTION_QUEUE, txnBody.RequestId, configs.BenchmarkLogsStatus.OK)
	txnBody.PublicKey = txe.SourceAccount().ToAccountId().Address()
	txnBody.SequenceNo = int64(txe.SeqNum())
	stellarRetriever.MapXDROperations(&txnBody, txe.Operations())
	txnBody.Status = "pending"
	if txnBody.TxnType == "5" {
		txnBody.Identifier = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations()[1].Body.ManageDataOp.DataValue), "&")
	}

	var datax = model.TransactionData{
		FOUser:        txnBody.AppAccount,
		XDR:           txnBody.XDR,
		AccountIssuer: txnBody.PublicKey,
	}
	responsexdr, errx := SubmitFOData(datax)
	if errx != nil {
		logrus.Error("response.Error.Code 400 for SubmitXDR", errx)
		return
	}

	display := stellarExecuter.ConcreteSubmitXDR{XDR: responsexdr}
	utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST, configs.BenchmarkLogsAction.FO_USER_XDR_SUBMITTING_TO_BLOCKCHAIN, txnBody.RequestId, configs.BenchmarkLogsStatus.SENDING)
	response := display.SubmitXDR(txnBody.TxnType)
	if response.Error.Code != 200 || response.Error.Message != "" {
		utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST, configs.BenchmarkLogsAction.FO_USER_XDR_SUBMITTING_TO_BLOCKCHAIN, txnBody.RequestId, configs.BenchmarkLogsStatus.ERROR)
		logrus.Error("Failed to submit the XDR ", " Error: ", response.Error.Message, " Timestamp: ", txnBody.Timestamp, " XDR: ",
			txnBody.XDR, "TXNType: ", txnBody.TxnType, " Identifier: ", txnBody.MapIdentifier, " Sequence No: ", txnBody.SequenceNo, " PublicKey: ", txnBody.PublicKey)
		deliver.Nack(false, true)
		return
	}

	txnBody.FOUserTXNHash = response.TXNID
	utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST, configs.BenchmarkLogsAction.FO_USER_XDR_SUBMITTING_TO_BLOCKCHAIN, txnBody.RequestId, configs.BenchmarkLogsStatus.SUCCESS)
	logrus.Info("Stellar FO user created TXN hash: ", response.TXNID, " Timestamp: ", txnBody.Timestamp, "TXNType: ", txnBody.TxnType, " Identifier: ",
		txnBody.MapIdentifier, " Sequence No: ", txnBody.SequenceNo, " PublicKey: ", txnBody.PublicKey)
	deliver.Ack(false)
	jsonStr, err := json.Marshal(txnBody)
	if err != nil {
		logrus.Error("Error in convert the struct to a JSON string using encoding/json:", err)
		return
	}
	PublishToQueue(configs.QueueBackLinks.Name, string(jsonStr), configs.QueueBackLinks.Method)
	utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST, configs.BenchmarkLogsAction.PUBLISH_TO_BACKLINK, txnBody.RequestId, configs.BenchmarkLogsStatus.OK)
	return
}

func SubmitBacklinksDataToStellar(deliver amqp091.Delivery) {
	object := dao.Connection{}
	log.Printf("@SubmitbacklinksDataToStellar Reciver a message to  : %s", deliver.Body)
	// Convert the JSON string to the txnBody struct
	var result model.TransactionCollectionBody
	err := json.Unmarshal(deliver.Body, &result)
	if err != nil {
		logrus.Error("Error Unmarshal @SubmitData " + err.Error())
		return
	}
	utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST, configs.BenchmarkLogsAction.CONSUMING_BACKLINK_QUEUE, result.RequestId, configs.BenchmarkLogsStatus.OK)
	var txe xdr.TransactionEnvelope
	// decode the XDR
	err1 := xdr.SafeUnmarshalBase64(result.XDR, &txe)
	if err1 != nil {
		logrus.Error("Error SafeUnmarshalBase64 @SubmitData " + err.Error())
		return
	}
	publicKey := constants.PublicKey
	secretKey := constants.SecretKey
	tracifiedAccount, err := keypair.ParseFull(secretKey)
	if err != nil {
		logrus.Error(err)
	}
	client := commons.GetHorizonClient()
	pubaccountRequest := horizonclient.AccountRequest{AccountID: publicKey}
	pubaccount, err := client.AccountDetail(pubaccountRequest)
	switch result.TxnType {
	case "0":
		// type 0 = Genesis
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
			Value: []byte(result.FOUserTXNHash),
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
			logrus.Println("Error while buliding XDR " + err.Error())
			break
		}
		// SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
		GatewayTXE, err := tx.Sign(commons.GetStellarNetwork(), tracifiedAccount)
		if err != nil {
			logrus.Println("Error while getting GatewayTXE by secretKey " + err.Error())
			break
		}
		// CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
		txeB64, err := GatewayTXE.Base64()
		if err != nil {
			logrus.Println("Error while converting GatewayTXE to base64 " + err.Error())
			break
		}
		var id apiModel.IdentifierModel
		id.MapValue = result.Identifier
		id.Identifier = result.MapIdentifier
		err3 := object.InsertIdentifier(id)
		if err3 != nil {
			logrus.Error("identifier map failed" + err3.Error())
		}
		// SUBMIT THE GATEWAY'S SIGNED XDR
		display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
		utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST_GENESIS, configs.BenchmarkLogsAction.BACKLINK_XDR_SUBMITTING_TO_BLOCKCHAIN, result.RequestId, configs.BenchmarkLogsStatus.SENDING)
		response1 := display1.SubmitXDR("G" + result.TxnType)
		if response1.Error.Code != 200 || response1.Error.Message != "" {
			utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST_GENESIS, configs.BenchmarkLogsAction.BACKLINK_XDR_SUBMITTING_TO_BLOCKCHAIN, result.RequestId, configs.BenchmarkLogsStatus.ERROR)
			logrus.Error("Failed to submit backlinks the XDR ", " Error: ", response1.Error.Message, " Timestamp: ", result.Timestamp, " XDR: ",
				result.XDR, "TXNType: ", result.TxnType, " Identifier: ", result.MapIdentifier, " Sequence No: ", result.SequenceNo, " PublicKey: ", result.PublicKey)
			deliver.Nack(false, true)
			break
		}
		utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST_GENESIS, configs.BenchmarkLogsAction.BACKLINK_XDR_SUBMITTING_TO_BLOCKCHAIN, result.RequestId, configs.BenchmarkLogsStatus.SUCCESS)
		result.TxnHash = response1.TXNID
		result.Status = "done"
		break
	case "2":
		// type 2 = Normal TDP
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
		CurrentTXNBuilder := txnbuild.ManageData{
			Name:  "CurrentTXN",
			Value: []byte(result.FOUserTXNHash),
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
			logrus.Println("Error while buliding XDR " + err.Error())
			break
		}
		// SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
		GatewayTXE, err := tx.Sign(commons.GetStellarNetwork(), tracifiedAccount)
		if err != nil {
			logrus.Println("Error while getting GatewayTXE by secretKey " + err.Error())
			break
		}
		// CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
		txeB64, err := GatewayTXE.Base64()
		if err != nil {
			logrus.Println("Error while converting GatewayTXE to base64 " + err.Error())
			break
		}
		// SUBMIT THE GATEWAY'S SIGNED XDR
		display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
		utilities.BenchmarkLog("tdp-request-normal", configs.BenchmarkLogsAction.BACKLINK_XDR_SUBMITTING_TO_BLOCKCHAIN, result.RequestId, configs.BenchmarkLogsStatus.SENDING)
		response1 := display1.SubmitXDR("G" + result.TxnType)
		if response1.Error.Code != 200 || response1.Error.Message != "" {
			utilities.BenchmarkLog("tdp-request-normal", configs.BenchmarkLogsAction.BACKLINK_XDR_SUBMITTING_TO_BLOCKCHAIN, result.RequestId, configs.BenchmarkLogsStatus.ERROR)
			logrus.Error("Failed to submit backlinks the XDR ", " Error: ", response1.Error.Message, " Timestamp: ", result.Timestamp, " XDR: ",
				result.XDR, "TXNType: ", result.TxnType, " Identifier: ", result.MapIdentifier, " Sequence No: ", result.SequenceNo, " PublicKey: ", result.PublicKey)
			deliver.Nack(false, true)
			break
		}
		utilities.BenchmarkLog("tdp-request-normal", configs.BenchmarkLogsAction.BACKLINK_XDR_SUBMITTING_TO_BLOCKCHAIN, result.RequestId, configs.BenchmarkLogsStatus.SUCCESS)
		result.TxnHash = response1.TXNID
		result.Status = "done"
		break
	case "9":
		// type 9 = Transfer
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
			Value: []byte(result.FOUserTXNHash),
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
			logrus.Error("Error while buliding XDR " + err.Error())
			break
		}
		// SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
		GatewayTXE, err := tx.Sign(commons.GetStellarNetwork(), tracifiedAccount)
		if err != nil {
			logrus.Error("Error while getting GatewayTXE " + err.Error())
			break
		}
		// CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
		txeB64, err := GatewayTXE.Base64()
		if err != nil {
			logrus.Error("Error while converting to base64 " + err.Error())
			break
		}
		// SUBMIT THE GATEWAY'S SIGNED XDR
		display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
		utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST_TRANSFER, configs.BenchmarkLogsAction.BACKLINK_XDR_SUBMITTING_TO_BLOCKCHAIN, result.RequestId, configs.BenchmarkLogsStatus.SENDING)
		response1 := display1.SubmitXDR("G" + result.TxnType)
		if response1.Error.Code != 200 || response1.Error.Message != "" {
			utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST_TRANSFER, configs.BenchmarkLogsAction.BACKLINK_XDR_SUBMITTING_TO_BLOCKCHAIN, result.RequestId, configs.BenchmarkLogsStatus.ERROR)
			logrus.Error("Failed to submit backlinks the XDR ", " Error: ", response1.Error.Message, " Timestamp: ", result.Timestamp, " XDR: ",
				result.XDR, "TXNType: ", result.TxnType, " Identifier: ", result.MapIdentifier, " Sequence No: ", result.SequenceNo, " PublicKey: ", result.PublicKey)
			deliver.Nack(false, true)
			break
		}
		utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST_TRANSFER, configs.BenchmarkLogsAction.BACKLINK_XDR_SUBMITTING_TO_BLOCKCHAIN, result.RequestId, configs.BenchmarkLogsStatus.SUCCESS)
		result.TxnHash = response1.TXNID
		result.Status = "done"
		break
	case "5":
		// type 5 = Split Parent
		var UserSplitTxnHashes string
		var PreviousTxn string
		// ParentIdentifier = Identifier
		pData, errAsnc := object.GetLastTransactionbyIdentifierAndTenantId(result.Identifier, result.TenantID).Then(func(data interface{}) interface{} {
			return data
		}).Await()
		if pData == nil || errAsnc != nil {
			logrus.Error("Error @GetLastTransactionbyIdentifier @SubmitSplit ")
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
		UserSplitTxnHashes = result.FOUserTXNHash
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
			logrus.Error("Error @tx.Sign @SubmitSplit " + err.Error())
			result.TxnHash = UserSplitTxnHashes
			result.Status = "Pending"
			break
		}
		// CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
		txeB64, err := GatewayTXE.Base64()
		if err != nil {
			logrus.Error("Error @GatewayTXE.Base64 @SubmitSplit " + err.Error())
			result.TxnHash = UserSplitTxnHashes
			result.Status = "Pending"
			break
		}
		// SUBMIT THE GATEWAY'S SIGNED XDR
		display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
		utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST_SPLIT_PARENT, configs.BenchmarkLogsAction.BACKLINK_XDR_SUBMITTING_TO_BLOCKCHAIN, result.RequestId, configs.BenchmarkLogsStatus.SENDING)
		response1 := display1.SubmitXDR("G" + result.TxnType)
		if response1.Error.Code != 200 || response1.Error.Message != "" {
			utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST_SPLIT_PARENT, configs.BenchmarkLogsAction.BACKLINK_XDR_SUBMITTING_TO_BLOCKCHAIN, result.RequestId, configs.BenchmarkLogsStatus.ERROR)
			logrus.Error("Failed to submit backlinks the XDR ", " Error: ", response1.Error.Message, " Timestamp: ", result.Timestamp, " XDR: ",
				result.XDR, "TXNType: ", result.TxnType, " Identifier: ", result.MapIdentifier, " Sequence No: ", result.SequenceNo, " PublicKey: ", result.PublicKey)
			deliver.Nack(false, true)
			break
		} else {
			utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST_SPLIT_PARENT, configs.BenchmarkLogsAction.BACKLINK_XDR_SUBMITTING_TO_BLOCKCHAIN, result.RequestId, configs.BenchmarkLogsStatus.SUCCESS)
			// UPDATE THE TRANSACTION COLLECTION WITH TXN HASH
			result.TxnHash = response1.TXNID
			if result.TxnType == "5" {
				PreviousTxn = response1.TXNID
			}
			break
		}
	case "6":
		// type 0 = Split Child
		var UserSplitTxnHashes string
		var PreviousTxn string
		var SplitParentProfile string
		var PreviousSplitProfile string
		/*
			When constructing a backlink transaction(put from gateway) for a split, it is important to exclude the split-parent transaction as its previous transaction.
			Instead, you should obtain the most recent transaction that is specific to the identifier and disregard the split-parent transaction.
		*/
		backlinkData, err := object.GetLastTransactionbyIdentifierNotSplitParent(result.FromIdentifier1, result.TenantID).Then(func(data interface{}) interface{} {
			return data
		}).Await()
		if backlinkData == nil || err != nil {
			logrus.Info("Can not find transaction form database ", "build Split")
		} else {
			result := backlinkData.(model.TransactionCollectionBody)
			PreviousTxn = result.TxnHash
			result.PreviousTxnHash = result.TxnHash
		}
		UserSplitTxnHashes = result.FOUserTXNHash
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
			logrus.Error("Error @tx.Sign @SubmitSplit " + err.Error())
			result.TxnHash = UserSplitTxnHashes
			result.Status = "Pending"
			break
		}
		// CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
		txeB64, err := GatewayTXE.Base64()
		if err != nil {
			logrus.Error("Error @GatewayTXE.Base64 @SubmitSplit " + err.Error())
			result.TxnHash = UserSplitTxnHashes
			result.Status = "Pending"
			break
		}
		// SUBMIT THE GATEWAY'S SIGNED XDR
		display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
		utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST_SPLIT_CHILD, configs.BenchmarkLogsAction.BACKLINK_XDR_SUBMITTING_TO_BLOCKCHAIN, result.RequestId, configs.BenchmarkLogsStatus.SENDING)
		response1 := display1.SubmitXDR("G" + result.TxnType)
		if response1.Error.Code != 200 || response1.Error.Message != "" {
			utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST_SPLIT_CHILD, configs.BenchmarkLogsAction.BACKLINK_XDR_SUBMITTING_TO_BLOCKCHAIN, result.RequestId, configs.BenchmarkLogsStatus.ERROR)
			logrus.Error("Failed to submit backlinks the XDR ", " Error: ", response1.Error.Message, " Timestamp: ", result.Timestamp, " XDR: ",
				result.XDR, "TXNType: ", result.TxnType, " Identifier: ", result.MapIdentifier, " Sequence No: ", result.SequenceNo, " PublicKey: ", result.PublicKey)
			deliver.Nack(false, true)
			break
		} else {
			utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST_SPLIT_CHILD, configs.BenchmarkLogsAction.BACKLINK_XDR_SUBMITTING_TO_BLOCKCHAIN, result.RequestId, configs.BenchmarkLogsStatus.SUCCESS)
			// UPDATE THE TRANSACTION COLLECTION WITH TXN HASH
			result.TxnHash = response1.TXNID
			if result.TxnType == "5" {
				PreviousTxn = response1.TXNID
			}
			var PreviousProfile string
			pData1, errorAsync1 := object.GetProfilebyIdentifier(result.FromIdentifier1).Then(func(data interface{}) interface{} {
				return data
			}).Await()
			if pData1 == nil || errorAsync1 != nil {
				logrus.Error("Error @GetProfilebyIdentifier @SubmitSplit ")
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
				logrus.Error("Error @InsertProfile @SubmitSplit " + err2.Error())
			}
			break
		}
	case "7":
		// type 7 = Merge
		var UserMergeTxnHashes string
		// FOR THE MERGE FIRST BLOCK RETRIEVE THE PREVIOUS TXN FROM GATEWAY DB
		if result.Identifier != result.FromIdentifier1 {
			pData, errorAsync := object.GetLastTransactionbyIdentifierAndTenantId(result.FromIdentifier1, result.TenantID).Then(func(data interface{}) interface{} {
				return data
			}).Await()
			if errorAsync != nil || pData == nil {
				logrus.Error("Error while GetLastTransactionbyIdentifier @SubmitMerge ")
				// ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
				// DUE TO THE CHILD HAVING A NEW IDENTIFIER
				result.PreviousTxnHash = ""
			} else {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER
				result1 := pData.(model.TransactionCollectionBody)
				result.PreviousTxnHash = result1.TxnHash
				logrus.Debug(result.PreviousTxnHash)
			}

			pData2, errorAsync2 := object.GetLastTransactionbyIdentifierAndTenantId(result.FromIdentifier2, result.TenantID).Then(func(data interface{}) interface{} {
				return data
			}).Await()

			if errorAsync2 != nil || pData2 == nil {
				logrus.Error("Error while GetLastTransactionbyIdentifier @SubmitMerge " + errorAsync.Error())
				result.PreviousTxnHash2 = ""
			} else {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER
				result2 := pData2.(model.TransactionCollectionBody)
				result.PreviousTxnHash2 = result2.TxnHash
				logrus.Debug(result.PreviousTxnHash)
			}
		} else {
			pData3, errorAsync3 := object.GetLastTransactionbyIdentifierAndTenantId(result.FromIdentifier2, result.TenantID).Then(func(data interface{}) interface{} {
				return data
			}).Await()

			if errorAsync3 != nil || pData3 == nil {
				logrus.Error("Error while GetLastTransactionbyIdentifier @SubmitMerge ")
				result.PreviousTxnHash2 = ""
			} else {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER
				result3 := pData3.(model.TransactionCollectionBody)
				result.PreviousTxnHash2 = result3.TxnHash
				logrus.Debug(result.PreviousTxnHash2)
			}
		}
		// SUBMIT THE FIRST XDR SIGNED BY THE USER
		UserMergeTxnHashes = result.FOUserTXNHash
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
				logrus.Error("Error while GetLastTransactionbyIdentifier @@SubmitMerge Identifier ", result.Identifier, " mergeBlock: ", result.MergeBlock-1)
			} else if previousTxn == nil {
				logrus.Error("Can not find GetLastTransactionbyIdentifier @SubmitMerge Identifier ", result.Identifier, " mergeBlock: ", result.MergeBlock-1)
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
				logrus.Error("Error while GetLastTransactionbyIdentifier @SubmitMerge " + errorAsync.Error())
				result.MergeID = ""
			} else {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER
				result4 := pData.(model.TransactionCollectionBody)
				// MergeID = result.TxnHash
				result.MergeID = result4.TxnHash
				logrus.Error(result.MergeID)
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
			logrus.Error("Error while build Transaction @SubmitMerge " + err.Error())
			result.TxnHash = UserMergeTxnHashes
			result.Status = "Pending"
			break
		}
		// CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
		txeB64, err := GatewayTXE.Base64()
		if err != nil {
			logrus.Error("Error while convert GatewayTXE to base64 @SubmitMerge " + err.Error())
			result.TxnHash = UserMergeTxnHashes
			result.Status = "Pending"
			break
		}

		// SUBMIT THE GATEWAY'S SIGNED XDR
		display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
		utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST_MERGE_7, configs.BenchmarkLogsAction.BACKLINK_XDR_SUBMITTING_TO_BLOCKCHAIN, result.RequestId, configs.BenchmarkLogsStatus.SENDING)
		response1 := display1.SubmitXDR("G" + result.TxnType)
		if response1.Error.Code != 200 || response1.Error.Message != "" {
			utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST_MERGE_7, configs.BenchmarkLogsAction.BACKLINK_XDR_SUBMITTING_TO_BLOCKCHAIN, result.RequestId, configs.BenchmarkLogsStatus.ERROR)
			logrus.Error("Failed to submit backlinks the XDR ", " Error: ", response1.Error.Message, " Timestamp: ", result.Timestamp, " XDR: ",
				result.XDR, "TXNType: ", result.TxnType, " Identifier: ", result.MapIdentifier, " Sequence No: ", result.SequenceNo, " PublicKey: ", result.PublicKey)
			deliver.Nack(false, true)
			break
		} else {
			utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST_MERGE_7, configs.BenchmarkLogsAction.BACKLINK_XDR_SUBMITTING_TO_BLOCKCHAIN, result.RequestId, configs.BenchmarkLogsStatus.SUCCESS)
			// UPDATE THE TRANSACTION COLLECTION WITH TXN HASH
			result.TxnHash = response1.TXNID
			var PreviousProfile string
			pData, errorAsync := object.GetProfilebyIdentifier(result.FromIdentifier1).Then(func(data interface{}) interface{} {
				return data
			}).Await()

			if errorAsync != nil || pData == nil {
				logrus.Error("Error while GetProfilebyIdentifier @SubmitMerge" + errorAsync.Error())
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
				logrus.Error("Error while InsertProfile @SubmitMerge " + err3.Error())
			}
			break
		}
	case "8":
		// type 8 = Merge
		var UserMergeTxnHashes string
		// var PreviousTxn string
		// FOR THE MERGE FIRST BLOCK RETRIEVE THE PREVIOUS TXN FROM GATEWAY DB
		if result.Identifier != result.FromIdentifier1 {
			pData, errorAsync := object.GetLastTransactionbyIdentifierAndTenantId(result.FromIdentifier1, result.TenantID).Then(func(data interface{}) interface{} {
				return data
			}).Await()

			if errorAsync != nil || pData == nil {
				logrus.Error("Error while GetLastTransactionbyIdentifier @SubmitMerge ")
				///ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
				//DUE TO THE CHILD HAVING A NEW IDENTIFIER
				result.PreviousTxnHash = ""
			} else {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER
				result1 := pData.(model.TransactionCollectionBody)
				result.PreviousTxnHash = result1.TxnHash
				logrus.Debug(result.PreviousTxnHash)
			}

			pData2, errorAsync2 := object.GetLastTransactionbyIdentifierAndTenantId(result.FromIdentifier2, result.TenantID).Then(func(data interface{}) interface{} {
				return data
			}).Await()

			if errorAsync2 != nil || pData2 == nil {
				logrus.Error("Error while GetLastTransactionbyIdentifier @SubmitMerge ")
				///ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
				//DUE TO THE CHILD HAVING A NEW IDENTIFIER
				result.PreviousTxnHash2 = ""
			} else {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER
				result2 := pData2.(model.TransactionCollectionBody)
				result.PreviousTxnHash2 = result2.TxnHash
				logrus.Debug(result.PreviousTxnHash)
			}
		} else {
			pData3, errorAsync3 := object.GetLastTransactionbyIdentifierAndTenantId(result.FromIdentifier2, result.TenantID).Then(func(data interface{}) interface{} {
				return data
			}).Await()

			if errorAsync3 != nil || pData3 == nil {
				logrus.Error("Error while GetLastTransactionbyIdentifier @SubmitMerge ")
				///ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
				//DUE TO THE CHILD HAVING A NEW IDENTIFIER
				result.PreviousTxnHash2 = ""
			} else {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER
				result3 := pData3.(model.TransactionCollectionBody)
				result.PreviousTxnHash2 = result3.TxnHash
				logrus.Debug(result.PreviousTxnHash2)
			}
		}
		UserMergeTxnHashes = result.FOUserTXNHash
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
				logrus.Error("Error while GetLastTransactionbyIdentifier @@SubmitMerge Identifier ", result.Identifier, " mergeBlock: ", result.MergeBlock-1)
			} else if previousTxn == nil {
				logrus.Error("Can not find GetLastTransactionbyIdentifier @SubmitMerge Identifier ", result.Identifier, " mergeBlock: ", result.MergeBlock-1)
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
				logrus.Error("Error while GetLastTransactionbyIdentifier @SubmitMerge ")
				///ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
				//DUE TO THE CHILD HAVING A NEW IDENTIFIER
				result.MergeID = ""
			} else {
				///ASSIGN PREVIOUS MANAGE DATA BUILDER
				result4 := pData.(model.TransactionCollectionBody)
				// MergeID = result.TxnHash
				result.MergeID = result4.TxnHash
				logrus.Error(result.MergeID)
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
			logrus.Error("Error while build Transaction @SubmitMerge " + err.Error())
			result.TxnHash = UserMergeTxnHashes
			result.Status = "Pending"
			break
		}
		// CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
		txeB64, err := GatewayTXE.Base64()
		if err != nil {
			logrus.Error("Error while convert GatewayTXE to base64 @SubmitMerge " + err.Error())
			result.TxnHash = UserMergeTxnHashes
			result.Status = "Pending"
			break
		}
		// SUBMIT THE GATEWAY'S SIGNED XDR
		display1 := stellarExecuter.ConcreteSubmitXDR{XDR: txeB64}
		utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST_MERGE_8, configs.BenchmarkLogsAction.BACKLINK_XDR_SUBMITTING_TO_BLOCKCHAIN, result.RequestId, configs.BenchmarkLogsStatus.SENDING)
		response1 := display1.SubmitXDR("G" + result.TxnType)
		if response1.Error.Code != 200 || response1.Error.Message != "" {
			utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST_MERGE_8, configs.BenchmarkLogsAction.BACKLINK_XDR_SUBMITTING_TO_BLOCKCHAIN, result.RequestId, configs.BenchmarkLogsStatus.ERROR)
			logrus.Error("Failed to submit backlinks the XDR ", " Error: ", response1.Error.Message, " Timestamp: ", result.Timestamp, " XDR: ",
				result.XDR, "TXNType: ", result.TxnType, " Identifier: ", result.MapIdentifier, " Sequence No: ", result.SequenceNo, " PublicKey: ", result.PublicKey)
			deliver.Nack(false, true)
			break
		} else {
			utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST_MERGE_8, configs.BenchmarkLogsAction.BACKLINK_XDR_SUBMITTING_TO_BLOCKCHAIN, result.RequestId, configs.BenchmarkLogsStatus.SUCCESS)
			// UPDATE THE TRANSACTION COLLECTION WITH TXN HASH
			result.TxnHash = response1.TXNID
			var PreviousProfile string
			pData, errorAsync := object.GetProfilebyIdentifier(result.FromIdentifier1).Then(func(data interface{}) interface{} {
				return data
			}).Await()
			if errorAsync != nil || pData == nil {
				logrus.Error("Error while GetProfilebyIdentifier @SubmitMerge")
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
				logrus.Error("Error while InsertProfile @SubmitMerge " + err3.Error())
			}
			break
		}
	}
	errInsertTransaction := object.InsertTransaction(result)
	if errInsertTransaction != nil {
		logrus.Error("Error while @InsertTransaction " + err1.Error())
		deliver.Nack(false, true)
		return
	}
	utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST, configs.BenchmarkLogsAction.TDP_REQUEST_SUBMITTED, result.RequestId, configs.BenchmarkLogsStatus.COMPLETE)
	logrus.Info("Stellar back-link TXN hash: ", result.TxnHash, " Timestamp: ", result.Timestamp, "TXNType: ", result.TxnType,
		" Identifier: ", result.MapIdentifier, " Sequence No: ", result.SequenceNo, " PublicKey: ", result.PublicKey)
	deliver.Ack(false)
	return
}
