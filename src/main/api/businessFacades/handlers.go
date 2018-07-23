package businessFacades

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	// "io/ioutil"

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
	vars := mux.Vars(r)

	var response model.POC
	lol:= []model.Current{}

	TraceTree:=vars["dbTree"]
	decoded, err := base64.StdEncoding.DecodeString(TraceTree)
	if err != nil {
		fmt.Println("decode error:", err)
	}else{
		var raw map[string]interface{}
		json.Unmarshal(decoded, &raw)
		// raw["count"] = 2
		out, _ := json.Marshal(raw["Chain"])
	
		keysBody := out
		keys := make([]model.Current, 0)
		json.Unmarshal(keysBody, &keys)
		// var lol
		
		for i:=0;i<len(keys);i++{
			lo:=model.Current{keys[i].TDPID,keys[i].Hash}
			lol= append(lol, lo)
		}
		fmt.Println(lol)

	}

	

	output := []model.Current{}
	display := &stellarRetriever.ConcretePOC{Txn: vars["Txn"], ProfileID: vars["PID"], DBTree: lol, BCTree: output}
	response = display.InterpretPOC(display)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(response.RetrievePOC.Error.Code)
	// result := apiModel.PoeSuccess{Message: response.RetrievePOC.Error.Message, TxNHash: response.RetrievePOC.Txn}
	result := apiModel.PocSuccess{Message: response.RetrievePOC.Error.Message, Chain: response.RetrievePOC.DBHash}
	json.NewEncoder(w).Encode(result)
	return

	// json.NewEncoder(w).Encode("result")
	// return

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


type PublicKey struct {
	Name  string
	Value string
}

type KeysResponse struct {
	Collection []PublicKey
}
