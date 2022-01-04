package businessFacades

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/gorilla/mux"
)

func GetProofPresentationProtocolByProofName(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	object := dao.Connection{}
	
	p := object.GetProofProtocolByProofName(vars["proofname"]) 
	p.Then(func(data interface{}) interface{}{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
		return data
	}).Catch(func(err error) error {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Proof protocol cannot be found in Gateway datastore",
		}
		json.NewEncoder(w).Encode(result)
		return err
	})
	p.Await()
}

func InsertProofPresentationProtocol(w http.ResponseWriter, r *http.Request){
	var newProtocolObj model.ProofProtocol

	err := json.NewDecoder(r.Body).Decode(&newProtocolObj)
	if err != nil{
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error while decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		return
	}

	object := dao.Connection{}
	err1 := object.InsertProofProtocol(newProtocolObj)
	if err1 != nil{
		fmt.Println(err1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error when inserting protocol to the Datastore",
		}
		json.NewEncoder(w).Encode(result)
		return
	}else{
		fmt.Println(err1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Protocol inserted to the Datastore",
		}
		json.NewEncoder(w).Encode(result)
		return
	}
}

func UpdateProofPresentationProtocol(w http.ResponseWriter, r *http.Request){
	var Obj model.ProofProtocol
	var selection model.ProofProtocol

	err := json.NewDecoder(r.Body).Decode(&Obj)
	if err != nil{
		w.Header().Set("Content-Type", "application/json;")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while decoding the body")
		fmt.Println(err)
		return
	}

	object := dao.Connection{}
	_, err1 := object.GetProofProtocolByProofName(Obj.ProofName).Then(func(data interface{}) interface{}{
		selection = data.(model.ProofProtocol)
		err2 := object.UpdateProofPresesntationProtocol(selection, Obj)
		if err2 != nil{
			w.Header().Set("Content-Type", "application/json;")
			w.WriteHeader(http.StatusBadRequest)
			result := apiModel.SubmitXDRSuccess{
				Status: "Error when updating the protocol",
			}
			json.NewEncoder(w).Encode(result)
		}else{
			w.Header().Set("Content-Type", "application/json;")
			w.WriteHeader(http.StatusOK)
			result := apiModel.SubmitXDRSuccess{
				Status: "Protocol updated successfully",
			}
		json.NewEncoder(w).Encode(result)
	}
	return data
	}).Await()

	if err1 != nil{
		w.Header().Set("Content-Type", "application/json;")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error when fetching the protocol from Datastore or protocol does not exists in the Datastore",
		}
		fmt.Println(err1)
		json.NewEncoder(w).Encode(result)
	}
}

func DeleteProofPresentationProtocolByProofName(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	object := dao.Connection{}

	err := object.DeleteProofPresentationProtocolByProofName(vars["proofname"])
	if err != nil{
		w.Header().Set("Content-Type", "application/json;")
			w.WriteHeader(http.StatusBadRequest)
			result := apiModel.SubmitXDRSuccess{
				Status: "Error when deleting the protocol",
			}
			json.NewEncoder(w).Encode(result)
	}else{
		w.Header().Set("Content-Type", "application/json;")
			w.WriteHeader(http.StatusOK)
			result := apiModel.SubmitXDRSuccess{
				Status: "Protocol deleted successfully",
			}
		json.NewEncoder(w).Encode(result)
	}
}