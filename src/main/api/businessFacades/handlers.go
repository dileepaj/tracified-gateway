package businessFacades

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gorilla/mux"

	"main/api/apiModel"
	"main/model"
	"main/proofs/builder"
	"main/proofs/retriever/stellarRetriever"
)

//To be implemented
func SaveDataHash(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	result := builder.InsertTDP(vars["hash"], vars["secret"], vars["profileId"], vars["rootHash"])

	//test case
	// err1 := Error1{Code: 0, Message: "no root found"}
	// result := RootTree{Hash: "", Error: err1}

	//log the results
	fmt.Println(result, "result!!!")

	if result.Hash != "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return
	} else {
		// w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// switch result.Error.Code {
		// case 0:
		// 	w.WriteHeader(http.StatusNotFound)
		// 	json.NewEncoder(w).Encode(apiModel.JsonErr{StatusCode: http.StatusNotFound, Error: "No root"})
		// case 1:
		// 	w.WriteHeader(http.StatusNotFound)
		// 	json.NewEncoder(w).Encode(apiModel.JsonErr{StatusCode: http.StatusNotFound, Error: "Not Found"})
		// default:
		// 	w.WriteHeader(http.StatusNotFound)
		// 	json.NewEncoder(w).Encode(apiModel.JsonErr{StatusCode: http.StatusNotFound, Error: "Not Found"})
		// }

	}

}

//To be implemented
func CheckPOC(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)

	var response model.POC

	display := &stellarRetriever.ConcretePOC{Txn: vars["Txn"], ProfileID: vars["PID"], DBTree: vars["dbTree"]}
	response = display.InterpretPOC(display)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(response.RetrievePOC.Error.Code)
	result := apiModel.PoeSuccess{Message: response.RetrievePOC.Error.Message, TxNHash: response.RetrievePOC.Txn}
	json.NewEncoder(w).Encode(result)
	return

}

func CheckPOE(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var response model.POE

	display := &stellarRetriever.ConcretePOE{Txn: vars["Txn"], ProfileID: vars["PID"], Hash: vars["Hash"]}
	response = display.InterpretPOE(display)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(response.RetrievePOE.Error.Code)
	result := apiModel.PoeSuccess{Message: response.RetrievePOE.Error.Message, TxNHash: response.RetrievePOE.Txn}
	json.NewEncoder(w).Encode(result)
	return

}
