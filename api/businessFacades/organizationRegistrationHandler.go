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

	var accept xdr.Transaction
	var reject xdr.Transaction

	err = xdr.SafeUnmarshalBase64(Obj.AcceptXDR, &accept)
	if err != nil {
		fmt.Println(err)
	}

	if commons.GoDotEnvVariable("NETWORKPASSPHRASE") == "test" {

		acceptBuild := build.TransactionBuilder{TX: &accept, NetworkPassphrase: build.TestNetwork.Passphrase}

		acc, _ := acceptBuild.Hash()
		validAccept := fmt.Sprintf("%x", acc)
		fmt.Println(acc)
		err = xdr.SafeUnmarshalBase64(Obj.RejectXDR, &reject)
		if err != nil {
			fmt.Println(err)
		}

		rejectBuild := build.TransactionBuilder{TX: &reject, NetworkPassphrase: build.TestNetwork.Passphrase}
		fmt.Println(build.TestNetwork.Passphrase)

		rej, _ := rejectBuild.Hash()
		validReject := fmt.Sprintf("%x", rej)

		Obj.AcceptTxn = validAccept
		Obj.RejectTxn = validReject

	} else if commons.GoDotEnvVariable("NETWORKPASSPHRASE") == "public" {
		acceptBuild := build.TransactionBuilder{TX: &accept, NetworkPassphrase: build.PublicNetwork.Passphrase}

		acc, _ := acceptBuild.Hash()
		validAccept := fmt.Sprintf("%x", acc)
		fmt.Println(acc)
		err = xdr.SafeUnmarshalBase64(Obj.RejectXDR, &reject)

		if err != nil {
			fmt.Println(err)
		}

		rejectBuild := build.TransactionBuilder{TX: &reject, NetworkPassphrase: build.PublicNetwork.Passphrase}
		fmt.Println(build.TestNetwork.Passphrase)

		rej, _ := rejectBuild.Hash()
		validReject := fmt.Sprintf("%x", rej)

		Obj.AcceptTxn = validAccept
		Obj.RejectTxn = validReject
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
		w.WriteHeader(http.StatusOK)
		result := model.TestimonialOrganizationResponse{
			AcceptTxn: Obj.AcceptTxn,
			AcceptXDR: Obj.AcceptXDR,
			RejectTxn: Obj.RejectTxn,
			RejectXDR: Obj.RejectXDR,
			Status:    Obj.Status}
		json.NewEncoder(w).Encode(result)
		return
	}

}

func GetAllOrganizations(w http.ResponseWriter, r *http.Request) {
	object := dao.Connection{}
	p := object.GetAllApprovedOrganizations()
	p.Then(func(data interface{}) interface{} {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
		return data
	}).Catch(func(error error) error {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "No Approved organizations were found",
		}
		json.NewEncoder(w).Encode(result)
		return error
	})
	p.Await()
}

func GetOrganizationByPublicKey(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	object := dao.Connection{}
	p := object.GetOrganizationByAuthor(vars["PK"])
	p.Then(func(data interface{}) interface{} {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
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
	case "Accepted":

		p := object.GetOrganizationByAcceptTxn(Obj.AcceptTxn)
		p.Then(func(data interface{}) interface{} {
			selection = data.(model.TestimonialOrganization)
			display := &deprecatedBuilder.AbstractTDPInsert{XDR: Obj.AcceptXDR}
			response := display.TDPInsert()

			if response.Error.Code == 400 {
				w.WriteHeader(400)
				result := apiModel.SubmitXDRSuccess{
					Status: "Failed"}
				json.NewEncoder(w).Encode(result)
			} else {

				Obj.TxnHash = response.TXNID
				fmt.Println(response.TXNID)

				err1 := object.UpdateOrganization(selection, Obj)
				if err1 != nil {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(400)
					result := apiModel.SubmitXDRSuccess{
						Status: "Failed"}
					json.NewEncoder(w).Encode(result)

				} else {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusOK)

					result := apiModel.SubmitXDRSuccess{
						Status: "Success"}
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
	case "Rejected":
		p := object.GetOrganizationByRejectTxn(Obj.RejectTxn)
		p.Then(func(data interface{}) interface{} {

			selection = data.(model.TestimonialOrganization)

			display := &deprecatedBuilder.AbstractTDPInsert{XDR: Obj.RejectXDR}
			response := display.TDPInsert()

			if response.Error.Code == 400 {
				w.WriteHeader(400)
				result := apiModel.SubmitXDRSuccess{
					Status: "Failed"}
				json.NewEncoder(w).Encode(result)
			} else {
				Obj.TxnHash = response.TXNID
				err1 := object.UpdateOrganization(selection, Obj)
				if err1 == nil {

					result := apiModel.SubmitXDRSuccess{Status: "Success"}
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
		break

	default:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400)
		result := apiModel.InsertCOCCollectionResponse{
			Message: "Failed, Status invalid"}
		json.NewEncoder(w).Encode(result)
	}

	return
}
