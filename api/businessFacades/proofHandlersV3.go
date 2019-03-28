package businessFacades

import (
	"crypto/sha256"
	"encoding/base64"
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
	"github.com/stellar/go/xdr"
)

type PublicKey struct {
	Name  string
	Value string
}

type KeysResponse struct {
	Collection []PublicKey
}

//CheckPOEV3 - WORKING MODEL
func CheckPOEV3(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	var result model.TransactionCollectionBody
	object := dao.Connection{}
	var CurrentTxn string
	p := object.GetTransactionForTdpId(vars["Txn"])
	p.Then(func(data interface{}) interface{} {
		result = data.(model.TransactionCollectionBody)
		return nil
	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "TDPID NOT FOUND IN DATASTORE"}
		json.NewEncoder(w).Encode(response)
		fmt.Println(response)
		return error
	})
	p.Await()

	result1, err := http.Get("https://horizon-testnet.stellar.org/transactions/" + result.TxnHash + "/operations")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Txn for the TXN does not exist in the Blockchain"}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, _ := ioutil.ReadAll(result1.Body)

	// if result1.StatusCode == 200 {
	var raw map[string]interface{}
	json.Unmarshal(data, &raw)
	// raw["count"] = 2
	out, _ := json.Marshal(raw["_embedded"])
	var raw1 map[string]interface{}
	json.Unmarshal(out, &raw1)
	out1, _ := json.Marshal(raw1["records"])

	keysBody := out1
	keys := make([]PublicKey, 0)
	json.Unmarshal(keysBody, &keys)

	byteData, _ := base64.StdEncoding.DecodeString(keys[2].Value)

	CurrentTxn = string(byteData)
	fmt.Println("THE TXN OF THE USER TXN: " + CurrentTxn)
	// }F

	// fmt.Println(result)
	var response model.POE
	// url := "http://localhost:3001/api/v2/dataPackets/raw?id=5c9141b2618cf404ec5e105d"
	url := constants.TracifiedBackend + constants.RawTDP + vars["Txn"]

	bearer := "Bearer " + constants.BackendToken
	// Create a new request using http
	req, er := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	resq, er := client.Do(req)

	if er != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Connection to the Traceability DataStore was interupted"}
		json.NewEncoder(w).Encode(response)
		return
	}

	body, _ := ioutil.ReadAll(resq.Body)
	// if resq.StatusCode == 200 {
	var raw2 map[string]interface{}
	json.Unmarshal(body, &raw2)
	// fmt.Println(string(raw["Data"]))
	// fmt.Println(body)

	h := sha256.New()
	lol := raw2["data"]

	fmt.Println(raw2["data"])

	h.Write([]byte(fmt.Sprintf("%s", lol) + result.Identifier))
	fmt.Println("RAW BASE64 + IDENTIFIER")

	fmt.Printf("%x\n", h.Sum(nil))

	poeStructObj := apiModel.POEStruct{Txn: result.TxnHash,
		Hash: strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))}
	display := &interpreter.AbstractPOE{POEStruct: poeStructObj}
	response = display.InterpretPOE()

	w.WriteHeader(response.RetrievePOE.Error.Code)
	json.NewEncoder(w).Encode(response.RetrievePOE)
	// w.WriteHeader(http.StatusOK)
	// res := model.Error{Message: fmt.Sprintf("%s", lol)}
	// json.NewEncoder(w).Encode(res)
	// }

	return

}

//CheckPOCV3 - NEEDS TO BE TESTED
func CheckPOCV3(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	object := dao.Connection{}
	var response model.POC
	var pocStructObj apiModel.POCStruct

	p := object.GetTransactionForTdpId(vars["Txn"])
	p.Then(func(data interface{}) interface{} {

		result := data.(model.TransactionCollectionBody)
		pocStructObj.DBTree = []model.Current{}
		// fmt.Println(result)
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

						w.WriteHeader(http.StatusBadRequest)
						response := model.Error{Message: "Connection to the DataStore was interupted"}
						json.NewEncoder(w).Encode(response)
					} else {
						// fmt.Println(req)
						body, _ := ioutil.ReadAll(resq.Body)
						var raw map[string]interface{}
						json.Unmarshal(body, &raw)

						h := sha256.New()
						base64 := raw["data"]
						// fmt.Println(base64)

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

			// fmt.Println(response.RetrievePOC.Error.Message)

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(200)
			// w.WriteHeader(http.StatusBadRequest)

			// result := apiModel.PoeSuccess{Message: "response.RetrievePOC.Error.Message", TxNHash: "response.RetrievePOC.Txn"}

			result := apiModel.PocSuccess{Message: response.RetrievePOC.Error.Message, Chain: pocStructObj.DBTree}
			fmt.Println(result)
			fmt.Println(response.RetrievePOC.Error.Message)

			json.NewEncoder(w).Encode(result)
			// 		return

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

//CheckFullPOCV3 - NEEDS TO BE TESTED
func CheckFullPOCV3(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	object := dao.Connection{}
	var response model.POC
	var pocStructObj apiModel.POCStruct

	p := object.GetTransactionForTdpId(vars["Txn"])
	p.Then(func(data interface{}) interface{} {

		result := data.(model.TransactionCollectionBody)
		pocStructObj.DBTree = []model.Current{}
		// fmt.Println(result)
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
					// fmt.Println(req)
					body, _ := ioutil.ReadAll(resq.Body)
					var raw map[string]interface{}
					json.Unmarshal(body, &raw)

					h := sha256.New()
					base64 := raw["data"]
					// fmt.Println(base64)

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

			// fmt.Println(response.RetrievePOC.Error.Message)

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(200)
			// w.WriteHeader(http.StatusBadRequest)

			// result := apiModel.PoeSuccess{Message: "response.RetrievePOC.Error.Message", TxNHash: "response.RetrievePOC.Txn"}

			result := apiModel.PocSuccess{Message: response.RetrievePOC.Error.Message, Chain: pocStructObj.DBTree}
			fmt.Println(result)
			fmt.Println(response.RetrievePOC.Error.Message)

			json.NewEncoder(w).Encode(result)
			// 		return

			return data
		}).Catch(func(error error) error {
			return error
		})
		g.Await()

		return data

	}).Catch(func(error error) error {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		w.WriteHeader(http.StatusOK)
		response := model.Error{Message: "Identifier Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error

	})
	p.Await()

	// return

}

//CheckPOGV3 - WORKING MODEL
func CheckPOGV3(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var response model.POG
	var UserGenesis string

	object := dao.Connection{}
	p := object.GetLastTransactionbyIdentifier(vars["Identifier"])
	p.Then(func(data interface{}) interface{} {

		LastTxn := data.(model.TransactionCollectionBody)
		fmt.Println(LastTxn)
		g := object.GetFirstTransactionbyIdentifier(vars["Identifier"])
		g.Then(func(data interface{}) interface{} {

			FirstTxnGateway := data.(model.TransactionCollectionBody)

			// fmt.Println(FirstTxnGateway)
			fmt.Println("First TXN SIGNED BY GATEWAY IS USED TO REQUEST THE USER's GENESIS")
			result1, err := http.Get("https://horizon-testnet.stellar.org/transactions/" + FirstTxnGateway.TxnHash + "/operations")
			if err != nil {

			} else {
				data, _ := ioutil.ReadAll(result1.Body)

				if result1.StatusCode == 200 {
					var raw map[string]interface{}
					json.Unmarshal(data, &raw)
					// raw["count"] = 2
					out, _ := json.Marshal(raw["_embedded"])
					var raw1 map[string]interface{}
					json.Unmarshal(out, &raw1)
					out1, _ := json.Marshal(raw1["records"])

					keysBody := out1
					keys := make([]PublicKey, 0)
					json.Unmarshal(keysBody, &keys)

					//GET THE USER SIGNED GENESIS TXN
					byteData, _ := base64.StdEncoding.DecodeString(keys[2].Value)
					UserGenesis = string(byteData)
					fmt.Println("THE TXN OF THE USER TXN: " + UserGenesis)
				}
			}

			pogStructObj := apiModel.POGStruct{LastTxn: LastTxn.TxnHash, POGTxn: UserGenesis, Identifier: vars["Identifier"]}
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

	// fmt.Println("response.RetrievePOG.Error.Code")
	// fmt.Println(response.RetrievePOG.Error.Code)

	return

}

//CheckPOCOCV3 - Developing
func CheckPOCOCV3(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var txe xdr.Transaction
	var COC model.COCCollectionBody
	var COCAvailable bool
	vars := mux.Vars(r)
	object := dao.Connection{}
	p := object.GetCOCbyAcceptTxn(vars["TxnId"])
	p.Then(func(data interface{}) interface{} {
		COCAvailable = true
		COC = data.(model.COCCollectionBody)
		fmt.Println(COC)
		return data

	}).Catch(func(error error) error {
		COCAvailable = false
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "COCTXN NOT FOUND IN GATEWAY DATASTORE"}
		json.NewEncoder(w).Encode(response)
		fmt.Println(response)
		return error

	})
	p.Await()

	if COCAvailable {
		err := xdr.SafeUnmarshalBase64(COC.AcceptXdr, &txe)
		if err != nil {
			//ignore error
		}
		display := &interpreter.AbstractPOCOC{Txn: vars["TxnId"], DBCOC: txe}
		display.InterpretPOCOC(w, r)
	}

	return

}
