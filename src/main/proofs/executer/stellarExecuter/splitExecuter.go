package stellarExecuter

import (
	// "encoding/base64"
	// "encoding/json"
	"fmt"
	"main/api/apiModel"
	"main/model"
	"net/http"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

type ConcreteSplit struct {
	SplitProfileStruct apiModel.SplitProfileStruct
	// PreviousProfileID string
	// PreviousTXNID string
	// Identifier string
	// InsertType    string
	// ProfileID     string
	CurAssets string
	// Code          string
}

func (cd *ConcreteSplit) InsertSplit() model.SplitProfileResponse {

	// publicKey := "GAEO4AVTWOD6YRC3WFYYXFR6EYYRD2MYKLBB6XTHC3YDUPIEXEIKD5C3"
	secretKey := "SBSEIZJJXYL6SIC5Y2RDYEQYSBBSRTPSAPGBQPKXGLHC5TZZBC3TSYLC"

	var response model.SplitProfileResponse

	// save data
	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{secretKey},
		build.AutoSequence{horizon.DefaultTestNetClient},
		build.SetData("TransactionType", []byte(cd.SplitProfileStruct.InsertProfileStruct.Type)),
		build.SetData("PreviousTXNID", []byte(cd.SplitProfileStruct.InsertProfileStruct.PreviousTXNID)),
		build.SetData("ProfileID", []byte(cd.SplitProfileStruct.ProfileID)),
		build.SetData("Assets", []byte(cd.CurAssets)),
		build.SetData("Code", []byte(cd.SplitProfileStruct.Code)),
	)

	if err != nil {
		// panic(err)
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "The HTTP request failed for SplitProfile "
		fmt.Println(err)
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
	response.Txn = resp.Hash
	response.PreviousTXNID = cd.SplitProfileStruct.InsertProfileStruct.PreviousTXNID

	return response

}
