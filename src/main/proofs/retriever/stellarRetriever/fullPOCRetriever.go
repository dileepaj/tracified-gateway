package stellarRetriever

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/model"

	"net/http"
)

// type PublicKeyPOC struct {
// 	Name  string
// 	Value string
// }

// type KeysResponsePOC struct {
// 	Collection []PublicKeyPOC
// }

// type ConcretePOC struct {
// 	// *interpreter.AbstractPOC
// 	Txn       string
// 	ProfileID string
// 	DBTree    []model.Current
// 	BCTree    []model.Current
// }

func (db *ConcretePOC) RetrieveFullPOC() model.RetrievePOC {
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
			// var dataHash string

			var transactionType string
			var TDPHash []string
			if keys[0] != (PublicKeyPOC{}) {
				transactionType = Base64DecEnc("Decode", keys[0].Value)
				fmt.Println("transactionType")
				fmt.Println(transactionType)
			}

			if transactionType == "2" {
				if keys[3] != (PublicKeyPOC{}) {
					TDPHash = append(TDPHash, Base64DecEnc("Decode", keys[3].Value))
					fmt.Println("TDPHash")
					fmt.Println(TDPHash)
				}

			}

			if transactionType == "6" {
				if keys[3] != (PublicKeyPOC{}) {
					mergeID := Base64DecEnc("Decode", keys[3].Value)
					result, err := http.Get("https://horizon-testnet.stellar.org/transactions/" + mergeID + "/operations")
					if err != nil {
						Rerr.Code = result.StatusCode
						Rerr.Message = "The HTTP request failed for Merge ID"
						response.Txn = db.Txn
						response.Error = Rerr
						return response
					} else {

						if result.StatusCode == 200 {
							TDPHash = append(TDPHash, Base64DecEnc("Decode", keys[3].Value))
							fmt.Println("TDPHash")
							fmt.Println(TDPHash)
						}

					}

				}

			}

			if keys[1] != (PublicKeyPOC{}) {
				bcPreHash = Base64DecEnc("Decode", keys[1].Value)
				fmt.Println("bcPreHash")
				fmt.Println(bcPreHash)
			}

			temp := model.Current{TXNID: db.Txn, TType: transactionType, DataHash: TDPHash}
			db.BCTree = append(db.BCTree, temp)

			Rerr.Code = result.StatusCode
			Rerr.Message = "The Blockchain Tree Retrieved successfully"
			response.Txn = db.Txn
			response.BCHash = db.BCTree
			response.DBHash = db.DBTree
			response.Error = Rerr

			if keys[1].Value != "" {
				object := ConcretePOC{Txn: bcPreHash, BCTree: db.BCTree, DBTree: db.DBTree, ProfileID: db.ProfileID}
				response = object.RetrieveFullPOC()
			}

			// fmt.Print(response)
		}

		return response
	}

}
