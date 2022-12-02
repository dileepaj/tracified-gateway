package ethereuemmetricbind

import (
	"encoding/base64"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/oklog/ulid"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	contractName = ``
)

/*
Generate smart contract for metric binding
*/
func SmartContractGeneratorForMetric(w http.ResponseWriter, r *http.Request, metricBindJson model.MetricDataBindingRequest) {
	object := dao.Connection{}
	var status string

	metricDetails, errWhenGettingMetricStatus := object.GetEthMetricStatus(metricBindJson.Metric.ID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errWhenGettingMetricStatus != nil {
		logrus.Error("An error occurred when getting metric status ", errWhenGettingMetricStatus)
	}
	if metricDetails == nil {
		status = ""
	} else {
		status = metricDetails.(model.EthereumMetricBind).Status
		logrus.Info("Deploy status : ", status)
	}

	if status == "SUCCESS" {
		logrus.Info("Contract for metric " + metricBindJson.Metric.Name + " has been added to the blockchain")
		commons.JSONErrorReturn(w, r, "Status : "+status, 400, "Requested metric is in the blockchain")
		return
	} else if status == "QUEUE" {
		logrus.Info("Requested metric is in the queue, please try again")
		commons.JSONErrorReturn(w, r, "Status : "+status, 400, "Requested metric is in the queue, please try again")
		return
	} else if status == "" || status == "FAILED" {
		if status == "FAILED" {
			logrus.Info("Requested metric is in the failed status, trying to redeploy")
		} else {
			logrus.Info("New metric bind request, initiating new deployment")
		}

		ethMetricObj := model.EthereumMetricBind{
			MetricID:          metricBindJson.Metric.ID,
			MetricName:        metricBindJson.Metric.Name,
			Metric:            metricBindJson.Metric,
			ContractName:      "",
			TemplateString:    "",
			BINstring:         "",
			ABIstring:         "",
			Timestamp:         time.Now().String(),
			ContractAddress:   "",
			TransactionHash:   "",
			TransactionCost:   "",
			TransactionTime:   "",
			TransactionUUID:   "",
			TransactionSender: commons.GoDotEnvVariable("ETHEREUMPUBKEY"),
			User:              metricBindJson.User,
			ErrorMessage:      "",
			Status:            "",
		}

		if status == "" {
			//generate transaction UUID
			timeNow := time.Now().UTC()
			entropy := rand.New(rand.NewSource(timeNow.UnixNano()))
			id := ulid.MustNew(ulid.Timestamp(timeNow), entropy)
			logrus.Info("TXN UUID : ", id)
			ethMetricObj.TransactionUUID = id.String()
		} else {
			ethMetricObj.TransactionUUID = metricDetails.(model.EthereumExpertFormula).TransactionUUID
		}

		//setting up the contract name and starting the contract
		contractName = cases.Title(language.English).String(metricBindJson.Metric.Name)
		contractName = strings.ReplaceAll(contractName, " ", "")

		//TODO Contract writer components

		template := ""
		b64Template := base64.StdEncoding.EncodeToString([]byte(template))
		ethMetricObj.TemplateString = b64Template

		//Write contract template into a file
		//fo, errInOutput := os.Create(commons.GoDotEnvVariable("METRICCONTRACTLOCATION") + "/" + contractName + `.sol`)

	} else {
		logrus.Info("Invalid metric status " + status)
		commons.JSONErrorReturn(w, r, status, http.StatusInternalServerError, "Invalid metric status : ")
		return
	}

}
