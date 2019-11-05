package interpreter

import (
	"encoding/json"
	"fmt"

	"strings"

	"github.com/dileepaj/tracified-gateway/api/apiModel"

	// "encoding/json"
	"github.com/stellar/go/xdr"

	"net/http"

	"github.com/dileepaj/tracified-gateway/proofs/retriever/stellarRetriever"
)

type AbstractPOCOC struct {
	Txn   string
	DBCOC xdr.Transaction
	BCCOC xdr.Transaction
	XDR   string
}

/*InterpretPOCOC - Working Model
@author - Azeem Ashraf
@desc - Interprets All the required fields necessary to perform POCOC
@params - ResponseWriter,Request
*/
func (AP *AbstractPOCOC) InterpretPOCOC(w http.ResponseWriter, r *http.Request) {

	object := stellarRetriever.ConcretePOCOC{Txn: AP.Txn}
	bcCOC := object.RetrievePOCOC()
	AP.BCCOC = bcCOC

	fmt.Println(AP.BCCOC.SourceAccount.Address())
	fmt.Println(AP.DBCOC.SourceAccount.Address())

	w.WriteHeader(http.StatusOK)
	result := compareCOC(AP.DBCOC, AP.BCCOC)

	response := CocSpecialResponse{
		Status: result.Status,
		Txn:    AP.Txn,
		Xdr:    AP.XDR,
	}
	json.NewEncoder(w).Encode(response)
	return

}

func compareCOC(db xdr.Transaction, bc xdr.Transaction) apiModel.SubmitXDRSuccess {
	var result apiModel.SubmitXDRSuccess

	if db.SourceAccount.Address() != bc.SourceAccount.Address() {
		result.Status = "Failed, Source Address in Gateway and Blockchain Doesn't match"
	} else if strings.TrimLeft(fmt.Sprintf("%s", db.Operations[0].Body.ManageDataOp.DataValue), "&") !=
		strings.TrimLeft(fmt.Sprintf("%s", bc.Operations[0].Body.ManageDataOp.DataValue), "&") {
		result.Status = "Failed, Txn Type in Gateway and Blockchain Doesn't match"
	} else if strings.TrimLeft(fmt.Sprintf("%s", db.Operations[1].Body.ManageDataOp.DataValue), "&") !=
		strings.TrimLeft(fmt.Sprintf("%s", bc.Operations[1].Body.ManageDataOp.DataValue), "&") {
		result.Status = "Failed, Identifier in Gateway and Blockchain Doesn't match"
	} else if fmt.Sprintf("%s", db.Operations[3].Body.PaymentOp.Asset.AlphaNum12.AssetCode) !=
		fmt.Sprintf("%s", bc.Operations[3].Body.PaymentOp.Asset.AlphaNum12.AssetCode) {
		result.Status = "Failed, Asset Code in Gateway and Blockchain Doesn't match"
	} else if db.Operations[3].Body.PaymentOp.Amount != bc.Operations[3].Body.PaymentOp.Amount {
		result.Status = "Failed, Asset Amount in Gateway and Blockchain Doesn't match"
	} else if db.Operations[3].Body.PaymentOp.Destination.Address() != bc.Operations[3].Body.PaymentOp.Destination.Address() {
		result.Status = "Failed, Destination Address in Gateway and Blockchain Doesn't match"
	} else {
		result.Status = "Success, COC in Gateway and Blockchain matches"
	}

	//temporary creation of the coc txn details for user view

	// var opArray []operation

	// op:=operation{
	// 	name:strings.TrimLeft(fmt.Sprintf("%s", db.Operations[0].Body.ManageDataOp.DataName), "&"),
	// 	value:strings.TrimLeft(fmt.Sprintf("%s", db.Operations[0].Body.ManageDataOp.DataValue), "&"),
	// }

	// op1:=operation{
	// 	name:strings.TrimLeft(fmt.Sprintf("%s", db.Operations[1].Body.ManageDataOp.DataName), "&"),
	// 	value:strings.TrimLeft(fmt.Sprintf("%s", db.Operations[1].Body.ManageDataOp.DataValue), "&"),
	// }

	// n := bytes.Index(db.Operations[3].Body.PaymentOp.Asset.AlphaNum12.AssetCode, []byte{0})

	// s := string(byteArray[:n])

	// op2:=operation{
	// 	name:"Asset Code",
	// 	value:strconv.Itoa(),
	// }
	// op3:=operation{
	// 	name:"Asset Ammount",
	// 	value:db.Operations[3].Body.PaymentOp.Amount,
	// }

	// opArray = append(opArray, op);
	// opArray = append(opArray, op1);
	// opArray = append(opArray, op2);
	// opArray = append(opArray, op3);

	// txe:=transaction{
	// 	SourceAccount:db.SourceAccount.Address(),
	// 	Operations:opArray,
	// 	DestinationAccount:db.Operations[3].Body.PaymentOp.Destination.Address(),

	// }

	return result
}

type CocSpecialResponse struct {
	Status string
	Txn    string
	Xdr    string
}

type transaction struct {
	SourceAccount      string
	DestinationAccount string
	Fee                int32
	SeqNum             int
	TimeBounds         []int
	Operations         []operation
}

type operation struct {
	name  string
	value string
}
