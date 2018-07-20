package businessFacades

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gorilla/mux"

	"main/api/apiModel"
	"main/model"
	"main/proofs/executer/stellarExecuter"
	"main/proofs/retriever/stellarRetriever"
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

//To be implemented
func CheckPOC(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

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
