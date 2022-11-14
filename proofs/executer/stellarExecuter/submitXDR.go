package stellarExecuter

import (
	"github.com/dileepaj/tracified-gateway/commons"
	log "github.com/sirupsen/logrus"

	"net/http"
	"time"

	"github.com/stellar/go/clients/horizonclient"

	// "github.com/dileepaj/tracified-gateway/api/apiModel"
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
func (cd *ConcreteSubmitXDR) SubmitXDR(tType string) model.SubmitXDRResponse {
	// log.Debug("=========================== submitXDR.go SubmitXDR =============================")
	horizonClient := commons.GetHorizonClient()
	var response model.SubmitXDRResponse
	//s := time.Now().UTC().String()

	//f, err := os.OpenFile("GatewayLogs"+s[:10], os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//if err != nil {
	//	log.Fatalf("error opening file: %v", err)
	//}
	//defer f.Close()

	//log.SetOutput(f)
	// log.Println("This is a test log entry")
	resp, err := horizonClient.SubmitTransactionXDR(cd.XDR)
	if err != nil {
		log.Error("Error while SubmitTransaction to stellar test net " + err.Error())
		error1 := err.(*horizonclient.Error)
		log.Error(error1.Problem.Detail)
		TC, err := error1.ResultCodes()
		if err != nil {
			log.Error("Error while getting ResultCodes from horizon.Error")
		}

		response.Error.Message = TC.TransactionCode

		// log.SetOutput(f)

		log.Error(time.Now().UTC().String() + "- TXNType:" + tType + " " + response.Error.Message + "  ", cd.XDR)
		log.Error(time.Now().UTC().String() + "- TXNType:" + tType + " " + response.Error.Message + "  ")
		log.Error(time.Now().UTC().String() + "- TXNType:" + tType + " " + response.Error.Message + "  ")
		response.Error.Code = http.StatusBadRequest
		// response.Error.Message = err.Error()
		return response
	}

	// fmt.Println("Successful Transaction:")
	// fmt.Println("Ledger:", resp.Ledger)
	log.Info(time.Now().UTC().String() + "- TXNType:" + tType + " Hash:" + resp.Hash)
	log.Info(time.Now().UTC().String() + "- TXNType:" + tType + " Hash:" + resp.Hash)

	response.Error.Code = http.StatusOK
	response.Error.Message = "Transaction performed in the blockchain."
	log.Info("Transaction performed in the blockchain.")
	response.TXNID = resp.Hash

	return response

}
