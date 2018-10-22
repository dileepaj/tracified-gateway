package businessFacades

import (

	// "main/proofs/retriever/stellarRetriever"

	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gorilla/mux"

	"main/api/apiModel"
	"main/model"

	// "main/proofs/builder"
	"main/proofs/interpreter"
)

func CheckPOC(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var response model.POC
	var TObj apiModel.POCOBJ
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

		pocStructObj := apiModel.POCStruct{Txn: vars["Txn"], ProfileID: vars["PID"], DBTree: TObj.Chain}
		display := &interpreter.AbstractPOC{POCStruct: pocStructObj}
		response = display.InterpretPOC()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(response.RetrievePOC.Error.Code)
		// w.WriteHeader(http.StatusBadRequest)

		// result := apiModel.PoeSuccess{Message: "response.RetrievePOC.Error.Message", TxNHash: "response.RetrievePOC.Txn"}
		result := apiModel.PocSuccess{Message: response.RetrievePOC.Error.Message, Chain: TObj.Chain}
		json.NewEncoder(w).Encode(result)
		return
	}
	return
}

func CheckFullPOC(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var response model.POC
	var TObj apiModel.POCOBJ
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

		pocStructObj := apiModel.POCStruct{
			Txn:       vars["Txn"],
			ProfileID: vars["PID"],
			DBTree:    TObj.Chain}
		display := &interpreter.AbstractPOC{POCStruct: pocStructObj}
		response = display.InterpretFullPOC()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(response.RetrievePOC.Error.Code)
		// w.WriteHeader(http.StatusBadRequest)

		// result := apiModel.PoeSuccess{Message: "response.RetrievePOC.Error.Message", TxNHash: "response.RetrievePOC.Txn"}
		result := apiModel.PocSuccess{
			Message: response.RetrievePOC.Error.Message,
			Chain:   TObj.Chain}
		json.NewEncoder(w).Encode(result)

		return
	}
	return
}

func CheckPOG(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var response model.POG
	pogStructObj := apiModel.POGStruct{LastTxn: vars["LastTxn"], POGTxn: vars["POGTxn"], Identifier: vars["Identifier"]}
	display := &interpreter.AbstractPOG{POGStruct: pogStructObj}
	response = display.InterpretPOG()

	// fmt.Println("response.RetrievePOG.Error.Code")
	// fmt.Println(response.RetrievePOG.Error.Code)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(response.RetrievePOG.Error.Code)
	// result := apiModel.PoeSuccess{Message: response.RetrievePOG.Error.Message, TxNHash: response.RetrievePOG.CurTxn}
	json.NewEncoder(w).Encode(response)
	return

}

func CheckPOE(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var response model.POE
	poeStructObj := apiModel.POEStruct{Txn: vars["Txn"], ProfileID: vars["PID"], Hash: vars["Hash"]}
	display := &interpreter.AbstractPOE{POEStruct: poeStructObj}
	response = display.InterpretPOE()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(response.RetrievePOE.Error.Code)
	json.NewEncoder(w).Encode(response.RetrievePOE)
	return

}
