package businessFacades

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dileepaj/tracified-gateway/authentication"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols"
	"github.com/dileepaj/tracified-gateway/validations"
	"github.com/relvacode/iso8601"
	"github.com/sirupsen/logrus"
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
		logrus.Error(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while decoding the body")
		return
	}

	errInJsonValidation := validations.ValidateFormulaBuilder(formulaJSON)
	if errInJsonValidation != nil {
		logrus.Error("Request body failed the validation check : ", errInJsonValidation)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Request body is invalid, Error : " + errInJsonValidation.Error())
		return
	} else {
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

//BindMetric method : binds the metric with mutiple formulas
func BindMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var metricBindJSON model.MetricBindingRequest

	err := json.NewDecoder(r.Body).Decode(&metricBindJSON)
	if err != nil {
		logrus.Error(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while decoding the body")
		return
	}

	errInJsonValidationInMetricBind := validations.ValidateMetricBindingRequest(metricBindJSON)
	if errInJsonValidationInMetricBind != nil {
		logrus.Error("Request body failed the validation check : ", errInJsonValidationInMetricBind)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Request body is invalid, Error : " + errInJsonValidationInMetricBind.Error())
		return
	} else {
		// currentTime := time.Now()
		// layout := "2006-01-02T15:04:05Z"
		timein := time.Now().Local().Add(time.Hour*time.Duration(0) +
			time.Minute*time.Duration(5) +
			time.Second*time.Duration(0))
		convertedFromTime, _ := iso8601.ParseString(time.Now().String())
		convertedToTime, _ := iso8601.ParseString(timein.String())
		fmt.Println("Time from ", convertedFromTime, " time to ", convertedToTime)
		// //build the abstract and call the relevent blockchain
		// authentication.API_Throttler("TestEntityType2", "TestEntity2", 5, 1)
		apiReq := model.API_ThrottlerRequest{
			RequestEntityType: "Test",
			RequestEntity:     "PK",
			FormulaID:         "234",
			AllowedAmount:     4,
			FromTime:          convertedFromTime,
			ToTime:            convertedToTime,
		}

		authentication.API_Throttler(apiReq)
	}
}
