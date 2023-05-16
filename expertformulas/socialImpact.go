package expertformulas

import (
	"encoding/json"
	"net/http"

	ethereuemmetricbind "github.com/dileepaj/tracified-gateway/expertformulas/ethereum/codeGenerator/ethereuemMetricBind"
	ethereumExpertFormula "github.com/dileepaj/tracified-gateway/expertformulas/ethereum/codeGenerator/ethereumExpertFormula"
	expertFormula "github.com/dileepaj/tracified-gateway/expertformulas/stellarprotocols/expertFormula"
	"github.com/dileepaj/tracified-gateway/expertformulas/stellarprotocols/metricBinding"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

// FieldCount - unique formula element Count -formula element count, can be VARIABLE, SEMATICCONSTANT, REFERREDCONSTANT
type AbstractSocialImpact struct {
	Blockchain    []int
	FormulaJSON   model.FormulaBuildingRequest
	FieldCount    int
	VariableCount int
	ExpertId      string
}

type AbstractSocialImpactMetricBinding struct {
	Blockchain     string
	MetricBindJSON model.MetricDataBindingRequest
}

/*
SocialImpactExpertFormula
des-This method check the blockchain Type and call the relevant method according the blockchain
1 - Stellar
2 - Ethereum
*/
func (socialImpact *AbstractSocialImpact) SocialImpactExpertFormula(w http.ResponseWriter, r *http.Request) {
	if socialImpact.Blockchain == "STELLAR" {
		expertFormula.StellarExpertFormulaBuilder(w, r, socialImpact.FormulaJSON, socialImpact.FieldCount, socialImpact.VariableCount, socialImpact.ExpertId)
	} else if socialImpact.Blockchain == "ETHEREUM" {
		ethereumExpertFormula.SmartContractGeneratorForFormula(w, r, socialImpact.FormulaJSON, socialImpact.FieldCount)
	} else {
		logrus.Error("Blockchain type issue")
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Code: http.StatusBadRequest, Message: "Can not support " + socialImpact.Blockchain + " yet"}
		json.NewEncoder(w).Encode(response)
		return
	}
}

func (socialImpact *AbstractSocialImpactMetricBinding) SocialImpactMetricBinding(w http.ResponseWriter, r *http.Request) {
	if socialImpact.Blockchain == "STELLAR" {
		metricBinding.StellarMetricBinding(w, r, socialImpact.MetricBindJSON)
	} else if socialImpact.Blockchain == "ETHEREUM" {
		ethereuemmetricbind.SmartContractHandlerForMetric(w, r, socialImpact.MetricBindJSON)
	} else {
		logrus.Error("Blockchain type issue")
		w.WriteHeader(http.StatusNoContent)
		response := model.Error{Code: http.StatusNoContent, Message: "Can not supported " + socialImpact.Blockchain + " yet"}
		json.NewEncoder(w).Encode(response)
		return
	}
}
