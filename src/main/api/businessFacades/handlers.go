package businessFacades

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gorilla/mux"

	"main/api/apiModel"
	"main/proofs/builder"
	"main/proofs/interpreter"
	"main/proofs/retriever/stellarRetriever"
)

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
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		switch result.Error.Code {
		case 0:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(apiModel.JsonErr{StatusCode: http.StatusNotFound, Error: "No root"})
		case 1:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(apiModel.JsonErr{StatusCode: http.StatusNotFound, Error: "Not Found"})
		default:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(apiModel.JsonErr{StatusCode: http.StatusNotFound, Error: "Not Found"})
		}

	}

}

func CheckPOC(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	result := interpreter.InterpretPOC(vars["rootHash"], vars["treeObj"])

	//log the results
	fmt.Println(result, "result!!!")

	// if result.Previous != "" && result.Current != "" {
	// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// 	w.WriteHeader(http.StatusOK)
	// 	json.NewEncoder(w).Encode(result)
	// 	return
	// }
	// else {
	// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// 	switch result.Error.Code {
	// 	case 0:
	// 		w.WriteHeader(http.StatusNotFound)
	// 		json.NewEncoder(w).Encode(apiModel.JsonErr{Code: http.StatusNotFound, Text: "No PRE"})
	// 	case 1:
	// 		w.WriteHeader(http.StatusNotFound)
	// 		json.NewEncoder(w).Encode(apiModel.JsonErr{Code: http.StatusNotFound, Text: "Not Found"})
	// 	default:
	// 		w.WriteHeader(http.StatusNotFound)
	// 		json.NewEncoder(w).Encode(apiModel.JsonErr{Code: http.StatusNotFound, Text: "Not Found"})
	// 	}

	// }

}

func CheckPOE(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// display := &stellarRetriever.ConcretePOE{Txn: "e903f5ef813002295e97c0f08cf26d1fd411615e18384890395f6b0943ed83b5", ProfileID: "ProfileID001", Hash: "cf68e34967e10837d629b941bb8ec85d0ef016bc324340bd54e0ccae08a30b7a"}
	display := &stellarRetriever.ConcretePOE{Txn: vars["Txn"], ProfileID: vars["PID"], Hash: vars["Hash"]}
	result, Txn := display.InterpretPOE(display)
	fmt.Println(result)
	fmt.Println("Success")
	fmt.Println(Txn)

	response := apiModel.PoeSuccess{Message: "Success", TxNHash: Txn}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
