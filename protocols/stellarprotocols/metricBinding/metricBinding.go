package metricBinding

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"

	"github.com/sirupsen/logrus"
)

func StellarMetricBinding(w http.ResponseWriter, r *http.Request, metricBindJson model.MetricDataBindingRequest) {
	w.Header().Set("Content-Type", "application/json")
	//var status string

	//call metric id database method

	//noOfFormula := int32(len(metricBindJson.Activity))
	//building memo
	//memo, errInMemoBuilder := metricbinding.BuildMemo(metricBindJson.ID, metricBindJson.Name, metricBindJson.TenantId, noOfFormula)

	stellarprotocol := stellarprotocols.StellarTrasaction{
		PublicKey:  constants.PublicKey,
		SecretKey:  constants.SecretKey,
		// Operations: {},
		Memo:       "helllooooooo",
	}

	//submit transaction
	err, errCode, hash := stellarprotocol.SubmitToStellerBlockchain()
	if err != nil {
		//status = "Failed"
		logrus.Error("Error when submitting transaction to blockchain  ", err)
		w.WriteHeader(errCode)
		response := model.Error{Code: errCode, Message: "Error when submitting transaction to blockchain  " + err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	logrus.Info("Transaction Hash ", hash)
}
