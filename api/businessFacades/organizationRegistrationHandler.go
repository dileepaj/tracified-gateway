package businessFacades

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/deprecatedBuilder"
	"github.com/dileepaj/tracified-gateway/utilities"
	//"github.com/go-openapi/runtime/logger"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

func InsertOrganization(w http.ResponseWriter, r *http.Request) {

	var Obj model.TestimonialOrganization
	b, err := ioutil.ReadAll(r.Body)
	strBody,_:=json.Marshal(r.Body)
	logger := utilities.NewCustomLogger()
	logger.LogWriter("request body   :"+string(strBody), constants.INFO)
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
	strObj,_:=json.Marshal(Obj)
	strPGPDataObj,_:=json.Marshal(Obj.PGPData)

	logger.LogWriter("Decoded object :"+string(strObj),constants.INFO)
	logger.LogWriter("PGP info   :"+string(strPGPDataObj),constants.INFO)

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
		logger.LogWriter("Error while calling SafeUnmarshalBase64(Obj.AcceptXDR, &accept)   :"+err.Error(), constants.ERROR)
	}

	acceptBuild, _ := txnbuild.TransactionFromXDR(Obj.AcceptXDR)

	acc, _ := acceptBuild.Hash(network.TestNetworkPassphrase)
	validAccept := fmt.Sprintf("%x", acc)

	err = xdr.SafeUnmarshalBase64(Obj.RejectXDR, &reject)
	if err != nil {
		logger.LogWriter("Error while calling SafeUnmarshalBase64(Obj.RejectXDR, &reject)  :"+err.Error(), constants.ERROR)
	}

	rejectBuild, _ := txnbuild.TransactionFromXDR(Obj.RejectXDR)

	rej, _ := rejectBuild.Hash(network.TestNetworkPassphrase)
	validReject := fmt.Sprintf("%x", rej)

	Obj.AcceptTxn = validAccept
	Obj.RejectTxn = validReject

	var txe xdr.Transaction
	err1 := xdr.SafeUnmarshalBase64(Obj.AcceptXDR, &txe)
	if err1 != nil {
		logger.LogWriter("Error occured while SafeUnmarshalBase64  :"+err1.Error(),constants.ERROR)
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
	logger := utilities.NewCustomLogger()
	b, err := ioutil.ReadAll(r.Body)
	strBody,_:=json.Marshal(r.Body)
	logger.LogWriter("request Body :"+string(strBody),constants.INFO)
	defer r.Body.Close()
	//err := json.NewDecoder(r.Body).Decode(&Obj)

	logger.LogWriter("-----------------end------------------------",constants.DEBUG)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while Decoding the body")
		logger.LogWriter("Error while Decoding the body"+err.Error(),constants.ERROR)
		return
	}
	errjson := json.Unmarshal(b, &Obj)
	if errjson != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	strObj,_:=json.Marshal(Obj)
	logger.LogWriter("Decoded Object :"+string(strObj),constants.INFO)

	object := dao.Connection{}
	logger.LogWriter("case STATUS :"+Obj.Status,constants.INFO)
	logger.LogWriter("APPROVED STATUS :"+model.Approved.String(),constants.INFO)
	switch Obj.Status {
	case model.Approved.String():

		_, err := object.GetOrganizationByAcceptTxn(Obj.AcceptTxn).Then(func(data interface{}) interface{} {
			selection = data.(model.TestimonialOrganization)
			display := &deprecatedBuilder.AbstractTDPInsert{XDR: Obj.AcceptXDR}
			response := display.TDPInsert()
			logger.LogWriter("Selection Sequence  :"+selection.SequenceNo,constants.INFO)
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
				logger.LogWriter("response.TXNID  :"+response.TXNID,constants.INFO)

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
					logger.LogWriter("Error occured while UpdateOrganizationInfo   :"+err2.Error(),constants.ERROR)

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
