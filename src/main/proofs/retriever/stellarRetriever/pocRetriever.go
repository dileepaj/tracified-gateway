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
	DBTree    []model.Current
	BCTree    []model.Current
}

func (db *ConcretePOC) RetrievePOC() model.RetrievePOC {

	return db.Circle()
}

func (db *ConcretePOC) Circle() model.RetrievePOC {

	var response model.RetrievePOC
	var Rerr model.Error
	// output := make([]string, 20)
	result, err := http.Get("https://horizon-testnet.stellar.org/transactions/" + db.Txn + "/operations")
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

			var bcPreHash string
			if keys[0] != (PublicKeyPOC{}) {
				bcPreHash = Base64DecEnc("Decode", keys[0].Value)

			}
			temp := model.Current{db.Txn, Base64DecEnc("Decode", keys[1].Value)}
			db.BCTree = append(db.BCTree, temp)

			// bcPreHash := Base64DecEnc("Decode", keys[0].Value)
			profile := Base64DecEnc("Decode", keys[2].Value)
			db.ProfileID = profile
			// db.BCTree = BuildTree(bcCurrentHash, db.BCTree)
			// fmt.Println(Base64DecEnc("Decode", nil))

			for keys[0].Value != "" {
				db.Txn = bcPreHash
				result, _ = http.Get("https://horizon-testnet.stellar.org/transactions/" + db.Txn + "/operations")

				data, _ = ioutil.ReadAll(result.Body)

				if result.StatusCode == 200 {
					var raw map[string]interface{}
					json.Unmarshal(data, &raw)
					// raw["count"] = 2
					out, _ = json.Marshal(raw["_embedded"])

					var raw1 map[string]interface{}
					json.Unmarshal(out, &raw1)

					out1, _ = json.Marshal(raw1["records"])

					keysBody = out1
					keys = make([]PublicKeyPOC, 0)
					json.Unmarshal(keysBody, &keys)
					// fmt.Printf("%#v", keys[0].Name)
					// fmt.Printf("%#v", keys[0].Value)
					// bcCurrentHash := db.Txn
					if keys[0].Value != "" {
						bcPreHash = Base64DecEnc("Decode", keys[0].Value)
					}

					profile = Base64DecEnc("Decode", keys[2].Value)
					db.ProfileID = profile
					temp = model.Current{db.Txn, Base64DecEnc("Decode", keys[1].Value)}
					db.BCTree = append(db.BCTree, temp)
				}
			}
		}
		Rerr.Code = result.StatusCode
		Rerr.Message = "The Blockchain Tree Retrieved successfully"
		response.Txn = db.Txn
		response.BCHash = db.BCTree
		response.DBHash = db.DBTree
		response.Error = Rerr
		return response
	}

}
