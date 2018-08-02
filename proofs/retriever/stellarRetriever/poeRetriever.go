package stellarRetriever

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tracified-gateway/model"
	"github.com/tracified-gateway/proofs/interpreter"
)

type PublicKey struct {
	Name  string
	Value string
}

type KeysResponse struct {
	Collection []PublicKey
}

type ConcretePOE struct {
	*interpreter.AbstractPOE
	Txn       string
	ProfileID string
	Hash      string
}

func (db *ConcretePOE) RetrievePOE() model.RetrievePOE {
	var bcHash string
	var response model.RetrievePOE
	var Rerr model.Error

	result, err := http.Get("https://horizon-testnet.stellar.org/transactions/" + db.Txn + "/operations")
	if err != nil {
		Rerr.Code = result.StatusCode
		Rerr.Message = "The HTTP request failed for RetrievePOE"
		response.Txn = db.Txn

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
			bcHash = Base64DecEnc("Decode", keys[0].Value)

			Rerr.Code = http.StatusOK
			Rerr.Message = "Txn Hash retrieved from the blockchain."
			response.Error = Rerr
			response.Txn = db.Txn
			response.DBHash = db.Hash
			response.BCHash = bcHash

		} else {
			Rerr.Code = http.StatusOK
			Rerr.Message = "Txn Hash does not exist in the blockchain."
			response.Txn = db.Txn

			response.Error = Rerr
			return response
		}

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
