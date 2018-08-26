package businessFacades

import (
	"io/ioutil"

	// "encoding/base64"

	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gorilla/mux"

	"main/api/apiModel"
	"main/model"
	"main/proofs/builder"
	"main/proofs/interpreter"
)

func SaveData(w http.ResponseWriter, r *http.Request) {
	// 	vars := mux.Vars(r)
	// 	response := model.InsertDataResponse{}

	// 	display := &stellarExecuter.ConcreteInsertData{Hash: vars["hash"], InsertType: vars["type"], PreviousTXNID: vars["PreviousTXNID"], ProfileId: vars["profileId"]}
	// 	response = display.TDPInsert(display)

	// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// 	w.WriteHeader(response.Error.Code)
	// 	result := apiModel.InsertSuccess{Message: response.Error.Message, TxNHash: response.Txn, ProfileID: response.ProfileID, Type: response.TxnType}
	// 	json.NewEncoder(w).Encode(result)

	// return

}

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
			DataHash: keys[i].DataHash}
		dbTree = append(dbTree, temp)
	}

	fmt.Println(dbTree)

	// output := []model.Current{}
	display := &interpreter.AbstractPOC{Txn: vars["Txn"], ProfileID: vars["PID"], DBTree: dbTree}
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

	// output := []model.Current{}
	display := &interpreter.AbstractPOC{
		Txn:       vars["Txn"],
		ProfileID: vars["PID"],
		DBTree:    dbTree}
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

	display := &interpreter.AbstractPOG{LastTxn: vars["LastTxn"], POGTxn: vars["POGTxn"], Identifier: vars["Identifier"]}
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

	display := &interpreter.AbstractPOE{Txn: vars["Txn"], ProfileID: vars["PID"], Hash: vars["Hash"]}
	response = display.InterpretPOE()

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
		case "0":
			result := model.InsertGenesisResponse{}

			display := &builder.AbstractGenesisInsert{Identifiers: TObj.Identifiers[0], InsertType: TType}
			result = display.GenesisInsert()

			w.WriteHeader(result.Error.Code)
			result2 := apiModel.GenesisSuccess{Message: result.Error.Message, TxnHash: result.Txn, GenesisTxn: result.GenesisTxn, Identifiers: result.Identifiers, Type: result.TxnType}
			json.NewEncoder(w).Encode(result2)

		case "1":
			response := model.InsertProfileResponse{}

			display := &builder.AbstractProfileInsert{Identifiers: TObj.Identifiers[0], InsertType: TType, PreviousTXNID: TObj.PreviousTXNID[0], PreviousProfileID: TObj.PreviousProfileID[0]}
			response = display.ProfileInsert()

			w.WriteHeader(response.Error.Code)
			result := apiModel.ProfileSuccess{Message: response.Error.Message, TxNHash: response.Txn, PreviousTXNID: response.PreviousTXNID, PreviousProfileID: response.PreviousProfileID, Identifiers: response.Identifiers, Type: response.TxnType}
			json.NewEncoder(w).Encode(result)
		case "2":
			response := model.InsertDataResponse{}

			display := &builder.AbstractTDPInsert{Hash: TObj.Data[0], InsertType: TType, PreviousTXNID: TObj.PreviousTXNID[0], ProfileId: TObj.ProfileID[0]}
			response = display.TDPInsert()

			w.WriteHeader(response.Error.Code)
			result := apiModel.InsertSuccess{Message: response.Error.Message, TxNHash: response.Txn, ProfileID: response.ProfileID, Type: response.TxnType}
			json.NewEncoder(w).Encode(result)
		case "5":
			// var SplitProfiles []string
			response := model.SplitProfileResponse{}

			// for i := 0; i < len(TObj.Identifiers); i++ {

			display := &builder.AbstractSplitProfile{
				Identifiers:       TObj.Identifier,
				SplitIdentifiers:  TObj.Identifiers,
				InsertType:        TType,
				PreviousTXNID:     TObj.PreviousTXNID[0],
				PreviousProfileID: TObj.ProfileID[0]}
			response = display.ProfileSplit()
			// 	SplitProfiles = append(SplitProfiles, response.Txn)
			// }

			w.WriteHeader(response.Error.Code)
			result := apiModel.SplitSuccess{
				Message:          response.Error.Message,
				TxnHash:          response.Txn,
				PreviousTXNID:    response.PreviousTXNID,
				SplitProfiles:    response.SplitProfiles,
				SplitTXN:         response.SplitTXN,
				Identifier:       TObj.Identifier,
				SplitIdentifiers: TObj.Identifiers,
				Type:             TType}
			json.NewEncoder(w).Encode(result)
		case "6":

			response := model.MergeProfileResponse{}

			display := &builder.AbstractMergeProfile{
				Identifiers:        TObj.Identifier,
				InsertType:         TType,
				PreviousTXNID:      TObj.PreviousTXNID[0],
				PreviousProfileID:  TObj.ProfileID[0],
				MergingTXNs:        TObj.MergingTXNs,
				ProfileID:          TObj.ProfileID[0],
				MergingIdentifiers: TObj.Identifiers}
			response = display.ProfileMerge()

			w.WriteHeader(response.Error.Code)
			result := apiModel.MergeSuccess{
				Message:            response.Error.Message,
				TxnHash:            response.Txn,
				PreviousTXNID:      response.PreviousTXNID,
				ProfileID:          response.ProfileID,
				Identifier:         TObj.Identifier,
				Type:               TType,
				MergingIdentifiers: response.PreviousIdentifiers,
				MergeTXNs:          response.MergeTXNs}
			json.NewEncoder(w).Encode(result)
		default:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Please send a valid Transaction Type")
			return
		}

	}

	return

}
