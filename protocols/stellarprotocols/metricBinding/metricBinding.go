package metricBinding

import (
	b64 "encoding/base64"
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

func StellarMetricBinding(w http.ResponseWriter, r *http.Request, metricBindJson model.MetricDataBindingRequest) {
	metricBinding := MetricBinding{}
	stellarProtocol := expertformula.ExpertFormula{}
	var manageDataOpArray []txnbuild.Operation
	object := dao.Connection{}
	metricBindingStore := metricBindJson
	metricBindingStore.Status = "Failed"
	// building memo
	// mapMetricId uint64, metricName string, tenantId uint32, noOfFormula int32
	metricMapID, errCode, err := InsertAndFindMetricID(metricBindJson.ID, metricBindJson.Name)
	if err != nil {
		metricBindingStore.ErrorMessage = err.Error()
		_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
		if errResult != nil {
			logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
		}
		commons.JSONErrorReturn(w, r, err.Error(), errCode, "InsertAndFindMetricID ")
		return
	}
	tenantMapId, err := InsertAndFindTenentID(metricBindJson.TenantId)
	if err != nil {
		metricBindingStore.ErrorMessage = err.Error()
		_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
		if errResult != nil {
			logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
		}
		commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "InsertAndFindMetricID ")
		return
	}
	memo, errInMemoBuilder := metricBinding.BuildMemo(metricMapID, uint32(tenantMapId), uint16(len(metricBindJson.Activities)), uint16(236))
	if errInMemoBuilder != nil {
		metricBindingStore.ErrorMessage = errInMemoBuilder.Error()
		_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
		if errResult != nil {
			logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
		}
		commons.JSONErrorReturn(w, r, errInMemoBuilder.Error(), http.StatusInternalServerError, "")
		return
	}
	// metric name builder
	// base64 encode the metric name to overcome stellar UTE-8 issues
	metricNameB64 := b64.StdEncoding.EncodeToString([]byte(metricBindJson.Name))
	metricName, errWhenBuildingMetricName := metricBinding.BuildMetricName(metricNameB64)
	if errWhenBuildingMetricName != nil {
		metricBindingStore.ErrorMessage = errWhenBuildingMetricName.Error()
		_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
		if errResult != nil {
			logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
		}
		commons.JSONErrorReturn(w, r, errWhenBuildingMetricName.Error(), http.StatusInternalServerError, "")
		return
	}
	manageDataOpArray = append(manageDataOpArray, &metricName)

	publisherIdentity, err := stellarProtocol.BuildPublisherManageData(metricBindJson.User.Publickey)
	if err != nil {
		metricBindingStore.ErrorMessage = err.Error()
		_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
		if errResult != nil {
			logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
		}
		commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "BuildAuthorManageData ")
		return
	}
	manageDataOpArray = append(manageDataOpArray, &publisherIdentity)
	// manage data opration order counter
	c := 1
	for i, activity := range metricBindJson.Activities {
		// stageID, err := strconv.Atoi(activity.StageID)
		if err != nil {
			metricBindingStore.ErrorMessage = err.Error()
			_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
			if errResult != nil {
				logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
			}
			commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Stage ID not a int string ")
			return
		}
		// checked whether given formulaID already in the database or not
		formulaMapID, err := object.GetFormulaMapID(activity.MetricFormula.MetricExpertFormula.ID).Then(func(data interface{}) interface{} {
			return data
		}).Await()
		formulaDetails := formulaMapID.(model.FormulaIDMap)
		if err != nil {
			metricBindingStore.ErrorMessage = err.Error()
			_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
			if errResult != nil {
				logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
			}
			commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Can not find the Formula in database ")
			return
		}
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
		formulaDefinition, keyFD, valueFD, err := metricBinding.BuildFormulaDefinition(formulaDetails.MapID, activityMapId, uint16(len(activity.MetricFormula.Formula)))
		if err != nil {
			metricBindingStore.ErrorMessage = err.Error()
			_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
			if errResult != nil {
				logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
			}
			commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "BuildFormulaDefinition ")
			return
		}
		c++
		metricBindingStore.Activities[i].ActivityFormulaDefinitionManageData = model.ActivityFormulaDefinitionManageData{
			ManageDataOrder:     c,
			ManageDataType:      "FORMULA DEFINITION",
			ActivityMapID:       activityMapId,
			NoOfDynamicVariable: uint32(len(activity.MetricFormula.Formula)),
			Key:                 keyFD,
			Value:               valueFD,
		}
		manageDataOpArray = append(manageDataOpArray, &formulaDefinition)
		for j, formula := range activity.MetricFormula.Formula {
			if formula.ArtifactTemplateID == "" {
				bindValue := model.ValueBuilder{
					ValueUUID:           formula.ID,
					WorkflowID:          "test",
					StageID:             activity.StageID,
					TracabilityDataType: 1,
					BindingType:         1,
				}
				valueDefinition, keyVD, valueVD, err := metricBinding.BuildGeneralValueManageData(bindValue)
				if err != nil {
					metricBindingStore.ErrorMessage = err.Error()
					_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
					if errResult != nil {
						logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
					}
					commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "BuildFormulaDefinition ")
					return
				}
				c++
				metricBindingStore.Activities[i].MetricFormula.Formula[j].BindManageData = model.BindManageData{
					ManageDataOrder: c,
					ManageDataType:  "STAGE BINDDATA",
					BindData:        bindValue,
					Key:             keyVD,
					Value:           valueVD,
				}
				manageDataOpArray = append(manageDataOpArray, &valueDefinition)
			} else {
				//Master data type
				bindValue := model.ValueBuilder{
					ValueUUID:           formula.ID,
					WorkflowID:          "test",
					StageID:             activity.StageID,
					TracabilityDataType: 7,
					BindingType:         1,
				}
				valueDefinition, keyVD, valueVD, err := metricBinding.BuildGeneralValueManageData(bindValue)
				if err != nil {
					metricBindingStore.ErrorMessage = err.Error()
					_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
					if errResult != nil {
						logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
					}
					commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "BuildFormulaDefinition ")
					return
				}
				c++
				metricBindingStore.Activities[i].MetricFormula.Formula[j].BindManageData = model.BindManageData{
					ManageDataOrder: c,
					ManageDataType:  "STAGE BINDDATA",
					BindData:        bindValue,
					Key:             keyVD,
					Value:           valueVD,
				}
				manageDataOpArray = append(manageDataOpArray, &valueDefinition)

				//Meta data builder
				metaDataB64 := b64.StdEncoding.EncodeToString([]byte("Artifact template name"))
				metaDataBuilder, errInMetaDataBuilder := metricBinding.BuildMetaData(metaDataB64)
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
				manageDataOpArray = append(manageDataOpArray, &metaDataBuilder)

				//Primary key column name builder
				pkColumnNameB64 := b64.StdEncoding.EncodeToString([]byte("Primary Key Column Name"))
				pkColBuilder, errInPKColBuilder := metricBinding.BuildPrimaryKeyColumn(pkColumnNameB64)
				if errInPKColBuilder != nil {
					logrus.Error("Buidling primary key column failed ", errInPKColBuilder.Error())
					metricBindingStore.ErrorMessage = errInMetaDataBuilder.Error()
					_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
					if errResult != nil {
						logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
					}
					commons.JSONErrorReturn(w, r, errInPKColBuilder.Error(), http.StatusInternalServerError, "BuildFormulaDefinition ")
					return
				}
				manageDataOpArray = append(manageDataOpArray, &pkColBuilder)

			}
		}
	}
	metricBindingStore.Transactions = model.TransacionDetailsMetricBinding{
		Memo:                []byte(memo),
		MetricID:            metricBindJson.ID,
		MtericMapId:         metricMapID,
		TenantMapId:         uint32(tenantMapId),
		MetricName:          metricBindJson.Name,
		NoOfActivityFormula: len(metricBindJson.Activities),
		PublisherIdentityManageData: model.PublisherIdentity{
			ManageDataOrder: 1,
			ManageDataType:  "PUBLISHER IDENTITY",
			UserID:          metricBindJson.User.ID,
			PublicKey:       metricBindJson.User.Publickey,
			ManageDataKey:   publisherIdentity.Name,
			ManageDataValue: publisherIdentity.Value,
		},
		TotalNumberOfManageData:   len(manageDataOpArray),
		MaxNumOfManageDatePerHash: 25,
	}
	// split manage data in to 25 length sub arrays
	manageData2dArray := commons.ChunkSlice(manageDataOpArray, 25)
	var errorInXDR error
	// loop the manage data opration2d array and build trasacion
	for _, manadataOperationArray := range manageData2dArray {
		stellarprotocol := stellarprotocols.StellarTrasaction{
			PublicKey:  constants.PublicKey,
			SecretKey:  constants.SecretKey,
			Operations: manadataOperationArray,
			Memo:       memo,
		}
		// submit transaction
		err, errCode, hash := stellarprotocol.SubmitToStellerBlockchain()
		errorInXDR = err
		if err != nil {
			// status = "Failed"
			metricBindingStore.ErrorMessage = err.Error()
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
		metricBindingStore.Transactions.TXNHashes = append(metricBindingStore.Transactions.TXNHashes, hash)
		logrus.Info("Transaction Hash ", hash)
	}
	if errorInXDR != nil {
		metricBindingStore.ErrorMessage = err.Error()
		_, errResult := object.InsertMetricBindingFormula(metricBindingStore)
		if errResult != nil {
			logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
		}
		commons.JSONErrorReturn(w, r, errorInXDR.Error(), http.StatusInternalServerError, "Error when submitting transaction to blockchain ")
		return
	} else {
		metricBindingStore.Status = "SAVED IN BLOCKCHAIN"
		id, errResult := object.InsertMetricBindingFormula(metricBindingStore)
		if errResult != nil {
			logrus.Error("Error while inserting the metric binding formula into DB: ", errResult)
		}
		w.WriteHeader(http.StatusOK)
		response := model.SuccessResponseMetricBinding{
			Code:              http.StatusOK,
			ID:                id,
			MetricID:          metricBindJson.ID,
			TransactionHashes: metricBindingStore.Transactions.TXNHashes,
		}
		json.NewEncoder(w).Encode(response)
		return
	}
}
