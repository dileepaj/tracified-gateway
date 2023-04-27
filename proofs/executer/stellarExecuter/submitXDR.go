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
	log.Debug("=========================== submitXDR.go SubmitXDR =============================")
	horizonClient := commons.GetHorizonClient()
	var response model.SubmitXDRResponse
	resp, err := horizonClient.SubmitTransactionXDR(cd.XDR)
	if err != nil {
		log.Error("Error while SubmitTransaction to stellar test net " + err.Error())
		error1 := err.(*horizonclient.Error)
		log.Error(" Error Message ", error1.Problem.Detail)
		log.Error(time.Now().UTC().String() + "- Error code: " + strconv.Itoa(error1.Response.StatusCode))
		log.Error(time.Now().UTC().String() + "- Error code: " + strconv.Itoa(error1.Problem.Status))
		RC, err := error1.ResultCodes()
		if err != nil {
			log.Error("Error while getting ResultCodes from horizon.Error")
		}
		if RC != nil {
			response.Error.Message = RC.TransactionCode
		}
		log.Error(" Error Response ", response.Error.Message)
		if (error1.Response.StatusCode == 504 && error1.Problem.Status == 504 && count < 10) || (response.Error.Message == "tx_insufficient_fee" && count < 10) {
			log.Info("Resubmitting transaction", cd.XDR)
			display := ConcreteSubmitXDR{XDR: cd.XDR}
			count = count + 1
			display.SubmitXDR(tType)
		} else if (error1.Response.StatusCode == 504 && error1.Problem.Status == 504 && count >= 10) || (response.Error.Message == "tx_insufficient_fee" && count >= 10) {
			log.Error("Error while submitting xdr to blockchain after retrying 10 times")
			return response
		} else {
			TC, err := error1.ResultCodes()
			if err != nil {
				log.Error("Error while getting ResultCodes from horizon.Error")
			}
			if TC != nil {
				response.Error.Message = TC.TransactionCode
			}
			log.Error(time.Now().UTC().String()+"- TXNType:"+tType+" "+response.Error.Message+"  ", cd.XDR)
			log.Error(time.Now().UTC().String() + "- Error Message: " + response.Error.Message)
			response.Error.Code = http.StatusBadRequest
			return response
		}
	} else {
		log.Info("Ledger:", resp.Ledger)
		log.Info(time.Now().UTC().String() + "- TXNType:" + tType + " Hash:" + resp.Hash)
		response.Error.Code = http.StatusOK
		response.Error.Message = "Transaction performed in the blockchain."
		log.Info("Transaction performed in the blockchain.  " + resp.Hash)
		response.TXNID = resp.Hash
		return response
	}
	return response
}
