package businessFacades

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"main/api/apiModel"
	"main/proofs/builder"
	"main/proofs/interpreter"
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

	result := interpreter.InterpretPOE(vars["hash"], vars["TDPId"], vars["rootHash"])

	//log the results
	fmt.Println(result, "result!!!")

	if result.TxNHash != "" && result.RootHash != "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return
	}
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
