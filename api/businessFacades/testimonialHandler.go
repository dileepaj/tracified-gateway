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

	var accept xdr.Transaction
	var reject xdr.Transaction

	if commons.GoDotEnvVariable("NETWORKPASSPHRASE") == "test" {

		err = xdr.SafeUnmarshalBase64(Obj.AcceptXDR, &accept)
		if err != nil {
			fmt.Println(err)
		}

		acceptBuild := build.TransactionBuilder{TX: &accept, NetworkPassphrase: build.TestNetwork.Passphrase}

		acc, _ := acceptBuild.Hash()
		validAccept := fmt.Sprintf("%x", acc)

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

		err = xdr.SafeUnmarshalBase64(Obj.AcceptXDR, &accept)
		if err != nil {
			fmt.Println(err)
		}

		acceptBuild := build.TransactionBuilder{TX: &accept, NetworkPassphrase: build.PublicNetwork.Passphrase}

		acc, _ := acceptBuild.Hash()
		validAccept := fmt.Sprintf("%x", acc)

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
	err2 := object.InsertTestimonial(Obj)

	if err2 != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{Status: "Failed"}
		json.NewEncoder(w).Encode(result)
		return
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		result := model.TestimonialResponse{
			Status:      Obj.Status,
			Testimonial: Obj.Testimonial}
		json.NewEncoder(w).Encode(result)
		return
	}

}

func GetTestimonialBySender(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	object := dao.Connection{}
	p := object.GetTestimonialBySenderPublickey(vars["PK"])
	p.Then(func(data interface{}) interface{} {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
		return data
	}).Catch(func(error error) error {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Sender PublicKey Not Found in Gateway DataStore",
		}
		json.NewEncoder(w).Encode(result)
		return error
	})
	p.Await()

}

func GetTestimonialByReciever(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	object := dao.Connection{}
	p := object.GetTestimonialByRecieverPublickey(vars["PK"])
	p.Then(func(data interface{}) interface{} {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
		return data
	}).Catch(func(error error) error {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Reciever PublicKey Not Found in Gateway DataStore",
		}
		json.NewEncoder(w).Encode(result)
		return error
	})
	p.Await()

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
	case "Approved":

		p := object.GetTestimonialByAcceptTxn(Obj.AcceptTxn)
		p.Then(func(data interface{}) interface{} {
			selection = data.(model.Testimonial)
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
						TxnHash:     response.TXNID,
						Status:      Obj.Status,
						Testimonial: Obj.Testimonial,
					}
					json.NewEncoder(w).Encode(result)
				}
			}
			return data
		}).Catch(func(error error) error {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(400)
			result := apiModel.InsertTestimonialCollectionResponse{
				Message: "Error while fetch data from db or AcceptTxn Not exist in DB",
			}
			json.NewEncoder(w).Encode(result)
			return error
		})
		p.Await()
		break
	case "Rejected":
		p := object.GetTestimonialByRejectTxn(Obj.RejectTxn)
		p.Then(func(data interface{}) interface{} {

			selection = data.(model.Testimonial)

			display := &deprecatedBuilder.AbstractTDPInsert{XDR: Obj.RejectXDR}
			response := display.TDPInsert()

			if response.Error.Code == 400 {
				w.WriteHeader(400)
				result := apiModel.SubmitXDRSuccess{
					Status: "Failed"}
				json.NewEncoder(w).Encode(result)
			} else {
				Obj.TxnHash = response.TXNID
				err1 := object.UpdateTestimonial(selection, Obj)
				if err1 == nil {

					result := model.TestimonialResponse{
						TxnHash:     response.TXNID,
						Status:      Obj.Status,
						Testimonial: Obj.Testimonial,
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
		}).Catch(func(error error) error {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			result := apiModel.InsertTestimonialCollectionResponse{
				Message: "Error while fetch data from db or RejectTxn Not exist in DB",
			}
			json.NewEncoder(w).Encode(result)
			return error
		})
		p.Await()
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
