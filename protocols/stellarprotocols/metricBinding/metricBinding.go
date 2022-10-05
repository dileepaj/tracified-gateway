package metricBinding

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
	expertformula "github.com/dileepaj/tracified-gateway/protocols/stellarprotocols/expertFormula"
	"github.com/stellar/go/txnbuild"

	"github.com/sirupsen/logrus"
)

func StellarMetricBinding(w http.ResponseWriter, r *http.Request, metricBindJson model.MetricDataBindingRequest) {
	metricBinding := MetricBinding{}
	stellarProtocol := expertformula.ExpertFormula{}
	var manageDataOpArray []txnbuild.Operation
	object := dao.Connection{}

	// building memo
	// mapMetricId uint64, metricName string, tenantId uint32, noOfFormula int32
	metricMapID, err := InsertAndFindMetricID(metricBindJson.ID, metricBindJson.Name)
	if err != nil {
		commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "InsertAndFindMetricID ")
		return
	}
	tenantMapId, err := InsertAndFindTenentID(metricBindJson.TenantId)
	if err != nil {
		commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "InsertAndFindMetricID ")
		return
	}

	memo, errInMemoBuilder := metricBinding.BuildMemo(metricMapID, metricBindJson.Name, uint32(tenantMapId), int32(len(metricBindJson.Activities)))
	if errInMemoBuilder != nil {
		commons.JSONErrorReturn(w, r, errInMemoBuilder.Error(), http.StatusInternalServerError, "")
		return
	}

	AuthorIdentity, err := stellarProtocol.BuildAuthorManageData(metricBindJson.UserPK)
	if err != nil {
		commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "BuildAuthorManageData ")
		return
	}
	manageDataOpArray = append(manageDataOpArray, &AuthorIdentity)
	fmt.Println("FFFFF",metricBindJson.Activities)
	for _, activity := range metricBindJson.Activities {
		fmt.Println("FFFFF")
		for _, formula := range activity.MetricFormula.Formula {
			stageID, err := strconv.Atoi(activity.StageID)
			if err != nil {
				commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Stage ID not a int string ")
				return
			}
			// checked whether given formulaID already in the database or not
			formulaMapID, err := object.GetFormulaMapID(activity.MetricFormula.MetricExpertFormula.ID).Then(func(data interface{}) interface{} {
				return data
			}).Await()
			formulaDetails := formulaMapID.(model.FormulaIDMap)
			if err != nil {
				commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Can not find the Formula in database ")
				return
			}
			activityMapId, err := InsertAndFindActivityID(activity.ID, activity.Name, activity.MetricID, activity.StageID)
			if err != nil {
				commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "InsertAndFindMetricID ")
				return
			}
			formulaDefinition, err := metricBinding.BuildFormulaDefinition(formulaDetails.MapID, activityMapId, uint32(stageID), uint32(len(activity.MetricFormula.Formula)), activity.Name)
			if err != nil {
				commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "BuildFormulaDefinition ")
				return
			}
			manageDataOpArray = append(manageDataOpArray, &formulaDefinition)

			if formula.ArtifactTemplateID == "" {
				valueDetails, err := object.GetValueMapID(formula.ID).Then(func(data interface{}) interface{} {
					return data
				}).Await()
				bindValue := model.GeneralValueDefBuildRequest{
					ResourceType: "STAGE",
					ResourceName: formula.Key,
					Key:          formula.Key,
					VariableUUID: formula.ID,
					VariableName: valueDetails.(model.ValueIDMap).ValueName,
					BindingType:  0,
					ResourceID:   activity.StageID,
					FormulaID:    formulaDetails.MapID,
				}
				valueDefinition, err := metricBinding.ValueDefinitionBuilder(bindValue)
				if err != nil {
					commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "BuildFormulaDefinition ")
					return
				}
				manageDataOpArray = append(manageDataOpArray, &valueDefinition)
			}else{
				valueDetails, err := object.GetValueMapID(formula.ID).Then(func(data interface{}) interface{} {
					return data
				}).Await()
				bindValue := model.GeneralValueDefBuildRequest{
					ResourceType: "MASTR",
					ResourceName: formula.Field,
					Key:          formula.Key,
					VariableUUID: formula.ID,
					VariableName: valueDetails.(model.ValueIDMap).ValueName,
					BindingType:  1,
					ResourceID:   formula.ArtifactTemplateID,
					FormulaID:    formulaDetails.MapID,
				}
				valueDefinition, err := metricBinding.ValueDefinitionBuilder(bindValue)
				if err != nil {
					commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "BuildFormulaDefinition ")
					return
				}
				manageDataOpArray = append(manageDataOpArray, &valueDefinition)
			}
			

		}
	}


	stellarprotocol := stellarprotocols.StellarTrasaction{
		PublicKey:  constants.PublicKey,
		SecretKey:  constants.SecretKey,
		Operations: manageDataOpArray,
		Memo:       memo,
	}

	// submit transaction
	err, errCode, hash := stellarprotocol.SubmitToStellerBlockchain()
	if err != nil {
		// status = "Failed"
		logrus.Error("Error when submitting transaction to blockchain  ", err)
		w.WriteHeader(errCode)
		response := model.Error{Code: errCode, Message: "Error when submitting transaction to blockchain  " + err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	logrus.Info("Transaction Hash ", hash)
}
