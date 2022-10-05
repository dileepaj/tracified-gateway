package metricBinding

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
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

	// building memo
	// mapMetricId uint64, metricName string, tenantId uint32, noOfFormula int32
	// metricmapID,err:= InsertAndFindMetricID(metricBindJson.ID,metricBindJson.Name,)

	memo, errInMemoBuilder := metricBinding.BuildMemo(27, metricBindJson.Name, 54, int32(len(metricBindJson.Activity)))
	if errInMemoBuilder != nil {
		commons.JSONErrorReturn(w, r, errInMemoBuilder.Error(), http.StatusInternalServerError, "")
		return
	}


	AuthorIdentity,err:= stellarProtocol.BuildAuthorManageData(metricBindJson.UserPK)
	if err != nil {
		commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "BuildAuthorManageData ")
		return
	}
	manageDataOpArray=append(manageDataOpArray, &AuthorIdentity)

	stellarprotocol := stellarprotocols.StellarTrasaction{
		PublicKey: constants.PublicKey,
		SecretKey: constants.SecretKey,
		Operations:  manageDataOpArray,
		Memo: memo,
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
