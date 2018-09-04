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
	"main/proofs/retriever/stellarRetriever"
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
			DataHash: keys[i].DataHash,
			MergedID: keys[i].MergedID}
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
			result2 := apiModel.GenesisSuccess{Message: result.Error.Message, ProfileTxn: result.ProfileTxn, GenesisTxn: result.GenesisTxn, Identifiers: result.Identifiers, Type: result.TxnType}
			json.NewEncoder(w).Encode(result2)

		case "1":
			response := model.InsertProfileResponse{}

			display := &builder.AbstractProfileInsert{Identifiers: TObj.Identifiers[0], InsertType: TType, PreviousTXNID: TObj.PreviousTXNID[0], PreviousProfileID: TObj.PreviousProfileID[0]}
			response = display.ProfileInsert()

			w.WriteHeader(response.Error.Code)
			result := apiModel.ProfileSuccess{Message: response.Error.Message, ProfileTxn: response.ProfileTxn, PreviousTXNID: response.PreviousTXNID, PreviousProfileID: response.PreviousProfileID, Identifiers: response.Identifiers, Type: response.TxnType}
			json.NewEncoder(w).Encode(result)
		case "2":
			response := model.InsertDataResponse{}

			display := &builder.AbstractTDPInsert{Hash: TObj.Data, InsertType: TType, PreviousTXNID: TObj.PreviousTXNID[0], ProfileId: TObj.ProfileID[0]}
			response = display.TDPInsert()

			w.WriteHeader(response.Error.Code)
			result := apiModel.InsertSuccess{Message: response.Error.Message, TxNHash: response.TDPID, ProfileID: response.ProfileID, Type: response.TxnType}
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
				PreviousProfileID: TObj.ProfileID[0],
				Assets:            TObj.Assets,
				Code:              TObj.Code}
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

func CreateTrust(w http.ResponseWriter, r *http.Request) {

	// var response model.POE
	var TObj apiModel.CreateTrustLine
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
		display := &builder.AbstractTrustline{Code: TObj.Code, Limit: TObj.Limit, Issuerkey: TObj.Issuerkey, Signerkey: TObj.Signerkey}
		result := display.Trustline()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		result2 := apiModel.PoeSuccess{Message: "TrustLine Created", TxNHash: result}
		json.NewEncoder(w).Encode(result2)
		return
	}
	return
}
func SendAssests(w http.ResponseWriter, r *http.Request) {

	// var response model.POE
	var TObj apiModel.SendAssest
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
		display := &builder.AbstractAssetTransfer{SendAssest: TObj}
		response := display.AssetTransfer()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(response.Error.Code)
		result := apiModel.SendAssetRes{Message: response.Error.Message, PreviousTXNID: response.PreviousTXNID, PreviousProfileID: response.PreviousProfileID, Code: response.Code, Amount: response.Amount, Txn: response.Txn, To: response.To, From: response.From}
		json.NewEncoder(w).Encode(result)
		return
	}
	return
}

func MultisigAccount(w http.ResponseWriter, r *http.Request) {

	var TObj apiModel.RegistrarAccount
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

		// var response model.POE

		display := &builder.AbstractCreateRegistrar{RegistrarAccount: TObj}
		result := display.CreateRegistrarAcc()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		result2 := apiModel.PoeSuccess{Message: "Success", TxNHash: result}
		json.NewEncoder(w).Encode(result2)
		return
	}
	return

}

func AppointRegistrar(w http.ResponseWriter, r *http.Request) {

	var TObj apiModel.AppointRegistrar
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

		display := &builder.AbstractAppointRegistrar{AppointRegistrar: TObj}
		result := display.AppointReg()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		result2 := apiModel.RegSuccess{Message: "Success", Xdr: result}
		json.NewEncoder(w).Encode(result2)
		return
	}
	return
}
func TransformV2(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)

	// var response model.POE
	var TObj apiModel.AssestTransfer
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
		display := &builder.AbstractTransformAssets{AssestTransfer: TObj}
		result := display.TransformAssets()
		// display := &builder.AbstractTransformAssets{Code1: vars["code1"], Limit1: vars["limit1"], Code2: vars["code2"], Limit2: vars["limit2"], Code3: vars["code3"], Limit3: vars["limit3"], Code4: vars["code4"], Limit4: vars["limit4"]}
		// display := &builder.AbstractTransformAssets{Code1: TObj.Asset[0].Code, Limit1: TObj.Asset[0].Limit, Code2: TObj.Asset[1].Code, Limit2: TObj.Asset[1].Limit, Code3: TObj.Asset[2].Code, Limit3: TObj.Asset[2].Limit, Code4: TObj.Asset[3].Code, Limit4: TObj.Asset[3].Limit}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		result2 := apiModel.RegSuccess{Message: "Success", Xdr: result}
		json.NewEncoder(w).Encode(result2)
		return
	}
	return

}

func COC(w http.ResponseWriter, r *http.Request) {

	// var response model.POE
	var TObj apiModel.ChangeOfCustody
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
		display := &builder.AbstractCoCTransaction{ChangeOfCustody: TObj}
		response := display.CoCTransaction()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(response.Error.Code)
		result2 := apiModel.COCRes{Message: response.Error.Message, PreviousTXNID: response.PreviousTXNID, PreviousProfileID: response.PreviousProfileID, Code: response.Code, Amount: response.Amount, To: response.To, From: response.From, TxnXDR: response.TxnXDR}
		json.NewEncoder(w).Encode(result2)
		return
	}
	return
}

func COCLink(w http.ResponseWriter, r *http.Request) {

	// var response model.POE
	var TObj apiModel.ChangeOfCustodyLink
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
		display := &builder.AbstractcocLink{ChangeOfCustodyLink: TObj}
		result := display.CoCLink()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		result2 := apiModel.PoeSuccess{Message: "Success", TxNHash: result}
		json.NewEncoder(w).Encode(result2)
		return
	}
	return
}

func DeveloperRetriever(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var response model.POC

	display := &stellarRetriever.ConcretePOC{Txn: vars["Txn"]}
	response.RetrievePOC = display.RetrieveFullPOC()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(response.RetrievePOC.Error.Code)
	// w.WriteHeader(http.StatusBadRequest)

	// result := apiModel.PoeSuccess{Message: "response.RetrievePOC.Error.Message", TxNHash: "response.RetrievePOC.Txn"}
	result := apiModel.PocSuccess{
		Chain: response.RetrievePOC.BCHash}
	json.NewEncoder(w).Encode(result)

	return

}
