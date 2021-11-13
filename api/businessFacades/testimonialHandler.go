package businessFacades

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/deprecatedBuilder"
	"github.com/gorilla/mux"
	"github.com/stellar/go/build"
	"github.com/stellar/go/xdr"
)

func InsertTestimonial(w http.ResponseWriter, r *http.Request) {

	var Obj model.Testimonial
	err := json.NewDecoder(r.Body).Decode(&Obj)
	if err != nil {
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error while Decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		return
	}

	if Obj.Status != model.Pending.String() {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "invalid Status",
		}
		json.NewEncoder(w).Encode(result)
		return
	}

	var accept xdr.Transaction
	var reject xdr.Transaction

	if commons.GoDotEnvVariable("NETWORKPASSPHRASE") == "test" {

		err = xdr.SafeUnmarshalBase64(Obj.AcceptXDR, &accept)
		if err != nil {
			fmt.Println(err)
		}

		acceptBuild := build.TransactionBuilder{TX: &accept, NetworkPassphrase: commons.GetHorizonNetwork().Passphrase}

		acc, _ := acceptBuild.Hash()
		validAccept := fmt.Sprintf("%x", acc)

		err = xdr.SafeUnmarshalBase64(Obj.RejectXDR, &reject)
		if err != nil {
			fmt.Println(err)
		}

		rejectBuild := build.TransactionBuilder{TX: &reject, NetworkPassphrase: commons.GetHorizonNetwork().Passphrase}
		fmt.Println(commons.GetHorizonNetwork().Passphrase)

		rej, _ := rejectBuild.Hash()
		validReject := fmt.Sprintf("%x", rej)

		Obj.AcceptTxn = validAccept
		Obj.RejectTxn = validReject

	} else if commons.GoDotEnvVariable("NETWORKPASSPHRASE") == "public" {

		err = xdr.SafeUnmarshalBase64(Obj.AcceptXDR, &accept)
		if err != nil {
			fmt.Println(err)
		}

		acceptBuild := build.TransactionBuilder{TX: &accept, NetworkPassphrase: commons.GetHorizonNetwork().Passphrase}

		acc, _ := acceptBuild.Hash()
		validAccept := fmt.Sprintf("%x", acc)

		err = xdr.SafeUnmarshalBase64(Obj.RejectXDR, &reject)

		if err != nil {
			fmt.Println(err)
		}

		rejectBuild := build.TransactionBuilder{TX: &reject, NetworkPassphrase: commons.GetHorizonNetwork().Passphrase}
		fmt.Println(commons.GetHorizonNetwork().Passphrase)

		rej, _ := rejectBuild.Hash()
		validReject := fmt.Sprintf("%x", rej)

		Obj.AcceptTxn = validAccept
		Obj.RejectTxn = validReject
	}

	var txe xdr.Transaction
	err1 := xdr.SafeUnmarshalBase64(Obj.AcceptXDR, &txe)
	if err1 != nil {
		fmt.Println(err1)
	}
	useSentSequence := false

	for i := 0; i < len(txe.Operations); i++ {

		if txe.Operations[i].Body.Type == xdr.OperationTypeBumpSequence {

			v := fmt.Sprint(txe.Operations[i].Body.BumpSequenceOp.BumpTo)

			Obj.SequenceNo = v
			useSentSequence = true

		}
	}
	if !useSentSequence {

		v := fmt.Sprint(txe.SeqNum)
		Obj.SequenceNo = v
	}

	object := dao.Connection{}
	err2 := object.InsertTestimonial(Obj)

	if err2 != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{Status: "Failed"}
		json.NewEncoder(w).Encode(result)
		return
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(201)
		result := model.TestimonialResponse{
			SequenceNo:  Obj.SequenceNo,
			Status:      Obj.Status,
			Testimonial: Obj.Testimonial}
		json.NewEncoder(w).Encode(result)
		return
	}

}

func GetTestimonialBySender(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	object := dao.Connection{}

	_, err := object.GetTestimonialBySenderPublickey(vars["PK"]).Then(func(data interface{}) interface{} {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
		return data
	}).Await()
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(204)
		result := apiModel.SubmitXDRSuccess{
			Status: "Sender PublicKey Not Found in Gateway DataStore",
		}
		json.NewEncoder(w).Encode(result)
	}
}

func GetTestimonialByReciever(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	object := dao.Connection{}

	_, err := object.GetTestimonialByRecieverPublickey(vars["PK"]).Then(func(data interface{}) interface{} {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
		return data
	}).Await()

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(204)
		result := apiModel.SubmitXDRSuccess{
			Status: "Reciever PublicKey Not Found in Gateway DataStore",
		}
		json.NewEncoder(w).Encode(result)
	}

}

func UpdateTestimonial(w http.ResponseWriter, r *http.Request) {
	var Obj model.Testimonial
	var selection model.Testimonial

	err := json.NewDecoder(r.Body).Decode(&Obj)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while Decoding the body")
		fmt.Println(err)
		return
	}
	object := dao.Connection{}

	switch Obj.Status {
	case model.Approved.String():

		_, err := object.GetTestimonialByAcceptTxn(Obj.AcceptTxn).Then(func(data interface{}) interface{} {
			selection = data.(model.Testimonial)
			display := &deprecatedBuilder.AbstractTDPInsert{XDR: Obj.AcceptXDR}
			response := display.TDPInsert()

			if response.Error.Code == 400 {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(502)
				errors_string := strings.ReplaceAll(response.Error.Message, "op_success? ", "")
				result := map[string]interface{}{
					"message":    "Failed to submit the Blockchain transaction",
					"error_code": errors_string,
				}
				json.NewEncoder(w).Encode(result)
			} else {

				Obj.TxnHash = response.TXNID
				fmt.Println(response.TXNID)

				err1 := object.UpdateTestimonial(selection, Obj)
				if err1 != nil {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(400)
					result := apiModel.SubmitXDRSuccess{
						Status: "Failed"}
					json.NewEncoder(w).Encode(result)

				} else {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusOK)

					result := model.TestimonialResponse{
						SequenceNo:  selection.SequenceNo,
						TxnHash:     response.TXNID,
						Status:      Obj.Status,
						Testimonial: selection.Testimonial,
					}
					json.NewEncoder(w).Encode(result)
				}
			}
			return data
		}).Await()

		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(204)
			result := apiModel.InsertTestimonialCollectionResponse{
				Message: "Error while fetch data from db or AcceptTxn Not exist in DB",
			}
			json.NewEncoder(w).Encode(result)
		}
		break
	case model.Rejected.String():

		_, err := object.GetTestimonialByRejectTxn(Obj.RejectTxn).Then(func(data interface{}) interface{} {

			selection = data.(model.Testimonial)

			display := &deprecatedBuilder.AbstractTDPInsert{XDR: Obj.RejectXDR}
			response := display.TDPInsert()

			if response.Error.Code == 400 {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(502)
				errors_string := strings.ReplaceAll(response.Error.Message, "op_success? ", "")
				result := map[string]interface{}{
					"message":    "Failed to submit the Blockchain transaction",
					"error_code": errors_string,
				}
				json.NewEncoder(w).Encode(result)
			} else {
				Obj.TxnHash = response.TXNID
				err1 := object.UpdateTestimonial(selection, Obj)
				if err1 == nil {

					result := model.TestimonialResponse{
						SequenceNo:  selection.SequenceNo,
						TxnHash:     response.TXNID,
						Status:      Obj.Status,
						Testimonial: selection.Testimonial,
					}
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusOK)
					err2 := json.NewEncoder(w).Encode(result)
					fmt.Println(err2)

				} else {

					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(400)
					result := apiModel.SubmitXDRSuccess{Status: "Failed"}
					json.NewEncoder(w).Encode(result)

				}
			}
			return data
		}).Await()
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(204)
			result := apiModel.InsertTestimonialCollectionResponse{
				Message: "Error while fetch data from db or RejectTxn Not exist in DB",
			}
			json.NewEncoder(w).Encode(result)
		}
		break
	default:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400)
		result := apiModel.InsertTestimonialCollectionResponse{
			Message: "Failed, Status invalid"}
		json.NewEncoder(w).Encode(result)
	}

	return

}
