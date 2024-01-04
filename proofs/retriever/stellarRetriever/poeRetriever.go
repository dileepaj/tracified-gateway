package stellarRetriever

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
)

type PublicKey struct {
	Name  string
	Value string
}

type KeysResponse struct {
	Collection []PublicKey
}

type ConcretePOE struct {
	POEStruct apiModel.POEStruct
	// Txn       string
	// ProfileID string
	// Hash      string
}

/*RetrievePOE - WORKING MODEL
@author - Azeem Ashraf
@desc - Retrieves Data TXN from stellar using the TXN ID
@params - XDR
*/
func (db *ConcretePOE) RetrievePOE() model.RetrievePOE {
	var bcHash string
	var response model.RetrievePOE
	var Rerr model.Error
	var CurrentTxn string

	//RETRIEVE GATEWAY SIGNED TXN
	result, err := http.Get(commons.GetHorizonClient().HorizonURL + "transactions/" + db.POEStruct.Txn + "/operations?limit=30")
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
			result, err := http.Get(commons.GetHorizonClient().HorizonURL + "transactions/" + CurrentTxn + "/operations?limit=30")
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
					bcHash = Base64DecEnc("Decode", keys[4].Value)
					Identifier := Base64DecEnc("Decode", keys[3].Value)

					Rerr.Code = http.StatusOK
					Rerr.Message = "Txn Hash retrieved from the blockchain."
					response.Error = Rerr
					response.TdpId = CurrentTxn
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

/**
*
*Decode or Encode a String from/to Base64
*@param typ
*@param msg
 */
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
