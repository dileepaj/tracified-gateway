package stellarRetriever

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"main/proofs/interpreter"
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

func (db *ConcretePOE) RetrievePOETest() (string, string, string) {
	var bcHash string
	response, err := http.Get("https://horizon-testnet.stellar.org/transactions/" + db.Txn + "/operations")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)

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
		// fmt.Println(bcHash)

	}
	return db.Txn, bcHash, db.Hash

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
