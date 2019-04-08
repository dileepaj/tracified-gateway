package stellarExecuter

import (
	"fmt"
	"net/http"

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
func (cd *ConcreteSubmitXDR) SubmitXDR() model.SubmitXDRResponse {

	var response model.SubmitXDRResponse

	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(cd.XDR)
	if err != nil {
		// panic(err)
		fmt.Println(err.Error())
		response.Error.Code = http.StatusBadRequest
		response.Error.Message = err.Error()
		return response
	}

	// fmt.Println("Successful Transaction:")
	// fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)

	response.Error.Code = http.StatusOK
	response.Error.Message = "Transaction performed in the blockchain."
	response.TXNID = resp.Hash

	return response

}
