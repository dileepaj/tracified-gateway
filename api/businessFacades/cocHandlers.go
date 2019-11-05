package businessFacades

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dileepaj/tracified-gateway/proofs/deprecatedBuilder"

	"github.com/stellar/go/build"
	"github.com/stellar/go/xdr"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/builder"
	"github.com/gorilla/mux"
)

/*GetCocBySender - WORKING MODEL
@author - Azeem Ashraf
@desc - Returns the COC Collection by querying the gateway DB by Sender Public Key
@params - ResponseWriter,Request
*/
func GetCocBySender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	object := dao.Connection{}
	p := object.GetCOCbySender(vars["Sender"])
	p.Then(func(data interface{}) interface{} {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		// result := apiModel.GetMultiCOCCollection{
		// 	Collection: data}
		json.NewEncoder(w).Encode(data)
		return data
	}).Catch(func(error error) error {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "PublicKey Not Found in Gateway DataStore",
		}
		json.NewEncoder(w).Encode(result)
		return error
	})
	p.Await()

}

/*GetCocByReceiver - WORKING MODEL
@author - Azeem Ashraf
@desc - Returns the COC Collection by querying the gateway DB by Receiver Public Key
@params - ResponseWriter,Request
*/
func GetCocByReceiver(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	object := dao.Connection{}
	p := object.GetCOCbyReceiver(vars["Receiver"])
	p.Then(func(data interface{}) interface{} {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		// result := apiModel.GetMultiCOCCollection{
		// 	Collection: data}
		json.NewEncoder(w).Encode(data)
		return data
	}).Catch(func(error error) error {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "PublicKey Not Found in Gateway DataStore",
		}
		json.NewEncoder(w).Encode(result)
		return error
	})
	p.Await()

}

/*InsertCocCollection - WORKING MODEL
@author - Azeem Ashraf
@desc - Inserts a COC Collection received by the wallet application with sender's signature
@params - ResponseWriter,Request
*/
func InsertCocCollection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var GObj model.COCCollectionBody
	err := json.NewDecoder(r.Body).Decode(&GObj)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error while Decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		return
	}
	var accept xdr.Transaction
	var reject xdr.Transaction
	err = xdr.SafeUnmarshalBase64(GObj.AcceptXdr, &accept)
	if err != nil {
		fmt.Println(err)
	}

	brr := build.TransactionBuilder{TX: &accept, NetworkPassphrase: build.PublicNetwork.Passphrase}
	fmt.Println(build.PublicNetwork.Passphrase)
	t, _ := brr.Hash()
	test := fmt.Sprintf("%x", t)

	err = xdr.SafeUnmarshalBase64(GObj.RejectXdr, &reject)
	if err != nil {
		fmt.Println(err)
	}

	brr1 := build.TransactionBuilder{TX: &reject, NetworkPassphrase: build.PublicNetwork.Passphrase}
	fmt.Println(build.PublicNetwork.Passphrase)
	t1, _ := brr1.Hash()
	test1 := fmt.Sprintf("%x", t1)

	var txe xdr.Transaction
	err1 := xdr.SafeUnmarshalBase64(GObj.AcceptXdr, &txe)
	if err1 != nil {
		fmt.Println(err1)
	}
	useSentSequence := false

	for i := 0; i < len(txe.Operations); i++ {

		if txe.Operations[i].Body.Type == xdr.OperationTypeBumpSequence {
			fmt.Println("HAHAHAHA BUMPY")
			v := fmt.Sprint(txe.Operations[i].Body.BumpSequenceOp.BumpTo)
			fmt.Println(v)
			GObj.SequenceNo = v
			useSentSequence = true

		}
	}
	if !useSentSequence {
		fmt.Println("seq")
		fmt.Println(txe.SeqNum)
		v := fmt.Sprint(txe.SeqNum)
		GObj.SequenceNo = v
	}
	fmt.Println("SubAcc")
	fmt.Println(GObj.SubAccount)
	GObj.AcceptTxn = test
	GObj.RejectTxn = test1
	fmt.Println(GObj)
	object := dao.Connection{}
	err2 := object.InsertCoc(GObj)

	if err2 != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.InsertCOCCollectionResponse{
			Message: "Failed"}
		json.NewEncoder(w).Encode(result)
		return
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		result := apiModel.InsertCOCCollectionResponse{
			Message: "Success", Body: GObj}
		json.NewEncoder(w).Encode(result)
		return
	}
}

/*UpdateCocCollection - WORKING MODEL
@author - Azeem Ashraf
@desc - Handles the Proof of Existance by retrieving the Raw Data from the Traceability Data Store
and Retrieves the TXN ID and calls POE Interpreter
Finally Returns the Response given by the POE Interpreter
@params - ResponseWriter,Request
*/
func UpdateCocCollection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var GObj model.COCCollectionBody
	var selection model.COCCollectionBody
	var result apiModel.InsertCOCCollectionResponse

	err := json.NewDecoder(r.Body).Decode(&GObj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while Decoding the body")
		fmt.Println(err)
		return
	}
	fmt.Println(GObj)
	object := dao.Connection{}
	switch GObj.Status {
	case "accepted":
		p := object.GetCOCbyAcceptTxn(GObj.AcceptTxn)
		p.Then(func(data interface{}) interface{} {

			selection = data.(model.COCCollectionBody)

			var TXNS []model.TransactionCollectionBody
			TXN := model.TransactionCollectionBody{
				XDR: GObj.AcceptXdr,
			}
			TXNS = append(TXNS, TXN)
			fmt.Println(TXNS)
			status, response := builder.XDRSubmitter(TXNS)

			if !status {
				w.WriteHeader(400)
				result = apiModel.InsertCOCCollectionResponse{
					Message: "Failed"}
				json.NewEncoder(w).Encode(result)
			} else {

				GObj.TxnHash = response.TXNID
				fmt.Println(response.TXNID)

				err1 := object.UpdateCOC(selection, GObj)
				if err1 != nil {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(400)
					result = apiModel.InsertCOCCollectionResponse{
						Message: "Failed"}
					json.NewEncoder(w).Encode(result)

				} else {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusOK)
					body := GObj
					body.AcceptTxn = GObj.AcceptTxn
					body.AcceptXdr = GObj.AcceptXdr
					body.Status = GObj.Status
					result = apiModel.InsertCOCCollectionResponse{
						Message: "Success", Body: body}
					json.NewEncoder(w).Encode(result)
				}
			}
			return data
		}).Catch(func(error error) error {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(error)
			return error
		})
		p.Await()
		break
	case "rejected":
		p := object.GetCOCbyRejectTxn(GObj.RejectTxn)
		p.Then(func(data interface{}) interface{} {
			selection = data.(model.COCCollectionBody)
			display := &deprecatedBuilder.AbstractTDPInsert{XDR: GObj.RejectXdr}
			response := display.TDPInsert()

			if response.Error.Code == 400 {
				w.WriteHeader(400)
				result = apiModel.InsertCOCCollectionResponse{
					Message: "Failed"}
				json.NewEncoder(w).Encode(result)
			} else {
				GObj.TxnHash = response.TXNID
				fmt.Println(response.TXNID)
				err1 := object.UpdateCOC(selection, GObj)
				if err1 != nil {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(400)
					result = apiModel.InsertCOCCollectionResponse{
						Message: "Failed"}
					json.NewEncoder(w).Encode(result)

				} else {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusOK)
					body := GObj
					body.RejectTxn = GObj.RejectTxn
					body.RejectXdr = GObj.RejectXdr
					body.Status = GObj.Status
					result = apiModel.InsertCOCCollectionResponse{
						Message: "Success", Body: body}
					json.NewEncoder(w).Encode(result)

				}
			}

			return data
		}).Catch(func(error error) error {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(error)
			return error
		})
		p.Await()
		break

	default:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400)
		result = apiModel.InsertCOCCollectionResponse{
			Message: "Failed, Status invalid"}
		json.NewEncoder(w).Encode(result)
	}
	return
}

/*CheckAccountsStatus - WORKING MODEL
@author - Azeem Ashraf
@desc - Checks all the available COCs in the gateway datastore
and retrieves them by the sender's publickey and returns the status and sequence number.
@params - ResponseWriter,Request
*/
func CheckAccountsStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var GObj apiModel.GetSubAccountStatus
	var result []apiModel.GetSubAccountStatusResponse

	err := json.NewDecoder(r.Body).Decode(&GObj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while Decoding the body")
		fmt.Println(err)
		return
	}
	fmt.Println(GObj)
	object := dao.Connection{}
	for i := 0; i < len(GObj.SubAccounts); i++ {

		p := object.GetLastCOCbySubAccount(GObj.SubAccounts[i])
		p.Then(func(data interface{}) interface{} {
			result = append(result, data.(apiModel.GetSubAccountStatusResponse))
			return data
		}).Catch(func(error error) error {
			result = append(result, apiModel.GetSubAccountStatusResponse{SubAccount: GObj.SubAccounts[i], Available: true})

			return error
		})
		p.Await()

	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
	return

}

/*LastCOC - WORKING MODEL
@author - Azeem Ashraf
@desc - Returns the Txn ID of the last COC Txn
@params - ResponseWriter,Request
*/
func LastCOC(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	object := dao.Connection{}
	p := object.GetLastCOCbyIdentifier(vars["Identifier"])
	p.Then(func(data interface{}) interface{} {

		result := data.(model.COCCollectionBody)
		// res := model.LastTxnResponse{LastTxn: result.TxnHash}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return nil
	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Identifier Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error
	})
	p.Await()

}
