package rabbitmq

import (
	"bytes"
	"encoding/json"
	"fmt"

	"time"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
	ethereumservices "github.com/dileepaj/tracified-gateway/services/ethereumServices"
	"github.com/dileepaj/tracified-gateway/services/ethereumServices/dbCollectionHandler"
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
					MetricId:              queue.MetricBinding.Metric.ID,
					MetricMapID:           queue.MetricBinding.MetricMapID,
					Metric:                queue.MetricBinding.Metric,
					User:                  queue.MetricBinding.User,
					TotalNoOfManageData:   queue.MetricBinding.TotalNoOfManageData + (queue.MetricBinding.TotalNoOfManageData / 25) + 1,
					NoOfManageDataInTxn:   queue.MetricBinding.NoOfManageDataInTxn + 1, // with previous transaction back-link
					TransactionTime:       convertedTime,
					TransactionCost:       convertedCost,
					Memo:                  queue.Memo,
					Status:                "SUCCESS",
					XDR:                   xdr,
					TxnSenderPK:           senderPK,
					Timestamp:             time.Now().String(),
					TxnUUID:               queue.MetricBinding.TxnUUID,
					ActivityManageDataMap: queue.MetricBinding.ActivityManageDataMap,
					TransactionOrderCount: queue.TransactionCount + 1,
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
				// use deployment strategy
				expertDeployer := &ethereumservices.ContractDeployerContext{}
				expertDeployer.SetContractDeploymentStrategy(&ethereumservices.AbstractContractDeployment{
					ABI:          queue.EthereumExpertFormula.ABIstring,
					BIN:          queue.EthereumExpertFormula.BINstring,
					Identifier:   queue.EthereumExpertFormula.TransactionUUID,
					ContractType: "ETHEXPERTFORMULA",
					OtherParams:  []any{queue.EthereumExpertFormula},
				})
				//Call the deploy method
				address, txnHash, deploymentCost, _, _, _, errWhenDeploying := expertDeployer.ExecuteContractDeployment()
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
					TransactionUUID:     queue.EthereumExpertFormula.TransactionUUID,
					TransactionSender:   queue.EthereumExpertFormula.TransactionSender,
					Verify:              queue.EthereumExpertFormula.Verify,
					ErrorMessage:        "",
					Status:              "PENDING",
					ActualStatus:        queue.EthereumExpertFormula.ActualStatus,
				}
				if errWhenDeploying != nil {
					//Insert to DB with FAILED status
					ethExpertFormulaObj.Status = "FAILED"
					ethExpertFormulaObj.ErrorMessage = errWhenDeploying.Error()
					ethExpertFormulaObj.ActualStatus = 111 // DEPLOYMENT_FAILED
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
					ethExpertFormulaObj.ActualStatus = 112 // DEPLOYMENT_TRANSACTION_PENDING
					errWhenUpdatingStatus := object.UpdateEthereumFormulaStatus(queue.EthereumExpertFormula.FormulaID, queue.EthereumExpertFormula.TransactionUUID, ethExpertFormulaObj)
					if errWhenUpdatingStatus != nil {
						logrus.Error("Error when updating the status of formula status for Eth , formula ID " + ethExpertFormulaObj.FormulaID)
					}
					logrus.Info("Formula update called with status ", ethExpertFormulaObj.Status)
					logrus.Info("-------------------------------------------------------------------------------------------------------------------------------------")
					logrus.Info("Deployed expert expert formula smart contract to blockchain")
					logrus.Info("Contract address : " + address)
					logrus.Info("Transaction hash : " + txnHash)
					logrus.Info("-------------------------------------------------------------------------------------------------------------------------------------")

				}
			} else if queue.Type == "ETHMETRICBIND" {
				logrus.Info("Received mgs Type (ETHMETRICBIND)")
				// use deployment strategy
				metricDeployer := &ethereumservices.ContractDeployerContext{}
				metricDeployer.SetContractDeploymentStrategy(&ethereumservices.AbstractContractDeployment{
					ABI:          queue.EthereumMetricBind.ABIstring,
					BIN:          queue.EthereumMetricBind.BINstring,
					Identifier:   queue.EthereumMetricBind.TransactionUUID,
					ContractType: "ETHMETRICBIND",
					OtherParams:  []any{queue.EthereumMetricBind},
				})
				//Call the deploy method
				address, txnHash, deploymentCost, _, _, _, errWhenDeploying := metricDeployer.ExecuteContractDeployment()
				ethMetricObj := model.EthereumMetricBind{
					MetricID:          queue.EthereumMetricBind.MetricID,
					MetricName:        queue.EthereumMetricBind.MetricName,
					Metric:            queue.EthereumMetricBind.Metric,
					ContractName:      queue.EthereumMetricBind.ContractName,
					TemplateString:    queue.EthereumMetricBind.TemplateString,
					BINstring:         queue.EthereumMetricBind.BINstring,
					ABIstring:         queue.EthereumMetricBind.ABIstring,
					Timestamp:         time.Now().UTC().String(),
					ContractAddress:   address,
					TransactionHash:   txnHash,
					TransactionCost:   deploymentCost,
					TransactionUUID:   queue.EthereumMetricBind.TransactionUUID,
					TransactionSender: queue.EthereumMetricBind.TransactionSender,
					User:              queue.EthereumMetricBind.User,
					ErrorMessage:      "",
					Status:            "PENDING",
					FormulaIDs:        queue.EthereumMetricBind.FormulaIDs,
					ValueIDs:          queue.EthereumMetricBind.ValueIDs,
					Type:              queue.EthereumMetricBind.Type,
					FormulaID:         queue.EthereumMetricBind.FormulaID,
					ActualStatus:      queue.EthereumMetricBind.ActualStatus,
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
					pendingContract := model.PendingContracts{
						ContractAddress: ethMetricObj.ContractAddress,
						ContractType:    "ETHMETRICBIND",
						Identifier:      ethMetricObj.TransactionUUID,
						TransactionHash: ethMetricObj.TransactionHash,
						Status:          "FAILED",
						ErrorMessage:    ethMetricObj.ErrorMessage,
					}
					errWhenInvalidatingMetric := dbCollectionHandler.InvalidateMetric(pendingContract, ethMetricObj.Status, ethMetricObj.ErrorMessage)
					if errWhenInvalidatingMetric != nil {
						logrus.Error("Error when invalidating the metric : " + queue.EthereumMetricBind.MetricID)
					}
					logrus.Info("Metric update called with FAILED status. Type: " + ethMetricObj.Type)
					logrus.Info("Contract deployment unsuccessful")
				} else {
					//if deploy method is success update the status into success
					errWhenUpdatingStatus := object.UpdateEthereumMetricStatus(queue.EthereumMetricBind.MetricID, queue.EthereumMetricBind.TransactionUUID, ethMetricObj)
					if errWhenUpdatingStatus != nil {
						logrus.Error("Error when updating the status of metric status for Eth , formula ID " + ethMetricObj.MetricID)
					}
					logrus.Info("Metric update called with status " + ethMetricObj.Status + ". Type: " + ethMetricObj.Type)
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

// rabbitmq queue message listener
func ReleaseLock() (error, string) {
	logrus.Info("-------------------------------Release Lock")
	object := dao.Connection{}
	var message string
	rabbitConnection := `amqp://` + USER + `:` + PASSWORD + `@` + HOSTNAME + `:` + PORT + `/`
	conn, err := amqp.Dial(rabbitConnection)
	if err != nil {
		logrus.Error("rabbitmq connection issue ", err)
		return err, ""
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		logrus.Error("%s: %s", "Failed to open a channel", err)
	}
	defer channel.Close()

	// Declare the queue
	queuelock, err := channel.QueueDeclare(
		"buyingnfts",
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err, ""
	}

	// Acknowledge the message to release the lock
	msgs, err := channel.Consume(
		queuelock.Name,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err, ""
	}
	var forever chan struct{}
	go func() {
		for d := range msgs {
			var pendingnft model.PendingNFTS
			if err := json.Unmarshal(d.Body, &pendingnft); err != nil {
				logrus.Error("Unmarshal in rabbitmq reciverRmq ", err.Error())
			}
			pendingNFTS := model.PendingNFTS{
				Blockchain:    pendingnft.Blockchain,
				NFTIdentifier: pendingnft.NFTIdentifier,
				Status:        pendingnft.Status,
				ImageBase64:   pendingnft.ImageBase64,
				User:          pendingnft.User,
				Version:       pendingnft.Version,
			}
			nfts, errWhenGettingNFTS := object.GetNFTById1Id2Id3Id4(pendingNFTS.Blockchain, pendingNFTS.NFTIdentifier, pendingNFTS.ImageBase64, pendingNFTS.Version).Then(func(data interface{}) interface{} {
				return data
			}).Await()
			logrus.Info("after getting", errWhenGettingNFTS)
			if errWhenGettingNFTS != nil {
				logrus.Error("error when retrieving the nfts", nfts)
				logrus.Info("retreived nfts: ", pendingnft)
				pendingNFTS := model.PendingNFTS{
					Blockchain:    pendingnft.Blockchain,
					NFTIdentifier: pendingnft.NFTIdentifier,
					Status:        "PROCESSED",
					ImageBase64:   pendingnft.ImageBase64,
					User:          pendingnft.User,
					Version:       pendingnft.Version,
				}
				errWhenUpdatingNFTStatus := object.InsertToNFTStatus(pendingNFTS)
				if errWhenUpdatingNFTStatus != nil {
					logrus.Error("Error when updating the status of the NFT ")
				}
				logrus.Info("NFT update called with status 200")
				message = pendingNFTS.User
			}
			if pendingNFTS.Status == "PROCESSING" {
				message = "none"
			}
		}
	}()
	logrus.Info(" [*] release lock Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil, message
}

// push element to rabbitmq queue
func LockRequest(pendingNFTS model.PendingNFTS) error {
	rabbitConnection := `amqp://` + USER + `:` + PASSWORD + `@` + HOSTNAME + `:` + PORT + `/`
	conn, err := amqp.Dial(rabbitConnection)
	if err != nil {
		logrus.Error("rabbitmq connection issue LockRequest", err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logrus.Error("%s: %s", "Failed to open a channel LockRequest", err)
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"buyingnfts", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		logrus.Error("%s: %s", "Failed to declare a queue LockRequest", err)
	}
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(pendingNFTS)
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
		logrus.Error("%s: %s", "Failed to publish a message LockRequest", err)
	}

	return nil
}
