package businessFacades

import (
	"encoding/base64"
	"io/ioutil"
	"sort"

	"github.com/dileepaj/tracified-gateway/commons"
	log "github.com/sirupsen/logrus"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/stellar/go/xdr"

	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"strconv"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/gorilla/mux"
	"github.com/stellar/go/strkey"
)

type transaction struct {
	transaction string
}

type tdpToTransaction struct {
	TdpId string `json:"tdpId"`
}

func GetTransactionId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	object := dao.Connection{}
	p := object.GetTransactionForTdpId(vars["id"])
	fmt.Println(vars["id"])
	p.Then(func(data interface{}) interface{} {
		TxnHash := (data.(model.TransactionCollectionBody)).TxnHash

		mapD := map[string]string{"transaction": TxnHash}
		mapB, _ := json.Marshal(mapD)
		// fmt.Println(string(mapB))
		// trans := transaction{transaction:TxnHash}
		// s := fmt.Sprintf("%v", trans)

		encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
		text := encoded
		result := model.TransactionId{Txnhash: TxnHash,
			Url: commons.GetHorizonClient().URL + "/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
				text + "%3D%3D&network=public"}

		// res := TDP{TdpId: result.TdpId}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return nil
	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "TDP ID Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error
	})
	p.Await()

}

func GetTransactionsForTDP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	var result []model.TransactionIds
	object := dao.Connection{}
	p := object.GetAllTransactionForTdpId(vars["id"])
	fmt.Println(vars["id"])
	p.Then(func(data interface{}) interface{} {
		res := data.([]model.TransactionCollectionBody)
		for _, TxnBody := range res {
			TxnHash := TxnBody.TxnHash
			mapD := map[string]string{"transaction": TxnHash}
			mapB, err := json.Marshal(mapD)
			if err != nil {
				log.Error("Error while json.Marshal(mapD) " + err.Error())
			}
			// fmt.Println(string(mapB))
			// trans := transaction{transaction:TxnHash}
			// s := fmt.Sprintf("%v", trans)
			encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
			text := encoded
			temp := model.TransactionIds{Txnhash: TxnHash,
				Url: commons.GetHorizonClient().URL + "/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
					text + "%3D%3D&network=public",
				Identifier: TxnBody.Identifier}
			result = append(result, temp)
		}
		// res := TDP{TdpId: result.TdpId}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return nil
	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "TDP ID Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error
	})
	p.Await()
}

func GetTransactionsForTdps(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	fmt.Println("lol")
	var TDPs apiModel.GetTransactionId

	if r.Header == nil {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "No Header present!",
		}
		json.NewEncoder(w).Encode(result)

		return
	}

	if r.Header.Get("Content-Type") == "" {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "No Content-Type present!",
		}
		json.NewEncoder(w).Encode(result)

		return
	}

	// fmt.Println(TDP)
	err := json.NewDecoder(r.Body).Decode(&TDPs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error while Decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		fmt.Println(err)
		return
	}
	// fmt.Println(TDPs)
	object := dao.Connection{}

	var resultArray []model.TransactionIds
	var identifer string
	var arrSize = len(TDPs.TdpID)
	for i := 0; i < arrSize; i++ {

		p := object.GetTransactionForTdpId(TDPs.TdpID[i])
		fmt.Println(TDPs.TdpID[i])
		p.Then(func(data interface{}) interface{} {

			if TDPs.TdpID[i] == "" {
				temp := model.TransactionIds{Txnhash: "Not Found", Url: "Not Found", Identifier: "Not Found", TdpId: TDPs.TdpID[i]}
				resultArray = append(resultArray, temp)
			} else {
				Txn := data.(model.TransactionCollectionBody)
				// mapD := map[string]string{"transaction": Txn.TxnHash}
				// mapB, _ := json.Marshal(mapD)
				// fmt.Println(Txn.ProfileID)
				// trans := transaction{transaction:TxnHash}
				// s := fmt.Sprintf("%v", trans)
				// identifer = Txn.Identifier
				// encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
				// text := encoded
				temp := model.TransactionIds{Txnhash: Txn.TxnHash,
					Url: commons.GetHorizonClient().URL + "/transactions/" + Txn.TxnHash, Identifier: Txn.Identifier, TdpId: TDPs.TdpID[i]}

				resultArray = append(resultArray, temp)
			}

			// res := TDP{TdpId: result.TdpId}
			return nil
		}).Catch(func(error error) error {
			temp := model.TransactionIds{Txnhash: "Not Found", Url: "Not Found", Identifier: "Not Found", TdpId: TDPs.TdpID[i]}
			resultArray = append(resultArray, temp)

			return error
		})
		p.Await()

		q := object.GetPogTransaction(identifer)
		fmt.Println(identifer)
		q.Then(func(data interface{}) interface{} {

			// if TDPs.TdpID[i] == "" {
			// 	temp := model.TransactionIds{Txnhash: "Not Found", Url: "Not Found", Identifier: "Not Found", TdpId: TDPs.TdpID[i]}
			// 	resultArray = append(resultArray, temp)
			// } else {
			Txn := data.(model.TransactionCollectionBody)
			// mapD := map[string]string{"transaction": Txn.TxnHash}
			// mapB, _ := json.Marshal(mapD)
			// fmt.Println(Txn.TxnHash)
			// trans := transaction{transaction:TxnHash}
			// s := fmt.Sprintf("%v", trans)

			// encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
			// text := encoded
			temp := model.TransactionIds{Txnhash: Txn.TxnHash,
				Url: commons.GetHorizonClient().URL + "/transactions/" + Txn.TxnHash, Identifier: Txn.Identifier, TdpId: TDPs.TdpID[i]}
			resultArray = append(resultArray, temp)
			return nil
		}).Catch(func(error error) error {
			temp := model.TransactionIds{Txnhash: "Not Found", Url: "Not Found", Identifier: "Not Found", TdpId: TDPs.TdpID[i]}
			resultArray = append(resultArray, temp)
			return error
		})
		q.Await()
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultArray)

}

func GetTransactionsForPK(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	var result []model.TransactionIds
	object := dao.Connection{}
	p := object.GetAllTransactionForPK(vars["id"])
	// fmt.Println(vars["id"])
	p.Then(func(data interface{}) interface{} {
		res := data.([]model.TransactionCollectionBody)
		for _, TxnBody := range res {
			TxnHash := TxnBody.TxnHash
			mapD := map[string]string{"transaction": TxnHash}
			mapB, _ := json.Marshal(mapD)
			// fmt.Println(string(mapB))
			// trans := transaction{transaction:TxnHash}
			// s := fmt.Sprintf("%v", trans)
			encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
			text := encoded
			temp := model.TransactionIds{Txnhash: TxnHash,
				Url: commons.GetHorizonClient().URL + "/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
					text + "%3D%3D&network=public",
				Identifier: TxnBody.Identifier, TdpId: TxnBody.TdpId}
			result = append(result, temp)
		}
		// res := TDP{TdpId: result.TdpId}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return nil
	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "TDP ID Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error
	})
	p.Await()
}

func QueryTransactionsByKey(w http.ResponseWriter, r *http.Request) {
	//log.Debug("----------------------------- QueryTransactionsByKey --------------------------------")
	var response model.Error;
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var result []model.PrevTxnResponse
	key1, error := r.URL.Query()["perPage"]
	if !error || len(key1[0]) < 1 {
		log.Println("Url Param 'perPage' is missing")
		return
	}
	key2, error := r.URL.Query()["page"]
	if !error || len(key2[0]) < 1 {
		log.Println("Url Param 'page' is missing")
		return
	}
	key3, error := r.URL.Query()["txn"]
	if !error || len(key2[0]) < 1 {
		log.Println("Url Param 'txn' is missing")
		return
	}
	perPage, err := strconv.Atoi(key1[0])
	if err != nil {
		log.Error("Error while read limit " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		response = model.Error{Message: "The parameter should be an integer " + err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	page, err := strconv.Atoi(key2[0])
	if err != nil {
		log.Error("Error while read limit " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		response = model.Error{Message: "The parameter should be an integer " + err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	txn := key3[0]
	object := dao.Connection{}
	log.Info("SearchBy:  ",checkValidVersionByte(txn))
	switch  checkValidVersionByte(txn) {
	case "pk":
		qdata, err := object.GetAllTransactionForPK_Paginated(txn, page, perPage).Then(func(data interface{}) interface{} {
			return data
		}).Await()
		if err != nil {
			log.Error("Unable to connect gateway datastore")
			w.WriteHeader(http.StatusNotFound)
			response = model.Error{Code:http.StatusNotFound,Message: "Unable to connect gateway datastore"}
			json.NewEncoder(w).Encode(response)
			return
		}
		if qdata == nil {
			log.Error("Identifier is not found in gateway datastore")
			w.WriteHeader(http.StatusNoContent)
			response = model.Error{Code:http.StatusNoContent,Message: "Identifier is not found in gateway datastore"}
			json.NewEncoder(w).Encode(response)
			return
		}
		res := qdata.(model.TransactionCollectionBodyWithCount)
		count := strconv.Itoa(int(res.Count))
		for _, TxnBody := range res.Transactions {
			TxnHash := TxnBody.TxnHash
			var txe xdr.Transaction
			status := "success"
			timestamp := ""
			ledger := ""
			feePaid := ""
			from := ""
			to := ""
			assetcode := ""
			result1, err := http.Get(commons.GetHorizonClient().URL + "/transactions/" + TxnHash)
			if err != nil {
				log.Error("Unable to reach Stellar network in result1")
				status = "Unable to reach Stellar network"
				w.WriteHeader(http.StatusNotFound)		
				response=model.Error{Code:http.StatusNotFound,Message: "Unable to reach Stellar network"}
				json.NewEncoder(w).Encode(response)
			}
			if result1.StatusCode != 200 {
				log.Error("Transaction could not be retrieved from Stellar Network in result1")
				status = "Transaction could not be retrieved from Stellar Network"
				w.WriteHeader(http.StatusNoContent)	
				response=model.Error{Code:http.StatusNoContent,Message: "Transaction could not be retrieved from Stellar Network"}
				json.NewEncoder(w).Encode(response)
			}
			data, err := ioutil.ReadAll(result1.Body)
			if err != nil {
				log.Error("Unable to reading response")
				response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to reading response"}
			}
			if status == "success" {
				var raw map[string]interface{}
				json.Unmarshal(data, &raw)
				timestamp = fmt.Sprintf("%s", raw["created_at"])
				ledger = fmt.Sprintf("%.0f", raw["ledger"])
				feePaid = fmt.Sprintf("%s", raw["fee_charged"])
				from = fmt.Sprintf("%s", raw["source_account"])
				to = fmt.Sprintf("%s", raw["source_account"])
				errXDR := xdr.SafeUnmarshalBase64(fmt.Sprintf("%s", raw["envelope_xdr"]), &txe)
				if errXDR != nil {
					//log.Error("Error while SafeUnmarshalBase64 " + errXDR.Error())
				}
				if TxnBody.TxnType == "10" {
					result1, err := http.Get(commons.GetHorizonClient().URL + "/transactions/" + TxnHash + "/operations")
					if err != nil {
						log.Error("Error while getting transactions by txnhash ")
						w.WriteHeader(http.StatusNotFound)
						response := model.Error{Code:http.StatusNotFound,Message: "Txn for the TXN does not exist in the Blockchain "}
						json.NewEncoder(w).Encode(response)
					}
					if result1.StatusCode != 200 {
						log.Error("Transaction could not be retrieved from Stellar Network in acceptresult1")
						status = "Transaction could not be retrieved from Stellar Network"
						w.WriteHeader(http.StatusNoContent)	
						response=model.Error{Code:http.StatusNoContent,Message: "Transaction could not be retrieved from Stellar Network"}
						 json.NewEncoder(w).Encode(response)
					}
					data, err := ioutil.ReadAll(result1.Body)
					if err != nil {
						log.Error("Unable to reading response")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to reading response"}
					}
					var raw map[string]interface{}
					err = json.Unmarshal(data, &raw)
					if err != nil {
						log.Error("Unable to unmarshal row data")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal data"}
					}
					out, err := json.Marshal(raw["_embedded"])
					if err != nil {
						log.Error("Unable to marshal embedded")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to marshal embedded"}
					}
					var raw1 map[string]interface{}
					err = json.Unmarshal(out, &raw1)
					if err != nil {
						log.Error("Unable to unmarshal  json.Unmarshal(out, &raw1)")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal data"}	
					}
					out1, err := json.Marshal(raw1["records"])
					if err != nil {
						log.Error("Unable to marshal records ")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to marshal records"}
					}
					keysBody := out1
					keys := make([]PublicKeyPOCOC, 0)
					err = json.Unmarshal(keysBody, &keys)
					if err != nil {
						log.Error("Unable to unmarshal keys data")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal key's data"}
					}
					acceptTxn_byteData, err := base64.StdEncoding.DecodeString(keys[2].Value)
					if err != nil {
						log.Error("Unable to base64 decoding")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to base64 decoding"}
					}
					acceptTxn := string(acceptTxn_byteData)
					//log.Info("acceptTxn: " + acceptTxn)
					acceptresult1, err := http.Get(commons.GetHorizonClient().URL + "/transactions/" + acceptTxn + "/operations")
					if err != nil {
						log.Error("Unable to reach Stellar network in acceptresult1")
						status = "Unable to reach Stellar network"
						w.WriteHeader(http.StatusNotFound)		
						response=model.Error{Code:http.StatusNotFound,Message: "Unable to reach Stellar network"}
						 json.NewEncoder(w).Encode(response)			
					}
					if acceptresult1.StatusCode != 200 {
						log.Error("Transaction could not be retrieved from Stellar Network in acceptresult1")
						status = "Transaction could not be retrieved from Stellar Network"
						w.WriteHeader(http.StatusNoContent)	
						response=model.Error{Code:http.StatusNoContent,Message: "Transaction could not be retrieved from Stellar Network"}
						 json.NewEncoder(w).Encode(response)
					}
					acceptdata, err := ioutil.ReadAll(acceptresult1.Body)
					if err != nil {
						log.Error("Unable to reading response",err)
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to reading response"}
					}
					var acceptraw map[string]interface{}
					err = json.Unmarshal(acceptdata, &acceptraw)
					if err != nil {
						log.Error("Unable to unmarshal data")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal data"}
					}
					acceptout, err := json.Marshal(acceptraw["_embedded"])
					if err != nil {
						log.Error("Unable to marshal _embedded ")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to marshal _embedded"}
					}
					var acceptraw1 map[string]interface{}
					err = json.Unmarshal(acceptout, &acceptraw1)
					if err != nil {
						log.Error("Unable to unmarshal data")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal data"}
					}
					acceptout1, err := json.Marshal(acceptraw1["records"])
					if err != nil {
						log.Error("Unable to marshal records")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to marshal records"}
					}
					acceptkeysBody := acceptout1
					acceptkeys := make([]PublicKeyPOCOC, 0)
					err = json.Unmarshal(acceptkeysBody, &acceptkeys)
					if err != nil {
						log.Error("Unable to unmarshal key's data")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal key's data"}
					}

					to = string(acceptkeys[3].To)
					log.Info("Destination: " + to)

					from = string(acceptkeys[3].From)
					log.Info("Source: " + from)

					assetcode = string(acceptkeys[3].Asset_code)
					log.Info("Asset Code: " + assetcode)
				}
			}
			mapD := map[string]string{"transaction": TxnHash}
			mapB, _ := json.Marshal(mapD)
			// fmt.Println(string(mapB))
			// trans := transaction{transaction:TxnHash}
			// s := fmt.Sprintf("%v", trans)
			encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
			text := encoded
			temp := model.PrevTxnResponse{
				Status: status, Txnhash: TxnHash,
				Url: commons.GetHorizonClient().URL + "/transactions/" + TxnHash + "/operations",
				LabUrl: commons.GetStellarLaboratoryClient() + "/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
					text + "%3D%3D&network=" + commons.GetHorizonClientNetworkName(),
				Identifier:     TxnBody.Identifier,
				TdpId:          TxnBody.TdpId,
				DataHash:       TxnBody.DataHash,
				Blockchain:		"Stellar",
				Timestamp:      timestamp,
				TxnType:        GetTransactiontype(TxnBody.TxnType),
				FeePaid:        feePaid,
				Ledger:         ledger,
				SourceAccount:  TxnBody.PublicKey,
				From:           from,
				SequenceNo:     TxnBody.SequenceNo,
				AvailableProof: GetProofName(TxnBody.TxnType),
				To:             to,
				ProductName:    TxnBody.ProductName,
				Itemcount:      count,
				AssetCode:      assetcode}
				if(response.Message!=""&&response.Code!=0){
					log.Error(response.Message)
					w.WriteHeader(response.Code)
					json.NewEncoder(w).Encode(response)
					}else{
						result = append(result, temp)
					}
			result = append(result, temp)
		}
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].Timestamp > result[j].Timestamp
		})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return

	case "txn":
		qdata, err := object.GetAllTransactionForTxId(txn).Then(func(data interface{}) interface{} {
			return data
		}).Await()
		if err != nil {
			log.Error("Unable to connect gateway datastore")
			w.WriteHeader(http.StatusNotFound)
			response = model.Error{Code:http.StatusNotFound,Message: "Unable to connect gateway datastaore"}
			json.NewEncoder(w).Encode(response)
			return
		}
		if qdata == nil {
			log.Error("Identifier is not found in gateway datastore")
			w.WriteHeader(http.StatusNoContent)
			response = model.Error{Code:http.StatusNoContent,Message: "Identifier is not found in gateway datastore"}
			json.NewEncoder(w).Encode(response)
			return
		}
		res := qdata.([]model.TransactionCollectionBody)
		for _, TxnBody := range res {
			TxnHash := TxnBody.TxnHash
			var txe xdr.Transaction
			status := "success"
			timestamp := ""
			ledger := ""
			feePaid := ""
			from := ""
			to := ""
			assetcode := ""
			result1, err := http.Get(commons.GetHorizonClient().URL + "/transactions/" + TxnHash)
			if err != nil {
				log.Error("Unable to reach Stellar network in result1")
				status = "Unable to reach Stellar network"
				w.WriteHeader(http.StatusNotFound)		
				response=model.Error{Code:http.StatusNotFound,Message: "Unable to reach Stellar network"}
				json.NewEncoder(w).Encode(response)
			}
			if result1.StatusCode != 200 {
				log.Error("Transaction could not be retrieved from Stellar Network in result1")
				status = "Transaction could not be retrieved from Stellar Network"
				w.WriteHeader(http.StatusNoContent)
				response=model.Error{Code:http.StatusNoContent,Message: "Transaction could not be retrieved from Stellar Network"}
				json.NewEncoder(w).Encode(response)
			}
			data, _ := ioutil.ReadAll(result1.Body)
			if status == "success" {
				var raw map[string]interface{}
				json.Unmarshal(data, &raw)
				timestamp = fmt.Sprintf("%s", raw["created_at"])
				ledger = fmt.Sprintf("%.0f", raw["ledger"])
				feePaid = fmt.Sprintf("%s", raw["fee_charged"])
				from = fmt.Sprintf("%s", raw["source_account"])
				to = fmt.Sprintf("%s", raw["source_account"])
				errXDR := xdr.SafeUnmarshalBase64(fmt.Sprintf("%s", raw["envelope_xdr"]), &txe)
				if errXDR != nil {
					//log.Error("Error SafeUnmarshalBase64 " + errXDR.Error())
				}
				if TxnBody.TxnType == "10" {
					result1, err := http.Get(commons.GetHorizonClient().URL + "/transactions/" + TxnHash + "/operations")
					if err != nil {
						log.Error("Unable to reach Stellar network in result1")
						status = "Unable to reach Stellar network"
						w.WriteHeader(http.StatusNotFound)		
						response=model.Error{Code:http.StatusNotFound,Message: "Unable to reach Stellar network"}
						json.NewEncoder(w).Encode(response)
					}
					if result1.StatusCode != 200 {
						log.Error("Transaction could not be retrieved from Stellar Network in result1")
						status = "Transaction could not be retrieved from Stellar Network"
						w.WriteHeader(http.StatusNoContent)	
						response=model.Error{Code:http.StatusNoContent,Message: "Transaction could not be retrieved from Stellar Network"}
						json.NewEncoder(w).Encode(response)
					}
					data, err := ioutil.ReadAll(result1.Body)
					if err != nil {
						log.Error("Unable to reading response")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to reading response"}
					}
					var raw map[string]interface{}
					err = json.Unmarshal(data, &raw)
					if err != nil {
						log.Error("Unable to unmarshal row data")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal data"}
					}
					out, err := json.Marshal(raw["_embedded"])
					if err != nil {
						log.Error("Unable to marshal embedded")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to marshal embedded"}
					}
					var raw1 map[string]interface{}
					err = json.Unmarshal(out, &raw1)
					if err != nil {
						log.Error("Unable to unmarshal  json.Unmarshal(out, &raw1)")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal data"}	
					}
					out1, err := json.Marshal(raw1["records"])
					if err != nil {
						log.Error("Unable to marshal records ")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to marshal records"}
					}
					keysBody := out1
					keys := make([]PublicKeyPOCOC, 0)
					err = json.Unmarshal(keysBody, &keys)
					if err != nil {
						log.Error("Unable to unmarshal keys data")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal key's data"}
					}
					acceptTxn_byteData, err := base64.StdEncoding.DecodeString(keys[2].Value)
					if err != nil {
						log.Error("Unable to base64 decoding")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to base64 decoding"}
					}
					acceptTxn := string(acceptTxn_byteData)
					//log.Info("acceptTxn: " + acceptTxn)
					acceptresult1, err := http.Get(commons.GetHorizonClient().URL + "/transactions/" + acceptTxn + "/operations")
					if err != nil {
						log.Error("Unable to reach Stellar network in acceptresult1")
						status = "Unable to reach Stellar network"
						w.WriteHeader(http.StatusNotFound)		
						response=model.Error{Code:http.StatusNotFound,Message: "Unable to reach Stellar network"}
						json.NewEncoder(w).Encode(response)			
					}
					if acceptresult1.StatusCode != 200 {
						log.Error("Transaction could not be retrieved from Stellar Network in acceptresult1")
						status = "Transaction could not be retrieved from Stellar Network"
						w.WriteHeader(http.StatusNoContent)		
						response=model.Error{Code:http.StatusNoContent,Message: "Transaction could not be retrieved from Stellar Network"}
						json.NewEncoder(w).Encode(response)
					}
					acceptdata, err := ioutil.ReadAll(acceptresult1.Body)
					if err != nil {
						log.Error("Unable to reading response",err)
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to reading response"}
					}
					var acceptraw map[string]interface{}
					err = json.Unmarshal(acceptdata, &acceptraw)
					if err != nil {
						log.Error("Unable to unmarshal data")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal data"}
					}
					acceptout, err := json.Marshal(acceptraw["_embedded"])
					if err != nil {
						log.Error("Unable to marshal _embedded ")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to marshal _embedded"}
					}
					var acceptraw1 map[string]interface{}
					err = json.Unmarshal(acceptout, &acceptraw1)
					if err != nil {
						log.Error("Unable to unmarshal data")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal data"}
					}
					acceptout1, err := json.Marshal(acceptraw1["records"])
					if err != nil {
						log.Error("Unable to marshal records")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to marshal records"}
					}
					acceptkeysBody := acceptout1
					acceptkeys := make([]PublicKeyPOCOC, 0)
					err = json.Unmarshal(acceptkeysBody, &acceptkeys)
					if err != nil {
						log.Error("Unable to unmarshal key's data")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal key's data"}
					}

					to = string(acceptkeys[3].To)
					log.Info("Destination: " + to)

					from = string(acceptkeys[3].From)
					log.Info("Source: " + from)

					assetcode = string(acceptkeys[3].Asset_code)
					log.Info("Asset Code: " + assetcode)
				}
			}
			mapD := map[string]string{"transaction": TxnHash}
			mapB, _ := json.Marshal(mapD)
			// fmt.Println(string(mapB))
			// trans := transaction{transaction:TxnHash}
			// s := fmt.Sprintf("%v", trans)

			encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
			text := encoded
			temp := model.PrevTxnResponse{
				Status: status, Txnhash: TxnHash,
				Url: commons.GetHorizonClient().URL + "/transactions/" + TxnHash + "/operations",
				LabUrl: commons.GetStellarLaboratoryClient() + "/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
					text + "%3D%3D&network=" + commons.GetHorizonClientNetworkName(),
				Identifier:     TxnBody.Identifier,
				TdpId:          TxnBody.TdpId,
				Blockchain:		"Stellar",
				DataHash:       TxnBody.DataHash,
				Timestamp:      timestamp,
				TxnType:        GetTransactiontype(TxnBody.TxnType),
				FeePaid:        feePaid,
				Ledger:         ledger,
				SourceAccount:  TxnBody.PublicKey,
				From:           from,
				SequenceNo:     TxnBody.SequenceNo,
				AvailableProof: GetProofName(TxnBody.TxnType),
				To:             to,
				ProductName:    TxnBody.ProductName,
				AssetCode:      assetcode}
				if(response.Message!=""&&response.Code!=0){
					log.Error(response.Message)
					w.WriteHeader(response.Code)
					json.NewEncoder(w).Encode(response)
					}else{
						result = append(result, temp)
					}	
		}
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return

	case "tdpid":
		qdata, err := object.GetAllTransactionForTdpId_Paginated(txn, page, perPage).Then(func(data interface{}) interface{} {
			return data
		}).Await()
		if err != nil {
			log.Error("Unable to connect gateway datastore")
			w.WriteHeader(http.StatusNotFound)
			response = model.Error{Code:http.StatusNotFound,Message: "Unable to connect gateway datastaore"}
			json.NewEncoder(w).Encode(response)
			return
		}
		if qdata == nil {
			log.Error("Identifier is not found in gateway datastore")
			w.WriteHeader(http.StatusNoContent)
			response = model.Error{Code:http.StatusNoContent,Message: "Identifier is not found in gateway datastore"}
			json.NewEncoder(w).Encode(response)
			return
		}
		res := qdata.(model.TransactionCollectionBodyWithCount)
		count := strconv.Itoa(int(res.Count))
		for _, TxnBody := range res.Transactions {
			TxnHash := TxnBody.TxnHash
			var txe xdr.Transaction
			status := "success"
			timestamp := ""
			ledger := ""
			feePaid := ""
			from := ""
			to := ""
			assetcode := ""
			result1, err := http.Get(commons.GetHorizonClient().URL + "/transactions/" + TxnHash)
			if err != nil {
				log.Error("Unable to reach Stellar network in result1")
				status = "Unable to reach Stellar network"
				w.WriteHeader(http.StatusNotFound)		
				response=model.Error{Code:http.StatusNotFound,Message: "Unable to reach Stellar network"}
				json.NewEncoder(w).Encode(response)
			}
			if result1.StatusCode != 200 {
				log.Error("Transaction could not be retrieved from Stellar Network in result1")
				status = "Transaction could not be retrieved from Stellar Network"
				w.WriteHeader(http.StatusNoContent)	
				response=model.Error{Code:http.StatusNoContent,Message: "Transaction could not be retrieved from Stellar Network"}
				json.NewEncoder(w).Encode(response)
			}
			data, _ := ioutil.ReadAll(result1.Body)
			if status == "success" {
				var raw map[string]interface{}
				json.Unmarshal(data, &raw)
				timestamp = fmt.Sprintf("%s", raw["created_at"])
				ledger = fmt.Sprintf("%.0f", raw["ledger"])
				feePaid = fmt.Sprintf("%s", raw["fee_charged"])
				// from = fmt.Sprintf("%s", raw["source_account"])
				// to = fmt.Sprintf("%s", raw["source_account"])
				errXDR := xdr.SafeUnmarshalBase64(fmt.Sprintf("%s", raw["envelope_xdr"]), &txe)
				if errXDR != nil {
					//log.Error("Error SafeUnmarshalBase64" + errXDR.Error())
				}
				if TxnBody.TxnType == "10" {
					result1, err := http.Get(commons.GetHorizonClient().URL + "/transactions/" + TxnHash + "/operations")
					if err != nil {
						log.Error("Unable to reach Stellar network in result1")
						status = "Unable to reach Stellar network"
						w.WriteHeader(http.StatusNotFound)		
						response=model.Error{Code:http.StatusNotFound,Message: "Unable to reach Stellar network"}
						json.NewEncoder(w).Encode(response)
					}
					if result1.StatusCode != 200 {
						log.Error("Transaction could not be retrieved from Stellar Network in result1")
						status = "Transaction could not be retrieved from Stellar Network"
						w.WriteHeader(http.StatusNoContent)	
						response=model.Error{Code:http.StatusNoContent,Message: "Transaction could not be retrieved from Stellar Network"}
						json.NewEncoder(w).Encode(response)
					}
					data, err := ioutil.ReadAll(result1.Body)
					if err != nil {
						log.Error("Unable to reading response")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to reading response"}
					}
					var raw map[string]interface{}
					err = json.Unmarshal(data, &raw)
					if err != nil {
						log.Error("Unable to unmarshal row data")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal data"}
					}
					out, err := json.Marshal(raw["_embedded"])
					if err != nil {
						log.Error("Unable to marshal embedded")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to marshal embedded"}
					}
					var raw1 map[string]interface{}
					err = json.Unmarshal(out, &raw1)
					if err != nil {
						log.Error("Unable to unmarshal  json.Unmarshal(out, &raw1)")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal data"}	
					}
					out1, err := json.Marshal(raw1["records"])
					if err != nil {
						log.Error("Unable to marshal records ")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to marshal records"}
					}
					keysBody := out1
					keys := make([]PublicKeyPOCOC, 0)
					err = json.Unmarshal(keysBody, &keys)
					if err != nil {
						log.Error("Unable to unmarshal keys data")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal key's data"}
					}
					acceptTxn_byteData, err := base64.StdEncoding.DecodeString(keys[2].Value)
					if err != nil {
						log.Error("Unable to base64 decoding")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to base64 decoding"}
					}
					acceptTxn := string(acceptTxn_byteData)
					//log.Info("acceptTxn: " + acceptTxn)
					acceptresult1, err := http.Get(commons.GetHorizonClient().URL + "/transactions/" + acceptTxn + "/operations")
					if err != nil {
						log.Error("Unable to reach Stellar network in acceptresult1")
						status = "Unable to reach Stellar network"
						w.WriteHeader(http.StatusNotFound)		
						response=model.Error{Code:http.StatusNotFound,Message: "Unable to reach Stellar network"}
						 json.NewEncoder(w).Encode(response)			
					}
					if acceptresult1.StatusCode != 200 {
						log.Error("Transaction could not be retrieved from Stellar Network in acceptresult1")
						status = "Transaction could not be retrieved from Stellar Network"
						w.WriteHeader(http.StatusNoContent)	
						response=model.Error{Code:http.StatusNoContent,Message: "Transaction could not be retrieved from Stellar Network"}
						 json.NewEncoder(w).Encode(response)
					}
					acceptdata, err := ioutil.ReadAll(acceptresult1.Body)
					if err != nil {
						log.Error("Unable to reading response",err)
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to reading response"}
					}
					var acceptraw map[string]interface{}
					err = json.Unmarshal(acceptdata, &acceptraw)
					if err != nil {
						log.Error("Unable to unmarshal data")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal data"}
					}
					acceptout, err := json.Marshal(acceptraw["_embedded"])
					if err != nil {
						log.Error("Unable to marshal _embedded ")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to marshal _embedded"}
					}
					var acceptraw1 map[string]interface{}
					err = json.Unmarshal(acceptout, &acceptraw1)
					if err != nil {
						log.Error("Unable to unmarshal data")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal data"}
					}
					acceptout1, err := json.Marshal(acceptraw1["records"])
					if err != nil {
						log.Error("Unable to marshal records")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to marshal records"}
					}
					acceptkeysBody := acceptout1
					acceptkeys := make([]PublicKeyPOCOC, 0)
					err = json.Unmarshal(acceptkeysBody, &acceptkeys)
					if err != nil {
						log.Error("Unable to unmarshal key's data")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal key's data"}
					}

					to = string(acceptkeys[3].To)
					log.Info("Destination: " + to)

					from = string(acceptkeys[3].From)
					log.Info("Source: " + from)

					assetcode = string(acceptkeys[3].Asset_code)
					log.Info("Asset Code: " + assetcode)
				}
			}
			mapD := map[string]string{"transaction": TxnHash}
			mapB, _ := json.Marshal(mapD)
			// fmt.Println(string(mapB))
			// trans := transaction{transaction:TxnHash}
			// s := fmt.Sprintf("%v", trans)
			encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
			text := encoded
			temp := model.PrevTxnResponse{
				Status: status, Txnhash: TxnHash,
				Url: commons.GetHorizonClient().URL + "/transactions/" + TxnHash + "/operations",
				LabUrl: commons.GetStellarLaboratoryClient() + "/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
					text + "%3D%3D&network=" + commons.GetHorizonClientNetworkName(),
				Identifier:     TxnBody.Identifier,
				Blockchain:		"Stellar",
				TdpId:          TxnBody.TdpId,
				DataHash:       TxnBody.DataHash,
				Timestamp:      timestamp,
				TxnType:        GetTransactiontype(TxnBody.TxnType),
				FeePaid:        feePaid,
				Ledger:         ledger,
				SourceAccount:  TxnBody.PublicKey,
				From:           from,
				SequenceNo:     TxnBody.SequenceNo,
				AvailableProof: GetProofName(TxnBody.TxnType),
				To:             to,
				ProductName:    TxnBody.ProductName,
				Itemcount:      count,
				AssetCode:      assetcode}
				if(response.Message!=""&&response.Code!=0){
					log.Error(response.Message)
					w.WriteHeader(response.Code)
					json.NewEncoder(w).Encode(response)
					}else{
						result = append(result, temp)
					}	
		}
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].Timestamp > result[j].Timestamp
		})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return

	case "":
		qdata, err := object.GetTransactionsbyIdentifier_Paginated(txn, page, perPage).Then(func(data interface{}) interface{} {
			return data
		}).Await()
		if err != nil {
			log.Error("Unable to connect gateway datastore")
			w.WriteHeader(http.StatusNotFound)
			response = model.Error{Code:http.StatusNotFound,Message: "Unable to connect gateway datastaore"}
			json.NewEncoder(w).Encode(response)
			return
		}
		if qdata == nil {
			log.Error("Identifier is not found in gateway datastore")
			w.WriteHeader(http.StatusNoContent)
			response = model.Error{Code:http.StatusNoContent,Message: "Identifier is not found in gateway datastore"}
			json.NewEncoder(w).Encode(response)
			return
		}
		res := qdata.(model.TransactionCollectionBodyWithCount)
		count := strconv.Itoa(int(res.Count))
		for _, TxnBody := range res.Transactions {
			TxnHash := TxnBody.TxnHash
			var txe xdr.Transaction
			status := "success"
			timestamp := ""
			ledger := ""
			feePaid := ""
			from := ""
			to := ""
			assetcode := ""
			result1, err := http.Get(commons.GetHorizonClient().URL + "/transactions/" + TxnHash)
			if err != nil {
				log.Error("Unable to reach Stellar network in result1")
				status = "Unable to reach Stellar network"
				w.WriteHeader(http.StatusBadRequest)		
				response=model.Error{Code:http.StatusBadRequest,Message: "Unable to reach Stellar network"}
				json.NewEncoder(w).Encode(response)
			}
			if result1.StatusCode != 200 {
				log.Error("Transaction could not be retrieved from Stellar Network in result1")
				status = "Transaction could not be retrieved from Stellar Network"
				w.WriteHeader(http.StatusNoContent)	
				response=model.Error{Code:http.StatusNoContent,Message: "Transaction could not be retrieved from Stellar Network"}
				json.NewEncoder(w).Encode(response)
			}
			data, _ := ioutil.ReadAll(result1.Body)
			if status == "success" {
				var raw map[string]interface{}
				json.Unmarshal(data, &raw)
				timestamp = fmt.Sprintf("%s", raw["created_at"])
				ledger = fmt.Sprintf("%.0f", raw["ledger"])
				feePaid = fmt.Sprintf("%s", raw["fee_charged"])
				from = fmt.Sprintf("%s", raw["source_account"])
				to = fmt.Sprintf("%s", raw["source_account"])
				errXDR := xdr.SafeUnmarshalBase64(fmt.Sprintf("%s", raw["envelope_xdr"]), &txe)
				if errXDR != nil {
					//log.Error("Error SafeUnmarshalBase64" + errXDR.Error())
				}
				if TxnBody.TxnType == "10" {
					result1, err := http.Get(commons.GetHorizonClient().URL + "/transactions/" + TxnHash + "/operations")
					if err != nil {
						log.Error("Unable to reach Stellar network in result1")
						status = "Unable to reach Stellar network"
						w.WriteHeader(http.StatusNotFound)		
						response=model.Error{Code:http.StatusNotFound,Message: "Unable to reach Stellar network"}
						json.NewEncoder(w).Encode(response)
					}
					if result1.StatusCode != 200 {
						log.Error("Transaction could not be retrieved from Stellar Network in result1")
						status = "Transaction could not be retrieved from Stellar Network"
						w.WriteHeader(http.StatusNoContent)	
						response=model.Error{Code:http.StatusNoContent,Message: "Transaction could not be retrieved from Stellar Network"}
						json.NewEncoder(w).Encode(response)
					}
					data, err := ioutil.ReadAll(result1.Body)
					if err != nil {
						log.Error("Unable to reading response")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to reading response"}
					}
					var raw map[string]interface{}
					err = json.Unmarshal(data, &raw)
					if err != nil {
						log.Error("Unable to unmarshal row data")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal data"}
					}
					out, err := json.Marshal(raw["_embedded"])
					if err != nil {
						log.Error("Unable to marshal embedded")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to marshal embedded"}
					}
					var raw1 map[string]interface{}
					err = json.Unmarshal(out, &raw1)
					if err != nil {
						log.Error("Unable to unmarshal  json.Unmarshal(out, &raw1)")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal data"}	
					}
					out1, err := json.Marshal(raw1["records"])
					if err != nil {
						log.Error("Unable to marshal records ")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to marshal records"}
					}
					keysBody := out1
					keys := make([]PublicKeyPOCOC, 0)
					err = json.Unmarshal(keysBody, &keys)
					if err != nil {
						log.Error("Unable to unmarshal keys data")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal key's data"}
					}
					acceptTxn_byteData, err := base64.StdEncoding.DecodeString(keys[2].Value)
					if err != nil {
						log.Error("Unable to base64 decoding")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to base64 decoding"}
					}
					acceptTxn := string(acceptTxn_byteData)
					//log.Info("acceptTxn: " + acceptTxn)
					acceptresult1, err := http.Get(commons.GetHorizonClient().URL + "/transactions/" + acceptTxn + "/operations")
					if err != nil {
						log.Error("Unable to reach Stellar network in acceptresult1")
						status = "Unable to reach Stellar network"
						w.WriteHeader(http.StatusNotFound)		
						response=model.Error{Code:http.StatusNotFound,Message: "Unable to reach Stellar network"}
						 json.NewEncoder(w).Encode(response)			
					}
					if acceptresult1.StatusCode != 200 {
						log.Error("Transaction could not be retrieved from Stellar Network in acceptresult1")
						status = "Transaction could not be retrieved from Stellar Network"
						w.WriteHeader(http.StatusNoContent)		
						response=model.Error{Code:http.StatusNoContent,Message: "Transaction could not be retrieved from Stellar Network"}
						 json.NewEncoder(w).Encode(response)
					}
					acceptdata, err := ioutil.ReadAll(acceptresult1.Body)
					if err != nil {
						log.Error("Unable to reading response",err)
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to reading response"}
					}
					var acceptraw map[string]interface{}
					err = json.Unmarshal(acceptdata, &acceptraw)
					if err != nil {
						log.Error("Unable to unmarshal data")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal data"}
					}
					acceptout, err := json.Marshal(acceptraw["_embedded"])
					if err != nil {
						log.Error("Unable to marshal _embedded ")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to marshal _embedded"}
					}
					var acceptraw1 map[string]interface{}
					err = json.Unmarshal(acceptout, &acceptraw1)
					if err != nil {
						log.Error("Unable to unmarshal data")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal data"}
					}
					acceptout1, err := json.Marshal(acceptraw1["records"])
					if err != nil {
						log.Error("Unable to marshal records")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to marshal records"}
					}
					acceptkeysBody := acceptout1
					acceptkeys := make([]PublicKeyPOCOC, 0)
					err = json.Unmarshal(acceptkeysBody, &acceptkeys)
					if err != nil {
						log.Error("Unable to unmarshal key's data")
						response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal key's data"}
					}

					to = string(acceptkeys[3].To)
					log.Info("Destination: " + to)

					from = string(acceptkeys[3].From)
					log.Info("Source: " + from)

					assetcode = string(acceptkeys[3].Asset_code)
					log.Info("Asset Code: " + assetcode)
				}
			}
			mapD := map[string]string{"transaction": TxnHash}
			mapB, _ := json.Marshal(mapD)
			// fmt.Println(string(mapB))
			// trans := transaction{transaction:TxnHash}
			// s := fmt.Sprintf("%v", trans)
			encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
			text := encoded
			temp := model.PrevTxnResponse{
				Status: status, Txnhash: TxnHash,
				Url: commons.GetHorizonClient().URL + "/transactions/" + TxnHash + "/operations",
				LabUrl: commons.GetStellarLaboratoryClient() + "/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
					text + "%3D%3D&network=" + commons.GetHorizonClientNetworkName(),
				Identifier:     TxnBody.Identifier,
				Blockchain:		"Stellar",
				TdpId:          TxnBody.TdpId,
				DataHash:       TxnBody.DataHash,
				Timestamp:      timestamp,
				TxnType:        GetTransactiontype(TxnBody.TxnType),
				FeePaid:        feePaid,
				Ledger:         ledger,
				SourceAccount:  TxnBody.PublicKey,
				From:           from,
				SequenceNo:     TxnBody.SequenceNo,
				AvailableProof: GetProofName(TxnBody.TxnType),
				To:             to,
				ProductName:    TxnBody.ProductName,
				Itemcount:      count,
				AssetCode:      assetcode}
				if(response.Message!=""&&response.Code!=0){
					log.Error(response.Message)
					w.WriteHeader(response.Code)
					json.NewEncoder(w).Encode(response)
					}else{
						result = append(result, temp)
					}
			result = append(result, temp)
		}
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].Timestamp > result[j].Timestamp
		})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return
	}
}

func RetriveTransactionId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	var result []model.PrevTxnResponse
	object := dao.Connection{}
	p := object.GetAllTransactionForTxId(vars["id"])
	fmt.Println(vars["id"])
	p.Then(func(data interface{}) interface{} {
		res := data.([]model.TransactionCollectionBody)
		for _, TxnBody := range res {
			TxnHash := TxnBody.TxnHash

			var txe xdr.Transaction
			status := "success"
			timestamp := ""
			ledger := ""
			feePaid := ""
			from := ""
			to := ""

			result1, err := http.Get(commons.GetHorizonClient().URL + "/transactions/" + TxnHash)
			if err != nil {
				status = "Txn Id Not Found in Stellar Public Net"
			}
			data, _ := ioutil.ReadAll(result1.Body)
			if result1.StatusCode != 200 {
				status = "Txn Id Not Found in Stellar Public Net"
			}

			if status == "success" {

				var raw map[string]interface{}
				json.Unmarshal(data, &raw)
				timestamp = fmt.Sprintf("%s", raw["created_at"])
				ledger = fmt.Sprintf("%.0f", raw["ledger"])
				feePaid = fmt.Sprintf("%s", raw["fee_charged"])
				from = fmt.Sprintf("%s", raw["source_account"])
				to = fmt.Sprintf("%s", raw["source_account"])

				errXDR := xdr.SafeUnmarshalBase64(fmt.Sprintf("%s", raw["envelope_xdr"]), &txe)

				if errXDR != nil {
					//ignore error
				}

				if TxnBody.TxnType == "10" {
					to = txe.Operations[3].Body.PaymentOp.Destination.Address()
				}
			}

			mapD := map[string]string{"transaction": TxnHash}
			mapB, _ := json.Marshal(mapD)
			// fmt.Println(string(mapB))
			// trans := transaction{transaction:TxnHash}
			// s := fmt.Sprintf("%v", trans)

			encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
			text := encoded

			temp := model.PrevTxnResponse{
				Status: status, Txnhash: TxnHash,
				Url: commons.GetHorizonClient().URL + "/transactions/" + TxnHash + "/operations",
				LabUrl: commons.GetStellarLaboratoryClient() + "/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
					text + "%3D%3D&network=" + commons.GetHorizonClientNetworkName(),
				Identifier:     TxnBody.Identifier,
				TdpId:          TxnBody.TdpId,
				DataHash:       TxnBody.DataHash,
				Timestamp:      timestamp,
				TxnType:        GetTransactiontype(TxnBody.TxnType),
				FeePaid:        feePaid,
				Ledger:         ledger,
				SourceAccount:  TxnBody.PublicKey,
				From:           from,
				SequenceNo:     TxnBody.SequenceNo,
				AvailableProof: GetProofName(TxnBody.TxnType),
				To:             to}

			result = append(result, temp)
		}

		// res := TDP{TdpId: result.TdpId}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return nil
	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Txn Id Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error
	})
	p.Await()

}

/*GetCOCByTxn - WORKING MODEL
@author - Azeem Ashraf
@desc - Returns the Txn ID of the last COC Txn
@params - ResponseWriter,Request
*/
func GetCOCByTxn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	object := dao.Connection{}
	p := object.GetCOCByTxn(vars["txn"])
	p.Then(func(data interface{}) interface{} {

		result := data.(model.COCCollectionBody)
		// res := model.LastTxnResponse{LastTxn: result.TxnHash}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return nil
	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Txn Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error
	})
	p.Await()

}

func checkValidVersionByte(key string) string {

	version, er := strkey.Version(key)
	if er != nil {
	}

	if version == strkey.VersionByteAccountID {
		return "pk"
	}

	if version == strkey.VersionByteSeed {
		return "sk"
	}

	// if version == strkey.VersionByteHashTx {
	// 	return "txn"
	// }

	// if version == strkey.VersionByteHashX {
	// return "hash"
	// }

	matched, err := regexp.MatchString(`^[0-9a-f]{64}$`, key)
	if err != nil {
	}
	if matched {
		return "txn"
	}
	matched1, err1 := regexp.MatchString(`^[0-9a-f]{24}$`, key)
	if err1 != nil {
	}
	if matched1 {
		return "tdpid"
	}
	return ""
}

func RetrievePreviousTranasctionsCount(w http.ResponseWriter, r *http.Request) {
	object := dao.Connection{}
	var totalRecordCount model.TotalTransaction;
	_, err := object.GetTransactionCount().Then(func(data interface{}) interface{} {
		totalRecordCount = model.TotalTransaction{TotalTransactionCount:data.(int64)}
		w.WriteHeader(http.StatusOK)
		return json.NewEncoder(w).Encode(totalRecordCount)
	}).Await()
	if err != nil {
		log.Error("Unable Reach to Database",err)
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(model.Error{Code:http.StatusNoContent,Message:"Unable Reach to Database"})
		return
	}
}

//RetrievePreviousTranasctions ...
func RetrievePreviousTranasctions(w http.ResponseWriter, r *http.Request) {
	var response model.Error;
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	key1, error := r.URL.Query()["perPage"]

	if !error || len(key1[0]) < 1 {
		log.Error("Url Parameter 'perPage' is missing")
		return
	}

	key2, error := r.URL.Query()["page"]

	if !error || len(key2[0]) < 1 {
		log.Error("Url Parameter 'page' is missing")
		return
	}

	key3, error := r.URL.Query()["NoPage"]

	if !error || len(key3[0]) < 1 {
		log.Error("Url Parameter 'NoPage' is missing")
		return
	}

	perPage, err := strconv.Atoi(key1[0])
	if err != nil {
		log.Error("Query parameter error" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		response = model.Error{Code:http.StatusBadRequest,Message: "The parameter should be an integer " + err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	page, err := strconv.Atoi(key2[0])
	if err != nil {
		log.Error("Query parameter error" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		response = model.Error{Code:http.StatusBadRequest,Message: "The parameter should be an integer " + err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	NoPage, err := strconv.Atoi(key3[0])
	if err != nil {
		log.Error("Query parameter error" + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		response = model.Error{Code:http.StatusBadRequest,Message: "The parameter should be an integer " + err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	var result []model.PrevTxnResponse
	object := dao.Connection{}

	_, err = object.GetPreviousTransactions(perPage, page, NoPage).Then(func(data interface{}) interface{} {
		res := data.([]model.TransactionCollectionBody)
		for _, TxnBody := range res {
			if TxnBody.TxnType != "11"{
				TxnHash := TxnBody.TxnHash
				var txe xdr.Transaction
				status := "success"
				timestamp := ""
				ledger := ""
				feePaid := ""
				from := ""
				to := ""
				assetcode := ""
				result1, err := http.Get(commons.GetHorizonClient().URL + "/transactions/" + TxnHash)
				if err!=nil {
					log.Error("Unable to reach Stellar network in result1")
					status = "Unable to reach Stellar network"
					w.WriteHeader(http.StatusNotFound)		
					response=model.Error{Code:http.StatusNotFound,Message: "Unable to reach Stellar network"}
					return json.NewEncoder(w).Encode(response)
				}	
				if result1.StatusCode != 200 {
					log.Error("Transaction could not be retrieved from Stellar Network in result1")
					status = "Transaction could not be retrieved from Stellar Network"
					w.WriteHeader(http.StatusNoContent)		
					response=model.Error{Code:http.StatusNoContent,Message: "Transaction could not be retrieved from Stellar Network"}
					return json.NewEncoder(w).Encode(response)
				}
				data, _ := ioutil.ReadAll(result1.Body)
				if status == "success" {
					var raw map[string]interface{}
					json.Unmarshal(data, &raw)
					timestamp = fmt.Sprintf("%s", raw["created_at"])
					ledger = fmt.Sprintf("%.0f", raw["ledger"])
					feePaid = fmt.Sprintf("%s", raw["fee_charged"])
					from = fmt.Sprintf("%s", raw["source_account"])
					to = fmt.Sprintf("%s", raw["source_account"])
					errXDR := xdr.SafeUnmarshalBase64(fmt.Sprintf("%s", raw["envelope_xdr"]), &txe)
					if errXDR != nil {
						//log.Error("Error SafeUnmarshalBase64 " + errXDR.Error())
					}
					if TxnBody.TxnType == "10" {
						result1, err := http.Get(commons.GetHorizonClient().URL + "/transactions/" + TxnHash + "/operations")
						if err != nil {
							log.Error("Unable to reach Stellar network in result1")
							status = "Unable to reach Stellar network"
							w.WriteHeader(http.StatusNotFound)		
							response=model.Error{Code:http.StatusNotFound,Message: "Unable to reach Stellar network"}
							return json.NewEncoder(w).Encode(response)			
						}
						if result1.StatusCode != 200 {
							log.Error("Transaction could not be retrieved from Stellar Network in result1")
							status = "Transaction could not be retrieved from Stellar Network"
							w.WriteHeader(http.StatusNoContent)		
							response=model.Error{Code:http.StatusNoContent,Message: "Transaction could not be retrieved from Stellar Network"}
							return json.NewEncoder(w).Encode(response)
						}
						data, err := ioutil.ReadAll(result1.Body)
						if err != nil {
							log.Error("Unable to reading response")
							response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to reading response"}
						}
						var raw map[string]interface{}
						err = json.Unmarshal(data, &raw)
						if err != nil {
							log.Error("Unable to unmarshal row data")
							response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal data"}
						}
						out, err := json.Marshal(raw["_embedded"])
						if err != nil {
							log.Error("Unable to marshal embedded")
							response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to marshal embedded"}
						}
						var raw1 map[string]interface{}
						err = json.Unmarshal(out, &raw1)
						if err != nil {
							log.Error("Unable to unmarshal  json.Unmarshal(out, &raw1)")
							response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal data"}	
						}
						out1, err := json.Marshal(raw1["records"])
						if err != nil {
							log.Error("Unable to marshal records ")
							response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to marshal records"}
						}
						keysBody := out1
						keys := make([]PublicKeyPOCOC, 0)
						err = json.Unmarshal(keysBody, &keys)
						if err != nil {
							log.Error("Unable to unmarshal keys data")
							response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal key's data"}
						}
						acceptTxn_byteData, err := base64.StdEncoding.DecodeString(keys[2].Value)
						if err != nil {
							log.Error("Unable to base64 decoding")
							response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to base64 decoding"}
						}
						acceptTxn := string(acceptTxn_byteData)
						log.Info("acceptTxn: " + acceptTxn)
						acceptresult1, err := http.Get(commons.GetHorizonClient().URL + "/transactions/" + acceptTxn + "/operations")
						if err != nil {
							log.Error("Unable to reach Stellar network in acceptresult1")
							status = "Unable to reach Stellar network"
							w.WriteHeader(http.StatusNotFound)		
							response=model.Error{Code:http.StatusNotFound,Message: "Unable to reach Stellar network"}
							return json.NewEncoder(w).Encode(response)			
						}
						if acceptresult1.StatusCode != 200 {
							log.Error("Transaction could not be retrieved from Stellar Network in acceptresult1")
							status = "Transaction could not be retrieved from Stellar Network"
							w.WriteHeader(http.StatusNoContent)		
							response=model.Error{Code:http.StatusNoContent,Message: "Transaction could not be retrieved from Stellar Network"}
							return json.NewEncoder(w).Encode(response)
						}
						acceptdata, err := ioutil.ReadAll(acceptresult1.Body)
						if err != nil {
							log.Error("Unable to reading response",err)
							response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to reading response"}
						}
						var acceptraw map[string]interface{}
						err = json.Unmarshal(acceptdata, &acceptraw)
						if err != nil {
							log.Error("Unable to unmarshal data")
							response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal data"}
						}
						acceptout, err := json.Marshal(acceptraw["_embedded"])
						if err != nil {
							log.Error("Unable to marshal _embedded ")
							response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to marshal _embedded"}
						}
						var acceptraw1 map[string]interface{}
						err = json.Unmarshal(acceptout, &acceptraw1)
						if err != nil {
							log.Error("Unable to unmarshal data")
							response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal data"}
						}
						acceptout1, err := json.Marshal(acceptraw1["records"])
						if err != nil {
							log.Error("Unable to marshal records")
							response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to marshal records"}
						}
						acceptkeysBody := acceptout1
						acceptkeys := make([]PublicKeyPOCOC, 0)
						err = json.Unmarshal(acceptkeysBody, &acceptkeys)
						if err != nil {
							log.Error("Unable to unmarshal key's data")
							response=model.Error{Code:http.StatusInternalServerError,Message: "Unable to unmarshal key's data"}
						}
						to = string(acceptkeys[3].To)
						log.Info("Destination: " + to)

						from = string(acceptkeys[3].From)
						log.Info("Source: " + from)

						assetcode = string(acceptkeys[3].Asset_code)
						log.Info("Source: " + assetcode)
					}
				}
				//mapD := map[string]string{"transaction": TxnHash}
				//mapB, err := json.Marshal(mapD)
				//if err != nil {
				//log.Error("Error while json.Marshal(mapD) " + err.Error())
				//}
				// fmt.Println(string(mapB))
				// trans := transaction{transaction:TxnHash}
				// s := fmt.Sprintf("%v", trans)
				//encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
				//text := encoded
				temp := model.PrevTxnResponse{
					Status: status, Txnhash: TxnHash,
					Url:            commons.GetHorizonClient().URL + "/transactions/" + TxnHash + "/operations",
					Identifier:     TxnBody.Identifier,
					TdpId:          TxnBody.TdpId,
					DataHash:       TxnBody.DataHash,
					Timestamp:      timestamp,
					TxnType:        GetTransactiontype(TxnBody.TxnType),
					FeePaid:        feePaid,
					Ledger:         ledger,
					SourceAccount:  TxnBody.PublicKey,
					From:           from,
					SequenceNo:     TxnBody.SequenceNo,
					AvailableProof: GetProofName(TxnBody.TxnType),
					To:             to,
					ProductName:    TxnBody.ProductName,
					AssetCode:      assetcode}
				result = append(result, temp)
			}
		}
		// res := TDP{TdpId: result.TdpId}
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].Timestamp > result[j].Timestamp
		})
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return nil
	}).Await()

	if err != nil {
		if(response.Message!=""&&response.Code!=0){
		log.Error(response.Message)
		w.WriteHeader(response.Code)
		json.NewEncoder(w).Encode(response)
		}else{
			log.Error("No Transactions Found in Gateway DataStore")
			w.WriteHeader(http.StatusNoContent)
			json.NewEncoder(w).Encode(model.Error{Code:http.StatusNoContent,Message:"No Transactions Found in Gateway DataStore"})
		}
	}
}

func GetTransactiontype(Type string) string {
	switch Type {
	case "0":
		return "genesis"
	case "2":
		return "tdp"
	case "5":
		return "splitParent"
	case "6":
		return "splitChild"
	case "7":
		return "merge"
	case "8":
		return "merge"
	case "10":
		return "coc"
	case "11":
		return "cocProof"
	}
	return Type
}

func GetProofName(Type string) []string {
	var result []string
	switch Type {
	case "0":
		result = append(result, "pog")
	case "2":
		result = append(result, "poe")
		result = append(result, "poc")
	case "5":
		result = append(result, "poc")
	case "6":
		result = append(result, "poc")
	case "7":
		result = append(result, "poc")
	case "8":
		result = append(result, "poc")
	case "10":
		result = append(result, "pococ")
		result = append(result, "poc")

	case "11":
		result = append(result, "pococ")
	}
	return result
}
