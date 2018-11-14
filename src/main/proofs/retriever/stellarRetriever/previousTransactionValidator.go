package stellarRetriever

import (
	"encoding/json"
	"io/ioutil"
	// "main/api/apiModel"
	"main/model"

	"net/http"
)

type ConcretePrevious struct {
	POCStruct model.RetrievePrevious
	Count     int
	Err       error
}

func (db *ConcretePrevious) RetrievePrevious8Transactions(txn string) (model.RetrievePrevious, error) {

	db.Count++
	result, err := http.Get("https://horizon-testnet.stellar.org/transactions/" + txn + "/operations")
	if err != nil {
		db.Err = err
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
			keys := make([]PublicKeyPOC, 0)
			json.Unmarshal(keysBody, &keys)

			var bcPreHash string
			var transactionType string
			var temp model.Current
			temp.TXNID = txn
			if keys[0] != (PublicKeyPOC{}) {
				transactionType = Base64DecEnc("Decode", keys[0].Value)
				temp.TType = transactionType
			}
			if keys[1] != (PublicKeyPOC{}) {
				bcPreHash = Base64DecEnc("Decode", keys[1].Value)
			}

			db.POCStruct.HashList = append(db.POCStruct.HashList, temp)

			if bcPreHash != "" {
				if db.Count == 8 {
					return db.POCStruct, db.Err
				} else {
					return db.RetrievePrevious8Transactions(bcPreHash)
				}
			} else {
				return db.POCStruct, db.Err
			}

		}

	}
	return db.POCStruct, db.Err

}
