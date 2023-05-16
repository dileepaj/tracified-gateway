package controller

import (
	"net/http"

	"github.com/dileepaj/tracified-gateway/apiDemo/businesslogic"
	"github.com/dileepaj/tracified-gateway/apiDemo/model/dtos/request"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/dileepaj/tracified-gateway/validations"
)

func SubmitGenesis(w http.ResponseWriter, r *http.Request) {
	var requestBody []request.TransactionCollectionBodyGenesis
	customLogger := utilities.NewCustomLogger()

	if !commons.DecodeJSONRequestBody(w, r, &requestBody) {
		return
	}
	err := validations.ValidateGenesisTDPRequest(requestBody)
	if err != nil {
		customLogger.LogWriter("Error when validating TDP request : "+err.Error(), 3)
	}
	// Use the make function to preallocate the memory
	tdps := make([]model.TransactionCollectionBody, len(requestBody))

	for i := 0; i < len(requestBody); i++ {
		tdps[i] = model.TransactionCollectionBody{
			Identifier: requestBody[i].Identifier,
			SequenceNo: requestBody[i].SequenceNo, ToIdentifier: requestBody[i].ToIdentifier,
			MapIdentifier: requestBody[i].MapIdentifier, XDR: requestBody[i].XDR,
			Orphan: requestBody[i].Orphan, MergeBlock: requestBody[i].MergeBlock,
		}
	}

	display := &businesslogic.AbstractXDR{TxnBody: tdps}
	err, code, status := display.GenesisAndDataXDRToCron()
	if err == nil && code == http.StatusOK {
		utilities.SuccessResponse[string](w, status)
		return
	} else {
		utilities.HandleError(w, "Error when getting genesis and data XDR to cron : "+err.Error(), http.StatusGatewayTimeout)
		return
	}
}