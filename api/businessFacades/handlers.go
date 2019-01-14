package businessFacades

import (
	"strings"
	"crypto/sha256"
	"github.com/dileepaj/tracified-gateway/constants"
	"io/ioutil"
	"github.com/dileepaj/tracified-gateway/dao"
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

func GatewayRetriever(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// var response model.POC

	object := dao.Connection{}
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
				if res[i].TxnType=="2"{
					// url := "http://localhost:3001/api/v1/dataPackets/raw?id=" + res[i].TdpId
					url := constants.TracifiedBackend + constants.RawTDP + res[i].TdpId

					bearer := "Bearer " + "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb21wYW55IjoiVGVzdCAiLCJ1c2VybmFtZSI6Imhwa2F2aW5kQGdtYWlsLmNvbSIsImxvY2FsZSI6IlNyaSBMYW5rYSIsInBlcm1pc3Npb25zIjp7IjAiOlsiMTAiLCI3IiwiOCIsIjkiXSwiMDAyMDgiOlsiMSJdfSwidHlwZSI6IkFkbWluIiwidGVuYW50SUQiOiI0OTk4NDZkMC0yZDlhLTExZTgtODhmMy0wMzEyMmJkNDA1ZTEiLCJhdXRoX3RpbWUiOjE1NDIyNzI4ODYsIm5hbWUiOiJTYWFyYWtldGhhIHRlc3QgYWNjb3VudCAgIiwic3RhZ2VzIjpbIjAwMjAxIiwiMDAyMDIiLCIwMDIwMyIsIjAwMjAzIiwiMDAyMDQiLCIwMDIwNSIsIjAwMjA2IiwiMDAyMDciLCIwMDIwOCIsIjAwMjA5Il0sInBob25lX251bWJlciI6Iis5NDc3OTI5OTU5MCIsImVtYWlsIjoiaHBrYXZpbmRAZ21haWwuY29tIiwiYWRkcmVzcyI6eyJmb3JtYXR0ZWQiOiI5OXggdGVjaCJ9LCJkb21haW4iOiJEYWlyeSIsImRpc3BsYXlJbWFnZSI6Imh0dHBzOi8vdHJhY2lmaWVkLXByb2ZpbGUtaW1hZ2VzLnMzLmFwLXNvdXRoLTEuYW1hem9uYXdzLmNvbS9ocGthdmluZCU0MGdtYWlsLmNvbTE2Y2Q4OTYwLWU3ZjYtMTFlOC1iNzhlLTJkODAyZDQ2ZjlhNi5qcGVnIiwiaWF0IjoxNTQyMjcyODg1LCJleHAiOjE5OTI0NDU2ODV9.zLuscboIwwEmxB2-YLOiNb2NhxTBKkhKLZwM9Qrahtk"
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
						// fmt.Println(req)
						body, _ := ioutil.ReadAll(resq.Body)
						var raw map[string]interface{}
						json.Unmarshal(body, &raw)
	
						h := sha256.New()
						base64 := raw["data"]
						// fmt.Println(base64)
	
						h.Write([]byte(fmt.Sprintf("%s", base64)+result.Identifier))
						// fmt.Printf("%x", h.Sum(nil))
	
						DataStoreTXN := model.Current{
							TType:      res[i].TxnType,
							TXNID:      res[i].TxnHash,
							Identifier: res[i].Identifier,
							DataHash:   strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))}
	
						pocStructObj.DBTree = append(pocStructObj.DBTree, DataStoreTXN)
					}
				}else{
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
			// fmt.Println(result)
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


	return

}

func GatewayRetrieverWithIdentifier(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// var response model.POC

	object := dao.Connection{}
	var pocStructObj apiModel.POCStruct

	// p := object.GetTransactionForTdpId(vars["Txn"])
	// p.Then(func(data interface{}) interface{} {

	// 	result := data.(model.TransactionCollectionBody)
		pocStructObj.DBTree = []model.Current{}
		// fmt.Println(result)
		g := object.GetTransactionsbyIdentifier(vars["Identifier"])
		g.Then(func(data interface{}) interface{} {
			res := data.([]model.TransactionCollectionBody)
			pocStructObj.Txn = res[len(res)-1].TxnHash

			for i := len(res) - 1; i >= 0; i-- {
				if res[i].TxnType=="2"{
					// url := "http://localhost:3001/api/v1/dataPackets/raw?id=" + res[i].TdpId
					url := constants.TracifiedBackend + constants.RawTDP + res[i].TdpId

					bearer := "Bearer " + "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb21wYW55IjoiVGVzdCAiLCJ1c2VybmFtZSI6Imhwa2F2aW5kQGdtYWlsLmNvbSIsImxvY2FsZSI6IlNyaSBMYW5rYSIsInBlcm1pc3Npb25zIjp7IjAiOlsiMTAiLCI3IiwiOCIsIjkiXSwiMDAyMDgiOlsiMSJdfSwidHlwZSI6IkFkbWluIiwidGVuYW50SUQiOiI0OTk4NDZkMC0yZDlhLTExZTgtODhmMy0wMzEyMmJkNDA1ZTEiLCJhdXRoX3RpbWUiOjE1NDIyNzI4ODYsIm5hbWUiOiJTYWFyYWtldGhhIHRlc3QgYWNjb3VudCAgIiwic3RhZ2VzIjpbIjAwMjAxIiwiMDAyMDIiLCIwMDIwMyIsIjAwMjAzIiwiMDAyMDQiLCIwMDIwNSIsIjAwMjA2IiwiMDAyMDciLCIwMDIwOCIsIjAwMjA5Il0sInBob25lX251bWJlciI6Iis5NDc3OTI5OTU5MCIsImVtYWlsIjoiaHBrYXZpbmRAZ21haWwuY29tIiwiYWRkcmVzcyI6eyJmb3JtYXR0ZWQiOiI5OXggdGVjaCJ9LCJkb21haW4iOiJEYWlyeSIsImRpc3BsYXlJbWFnZSI6Imh0dHBzOi8vdHJhY2lmaWVkLXByb2ZpbGUtaW1hZ2VzLnMzLmFwLXNvdXRoLTEuYW1hem9uYXdzLmNvbS9ocGthdmluZCU0MGdtYWlsLmNvbTE2Y2Q4OTYwLWU3ZjYtMTFlOC1iNzhlLTJkODAyZDQ2ZjlhNi5qcGVnIiwiaWF0IjoxNTQyMjcyODg1LCJleHAiOjE5OTI0NDU2ODV9.zLuscboIwwEmxB2-YLOiNb2NhxTBKkhKLZwM9Qrahtk"
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
						// fmt.Println(req)
						body, _ := ioutil.ReadAll(resq.Body)
						var raw map[string]interface{}
						json.Unmarshal(body, &raw)
	
						h := sha256.New()
						base64 := raw["data"]
						// fmt.Println(base64)
	
						h.Write([]byte(fmt.Sprintf("%s", base64)+res[i].Identifier))
						// fmt.Printf("%x", h.Sum(nil))
	
						DataStoreTXN := model.Current{
							TType:      res[i].TxnType,
							TXNID:      res[i].TxnHash,
							Identifier: res[i].Identifier,
							DataHash:   strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))}
	
						pocStructObj.DBTree = append(pocStructObj.DBTree, DataStoreTXN)
					}
				}else{
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
			// fmt.Println(result)
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
		})
		g.Await()

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


