package stellarRetriever

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"
)

type ConcretePOG struct {
	POGStruct apiModel.POGStruct
}

func (db *ConcretePOG) RetrievePOG() model.RetrievePOG {
	CurrentTxn := db.POGStruct.LastTxn
	// var response model.RetrievePOG
	response := model.RetrievePOG{CurrentTxn, "", "", model.Error{200, ""}}
	var Rerr model.Error
	var PreviousTxn string

	result, err := http.Get("https://horizon-testnet.stellar.org/transactions/" + db.POGStruct.LastTxn + "/operations")
	if err != nil {
		Rerr.Code = result.StatusCode
		Rerr.Message = "The HTTP request failed for RetrievePOG"
		response.CurTxn = CurrentTxn
		response.Message = Rerr
		return response
	} else {
		data, _ := ioutil.ReadAll(result.Body)
		if result.StatusCode == 200 {
			var raw map[string]interface{}
			json.Unmarshal(data, &raw)
			out, _ := json.Marshal(raw["_embedded"])

			var raw1 map[string]interface{}
			json.Unmarshal(out, &raw1)
			out1, _ := json.Marshal(raw1["records"])

			keysBody := out1
			keys := make([]PublicKey, 0)
			json.Unmarshal(keysBody, &keys)
			fmt.Println("keys map => ", keys)

			PreviousTxn=Base64DecEnc("Decode", keys[1].Value)
			//Retrive Current TXN DATA
			result1, err1 := http.Get("https://horizon-testnet.stellar.org/transactions/" + Base64DecEnc("Decode", keys[1].Value) + "/operations")
			if err1 != nil {
				Rerr.Code = result.StatusCode
				Rerr.Message = "The HTTP request failed for RetrievePOG"
				response.CurTxn = CurrentTxn
				response.PreTxn = Base64DecEnc("Decode", keys[0].Value)
				response.Message = Rerr
				return response
			} else {
				data, _ := ioutil.ReadAll(result1.Body)
				if result1.StatusCode == 200 {
					var raw map[string]interface{}
					json.Unmarshal(data, &raw)
					out, _ := json.Marshal(raw["_embedded"])

					var raw1 map[string]interface{}
					json.Unmarshal(out, &raw1)
					out1, _ := json.Marshal(raw1["records"])

					keysBody := out1
					keys := make([]PublicKey, 0)
					json.Unmarshal(keysBody, &keys)
					fmt.Println("keys map => ", keys)

					//Use TXN TYPE TO RETRIEVE POG VALUES
					TxnType := Base64DecEnc("Decode", keys[0].Value)

					if TxnType == "0" {
						Rerr.Code = http.StatusOK
						Rerr.Message = "Txn Hash retrieved from the blockchain."
						response.Message = Rerr
						response.CurTxn = CurrentTxn
						response.Identifier = Base64DecEnc("Decode", keys[1].Value)

						return response

					} else if PreviousTxn != "" || PreviousTxn != "0" {
						// PreviousTxn = Base64DecEnc("Decode", keys[1].Value)

						pogStruct := apiModel.POGStruct{LastTxn: PreviousTxn}
						object := ConcretePOG{POGStruct: pogStruct}

						response = object.RetrievePOG()
					} else {
						Rerr.Code = http.StatusOK
						Rerr.Message = "Genesis Transaction not found."
						response.Message = Rerr

						return response
					}
				}
			}
		} else {
			Rerr.Code = http.StatusOK
			Rerr.Message = "Txn Hash does not exist in the blockchain."
			response.CurTxn = CurrentTxn
			response.Message = Rerr

			return response
		}
	}
	return response
}
