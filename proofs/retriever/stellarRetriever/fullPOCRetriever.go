package stellarRetriever

import (
	"fmt"
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
	sr := ConcreteStellarTransaction{Txnhash: db.POCStruct.Txn}
	oprns, err1 := sr.RetrieveOperations()

	if err1 != nil {
		Rerr.Code = http.StatusBadRequest
		Rerr.Message = "The HTTP request failed for RetrievePOC"
		response.Txn = db.POCStruct.Txn
		response.Error = Rerr
		return response
	}

	Current := Base64DecEnc("Decode", oprns.Embedded.Records[2].Value)
	GatewayTXNType := Base64DecEnc("Decode", oprns.Embedded.Records[0].Value)

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

		p1Data, errP1Asnc := object.GetTransactionByTxnhash(db.POCStruct.Txn).Then(func(data interface{}) interface{} {
			return data
		}).Await()

		if errP1Asnc != nil || p1Data == nil {
			///ASSIGN PREVIOUS MANAGE DATA BUILDER - THIS WILL BE THE CASE TO ANY SPLIT CHILD
			//DUE TO THE CHILD HAVING A NEW IDENTIFIER
			///ASSIGN PREVIOUS MANAGE DATA BUILDER
		} else {
			result := p1Data.(model.TransactionCollectionBody)
			bcPreHash = result.PreviousTxnHash
			fmt.Println(result.PreviousTxnHash)
		}
	default:
		if oprns.Embedded.Records[1].Value != "" {
			bcPreHash = Base64DecEnc("Decode", oprns.Embedded.Records[1].Value)
		}

	}

	//calls stellar to retrieve the user transaction

	sr1 := ConcreteStellarTransaction{Txnhash: db.POCStruct.Txn}
	oprns1, err := sr1.RetrieveOperations()

	if err != nil {
		Rerr.Code = http.StatusBadRequest
		Rerr.Message = "The HTTP request failed for RetrievePOC"
		response.Txn = db.POCStruct.Txn
		response.Error = Rerr
		return response

	} else {

		var transactionType string
		var TDPHash string
		var Profile string

		//initially checks for the transaction type then decideds on the other fields
		if oprns1.Embedded.Records[0].Value != "" {
			transactionType = Base64DecEnc("Decode", oprns1.Embedded.Records[0].Value)
			// fmt.Println("transactionType")
			// fmt.Println(transactionType)
		}
		switch transactionType {
		//genesis
		case "G0":
			identifier := Base64DecEnc("Decode", oprns1.Embedded.Records[1].Value)
			temp = model.Current{
				TXNID:      Current,
				TType:      transactionType,
				Identifier: identifier}
		//dataPacket
		case "G2":

			// Profile = Base64DecEnc("Decode", oprns1.Embedded.Records[2].Value)
			identifier := Base64DecEnc("Decode", oprns1.Embedded.Records[1].Value)
			TDPHash = Base64DecEnc("Decode", oprns1.Embedded.Records[2].Value)

			temp = model.Current{
				TXNID:      Current,
				TType:      transactionType,
				DataHash:   TDPHash,
				ProfileID:  Profile,
				Identifier: identifier}

		//split parent
		case "G5":
			identifier := Base64DecEnc("Decode", oprns1.Embedded.Records[1].Value)

			temp = model.Current{
				TXNID:      Current,
				TType:      transactionType,
				Identifier: identifier}

		//split child
		case "G6":
			identifier := Base64DecEnc("Decode", oprns1.Embedded.Records[1].Value)

			temp = model.Current{
				TXNID:      Current,
				TType:      transactionType,
				Identifier: identifier}

		case "G10":
			identifier := Base64DecEnc("Decode", oprns1.Embedded.Records[1].Value)

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
		Rerr.Code = http.StatusOK
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

		return response
	}

}
