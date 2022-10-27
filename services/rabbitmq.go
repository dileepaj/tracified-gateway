package services

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

var (
	USER     = commons.GoDotEnvVariable("RABBITUSER")
	PASSWORD = commons.GoDotEnvVariable("RABBITPASSWORD")
	HOSTNAME = commons.GoDotEnvVariable("RABBITHOSTNAME")
	PORT     = commons.GoDotEnvVariable("RABBITPORT")
)

func ReciverRmq() error {
	rabbitCoonection := `amqp://` + USER + `:` + PASSWORD + `@` + HOSTNAME + `:` + PORT + `/`
	conn, err := amqp.Dial(rabbitCoonection)
	if err != nil {
		logrus.Error("Failed to connect to RabbitMQ ", err)
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logrus.Error("%s: %s", "Failed to open a channel", err)
		return err
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
		return err
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
		return err
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
			for _, manageData := range queue.Operations {
				manageDataOprations = append(manageDataOprations, &manageData)
			}
			logrus.Info("Received to queue")
			if queue.Type == "METRICBIND" {
				logrus.Info("Received mgs Type (METRICBIND)")
				stellarprotocol := stellarprotocols.StellarTrasaction{Operations: manageDataOprations, Memo: string(queue.Memo)}
				err, errCode, hash, sequenceNo, xdr, senderPK := stellarprotocol.SubmitToStellerBlockchain()
				metricBindingStore := model.MetricBindingStore{
					MetricId:    queue.MetricBinding.Metric.ID,
					MetricMapID: queue.MetricBinding.MetricMapID,
					Metric:      queue.MetricBinding.Metric,
					User:        queue.MetricBinding.User,
					Memo:        queue.Memo,
					Status:      "SUCCESS",
					XDR:         xdr,
					TxnSenderPK: senderPK,
					Timestamp:   time.Now().String(),
				}
				if err != nil {
					metricBindingStore.ErrorMessage = err.Error()
					metricBindingStore.ErrorMessage = err.Error()
					metricBindingStore.Status = "Falied"
					_, err := object.InsertMetricBindingFormula(metricBindingStore)
					if err != nil {
						logrus.Error("Error while inserting the metric binding formula into DB: ", err)
					}
					logrus.Error("Stellar transacion submitting issue in queue", err, " error code ", errCode)
					logrus.Println("XDR  ", xdr)
				} else {
					metricBindingStore.SequenceNo = sequenceNo
					metricBindingStore.TxnHash = hash
					_, err := object.InsertMetricBindingFormula(metricBindingStore)
					if err != nil {
						logrus.Error("Error while inserting the metric binding formula into DB: ", err)
					}
					logrus.Info("-------------------------------------------------------------------------------------------------------------------------------------")
					logrus.Info("Stellar transacion submitting to blockchain (METRICBINDINIG) , Transaction Hash : ", hash)
					logrus.Info("-------------------------------------------------------------------------------------------------------------------------------------")
				}
			} else if queue.Type == "EXPERTFORMULA" {
				logrus.Info("Received mgs Type (EXPERTFORMULA)")
				stellarprotocol := stellarprotocols.StellarTrasaction{Operations: manageDataOprations, Memo: string(queue.Memo)}
				err, errCode, hash, sequenceNo, xdr, senderPK := stellarprotocol.SubmitToStellerBlockchain()
				metricBindingStore := model.MetricBindingStore{
					MetricId:    queue.MetricBinding.Metric.ID,
					MetricMapID: queue.MetricBinding.MetricMapID,
					Metric:      queue.MetricBinding.Metric,
					User:        queue.MetricBinding.User,
					Memo:        queue.Memo,
					Status:      "SUCCESS",
					XDR:         xdr,
					TxnSenderPK: senderPK,
					Timestamp:   time.Now().String(),
				}
				if err != nil {
					metricBindingStore.ErrorMessage = err.Error()
					metricBindingStore.ErrorMessage = err.Error()
					metricBindingStore.Status = "Falied"
					_, err := object.InsertMetricBindingFormula(metricBindingStore)
					if err != nil {
						logrus.Error("Error while inserting the metric binding formula into DB: ", err)
					}
					logrus.Error("Stellar transacion submitting issue in queue", err, " error code ", errCode)
					logrus.Println("XDR  ", xdr)
				} else {
					metricBindingStore.SequenceNo = sequenceNo
					metricBindingStore.TxnHash = hash
					_, err := object.InsertMetricBindingFormula(metricBindingStore)
					if err != nil {
						logrus.Error("Error while inserting the metric binding formula into DB: ", err)
					}
					logrus.Info("-------------------------------------------------------------------------------------------------------------------------------------")
					logrus.Info("Stellar transacion submitting to blockchain (METRICBINDINIG) , Transaction Hash : ", hash)
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
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logrus.Error("%s: %s", "Failed to open a channel", err)
		return err
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
		return err
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
		return err
	}
	return nil
}
