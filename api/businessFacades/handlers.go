package businessFacades

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	// "fmt"

	"net/http"

	"github.com/gorilla/mux"

	"github.com/tracified-gateway/api/apiModel"
	"github.com/tracified-gateway/model"
	"github.com/tracified-gateway/proofs/executer/stellarExecuter"
	"github.com/tracified-gateway/proofs/retriever/stellarRetriever"
)

func SaveData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response := model.InsertDataResponse{}

	display := &stellarExecuter.ConcreteInsertData{Hash: vars["hash"], InsertType: vars["type"], PreviousTDPID: vars["previousTDPID"], ProfileId: vars["profileId"]}
	response = display.TDPInsert(display)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(response.Error.Code)
	result := apiModel.InsertSuccess{Message: response.Error.Message, TxNHash: response.Txn, ProfileID: response.ProfileID, Type: response.TxnType}
	json.NewEncoder(w).Encode(result)

	return

}

// To be implemented
func CheckPOC(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var response model.POC
	lol := []model.Current{}

	TraceTree := vars["dbTree"]
	decoded, err := base64.StdEncoding.DecodeString(TraceTree)
	if err != nil {
		fmt.Println("decode error:", err)
	} else {
		var raw map[string]interface{}
		json.Unmarshal(decoded, &raw)
		// raw["count"] = 2
		out, _ := json.Marshal(raw["Chain"])

		keysBody := out
		keys := make([]model.Current, 0)
		json.Unmarshal(keysBody, &keys)
		// var lol

		for i := 0; i < len(keys); i++ {
			lo := model.Current{keys[i].TDPID, keys[i].Hash}
			lol = append(lol, lo)
		}
		// fmt.Println(lol)

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

func Transaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	TType := (vars["TType"])
	var TObj apiModel.TransactionStruct
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Please send a request body")
		return
	} else {
		err := json.NewDecoder(r.Body).Decode(&TObj)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Error while Decoding the body")
			fmt.Println(err)
			return
		}
		fmt.Println(TObj)

		switch TType {
		// case "0":
		// 	w.WriteHeader(http.StatusNotFound)
		// 	json.NewEncoder(w).Encode(apiModel.JsonErr{StatusCode: http.StatusNotFound, Error: "No root"})
		case "1":
			// response := model.InsertDataResponse{}

			// display := &stellarExecuter.ConcreteInsertData{Hash: TObj.Data[0], InsertType: TObj.TType, PreviousTDPID: TObj.PreviousTDPID, ProfileId: TObj.ProfileID[0]}
			// response = display.TDPInsert(display)

			// w.WriteHeader(response.Error.Code)
			// result := apiModel.InsertSuccess{Message: response.Error.Message, TxNHash: response.Txn, ProfileID: response.ProfileID, Type: response.TxnType}
			// json.NewEncoder(w).Encode(result)
		case "2":
			response := model.InsertDataResponse{}

			display := &stellarExecuter.ConcreteInsertData{Hash: TObj.Data[0], InsertType: TObj.TType, PreviousTDPID: TObj.PreviousTDPID, ProfileId: TObj.ProfileID[0]}
			response = display.TDPInsert(display)

			w.WriteHeader(response.Error.Code)
			result := apiModel.InsertSuccess{Message: response.Error.Message, TxNHash: response.Txn, ProfileID: response.ProfileID, Type: response.TxnType}
			json.NewEncoder(w).Encode(result)
		// case "5":
		// 	w.WriteHeader(http.StatusNotFound)
		// 	json.NewEncoder(w).Encode(apiModel.JsonErr{StatusCode: http.StatusNotFound, Error: "No root"})
		// case 6:
		// 	w.WriteHeader(http.StatusNotFound)
		// 	json.NewEncoder(w).Encode(apiModel.JsonErr{StatusCode: http.StatusNotFound, Error: "No root"})

		default:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Please send a valid Transaction Type")
			return
		}

	}

	return

}
