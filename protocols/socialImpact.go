package protocols

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/model"
	expertformula "github.com/dileepaj/tracified-gateway/protocols/stellarprotocols/expertFormula"
	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols/metricBinding"
	"github.com/sirupsen/logrus"
)

// FieldCount - unique formula element Count -formula element count, can be VARIABLE, SEMATICCONSTANT, REFERREDCONSTANT
type AbstractSocialImpact struct {
	Blockchain    string
	FormulaJSON   model.FormulaBuildingRequest
	FieldCount    int
	VariableCount int
}

type AbstractSocialImpactMetricBinding struct {
	Blockchain     string
	MetricBindJSON model.MetricDataBindingRequest
}

/*
SocialImpactExpertFormula
des-This mothod check the blockchain Type and call the relevent method according the blockchain
*/
func (soacialImpact *AbstractSocialImpact) SocialImpactExpertFormula(w http.ResponseWriter, r *http.Request) {
	if soacialImpact.Blockchain == "STELLAR" {
		expertformula.StellarExpertFormulBuilder(w, r, soacialImpact.FormulaJSON, soacialImpact.FieldCount, soacialImpact.VariableCount)
	} else {
		logrus.Error("Blockchain type issue")
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Code: http.StatusBadRequest, Message: "Can not support " + soacialImpact.Blockchain + " yet"}
		json.NewEncoder(w).Encode(response)
		return
	}
}

func (soacialImpact *AbstractSocialImpactMetricBinding) SocialImpactMetricBinding(w http.ResponseWriter, r *http.Request) {
	if soacialImpact.Blockchain == "STELLAR" {
		metricBinding.StellarMetricBinding(w, r, soacialImpact.MetricBindJSON)
	} else {
		logrus.Error("Blockchain type issue")
		w.WriteHeader(http.StatusNoContent)
		response := model.Error{Code: http.StatusNoContent, Message: "Can notsupport " + soacialImpact.Blockchain + " yet"}
		json.NewEncoder(w).Encode(response)
		return
	}
}
