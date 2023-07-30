package services

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
	"github.com/dileepaj/tracified-gateway/proofs/retriever/stellarRetriever"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/xdr"
)

func SubmitData(deliver amqp091.Delivery) {
	// object := dao.Connection{}
	log.Printf("Reciver a message: %s", deliver.Body)
	// Convert the JSON string to the txnBody struct
	var txnBody model.TransactionCollectionBody
	err := json.Unmarshal(deliver.Body, &txnBody)
	if err != nil {
		logrus.Error("Error Unmarshal @SubmitData " + err.Error())
		return
	}
	var txe xdr.TransactionEnvelope
	// decode the XDR
	err1 := xdr.SafeUnmarshalBase64(txnBody.XDR, &txe)
	if err1 != nil {
		logrus.Error("Error SafeUnmarshalBase64 @SubmitData " + err.Error())
		return
	}
	txnBody.PublicKey = txe.SourceAccount().ToAccountId().Address()
	txnBody.SequenceNo = int64(txe.SeqNum())
	stellarRetriever.MapXDROperations(&txnBody, txe.Operations())
	txnBody.Status = "pending"
	if txnBody.TxnType == "5" {
		txnBody.Identifier = strings.TrimLeft(fmt.Sprintf("%s", txe.Operations()[1].Body.ManageDataOp.DataValue), "&")
	}
	logrus.Debug("Identifier: ", txnBody.Identifier, " XDR: ", txnBody.XDR)

	switch txnBody.TxnType {
	case "0":
		display := stellarExecuter.ConcreteSubmitXDR{XDR: txnBody.XDR}
		response := display.SubmitXDR(txnBody.TxnType)
		fmt.Println("------------------------------------------", response)
		if response.Error.Code != 200 || response.Error.Message != "" {
			logrus.Error("Failed to submit the XDR ", " Error: ", response.Error.Message, " Timestamp: ", txnBody.Timestamp, " XDR: ",
				txnBody.XDR, "TXNType: ", txnBody.TxnType, " Identifier: ", txnBody.MapIdentifier, " Sequence No: ", txnBody.SequenceNo, " PublicKey: ", txnBody.PublicKey)
			deliver.Ack(true)
			break
		}
		userCreatedXDRTxnHash := response.TXNID
		fmt.Println("-------------------------------------------", userCreatedXDRTxnHash)
		deliver.Ack(false)
		break
	case "2":
		// data, errorLastTXN := object.GetLastTransactionbyIdentifierAndTenantId(txnBody.Identifier, txnBody.TenantID).Then(func(data interface{}) interface{} {
		// 	return data
		// }).Await()

		display := stellarExecuter.ConcreteSubmitXDR{XDR: txnBody.XDR}
		response := display.SubmitXDR(txnBody.TxnType)
		fmt.Println("------------------------------------------", response)
		userCreatedXDRTxnHash := response.TXNID
		fmt.Println("-------------------------------------------", userCreatedXDRTxnHash)
		if response.Error.Code != 200 || response.Error.Message != "" {
			logrus.Error("Failed to submit the XDR ", " Error: ", response.Error.Message, " Timestamp: ", txnBody.Timestamp, " XDR: ",
				txnBody.XDR, "TXNType: ", txnBody.TxnType, " Identifier: ", txnBody.MapIdentifier, " Sequence No: ", txnBody.SequenceNo, " PublicKey: ", txnBody.PublicKey)
			deliver.Ack(true)
			break
		}
		deliver.Ack(false)
		break
	}
	return
}
