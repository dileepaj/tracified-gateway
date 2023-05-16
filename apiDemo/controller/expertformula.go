package controller

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

func BuildSocialImpactExpertFormula(w http.ResponseWriter, r *http.Request) {
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