package stellarRetriever

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"

	"net/http"
)

/*RetrieveFullPOC - WORKING MODEL
@author - Azeem Ashraf
@desc - Retrieves the whole tree from stellar using the last TXN in the chain
@params - XDR
*/
func (db *ConcretePOC) RetrieveFullPOC() model.RetrievePOC {
	object := dao.Connection{}

	var response model.RetrievePOC
	var Rerr model.Error
	// var mergeTree []model.Current
	var temp model.Current
	var bcPreHash string

	// output := make([]string, 20)

	//calls stellar to retrieve the gateway transaction
	result1, err1 := http.Get("https://horizon.stellar.org/transactions/" + db.POCStruct.Txn + "/operations")
	if err1 != nil {
		Rerr.Code = result1.StatusCode
		Rerr.Message = "The HTTP request failed for RetrievePOC"
		response.Txn = db.POCStruct.Txn
		response.Error = Rerr
		return response
	}

	data, _ := ioutil.ReadAll(result1.Body)
	var raw map[string]interface{}
	json.Unmarshal(data, &raw)
	out, _ := json.Marshal(raw["_embedded"])

	var raw1 map[string]interface{}
	json.Unmarshal(out, &raw1)

	out1, _ := json.Marshal(raw1["records"])

	keysBody := out1
	keys := make([]PublicKeyPOC, 0)
	json.Unmarshal(keysBody, &keys)

	Current := Base64DecEnc("Decode", keys[2].Value)
	GatewayTXNType := Base64DecEnc("Decode", keys[0].Value)

	//perform a check on the gateway txn
	switch GatewayTXNType {
	//for split child we realise that the identifier will change for parent
	//thus on the child we look at the profile created
	//and get the txn id of the previous profile and get the last TXN ID
	// case "G6":
	// 	var PreviousIdentifier string
	// 	PreviousProfile := Base64DecEnc("Decode", keys[4].Value)
	// 	p := object.GetProfilebyProfileID(PreviousProfile)
	// 	p.Then(func(data interface{}) interface{} {
	// 		result := data.(model.ProfileCollectionBody)
	// 		PreviousIdentifier = result.Identifier
	// 		return nil
	// 	}).Catch(func(error error) error {
	// 		PreviousIdentifier = ""
	// 		return error
	// 	})
	// 	p.Await()

	// 	p1 := object.GetLastTransactionbyIdentifier(PreviousIdentifier)
	// 	p1.Then(func(data interface{}) interface{} {
	// 		///ASSIGN PREVIOUS MANAGE DATA BUILDER
	// 		result := data.(model.TransactionCollectionBody)
	// 		bcPreHash = result.TxnHash
	// 		return nil
	// 	}).Catch(func(error error) error {
	// 		///ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
	// 		//DUE TO THE CHILD HAVING A NEW IDENTIFIER
	// 		return error
	// 	})
	// 	p1.Await()
	case "G5":

		p1 := object.GetTransactionByTxnhash(db.POCStruct.Txn)
		p1.Then(func(data interface{}) interface{} {
			///ASSIGN PREVIOUS MANAGE DATA BUILDER
			result := data.(model.TransactionCollectionBody)

			bcPreHash = result.PreviousTxnHash
			fmt.Println(result.PreviousTxnHash)
			return nil
		}).Catch(func(error error) error {
			///ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
			//DUE TO THE CHILD HAVING A NEW IDENTIFIER
			return error
		})
		p1.Await()
	default:
		if keys[1] != (PublicKeyPOC{}) {
			bcPreHash = Base64DecEnc("Decode", keys[1].Value)
			fmt.Println("bcPreHash")
			fmt.Println(bcPreHash)
		}

	}

	//calls stellar to retrieve the user transaction
	result, err := http.Get("https://horizon.stellar.org/transactions/" + Current + "/operations")
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
			out, _ := json.Marshal(raw["_embedded"])

			var raw1 map[string]interface{}
			json.Unmarshal(out, &raw1)

			out1, _ := json.Marshal(raw1["records"])

			keysBody := out1
			keys := make([]PublicKeyPOC, 0)
			json.Unmarshal(keysBody, &keys)

			var transactionType string
			var TDPHash string
			var Profile string

			//initially checks for the transaction type then decideds on the other fields
			if keys[0] != (PublicKeyPOC{}) {
				transactionType = Base64DecEnc("Decode", keys[0].Value)
				// fmt.Println("transactionType")
				// fmt.Println(transactionType)
			}
			switch transactionType {
			//genesis
			case "0":
				identifier := Base64DecEnc("Decode", keys[1].Value)
				temp = model.Current{
					TXNID:      Current,
					TType:      transactionType,
					Identifier: identifier}
			//dataPacket
			case "2":

				// Profile = Base64DecEnc("Decode", keys[2].Value)
				identifier := Base64DecEnc("Decode", keys[1].Value)
				TDPHash = Base64DecEnc("Decode", keys[2].Value)

				fmt.Println("TDPHash")
				fmt.Println(TDPHash)

				temp = model.Current{
					TXNID:      Current,
					TType:      transactionType,
					DataHash:   TDPHash,
					ProfileID:  Profile,
					Identifier: identifier}

			//split parent
			case "5":
				identifier := Base64DecEnc("Decode", keys[1].Value)

				temp = model.Current{
					TXNID:      Current,
					TType:      transactionType,
					Identifier: identifier}

			//split child
			case "6":
				identifier := Base64DecEnc("Decode", keys[1].Value)

				temp = model.Current{
					TXNID:      Current,
					TType:      transactionType,
					Identifier: identifier}

			case "10":
				identifier := Base64DecEnc("Decode", keys[1].Value)

				temp = model.Current{
					TXNID:      Current,
					TType:      transactionType,
					Identifier: identifier}
			default:
				// identifier := Base64DecEnc("Decode", keys[1].Value)

				// temp = model.Current{
				// 	TXNID:      Current,
				// 	TType:      transactionType,
				// 	Identifier: identifier}
			}

			db.POCStruct.BCTree = append(db.POCStruct.BCTree, temp)

			Rerr.Code = result.StatusCode
			Rerr.Message = "The Blockchain Tree Retrieved successfully"
			response.Txn = db.POCStruct.Txn
			response.BCHash = db.POCStruct.BCTree
			response.DBHash = db.POCStruct.DBTree
			response.Error = Rerr

			if bcPreHash == "0" {
				bcPreHash = ""
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
