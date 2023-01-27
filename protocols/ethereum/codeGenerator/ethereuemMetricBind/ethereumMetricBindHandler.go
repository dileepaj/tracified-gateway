package ethereuemmetricbind

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	activityWriters "github.com/dileepaj/tracified-gateway/protocols/ethereum/codeGenerator/ethereuemMetricBind/activityContractWriters"
	metadataWriters "github.com/dileepaj/tracified-gateway/protocols/ethereum/codeGenerator/ethereuemMetricBind/metadataWriters"
	"github.com/dileepaj/tracified-gateway/vendor/github.com/sirupsen/logrus"
	"github.com/oklog/ulid"
)

func SmartContractHandlerForMetric(w http.ResponseWriter, r *http.Request, metricBindJson model.MetricDataBindingRequest) {
	object := dao.Connection{}

	// create the metadata object
	metaDataObj := model.MetricMetadataReq{
		MetricId:     metricBindJson.Metric.ID,
		MetricName:   metricBindJson.Metric.Name,
		TenetId:      metricBindJson.Metric.TenantId,
		NoOfFormulas: len(metricBindJson.Metric.MetricActivities),
		TrustNetPK:   metricBindJson.User.Publickey,
		Metric:       metricBindJson.Metric,
		User:         metricBindJson.User,
	}

	// get the activity list of the metric
	activities := metricBindJson.Metric.MetricActivities

	// metric object for metadata smart contract to be stored in the database
	ethMetricObjForMetaData := model.EthereumMetricBind{
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
		Type:              "METADATA",
		FormulaID:         "",
	}

	// get metric map id
	metricMapId, errWhenGettingMetricMapId := GetMetricMapId(metricBindJson.Metric.ID)
	if errWhenGettingMetricMapId != nil {
		logrus.Info(errWhenGettingMetricMapId)
		return
	}
	metricMapIDString := strconv.FormatUint(metricMapId, 10)

	contractName := "Metric_" + metricMapIDString
	ethMetricObjForMetaData.ContractName = contractName

	// get the status of the metric metadata contract
	status, metricDetails, errWhenGettingMetadataContractStatus := GetMetricSmartContractStatus(metricBindJson.Metric.ID, "METADATA")
	if errWhenGettingMetadataContractStatus != nil {
		// TODO: handle the error
		logrus.Info(errWhenGettingMetadataContractStatus)
		return
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
	}

	if status == "" {
		//generate transaction UUID
		timeNow := time.Now().UTC()
		entropy := rand.New(rand.NewSource(timeNow.UnixNano()))
		id := ulid.MustNew(ulid.Timestamp(timeNow), entropy)
		logrus.Info("TXN UUID : ", id)
		ethMetricObjForMetaData.TransactionUUID = id.String()
	} else {
		// get the transaction UUID from the database
		ethMetricObjForMetaData.TransactionUUID = metricDetails.TransactionUUID
	}

	if status == "" || status == "FAILED" {
		// deploy the smart contract for meta data
		errWhenDeployingMetaDataSmartContract := metadataWriters.MetricMetadataContractDeployer(metaDataObj, metricMapIDString)
		if errWhenDeployingMetaDataSmartContract != nil {
			ethMetricObjForMetaData.ErrorMessage = errWhenDeployingMetaDataSmartContract.Error()
			ethMetricObjForMetaData.Status = "FAILED"
			if status == "" {
				// store the metric object in the database
				errWhenStoringMetricObj := object.InsertToEthMetricDetails(ethMetricObjForMetaData)
				if errWhenStoringMetricObj != nil {
					logrus.Info(errWhenStoringMetricObj)
					return
				}
			} else if status == "FAILED" {
				// update the metric object in the database
				errWhenUpdatingMetricObj := object.UpdateEthereumMetricStatus(ethMetricObjForMetaData.MetricID, ethMetricObjForMetaData.TransactionUUID, ethMetricObjForMetaData)
				if errWhenUpdatingMetricObj != nil {
					logrus.Info(errWhenUpdatingMetricObj)
					return
				}
			}
			logrus.Info(errWhenDeployingMetaDataSmartContract)
			return
		} else {
			ethMetricObjForMetaData.Status = "SUCCESS"
			// update the metric object in the database
			if status == "FAILED" {
				errWhenUpdatingMetricObj := object.UpdateEthereumMetricStatus(ethMetricObjForMetaData.MetricID, ethMetricObjForMetaData.TransactionUUID, ethMetricObjForMetaData)
				if errWhenUpdatingMetricObj != nil {
					logrus.Info(errWhenUpdatingMetricObj)
					return
				}
			} else if status == "" {
				// store the metric object in the database
				errWhenStoringMetricObj := object.InsertToEthMetricDetails(ethMetricObjForMetaData)
				if errWhenStoringMetricObj != nil {
					logrus.Info(errWhenStoringMetricObj)
					return
				}
			}
		}
	}

	// get the status of the metric metadata contract after deploying(/redeploying) the contract
	status, metricDetails, errWhenGettingMetadataContractStatus = GetMetricSmartContractStatus(metricBindJson.Metric.ID, "METADATA")
	if errWhenGettingMetadataContractStatus != nil {
		// TODO: handle the error
		logrus.Info(errWhenGettingMetadataContractStatus)
		return
	}

	// check if the status is SUCCESS or not, if SUCCESS then proceed to create the smart contract for the metric activities
	// TODO: handle status QUEUE
	if status == "SUCCESS" {
		if len(activities) > 0 {
			for _, activity := range activities {
				// TODO: validation to check whether the contract is already deployed or not(metric collection) using metric id and formula id
				// if it is deployed then skip the deployment
				// check the index of the loop, if it is not 0 then check the index-1 activity contract status from the metric collection
				// if the status is SUCCESS then deploy the activity contract

				errWhenDeployingActivityContract := activityWriters.ActivityContractDeployer(metricBindJson.Metric.ID, activity, metricBindJson.Metric.Name, metricBindJson.Metric, metricBindJson.User)
				if errWhenDeployingActivityContract != nil {
					// TODO: handle the error
					logrus.Info(errWhenDeployingActivityContract)
				}
			}
		}
	}

}
