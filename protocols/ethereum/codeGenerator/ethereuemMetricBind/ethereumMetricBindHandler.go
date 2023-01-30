package ethereuemmetricbind

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	activityWriters "github.com/dileepaj/tracified-gateway/protocols/ethereum/codeGenerator/ethereuemMetricBind/ActivityContractWriters"
	metadataWriters "github.com/dileepaj/tracified-gateway/protocols/ethereum/codeGenerator/ethereuemMetricBind/metadataWriters"
	"github.com/oklog/ulid"
	"github.com/sirupsen/logrus"
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

	metadataContractName := "Metric_" + metricMapIDString
	ethMetricObjForMetaData.ContractName = metadataContractName

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
			ethMetricObjForMetaData.Status = "QUEUE"
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
		//TODO : assign to a boolean variable and will be checked before deploying first activity
	}

	// check if the status is SUCCESS or not, if SUCCESS then proceed to create the smart contract for the metric activities
	// TODO: handle status QUEUE
	if status == "SUCCESS" {
		if len(activities) > 0 {

			for i := 0; i < len(activities); i++ {
				//check if the contract for the this metric ID + formula ID + type deployed
				formulaStatus, formulaDetails, errWhenGettingFormulaStatus := GetMetricSmartContractStatusForFormula(metricBindJson.Metric.ID, "ACTIVITY", activities[i].MetricFormula.MetricExpertFormula.ID)
				if errWhenGettingFormulaStatus != nil {
					//TODO: handle error
					logrus.Info(errWhenGettingMetadataContractStatus)
					return
				}

				logrus.Info(formulaDetails)

				//todo : handle the DB update and insert with UUID

				if formulaStatus == "SUCCESS" || formulaStatus == "QUEUE" {
					//skip this loop and go to next formula
					continue
				} else if formulaStatus == "" || formulaStatus == "FAILED" {
					//check the index of the loop to skip the checking of the previous formula deployment
					var canCallNextDeployment bool
					if i != 0 {
						//check the previous formula contract deployment status
						previousStatus, _, errWhenGettingPreviousStatus := GetMetricSmartContractStatusForFormula(metricBindJson.Metric.ID, "ACTIVITY", activities[i-1].MetricFormula.MetricExpertFormula.ID)
						if errWhenGettingPreviousStatus != nil {
							//TODO: handle error
							logrus.Info(errWhenGettingMetadataContractStatus)
							return
						}
						if previousStatus == "SUCCESS" {
							canCallNextDeployment = true
						} else if previousStatus == "FAILED" || previousStatus == "" || previousStatus == "QUEUE" {
							canCallNextDeployment = false
						}
					} else if i == 0 {
						canCallNextDeployment = true
					}

					if canCallNextDeployment {
						// get the formula map id form DB
						formulaMapID, errWhenGettingFormulaMapId := GetFormulaMapId(activities[i].MetricFormula.MetricExpertFormula.ID)
						if errWhenGettingFormulaMapId != nil {
							logrus.Info(errWhenGettingFormulaMapId)
							return
						}
						formulaMapIDString := strconv.FormatUint(formulaMapID, 10)
						activityContractName := "Metric_" + metricMapIDString + "_" + formulaMapIDString
						_ = activityContractName

						errWhenDeployingActivityContract := activityWriters.ActivityContractDeployer(metricMapIDString, formulaMapIDString, metricBindJson.Metric.ID, activities[i], metricBindJson.Metric.Name, metricBindJson.Metric, metricBindJson.User)
						if errWhenDeployingActivityContract != nil {
							// TODO: handle the error
							logrus.Info(errWhenDeployingActivityContract)
						}
					}
				}

			}
		}
	}

}
