package businessFacades

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/deprecatedBuilder"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

func InsertOrganization(w http.ResponseWriter, r *http.Request) {

	var Obj model.TestimonialOrganization
	b, err := ioutil.ReadAll(r.Body)
	log.Println("request body:", r.Body)
	defer r.Body.Close()
	//err := json.NewDecoder(r.Body).Decode(&Obj)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	errjson := json.Unmarshal(b, &Obj)
	if errjson != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	log.Println("Decoded Object: ", Obj)
	log.Println("PGP info", Obj.PGPData)
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

	acceptBuild, _ := txnbuild.TransactionFromXDR(Obj.AcceptXDR)

	acc, _ := acceptBuild.Hash(network.TestNetworkPassphrase)
	validAccept := fmt.Sprintf("%x", acc)

	err = xdr.SafeUnmarshalBase64(Obj.RejectXDR, &reject)
	if err != nil {
		fmt.Println(err)
	}

	rejectBuild, _ := txnbuild.TransactionFromXDR(Obj.RejectXDR)

	rej, _ := rejectBuild.Hash(network.TestNetworkPassphrase)
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
	b, err := ioutil.ReadAll(r.Body)
	log.Println("request body:", r.Body)
	defer r.Body.Close()
	//err := json.NewDecoder(r.Body).Decode(&Obj)

	log.Println("-----------------end------------------------")
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while Decoding the body")
		fmt.Println(err)
		return
	}
	errjson := json.Unmarshal(b, &Obj)
	if errjson != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	log.Println("Decoded obj: ", Obj)
	object := dao.Connection{}
	log.Println("case STATUS: ", Obj.Status)
	log.Println("APPROVED STATUS: ", model.Approved.String())
	switch Obj.Status {
	case model.Approved.String():

		_, err := object.GetOrganizationByAcceptTxn(Obj.AcceptTxn).Then(func(data interface{}) interface{} {
			selection = data.(model.TestimonialOrganization)
			display := &deprecatedBuilder.AbstractTDPInsert{XDR: Obj.AcceptXDR}
			response := display.TDPInsert()
			fmt.Println(selection.SequenceNo)
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

				err1 := object.UpdateOrganizationInfo(Obj)

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
				w.WriteHeader(502)
				errors_string := strings.ReplaceAll(response.Error.Message, "op_success? ", "")
				result := map[string]interface{}{
					"message":    "Failed to submit the Blockchain transaction",
					"error_code": errors_string,
				}
				json.NewEncoder(w).Encode(result)
			} else {
				Obj.TxnHash = response.TXNID
				err1 := object.UpdateOrganizationInfo(Obj)

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

func GetAllOrganizations_Paginated(w http.ResponseWriter, r *http.Request) {

	var response model.Error
	key1, error := r.URL.Query()["perPage"]

	if !error || len(key1[0]) < 1 {
		log.Error("Url Parameter 'perPage' is missing")
		return
	}

	key2, error := r.URL.Query()["page"]

	if !error || len(key2[0]) < 1 {
		log.Error("Url Parameter 'page' is missing")
		return
	}

	perPage, err := strconv.Atoi(key1[0])
	if err != nil {
		log.Error("Query parameter error" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		response = model.Error{Code: http.StatusBadRequest, Message: "The parameter should be an integer " + err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	page, err := strconv.Atoi(key2[0])
	if err != nil {
		log.Error("Query parameter error" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		response = model.Error{Code: http.StatusBadRequest, Message: "The parameter should be an integer " + err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	object := dao.Connection{}
	_, err = object.GetAllApprovedOrganizations_Paginated(perPage, page).Then(func(data interface{}) interface{} {
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
