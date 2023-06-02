package stellarExecuter

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dileepaj/tracified-gateway/commons"
	log "github.com/sirupsen/logrus"

	"github.com/stellar/go/clients/horizonclient"

	"github.com/dileepaj/tracified-gateway/model"
)

type ConcreteSubmitXDR struct {
	XDR string
}

/*SubmitXDR - WORKING MODEL

@author - Azeem Ashraf
@desc - Submits the XDR to stellar horizon api and returns the TXN hash.
@params - XDR
*/
var count int = 0

func (cd *ConcreteSubmitXDR) SubmitXDR(tType string) model.SubmitXDRResponse {
	var count int = 0
	log.Debug("=========================== submitXDR.go SubmitXDR =============================")
	horizonClient := commons.GetHorizonClient()
	var response model.SubmitXDRResponse
	resp, err := horizonClient.SubmitTransactionXDR(cd.XDR)
	if err != nil {
		error1 := err.(*horizonclient.Error)
		log.Error(" Error Message Problem ", error1.Problem)
		log.Error(" Error Message response ", error1.Response)
		log.Error(" Error Response ", response.Error.Message)
		// Something went wrong on the Horizon server’s end.
		if error1.Response.StatusCode == 500 && error1.Problem.Status == 500 {
			time.Sleep(10 * time.Second)
			log.Info(" Resubmitting transaction (Timeout issue) ", cd.XDR)
			display := ConcreteSubmitXDR{XDR: cd.XDR}
			return display.SubmitXDR(tType)
			// Timeout - this try to resubmit only 10 times each 20 second
		} else if (error1.Response.StatusCode == 504 && count < 10 || error1.Problem.Status == 504 && count < 10) || (response.Error.Message == "tx_insufficient_fee" && count < 10) {
			time.Sleep(20 * time.Second)
			log.Info(" Resubmitting transaction (Timeout issue) ", cd.XDR)
			display := ConcreteSubmitXDR{XDR: cd.XDR}
			count = count + 1
			return display.SubmitXDR(tType)
			// other errors
		} else {
			TC, err := error1.ResultCodes()
			if err != nil {
				log.Error(" Error while getting ResultCodes from horizon.Error ")
			}
			if TC != nil {
				response.Error.Message = TC.TransactionCode
			}
			log.Error(time.Now().UTC().String()+" - TXNType: "+tType+" "+response.Error.Message+" - Error Message: ", response.Error.Message, " - XDR ", cd.XDR)
			response.Error.Code = http.StatusBadRequest
			return response
		}
	} else {
		log.Info( " Ledger: ", resp.Ledger)
		log.Info(time.Now().UTC().String() + " - TXNType: " + tType + " Hash: " + resp.Hash)
		response.Error.Code = http.StatusOK
		response.Error.Message = " Transaction performed in the blockchain. "
		log.Info(" Transaction performed in the blockchain. " + resp.Hash)
		response.TXNID = resp.Hash
		return response
	}
}
