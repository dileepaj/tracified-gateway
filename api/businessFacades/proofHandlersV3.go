package businessFacades

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/stellar/go/support/log"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/interpreter"
	"github.com/dileepaj/tracified-gateway/proofs/retriever/stellarRetriever"
	"github.com/gorilla/mux"
)

type PublicKey struct {
	Name  string
	Value string
}

type KeysResponse struct {
	Collection []PublicKey
}

type Item struct {
	ItemID   string `json:"itemID"`
	ItemName string `json:"itemName"`
}

type TdpHeader struct {
	Identifiers      []string `json:"identifiers"`
	Item             Item     `json:"item"`
	StageID          string   `json:"stageID"`
	TimeStamp        string   `json:"timeStamp"`
	WorkflowRevision string   `json:"workflowRevision"`
}

type TdpData struct {
	Data      map[string]interface{} `json:"data"`
	TdpHeader TdpHeader              `json:"header"`
}

type Identifier struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type response struct {
	Status string
}

/*CheckPOEV3 - WORKING MODEL
@author - Azeem Ashraf, Jajeththanan Sabapathipillai
@desc - Handles the Proof of Existance by retrieving the Raw Data from the Traceability Data Store
and Retrieves the TXN ID and calls POE Interpreter
Finally Returns the Response given by the POE Interpreter
@params - ResponseWriter,Request
*/
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
		log.Error("Error while GetTransactionForTdpId " + error.Error())
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "TDPID NOT FOUND IN DATASTORE"}
		json.NewEncoder(w).Encode(response)
		fmt.Println(response)
		return error
	}).Await()

	result1, err := http.Get(commons.GetHorizonClient().HorizonURL + "transactions/" + result.TxnHash + "/operations")
	fmt.Println(commons.GetHorizonClient().HorizonURL + "transactions/" + result.TxnHash + "/operations")
	fmt.Println(result1,err)
	if err != nil {
		log.Error("Error while getting transactions by txnhash " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Txn for the TXN does not exist in the Blockchain " + err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := ioutil.ReadAll(result1.Body)
	if err != nil {
		log.Error("Error while read response " + err.Error())
	}
	var raw map[string]interface{}
	err = json.Unmarshal(data, &raw)
	if err != nil {
		log.Error("Error while json.Unmarshal(data, &raw) " + err.Error())
	}

	out, err := json.Marshal(raw["_embedded"])
	if err != nil {
		log.Error("Error while json marshal _embedded " + err.Error())
	}
	var raw1 map[string]interface{}
	err = json.Unmarshal(out, &raw1)
	if err != nil {
		log.Error("Error while json.Unmarshal(out, &raw1) " + err.Error())
	}
	out1, err := json.Marshal(raw1["records"])
	if err != nil {
		log.Error("Error while json marshal records " + err.Error())
	}
	keysBody := out1
	keys := make([]PublicKey, 0)
	err = json.Unmarshal(keysBody, &keys)
	if err != nil {
		log.Error("Error while json.Unmarshal(keysBody, &keys) " + err.Error())
	}
	byteData, err := base64.StdEncoding.DecodeString(keys[2].Value)
	if err != nil {
		log.Error("Error while base64.StdEncoding.DecodeString " + err.Error())
	}
	CurrentTxn = string(byteData)
	log.Info("THE TXN OF THE USER TXN: " + CurrentTxn)

	var finalResult []model.POEResponse

	var response model.POE
	// url := "http://localhost:3001/api/v2/dataPackets/raw?id=5c9141b2618cf404ec5e105d"
	url := constants.TracifiedBackend + constants.RawTDP + vars["Txn"]
	fmt.Println(url)
	bearer := "Bearer " + constants.BackendToken

	// Create a new request using http
	req, er := http.NewRequest("GET", url, nil)
	if er != nil {
		log.Error("Error while create new request using http " + er.Error())
	}
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	resq, er := client.Do(req)
	if er != nil {
		log.Error("Error while getting response " + er.Error())
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Connection to the Traceability DataStore was interupted " + er.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	body, err := ioutil.ReadAll(resq.Body)
	if err != nil {
		log.Error("Error while ioutil.ReadAll(resq.Body) " + err.Error())
	}
	var raw2 map[string]interface{}
	json.Unmarshal(body, &raw2)

	h := sha256.New()
	lol := raw2["data"]
	encodedString := raw2["data"].(string)
	decodedString, err := base64.StdEncoding.DecodeString(encodedString)
	var tdpJson TdpData
	json.Unmarshal(decodedString, &tdpJson)
	decodedIdentifier, err := base64.StdEncoding.DecodeString(tdpJson.TdpHeader.Identifiers[0])
	var identifierObj Identifier
	json.Unmarshal(decodedIdentifier, &identifierObj)
	h.Write([]byte(fmt.Sprintf("%s", lol) + identifierObj.Id))
	poeStructObj := apiModel.POEStruct{Txn: result.TxnHash, Hash: result.DataHash}
	display := &interpreter.AbstractPOE{POEStruct: poeStructObj}
	response = display.InterpretPOE()

	w.WriteHeader(response.RetrievePOE.Error.Code)

	//var txe xdr.Transaction
	TxnHash := CurrentTxn
	PublicKey := result.PublicKey

	result2, err2 := http.Get(commons.GetHorizonClient().HorizonURL + "transactions/" + TxnHash)
	fmt.Println(commons.GetHorizonClient().HorizonURL + "transactions/" + TxnHash)
	fmt.Println(result2,err)
	if err2 != nil {
		log.Error("Error while get transactions by TxnHash " + err2.Error())
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Txn Id Not Found in Stellar Public Net " + err2.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data2, err := ioutil.ReadAll(result2.Body)
	if err != nil {
		log.Error("Error while ioutil.ReadAll(result2.Body) " + err.Error())
	}
	if result2.StatusCode != 200 {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Txn Id Not Found in Stellar Public Net"}
		json.NewEncoder(w).Encode(response)
		return
	}
	var raw3 map[string]interface{}
	err = json.Unmarshal(data2, &raw3)
	if err != nil {
		log.Error("Error while json.Unmarshal(data2, &raw3) " + err.Error())
	}

	fmt.Println(raw3["envelope_xdr"])
	fmt.Println("HAHAHAHAAHAHAH")
	timestamp := fmt.Sprintf("%s", raw3["created_at"])
	ledger := fmt.Sprintf("%.0f", raw3["ledger"])
	feePaid := fmt.Sprintf("%s", raw3["fee_charged"])
	//envelopeXDR := fmt.Sprintf("%v", raw3["envelope_xdr"])
	//errXDR := xdr.SafeUnmarshalBase64(envelopeXDR, &txe)
	//if errXDR != nil {
	//	log.Error("Error while SafeUnmarshalBase64 "+errXDR.Error())
	//}

	mapD := map[string]string{"transaction": TxnHash}
	mapB, err := json.Marshal(mapD)
	if err != nil {
		log.Error("Error while json.Marshal(mapD) " + err.Error())
	}
	fmt.Println(string(mapB))

	encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
	text := encoded
	temp := model.POEResponse{
		Txnhash: TxnHash,
		Url:            commons.GetHorizonClient().HorizonURL + "transactions/" + TxnHash + "/operations",
		LabUrl:			commons.GetStellarLaboratoryClient() + "/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
		text + "%3D%3D&network=" + commons.GetHorizonClientNetworkName(),
		Identifier:     result.Identifier,
		SequenceNo:     result.SequenceNo,
		TxnType:        "tdp",
		Status:         response.RetrievePOE.Error.Message,
		BlockchainName: "Stellar",
		Timestamp:      timestamp,
		Ledger:         ledger,
		FeePaid:        feePaid,
		SourceAccount:  PublicKey,
		DbHash:         response.RetrievePOE.DBHash,
		BcHash:         response.RetrievePOE.BCHash}

	finalResult = append(finalResult, temp)

	json.NewEncoder(w).Encode(finalResult)

	return

}

/*CheckPOCV3 - Needs to be Tested
@author - Azeem Ashraf
@desc - Handles the Proof of Continuity by using the TXN ID in the PARAMS and
Creates the Complete tree using the gateway DB
and calls POC Interpreter sending the tree in as a Param
Finally Returns the Response given by the POC Interpreter
@params - ResponseWriter,Request
*/
func CheckPOCV3(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	object := dao.Connection{}
	var response model.POC
	var pocStructObj apiModel.POCStruct

	//checks the gateway DB for a TXN with the TdpID in the parameter
	pData, errAsnc := object.GetTransactionByTxnhash(vars["Txn"]).Then(func(data interface{}) interface{} {
		return data
	}).Await()

	if errAsnc != nil || pData == nil {
		log.Error("Error while GetTransactionByTxnhash " + errAsnc.Error())
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		response := model.Error{Message: "Txn Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return
	} else {
		result := pData.(model.TransactionCollectionBody)
		// fmt.Println(result)
		pocStructObj.Txn = result.TxnHash

		pocStructObj.DBTree = []model.Current{}
		// fmt.Println(result)

		//using the identifier retrieved from the gateway DB for the particular TdpID

		//retrieve all the transactions.
		// g := object.GetTransactionsbyIdentifier(result.Identifier)
		// g.Then(func(data interface{}) interface{} {
		// 	res := data.([]model.TransactionCollectionBody)
		// 	pocStructObj.Txn = res[len(res)-1].TxnHash

		// 	for i := len(res) - 1; i >= 0; i-- {

		// 		result1, err1 := http.Get("https://horizon.stellar.orgtransactions/" + res[i].TxnHash + "/operations")
		// 		if err1 != nil {
		// 			// Rerr.Code = result1.StatusCode
		// 			// Rerr.Message = "The HTTP request failed for RetrievePOC"
		// 			// response.Txn = db.POCStruct.Txn
		// 			// response.Error = Rerr
		// 			// return response
		// 		}

		// 		data, _ := ioutil.ReadAll(result1.Body)
		// 		var raw map[string]interface{}
		// 		json.Unmarshal(data, &raw)
		// 		out, _ := json.Marshal(raw["_embedded"])

		// 		var raw1 map[string]interface{}
		// 		json.Unmarshal(out, &raw1)

		// 		out1, _ := json.Marshal(raw1["records"])

		// 		keysBody := out1
		// 		keys := make([]PublicKeyPOC, 0)
		// 		json.Unmarshal(keysBody, &keys)

		// 		Current := Base64DecEnc("Decode", keys[2].Value)
		// 		// GatewayTXNType := Base64DecEnc("Decode", keys[0].Value)

		// 		if res[i].TxnType == "2" {
		// 			// url := "http://localhost:3001/api/v1/dataPackets/raw?id=" + res[i].TdpId
		// 			url := constants.TracifiedBackend + constants.RawTDP + res[i].TdpId

		// 			bearer := "Bearer " + constants.BackendToken
		// 			// Create a new request using http
		// 			req, er := http.NewRequest("GET", url, nil)

		// 			req.Header.Add("Authorization", bearer)
		// 			client := &http.Client{}
		// 			resq, er := client.Do(req)

		// 			if er != nil {

		// 				w.WriteHeader(http.StatusBadRequest)
		// 				response := model.Error{Message: "Connection to the DataStore was interupted"}
		// 				json.NewEncoder(w).Encode(response)
		// 			} else {
		// 				// fmt.Println(req)
		// 				body, _ := ioutil.ReadAll(resq.Body)
		// 				var raw map[string]interface{}
		// 				json.Unmarshal(body, &raw)

		// 				h := sha256.New()
		// 				base64 := raw["data"]
		// 				// fmt.Println(base64)

		// 				h.Write([]byte(fmt.Sprintf("%s", base64) + result.Identifier))
		// 				// fmt.Printf("%x", h.Sum(nil))

		// 				DataStoreTXN := model.Current{
		// 					TType:      res[i].TxnType,
		// 					TXNID:      Current,
		// 					Identifier: res[i].Identifier,
		// 					DataHash:   strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))}

		// 				pocStructObj.DBTree = append(pocStructObj.DBTree, DataStoreTXN)
		// 			}
		// 		} else {

		// 			//this should be wear all future TXN types and their fields should be assigned

		// 			//when retrieving from the gateway DB
		// 			DataStoreTXN := model.Current{
		// 				TType:      res[i].TxnType,
		// 				TXNID:      Current,
		// 				Identifier: res[i].Identifier,
		// 			}
		// 			pocStructObj.DBTree = append(pocStructObj.DBTree, DataStoreTXN)
		// 		}
		// 	}

		// pocStructObj = apiModel.POCStruct{

		// }
		display := &interpreter.AbstractPOC{POCStruct: pocStructObj}
		response = display.InterpretFullPOC()
		// fmt.Println(response.RetrievePOC.BCHash)
		var POCTree []model.POCResponse
		for _, tree := range response.RetrievePOC.BCHash {
			sr := stellarRetriever.ConcreteStellarTransaction{Txnhash: tree.TXNID}
			txn, err := sr.RetrieveTransaction()
			if err != nil {
				// w.WriteHeader(http.StatusBadRequest)
				// response := model.Error{Message: "Txn Id Not Found in Stellar Public Net"}
				// json.NewEncoder(w).Encode(response)
				// return nil
			}
			timestamp := fmt.Sprintf("%s", txn.CreatedAt)
			ledger := strconv.Itoa(txn.Ledger)
			feePaid := fmt.Sprintf("%s", txn.FeeCharged)
			if err != nil {
				log.Error("Error while json.Marshal(mapD) " + err.Error())
			}
			// fmt.Println(string(mapB))
			// trans := transaction{transaction:TxnHash}
			// s := fmt.Sprintf("%v", trans)

			//GET THE USER SIGNED GENESIS TXN
			oprn, err := sr.RetrieveOperations()
			if err != nil {

			}
			Type, _ := base64.StdEncoding.DecodeString(oprn.Embedded.Records[0].Value)
			Identifier, _ := base64.RawStdEncoding.DecodeString(oprn.Embedded.Records[1].Value)
			SourceAccount := txn.SourceAccount
			sequenceNo, err := strconv.Atoi(txn.SourceAccountSequence)
			temp := model.POCResponse{
				Status:         "success",
				BlockchainName: "Stellar",
				Txnhash:        tree.TXNID,
				TxnType:        GetTransactiontype(string(Type)),
				AvailableProof: GetProofName(string(Type)),
				Url:            txn.Links.Self.Href + "/operations",
				DataHash:       tree.DataHash,
				Timestamp:      timestamp,
				Ledger:         ledger,
				FeePaid:        feePaid,
				Identifier:     string(Identifier),
				SourceAccount:  SourceAccount,
				SequenceNo:     int64(sequenceNo),
			}

			POCTree = append(POCTree, temp)
		}

		// fmt.Println(response.RetrievePOC.Error.Message)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(200)
		// w.WriteHeader(http.StatusBadRequest)
		// res := apiModel.PocSuccess{Message: "Tree retrieved", Chain: POCTree}

		// res := apiModel.PocSuccess{Message: response.RetrievePOC.Error.Message, Chain: response.RetrievePOC.BCHash}
		// fmt.Println(result)
		fmt.Println(response.RetrievePOC.Error.Message)

		json.NewEncoder(w).Encode(POCTree)
		//return

		// return data
		// }).Catch(func(error error) error {
		// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// 	w.WriteHeader(http.StatusOK)
		// 	response := model.Error{Message: "Identifier for the TDP ID Not Found in Gateway DataStore"}
		// 	json.NewEncoder(w).Encode(response)
		// 	return error
		// })
		// g.Await()
		return
	}
	// return
}

/*CheckFullPOCV3 - Needs to be Tested
@author - Azeem Ashraf
@desc - Handles the Full Proof of Continuity by using the TXN ID in the PARAMS and
Creates the Complete tree using the gateway DB
and calls FullPOC Interpreter sending the tree in as a Param
Finally Returns the Response given by the FullPOC Interpreter
@params - ResponseWriter,Request
*/
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
		}).Await()

		return data

	}).Catch(func(error error) error {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Identifier Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error

	}).Await()

	// return

}

/*CheckPOGV3 - WORKING MODEL
@author - Azeem Ashraf, Jajeththanan Sabapathipillai
@desc - Handles the Proof of Genesis  Retrieves the TXN ID and calls POG Interpreter
Creates the Complete tree using the gateway DB
and calls FullPOC Interpreter sending the tree in as a Param
Finally Returns the Response given by the FullPOC Interpreter
@params - ResponseWriter,Request
*/
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

			//First TXN SIGNED BY GATEWAY IS USED TO REQUEST THE USER's GENESIS
			result1, err := http.Get("https://horizon.stellar.orgtransactions/" + FirstTxnGateway.TxnHash + "/operations")
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

/*CheckPOGV3Rewrite - WORKING MODEL
@author - Azeem Ashraf, Jajeththanan Sabapathipillai
@desc - Handles the Proof of Genesis  Retrieves the TXN ID and calls POG Interpreter
Creates the Complete tree using the gateway DB
and calls FullPOC Interpreter sending the tree in as a Param
Finally Returns the Response given by the FullPOC Interpreter
@params - ResponseWriter,Request
*/
func CheckPOGV3Rewrite(writer http.ResponseWriter, r *http.Request) {
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// var txe TransactionEnvelope

	vars := mux.Vars(r)
	var result []model.POGResponse
	var res model.TransactionCollectionBody
	object := dao.Connection{}
	p := object.GetTransactionByTxnhash(vars["Txn"])
	resData, err := p.Then(func(data interface{}) interface{} {
		return data
	}).Await()

	if err != nil || resData == nil {
		log.Error("Error while GetTransactionForTdpId " + err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "TDPID NOT FOUND IN DATASTORE"}
		json.NewEncoder(writer).Encode(response)
		fmt.Println(response)
		return
	}

	res = resData.(model.TransactionCollectionBody)

	TxnHash := res.TxnHash
	PublicKey := res.PublicKey
	result1, err := http.Get(commons.GetHorizonClient().HorizonURL + "transactions/" + TxnHash)
	if err != nil {
		log.Error("Error while getting transactions by TxnHash " + err.Error())
		response := model.Error{Message: "Txn Id Not Found in Stellar Public Net " + err.Error()}
		json.NewEncoder(writer).Encode(response)
		return
	}

	data, err := ioutil.ReadAll(result1.Body)
	if err != nil {
		log.Error("Error while ReadAll " + err.Error())
	}

	if result1.StatusCode != 200 {
		writer.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Txn Id Not Found in Stellar Public Net"}
		json.NewEncoder(writer).Encode(response)
		return
	}
	var raw map[string]interface{}
	json.Unmarshal(data, &raw)

	timestamp := fmt.Sprintf("%s", raw["created_at"])
	ledger := fmt.Sprintf("%.0f", raw["ledger"])
	feePaid := fmt.Sprintf("%s", raw["fee_charged"])

	// envelopeXDR := fmt.Sprintf("%s", raw["envelope_xdr"])
	// errXDR := xdr.SafeUnmarshalBase64(envelopeXDR, &txe)
	// if errXDR != nil {
	// 	log.Error("Error while SafeUnmarshalBase64 " + errXDR.Error())
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	response := model.Error{Message: "Error while SafeUnmarshalBase64 " + errXDR.Error()}
	// 	json.NewEncoder(w).Encode(response)
	// 	return nil
	// }

	result2, _ := http.Get(commons.GetHorizonClient().HorizonURL + "transactions/" + TxnHash + "/operations")
	data2, _ := ioutil.ReadAll(result2.Body)
	var raw2 map[string]interface{}
	var raw4 map[string]interface{}
	var raw5 map[string]interface{}
	json.Unmarshal(data2, &raw2)
	out1, _ := json.Marshal(raw2["_embedded"])
	json.Unmarshal(out1, &raw2)

	var raw3 []interface{}
	out3, _ := json.Marshal(raw2["records"])
	json.Unmarshal(out3, &raw3)

	out5, _ := json.Marshal(raw3[0])
	out4, _ := json.Marshal(raw3[1])
	out6, _ := json.Marshal(raw3[2])

	json.Unmarshal(out5, &raw4)
	json.Unmarshal(out4, &raw2)
	json.Unmarshal(out6, &raw5)

	//GET THE USER SIGNED GENESIS TXN
	Type := strings.TrimLeft(fmt.Sprintf("%s", raw4["value"]), "&")
	Previous := strings.TrimLeft(fmt.Sprintf("%s", raw2["value"]), "&")
	CurrentTxn := strings.TrimLeft(fmt.Sprintf("%s", raw5["value"]), "&")

	TypeDecoded, _ := base64.StdEncoding.DecodeString(Type)
	PreviousDecoded, _ := base64.StdEncoding.DecodeString(Previous)
	CurrentTxnDecoded, _ := base64.StdEncoding.DecodeString(CurrentTxn)

	if string(TypeDecoded) != "G0" || string(PreviousDecoded) != "" {
		writer.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "This Transaction is not a Genesis Txn"}
		json.NewEncoder(writer).Encode(response)
		return
	}

	result3, err4 := http.Get(commons.GetHorizonClient().HorizonURL + "transactions/" + string(CurrentTxnDecoded) + "/operations")
	if err4 != nil {
		log.Error("Error while getting the current transaction by TxnHash " + err.Error())
		response := model.Error{Message: "Current Txn Id Not Found in Stellar Public Net " + err.Error()}
		json.NewEncoder(writer).Encode(response)
		return
	}
	data5, _ := ioutil.ReadAll(result3.Body)
	var raw6 map[string]interface{}
	json.Unmarshal(data5, &raw6)

	var raw7 []interface{}
	json.Unmarshal(data2, &raw2)
	out8, _ := json.Marshal(raw6["_embedded"])
	json.Unmarshal(out8, &raw6)
	out8, _ = json.Marshal(raw6["records"])
	json.Unmarshal(out8, &raw7)

	ProductName := "N/A"
	ProductId := "N/A"

	if len(raw7) > 3 {
		out9, _ := json.Marshal(raw7[2])
		out10, _ := json.Marshal(raw7[3])

		var raw20 map[string]interface{}
		var raw40 map[string]interface{}

		json.Unmarshal(out9, &raw20)
		json.Unmarshal(out10, &raw40)

		ProductNameEncoded := strings.TrimLeft(fmt.Sprintf("%s", raw20["value"]), "&")
		ProductIdEncoded := strings.TrimLeft(fmt.Sprintf("%s", raw40["value"]), "&")

		ProductNameDecoded, _ := base64.StdEncoding.DecodeString(ProductNameEncoded)
		ProductIdDecoded, _ := base64.StdEncoding.DecodeString(ProductIdEncoded)

		ProductName = string(ProductNameDecoded)
		ProductId = string(ProductIdDecoded)
	}

	mapD := map[string]string{"transaction": TxnHash}
	mapB, err := json.Marshal(mapD)
	if err != nil {
		log.Error("Error while json.Marshal(mapD) " + err.Error())
	}
	fmt.Println(string(mapB))

	encoded := base64.StdEncoding.EncodeToString([]byte(string(mapB)))
	text := encoded

	temp := model.POGResponse{
		Txnhash:        TxnHash,
		Url:            commons.GetHorizonClient().HorizonURL + "transactions/" + string(CurrentTxnDecoded) + "/operations",
		LabUrl:			commons.GetStellarLaboratoryClient() + "/laboratory/#explorer?resource=operations&endpoint=for_transaction&values=" +
		text + "%3D%3D&network=" + commons.GetHorizonClientNetworkName(),
		Identifier:     res.Identifier,
		SequenceNo:     res.SequenceNo,
		TxnType:        "genesis",
		Status:         "Success",
		BlockchainName: "Stellar",
		Timestamp:      timestamp,
		Ledger:         ledger,
		FeePaid:        feePaid,
		SourceAccount:  PublicKey,
		ProductName:    ProductName,
		ProductId:      ProductId,
	}
	result = append(result, temp)
	// writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(result)
	return

	// log.Error("Error while GetTransactionByTxnhash " + error.Error())
	// // writer.WriteHeader(http.StatusBadRequest)
	// response := model.Error{Message: "Txn Id Not Found in Gateway DataStore"}
	// json.NewEncoder(writer).Encode(response)
	// return

}

/*CheckPOCOCV3 - WORKING MODEL
@author - Azeem Ashraf
@desc - Handles the Proof of Change of Custody by using the last COC TXN ID as Param,
retrieves the COC object from the Gateway DB
and calls POCOC Interpreter COC Transaction in the Stellar Transaction Format
Finally Returns the Response given by the POCOC Interpreter
@params - ResponseWriter,Request
*/
func CheckPOCOCV3(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//var txe xdr.Transaction
	var COC model.COCCollectionBody
	var COCAvailable bool
	var txe interpreter.XDR
	vars := mux.Vars(r)

	result1, err := http.Get(commons.GetHorizonClient().HorizonURL + "transactions/" + vars["TxnId"] + "/operations")
	if err != nil {
		log.Error("Error while getting transactions by txnhash " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Txn for the TXN does not exist in the Blockchain " + err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	data, err := ioutil.ReadAll(result1.Body)
	if err != nil {
		log.Error("Error while read response " + err.Error())
	}
	var raw map[string]interface{}
	err = json.Unmarshal(data, &raw)
	if err != nil {
		log.Error("Error while json.Unmarshal(data, &raw) " + err.Error())
	}

	out, err := json.Marshal(raw["_embedded"])
	if err != nil {
		log.Error("Error while json marshal _embedded " + err.Error())
	}
	var raw1 map[string]interface{}
	err = json.Unmarshal(out, &raw1)
	if err != nil {
		log.Error("Error while json.Unmarshal(out, &raw1) " + err.Error())
	}
	out1, err := json.Marshal(raw1["records"])
	if err != nil {
		log.Error("Error while json marshal records " + err.Error())
	}
	keysBody := out1
	keys := make([]PublicKeyPOCOC, 0)
	err = json.Unmarshal(keysBody, &keys)
	if err != nil {
		log.Error("Error while json.Unmarshal(keysBody, &keys) " + err.Error())
	}
	acceptTxn_byteData, err := base64.StdEncoding.DecodeString(keys[2].Value)
	if err != nil {
		log.Error("Error while base64.StdEncoding.DecodeString " + err.Error())
	}
	acceptTxn := string(acceptTxn_byteData)
	log.Info("acceptTxn: " + acceptTxn)

	object := dao.Connection{}
	_, err = object.GetCOCbyAcceptTxn(acceptTxn).Then(func(data interface{}) interface{} {
		COCAvailable = true
		COC = data.(model.COCCollectionBody)
		fmt.Println(COC)
		return data
	}).Await()

	if err != nil {
		log.Error("Error while GetCOCbyTxn " + err.Error())
		COCAvailable = false
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "COCTXN NOT FOUND IN GATEWAY DATASTORE " + err.Error()}
		json.NewEncoder(w).Encode(response)
		fmt.Println(response)
	}

	if COC.Status == model.Rejected.String() || COC.Status == model.Expired.String() || COC.Status == model.Pending.String() {

		w.WriteHeader(http.StatusBadRequest)
		COCAvailable = false
		response := response{Status: COC.Status}
		json.NewEncoder(w).Encode(response)
	}

	if COCAvailable {
		/*err := xdr.SafeUnmarshalBase64(COC.AcceptXdr, &txe)
		if err != nil {
			log.Error("Error SafeUnmarshalBase64 COC " + err.Error())
		}
		proofhash := strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[2].Body.ManageDataOp.DataValue), "&")
		fmt.Println(proofhash)
		COCStatus := COC.Status
		display := &interpreter.AbstractPOCOC{Txn: vars["TxnId"], DBCOC: txe, XDR: COC.AcceptXdr, ProofHash: proofhash, COCStatus: COCStatus, SequenceNo: COC.SequenceNo}
		display.InterpretPOCOC(w, r)
		*/
		result1, err := http.Get(commons.GetHorizonClient().HorizonURL + "transactions/" + acceptTxn + "/operations")
		if err != nil {
			log.Error("Error while getting transactions by txnhash " + err.Error())
			w.WriteHeader(http.StatusBadRequest)
			response := model.Error{Message: "Txn for the TXN does not exist in the Blockchain " + err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}
		data, err := ioutil.ReadAll(result1.Body)
		if err != nil {
			log.Error("Error while read response " + err.Error())
		}
		var raw map[string]interface{}
		err = json.Unmarshal(data, &raw)
		if err != nil {
			log.Error("Error while json.Unmarshal(data, &raw) " + err.Error())
		}

		out, err := json.Marshal(raw["_embedded"])
		if err != nil {
			log.Error("Error while json marshal _embedded " + err.Error())
		}
		var raw1 map[string]interface{}
		err = json.Unmarshal(out, &raw1)
		if err != nil {
			log.Error("Error while json.Unmarshal(out, &raw1) " + err.Error())
		}
		out1, err := json.Marshal(raw1["records"])
		if err != nil {
			log.Error("Error while json marshal records " + err.Error())
		}
		keysBody := out1
		keys := make([]PublicKeyPOCOC, 0)
		err = json.Unmarshal(keysBody, &keys)
		if err != nil {
			log.Error("Error while json.Unmarshal(keysBody, &keys) " + err.Error())
		}
		ProofHash_byteData, err := base64.StdEncoding.DecodeString(keys[2].Value)
		if err != nil {
			log.Error("Error while base64.StdEncoding.DecodeString " + err.Error())
		}
		Proofhash := string(ProofHash_byteData)
		log.Info("ProofHash: " + Proofhash)

		txe.SourceAccount = string(keys[1].Source_account)
		log.Info("Source Account: " + txe.SourceAccount)

		txe.AssetCode = string(keys[3].Asset_code)
		log.Info("Asset Code: " + txe.AssetCode)

		txe.AssetAmount, err = strconv.ParseFloat(string(keys[3].Amount), 64)
		log.Info("Asset Amount: " + fmt.Sprintf("%f", txe.AssetAmount))

		Identifier_byteData, err := base64.StdEncoding.DecodeString(keys[1].Value)
		if err != nil {
			log.Error("Error while base64.StdEncoding.DecodeString " + err.Error())
		}

		txe.Identifier = string(Identifier_byteData)
		log.Info("Identifier: " + txe.Identifier)

		txe.Destination = string(keys[3].To)
		log.Info("Asset Code: " + txe.AssetCode)

		COCStatus := COC.Status
		display := &interpreter.AbstractPOCOCNew{Txn: acceptTxn, DBCOC: txe, XDR: COC.AcceptXdr, ProofHash: Proofhash, COCStatus: COCStatus, SequenceNo: COC.SequenceNo}
		display.InterpretPOCOCNew(w, r)

	}
	return
}

type PublicKeyPOC struct {
	Name  string
	Value string
}

type KeysResponsePOC struct {
	Collection []PublicKeyPOC
}

func Base64DecEnc(typ string, msg string) string {
	var text string

	if typ == "Encode" {
		encoded := base64.StdEncoding.EncodeToString([]byte(msg))
		text = (string(encoded))

	} else if typ == "Decode" {
		decoded, err := base64.StdEncoding.DecodeString(msg)
		if err != nil {
			fmt.Println("decode error:", err)
		} else {
			text = string(decoded)
		}

	} else {
		text = "Typ has to be either Encode or Decode!"
	}

	return text
}

type PublicKeyPOCOC struct {
	Name           string
	Value          string
	Source_account string
	Asset_code     string
	Amount         string
	To             string
	From           string
}
