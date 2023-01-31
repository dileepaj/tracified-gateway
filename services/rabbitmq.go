package services

import (
	"bytes"
	"encoding/json"
	"fmt"

	"time"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	deploy "github.com/dileepaj/tracified-gateway/protocols/ethereum/deploy"
	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

var (
	USER     = commons.GoDotEnvVariable("RABBITUSER")
	PASSWORD = commons.GoDotEnvVariable("RABBITPASSWORD")
	HOSTNAME = commons.GoDotEnvVariable("RABBITMQ_SERVICE_HOST")
	PORT     = commons.GoDotEnvVariable("RABBITPORT")
)

func ReceiverRmq() error {
	previousTxnHash := ""
	rabbitConnection := `amqp://` + USER + `:` + PASSWORD + `@` + HOSTNAME + `:` + PORT + `/`
	conn, err := amqp.Dial(rabbitConnection)
	if err != nil {
		logrus.Error("Failed to connect to RabbitMQ ", err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logrus.Error("%s: %s", "Failed to open a channel", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"socialimpact", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		logrus.Error("%s: %s", "Failed to declare a queue", err)
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		logrus.Error("%s: %s", "Failed to register a consumer", err)
	}
	var forever chan struct{}

	go func() {
		for d := range msgs {

			object := dao.Connection{}
			var queue model.SendToQueue
			var manageDataOprations []txnbuild.Operation
			if err := json.Unmarshal(d.Body, &queue); err != nil {
				logrus.Error("Unmarshal in rabbitmq reciverRmq ", err.Error())
			}
			if queue.TransactionCount == 0 {
				previousTxnHash = ""
			}
			for i := range queue.Operations {
				manageDataOprations = append(manageDataOprations, &queue.Operations[i])
			}
			logrus.Info("Received to queue")
			previousTxnHashBuilder := txnbuild.ManageData{
				Name:  "PREVIOUS TRANSACTION",
				Value: []byte(previousTxnHash),
			}
			manageDataOprations = append(manageDataOprations, &previousTxnHashBuilder)

			if queue.Type == "METRICBIND" {
				logrus.Info("Received mgs Type (METRICBIND)")
				startTime := time.Now()
				stellarprotocol := stellarprotocols.StellarTransaction{Operations: manageDataOprations, Memo: string(queue.Memo)}
				err, errCode, hash, sequenceNo, xdr, senderPK := stellarprotocol.SubmitToStellarBlockchain()
				endTime := time.Now()
				convertedTime := fmt.Sprintf("%f", endTime.Sub(startTime).Seconds())
				convertedCost := fmt.Sprintf("%f", 0.00001*float32(queue.ExpertFormula.NoOfManageDataInTxn+1))
				metricBindingStore := model.MetricBindingStore{
					MetricId:            queue.MetricBinding.Metric.ID,
					MetricMapID:         queue.MetricBinding.MetricMapID,
					Metric:              queue.MetricBinding.Metric,
					User:                queue.MetricBinding.User,
					TotalNoOfManageData: queue.MetricBinding.TotalNoOfManageData + (queue.MetricBinding.TotalNoOfManageData / 25) + 1,
					NoOfManageDataInTxn: queue.MetricBinding.NoOfManageDataInTxn + 1, // with previous transaction back-link
					TransactionTime:     convertedTime,
					TransactionCost:     convertedCost,
					Memo:                queue.Memo,
					Status:              "SUCCESS",
					XDR:                 xdr,
					TxnSenderPK:         senderPK,
					Timestamp:           time.Now().String(),
					TxnUUID:             queue.MetricBinding.TxnUUID,
				}
				if err != nil {
					metricBindingStore.ErrorMessage = err.Error()
					metricBindingStore.Status = "FAILED"
					errWhenUpdatingMetricBind := object.UpdateMetricBindStatus(queue.MetricBinding.MetricId, queue.MetricBinding.TxnUUID, metricBindingStore) // update -> metric id + txnUUID
					if errWhenUpdatingMetricBind != nil {
						logrus.Error("Error while updating the metric binding formula into DB: ", errWhenUpdatingMetricBind)
					}
					logrus.Info("Metric update called with failed status")
					logrus.Error("Stellar transacion submitting issue in queue", err, " error code ", errCode)
					logrus.Println("XDR  ", xdr)
				} else {
					metricBindingStore.SequenceNo = sequenceNo
					metricBindingStore.TxnHash = hash
					previousTxnHash = hash
					errWhenUpdatingMetricBind := object.UpdateMetricBindStatus(queue.MetricBinding.MetricId, queue.MetricBinding.TxnUUID, metricBindingStore)
					if errWhenUpdatingMetricBind != nil {
						logrus.Error("Error while updating the metric binding formula into DB: ", errWhenUpdatingMetricBind)
					}
					logrus.Info("Metric update called with success status")
					logrus.Info("-------------------------------------------------------------------------------------------------------------------------------------")
					logrus.Info("Stellar transacion submitting to blockchain (METRICBINDINIG) , Transaction Hash : ", hash)
					logrus.Info("-------------------------------------------------------------------------------------------------------------------------------------")
				}
			} else if queue.Type == "EXPERTFORMULA" {
				logrus.Info("Received mgs Type (EXPERTFORMULA)")
				startTime := time.Now()
				stellarprotocol := stellarprotocols.StellarTransaction{Operations: manageDataOprations, Memo: string(queue.Memo)}
				err, errCode, hash, sequenceNo, xdr, senderPK := stellarprotocol.SubmitToStellarBlockchain()
				expertFormulaStore := queue.ExpertFormula
				endTime := time.Now()
				convertedTime := fmt.Sprintf("%f", endTime.Sub(startTime).Seconds())
				convertedCost := fmt.Sprintf("%f", 0.00001*float32(queue.ExpertFormula.NoOfManageDataInTxn))
				expertFormulaStore = model.FormulaStore{
					MetricExpertFormula: queue.ExpertFormula.MetricExpertFormula,
					Verify:              queue.ExpertFormula.Verify,
					FormulaID:           queue.ExpertFormula.FormulaID,
					FormulaMapID:        queue.ExpertFormula.FormulaMapID,
					VariableCount:       queue.ExpertFormula.VariableCount,
					ExecutionTemplate:   queue.ExpertFormula.ExecutionTemplate,
					TotalNoOfManageData: queue.ExpertFormula.TotalNoOfManageData,
					NoOfManageDataInTxn: queue.ExpertFormula.NoOfManageDataInTxn,
					Memo:                queue.Memo,
					TxnHash:             hash,
					TxnSenderPK:         senderPK,
					XDR:                 xdr,
					SequenceNo:          sequenceNo,
					Status:              "SUCCESS",
					Timestamp:           time.Now().String(),
					TransactionTime:     convertedTime,
					TransactionCost:     convertedCost,
					ErrorMessage:        "",
					TxnUUID:             queue.ExpertFormula.TxnUUID,
				}
				if err != nil {
					expertFormulaStore.ErrorMessage = err.Error()
					expertFormulaStore.Status = "FAILED"
					errWhenUpdateingFormulaSatus := object.UpdateFormulaStatus(queue.ExpertFormula.FormulaID, queue.ExpertFormula.TxnUUID, expertFormulaStore) // update
					if errWhenUpdateingFormulaSatus != nil {
						logrus.Error("Error while updating the expert formula into DB: ", err)
					}
					logrus.Info("Formula update called with failed status")
					logrus.Error("Stellar transactions submitting issue in queue (EXPERTFORMULA)", err, " error code ", errCode)
					logrus.Println("XDR  ", xdr)
				} else {
					expertFormulaStore.SequenceNo = sequenceNo
					expertFormulaStore.TxnHash = hash
					previousTxnHash = hash
					errWhenUpdateingFormulaSatus := object.UpdateFormulaStatus(queue.ExpertFormula.FormulaID, queue.ExpertFormula.TxnUUID, expertFormulaStore) // update
					if errWhenUpdateingFormulaSatus != nil {
						logrus.Error("Error while updating the expert formula into DB: ", err)
					}
					logrus.Info("Formula update called with success status")
					logrus.Info("-------------------------------------------------------------------------------------------------------------------------------------")
					logrus.Info("Stellar transaction submitting to blockchain (EXPERTFORMULA) , Transaction Hash : ", hash)
					logrus.Info("-------------------------------------------------------------------------------------------------------------------------------------")
				}
			} else if queue.Type == "ETHEXPERTFORMULA" {
				logrus.Info("Received mgs Type (ETHEXPERTFORMULA)")
				startTime := time.Now()
				//Call the deploy method
				address, txnHash, deploymentCost, errWhenDeploying := deploy.DeployContract(queue.EthereumExpertFormula.ABIstring, queue.EthereumExpertFormula.BINstring)
				endTime := time.Now()
				convertedTime := fmt.Sprintf("%f", endTime.Sub(startTime).Seconds())
				ethExpertFormulaObj := model.EthereumExpertFormula{
					FormulaID:           queue.EthereumExpertFormula.FormulaID,
					FormulaName:         queue.EthereumExpertFormula.FormulaName,
					ExecutionTemplate:   queue.EthereumExpertFormula.ExecutionTemplate,
					MetricExpertFormula: queue.EthereumExpertFormula.MetricExpertFormula,
					VariableCount:       queue.EthereumExpertFormula.VariableCount,
					TemplateString:      queue.EthereumExpertFormula.TemplateString,
					BINstring:           queue.EthereumExpertFormula.BINstring,
					ABIstring:           queue.EthereumExpertFormula.ABIstring,
					GOstring:            queue.EthereumExpertFormula.GOstring,
					SetterNames:         queue.EthereumExpertFormula.SetterNames,
					ContractName:        queue.EthereumExpertFormula.ContractName,
					ContractAddress:     address,
					Timestamp:           time.Now().String(),
					TransactionHash:     txnHash,
					TransactionCost:     deploymentCost, //add after deploy
					TransactionTime:     convertedTime,
					TransactionUUID:     queue.EthereumExpertFormula.TransactionUUID,
					TransactionSender:   queue.EthereumExpertFormula.TransactionSender,
					User:                queue.EthereumExpertFormula.User,
					ErrorMessage:        "",
					Status:              "SUCCESS",
				}
				if errWhenDeploying != nil {
					//Insert to DB with FAILED status
					ethExpertFormulaObj.Status = "FAILED"
					ethExpertFormulaObj.ErrorMessage = errWhenDeploying.Error()
					logrus.Error("Error when deploying the expert formula smart contract : " + errWhenDeploying.Error())
					//if deploy method is success update the status into success
					errWhenUpdatingStatus := object.UpdateEthereumFormulaStatus(queue.EthereumExpertFormula.FormulaID, queue.EthereumExpertFormula.TransactionUUID, ethExpertFormulaObj)
					if errWhenUpdatingStatus != nil {
						logrus.Error("Error when updating the status of formula status for Eth , formula ID " + ethExpertFormulaObj.FormulaID)
					}
					logrus.Info("Formula update called with FAILED status")
					logrus.Info("Contract deployment unsuccessful")
				} else {
					//if deploy method is success update the status into success
					errWhenUpdatingStatus := object.UpdateEthereumFormulaStatus(queue.EthereumExpertFormula.FormulaID, queue.EthereumExpertFormula.TransactionUUID, ethExpertFormulaObj)
					if errWhenUpdatingStatus != nil {
						logrus.Error("Error when updating the status of formula status for Eth , formula ID " + ethExpertFormulaObj.FormulaID)
					}
					logrus.Info("Formula update called with SUCCESS status")
					logrus.Info("-------------------------------------------------------------------------------------------------------------------------------------")
					logrus.Info("Deployed expert expert formula smart contract to blockchain")
					logrus.Info("Contract address : " + address)
					logrus.Info("Transaction hash : " + txnHash)
					logrus.Info("-------------------------------------------------------------------------------------------------------------------------------------")

				}
			} else if queue.Type == "ETHMETRICBIND" {
				logrus.Info("Received mgs Type (ETHMETRICBIND)")
				startTime := time.Now()
				//Call the deploy method
				address, txnHash, deploymentCost, errWhenDeploying := deploy.DeployContract(queue.EthereumMetricBind.ABIstring, queue.EthereumMetricBind.BINstring)
				endTime := time.Now()
				convertedTime := fmt.Sprintf("%f", endTime.Sub(startTime).Seconds())
				ethMetricObj := model.EthereumMetricBind{
					MetricID:          queue.EthereumMetricBind.MetricID,
					MetricName:        queue.EthereumMetricBind.MetricName,
					Metric:            queue.EthereumMetricBind.Metric,
					ContractName:      queue.EthereumMetricBind.ContractName,
					TemplateString:    queue.EthereumMetricBind.TemplateString,
					BINstring:         queue.EthereumMetricBind.BINstring,
					ABIstring:         queue.EthereumMetricBind.ABIstring,
					Timestamp:         time.Now().String(),
					ContractAddress:   address,
					TransactionHash:   txnHash,
					TransactionCost:   deploymentCost,
					TransactionTime:   convertedTime,
					TransactionUUID:   queue.EthereumMetricBind.TransactionUUID,
					TransactionSender: queue.EthereumMetricBind.TransactionSender,
					User:              queue.EthereumMetricBind.User,
					ErrorMessage:      "",
					Status:            "SUCCESS",
					FormulaIDs:        queue.EthereumMetricBind.FormulaIDs,
					ValueIDs:          queue.EthereumMetricBind.ValueIDs,
					Type:              queue.EthereumMetricBind.Type,
					FormulaID:         queue.EthereumExpertFormula.FormulaID,
				}
				if errWhenDeploying != nil {
					//Insert to DB with FAILED status
					ethMetricObj.Status = "FAILED"
					ethMetricObj.ErrorMessage = errWhenDeploying.Error()
					logrus.Error("Error when deploying the metric bind smart contract : " + errWhenDeploying.Error())
					//if deploy method is success update the status into success
					errWhenUpdatingStatus := object.UpdateEthereumMetricStatus(queue.EthereumMetricBind.MetricID, queue.EthereumMetricBind.TransactionUUID, ethMetricObj)
					if errWhenUpdatingStatus != nil {
						logrus.Error("Error when updating the status of metric status for Eth , formula ID " + ethMetricObj.MetricID)
					}
					logrus.Info("Metric update called with FAILED status")
					logrus.Info("Contract deployment unsuccessful")
				} else {
					//if deploy method is success update the status into success
					errWhenUpdatingStatus := object.UpdateEthereumMetricStatus(queue.EthereumMetricBind.MetricID, queue.EthereumMetricBind.TransactionUUID, ethMetricObj)
					if errWhenUpdatingStatus != nil {
						logrus.Error("Error when updating the status of metric status for Eth , formula ID " + ethMetricObj.MetricID)
					}
					logrus.Info("Metric update called with SUCCESS status")

					insertObj := model.MetricLatestContract{
						MetricID:        ethMetricObj.MetricID,
						ContractAddress: address,
						Type:            ethMetricObj.Type,
					}
					if ethMetricObj.Type == "METADATA" {
						//insert the latest contract address in DB
						errWhenInsertingToLatest := object.EthereumInsertToMetricLatestContract(insertObj)
						if errWhenInsertingToLatest != nil {
							logrus.Error("Error when inserting to latest contract to DB: ", errWhenInsertingToLatest)
						}
						logrus.Info("Added " + address + " to latest contract collection")
					} else if ethMetricObj.Type == "ACTIVITY" {
						//update the latest contract address in DB
						errWhenUpdatingLatest := object.UpdateEthereumMetricLatestContract(ethMetricObj.MetricID, insertObj)
						if errWhenUpdatingLatest != nil {
							logrus.Errorf("Error when updating latest contract address in DB: ", errWhenUpdatingLatest)
						}
						logrus.Info("Updated " + address + " as latest contract")
					}
					logrus.Info("-------------------------------------------------------------------------------------------------------------------------------------")
					logrus.Info("Deployed expert metric bind smart contract to blockchain")
					logrus.Info("Contract address : " + address)
					logrus.Info("Transaction hash : " + txnHash)
					logrus.Info("-------------------------------------------------------------------------------------------------------------------------------------")
				}
			}
		}
	}()

	logrus.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}

func SendToQueue(queue model.SendToQueue) error {
	rabbitConnection := `amqp://` + USER + `:` + PASSWORD + `@` + HOSTNAME + `:` + PORT + `/`
	conn, err := amqp.Dial(rabbitConnection)
	if err != nil {
		logrus.Error("rabbitmq connection issue ", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logrus.Error("%s: %s", "Failed to open a channel", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"socialimpact", // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		logrus.Error("%s: %s", "Failed to declare a queue", err)
	}
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(queue)
	body := reqBodyBytes.Bytes()
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         []byte(body),
		})

	if err != nil {
		logrus.Error("%s: %s", "Failed to publish a message", err)
	}
	return nil
}
