package stellarExecuter

import (
	"net/http"
	"time"

	"github.com/dileepaj/tracified-gateway/commons"
	log "github.com/sirupsen/logrus"

	"github.com/stellar/go/clients/horizonclient"

	"github.com/dileepaj/tracified-gateway/model"
)

type ConcreteSubmitXDR struct {
	XDR string
}

/*
SubmitXDR - WORKING MODEL

@author - Azeem Ashraf
@desc - Submits the XDR to stellar horizon api and returns the TXN hash.
@params - XDR
*/
var count int = 0

func (cd *ConcreteSubmitXDR) SubmitXDR(tType string) model.SubmitXDRResponse {
	log.Debug("=========================== submitXDR.go SubmitXDR =============================")
	horizonClient := commons.GetHorizonClient()
	var response model.SubmitXDRResponse
	resp, err := horizonClient.SubmitTransactionXDR(cd.XDR)
	if err != nil {
		error1 := err.(*horizonclient.Error)
		log.Error("Error Message Problem ", error1.Problem)
		log.Error("Error Message response ", error1.Response)
		// Something went wrong on the Horizon server’s end.
		if error1.Response.StatusCode == 500 && error1.Problem.Status == 500 {
			response.Error.Code = http.StatusInternalServerError
			time.Sleep(40 * time.Second)
			log.Info("Resubmitting transaction (Timeout issue) ", cd.XDR)
			display := ConcreteSubmitXDR{XDR: cd.XDR}
			count = count + 1
			if count >= 10 {
				count = 0
				response.Error.Message = "Error in resubmitting transaction (Internal server Error)"
				return response
			}
			return display.SubmitXDR(tType)
			// Timeout - this try to resubmit only 10 times each 20 second
		} else if (error1.Response.StatusCode == 504 && count < 10 || error1.Problem.Status == 504 && count < 10) || (response.Error.Message == "tx_insufficient_fee" && count < 10) {
			response.Error.Code = http.StatusGatewayTimeout
			time.Sleep(40 * time.Second)
			log.Info("Resubmitting transaction (Timeout issue) ", cd.XDR)
			display := ConcreteSubmitXDR{XDR: cd.XDR}
			count = count + 1
			if count >= 10 {
				count = 0
				response.Error.Message = "Error in resubmitting transaction (Timeout issue)"
				return response
			}
			return display.SubmitXDR(tType)
			// other errors
		} else {
			TC, err := error1.ResultCodes()
			if err != nil {
				log.Error("Error while getting ResultCodes from horizon.Error ")
				response.Error.Message = "Error while getting ResultCodes from horizon.Error"
			}
			if TC != nil {
				response.Error.Message = TC.TransactionCode
			}
			log.Error(time.Now().UTC().String()+" - TXNType: "+tType+" "+response.Error.Message+" - Error Message: ", response.Error.Message, " - XDR ", cd.XDR)
			response.Error.Code = http.StatusBadRequest
			return response
		}
	} else {
		log.Info(time.Now().UTC().String()+" - TXNType: "+tType+" Hash: "+resp.Hash+" Ledger: ", resp.Ledger)
		log.Info("Transaction performed in the blockchain. " + resp.Hash)
		response.Error.Code = http.StatusOK
		response.TXNID = resp.Hash
		return response
	}
}
