package stellarExecuter

import (
	"fmt"
	"net/http"

	"github.com/tracified-gateway/api/apiModel"
	"github.com/tracified-gateway/model"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

type ConcreteProfile struct {
	InsertProfileStruct apiModel.InsertProfileStruct
	// Identifiers       string
	// InsertType        string
	// PreviousTXNID     string
	// PreviousProfileID string
}

func (cd *ConcreteProfile) InsertProfile() model.InsertProfileResponse {

	publicKey := "GD3EEFYWEP2XLLHONN2TRTQV4H5GSXJGCSUXZJGXGNZT4EFACOXEVLDJ"
	secretKey := "SA46OTS655ZDALIAODVCBWLWBXZWO6VUS6TU4U4GAIUVCKS2SYPDS7N4"
	var response model.InsertProfileResponse
	response.PreviousTXNID = cd.InsertProfileStruct.PreviousTXNID
	response.PreviousProfileID = cd.InsertProfileStruct.PreviousProfileID
	response.Identifiers = cd.InsertProfileStruct.Identifier
	response.TxnType = cd.InsertProfileStruct.Type

	// save data
	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{publicKey},
		build.AutoSequence{horizon.DefaultTestNetClient},
		build.SetData("Transaction Type", []byte(cd.InsertProfileStruct.Type)),
		build.SetData("PreviousTXNID", []byte(cd.InsertProfileStruct.PreviousTXNID)),
		build.SetData("PreviousProfileID", []byte(cd.InsertProfileStruct.PreviousProfileID)),
		build.SetData("Identifiers", []byte(cd.InsertProfileStruct.Identifier)),
	)

	if err != nil {
		// panic(err)
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "The HTTP request failed for InsertProfile "
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
	response.ProfileTxn = resp.Hash

	return response

}
