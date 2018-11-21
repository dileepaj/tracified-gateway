package businessFacades

import (
	"main/dao"

	// "main/proofs/retriever/stellarRetriever"
	"crypto/sha256"
	"net/http"

	"encoding/json"
	"fmt"
	"strings"

	// "net/http"

	"io/ioutil"

	"github.com/gorilla/mux"

	"main/api/apiModel"
	"main/model"

	// "main/proofs/builder"
	"main/proofs/interpreter"
)

// func CheckPOC(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	var response model.POC
// 	var TObj apiModel.POCOBJ
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	if r.Body == nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode("Please send a request body")
// 		return
// 	} else {
// 		err := json.NewDecoder(r.Body).Decode(&TObj)
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			json.NewEncoder(w).Encode("Error while Decoding the body")
// 			fmt.Println(err)
// 			return
// 		}

// 		fmt.Println(TObj)

// 		pocStructObj := apiModel.POCStruct{Txn: vars["Txn"], ProfileID: vars["PID"], DBTree: TObj.Chain}
// 		display := &interpreter.AbstractPOC{POCStruct: pocStructObj}
// 		response = display.InterpretPOC()

// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 		w.WriteHeader(response.RetrievePOC.Error.Code)
// 		// w.WriteHeader(http.StatusBadRequest)

// 		// result := apiModel.PoeSuccess{Message: "response.RetrievePOC.Error.Message", TxNHash: "response.RetrievePOC.Txn"}
// 		result := apiModel.PocSuccess{Message: response.RetrievePOC.Error.Message, Chain: TObj.Chain}
// 		json.NewEncoder(w).Encode(result)
// 		return
// 	}
// 	return
// }

func CheckFullPOC(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var response model.POC
	var TObj apiModel.POCOBJ
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

		pocStructObj := apiModel.POCStruct{
			Txn:       vars["Txn"],
			ProfileID: vars["PID"],
			DBTree:    TObj.Chain}
		display := &interpreter.AbstractPOC{POCStruct: pocStructObj}
		response = display.InterpretFullPOC()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(response.RetrievePOC.Error.Code)
		// w.WriteHeader(http.StatusBadRequest)

		// result := apiModel.PoeSuccess{Message: "response.RetrievePOC.Error.Message", TxNHash: "response.RetrievePOC.Txn"}
		result := apiModel.PocSuccess{
			Message: response.RetrievePOC.Error.Message,
			Chain:   TObj.Chain}
		json.NewEncoder(w).Encode(result)

		return
	}
	return
}

func CheckPOG(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var response model.POG
	pogStructObj := apiModel.POGStruct{LastTxn: vars["LastTxn"], POGTxn: vars["POGTxn"], Identifier: vars["Identifier"]}
	display := &interpreter.AbstractPOG{POGStruct: pogStructObj}
	response = display.InterpretPOG()

	// fmt.Println("response.RetrievePOG.Error.Code")
	// fmt.Println(response.RetrievePOG.Error.Code)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(response.RetrievePOG.Error.Code)
	// result := apiModel.PoeSuccess{Message: response.RetrievePOG.Error.Message, TxNHash: response.RetrievePOG.CurTxn}
	json.NewEncoder(w).Encode(response)
	return

}

// func CheckPOE(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)

// 	var response model.POE
// 	poeStructObj := apiModel.POEStruct{Txn: vars["Txn"], ProfileID: vars["PID"], Hash: vars["Hash"]}
// 	display := &interpreter.AbstractPOE{POEStruct: poeStructObj}
// 	response = display.InterpretPOE()

// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	w.WriteHeader(response.RetrievePOE.Error.Code)
// 	json.NewEncoder(w).Encode(response.RetrievePOE)
// 	return

// }

type test struct {
	Data string
}

func CheckPOE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	object := dao.Connection{}

	p := object.GetTransactionForTdpId(vars["Txn"])
	p.Then(func(data interface{}) interface{} {

		result := data.(model.TransactionCollectionBody)
		// fmt.Println(result)
		var response model.POE
		url := "http://localhost:3001/api/v1/dataPackets/raw?id=" + vars["Txn"]
		bearer := "Bearer " + "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb21wYW55IjoiVGVzdCAiLCJ1c2VybmFtZSI6Imhwa2F2aW5kQGdtYWlsLmNvbSIsImxvY2FsZSI6IlNyaSBMYW5rYSIsInBlcm1pc3Npb25zIjp7IjAiOlsiMTAiLCI3IiwiOCIsIjkiXSwiMDAyMDgiOlsiMSJdfSwidHlwZSI6IkFkbWluIiwidGVuYW50SUQiOiI0OTk4NDZkMC0yZDlhLTExZTgtODhmMy0wMzEyMmJkNDA1ZTEiLCJhdXRoX3RpbWUiOjE1NDIyNzI4ODYsIm5hbWUiOiJTYWFyYWtldGhhIHRlc3QgYWNjb3VudCAgIiwic3RhZ2VzIjpbIjAwMjAxIiwiMDAyMDIiLCIwMDIwMyIsIjAwMjAzIiwiMDAyMDQiLCIwMDIwNSIsIjAwMjA2IiwiMDAyMDciLCIwMDIwOCIsIjAwMjA5Il0sInBob25lX251bWJlciI6Iis5NDc3OTI5OTU5MCIsImVtYWlsIjoiaHBrYXZpbmRAZ21haWwuY29tIiwiYWRkcmVzcyI6eyJmb3JtYXR0ZWQiOiI5OXggdGVjaCJ9LCJkb21haW4iOiJEYWlyeSIsImRpc3BsYXlJbWFnZSI6Imh0dHBzOi8vdHJhY2lmaWVkLXByb2ZpbGUtaW1hZ2VzLnMzLmFwLXNvdXRoLTEuYW1hem9uYXdzLmNvbS9ocGthdmluZCU0MGdtYWlsLmNvbTE2Y2Q4OTYwLWU3ZjYtMTFlOC1iNzhlLTJkODAyZDQ2ZjlhNi5qcGVnIiwiaWF0IjoxNTQyMjcyODg1LCJleHAiOjE4NzI0NDU2ODV9.oiez4l8YlU0JmFl2e_kMkmAJTRe4u76Sz-mKmt-GNK0"
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
			// fmt.Println(string(raw["Data"]))
			// fmt.Println(body)

			h := sha256.New()
			lol := raw["Data"]
			fmt.Println(lol)

			h.Write([]byte(fmt.Sprintf("%s", lol)))

			fmt.Printf("%x", h.Sum(nil))

			poeStructObj := apiModel.POEStruct{Txn: result.TxnHash,
				Hash: strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))}
			display := &interpreter.AbstractPOE{POEStruct: poeStructObj}
			response = display.InterpretPOE()

			w.WriteHeader(response.RetrievePOE.Error.Code)
			json.NewEncoder(w).Encode(response.RetrievePOE)

		}

		return data

	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusNotFound)
		response := model.Error{Message: "Not Found"}
		json.NewEncoder(w).Encode(response)
		fmt.Println(response)
		return error

	})
	p.Await()

	// return

}

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
		// fmt.Println(result)
		g := object.GetTransactionsbyIdentifier(result.Identifier)
		g.Then(func(data interface{}) interface{} {
			res := data.([]model.TransactionCollectionBody)
			pocStructObj.Txn = "e2a4af62c8184d507b4f751d294df30b644c184e99aae9b1fd226c1fa90966b4"

			for i := len(res) - 1; i >= 0; i-- {
				url := "http://localhost:3001/api/v1/dataPackets/raw?id=" + res[i].TdpID
				bearer := "Bearer " + "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb21wYW55IjoiVGVzdCAiLCJ1c2VybmFtZSI6Imhwa2F2aW5kQGdtYWlsLmNvbSIsImxvY2FsZSI6IlNyaSBMYW5rYSIsInBlcm1pc3Npb25zIjp7IjAiOlsiMTAiLCI3IiwiOCIsIjkiXSwiMDAyMDgiOlsiMSJdfSwidHlwZSI6IkFkbWluIiwidGVuYW50SUQiOiI0OTk4NDZkMC0yZDlhLTExZTgtODhmMy0wMzEyMmJkNDA1ZTEiLCJhdXRoX3RpbWUiOjE1NDIyNzI4ODYsIm5hbWUiOiJTYWFyYWtldGhhIHRlc3QgYWNjb3VudCAgIiwic3RhZ2VzIjpbIjAwMjAxIiwiMDAyMDIiLCIwMDIwMyIsIjAwMjAzIiwiMDAyMDQiLCIwMDIwNSIsIjAwMjA2IiwiMDAyMDciLCIwMDIwOCIsIjAwMjA5Il0sInBob25lX251bWJlciI6Iis5NDc3OTI5OTU5MCIsImVtYWlsIjoiaHBrYXZpbmRAZ21haWwuY29tIiwiYWRkcmVzcyI6eyJmb3JtYXR0ZWQiOiI5OXggdGVjaCJ9LCJkb21haW4iOiJEYWlyeSIsImRpc3BsYXlJbWFnZSI6Imh0dHBzOi8vdHJhY2lmaWVkLXByb2ZpbGUtaW1hZ2VzLnMzLmFwLXNvdXRoLTEuYW1hem9uYXdzLmNvbS9ocGthdmluZCU0MGdtYWlsLmNvbTE2Y2Q4OTYwLWU3ZjYtMTFlOC1iNzhlLTJkODAyZDQ2ZjlhNi5qcGVnIiwiaWF0IjoxNTQyMjcyODg1LCJleHAiOjE5OTI0NDU2ODV9.zLuscboIwwEmxB2-YLOiNb2NhxTBKkhKLZwM9Qrahtk"
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
					base64 := raw["Data"]
					// fmt.Println(base64)

					h.Write([]byte(fmt.Sprintf("%s", base64)))
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
			json.NewEncoder(w).Encode(result)
			// 		return

			return data
		}).Catch(func(error error) error {
			return error
		})
		g.Await()

		return data

	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusNotFound)
		response := model.Error{Message: "Not Found"}
		json.NewEncoder(w).Encode(response)
		// fmt.Println(response)
		return error

	})
	p.Await()

	// return

}
