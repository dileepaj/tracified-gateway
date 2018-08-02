package stellarExecuter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Tracified-Gateway/model"
	"github.com/Tracified-Gateway/proofs/builder"
	"github.com/joho/godotenv"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

type ConcreteSplit struct {
	*builder.AbstractSplitProfile
	SplitProfiles string
	PreviousTXNID string
	Identifiers   string
	InsertType    string
	ProfileID     string
	Assets        string
	Code          string
}

func (cd *ConcreteSplit) InsertSplit() model.SplitProfileResponse {

	// // publicKey := "GAEO4AVTWOD6YRC3WFYYXFR6EYYRD2MYKLBB6XTHC3YDUPIEXEIKD5C3"
	// secretKey := "SBSEIZJJXYL6SIC5Y2RDYEQYSBBSRTPSAPGBQPKXGLHC5TZZBC3TSYLC"

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	secretKey := os.Getenv("SECRET_KEY")
	// publicKey := os.Getenv("PUBLIC_KEY")

	var response model.SplitProfileResponse

	data, _ := json.Marshal(&cd.SplitProfiles)
	sEnc := base64.StdEncoding.EncodeToString(data)
	fmt.Println(string(sEnc))

	// save data
	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{secretKey},
		build.AutoSequence{horizon.DefaultTestNetClient},
		build.SetData("TransactionType", []byte(cd.InsertType)),
		build.SetData("PreviousTXNID", []byte(cd.PreviousTXNID)),
		build.SetData("SplitProfiles", []byte(cd.SplitProfiles)),
		build.SetData("Assets", []byte(cd.Assets)),
		build.SetData("Code", []byte(cd.Code)),
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

	return response

}

func (cd *ConcreteSplit) InsertProfile() model.SplitProfileResponse {

	publicKey := "GAEO4AVTWOD6YRC3WFYYXFR6EYYRD2MYKLBB6XTHC3YDUPIEXEIKD5C3"
	// secretKey := "SBSEIZJJXYL6SIC5Y2RDYEQYSBBSRTPSAPGBQPKXGLHC5TZZBC3TSYLC"
	var response model.SplitProfileResponse
	response.PreviousTXNID = cd.PreviousTXNID
	response.PreviousProfileID = cd.ProfileID
	response.Identifiers = cd.Identifiers
	response.TxnType = cd.InsertType

	// save data
	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{publicKey},
		build.AutoSequence{horizon.DefaultTestNetClient},
		build.SetData("PreviousTXNID", []byte(cd.PreviousTXNID)),
		build.SetData("ProfileID", []byte(cd.ProfileID)),
		build.SetData("Identifiers", []byte(cd.Identifiers)),
	)

	if err != nil {
		// panic(err)
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "The HTTP request failed for InsertProfile "
		return response
	}

	// Sign the transaction to prove you are actually the person sending it.
	txe, err := tx.Sign(publicKey)
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

	return response

}
