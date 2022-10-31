package metricBinding

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	expertformula "github.com/dileepaj/tracified-gateway/protocols/stellarprotocols/expertFormula"
	"github.com/dileepaj/tracified-gateway/services"
	"github.com/oklog/ulid"
	"github.com/stellar/go/txnbuild"

	"github.com/sirupsen/logrus"
)

/*
StellarMetricBinding
des- This method build stellar trasaction for metric binding
	 steps
		* map the metric id and retrieve the mapped id
		* map the tenent id and retrieve the mapped id
		* build memo for the transaction
		* build publisher identity manage data operation
		* build formula definition manage data operation
			* map the activity id and retrieve the mapped id
		* loop through the formulaArray to build the value definition manage data operation
			* map the stage id and retrieve the mapped id
		* put managedata array to rabbitmq server
*/

var manageDataPerMetricBindingRequest int = 25

func StellarMetricBinding(w http.ResponseWriter, r *http.Request, metricBindJson model.MetricDataBindingRequest) {
	metricBinding := MetricBinding{}
	stellarProtocol := expertformula.ExpertFormula{}
	var manageDataOpArray []txnbuild.ManageData
	var metStatus string
	object := dao.Connection{}
	//find -> status -> Queue fail(retry) -> fail pass on -> sucess error and drop
	//get the metric curent status from the DB
	metricMapDetials, errIGettingMetricStatus := object.GetMetricStatus(metricBindJson.Metric.ID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errIGettingMetricStatus != nil {
		logrus.Error("An error occured when getting metric status ", errIGettingMetricStatus)
	}
	if metricMapDetials == nil {
		metStatus = ""
	}
	if metricMapDetials != nil {
		metricMapDet := metricMapDetials.(model.MetricBindingStore)
		metStatus = metricMapDet.Status
		logrus.Info("Status recorded : ", metStatus)
	}

	//check the status of the metric status
	if metStatus == "QUEUE" {
		//ask user to try again
		logrus.Info("Requested metric is in the queue, please try again")
		commons.JSONErrorReturn(w, r, metStatus, 400, "Requested metric is in the queue, please try again")
		return
	} else if metStatus == "FAILED" || metStatus == "" {
		logrus.Info("Requested metric id status is failed or a new binding request")
		//pass on the transaction
		metricBindingStore := model.MetricBindingStore{
			MetricId:  metricBindJson.Metric.ID,
			Metric:    metricBindJson.Metric,
			User:      metricBindJson.User,
			Timestamp: time.Now().String(),
		}
		// mapMetricId uint64, metricName string, tenantId uint32, noOfFormula int32
		metricMapID, errCode, err := InsertAndFindMetricID(metricBindJson.Metric.ID, metricBindJson.Metric.Name)
		if err != nil {
			metricBindingStore.ErrorMessage = err.Error()
			_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
			if errResult != nil {
				logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
			}
			commons.JSONErrorReturn(w, r, err.Error(), errCode, "InsertAndFindMetricID ")
			return
		}
		metricBindingStore.MetricMapID = metricMapID
		tenantMapId, err := InsertAndFindTenentID(metricBindJson.Metric.TenantId)
		if err != nil {
			metricBindingStore.ErrorMessage = err.Error()
			_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
			if errResult != nil {
				logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
			}
			commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "InsertAndFindMetricID ")
			return
		}

		// 1. Metric Name Definition (Compulsory MDO) metric name builder
		metricName, errWhenBuildingMetricName := metricBinding.BuildMetricNameManageData(metricBindJson.Metric.Name, "Metric name")
		if errWhenBuildingMetricName != nil {
			metricBindingStore.ErrorMessage = errWhenBuildingMetricName.Error()
			_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
			if errResult != nil {
				logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
			}
			commons.JSONErrorReturn(w, r, errWhenBuildingMetricName.Error(), http.StatusInternalServerError, "")
			return
		}
		manageDataOpArray = append(manageDataOpArray, metricName)

		// 2. Metric publisher public key definition (Compulsory MDO) 64 byte  hash256 of PGP's public key
		publisherIdentity, err := stellarProtocol.BuildPublicManageData(metricBindJson.User.Publickey)
		if err != nil {
			metricBindingStore.ErrorMessage = err.Error()
			_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
			if errResult != nil {
				logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
			}
			commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "BuildAuthorManageData ")
			return
		}
		manageDataOpArray = append(manageDataOpArray, publisherIdentity)
		// manage data opration order counter
		//! Formula definitions magane data strat
		for i, activity := range metricBindJson.Metric.Activities {
			// checked whether given formulaID already in the database or not
			formulaMapID, err := object.GetFormulaMapID(activity.MetricFormula.MetricExpertFormula.ID).Then(func(data interface{}) interface{} {
				return data
			}).Await()
			if err != nil {
				metricBindingStore.ErrorMessage = err.Error()
				_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
				if errResult != nil {
					logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
				}
				commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Can not find the Formula Id in the database ")
				return
			}
			formulaDetails := formulaMapID.(model.FormulaIDMap)
			activityMapId, err := InsertAndFindActivityID(activity.ID, activity.Name, activity.MetricID, activity.StageID)
			if err != nil {
				metricBindingStore.ErrorMessage = err.Error()
				_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
				if errResult != nil {
					logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
				}
				commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "InsertAndFindMetricID ")
				return
			}
			// 3. Formula definition with metadata (Compulsory MDO for each formula)
			formulaDefinition, _, _, err := metricBinding.BuildFormulaDefinition(formulaDetails.MapID, activityMapId, uint16(len(metricBindJson.Metric.Activities[i].MetricFormula.Formula)))
			if err != nil {
				metricBindingStore.ErrorMessage = err.Error()
				_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
				if errResult != nil {
					logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
				}
				commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "BuildFormulaDefinition ")
				return
			}
			manageDataOpArray = append(manageDataOpArray, formulaDefinition)
			// 4. Activity Name Definition (Compulsory MDO) metric name builder
			activityName, err := metricBinding.BuildMetricNameManageData(activity.Name, "Activity name")
			if err != nil {
				metricBindingStore.ErrorMessage = err.Error()
				_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
				if errResult != nil {
					logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
				}
				commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "")
				return
			}
			manageDataOpArray = append(manageDataOpArray, activityName)
			for j, formula := range activity.MetricFormula.Formula {
				if formula.ArtifactTemplateID == "" {
					bindValue := model.ValueBuilder{
						ValueUUID:           formula.ID,
						WorkflowID:          activity.WorkflowID,
						StageID:             activity.StageID,
						TracabilityDataType: uint(formula.Type),
						BindingType:         1,
					}
					// 5. General value definition
					valueDefinition, _, _, err := metricBinding.BuildGeneralValueManageData(bindValue, activity.MetricFormula.MetricExpertFormula.ID)
					if err != nil {
						metricBindingStore.ErrorMessage = err.Error()
						_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
						if errResult != nil {
							logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
						}
						commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "BuildFormulaDefinition ")
						return
					}
					manageDataOpArray = append(manageDataOpArray, valueDefinition)

					// 6. stage name builder
					stageNameBuilder, errInStageNameBuilder := metricBinding.BuildMetricNameManageData(activity.Stage.Name, "Stage Name")
					if errInStageNameBuilder != nil {
						logrus.Error("Buidling stage name failed ", errInStageNameBuilder.Error())
						metricBindingStore.ErrorMessage = errInStageNameBuilder.Error()
						_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
						if errResult != nil {
							logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
						}
						commons.JSONErrorReturn(w, r, errInStageNameBuilder.Error(), http.StatusInternalServerError, "BuildStageName ")
						return
					}
					manageDataOpArray = append(manageDataOpArray, stageNameBuilder)

					// 7. Key name manage data(workflow → revision number→ stage[] → stage id → traceability data → “key name” )
					keyNameBuilder, errInKeyNameBuilder := metricBinding.BuildMetricNameManageData(formula.Key, "Key Name")
					if errInKeyNameBuilder != nil {
						logrus.Error("Buidling key name failed ", errInKeyNameBuilder.Error())
						metricBindingStore.ErrorMessage = errInKeyNameBuilder.Error()
						_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
						if errResult != nil {
							logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
						}
						commons.JSONErrorReturn(w, r, errInKeyNameBuilder.Error(), http.StatusInternalServerError, "BuildKeyName ")
						return
					}
					manageDataOpArray = append(manageDataOpArray, keyNameBuilder)

				} else {
					// For Master data type Bind
					bindValue := model.ValueBuilder{
						ValueUUID:           formula.ID,
						WorkflowID:          activity.WorkflowID,
						StageID:             activity.StageID,
						TracabilityDataType: uint(formula.Type),
						BindingType:         2,
					}
					// 5. General value definition
					valueDefinition, _, _, err := metricBinding.BuildGeneralValueManageData(bindValue, activity.MetricFormula.MetricExpertFormula.ID)
					if err != nil {
						metricBindingStore.ErrorMessage = err.Error()
						_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
						if errResult != nil {
							logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
						}
						commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "BuildFormulaDefinition ")
						return
					}
					manageDataOpArray = append(manageDataOpArray, valueDefinition)

					// 6. stage name builder
					stageNameBuilder, errInStageNameBuilder := metricBinding.BuildMetricNameManageData(activity.Stage.Name, "Stage Name")
					if errInStageNameBuilder != nil {
						logrus.Error("Buidling stage name failed ", errInStageNameBuilder.Error())
						metricBindingStore.ErrorMessage = errInStageNameBuilder.Error()
						_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
						if errResult != nil {
							logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
						}
						commons.JSONErrorReturn(w, r, errInStageNameBuilder.Error(), http.StatusInternalServerError, "BuildStageName ")
						return
					}
					manageDataOpArray = append(manageDataOpArray, stageNameBuilder)

					// 7. key name builder
					keyNameBuilder, errInKeyNameBuilder := metricBinding.BuildMetricNameManageData(formula.Key, "Key Name")
					if errInKeyNameBuilder != nil {
						logrus.Error("Buidling key name failed ", errInKeyNameBuilder.Error())
						metricBindingStore.ErrorMessage = errInKeyNameBuilder.Error()
						_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
						if errResult != nil {
							logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
						}
						commons.JSONErrorReturn(w, r, errInKeyNameBuilder.Error(), http.StatusInternalServerError, "BuildKeyName ")
						return
					}
					manageDataOpArray = append(manageDataOpArray, keyNameBuilder)

					// Artifact ID Map
					artifactMapId, err := InsertAndFindArtifactID(metricBindJson.Metric.Activities[i].MetricFormula.Formula[j].ArtifactTemplateID)
					if err != nil {
						metricBindingStore.ErrorMessage = err.Error()
						_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
						if errResult != nil {
							logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
						}
						commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "InsertAndFindArtifactID ")
						return
					}
					// General master data info builder
					generalInfoBuilder, errInGeneralInfoBuilder := metricBinding.BuildGeneralMasterDataInfo(artifactMapId, uint(formula.Type))
					if errInGeneralInfoBuilder != nil {
						logrus.Error("Buidling general master data info failed ", errInGeneralInfoBuilder.Error())
						metricBindingStore.ErrorMessage = errInGeneralInfoBuilder.Error()
						_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
						if errResult != nil {
							logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
						}
						commons.JSONErrorReturn(w, r, errInGeneralInfoBuilder.Error(), http.StatusInternalServerError, "BuildGeneralMasterDataInfo ")
						return
					}
					manageDataOpArray = append(manageDataOpArray, generalInfoBuilder)

					// Metadata/ Artifact template name / Table name builder
					ArtifactNameBuilder, errInMetaDataBuilder := metricBinding.BuildMetricNameManageData(formula.ArtifactTemplate.Name, "Artifact template name")
					if errInMetaDataBuilder != nil {
						logrus.Error("Buidling metadata failed ", errInMetaDataBuilder.Error())
						metricBindingStore.ErrorMessage = errInMetaDataBuilder.Error()
						_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
						if errResult != nil {
							logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
						}
						commons.JSONErrorReturn(w, r, errInMetaDataBuilder.Error(), http.StatusInternalServerError, "BuildFormulaDefinition ")
						return
					}
					manageDataOpArray = append(manageDataOpArray, ArtifactNameBuilder)

					// Field key builder ---> Atrifact field key (which column) → key of field array in artifact template
					fieldKeyBuilder, errWhenBuildingFieldKey := metricBinding.BuildMetricNameManageData(formula.Field, "Atrifact field key")
					if errWhenBuildingFieldKey != nil {
						logrus.Error("Buidling field key failed ", errWhenBuildingFieldKey.Error())
						metricBindingStore.ErrorMessage = errWhenBuildingFieldKey.Error()
						_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
						if errResult != nil {
							logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
						}
						commons.JSONErrorReturn(w, r, errWhenBuildingFieldKey.Error(), http.StatusInternalServerError, "BuildFormulaDefinition ")
						return
					}
					manageDataOpArray = append(manageDataOpArray, fieldKeyBuilder)

					// Field name builder  Field name (which column) → name in field array in artifact template
					fieldNameBuilder, errWhenBuildingFieldName := metricBinding.BuildMetricNameManageData(formula.ArtifactTemplate.FieldName, "Atrifact field name")
					if errWhenBuildingFieldName != nil {
						logrus.Error("Buidling field name failed ", errWhenBuildingFieldName.Error())
						metricBindingStore.ErrorMessage = errWhenBuildingFieldName.Error()
						_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
						if errResult != nil {
							logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
						}
						commons.JSONErrorReturn(w, r, errWhenBuildingFieldName.Error(), http.StatusInternalServerError, "BuildFormulaDefinition ")
						return
					}
					manageDataOpArray = append(manageDataOpArray, fieldNameBuilder)
				}
			}
		}
		// split manage data in to 25 length sub arrays
		manageData2dArray := commons.ChunkSlice(manageDataOpArray, manageDataPerMetricBindingRequest)
		// loop the manage data opration2d array and build trasacion
		var memo string
		for i, managedataOperationArray := range manageData2dArray {
			// initial trasacion memo
			if i == 0 {
				memo0, errInMemoBuilder := metricBinding.BuildMemo(0, metricMapID, uint32(tenantMapId), uint16(len(metricBindJson.Metric.Activities)), uint8(len(managedataOperationArray)))
				if errInMemoBuilder != nil {
					metricBindingStore.ErrorMessage = errInMemoBuilder.Error()
					_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
					if errResult != nil {
						logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
					}
					commons.JSONErrorReturn(w, r, errInMemoBuilder.Error(), http.StatusInternalServerError, "")
					return
				}
				memo = memo0
			}

			if i != 0 {
				memo1, errInMemoBuilder := metricBinding.BuildMemo(1, metricMapID, uint32(tenantMapId), uint16(len(metricBindJson.Metric.Activities)), uint8(len(managedataOperationArray)))
				if errInMemoBuilder != nil {
					metricBindingStore.ErrorMessage = errInMemoBuilder.Error()
					_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
					if errResult != nil {
						logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
					}
					commons.JSONErrorReturn(w, r, errInMemoBuilder.Error(), http.StatusInternalServerError, "")
					return
				}
				memo = memo1
			}

			timeNow := time.Now().UTC()
			entropy := rand.New(rand.NewSource(timeNow.UnixNano()))
			id := ulid.MustNew(ulid.Timestamp(timeNow), entropy)
			logrus.Info("TXN UUID : ", id)
			metricBindingStore.TxnUUID = id.String()
			metricBindingStore.MetricMapID = metricMapID
			metricBindingStore.NoOfManageDataInTxn = len(managedataOperationArray)
			metricBindingStore.TotalNoOfManageData = len(manageDataOpArray)
			metricBindingStore.Status = "QUEUE"
			_, errResult := object.InsertMetricBindingFormula(metricBindingStore) //update
			if errResult != nil {
				logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
			}
			buildMetricBind := model.SendToQueue{
				MetricBinding: metricBindingStore,
				Type:          "METRICBIND",
				User:          metricBindJson.User,
				Memo:          []byte(memo),
				Operations:    managedataOperationArray,
			}
			err := services.SendToQueue(buildMetricBind)
			if err != nil {
				metricBindingStore.ErrorMessage = err.Error()
				_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
				if errResult != nil {
					logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
				}
				logrus.Error("Error when submitting managedata to queue  ", err)
				w.WriteHeader(errCode)
				response := model.Error{Code: errCode, Message: "Error when submitting managedata to queue  " + err.Error()}
				json.NewEncoder(w).Encode(response)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		response := model.SuccessResponseMetricBinding{
			Code:     http.StatusOK,
			MetricID: metricBindJson.Metric.ID,
			Message:  "Metric bind request send to queue",
		}
		json.NewEncoder(w).Encode(response)
		return

	} else if metStatus == "SUCCESS" {
		logrus.Info("Metic is already recorded in the blockchain and the gateway DB")
		//response indicating that metric is already recorded
		commons.JSONErrorReturn(w, r, metStatus, 400, "Metic is already recorded in the blockchain and the gateway DB")
		return
	} else {
		logrus.Info("Metric bind status is invalid : ", metStatus)
		commons.JSONErrorReturn(w, r, metStatus, 504, "Metric bind status is invalid, status : ")
		return
	}
}
