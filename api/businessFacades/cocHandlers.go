package businessFacades

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/proofs/deprecatedBuilder"
	"github.com/sirupsen/logrus"

	"github.com/stellar/go/txnbuild"
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
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	object := dao.Connection{}

	data, err := object.GetCOCbySender(vars["Sender"]).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: err.Error(),
		}
		json.NewEncoder(w).Encode(result)
		return
	} else if data == nil {
		w.WriteHeader(http.StatusNotFound)
		result := apiModel.SubmitXDRSuccess{
			Status: "PublicKey Not Found in Gateway DataStore",
		}
		json.NewEncoder(w).Encode(result)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
		return
	}
}

/*GetCocByReceiver - WORKING MODEL
@author - Azeem Ashraf
@desc - Returns the COC Collection by querying the gateway DB by Receiver Public Key
@params - ResponseWriter,Request
*/
func GetCocByReceiver(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	object := dao.Connection{}
	data, err := object.GetCOCbyReceiver(vars["Receiver"]).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: err.Error(),
		}
		json.NewEncoder(w).Encode(result)
		return
	} else if data == nil {
		w.WriteHeader(http.StatusNotFound)
		result := apiModel.SubmitXDRSuccess{
			Status: "PublicKey Not Found in Gateway DataStore",
		}
		json.NewEncoder(w).Encode(result)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
		return
	}
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
		logrus.Error(err)
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
		logrus.Error(err)
	}

	brr,_ := txnbuild.TransactionFromXDR(GObj.AcceptXdr)

	t, _ := brr.Hash(commons.GetStellarNetwork())
	test := fmt.Sprintf("%x", t)

	err = xdr.SafeUnmarshalBase64(GObj.RejectXdr, &reject)
	if err != nil {
		logrus.Error(err)
	}

	brr1,_ := txnbuild.TransactionFromXDR(GObj.AcceptXdr)

	t1, _ := brr1.Hash(commons.GetStellarNetwork())
	test1 := fmt.Sprintf("%x", t1)

	var txe xdr.Transaction
	err1 := xdr.SafeUnmarshalBase64(GObj.AcceptXdr, &txe)
	if err1 != nil {
		logrus.Error(err1)
	}
	useSentSequence := false
	if len(txe.Operations) > 0 {
		for i := 0; i < len(txe.Operations); i++ {
			if txe.Operations[i].Body.Type == xdr.OperationTypeBumpSequence {
				v := fmt.Sprint(txe.Operations[i].Body.BumpSequenceOp.BumpTo)
				GObj.SequenceNo = v
				useSentSequence = true
			}
	 	}	
	}
	if useSentSequence {
		v := fmt.Sprint(txe.SeqNum)
		GObj.SequenceNo = v
	}
	GObj.AcceptTxn = test
	GObj.RejectTxn = test1
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
		return
	}
	object := dao.Connection{}
	switch GObj.Status {
	case model.Accepted.String():
		_, err := object.GetCOCbyAcceptTxn(GObj.AcceptTxn).Then(func(data interface{}) interface{} {
			selection = data.(model.COCCollectionBody)
			var TXNS []model.TransactionCollectionBody
			TXN := model.TransactionCollectionBody{
				XDR: GObj.AcceptXdr,
				TxnType: "10",
				Identifier: selection.Identifier,
				PublicKey: selection.SubAccount,
			}
			if selection.TenantID != "" {
				TXN.TenantID = selection.TenantID
			}
			TXNS = append(TXNS, TXN)
			status, response := builder.XDRSubmitter(TXNS)
			if !status {
				w.WriteHeader(502)
				errors_string := strings.ReplaceAll(response.Error.Message, "op_success? ", "")
				result := map[string]interface{}{
					"message":    "Failed to submit the Blockchain transaction",
					"error_code": errors_string,
				}
				json.NewEncoder(w).Encode(result)
			} else {
				GObj.TxnHash = response.TXNID
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
		}).Await()
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(err)
		}
		break
	case model.Rejected.String():
		_, err := object.GetCOCbyRejectTxn(GObj.RejectTxn).Then(func(data interface{}) interface{} {
			selection = data.(model.COCCollectionBody)
			display := &deprecatedBuilder.AbstractTDPInsert{XDR: GObj.RejectXdr}
			response := display.TDPInsert()
			if response.Error.Code == 400 {
				w.WriteHeader(502)
				errors_string := strings.ReplaceAll(response.Error.Message, "op_success? ", "")
				result := map[string]interface{}{
					"message":    "Failed to submit the Blockchain transaction",
					"error_code": errors_string,
				}
				json.NewEncoder(w).Encode(result)
			} else {
				GObj.TxnHash = response.TXNID
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
		}).Await()
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(err)
		}
		break
	default:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400)
		result = apiModel.InsertCOCCollectionResponse{
			Message: "Failed, Status invalid"}
		json.NewEncoder(w).Encode(result)
	}
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

func CheckAccountsStatusExtended(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var GObj apiModel.GetSubAccountStatus
	var result []apiModel.GetSubAccountStatusResponse
	var Coc apiModel.GetSubAccountStatusResponse
	var Org apiModel.GetSubAccountStatusResponse
	var Testimonial apiModel.GetSubAccountStatusResponse

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
			//result = append(result, data.(apiModel.GetSubAccountStatusResponse))
			Coc = data.(apiModel.GetSubAccountStatusResponse)
			return data
		}).Catch(func(error error) error {
			//result = append(result, apiModel.GetSubAccountStatusResponse{SubAccount: GObj.SubAccounts[i], Available: true})

			Coc = apiModel.GetSubAccountStatusResponse{SubAccount: GObj.SubAccounts[i], Available: true, Operation: "COC"}
			return error
		})
		p.Await()

		p2 := object.GetLastOrganizationbySubAccount(GObj.SubAccounts[i])
		p2.Then(func(data interface{}) interface{} {
			//result = append(result, data.(apiModel.GetSubAccountStatusResponse))
			Org = data.(apiModel.GetSubAccountStatusResponse)
			return data
		}).Catch(func(error error) error {
			//result = append(result, apiModel.GetSubAccountStatusResponse{SubAccount: GObj.SubAccounts[i], Available: true})

			Org = apiModel.GetSubAccountStatusResponse{SubAccount: GObj.SubAccounts[i], Available: true, Operation: "Organization"}
			return error
		})
		p2.Await()

		p3 := object.GetLastTestimonialbySubAccount(GObj.SubAccounts[i])
		p3.Then(func(data interface{}) interface{} {
			//result = append(result, data.(apiModel.GetSubAccountStatusResponse))
			Testimonial = data.(apiModel.GetSubAccountStatusResponse)
			return data
		}).Catch(func(error error) error {
			//result = append(result, apiModel.GetSubAccountStatusResponse{SubAccount: GObj.SubAccounts[i], Available: true})

			Testimonial = apiModel.GetSubAccountStatusResponse{SubAccount: GObj.SubAccounts[i], Available: true, Operation: "Testimonial"}
			return error
		})
		p3.Await()

		Cocseq, err := strconv.Atoi(Coc.SequenceNo)
		if err != nil {
			println(err)
		}
		Testimonialseq, err := strconv.Atoi(Testimonial.SequenceNo)
		if err != nil {
			println(err)
		}
		Orgseq, err := strconv.Atoi(Org.SequenceNo)
		if err != nil {
			println(err)
		}

		if Coc.Available && Org.Available && Testimonial.Available {

			if (Cocseq > Testimonialseq) && (Cocseq >= Orgseq) {
				result = append(result, apiModel.GetSubAccountStatusResponse{SequenceNo: strconv.Itoa(Cocseq), SubAccount: GObj.SubAccounts[i], Available: true, Operation: Coc.Operation, Receiver: Coc.Receiver, Expiration: Coc.Expiration})
			} else if (Testimonialseq > Cocseq) && (Testimonialseq > Orgseq) {
				result = append(result, apiModel.GetSubAccountStatusResponse{SequenceNo: strconv.Itoa(Testimonialseq), SubAccount: GObj.SubAccounts[i], Available: true, Operation: Testimonial.Operation, Receiver: Testimonial.Receiver, Expiration: Testimonial.Expiration})
			} else if (Orgseq > Cocseq) && (Orgseq > Testimonialseq) {
				result = append(result, apiModel.GetSubAccountStatusResponse{SequenceNo: strconv.Itoa(Orgseq), SubAccount: GObj.SubAccounts[i], Available: true, Operation: Org.Operation, Receiver: Org.Receiver, Expiration: Org.Expiration})
			} else {
				result = append(result, apiModel.GetSubAccountStatusResponse{SubAccount: GObj.SubAccounts[i], Available: true})
			}
		} else {
			if Org.Available == false {
				result = append(result, apiModel.GetSubAccountStatusResponse{SubAccount: GObj.SubAccounts[i], Available: false, Operation: Org.Operation, Receiver: Org.Receiver, SequenceNo: Org.SequenceNo, Expiration: Org.Expiration})
			}

			if Testimonial.Available == false {
				result = append(result, apiModel.GetSubAccountStatusResponse{SubAccount: GObj.SubAccounts[i], Available: false, Operation: Testimonial.Operation, Receiver: Testimonial.Receiver, SequenceNo: Testimonial.SequenceNo, Expiration: Testimonial.Expiration})
			}

			if Coc.Available == false {
				result = append(result, apiModel.GetSubAccountStatusResponse{SubAccount: GObj.SubAccounts[i], Available: false, Operation: Coc.Operation, Receiver: Coc.Receiver, SequenceNo: Coc.SequenceNo, Expiration: Coc.Expiration})
			}

		}

	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
	return

}
