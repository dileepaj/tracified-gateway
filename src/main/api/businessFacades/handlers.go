package businessFacades

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"

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

	output := make([]model.Current, 10)
	display := &stellarRetriever.ConcretePOC{Txn: vars["Txn"], ProfileID: vars["PID"], DBTree: output, BCTree: output}
	response = display.InterpretPOC(display)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(response.RetrievePOC.Error.Code)
	// result := apiModel.PoeSuccess{Message: response.RetrievePOC.Error.Message, TxNHash: response.RetrievePOC.Txn}
	result := apiModel.PoeSuccess{Message: response.RetrievePOC.Error.Message, TxNHash: response.RetrievePOC.Txn}
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

func Base64DecEnc(typ string, msg string) string {
	var text string

	if typ == "Encode" {
		encoded := base64.StdEncoding.EncodeToString([]byte(msg))
		text = (string(encoded))

	} else if typ == "Decode" {
		decoded, err := base64.StdEncoding.DecodeString(msg)
		if err != nil {
			fmt.Println("decode error:", err)
		} else {
			text = string(decoded)
		}

	} else {
		text = "Typ has to be either Encode or Decode!"
	}

	return text
}

func doStuff(lol *http.Request) {
	data, _ := ioutil.ReadAll(lol.Body)

	var raw map[string]interface{}
	json.Unmarshal(data, &raw)
	// raw["count"] = 2
	out, _ := json.Marshal(raw["_embedded"])

	var raw1 map[string]interface{}
	json.Unmarshal(out, &raw1)

	out1, _ := json.Marshal(raw1["records"])

	keysBody := out1
	keys := make([]PublicKey, 0)
	json.Unmarshal(keysBody, &keys)
}

type PublicKey struct {
	Name  string
	Value string
}

type KeysResponse struct {
	Collection []PublicKey
}
