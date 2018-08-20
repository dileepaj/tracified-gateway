package businessFacades

import (
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

	// 	display := &stellarExecuter.ConcreteInsertData{Hash: vars["hash"], InsertType: vars["type"], PreviousTDPID: vars["previousTDPID"], ProfileId: vars["profileId"]}
	// 	response = display.TDPInsert(display)

	// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// 	w.WriteHeader(response.Error.Code)
	// 	result := apiModel.InsertSuccess{Message: response.Error.Message, TxNHash: response.Txn, ProfileID: response.ProfileID, Type: response.TxnType}
	// 	json.NewEncoder(w).Encode(result)

	// return

}

//To be implemented
func CheckPOC(w http.ResponseWriter, r *http.Request) {
	// 	vars := mux.Vars(r)

	// 	var response model.POC

	// 	display := &stellarRetriever.ConcretePOC{Txn: vars["Txn"], ProfileID: vars["PID"], DBTree: vars["dbTree"]}
	// 	response = display.InterpretPOC(display)

	// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// 	w.WriteHeader(response.RetrievePOC.Error.Code)
	// 	result := apiModel.PoeSuccess{Message: response.RetrievePOC.Error.Message, TxNHash: response.RetrievePOC.Txn}
	// 	json.NewEncoder(w).Encode(result)
	// 	return

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

			display := &builder.AbstractProfileInsert{Identifiers: TObj.Identifiers[0], InsertType: TType, PreviousTDPID: TObj.PreviousTDPID[0], PreviousProfileID: TObj.PreviousProfileID[0]}
			response = display.ProfileInsert()

			w.WriteHeader(response.Error.Code)
			result := apiModel.ProfileSuccess{Message: response.Error.Message, TxNHash: response.Txn, PreviousTDPID: response.PreviousTDPID, PreviousProfileID: response.PreviousProfileID, Identifiers: response.Identifiers, Type: response.TxnType}
			json.NewEncoder(w).Encode(result)
		case "2":
			response := model.InsertDataResponse{}

			display := &builder.AbstractTDPInsert{Hash: TObj.Data[0], InsertType: TType, PreviousTDPID: TObj.PreviousTDPID[0], ProfileId: TObj.ProfileID[0]}
			response = display.TDPInsert()

			w.WriteHeader(response.Error.Code)
			result := apiModel.InsertSuccess{Message: response.Error.Message, TxNHash: response.Txn, ProfileID: response.ProfileID, Type: response.TxnType}
			json.NewEncoder(w).Encode(result)
		case "5":
			// var SplitProfiles []string
			// response := model.InsertProfileResponse{}

			// for i := 0; i < len(TObj.Identifiers); i++ {

			// 	display := &stellarExecuter.ConcreteProfile{Identifiers: TObj.Identifiers[i], InsertType: TType, PreviousTDPID: TObj.PreviousTDPID[0], PreviousProfileID: TObj.ProfileID[0]}
			// 	response = display.ProfileInsert(display)
			// 	SplitProfiles = append(SplitProfiles, response.Txn)
			// }

			// w.WriteHeader(response.Error.Code)
			// result := apiModel.SplitSuccess{Message: response.Error.Message, TxnHash: response.Txn, PreviousTDPID: TObj.PreviousTDPID[0], ProfileID: TObj.ProfileID[0], SplitProfiles: SplitProfiles, Identifiers: TObj.Identifiers, Type: TType}
			// json.NewEncoder(w).Encode(result)
		case "6":
			// var MergeProfiles []string
			// response := model.InsertProfileResponse{}

			// for i := 0; i < len(TObj.Identifiers); i++ {

			// 	display := &stellarExecuter.ConcreteProfile{Identifiers: TObj.Identifiers[i], InsertType: TType, PreviousTDPID: TObj.PreviousTDPID[0], PreviousProfileID: TObj.ProfileID[0]}
			// 	response = display.ProfileInsert(display)
			// 	MergeProfiles = append(MergeProfiles, response.Txn)
			// }

			// display2 := &stellarExecuter.ConcreteMerge{MergeProfiles: MergeProfiles[0], PreviousTDPID: TObj.PreviousTDPID[0]}
			// response2 := display2.ProfileMerge(display2)

			// w.WriteHeader(response2.Error.Code)
			// result := apiModel.SplitSuccess{Message: response2.Error.Message, TxnHash: response2.Txn, PreviousTDPID: TObj.PreviousTDPID[0], ProfileID: TObj.ProfileID[0], Identifiers: TObj.Identifiers, Type: TType}
			// json.NewEncoder(w).Encode(result)
		default:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Please send a valid Transaction Type")
			return
		}

	}

	return

}
