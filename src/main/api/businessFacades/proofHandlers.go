package businessFacades

import (
	"io/ioutil"
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
	var dbTree []model.Current

	data, _ := ioutil.ReadAll(r.Body)

	var raw map[string]interface{}
	json.Unmarshal(data, &raw)
	// raw["count"] = 2
	out, _ := json.Marshal(raw["Chain"])

	keysBody := out
	keys := make([]model.Current, 0)
	json.Unmarshal(keysBody, &keys)
	// var lol

	for i := 0; i < len(keys); i++ {
		temp := model.Current{
			TType:    keys[i].TType,
			TXNID:    keys[i].TXNID,
			DataHash: keys[i].DataHash,
			MergedID: keys[i].MergedID}
		dbTree = append(dbTree, temp)
	}

	fmt.Println(dbTree)

	pocStructObj := apiModel.POCStruct{Txn: vars["Txn"], ProfileID: vars["PID"], DBTree: dbTree}
	display := &interpreter.AbstractPOC{POCStruct: pocStructObj}
	response = display.InterpretPOC()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(response.RetrievePOC.Error.Code)
	// w.WriteHeader(http.StatusBadRequest)

	// result := apiModel.PoeSuccess{Message: "response.RetrievePOC.Error.Message", TxNHash: "response.RetrievePOC.Txn"}
	result := apiModel.PocSuccess{Message: response.RetrievePOC.Error.Message, Chain: dbTree}
	json.NewEncoder(w).Encode(result)

	return
}

func CheckFullPOC(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var response model.POC
	var dbTree []model.Current

	data, _ := ioutil.ReadAll(r.Body)

	var raw map[string]interface{}
	json.Unmarshal(data, &raw)
	// raw["count"] = 2
	out, _ := json.Marshal(raw["Chain"])

	keysBody := out
	keys := make([]model.Current, 0)
	json.Unmarshal(keysBody, &keys)
	// var lol

	for i := 0; i < len(keys); i++ {
		temp := model.Current{
			TType:             keys[i].TType,
			TXNID:             keys[i].TXNID,
			DataHash:          keys[i].DataHash,
			ProfileID:         keys[i].ProfileID,
			PreviousProfileID: keys[i].PreviousProfileID,
			MergedID:          keys[i].MergedID,
			MergedChain:       keys[i].MergedChain,
			Identifier:        keys[i].Identifier,
			Assets:            keys[i].Assets,
			Time:              keys[i].Time}
		dbTree = append(dbTree, temp)
	}

	fmt.Println(dbTree)

	pocStructObj := apiModel.POCStruct{
		Txn: vars["Txn"], 
		ProfileID: vars["PID"], 
		DBTree: dbTree}
	display := &interpreter.AbstractPOC{POCStruct: pocStructObj}
	response = display.InterpretFullPOC()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(response.RetrievePOC.Error.Code)
	// w.WriteHeader(http.StatusBadRequest)

	// result := apiModel.PoeSuccess{Message: "response.RetrievePOC.Error.Message", TxNHash: "response.RetrievePOC.Txn"}
	result := apiModel.PocSuccess{
		Message: response.RetrievePOC.Error.Message,
		Chain:   dbTree}
	json.NewEncoder(w).Encode(result)

	return
}

func CheckPOG(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var response model.POG
	pogStructObj := apiModel.POGStruct{LastTxn: vars["LastTxn"], POGTxn: vars["POGTxn"], Identifier: vars["Identifier"]}
	display := &interpreter.AbstractPOG{POGStruct: pogStructObj}
	response = display.InterpretPOG()

	fmt.Println("response.RetrievePOG.Error.Code")
	fmt.Println(response.RetrievePOG.Error.Code)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(response.RetrievePOG.Error.Code)
	result := apiModel.PoeSuccess{Message: response.RetrievePOG.Error.Message, TxNHash: response.RetrievePOG.CurTxn}
	json.NewEncoder(w).Encode(result)
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
	result := apiModel.PoeSuccess{Message: response.RetrievePOE.Error.Message, TxNHash: response.RetrievePOE.Txn}
	json.NewEncoder(w).Encode(result)
	return

}
