package ethereuemmetricbind

import (
	"encoding/json"
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
	var canCallNextDeployment bool

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
		logrus.Error("Error when getting metric map ID : ", errWhenGettingMetricMapId)
	}
	metricMapIDString := strconv.FormatUint(metricMapId, 10)

	metadataContractName := "Metric_" + metricMapIDString
	ethMetricObjForMetaData.ContractName = metadataContractName

	// get the status of the metric metadata contract
	status, metricDetails, errWhenGettingMetadataContractStatus := GetMetricSmartContractStatus(metricBindJson.Metric.ID, "METADATA")
	if errWhenGettingMetadataContractStatus != nil {
		logrus.Error("Error when getting metadata contract status : ", errWhenGettingMetadataContractStatus)
	}

	if status == "SUCCESS" {
		logrus.Info("Contract for metric " + metricBindJson.Metric.Name + " has been added to the blockchain, checking the activity array")
		canCallNextDeployment = true
	} else if status == "QUEUE" {
		logrus.Info("Requested metric is in the queue, please try again")
		canCallNextDeployment = true
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
		errWhenDeployingMetaDataSmartContract := metadataWriters.MetricMetadataContractDeployer(metaDataObj, metricMapIDString, ethMetricObjForMetaData)
		if errWhenDeployingMetaDataSmartContract != nil {
			ethMetricObjForMetaData.ErrorMessage = errWhenDeployingMetaDataSmartContract.Error()
			ethMetricObjForMetaData.Status = "FAILED"
			if status == "" {
				// store the metric object in the database
				errWhenStoringMetricObj := object.InsertToEthMetricDetails(ethMetricObjForMetaData)
				if errWhenStoringMetricObj != nil {
					logrus.Info("Error when inserting to metric collection : ", errWhenStoringMetricObj)
					commons.JSONErrorReturn(w, r, errWhenStoringMetricObj.Error(), 500, "Error when inserting to metric collection : ")
					return
				}
			} else if status == "FAILED" {
				// update the metric object in the database
				errWhenUpdatingMetricObj := object.UpdateEthereumMetricStatus(ethMetricObjForMetaData.MetricID, ethMetricObjForMetaData.TransactionUUID, ethMetricObjForMetaData)
				if errWhenUpdatingMetricObj != nil {
					logrus.Info("Error when updating the metric collection : ", errWhenUpdatingMetricObj)
					commons.JSONErrorReturn(w, r, errWhenUpdatingMetricObj.Error(), 500, "Error when updating the metric collection : ")
					return
				}
			}
			logrus.Info("Error when deploying metadata metric contract : ", errWhenDeployingMetaDataSmartContract)
			commons.JSONErrorReturn(w, r, errWhenDeployingMetaDataSmartContract.Error(), 500, "Error when deploying metadata metric contract : ")
			return
		} else {
			ethMetricObjForMetaData.Status = "QUEUE"
			// update the metric object in the database
			if status == "FAILED" {
				errWhenUpdatingMetricObj := object.UpdateEthereumMetricStatus(ethMetricObjForMetaData.MetricID, ethMetricObjForMetaData.TransactionUUID, ethMetricObjForMetaData)
				if errWhenUpdatingMetricObj != nil {
					if errWhenUpdatingMetricObj != nil {
						logrus.Info("Error when updating the metric collection : ", errWhenUpdatingMetricObj)
						commons.JSONErrorReturn(w, r, errWhenUpdatingMetricObj.Error(), 500, "Error when updating the metric collection : ")
						return
					}
				}
			} else if status == "" {
				// store the metric object in the database
				errWhenStoringMetricObj := object.InsertToEthMetricDetails(ethMetricObjForMetaData)
				if errWhenStoringMetricObj != nil {
					if errWhenStoringMetricObj != nil {
						logrus.Info("Error when inserting to metric collection : ", errWhenStoringMetricObj)
						commons.JSONErrorReturn(w, r, errWhenStoringMetricObj.Error(), 500, "Error when inserting to metric collection : ")
						return
					}
				}
			}
		}
	}

	// get the status of the metric metadata contract after deploying(/redeploying) the contract
	status, metricDetails, errWhenGettingMetadataContractStatus = GetMetricSmartContractStatus(metricBindJson.Metric.ID, "METADATA")
	if errWhenGettingMetadataContractStatus != nil {
		logrus.Info("Error when getting metadata contract details from metric collection : ", errWhenGettingMetadataContractStatus)
		commons.JSONErrorReturn(w, r, errWhenGettingMetadataContractStatus.Error(), 500, "Error when getting metadata contract details from metric collection : ")
		return
	}

	// check if the status is SUCCESS or not, if SUCCESS then proceed to create the smart contract for the metric activities
	if status == "SUCCESS" || status == "QUEUE" {
		canCallNextDeployment = true
		if len(activities) > 0 {
			for i := 0; i < len(activities); i++ {
				//check if the contract for the this metric ID + formula ID + type deployed
				formulaStatus, formulaDetails, errWhenGettingFormulaStatus := GetMetricSmartContractStatusForFormula(metricBindJson.Metric.ID, "ACTIVITY", activities[i].MetricFormula.MetricExpertFormula.ID)
				if errWhenGettingFormulaStatus != nil {
					logrus.Info("Error when getting activity contract status : ", errWhenGettingMetadataContractStatus)
				}

				if formulaStatus == "SUCCESS" || formulaStatus == "QUEUE" {
					//skip this loop and go to next formula
					continue
				} else if formulaStatus == "" || formulaStatus == "FAILED" {
					//insert object for the formula
					ethMetricObjForFormula := model.EthereumMetricBind{
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
						Type:              "ACTIVITY",
						FormulaID:         activities[i].MetricFormula.MetricExpertFormula.ID,
					}
					//handle UUID
					if formulaStatus == "" {
						//generate transaction UUID
						timeNow := time.Now().UTC()
						entropy := rand.New(rand.NewSource(timeNow.UnixNano()))
						id := ulid.MustNew(ulid.Timestamp(timeNow), entropy)
						logrus.Info("TXN UUID : ", id)
						ethMetricObjForFormula.TransactionUUID = id.String()
					} else {
						// get the transaction UUID from the database
						ethMetricObjForFormula.TransactionUUID = formulaDetails.TransactionUUID
					}

					//check the index of the loop to skip the checking of the previous formula deployment
					if i != 0 {
						//check the previous formula contract deployment status
						previousStatus, _, errWhenGettingPreviousStatus := GetMetricSmartContractStatusForFormula(metricBindJson.Metric.ID, "ACTIVITY", activities[i-1].MetricFormula.MetricExpertFormula.ID)
						if errWhenGettingPreviousStatus != nil {
							ethMetricObjForFormula.ErrorMessage = errWhenGettingPreviousStatus.Error()
							ethMetricObjForFormula.Status = "FAILED"
							if formulaStatus == "" {
								//insert to DB
								errWhenInsertingFormulaMetricObj := object.InsertToEthMetricDetails(ethMetricObjForFormula)
								if errWhenInsertingFormulaMetricObj != nil {
									logrus.Info("Error when inserting to metric collection : ", errWhenInsertingFormulaMetricObj)
									commons.JSONErrorReturn(w, r, errWhenInsertingFormulaMetricObj.Error(), 500, "Error when inserting to metric collection : ")
									return
								}
							} else {
								//update collection
								errWhenUpdatingFormulaMetricObj := object.UpdateEthereumMetricStatus(ethMetricObjForFormula.MetricID, ethMetricObjForFormula.TransactionUUID, ethMetricObjForFormula)
								if errWhenUpdatingFormulaMetricObj != nil {
									logrus.Info("Error when updating the metric collection : ", errWhenUpdatingFormulaMetricObj)
									commons.JSONErrorReturn(w, r, errWhenUpdatingFormulaMetricObj.Error(), 500, "Error when updating the metric collection : ")
									return
								}
							}
							logrus.Info("Error when getting previous contract status : ", errWhenGettingPreviousStatus)
							commons.JSONErrorReturn(w, r, errWhenGettingPreviousStatus.Error(), 500, "Error when getting previous contract status : ")
							return
						}
						if previousStatus == "SUCCESS" {
							canCallNextDeployment = true
						} else if previousStatus == "FAILED" || previousStatus == "" || previousStatus == "QUEUE" {
							canCallNextDeployment = false
						}
					}

					if canCallNextDeployment {
						// get the formula map id form DB
						formulaMapID, errWhenGettingFormulaMapId := GetFormulaMapId(activities[i].MetricFormula.MetricExpertFormula.ID)
						if errWhenGettingFormulaMapId != nil {
							ethMetricObjForFormula.ErrorMessage = errWhenGettingFormulaMapId.Error()
							ethMetricObjForFormula.Status = "FAILED"
							if formulaStatus == "" {
								//insert to DB
								errWhenInsertingFormulaMetricObj := object.InsertToEthMetricDetails(ethMetricObjForFormula)
								if errWhenInsertingFormulaMetricObj != nil {
									logrus.Info("Error when inserting to metric collection : ", errWhenInsertingFormulaMetricObj)
									commons.JSONErrorReturn(w, r, errWhenInsertingFormulaMetricObj.Error(), 500, "Error when inserting to metric collection : ")
									return
								}
							} else {
								//update collection
								errWhenUpdatingFormulaMetricObj := object.UpdateEthereumMetricStatus(ethMetricObjForFormula.MetricID, ethMetricObjForFormula.TransactionUUID, ethMetricObjForFormula)
								if errWhenUpdatingFormulaMetricObj != nil {
									logrus.Info("Error when updating the metric collection : ", errWhenUpdatingFormulaMetricObj)
									commons.JSONErrorReturn(w, r, errWhenUpdatingFormulaMetricObj.Error(), 500, "Error when updating the metric collection : ")
									return
								}
							}
							logrus.Info("Error when getting formula map ID : ", errWhenGettingFormulaMapId)
						}
						formulaMapIDString := strconv.FormatUint(formulaMapID, 10)
						activityContractName := "Metric_" + metricMapIDString + "_Formula_" + formulaMapIDString
						ethMetricObjForFormula.ContractName = activityContractName

						errWhenDeployingActivityContract := activityWriters.ActivityContractDeployer(metricMapIDString, formulaMapIDString, metricBindJson.Metric.ID, activities[i], metricBindJson.Metric.Name, metricBindJson.Metric, metricBindJson.User, ethMetricObjForFormula)
						if errWhenDeployingActivityContract != nil {
							ethMetricObjForFormula.ErrorMessage = errWhenDeployingActivityContract.Error()
							ethMetricObjForFormula.Status = "FAILED"
							if formulaStatus == "" {
								//insert to DB
								errWhenInsertingFormulaMetricObj := object.InsertToEthMetricDetails(ethMetricObjForFormula)
								if errWhenInsertingFormulaMetricObj != nil {
									logrus.Info("Error when inserting to metric collection : ", errWhenInsertingFormulaMetricObj)
									commons.JSONErrorReturn(w, r, errWhenInsertingFormulaMetricObj.Error(), 500, "Error when inserting to metric collection : ")
									return
								}
							} else {
								//update collection
								errWhenUpdatingFormulaMetricObj := object.UpdateEthereumMetricStatus(ethMetricObjForFormula.MetricID, ethMetricObjForFormula.TransactionUUID, ethMetricObjForFormula)
								if errWhenUpdatingFormulaMetricObj != nil {
									logrus.Info("Error when updating the metric collection : ", errWhenUpdatingFormulaMetricObj)
									commons.JSONErrorReturn(w, r, errWhenUpdatingFormulaMetricObj.Error(), 500, "Error when updating the metric collection : ")
									return
								}
							}
							logrus.Info("Error when deploying activity contract : ", errWhenDeployingActivityContract)
							commons.JSONErrorReturn(w, r, errWhenDeployingActivityContract.Error(), 500, "Error when deploying activity contract : ")
							return
						}
					}
				}

			}

			//Once all the contracts are deployed
			w.WriteHeader(http.StatusOK)
			response := model.SuccessResponseMetricBinding{
				Code:     http.StatusOK,
				MetricID: ethMetricObjForMetaData.MetricID,
				Message:  "Metric binding request sent to queue",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

}
