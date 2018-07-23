package stellarRetriever

import (
	// "encoding/base64"
	"encoding/json"
	// "fmt"
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
	// var bcHash string
	// var Cerr error
	// var response model.RetrievePOC
	// var Rerr model.Error
	/////////////////////////////////
	// result, err := http.Get("https://horizon-testnet.stellar.org/transactions/" + db.Txn + "/operations")
	// //https://horizon-testnet.stellar.org/transactions/e903f5ef813002295e97c0f08cf26d1fd411615e18384890395f6b0943ed83b5/operations
	// if err != nil {
	// 	Rerr.Code = result.StatusCode
	// 	Rerr.Message = "The HTTP request failed for RetrievePOC"
	// 	response.Txn = db.Txn

	// 	response.Error = Rerr
	// 	return response

	// } else {
	// 	data, _ := ioutil.ReadAll(result.Body)

	// 	if result.StatusCode == 200 {
	// 		var raw map[string]interface{}
	// 		json.Unmarshal(data, &raw)
	// 		// raw["count"] = 2
	// 		out, _ := json.Marshal(raw["_embedded"])

	// 		var raw1 map[string]interface{}
	// 		json.Unmarshal(out, &raw1)

	// 		out1, _ := json.Marshal(raw1["records"])

	// 		keysBody := out1
	// 		keys := make([]PublicKeyPOC, 0)
	// 		json.Unmarshal(keysBody, &keys)
	// 		// fmt.Printf("%#v", keys[0].Name)
	// 		// fmt.Printf("%#v", keys[0].Value)
	// 		bcCurrentHash := Base64DecEnc("Decode", keys[1].Value)
	// 		bcPreHash := Base64DecEnc("Decode", keys[0].Value)
	// 		profile := Base64DecEnc("Decode", keys[2].Value)
	// 		db.Txn = bcPreHash
	// 		db.ProfileID = profile
	// 		db.BCTree = BuildTree(bcCurrentHash ,db.BCTree)

	// 		// display := &ConcretePOC{Txn: db.Txn, ProfileID: db.ProfileID, DBTree: db.BCTree, BCTree:db.BCTree}
	// 		// response = display.InterpretPOC(display)
	// 		// response=RetrievePOC()
	// 		Rerr.Code = http.StatusOK
	// 		Rerr.Message = "Txn Hash retrieved from the blockchain."
	// 		response.Error = Rerr
	// 		response.Txn = db.Txn
	// 		response.DBHash = db.DBTree
	// 		response.BCHash = db.BCTree

	// 		// if bcPreHash != "" {

	// 		// 	response = interpreter.POCInterface.RetrievePOC()
	// 		// }

	// 		response.Error = Rerr

	// 	} else {

	// 		Rerr.Code = http.StatusOK
	// 		Rerr.Message = "Txn Hash does not exist in the blockchain."
	// 		response.Txn = db.Txn

	// 		response.Error = Rerr
	// 		return response
	// 	}

	// }

	//////////////////////////

	return db.Circle()
}

// func BuildTree(bcCurrentHash string, Tree model.Node) model.Node {
// 	cur1 := model.Current{bcCurrentHash}
// 	preArr1 := []model.Node{Tree}
// 	pre1 := model.Node{preArr1, cur1}
// 	fmt.Println(Tree)
// 	fmt.Println(pre1)
// 	return pre1
// }

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
			temp:=model.Current{db.Txn,Base64DecEnc("Decode", keys[1].Value)}
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
					temp=model.Current{db.Txn,Base64DecEnc("Decode", keys[1].Value)}
					db.BCTree = append(db.BCTree, temp)
				}
			}
		}
		response.Txn = db.Txn
		response.BCHash = db.BCTree
		response.Error = Rerr
		return response
	}
	// Rerr.Code = result.StatusCode
	// Rerr.Message = "The HTTP request failed for RetrievePOC"

}
