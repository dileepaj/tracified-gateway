package stellarRetriever

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"main/api/apiModel"
	"main/model"
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

		response.Error = Rerr
		return response

	} else {
		data, _ := ioutil.ReadAll(result.Body)

		if result.StatusCode == 200 {
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
			// fmt.Printf("%#v", keys[0].Name)
			// fmt.Printf("%#v", keys[0].Value)
			fmt.Println("keys map => ", keys)
			TxnType := Base64DecEnc("Decode", keys[0].Value)

			if keys[0].Name == "Transaction Type" && TxnType == "0" {
				Rerr.Code = http.StatusOK
				Rerr.Message = "Txn Hash retrieved from the blockchain."
				response.Error = Rerr
				response.CurTxn = CurrentTxn
				response.PreTxn = Base64DecEnc("Decode", keys[1].Value)
				response.Identifier = Base64DecEnc("Decode", keys[2].Value)

				return response

			} else if keys[1].Value != "" {
				PreviousTxn = Base64DecEnc("Decode", keys[1].Value)

				pogStruct := apiModel.POGStruct{LastTxn: PreviousTxn}
				object := ConcretePOG{POGStruct: pogStruct}

				response = object.RetrievePOG()
			} else {
				Rerr.Code = http.StatusOK
				Rerr.Message = "Genesis Transaction not found."
				response.Error = Rerr

				return response
			}

		} else {
			Rerr.Code = http.StatusOK
			Rerr.Message = "Txn Hash does not exist in the blockchain."
			response.CurTxn = CurrentTxn
			response.Error = Rerr

			return response
		}

	}

	return response
}
