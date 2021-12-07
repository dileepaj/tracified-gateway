package businessFacades

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/deprecatedBuilder"
	"github.com/gorilla/mux"
	"github.com/stellar/go/build"
	"github.com/stellar/go/xdr"
)

func InsertOrganization(w http.ResponseWriter, r *http.Request) {

	var Obj model.TestimonialOrganization

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

	rej, _ := rejectBuild.Hash()
	validReject := fmt.Sprintf("%x", rej)

	Obj.AcceptTxn = validAccept
	Obj.RejectTxn = validReject

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
	err2 := object.InsertOrganization(Obj)

	if err2 != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.InsertorganizationCollectionResponse{
			Message: "Failed"}
		json.NewEncoder(w).Encode(result)
		return
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(201)
		result := model.TestimonialOrganizationResponse{
			AcceptTxn:  Obj.AcceptTxn,
			AcceptXDR:  Obj.AcceptXDR,
			RejectTxn:  Obj.RejectTxn,
			RejectXDR:  Obj.RejectXDR,
			SequenceNo: Obj.SequenceNo,
			Status:     Obj.Status}
		json.NewEncoder(w).Encode(result)
		return
	}

}

func GetAllOrganizations(w http.ResponseWriter, r *http.Request) {

	object := dao.Connection{}

	_, err := object.GetAllApprovedOrganizations().Then(func(data interface{}) interface{} {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
		return data
	}).Await()
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(204)
		result := apiModel.SubmitXDRSuccess{
			Status: "No Approved organizations were found",
		}
		json.NewEncoder(w).Encode(result)
	}

}

func GetOrganizationByPublicKey(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	object := dao.Connection{}

	_, err := object.GetOrganizationByAuthor(vars["PK"]).Then(func(data interface{}) interface{} {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
		return data
	}).Await()
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(204)
		result := apiModel.SubmitXDRSuccess{
			Status: "PublicKey Not Found in Gateway DataStore",
		}
		json.NewEncoder(w).Encode(result)
	}
}

func UpdateOrganization(w http.ResponseWriter, r *http.Request) {

	var Obj model.TestimonialOrganization
	var selection model.TestimonialOrganization

	err := json.NewDecoder(r.Body).Decode(&Obj)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while Decoding the body")
		fmt.Println(err)
		return
	}
	fmt.Println(Obj)
	object := dao.Connection{}

	switch Obj.Status {
	case model.Approved.String():

		_, err := object.GetOrganizationByAcceptTxn(Obj.AcceptTxn).Then(func(data interface{}) interface{} {
			selection = data.(model.TestimonialOrganization)
			display := &deprecatedBuilder.AbstractTDPInsert{XDR: Obj.AcceptXDR}
			response := display.TDPInsert()
			fmt.Println(selection.SequenceNo)
			if response.Error.Code == 400 {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(400)
				result := apiModel.SubmitXDRSuccess{
					Status: "Failed"}
				json.NewEncoder(w).Encode(result)
			} else {

				Obj.TxnHash = response.TXNID
				fmt.Println(response.TXNID)

				err1 := object.Updateorganization(selection, Obj)

				if err1 != nil {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(400)
					result := apiModel.SubmitXDRSuccess{
						Status: "Failed"}
					json.NewEncoder(w).Encode(result)

				} else {
					result := model.TestimonialOrganizationResponse{
						AcceptTxn:  selection.AcceptTxn,
						AcceptXDR:  Obj.AcceptXDR,
						RejectTxn:  selection.RejectTxn,
						RejectXDR:  Obj.RejectXDR,
						SequenceNo: selection.SequenceNo,
						TxnHash:    response.TXNID,
						Status:     Obj.Status}

					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(200)
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

		_, err := object.GetOrganizationByRejectTxn(Obj.RejectTxn).Then(func(data interface{}) interface{} {

			selection = data.(model.TestimonialOrganization)

			display := &deprecatedBuilder.AbstractTDPInsert{XDR: Obj.RejectXDR}
			response := display.TDPInsert()

			if response.Error.Code == 400 {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(400)
				result := apiModel.SubmitXDRSuccess{
					Status: "Failed"}
				json.NewEncoder(w).Encode(result)
			} else {
				Obj.TxnHash = response.TXNID
				err1 := object.Updateorganization(selection, Obj)

				if err1 == nil {

					result := model.TestimonialOrganizationResponse{
						AcceptTxn:  selection.AcceptTxn,
						AcceptXDR:  Obj.AcceptXDR,
						RejectTxn:  selection.RejectTxn,
						RejectXDR:  Obj.RejectXDR,
						TxnHash:    response.TXNID,
						SequenceNo: selection.SequenceNo,
						Status:     Obj.Status}

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
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(err)
		}
		break

	default:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400)
		result := apiModel.SubmitXDRSuccess{
			Status: "Failed, Status invalid"}
		json.NewEncoder(w).Encode(result)
	}

	return
}

func GetAllPendingAndRejectedOrganizations(w http.ResponseWriter, r *http.Request) {

	object := dao.Connection{}

	_, err := object.GetPendingAndRejectedOrganizations().Then(func(data interface{}) interface{} {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
		return data
	}).Await()
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(204)
		result := apiModel.SubmitXDRSuccess{
			Status: "No Pending or Rejected organizations were found",
		}
		json.NewEncoder(w).Encode(result)
	}
}
