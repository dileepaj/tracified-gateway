package businessFacades

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	//"go/constant"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/deprecatedBuilder"
	"github.com/dileepaj/tracified-gateway/proofs/retriever/stellarRetriever"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/gorilla/mux"
	// "github.com/hpcloud/tail"
	// "github.com/dileepaj/tracified-gateway/proofs/builder"
)

/*CreateTrust deprecated
@author - Sharmilan Somasundaram
*/
func CreateTrust(w http.ResponseWriter, r *http.Request) {

	// var response model.POE
	var TObj apiModel.TrustlineStruct
	logger := utilities.NewCustomLogger()
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
			logger.LogWriter("Error while Decoding the body  :"+err.Error(), constants.ERROR)
			return
		}
		display := &deprecatedBuilder.AbstractTrustline{TrustlineStruct: TObj}
		result := display.Trustline()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		result2 := apiModel.PoeSuccess{Message: "TrustLine Created", TxNHash: result}
		json.NewEncoder(w).Encode(result2)
		return
	}
	return
}

/*SendAssests deprecated
@author - Sharmilan Somasundaram
*/
func SendAssests(w http.ResponseWriter, r *http.Request) {

	// var response model.POE
	var TObj apiModel.SendAssest
	logger := utilities.NewCustomLogger()

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
			logger.LogWriter("Error while Decoding the body  :"+err.Error(), constants.ERROR)
			return
		}
		display := &deprecatedBuilder.AbstractAssetTransfer{SendAssest: TObj}
		response := display.AssetTransfer()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(response.Error.Code)
		result := apiModel.SendAssetRes{Message: response.Error.Message, PreviousTXNID: response.PreviousTXNID, PreviousProfileID: response.PreviousProfileID, Code: response.Code, Amount: response.Amount, Txn: response.Txn, To: response.To, From: response.From}
		json.NewEncoder(w).Encode(result)
		return
	}
	return
}

/*MultisigAccount deprecated
@author - Sharmilan Somasundaram
*/
func MultisigAccount(w http.ResponseWriter, r *http.Request) {

	var TObj apiModel.RegistrarAccount
	logger := utilities.NewCustomLogger()
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
			logger.LogWriter("Error while Decoding the body  :"+err.Error(), constants.ERROR)
			return
		}

		// var response model.POE

		display := &deprecatedBuilder.AbstractCreateRegistrar{RegistrarAccount: TObj}
		result := display.CreateRegistrarAcc()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		result2 := apiModel.PoeSuccess{Message: "Success", TxNHash: result}
		json.NewEncoder(w).Encode(result2)
		return
	}
	return

}

/*AppointRegistrar deprecated
@author - Sharmilan Somasundaram
*/
func AppointRegistrar(w http.ResponseWriter, r *http.Request) {

	var TObj apiModel.AppointRegistrar
	logger := utilities.NewCustomLogger()
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
			logger.LogWriter("Error while Decoding the body  :"+err.Error(), constants.ERROR)
			return
		}

		display := &deprecatedBuilder.AbstractAppointRegistrar{AppointRegistrar: TObj}
		result := display.AppointReg()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		result2 := apiModel.RegSuccess{Message: "Success", Xdr: result}
		json.NewEncoder(w).Encode(result2)
		return
	}
	return
}

/*TransformV2 deprecated
@author - Sharmilan Somasundaram
*/
func TransformV2(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)

	// var response model.POE
	var TObj apiModel.AssetTransfer
	logger := utilities.NewCustomLogger()
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
			logger.LogWriter("Error while Decoding the body  :"+err.Error(), constants.ERROR)
			return
		}
		display := &deprecatedBuilder.AbstractTransformAssets{AssetTransfer: TObj}
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

/*COC deprecated
@author - Sharmilan Somasundaram
*/
func COC(w http.ResponseWriter, r *http.Request) {

	// var response model.POE
	var TObj apiModel.ChangeOfCustody
	logger := utilities.NewCustomLogger()
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
			logger.LogWriter("Error while Decoding the body  :"+err.Error(), constants.ERROR)
			return
		}
		display := &deprecatedBuilder.AbstractCoCTransaction{ChangeOfCustody: TObj}
		response := display.CoCTransaction()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(response.Error.Code)
		result2 := apiModel.COCRes{Message: response.Error.Message, PreviousTXNID: response.PreviousTXNID, PreviousProfileID: response.PreviousProfileID, Code: response.Code, Amount: response.Amount, To: response.To, From: response.From, TxnXDR: response.TxnXDR}
		json.NewEncoder(w).Encode(result2)
		return
	}
	return
}

/*COCLink deprecated
@author - Sharmilan Somasundaram
*/
func COCLink(w http.ResponseWriter, r *http.Request) {

	// var response model.POE
	var TObj apiModel.ChangeOfCustodyLink
	logger := utilities.NewCustomLogger()
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
			logger.LogWriter("Error while Decoding the body  :"+err.Error(), constants.ERROR)
			return
		}
		display := &deprecatedBuilder.AbstractcocLink{ChangeOfCustodyLink: TObj}
		result := display.CoCLink()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		result2 := apiModel.PoeSuccess{Message: "Success", TxNHash: result}
		json.NewEncoder(w).Encode(result2)
		return
	}
	return
}

/*DeveloperRetriever testing
@author - Azeem Ashraf
*/
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

type logBody struct {
	Data []string
}

/*RetrieveLogsForToday testing
@author - Azeem Ashraf
*/
func RetrieveLogsForToday(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	s := time.Now().UTC().String()
	dat, err := ioutil.ReadFile("GatewayLogs" + s[:10])
	logger := utilities.NewCustomLogger()

	if err != nil {
		logger.LogWriter("Log File is not founnd  :"+err.Error(), constants.ERROR)
		w.WriteHeader(400)
		result := apiModel.SubmitXDRSuccess{
			Status: "Log File is not found",
		}
		json.NewEncoder(w).Encode(result)
		return

	} else {
		stuff := strings.Split(string(dat), "\n")
		w.WriteHeader(200)
		result := logBody{
			Data: stuff[:len(stuff)-1],
		}
		json.NewEncoder(w).Encode(result)
		return
	}

}

/*GatewayRetriever testing
@author - Azeem Ashraf
*/
func GatewayRetriever(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// var response model.POC

	object := dao.Connection{}
	var pocStructObj apiModel.POCStruct
	logger := utilities.NewCustomLogger()
	p := object.GetTransactionForTdpId(vars["Txn"])
	p.Then(func(data interface{}) interface{} {

		result := data.(model.TransactionCollectionBody)
		pocStructObj.DBTree = []model.Current{}
		strresult,_:=json.Marshal(result)
		logger.LogWriter("TransactionCollectionBody :"+string(strresult),constants.INFO)
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
						strreq,_:=json.Marshal(req)
						logger.LogWriter("Request :"+string(strreq),constants.INFO)
						
						body, _ := ioutil.ReadAll(resq.Body)
						var raw map[string]interface{}
						json.Unmarshal(body, &raw)

						h := sha256.New()
						base64 := raw["data"]
						strbase64,_:=json.Marshal(base64)
						logger.LogWriter("raw data :"+string(strbase64),constants.INFO)

						h.Write([]byte(fmt.Sprintf("%s", base64) + result.Identifier))
					
						logger.LogWriter( h.Sum(nil),constants.INFO)

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

			// // }
			// display := &interpreter.AbstractPOC{POCStruct: pocStructObj}
			// response = display.InterpretPOC()

			//fmt.Println(response.RetrievePOC.Error.Message)
		
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(200)
			// w.WriteHeader(http.StatusBadRequest)

			// result := apiModel.PoeSuccess{Message: "response.RetrievePOC.Error.Message", TxNHash: "response.RetrievePOC.Txn"}

			result := apiModel.PocSuccess{Chain: pocStructObj.DBTree}
			strresult,_:=json.Marshal(result)
			logger.LogWriter("PocSuccess result :"+string(strresult),constants.INFO)
			// fmt.Println(response.RetrievePOC.Error.Message)
			json.NewEncoder(w).Encode(result)
			// 		return

			return data
		}).Catch(func(error error) error {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")

			w.WriteHeader(http.StatusOK)
			response := model.Error{Message: "Identifier for the TDP ID Not Found in Gateway DataStore"}
			json.NewEncoder(w).Encode(response)
			return error
		}).Await()

		return data

	}).Catch(func(error error) error {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		w.WriteHeader(http.StatusOK)
		response := model.Error{Message: "TDP ID Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error

	}).Await()

	return

}

/*GatewayRetrieverWithIdentifier testing
@author - Azeem Ashraf
*/
func GatewayRetrieverWithIdentifier(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	logger := utilities.NewCustomLogger()

	// var response model.POC

	object := dao.Connection{}
	var pocStructObj apiModel.POCStruct

	// p := object.GetTransactionForTdpId(vars["Txn"])
	// p.Then(func(data interface{}) interface{} {

	// 	result := data.(model.TransactionCollectionBody)
	pocStructObj.DBTree = []model.Current{}
	// fmt.Println(result)
	
	gData, err := object.GetTransactionsbyIdentifier(vars["Identifier"]).Then(func(data interface{}) interface{} {
		return data
	}).Await()

	if err != nil || gData == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		response := model.Error{Message: "Identifier for the TDP ID Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return
	}

	res := gData.([]model.TransactionCollectionBody)
	pocStructObj.Txn = res[len(res)-1].TxnHash

	for i := len(res) - 1; i >= 0; i-- {
		if res[i].TxnType == "2" {
			url := "http://localhost:3001/api/v2/dataPackets/raw?id=" + res[i].TdpId
			// url := constants.TracifiedBackend + constants.RawTDP + res[i].TdpId

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
				strreq,_:=json.Marshal(req)
				logger.LogWriter("Request  :"+string(strreq),constants.INFO)
				body, _ := ioutil.ReadAll(resq.Body)
				var raw map[string]interface{}
				json.Unmarshal(body, &raw)

				h := sha256.New()
				base64 := raw["data"]

				h.Write([]byte(fmt.Sprintf("%s", base64) + res[i].Identifier))
				
				logger.LogWriter(h.Sum(nil), constants.INFO)
				logger.LogWriter(res[i].DataHash,constants.INFO)

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

	// // }
	// display := &interpreter.AbstractPOC{POCStruct: pocStructObj}
	// response = display.InterpretPOC()

	// fmt.Println(response.RetrievePOC.Error.Message)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	// w.WriteHeader(http.StatusBadRequest)

	// result := apiModel.PoeSuccess{Message: "response.RetrievePOC.Error.Message", TxNHash: "response.RetrievePOC.Txn"}

	result := apiModel.PocSuccess{Chain: pocStructObj.DBTree}
	strresult,_:=json.Marshal(result)
	logger.LogWriter("POC success result :"+string(strresult),constants.INFO)
	//fmt.Println(response.RetrievePOC.Error.Message)
	json.NewEncoder(w).Encode(result)
	// 		return



	// return data

	// }).Catch(func(error error) error {
	// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// 	w.WriteHeader(http.StatusOK)
	// 	response := model.Error{Message: "TDP ID Not Found in Gateway DataStore"}
	// 	json.NewEncoder(w).Encode(response)
	// 	return error

	// })
	// p.Await()

	return
}
