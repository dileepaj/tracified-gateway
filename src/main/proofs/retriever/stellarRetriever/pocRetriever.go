package stellarRetriever

import (
	// "fmt"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/api/apiModel"
	"main/model"

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
	POCStruct apiModel.POCStruct
	// Txn       string
	// ProfileID string
	// DBTree    []model.Current
	// BCTree    []model.Current
}

func (db *ConcretePOC) RetrievePOC() model.RetrievePOC {
	var response model.RetrievePOC
	var Rerr model.Error
	// output := make([]string, 20)
	result, err := http.Get("https://horizon-testnet.stellar.org/transactions/" + db.POCStruct.Txn + "/operations")
	if err != nil {
		Rerr.Code = result.StatusCode
		Rerr.Message = "The HTTP request failed for RetrievePOC"
		response.Txn = db.POCStruct.Txn
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
			var TDPHash string
			var mergeID string
			if keys[0] != (PublicKeyPOC{}) {
				transactionType = Base64DecEnc("Decode", keys[0].Value)
				fmt.Println("transactionType")
				fmt.Println(transactionType)
			}

			if transactionType == "2" {
				if keys[3] != (PublicKeyPOC{}) {
					TDPHash = Base64DecEnc("Decode", keys[4].Value)
					fmt.Println("TDPHash")
					fmt.Println(TDPHash)
				}

			}

			if transactionType == "6" {
				if keys[3] != (PublicKeyPOC{}) {
					mergeID = Base64DecEnc("Decode", keys[3].Value)
					result, err := http.Get("https://horizon-testnet.stellar.org/transactions/" + mergeID + "/operations")
					if err != nil {
						Rerr.Code = result.StatusCode
						Rerr.Message = "The HTTP request failed for Merge ID"
						response.Txn = db.POCStruct.Txn
						response.Error = Rerr
						return response
					} else {

						if result.StatusCode == 200 {
							mergeID = Base64DecEnc("Decode", keys[3].Value)
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

			temp := model.Current{TXNID: db.POCStruct.Txn, TType: transactionType, DataHash: TDPHash, MergedID: mergeID}
			db.POCStruct.BCTree = append(db.POCStruct.BCTree, temp)

			Rerr.Code = result.StatusCode
			Rerr.Message = "The Blockchain Tree Retrieved successfully"
			response.Txn = db.POCStruct.Txn
			response.BCHash = db.POCStruct.BCTree
			response.DBHash = db.POCStruct.DBTree
			response.Error = Rerr

			if keys[1].Value != "" {
				POCObject := apiModel.POCStruct{Txn: bcPreHash, BCTree: db.POCStruct.BCTree, DBTree: db.POCStruct.DBTree, ProfileID: db.POCStruct.ProfileID}
				object := ConcretePOC{POCStruct: POCObject}
				// object := ConcretePOC{Txn: bcPreHash, BCTree: db.POCStruct.BCTree, DBTree: db.POCStruct.DBTree, ProfileID: db.POCStruct.ProfileID}
				response = object.RetrievePOC()
			}

			// fmt.Print(response)
		}

		return response
	}

}
