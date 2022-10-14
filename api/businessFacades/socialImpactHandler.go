package businessFacades

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/authentication"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols"
	"github.com/dileepaj/tracified-gateway/validations"
)

/*
BuildSocialImpactFormula
des-This handler the expert formula bulding,
	In side this handler
	 * validate the JSON request body
	 * changethe formual elemet type according to the protocol
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
		// validation againt trust network
		errInTrustNetworkValidation := authentication.ValidateAgainstTrustNetwork(formulaJSON.User.Publickey)
		if errInTrustNetworkValidation != nil {
			commons.JSONErrorReturn(w, r, errInTrustNetworkValidation.Error(), http.StatusBadRequest, "Expert is not in the trust network ")
			return
		}
		authLayer := authentication.AuthLayer{
			FormulaId:    formulaJSON.ID,
			ExpertPK:     formulaJSON.User.Publickey,
			ExpertUserID: formulaJSON.User.ID,
			CiperText:    formulaJSON.CiperText,
			Plaintext:    formulaJSON.Formula,
		}
		err, errCode := authLayer.ValidateExpertRequest()
		if err != nil {
			commons.JSONErrorReturn(w, r, err.Error(), errCode, "Authentication Issue ")
			return
		}

		formulaArray := formulaJSON.Formula
		fieldCount := 0
		for i, element := range formulaJSON.Formula {
			if element.Type == "DATA" {
				formulaArray[i].Type = "VARIABLE"
			} else if element.Type == "CONSTANT" && element.MetricReferenceId != "" {
				formulaArray[i].Type = "REFERREDCONSTANT"
			} else if element.Type == "CONSTANT" && element.MetricReferenceId == "" {
				formulaArray[i].Type = "SEMANTICCONSTANT"
			}
			if element.Type != "OPERATOR" {
				fieldCount++
			}
		}
		formulaJSON.Formula = formulaArray
		// build the abstract struct and call the SocialImpactExpertFormula
		socialImpactBuilder := protocols.AbstractSocialImpact{
			Blockchain:  formulaJSON.Blockchain,
			FormulaJSON: formulaJSON,
			FieldCount:  fieldCount,
		}
		socialImpactBuilder.SocialImpactExpertFormula(w, r)
	}
}

// BindMetric method : binds the metric with mutiple formulas
func BindMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	permissionStatus := authentication.HasPermission(r.Header.Get("Authorization"))
	if !permissionStatus.Status || !permissionStatus.IsSubscriptionPaid {
		commons.JSONErrorReturn(w, r, "", http.StatusUnauthorized, "Status Unauthorized")
		return
	} else {
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
			metricBuilder := protocols.AbstractSocialImpactMetricBinding{
				Blockchain:     metricBindJSON.Metric.Blockchain,
				MetricBindJSON: metricBindJSON,
			}
			metricBuilder.SocialImpactMetricBinding(w, r)
		}
	}
}
