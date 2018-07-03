package businessFacades

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

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
			json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "No root"})
		case 1:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"})
		default:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"})
		}

	}

}

func CheckPOC(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	isValid, err := strconv.ParseBool(vars["isValid"])
	if err != nil {
		fmt.Println("Error in the Boolean!")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Error in the Boolean!"})
		return
	}
	result := interpreter.InterpretPOC(vars["rootHash"], isValid)

	//log the results
	fmt.Println(result, "result!!!\n")

	if result.Previous != "" && result.Current != "" {
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
	// 		json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "No PRE"})
	// 	case 1:
	// 		w.WriteHeader(http.StatusNotFound)
	// 		json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"})
	// 	default:
	// 		w.WriteHeader(http.StatusNotFound)
	// 		json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"})
	// 	}

	// }

}
