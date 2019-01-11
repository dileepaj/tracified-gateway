package businessFacades

import (
	// "io/ioutil"
	// "github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/proofs/retriever/stellarRetriever"

	// "gopkg.in/mgo.v2/bson"
	// "github.com/fanliao/go-promise"

	// "gopkg.in/mgo.v2"
	// "github.com/stellar/go/build"
	// "github.com/stellar/go/xdr"

	"encoding/json"
	"fmt"
	// "gopkg.in/mgo.v2"

	"net/http"

	// "github.com/fanliao/go-promise"
	"github.com/gorilla/mux"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	// "github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/builder"
	// "github.com/dileepaj/tracified-gateway/proofs/interpreter"
)

func CreateTrust(w http.ResponseWriter, r *http.Request) {

	// var response model.POE
	var TObj apiModel.TrustlineStruct
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
		display := &builder.AbstractTrustline{TrustlineStruct: TObj}
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
	var TObj apiModel.AssetTransfer
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
		display := &builder.AbstractTransformAssets{AssetTransfer: TObj}
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

	pocStructObj := apiModel.POCStruct{Txn: vars["Txn"]}
	display := &stellarRetriever.ConcretePOC{POCStruct: pocStructObj}
	// display := &stellarRetriever.ConcretePOC{Txn: vars["Txn"]}
	response.RetrievePOC = display.RetrieveFullPOC()
	// response.RetrievePOC = display.RetrievePOC()


	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	// w.WriteHeader(http.StatusBadRequest)

	// result := apiModel.PoeSuccess{Message: "response.RetrievePOC.Error.Message", TxNHash: "response.RetrievePOC.Txn"}
	result := apiModel.PocSuccess{
		Chain: response.RetrievePOC.BCHash}
	json.NewEncoder(w).Encode(result)

	return

}
