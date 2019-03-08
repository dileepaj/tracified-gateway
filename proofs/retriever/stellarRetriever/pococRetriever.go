package stellarRetriever

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/dileepaj/tracified-gateway/model"
)

func (db *ConcretePOE) RetrievePOCOC() model.RetrievePOE {
	var bcHash string
	var response model.RetrievePOE
	var Rerr model.Error
	var CurrentTxn string

	//RETRIEVE GATEWAY SIGNED TXN
	result, err := http.Get("https://horizon-testnet.stellar.org/transactions/" + db.POEStruct.Txn + "/operations")
	if err != nil {
		Rerr.Code = result.StatusCode
		Rerr.Message = "The HTTP request failed for RetrievePOE"
		response.TdpId = db.POEStruct.Txn
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

			// Gtype:=Base64DecEnc("Decode", keys[0].Value)
			// PreviousTxn = Base64DecEnc("Decode", keys[1].Value)
			CurrentTxn = Base64DecEnc("Decode", keys[2].Value) 
			//RETRIEVE THE USER SIGNED TXN USING THE CURRENT TXN IN GATEWAY SIGNED TRANSACTION
			result, err := http.Get("https://horizon-testnet.stellar.org/transactions/" + CurrentTxn + "/operations")
			if err != nil {
				Rerr.Code = result.StatusCode
				Rerr.Message = "The HTTP request failed for RetrievePOE"
				response.TdpId = db.POEStruct.Txn
				response.Error = Rerr
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
					bcHash = Base64DecEnc("Decode", keys[2].Value)
					Identifier := Base64DecEnc("Decode", keys[1].Value)

					Rerr.Code = http.StatusOK
					Rerr.Message = "Txn Hash retrieved from the blockchain."
					response.Error = Rerr
					response.TdpId = db.POEStruct.Txn
					response.DBHash = db.POEStruct.Hash
					response.BCHash = bcHash
					response.Identifier = Identifier

				} else {
					Rerr.Code = http.StatusOK
					Rerr.Message = "Txn Hash does not exist in the blockchain."
					response.TdpId = db.POEStruct.Txn

					response.Error = Rerr
					return response
				}

			}
		}
		// fmt.Printf("%#v", keys[0].Name)
		// fmt.Printf("%#v", keys[0].Value)

	}

	return response
}