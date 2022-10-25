package metricBinding

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
	expertformula "github.com/dileepaj/tracified-gateway/protocols/stellarprotocols/expertFormula"
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
		* load stellar account,build and sing the XDR
		* put XDR to stellar blockchain
*/

var manageDataPerMetricBindingRequest int = 25

func StellarMetricBinding(w http.ResponseWriter, r *http.Request, metricBindJson model.MetricDataBindingRequest) {
	metricBinding := MetricBinding{}
	stellarProtocol := expertformula.ExpertFormula{}
	var manageDataOpArray []txnbuild.Operation
	object := dao.Connection{}
	metricBindingStore := metricBindJson
	metricBindingStore.Metric.Status = "Failed"
	// building memo
	// mapMetricId uint64, metricName string, tenantId uint32, noOfFormula int32
	metricMapID, errCode, err := InsertAndFindMetricID(metricBindJson.Metric.ID, metricBindJson.Metric.Name)
	if err != nil {
		metricBindingStore.Metric.ErrorMessage = err.Error()
		_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
		if errResult != nil {
			logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
		}
		commons.JSONErrorReturn(w, r, err.Error(), errCode, "InsertAndFindMetricID ")
		return
	}
	tenantMapId, err := InsertAndFindTenentID(metricBindJson.Metric.TenantId)
	if err != nil {
		metricBindingStore.Metric.ErrorMessage = err.Error()
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
		metricBindingStore.Metric.ErrorMessage = errWhenBuildingMetricName.Error()
		_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
		if errResult != nil {
			logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
		}
		commons.JSONErrorReturn(w, r, errWhenBuildingMetricName.Error(), http.StatusInternalServerError, "")
		return
	}
	manageDataOpArray = append(manageDataOpArray, &metricName)

	// 2. Metric publisher public key definition (Compulsory MDO) 64 byte  hash256 of PGP's public key
	publisherIdentity, err := stellarProtocol.BuildPublicManageData(metricBindJson.User.Publickey)
	if err != nil {
		metricBindingStore.Metric.ErrorMessage = err.Error()
		_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
		if errResult != nil {
			logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
		}
		commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "BuildAuthorManageData ")
		return
	}
	manageDataOpArray = append(manageDataOpArray, &publisherIdentity)
	// manage data opration order counter
	//! Formula definitions magane data strat
	c := 2
	for i, activity := range metricBindJson.Metric.Activities {
		// checked whether given formulaID already in the database or not
		formulaMapID, err := object.GetFormulaMapID(activity.MetricFormula.MetricExpertFormula.ID).Then(func(data interface{}) interface{} {
			return data
		}).Await()
		if err != nil {
			metricBindingStore.Metric.ErrorMessage = err.Error()
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
			metricBindingStore.Metric.ErrorMessage = err.Error()
			_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
			if errResult != nil {
				logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
			}
			commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "InsertAndFindMetricID ")
			return
		}
		// 3. Formula definition with metadata (Compulsory MDO for each formula)
		formulaDefinition, keyFD, valueFD, err := metricBinding.BuildFormulaDefinition(formulaDetails.MapID, activityMapId, uint16(len(metricBindJson.Metric.Activities[i].MetricFormula.Formula)))
		if err != nil {
			metricBindingStore.Metric.ErrorMessage = err.Error()
			_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
			if errResult != nil {
				logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
			}
			commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "BuildFormulaDefinition ")
			return
		}
		c++
		metricBindingStore.Metric.Activities[i].ActivityFormulaDefinitionManageData = model.ActivityFormulaDefinitionManageData{
			ManageDataOrder:     c,
			ManageDataType:      "FORMULA DEFINITION",
			FormulaMapID:        formulaDetails.MapID,
			FormulaID:           activity.MetricFormula.MetricExpertFormula.ID,
			ActivityMapID:       activityMapId,
			ActivityID:          activity.ID,
			NoOfDynamicVariable: uint32(len(activity.MetricFormula.Formula)),
			Key:                 keyFD,
			Value:               valueFD,
		}
		manageDataOpArray = append(manageDataOpArray, &formulaDefinition)
		// 4. Activity Name Definition (Compulsory MDO) metric name builder
		activityName, err := metricBinding.BuildMetricNameManageData(activity.Name, "Activity name")
		if err != nil {
			metricBindingStore.Metric.ErrorMessage = err.Error()
			_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
			if errResult != nil {
				logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
			}
			commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "")
			return
		}
		c++
		metricBindingStore.Metric.Activities[i].ActivityNameMangeData = model.ManageDataActivityName{
			ManageDataOrder: c,
			ManageDataType:  "MATRIC NAME",
			ActivityName:    metricBindJson.Metric.Name,
			Name:            activityName.Name,
			Value:           activityName.Value,
		}
		manageDataOpArray = append(manageDataOpArray, &activityName)
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
				valueDefinition, keyVD, valueVD, err := metricBinding.BuildGeneralValueManageData(bindValue, activity.MetricFormula.MetricExpertFormula.ID)
				if err != nil {
					metricBindingStore.Metric.ErrorMessage = err.Error()
					_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
					if errResult != nil {
						logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
					}
					commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "BuildFormulaDefinition ")
					return
				}
				c++
				metricBindingStore.Metric.Activities[i].MetricFormula.Formula[j].BindManageData = model.BindManageData{
					ManageDataOrder: c,
					ManageDataType:  "STAGE BIND DATA",
					BindData:        bindValue,
					Key:             keyVD,
					Value:           valueVD,
				}
				manageDataOpArray = append(manageDataOpArray, &valueDefinition)

				// 6. stage name builder
				stageNameBuilder, errInStageNameBuilder := metricBinding.BuildMetricNameManageData(activity.Stage.Name, "Stage Name")
				if errInStageNameBuilder != nil {
					logrus.Error("Buidling stage name failed ", errInStageNameBuilder.Error())
					metricBindingStore.Metric.ErrorMessage = errInStageNameBuilder.Error()
					_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
					if errResult != nil {
						logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
					}
					commons.JSONErrorReturn(w, r, errInStageNameBuilder.Error(), http.StatusInternalServerError, "BuildStageName ")
					return
				}
				c++
				metricBindingStore.Metric.Activities[i].MetricFormula.Formula[j].BindManageData.Stage = model.ManageDataStageName{
					ManageDataOrder: c,
					ManageDataType:  "SATAGE NAME",
					StageName:       activity.Stage.Name,
					Name:            stageNameBuilder.Name,
					Value:           stageNameBuilder.Value,
				}
				manageDataOpArray = append(manageDataOpArray, &stageNameBuilder)

				// 7. Key name manage data(workflow → revision number→ stage[] → stage id → traceability data → “key name” )
				keyNameBuilder, errInKeyNameBuilder := metricBinding.BuildMetricNameManageData(formula.Key, "Key Name")
				if errInKeyNameBuilder != nil {
					logrus.Error("Buidling key name failed ", errInKeyNameBuilder.Error())
					metricBindingStore.Metric.ErrorMessage = errInKeyNameBuilder.Error()
					_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
					if errResult != nil {
						logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
					}
					commons.JSONErrorReturn(w, r, errInKeyNameBuilder.Error(), http.StatusInternalServerError, "BuildKeyName ")
					return
				}
				c++
				metricBindingStore.Metric.Activities[i].MetricFormula.Formula[j].BindManageData.KeyName = model.ManageDataKeyName{
					ManageDataOrder: c,
					ManageDataType:  "KEY NAME (WORKFLOW)",
					KeyName:         formula.Key,
					Name:            keyNameBuilder.Name,
					Value:           keyNameBuilder.Value,
				}
				manageDataOpArray = append(manageDataOpArray, &keyNameBuilder)

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
				valueDefinition, keyVD, valueVD, err := metricBinding.BuildGeneralValueManageData(bindValue, activity.MetricFormula.MetricExpertFormula.ID)
				if err != nil {
					metricBindingStore.Metric.ErrorMessage = err.Error()
					_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
					if errResult != nil {
						logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
					}
					commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "BuildFormulaDefinition ")
					return
				}
				c++
				metricBindingStore.Metric.Activities[i].MetricFormula.Formula[j].BindManageData = model.BindManageData{
					ManageDataOrder: c,
					ManageDataType:  "MASTER BIND DATA",
					BindData:        bindValue,
					Key:             keyVD,
					Value:           valueVD,
				}
				manageDataOpArray = append(manageDataOpArray, &valueDefinition)

				// 6. stage name builder
				stageNameBuilder, errInStageNameBuilder := metricBinding.BuildMetricNameManageData(activity.Stage.Name, "Stage Name")
				if errInStageNameBuilder != nil {
					logrus.Error("Buidling stage name failed ", errInStageNameBuilder.Error())
					metricBindingStore.Metric.ErrorMessage = errInStageNameBuilder.Error()
					_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
					if errResult != nil {
						logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
					}
					commons.JSONErrorReturn(w, r, errInStageNameBuilder.Error(), http.StatusInternalServerError, "BuildStageName ")
					return
				}
				c++
				metricBindingStore.Metric.Activities[i].MetricFormula.Formula[j].BindManageData.Stage = model.ManageDataStageName{
					ManageDataOrder: c,
					ManageDataType:  "SATAGE NAME",
					StageName:       activity.Stage.Name,
					Name:            stageNameBuilder.Name,
					Value:           stageNameBuilder.Value,
				}
				manageDataOpArray = append(manageDataOpArray, &stageNameBuilder)

				// 7. key name builder
				keyNameBuilder, errInKeyNameBuilder := metricBinding.BuildMetricNameManageData(formula.Key, "Key Name")
				if errInKeyNameBuilder != nil {
					logrus.Error("Buidling key name failed ", errInKeyNameBuilder.Error())
					metricBindingStore.Metric.ErrorMessage = errInKeyNameBuilder.Error()
					_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
					if errResult != nil {
						logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
					}
					commons.JSONErrorReturn(w, r, errInKeyNameBuilder.Error(), http.StatusInternalServerError, "BuildKeyName ")
					return
				}
				c++
				metricBindingStore.Metric.Activities[i].MetricFormula.Formula[j].BindManageData.KeyName = model.ManageDataKeyName{
					ManageDataOrder: c,
					ManageDataType:  "KEY NAME (WORKFLOW)",
					KeyName:         formula.Key,
					Name:            keyNameBuilder.Name,
					Value:           keyNameBuilder.Value,
				}
				manageDataOpArray = append(manageDataOpArray, &keyNameBuilder)

				// Artifact ID Map
				artifactMapId, err := InsertAndFindArtifactID(metricBindJson.Metric.Activities[i].MetricFormula.Formula[j].ArtifactTemplateID)
				if err != nil {
					metricBindingStore.Metric.ErrorMessage = err.Error()
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
					metricBindingStore.Metric.ErrorMessage = errInGeneralInfoBuilder.Error()
					_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
					if errResult != nil {
						logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
					}
					commons.JSONErrorReturn(w, r, errInGeneralInfoBuilder.Error(), http.StatusInternalServerError, "BuildGeneralMasterDataInfo ")
					return
				}
				c++
				metricBindingStore.Metric.Activities[i].MetricFormula.Formula[j].BindManageData.Master = model.ManageDataMasterDetails{
					ManageDataOrder:      c,
					ManageDataType:       "GENERAL MASTER DATA INFO",
					ArtifactID:           formula.ArtifactTemplateID,
					TracabilityDataType:  uint(formula.Type),
					Name:                 generalInfoBuilder.Name,
					Value:                generalInfoBuilder.Value,
					ArtifactTemplateName: model.ManageDataArtifactTemplateName{},
					ArtifactFieldKey:     model.ManageDataFieldKey{},
					ArtifactFieldName:    model.ManageDataFieldName{},
				}
				manageDataOpArray = append(manageDataOpArray, &generalInfoBuilder)

				// Metadata/ Artifact template name / Table name builder
				ArtifactNameBuilder, errInMetaDataBuilder := metricBinding.BuildMetricNameManageData(formula.ArtifactTemplate.Name, "Artifact template name")
				if errInMetaDataBuilder != nil {
					logrus.Error("Buidling metadata failed ", errInMetaDataBuilder.Error())
					metricBindingStore.Metric.ErrorMessage = errInMetaDataBuilder.Error()
					_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
					if errResult != nil {
						logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
					}
					commons.JSONErrorReturn(w, r, errInMetaDataBuilder.Error(), http.StatusInternalServerError, "BuildFormulaDefinition ")
					return
				}
				c++
				metricBindingStore.Metric.Activities[i].MetricFormula.Formula[j].BindManageData.Master = model.ManageDataMasterDetails{
					ManageDataOrder: c,
					ManageDataType:  "ARTIFACT TEMPLATE NAME",
					ArtifactID:      formula.ArtifactTemplateID,
					Name:            ArtifactNameBuilder.Name,
					Value:           ArtifactNameBuilder.Value,
				}
				manageDataOpArray = append(manageDataOpArray, &ArtifactNameBuilder)

				// Field key builder ---> Atrifact field key (which column) → key of field array in artifact template
				fieldKeyBuilder, errWhenBuildingFieldKey := metricBinding.BuildMetricNameManageData(formula.Field, "Atrifact field key")
				if errWhenBuildingFieldKey != nil {
					logrus.Error("Buidling field key failed ", errWhenBuildingFieldKey.Error())
					metricBindingStore.Metric.ErrorMessage = errWhenBuildingFieldKey.Error()
					_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
					if errResult != nil {
						logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
					}
					commons.JSONErrorReturn(w, r, errWhenBuildingFieldKey.Error(), http.StatusInternalServerError, "BuildFormulaDefinition ")
					return
				}
				c++
				metricBindingStore.Metric.Activities[i].MetricFormula.Formula[j].BindManageData.Master.ArtifactFieldKey = model.ManageDataFieldKey{
					ManageDataOrder: c,
					ManageDataType:  "ARTIFACT TEMPLATE FIELD Key",
					FieldKey:        formula.Field,
					Name:            fieldKeyBuilder.Name,
					Value:           fieldKeyBuilder.Value,
				}
				manageDataOpArray = append(manageDataOpArray, &fieldKeyBuilder)

				// Field name builder  Field name (which column) → name in field array in artifact template
				fieldNameBuilder, errWhenBuildingFieldName := metricBinding.BuildMetricNameManageData(formula.ArtifactTemplate.FieldName, "Atrifact field name")
				if errWhenBuildingFieldName != nil {
					logrus.Error("Buidling field name failed ", errWhenBuildingFieldName.Error())
					metricBindingStore.Metric.ErrorMessage = errWhenBuildingFieldName.Error()
					_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
					if errResult != nil {
						logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
					}
					commons.JSONErrorReturn(w, r, errWhenBuildingFieldName.Error(), http.StatusInternalServerError, "BuildFormulaDefinition ")
					return
				}
				c++
				metricBindingStore.Metric.Activities[i].MetricFormula.Formula[j].BindManageData.Master.ArtifactFieldKey = model.ManageDataFieldKey{
					ManageDataOrder: c,
					ManageDataType:  "ARTIFACT TEMPLATE FIELD NAME",
					FieldKey:        formula.ArtifactTemplate.FieldName,
					Name:            fieldNameBuilder.Name,
					Value:           fieldNameBuilder.Value,
				}
				manageDataOpArray = append(manageDataOpArray, &fieldNameBuilder)
			}
		}
	}
	metricBindingStore.Metric.Transactions = model.TransacionDetailsMetricBinding{
		MetricID:    metricBindJson.Metric.ID,
		MtericMapId: metricMapID,
		TenantMapId: uint32(tenantMapId),
		MetricName:  metricBindJson.Metric.Name,
		MetricNameMangeData: model.ManageDataMetricName{
			ManageDataOrder: 1,
			ManageDataType:  "MEREIC NAME",
			MetricName:      metricBindJson.Metric.Name,
			Name:            metricName.Name,
			Value:           metricName.Value,
		},
		NoOfActivityFormula:         len(metricBindJson.Metric.Activities),
		PublisherIdentityManageData: model.PublisherIdentity{ManageDataOrder: 2, ManageDataType: "PUBLISHER IDENTITY", UserID: metricBindJson.User.ID, PublicKeyHash: metricBindJson.User.Publickey, ManageDataKey: publisherIdentity.Name, ManageDataValue: publisherIdentity.Value},
		TotalNumberOfManageData:     len(manageDataOpArray),
		MaxNumOfManageDatePerHash:   manageDataPerMetricBindingRequest,
	}
	// split manage data in to 25 length sub arrays
	manageData2dArray := commons.ChunkSlice(manageDataOpArray, manageDataPerMetricBindingRequest)
	var errorInXDR error
	// loop the manage data opration2d array and build trasacion
	for i, manadataOperationArray := range manageData2dArray {
		// initial trasacion memo
		memo0, errInMemoBuilder := metricBinding.BuildMemo(0, metricMapID, uint32(tenantMapId), uint16(len(metricBindJson.Metric.Activities)), uint8(len(manageDataOpArray)))
		if errInMemoBuilder != nil {
			metricBindingStore.Metric.ErrorMessage = errInMemoBuilder.Error()
			_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
			if errResult != nil {
				logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
			}
			commons.JSONErrorReturn(w, r, errInMemoBuilder.Error(), http.StatusInternalServerError, "")
			return
		}
		stellarprotocol := stellarprotocols.StellarTrasaction{PublicKey: constants.PublicKey, SecretKey: constants.SecretKey, Operations: manadataOperationArray, Memo: memo0}
		transaction := model.TransactionHash{Order: i, Memo: []byte(memo0)}
		if i != 0 {
			memo1, errInMemoBuilder := metricBinding.BuildMemo(1, metricMapID, uint32(tenantMapId), uint16(len(metricBindJson.Metric.Activities)), uint8(len(manadataOperationArray)))
			if errInMemoBuilder != nil {
				metricBindingStore.Metric.ErrorMessage = errInMemoBuilder.Error()
				_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
				if errResult != nil {
					logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
				}
				commons.JSONErrorReturn(w, r, errInMemoBuilder.Error(), http.StatusInternalServerError, "")
				return
			}
			stellarprotocol.Memo = memo1
			transaction.Memo = []byte(memo1)
		}
		err, errCode, hash := stellarprotocol.SubmitToStellerBlockchain()
		errorInXDR = err
		if err != nil {
			// status = "Failed"
			metricBindingStore.Metric.ErrorMessage = err.Error()
			_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
			if errResult != nil {
				logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
			}
			logrus.Error("Error when submitting transaction to blockchain  ", err)
			w.WriteHeader(errCode)
			response := model.Error{Code: errCode, Message: "Error when submitting transaction to blockchain  " + err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}
		transaction.Hash = hash
		metricBindingStore.Metric.Transactions.TXNHashes = append(metricBindingStore.Metric.Transactions.TXNHashes, transaction)
		logrus.Info("Transaction Hash ", hash)
	}
	if errorInXDR != nil {
		metricBindingStore.Metric.ErrorMessage = err.Error()
		_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
		if errResult != nil {
			logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
		}
		commons.JSONErrorReturn(w, r, errorInXDR.Error(), http.StatusInternalServerError, "Error when submitting transaction to blockchain ")
		return
	} else {
		metricBindingStore.Metric.Status = "SAVED IN BLOCKCHAIN"
		id, errResult := object.InsertMetricBindingFormula(metricBindingStore)
		if errResult != nil {
			logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
		}
		w.WriteHeader(http.StatusOK)
		response := model.SuccessResponseMetricBinding{
			Code:              http.StatusOK,
			ID:                id,
			MetricID:          metricBindJson.Metric.ID,
			TransactionHashes: metricBindingStore.Metric.Transactions.TXNHashes,
		}
		json.NewEncoder(w).Encode(response)
		return
	}
}
