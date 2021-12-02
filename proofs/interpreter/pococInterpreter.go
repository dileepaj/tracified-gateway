package interpreter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/stellar/go/support/log"

	// "encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/proofs/retriever/stellarRetriever"
)

/*
type AbstractPOCOC struct {
	ProofHash  string
	Txn        string
	DBCOC      xdr.Transaction
	BCCOC      xdr.Transaction
	XDR        string
	COCStatus  string
	SequenceNo string
}
*/
/*InterpretPOCOC - Working Model
@author - Azeem Ashraf
@desc - Interprets All the required fields necessary to perform POCOC
@params - ResponseWriter,Request
*/
/*
func (AP *AbstractPOCOC) InterpretPOCOC(w http.ResponseWriter, r *http.Request) {

	var result []model.POCOCResponse
	FromSigned := true
	ToSigned := false
	if AP.COCStatus != "pending" && AP.COCStatus != "expired" {
		ToSigned = true
	}

	SequenceNo, errSequence := strconv.ParseInt(AP.SequenceNo, 10, 64)
	if errSequence == nil {
		fmt.Printf("%d of type %T", SequenceNo, SequenceNo)
	}
	mapD := map[string]string{"transaction": AP.Txn}
	mapB, _ := json.Marshal(mapD)
	fmt.Println(string(mapB))
	// trans := transaction{transaction:TxnHash}
	// s := fmt.Sprintf("%v", trans)

	encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
	text := (string(encoded))

	object := stellarRetriever.ConcretePOCOC{Txn: AP.ProofHash}
	bcCOC, state, timestamp, ledger, feePaid := object.RetrievePOCOC()
	if !state {
		// w.WriteHeader(http.StatusBadRequest)
		// response := model.Error{Message: "Failed to retrieve blockchain proof transaction"}
		// json.NewEncoder(w).Encode(response)

		w.WriteHeader(http.StatusOK)
		temp := model.POCOCResponse{
			Txnhash: AP.Txn,
			Url: "https://www.stellar.org/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
				text + "%3D%3D&network=public",
			Identifier:     strings.TrimLeft(fmt.Sprintf("%s", AP.DBCOC.Operations[1].Body.ManageDataOp.DataValue), "&"),
			From:           AP.DBCOC.Operations[1].SourceAccount.Address(),
			To:             AP.DBCOC.Operations[3].Body.PaymentOp.Destination.Address(),
			Timestamp:      timestamp,
			Ledger:         ledger,
			FeePaid:        feePaid,
			TxnType:        "coc",
			SequenceNo:     SequenceNo,
			Quantity:       strings.TrimRight(strconv.FormatInt(int64(AP.DBCOC.Operations[3].Body.PaymentOp.Amount), 10), "0"),
			AssetCode:      strings.TrimRight(fmt.Sprintf("%s", AP.DBCOC.Operations[3].Body.PaymentOp.Asset.AlphaNum12.AssetCode), "\u0000"),
			Status:         "Failed to retrieve blockchain proof transaction from Stellar",
			BlockchainName: "Stellar",
			FromSigned:     FromSigned,
			ToSigned:       ToSigned,
			COCStatus:      AP.COCStatus}
		result = append(result, temp)
		json.NewEncoder(w).Encode(result)

		return
	}
	AP.BCCOC = bcCOC

	fmt.Println(AP.BCCOC.SourceAccount.Address())
	fmt.Println(AP.DBCOC.Operations[0].SourceAccount.Address())

	res := compareCOC(AP.DBCOC, AP.BCCOC)

	w.WriteHeader(http.StatusOK)
	temp := model.POCOCResponse{
		Txnhash: AP.Txn,
		Url: "https://www.stellar.org/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
			text + "%3D%3D&network=public",
		Identifier:     strings.TrimLeft(fmt.Sprintf("%s", AP.DBCOC.Operations[1].Body.ManageDataOp.DataValue), "&"),
		From:           AP.DBCOC.Operations[1].SourceAccount.Address(),
		To:             AP.DBCOC.Operations[3].Body.PaymentOp.Destination.Address(),
		Timestamp:      timestamp,
		Ledger:         ledger,
		FeePaid:        feePaid,
		TxnType:        "coc",
		SequenceNo:     SequenceNo,
		Quantity:       strings.TrimRight(strconv.FormatInt(int64(AP.DBCOC.Operations[3].Body.PaymentOp.Amount), 10), "0"),
		AssetCode:      strings.TrimRight(fmt.Sprintf("%s", AP.DBCOC.Operations[3].Body.PaymentOp.Asset.AlphaNum12.AssetCode), "\u0000"),
		Status:         res.Status,
		BlockchainName: "Stellar",
		FromSigned:     FromSigned,
		ToSigned:       ToSigned,
		COCStatus:      AP.COCStatus}
	result = append(result, temp)
	json.NewEncoder(w).Encode(result)

	return

}
*/

/*
func compareCOC(db xdr.Transaction, bc xdr.Transaction) apiModel.SubmitXDRSuccess {
	var result apiModel.SubmitXDRSuccess

	fmt.Println(strings.TrimRight(strconv.FormatInt(int64(db.Operations[3].Body.PaymentOp.Amount), 10), "0"))
	fmt.Println(strings.TrimLeft(fmt.Sprintf("%s", bc.Operations[4].Body.ManageDataOp.DataValue), "&"))

	if db.Operations[1].SourceAccount.Address() != bc.SourceAccount.Address() {
		result.Status = "Failed, Source Address in Gateway and Blockchain Doesn't match"
		return result

		// } else if strings.TrimLeft(fmt.Sprintf("%s", db.Operations[0].Body.ManageDataOp.DataValue), "&") !=
		// 	strings.TrimLeft(fmt.Sprintf("%s", bc.Operations[0].Body.ManageDataOp.DataValue), "&") {
		// 	result.Status = "Failed, Txn Type in Gateway and Blockchain Doesn't match"
	} else if strings.TrimLeft(fmt.Sprintf("%s", db.Operations[1].Body.ManageDataOp.DataValue), "&") !=
		strings.TrimLeft(fmt.Sprintf("%s", bc.Operations[1].Body.ManageDataOp.DataValue), "&") {
		result.Status = "Failed, Identifier in Gateway and Blockchain Doesn't match"
		return result

	} else if !strings.Contains(
		strings.TrimRight(fmt.Sprintf("%s", db.Operations[3].Body.PaymentOp.Asset.AlphaNum12.AssetCode), "\u0000"),
		strings.TrimLeft(fmt.Sprintf("%s", bc.Operations[3].Body.ManageDataOp.DataValue), "&")) {
		result.Status = "Failed, Asset Code in Gateway and Blockchain Doesn't match"
		return result

	} else if strings.TrimRight(strconv.FormatInt(int64(db.Operations[3].Body.PaymentOp.Amount), 10), "0") !=
		strings.TrimLeft(fmt.Sprintf("%s", bc.Operations[4].Body.ManageDataOp.DataValue), "&") {
		result.Status = "Failed, Asset Amount in Gateway and Blockchain Doesn't match"
		return result

	} else if db.Operations[3].Body.PaymentOp.Destination.Address() !=
		strings.TrimLeft(fmt.Sprintf("%s", bc.Operations[2].Body.ManageDataOp.DataValue), "&") {
		result.Status = "Failed, Destination Address in Gateway and Blockchain Doesn't match"
		return result

	} else {
		result.Status = "Success"

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
*/
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

type AbstractPOCOCNew struct {
	ProofHash  string
	Txn        string
	DBCOC      XDR
	BCCOC      XDR
	XDR        string
	COCStatus  string
	SequenceNo string
}

type XDR struct {
	SourceAccount string
	Identifier    string
	AssetCode     string
	AssetAmount   float64
	Destination   string
}

func (AP *AbstractPOCOCNew) InterpretPOCOCNew(w http.ResponseWriter, r *http.Request) {

	var result []model.POCOCResponse
	FromSigned := true
	ToSigned := false

	if AP.COCStatus != model.Pending.String() && AP.COCStatus != model.Expired.String() {
		ToSigned = true
	}

	SequenceNo, errSequence := strconv.ParseInt(AP.SequenceNo, 10, 64)
	if errSequence == nil {
		fmt.Printf("%d of type %T", SequenceNo, SequenceNo)
	}

	object := stellarRetriever.ConcretePOCOCNew{Txn: AP.ProofHash}

	bcCOC, state, timestamp, ledger, feePaid := object.RetrievePOCOCNew()

	mapD := map[string]string{"transaction": AP.Txn}
	mapB, err := json.Marshal(mapD)
	if err != nil {
		log.Error("Error while json.Marshal(mapD) " + err.Error())
	}	
	encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
	text := encoded
	if !state {

		w.WriteHeader(http.StatusOK)
		temp := model.POCOCResponse{
			Txnhash: AP.Txn,
			Url: commons.GetHorizonClient().URL + "/transactions/" +
				AP.Txn + "/operations",
			LabUrl:	commons.GetStellarLaboratoryClient() + "/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
			text + "%3D%3D&network=" + commons.GetHorizonClientNetworkName(),
			Identifier:     AP.DBCOC.Identifier,
			From:           AP.DBCOC.SourceAccount,
			To:             AP.DBCOC.Destination,
			Timestamp:      timestamp,
			Ledger:         ledger,
			FeePaid:        feePaid,
			TxnType:        "coc",
			SequenceNo:     SequenceNo,
			Quantity:       fmt.Sprintf("%f", AP.DBCOC.AssetAmount),
			AssetCode:      AP.DBCOC.AssetCode,
			Status:         "Failed to retrieve blockchain proof transaction from Stellar",
			BlockchainName: "Stellar",
			FromSigned:     FromSigned,
			ToSigned:       ToSigned,
			COCStatus:      AP.COCStatus}
		result = append(result, temp)
		json.NewEncoder(w).Encode(result)

		return

	}

	AP.BCCOC.SourceAccount = bcCOC.SourceAccount
	AP.BCCOC.Identifier = bcCOC.Identifier
	AP.BCCOC.AssetCode = bcCOC.AssetCode
	AP.BCCOC.AssetAmount = bcCOC.AssetAmount
	AP.BCCOC.Destination = bcCOC.Destination

	res := compareCOCNew(AP.DBCOC, AP.BCCOC)

	w.WriteHeader(http.StatusOK)

	temp := model.POCOCResponse{
		Txnhash: AP.Txn,
		Url: commons.GetHorizonClient().URL + "/transactions/" +
			AP.Txn + "/operations",
		LabUrl:	commons.GetStellarLaboratoryClient() + "/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
		text + "%3D%3D&network=" + commons.GetHorizonClientNetworkName(),
		Identifier:     AP.DBCOC.Identifier,
		From:           AP.DBCOC.SourceAccount,
		To:             AP.DBCOC.Destination,
		Timestamp:      timestamp,
		Ledger:         ledger,
		FeePaid:        feePaid,
		TxnType:        "coc",
		SequenceNo:     SequenceNo,
		Quantity:       fmt.Sprintf("%f", AP.DBCOC.AssetAmount),
		AssetCode:      AP.DBCOC.AssetCode,
		Status:         res.Status,
		BlockchainName: "Stellar",
		FromSigned:     FromSigned,
		ToSigned:       ToSigned,
		COCStatus:      AP.COCStatus}
	result = append(result, temp)
	json.NewEncoder(w).Encode(result)

	return
}

func compareCOCNew(db XDR, bc XDR) apiModel.SubmitXDRSuccess {

	var result apiModel.SubmitXDRSuccess

	if db.SourceAccount != bc.SourceAccount {
		result.Status = "Failed, Source Address in Gateway and Blockchain Doesn't match"
		return result

	} else if db.AssetCode != bc.AssetCode {
		result.Status = "Failed, Asset Code in Gateway and Blockchain Doesn't match"
		return result

	} else if db.AssetAmount != bc.AssetAmount {
		result.Status = "Failed, Asset Amount in Gateway and Blockchain Doesn't match"
		return result

	} else if db.Destination != bc.Destination {
		result.Status = "Failed, Destination Address in Gateway and Blockchain Doesn't match"
		return result

	} else if db.Identifier != bc.Identifier {
		result.Status = "Failed, Identifier in Gateway and Blockchain Doesn't match"
		return result

	} else {
		result.Status = "Success"
	}
	return result
}
