package deprecatedStellarExecuter

import (
	"fmt"
	"net/http"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"
)

type ConcreteInsertPOA struct {
	InsertPOAStruct apiModel.InsertPOAStruct
}

func (cd *ConcreteInsertPOA) InsertPOAHash() model.InsertDataResponse {

	publicKey := "GD3EEFYWEP2XLLHONN2TRTQV4H5GSXJGCSUXZJGXGNZT4EFACOXEVLDJ"
	secretKey := "SA46OTS655ZDALIAODVCBWLWBXZWO6VUS6TU4U4GAIUVCKS2SYPDS7N4"
	var response model.InsertDataResponse
	response.ProfileID = cd.InsertPOAStruct.ProfileID
	response.TxnType = cd.InsertPOAStruct.Type

	// save data
	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{publicKey},
		build.AutoSequence{horizon.DefaultTestNetClient},
		build.SetData("Transaction Type", []byte(cd.InsertPOAStruct.Type)),
		build.SetData("PreviousTXNID", []byte(cd.InsertPOAStruct.PreviousTXNID)),
		build.SetData("ProfileID", []byte(cd.InsertPOAStruct.ProfileID)),
		build.SetData("Identifier", []byte(cd.InsertPOAStruct.Identifier[0])),
	)

	if err != nil {
		// panic(err)
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "The HTTP request failed for InsertPOAHash "
		return response
	}

	// Sign the transaction to prove you are actually the person sending it.
	txe, err := tx.Sign(secretKey)
	if err != nil {
		// panic(err)
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "signing request failed for the Transaction"
		return response
	}

	txeB64, err := txe.Base64()
	if err != nil {
		// panic(err)
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "Base64 conversion failed for the Transaction"
		return response
	}

	// And finally, send it off to Stellar!
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		// panic(err)
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "Test net client crashed"
		return response
	}

	fmt.Println("Successful Transaction:")
	fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)

	response.Error.Code = http.StatusOK
	response.Error.Message = "Transaction performed in the blockchain."
	response.TDPID = resp.Hash

	return response

}