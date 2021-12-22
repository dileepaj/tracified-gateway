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

func GetProofPresentationProtocol(w http.ResponseWriter, r *http.Request){

	vars := mux.Vars(r)
	object := dao.Connection{}

	p := object.GetProofProtocolByProof(vars["proof"]) //calls in dao retrieve
	p.Then(func(data interface{}) interface{}{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
		return data
	}).Catch(func(err error) error {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Proof protocol not found in the Datastore",
		}
		json.NewEncoder(w).Encode(result)
		return err
	})
	p.Await()
}

func InsertProofPresentationProtocol(w http.ResponseWriter, r *http.Request){

	var newProtocolObj model.ProofPresentation

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
			Status: "Error when inserting protocol",
		}

		json.NewEncoder(w).Encode(result)

		return
	}else{
		fmt.Println(err1)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		
		result := apiModel.SubmitXDRSuccess{
			Status: "Protocol inserted",
		}

		json.NewEncoder(w).Encode(result)

		return
	}

}

func UpdateProofPresesntationProtocol(w http.ResponseWriter, r *http.Request){

	var Obj model.ProofPresentation
	var selection model.ProofPresentation

	err := json.NewDecoder(r.Body).Decode(&Obj)
	if err != nil{
		w.Header().Set("Content-Type", "application/json;")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while decoding the body")
		fmt.Println(err)
		return
	}
	//fmt.Println(Obj)

	object := dao.Connection{}

	_, err1 := object.GetProofProtocolByProof(Obj.ProofName).Then(func(data interface{}) interface{}{

		selection = data.(model.ProofPresentation)
		//fmt.Println("------Selection-------")
		//fmt.Println(selection)

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
				Status: "Protocol updated",
			}

		json.NewEncoder(w).Encode(result)
	}

	return data

	}).Await()

	if err1 != nil{
		
		w.Header().Set("Content-Type", "application/json;")
		w.WriteHeader(http.StatusBadRequest)

		result := apiModel.SubmitXDRSuccess{
			Status: "Error when fetching the protocol from DB or protodcol does not exists in the DB",
		}

		fmt.Println(err1)
		json.NewEncoder(w).Encode(result)

	}

}