package ethereuemmetricbind

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	activityWriters "github.com/dileepaj/tracified-gateway/expertformulas/ethereum/codeGenerator/ethereuemMetricBind/ActivityContractWriters"
	metadataWriters "github.com/dileepaj/tracified-gateway/expertformulas/ethereum/codeGenerator/ethereuemMetricBind/metadataWriters"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/services/ethereumServices/dbCollectionHandler"
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
		Timestamp:         time.Now().UTC().String(),
		ContractAddress:   "",
		TransactionHash:   "",
		TransactionCost:   "",
		TransactionUUID:   "",
		TransactionSender: commons.GoDotEnvVariable("ETHEREUMPUBKEY"),
		User:              metricBindJson.User,
		ErrorMessage:      "",
		Status:            116,
		Type:              "METADATA",
		FormulaID:         "",
		ActualStatus: 	   101,	// SMART_CONTRACT_GENERATION_STARTED
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

	if status == 118 {
		logrus.Info("Contract for metric " + metricBindJson.Metric.Name + " has been added to the blockchain, checking the activity array")
		canCallNextDeployment = true
	} else if status == 116 {
		logrus.Info("Requested metric is in the queue, please try again")
		canCallNextDeployment = true
	} else if status == 0 || status == 119 {
		if status == 119 {
			logrus.Info("Requested metric is in the failed status, trying to redeploy")
			ethMetricObjForMetaData.Status = 119
		} else {
			logrus.Info("New metric bind request, initiating new deployment")
		}
	}

	if status == 0 {
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

	if status == 0 || status == 119 {
		if status == 0 { // NIL
			// store the metric object in the database
			errWhenStoringMetricObj := object.InsertToEthMetricDetails(ethMetricObjForMetaData)
			if errWhenStoringMetricObj != nil {
				logrus.Info("Error when inserting to metric collection : ", errWhenStoringMetricObj)
				commons.JSONErrorReturn(w, r, errWhenStoringMetricObj.Error(), 500, "Error when inserting to metric collection : ")
				return
			}
		}
		// deploy the smart contract for meta data
		errWhenDeployingMetaDataSmartContract := metadataWriters.MetricMetadataContractDeployer(metaDataObj, metricMapIDString, ethMetricObjForMetaData)
		if errWhenDeployingMetaDataSmartContract != nil {
			ethMetricObjForMetaData.ErrorMessage = errWhenDeployingMetaDataSmartContract.Error()
			ethMetricObjForMetaData.Status = 119 // FAILED
			ethMetricObjForMetaData.ActualStatus = 111	// DEPLOYMENT_FAILED
			// update the metric object in the database
			errWhenUpdatingMetricObj := object.UpdateSelectedEthMetricFields(ethMetricObjForMetaData.MetricID, ethMetricObjForMetaData.TransactionUUID, ethMetricObjForMetaData)
			if errWhenUpdatingMetricObj != nil {
				logrus.Info("Error when updating the metric collection : ", errWhenUpdatingMetricObj)
				commons.JSONErrorReturn(w, r, errWhenUpdatingMetricObj.Error(), 500, "Error when updating the metric collection : ")
				return
			}

			logrus.Info("Error when deploying metadata metric contract : ", errWhenDeployingMetaDataSmartContract)
			commons.JSONErrorReturn(w, r, errWhenDeployingMetaDataSmartContract.Error(), 500, "Error when deploying metadata metric contract : ")
			return
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
	if status == 118 || status == 117 || status == 116 {
		canCallNextDeployment = true
		if len(activities) > 0 {
			for i := 0; i < len(activities); i++ {
				//check if the contract for the this metric ID + formula ID + type deployed
				formulaStatus, formulaDetails, errWhenGettingFormulaStatus := GetMetricSmartContractStatusForFormula(metricBindJson.Metric.ID, "ACTIVITY", activities[i].MetricFormula.MetricExpertFormula.ID)
				if errWhenGettingFormulaStatus != nil {
					logrus.Info("Error when getting activity contract status : ", errWhenGettingMetadataContractStatus)
				}

				if formulaStatus == 118 || formulaStatus == 116 || formulaStatus == 117 {
					//skip this loop and go to next formula
					logrus.Info("Contract for formula already deployed or in queue : ", activities[i].MetricFormula.MetricExpertFormula.ID)
					continue
				} else if formulaStatus == 0 || formulaStatus == 119 {
					logrus.Info("New or failed activity contract deployment. Trying to deploying contract for formula : ", activities[i].MetricFormula.MetricExpertFormula.ID)
					//insert object for the formula
					ethMetricObjForFormula := model.EthereumMetricBind{
						MetricID:          metricBindJson.Metric.ID,
						MetricName:        metricBindJson.Metric.Name,
						Metric:            metricBindJson.Metric,
						ContractName:      "",
						TemplateString:    "",
						BINstring:         "",
						ABIstring:         "",
						Timestamp:         time.Now().UTC().String(),
						ContractAddress:   "",
						TransactionHash:   "",
						TransactionCost:   "",
						TransactionUUID:   "",
						TransactionSender: commons.GoDotEnvVariable("ETHEREUMPUBKEY"),
						User:              metricBindJson.User,
						ErrorMessage:      "",
						Status:            116,
						Type:              "ACTIVITY",
						FormulaID:         activities[i].MetricFormula.MetricExpertFormula.ID,
						ActualStatus:      101,	// SMART_CONTRACT_GENERATION_STARTED
					}
					//handle UUID
					if formulaStatus == 0 {
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

					//insert to activity contract details to the collection
					if formulaStatus == 119 {
						// set the status to FAILED
						ethMetricObjForFormula.Status = 119
					} else if formulaStatus == 0 {
						// store the metric object in the database
						errWhenStoringMetricObj := object.InsertToEthMetricDetails(ethMetricObjForFormula)
						if errWhenStoringMetricObj != nil {
							logrus.Info("Error when inserting to metric collection : ", errWhenStoringMetricObj)
							commons.JSONErrorReturn(w, r, errWhenStoringMetricObj.Error(), 500, "Error when inserting to metric collection : ")
							return
						}
					}

					//check the index of the loop to skip the checking of the previous formula deployment
					if i != 0 {
						//check the previous formula contract deployment status
						previousStatus, _, errWhenGettingPreviousStatus := GetMetricSmartContractStatusForFormula(metricBindJson.Metric.ID, "ACTIVITY", activities[i-1].MetricFormula.MetricExpertFormula.ID)
						if errWhenGettingPreviousStatus != nil {
							ethMetricObjForFormula.ErrorMessage = errWhenGettingPreviousStatus.Error()
							ethMetricObjForFormula.Status = 119	// FAILED
							ethMetricObjForFormula.ActualStatus = 102	// SMART_CONTRACT_GENERATION_FAILED
							if formulaStatus == 0 {
								//update collection
								errWhenUpdatingFormulaMetricObj := object.UpdateEthereumMetricStatus(ethMetricObjForFormula.MetricID, ethMetricObjForFormula.TransactionUUID, ethMetricObjForFormula)
								if errWhenUpdatingFormulaMetricObj != nil {
									logrus.Info("Error when updating the metric collection : ", errWhenUpdatingFormulaMetricObj)
									commons.JSONErrorReturn(w, r, errWhenUpdatingFormulaMetricObj.Error(), 500, "Error when updating the metric collection : ")
									return
								}
							}
							// update collections and invalidating
							pendingTransaction := model.PendingContracts{
								TransactionHash: "",
								ContractAddress: "",
								Status         : 119,	// FAILED
								CurrentIndex   : 0,
								ErrorMessage   : ethMetricObjForFormula.ErrorMessage,
								ContractType   : "ETHMETRICBIND",
								Identifier     : ethMetricObjForFormula.TransactionUUID,							
							}
							errorWheninvalidating := dbCollectionHandler.InvalidateMetric(pendingTransaction, ethMetricObjForFormula.Status, ethMetricObjForFormula.TransactionUUID)
							if errorWheninvalidating != nil {
								logrus.Info("Error when invalidating the metric collection : ", errorWheninvalidating)
								commons.JSONErrorReturn(w, r, errorWheninvalidating.Error(), 500, "Error when invalidating the metric collection : ")
								return
							}
							logrus.Info("Error when getting previous contract status : ", errWhenGettingPreviousStatus)
							commons.JSONErrorReturn(w, r, errWhenGettingPreviousStatus.Error(), 500, "Error when getting previous contract status : ")
							return
						}
						if previousStatus == 118 || previousStatus == 117 {
							canCallNextDeployment = true
						} else if previousStatus == 119 || previousStatus == 0 || previousStatus == 120 {
							canCallNextDeployment = false
						}
					}

					if canCallNextDeployment {
						// get the formula map id form DB
						formulaMapID, errWhenGettingFormulaMapId := GetFormulaMapId(activities[i].MetricFormula.MetricExpertFormula.ID)
						if errWhenGettingFormulaMapId != nil {
							ethMetricObjForFormula.ErrorMessage = errWhenGettingFormulaMapId.Error()
							ethMetricObjForFormula.Status = 119	// FAILED
							ethMetricObjForFormula.ActualStatus = 102	// SMART_CONTRACT_GENERATION_FAILED
							//update collection
							errWhenUpdatingFormulaMetricObj := object.UpdateEthereumMetricStatus(ethMetricObjForFormula.MetricID, ethMetricObjForFormula.TransactionUUID, ethMetricObjForFormula)
							if errWhenUpdatingFormulaMetricObj != nil {
								logrus.Info("Error when updating the metric collection : ", errWhenUpdatingFormulaMetricObj)
								commons.JSONErrorReturn(w, r, errWhenUpdatingFormulaMetricObj.Error(), 500, "Error when updating the metric collection : ")
								return
							}
							// update collections and invalidating
							pendingTransaction := model.PendingContracts{
								TransactionHash: "",
								ContractAddress: "",
								Status         : 119,	// FAILED
								CurrentIndex   : 0,
								ErrorMessage   : ethMetricObjForFormula.ErrorMessage,
								ContractType   : "ETHMETRICBIND",
								Identifier     : ethMetricObjForFormula.TransactionUUID,							
							}
							errorWheninvalidating := dbCollectionHandler.InvalidateMetric(pendingTransaction, ethMetricObjForFormula.Status, ethMetricObjForFormula.TransactionUUID)
							if errorWheninvalidating != nil {
								logrus.Info("Error when invalidating the metric collection : ", errorWheninvalidating)
								commons.JSONErrorReturn(w, r, errorWheninvalidating.Error(), 500, "Error when invalidating the metric collection : ")
								return
							}
							logrus.Info("Error when getting formula map ID : ", errWhenGettingFormulaMapId)
						}
						formulaMapIDString := strconv.FormatUint(formulaMapID, 10)
						activityContractName := "Metric_" + metricMapIDString + "_Formula_" + formulaMapIDString
						ethMetricObjForFormula.ContractName = activityContractName

						errWhenDeployingActivityContract := activityWriters.ActivityContractDeployer(metricMapIDString, formulaMapIDString, metricBindJson.Metric.ID, activities[i], metricBindJson.Metric.Name, metricBindJson.Metric, metricBindJson.User, ethMetricObjForFormula)
						if errWhenDeployingActivityContract != nil {
							ethMetricObjForFormula.ErrorMessage = errWhenDeployingActivityContract.Error()
							ethMetricObjForFormula.Status = 119	// FAILED
							//update collection
							errWhenUpdatingFormulaMetricObj := object.UpdateEthereumMetricStatus(ethMetricObjForFormula.MetricID, ethMetricObjForFormula.TransactionUUID, ethMetricObjForFormula)
							if errWhenUpdatingFormulaMetricObj != nil {
								logrus.Info("Error when updating the metric collection : ", errWhenUpdatingFormulaMetricObj)
								commons.JSONErrorReturn(w, r, errWhenUpdatingFormulaMetricObj.Error(), 500, "Error when updating the metric collection : ")
								return
							}
							// update collections and invalidating
							pendingTransaction := model.PendingContracts{
								TransactionHash: "",
								ContractAddress: "",
								Status         : 119,	// FAILED
								CurrentIndex   : 0,
								ErrorMessage   : ethMetricObjForFormula.ErrorMessage,
								ContractType   : "ETHMETRICBIND",
								Identifier     : ethMetricObjForFormula.TransactionUUID,							
							}
							errorWheninvalidating := dbCollectionHandler.InvalidateMetric(pendingTransaction, ethMetricObjForFormula.Status, ethMetricObjForFormula.TransactionUUID)
							if errorWheninvalidating != nil {
								logrus.Info("Error when invalidating the metric collection : ", errorWheninvalidating)
								commons.JSONErrorReturn(w, r, errorWheninvalidating.Error(), 500, "Error when invalidating the metric collection : ")
								return
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
