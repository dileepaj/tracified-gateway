package stellarExecuter

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/stellar/go/clients/horizon"

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
func (cd *ConcreteSubmitXDR) SubmitXDR(testnet bool, tType string) model.SubmitXDRResponse {
	fmt.Println("=========================== submitXDR.go SubmitXDR =============================")
	var response model.SubmitXDRResponse
	s := time.Now().UTC().String()

	f, err := os.OpenFile("GatewayLogs"+s[:10], os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	// log.Println("This is a test log entry")

	if testnet {
		fmt.Println("test net")
		resp, err := horizon.DefaultTestNetClient.SubmitTransaction(cd.XDR)
		fmt.Println(cd.XDR)
		if err != nil {
			fmt.Println("Error while SubmitTransaction to stellar test net " + err.Error())
			error1 := err.(*horizon.Error)
			fmt.Println(error1.Problem.Detail)
			TC, err := error1.ResultCodes()
			if err != nil{
				fmt.Println("Error while getting ResultCodes from horizon.Error")
			}
			for _, element := range TC.OperationCodes {
				response.Error.Message = response.Error.Message + element + "? "
			}

			// log.SetOutput(f)
			fmt.Println(time.Now().UTC().String() + "- TXNType:" + tType + " " + response.Error.Message)
			fmt.Println(time.Now().UTC().String() + "- TXNType:" + tType + " " + response.Error.Message)

			response.Error.Code = http.StatusBadRequest
			// response.Error.Message = err.Error()
			return response
		}

		// fmt.Println("Successful Transaction:")
		// fmt.Println("Ledger:", resp.Ledger)
		fmt.Println(time.Now().UTC().String() + "- TXNType:" + tType + " Hash:" + resp.Hash)
		fmt.Println(time.Now().UTC().String() + "- TXNType:" + tType + " Hash:" + resp.Hash)

		response.Error.Code = http.StatusOK
		response.Error.Message = "Transaction performed in the blockchain."
		fmt.Println("Transaction performed in the blockchain.")
		response.TXNID = resp.Hash
	} else {
		fmt.Println("public net")
		resp, err := horizon.DefaultPublicNetClient.SubmitTransaction(cd.XDR)
		if err != nil {
			fmt.Println("Error while SubmitTransaction to stellar public net " + err.Error())
			error1 := err.(*horizon.Error)
			TC, _ := error1.ResultCodes()
			for _, element := range TC.OperationCodes {
				response.Error.Message = response.Error.Message + element + "? "
			}

			// log.SetOutput(f)
			fmt.Println(time.Now().UTC().String() + "- TXNType:" + tType + " " + response.Error.Message)
			fmt.Println(time.Now().UTC().String() + "- TXNType:" + tType + " " + response.Error.Message)

			response.Error.Code = http.StatusBadRequest
			// response.Error.Message = err.Error()
			return response
		}

		// fmt.Println("Successful Transaction:")
		// fmt.Println("Ledger:", resp.Ledger)
		fmt.Println(time.Now().UTC().String() + "- TXNType:" + tType + " Hash:" + resp.Hash)
		fmt.Println(time.Now().UTC().String() + "- TXNType:" + tType + " Hash:" + resp.Hash)

		response.Error.Code = http.StatusOK
		response.Error.Message = "Transaction performed in the blockchain."
		fmt.Println("Transaction performed in the blockchain.")
		response.TXNID = resp.Hash
	}

	return response

}
