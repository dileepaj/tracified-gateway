package stellarRetriever

import (
	"encoding/json"
	"io/ioutil"
	"main/model"
	"main/proofs/interpreter"
	"net/http"
)

type PublicKeyPOC struct {
	Name  string
	Value string
}

type KeysResponsePOC struct {
	Collection []PublicKeyPOC
}

type ConcretePOC struct {
	*interpreter.AbstractPOC
	Txn       string
	ProfileID string
	DBTree    string
}

func (db *ConcretePOC) RetrievePOC() model.RetrievePOC {
	var bcHash string
	// var Cerr error
	var response model.RetrievePOC
	var Rerr model.Error

	result, err := http.Get("https://horizon-testnet.stellar.org/transactions/" + db.Txn + "/operations")
	//https://horizon-testnet.stellar.org/transactions/e903f5ef813002295e97c0f08cf26d1fd411615e18384890395f6b0943ed83b5/operations
	if err != nil {
		Rerr.Code = result.StatusCode
		Rerr.Message = "The HTTP request failed for RetrievePOC"
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
			keys := make([]PublicKeyPOC, 0)
			json.Unmarshal(keysBody, &keys)
			// fmt.Printf("%#v", keys[0].Name)
			// fmt.Printf("%#v", keys[0].Value)
			bcHash = Base64DecEnc("Decode", keys[0].Value)

			Rerr.Code = http.StatusOK
			Rerr.Message = "Txn Hash retrieved from the blockchain."
			response.Error = Rerr
			response.Txn = db.Txn
			response.DBHash = db.DBTree
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
