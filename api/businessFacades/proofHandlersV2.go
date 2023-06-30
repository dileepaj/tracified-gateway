package businessFacades

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/interpreter"
	"github.com/gorilla/mux"
)

type test struct {
	Data string
}

/*CheckPOE - Deprecated
@author - Azeem Ashraf, Jajeththanan Sabapathipillai
*/
func CheckPOE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	object := dao.Connection{}

	p := object.GetTransactionForTdpId(vars["Txn"])
	p.Then(func(data interface{}) interface{} {

		result := data.(model.TransactionCollectionBody)
		var response model.POE
		// url := "http://localhost:3001/api/v1/dataPackets/raw?id=" + vars["Txn"]
		url := constants.TracifiedBackend + constants.RawTDP + vars["Txn"]

		bearer := "Bearer " + constants.BackendToken
		// Create a new request using http
		req, er := http.NewRequest("GET", url, nil)

		req.Header.Add("Authorization", bearer)
		client := &http.Client{}
		resq, er := client.Do(req)

		if er != nil {

			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(er.Error)

		} else {
			body, _ := ioutil.ReadAll(resq.Body)
			var raw map[string]interface{}
			json.Unmarshal(body, &raw)

			h := sha256.New()
			lol := raw["data"]

			h.Write([]byte(fmt.Sprintf("%s", lol) + result.Identifier))

			fmt.Printf("%x", h.Sum(nil))

			poeStructObj := apiModel.POEStruct{Txn: result.TxnHash,
				Hash: strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))}
			display := &interpreter.AbstractPOE{POEStruct: poeStructObj}
			response = display.InterpretPOE("") //parameter was added here.

			w.WriteHeader(response.RetrievePOE.Error.Code)
			json.NewEncoder(w).Encode(response.RetrievePOE)

		}

		return data

	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusNotFound)
		response := model.Error{Message: "Not Found"}
		json.NewEncoder(w).Encode(response)
		return error

	})
	p.Await()

	// return

}

/*CheckPOC - Deprecated
@author - Azeem Ashraf
*/
func CheckPOC(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	object := dao.Connection{}
	var response model.POC
	var pocStructObj apiModel.POCStruct

	p := object.GetTransactionForTdpId(vars["Txn"])
	p.Then(func(data interface{}) interface{} {

		result := data.(model.TransactionCollectionBody)
		pocStructObj.DBTree = []model.Current{}
		g := object.GetTransactionsbyIdentifier(result.Identifier)
		g.Then(func(data interface{}) interface{} {
			res := data.([]model.TransactionCollectionBody)
			pocStructObj.Txn = res[len(res)-1].TxnHash

			for i := len(res) - 1; i >= 0; i-- {
				if res[i].TxnType == "2" {
					// url := "http://localhost:3001/api/v1/dataPackets/raw?id=" + res[i].TdpId
					url := constants.TracifiedBackend + constants.RawTDP + res[i].TdpId

					bearer := "Bearer " + constants.BackendToken
					// Create a new request using http
					req, er := http.NewRequest("GET", url, nil)

					req.Header.Add("Authorization", bearer)
					client := &http.Client{}
					resq, er := client.Do(req)

					if er != nil {
						w.Header().Set("Content-Type", "application/json; charset=UTF-8")

						w.WriteHeader(http.StatusOK)
						response := model.Error{Message: "Connection to the DataStore was interupted"}
						json.NewEncoder(w).Encode(response)
					} else {
						body, _ := ioutil.ReadAll(resq.Body)
						var raw map[string]interface{}
						json.Unmarshal(body, &raw)

						h := sha256.New()
						base64 := raw["data"]

						h.Write([]byte(fmt.Sprintf("%s", base64) + result.Identifier))
						// fmt.Printf("%x", h.Sum(nil))

						DataStoreTXN := model.Current{
							TType:      res[i].TxnType,
							TXNID:      res[i].TxnHash,
							Identifier: res[i].Identifier,
							DataHash:   strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))}

						pocStructObj.DBTree = append(pocStructObj.DBTree, DataStoreTXN)
					}
				} else {
					DataStoreTXN := model.Current{
						TType:      res[i].TxnType,
						TXNID:      res[i].TxnHash,
						Identifier: res[i].Identifier,
					}
					pocStructObj.DBTree = append(pocStructObj.DBTree, DataStoreTXN)
				}

			}

			// pocStructObj = apiModel.POCStruct{

			// }
			display := &interpreter.AbstractPOC{POCStruct: pocStructObj}
			response = display.InterpretPOC()

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(200)
			// w.WriteHeader(http.StatusBadRequest)

			// result := apiModel.PoeSuccess{Message: "response.RetrievePOC.Error.Message", TxNHash: "response.RetrievePOC.Txn"}

			result := apiModel.PocSuccess{Message: response.RetrievePOC.Error.Message, Chain: pocStructObj.DBTree}

			json.NewEncoder(w).Encode(result)
			return data
		}).Catch(func(error error) error {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")

			w.WriteHeader(http.StatusOK)
			response := model.Error{Message: "Identifier for the TDP ID Not Found in Gateway DataStore"}
			json.NewEncoder(w).Encode(response)
			return error
		})
		g.Await()

		return data

	}).Catch(func(error error) error {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		w.WriteHeader(http.StatusOK)
		response := model.Error{Message: "TDP ID Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error

	})
	p.Await()

	// return

}

/*CheckFullPOC - Deprecated
@author - Azeem Ashraf
*/
func CheckFullPOC(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	object := dao.Connection{}
	var response model.POC
	var pocStructObj apiModel.POCStruct

	p := object.GetTransactionForTdpId(vars["Txn"])
	p.Then(func(data interface{}) interface{} {

		result := data.(model.TransactionCollectionBody)
		pocStructObj.DBTree = []model.Current{}
		g := object.GetTransactionsbyIdentifier(result.Identifier)
		g.Then(func(data interface{}) interface{} {
			res := data.([]model.TransactionCollectionBody)
			pocStructObj.Txn = res[len(res)-1].TxnHash

			for i := len(res) - 1; i >= 0; i-- {
				// url := "http://localhost:3001/api/v1/dataPackets/raw?id=" + res[i].TdpId
				url := constants.TracifiedBackend + constants.RawTDP + res[i].TdpId

				bearer := "Bearer " + constants.BackendToken
				// Create a new request using http
				req, er := http.NewRequest("GET", url, nil)

				req.Header.Add("Authorization", bearer)
				client := &http.Client{}
				resq, er := client.Do(req)

				if er != nil {
					w.WriteHeader(http.StatusNotFound)
					json.NewEncoder(w).Encode(er.Error)
				} else {
					body, _ := ioutil.ReadAll(resq.Body)
					var raw map[string]interface{}
					json.Unmarshal(body, &raw)

					h := sha256.New()
					base64 := raw["data"]

					h.Write([]byte(fmt.Sprintf("%s", base64) + result.Identifier))
					// fmt.Printf("%x", h.Sum(nil))

					DataStoreTXN := model.Current{
						TType:      "2",
						TXNID:      res[i].TxnHash,
						Identifier: res[i].Identifier,
						DataHash:   strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))}

					pocStructObj.DBTree = append(pocStructObj.DBTree, DataStoreTXN)
				}
			}

			// pocStructObj = apiModel.POCStruct{

			// }
			display := &interpreter.AbstractPOC{POCStruct: pocStructObj}
			response = display.InterpretPOC()

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(200)
			// w.WriteHeader(http.StatusBadRequest)

			// result := apiModel.PoeSuccess{Message: "response.RetrievePOC.Error.Message", TxNHash: "response.RetrievePOC.Txn"}

			result := apiModel.PocSuccess{Message: response.RetrievePOC.Error.Message, Chain: pocStructObj.DBTree}
			json.NewEncoder(w).Encode(result)

			return data
		}).Catch(func(error error) error {
			return error
		}).Await()

		return data

	}).Catch(func(error error) error {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		w.WriteHeader(http.StatusOK)
		response := model.Error{Message: "Identifier Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error

	}).Await()

	// return

}

/*CheckPOG - Deprecated
@author - Azeem Ashraf, Jajeththanan Sabapathipillai
*/
func CheckPOG(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var response model.POG

	object := dao.Connection{}
	//RETRIVE LAST TRANSACTION HASH FOR THE IDENTIFIER
	p := object.GetLastTransactionbyIdentifier(vars["Identifier"])
	p.Then(func(data interface{}) interface{} {

		LastTxn := data.(model.TransactionCollectionBody)

		//RETRIVE FIRST TRANSACTION HASH FOR THE IDENTIFIER
		g := object.GetFirstTransactionbyIdentifier(vars["Identifier"])
		g.Then(func(data interface{}) interface{} {

			FirstTxn := data.(model.TransactionCollectionBody)

			pogStructObj := apiModel.POGStruct{LastTxn: LastTxn.TxnHash, POGTxn: FirstTxn.TxnHash, Identifier: vars["Identifier"]}
			display := &interpreter.AbstractPOG{POGStruct: pogStructObj}
			response = display.InterpretPOG()

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(response.RetrievePOG.Message.Code)
			json.NewEncoder(w).Encode(response.RetrievePOG)
			return nil
		}).Catch(func(error error) error {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")

			w.WriteHeader(http.StatusOK)
			response := model.Error{Message: "Identifier Not Found in Gateway DataStore"}
			json.NewEncoder(w).Encode(response)
			return error
		})
		g.Await()

		return nil
	}).Catch(func(error error) error {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		w.WriteHeader(http.StatusOK)
		response := model.Error{Message: "Identifier Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error
	})
	p.Await()

	return

}
