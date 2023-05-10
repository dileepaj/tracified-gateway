package controller

import (
	"net/http"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/apiDemo/businesslogic"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/dileepaj/tracified-gateway/utilities/httpresponse"
)

func SubmitGenesis(w http.ResponseWriter, r *http.Request) {
	var TDP []model.TransactionCollectionBody
	if !commons.DecodeJSONRequestBody(w, r, &TDP) {
		return
	}
	display := &businesslogic.AbstractXDR{TxnBody: TDP}
	err, code, status := display.GenesisAndDataXDRToCron()
	if err == nil && code == http.StatusOK {
		w.WriteHeader(http.StatusOK)
		result := apiModel.SubmitXDRSuccess{
			Status: status,
		}
		httpresponse.SuccessStatus[apiModel.SubmitXDRSuccess](w, result)
		return
	} else {
		utilities.HandleError(w, "Error when getting genesis and data XDR to cron : "+err.Error(), http.StatusGatewayTimeout)
		return
	}
}
