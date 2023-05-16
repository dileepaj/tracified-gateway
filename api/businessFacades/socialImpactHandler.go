package businessFacades

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/authentication"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/configs"
	protocols "github.com/dileepaj/tracified-gateway/expertformulas"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/validations"
	"github.com/sirupsen/logrus"
)

/*
BuildSocialImpactFormula
des-This handler the expert formula building,
	In side this handler
	 * validate the JSON request body
	 * change the formula element type according to the protocol
	 * call the SocialImpactExpertFormula method to build the formula
*/

func BuildSocialImpactExpertFormula(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var formulaJSON model.FormulaBuildingRequest
	err := json.NewDecoder(r.Body).Decode(&formulaJSON)
	if err != nil {
		commons.JSONErrorReturn(w, r, err.Error(), http.StatusBadRequest, "Error while decoding the body ")
		return
	}
	errInJsonValidation := validations.ValidateFormulaBuilder(formulaJSON)
	if errInJsonValidation != nil {
		commons.JSONErrorReturn(w, r, errInJsonValidation.Error(), http.StatusBadRequest, "Request body failed the validation check :")
		return
	} else {
		authLayer := authentication.AuthLayer{
			FormulaId: formulaJSON.MetricExpertFormula.ID,
			ExpertPK:  formulaJSON.Verify.PublicKey,
			Signature: formulaJSON.Verify.Signature,
			Plaintext: formulaJSON.Verify.Payload,
		}
		logrus.Info("Expert's public key  ", formulaJSON.Verify.PublicKey)
		err, errCode, id := authLayer.ValidateExpertRequest()
		if err != nil {
			commons.JSONErrorReturn(w, r, err.Error(), errCode, "Authentication Issue, ")
			return
		} else {
			formulaArray := formulaJSON.MetricExpertFormula.Formula
			fieldCount := 0
			variableCount := 0
			for i, element := range formulaJSON.MetricExpertFormula.Formula {
				if element.Type == "DATA" {
					formulaArray[i].Type = "VARIABLE"
					variableCount++
				} else if element.Type == "CONSTANT" && element.MetricReferenceId != "" {
					formulaArray[i].Type = "REFERREDCONSTANT"
				} else if element.Type == "CONSTANT" && element.MetricReferenceId == "" {
					formulaArray[i].Type = "SEMANTICCONSTANT"
				}
				if element.Type != "OPERATOR" {
					fieldCount++
				}
			}
			formulaJSON.MetricExpertFormula.Formula = formulaArray

			// if the blockchain is not provided in the request, then use the default blockchain
			if formulaJSON.MetricExpertFormula.Blockchain == "" {
				formulaJSON.MetricExpertFormula.Blockchain = configs.DefaultBlockchain
			}
			// build the abstract struct and call the SocialImpactExpertFormula
			socialImpactBuilder := protocols.AbstractSocialImpact{
				Blockchain:    formulaJSON.MetricExpertFormula.Blockchain,
				FormulaJSON:   formulaJSON,
				FieldCount:    fieldCount,
				VariableCount: variableCount,
				ExpertId:      id,
			}
			socialImpactBuilder.SocialImpactExpertFormula(w, r)
		}
	}
}

// BindMetric method : binds the metric with multiple formulas
func BindMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if configs.JWTAuthEnableBindMetricEndpoint {
		permissionStatus := authentication.HasPermission(r.Header.Get("Authorization"))
		if !permissionStatus.Status || !permissionStatus.IsSubscriptionPaid {
			commons.JSONErrorReturn(w, r, "", http.StatusUnauthorized, "Status Unauthorized")
			return
		}
	}
	var metricBindJSON model.MetricDataBindingRequest
	err := json.NewDecoder(r.Body).Decode(&metricBindJSON)
	if err != nil {
		commons.JSONErrorReturn(w, r, err.Error(), http.StatusBadRequest, "Error while decoding the body ")
		return
	}
	errInJsonValidationInMetricBind := validations.ValidateMetricDataBindingRequest(metricBindJSON)
	if errInJsonValidationInMetricBind != nil {
		commons.JSONErrorReturn(w, r, errInJsonValidationInMetricBind.Error(), http.StatusBadRequest, "Request body is invalid, Error :")
		return
	} else {
		if metricBindJSON.Metric.Blockchain == "" {
			metricBindJSON.Metric.Blockchain = configs.DefaultBlockchain
		}
		metricBuilder := protocols.AbstractSocialImpactMetricBinding{
			Blockchain:     metricBindJSON.Metric.Blockchain,
			MetricBindJSON: metricBindJSON,
		}
		metricBuilder.SocialImpactMetricBinding(w, r)
	}
}
