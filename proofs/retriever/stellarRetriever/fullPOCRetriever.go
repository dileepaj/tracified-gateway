package stellarRetriever

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"

	"net/http"
)
//RetrieveFullPOC ...
func (db *ConcretePOC) RetrieveFullPOC() model.RetrievePOC {
	var response model.RetrievePOC
	var Rerr model.Error
	var mergeTree []model.Current
	var temp model.Current
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
			var transactionType string
			var TDPHash string
			var Profile string
			if keys[0] != (PublicKeyPOC{}) {
				transactionType = Base64DecEnc("Decode", keys[0].Value)
				fmt.Println("transactionType")
				fmt.Println(transactionType)
			}
			switch transactionType {
			case "0":
				identifier := Base64DecEnc("Decode", keys[2].Value)
				temp = model.Current{TXNID: db.POCStruct.Txn, TType: transactionType, Identifier: identifier}
			case "1":
				// previousProfile := Base64DecEnc("Decode", keys[2].Value)
				identifier := Base64DecEnc("Decode", keys[2].Value)
				temp = model.Current{
					TXNID:             db.POCStruct.Txn,
					TType:             transactionType,
					Identifier:        identifier}
			case "2":

				// Profile = Base64DecEnc("Decode", keys[2].Value)
				identifier := Base64DecEnc("Decode", keys[2].Value)
				TDPHash = Base64DecEnc("Decode", keys[3].Value)

				fmt.Println("TDPHash")
				fmt.Println(TDPHash)

				temp = model.Current{TXNID: db.POCStruct.Txn, TType: transactionType, DataHash: TDPHash, ProfileID: Profile ,Identifier:identifier}
			case "3":
			case "4":
			case "5":
				identifier := Base64DecEnc("Decode", keys[2].Value)

				temp = model.Current{TXNID: db.POCStruct.Txn, TType: transactionType,Identifier:identifier}
			case "6":

				mergeID := Base64DecEnc("Decode", keys[3].Value)
				identifier := Base64DecEnc("Decode", keys[2].Value)
				// Profile = Base64DecEnc("Decode", keys[2].Value)
				result, err := http.Get("https://horizon-testnet.stellar.org/transactions/" + mergeID + "/operations")
				if err != nil {
					Rerr.Code = result.StatusCode
					Rerr.Message = "The HTTP request failed for Merge ID"
					response.Txn = db.POCStruct.Txn
					response.Error = Rerr
					return response
				} else {

					if result.StatusCode == 200 {
						POCObject1 := apiModel.POCStruct{Txn: mergeID, BCTree: mergeTree, DBTree: db.POCStruct.DBTree, ProfileID: Profile}

						object := ConcretePOC{POCStruct: POCObject1}
						// object := ConcretePOC{Txn: mergeID, BCTree: mergeTree, DBTree: db.POCStruct.DBTree, ProfileID: Profile}
						response = object.RetrieveFullPOC()

						mergeTree = response.BCHash
					}

				}

				temp = model.Current{
					TXNID:       db.POCStruct.Txn,
					TType:       transactionType,
					ProfileID:   Profile,
					MergedChain: mergeTree,
					MergedID:    mergeID,
					Identifier:identifier}
			case "7":

				temp = model.Current{
					TXNID: db.POCStruct.Txn,
					TType: transactionType}
			case "8":

				temp = model.Current{
					TXNID: db.POCStruct.Txn,
					TType: transactionType}
			case "9":

				temp = model.Current{
					TXNID: db.POCStruct.Txn,
					TType: transactionType}
			default:

			}

			if keys[1] != (PublicKeyPOC{}) {
				bcPreHash = Base64DecEnc("Decode", keys[1].Value)
				fmt.Println("bcPreHash")
				fmt.Println(bcPreHash)
			}

			db.POCStruct.BCTree = append(db.POCStruct.BCTree, temp)

			Rerr.Code = result.StatusCode
			Rerr.Message = "The Blockchain Tree Retrieved successfully"
			response.Txn = db.POCStruct.Txn
			response.BCHash = db.POCStruct.BCTree
			response.DBHash = db.POCStruct.DBTree
			response.Error = Rerr

			if bcPreHash=="0"{
				bcPreHash=""
			}
			if bcPreHash != "" {
				POCObject2 := apiModel.POCStruct{
					Txn:       bcPreHash,
					BCTree:    db.POCStruct.BCTree,
					DBTree:    db.POCStruct.DBTree,
					ProfileID: db.POCStruct.ProfileID}

				object := ConcretePOC{POCStruct: POCObject2}
				response = object.RetrieveFullPOC()
			}

			// fmt.Print(response)
		}

		return response
	}

}
