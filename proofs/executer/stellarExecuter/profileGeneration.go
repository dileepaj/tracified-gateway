package stellarExecuter

import (
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

type ConcreteProfile struct {
	*builder.AbstractProfileInsert
	Identifiers       string
	InsertType        string
	PreviousTXNID     string
	PreviousProfileID string
}

func (cd *ConcreteProfile) InsertProfile() model.InsertProfileResponse {

	// publicKey := "GAEO4AVTWOD6YRC3WFYYXFR6EYYRD2MYKLBB6XTHC3YDUPIEXEIKD5C3"
	// // secretKey := "SBSEIZJJXYL6SIC5Y2RDYEQYSBBSRTPSAPGBQPKXGLHC5TZZBC3TSYLC"

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	secretKey := os.Getenv("TRAC_SECRET_KEY")
	publicKey := os.Getenv("TRAC_PUBLIC_KEY")

	var response model.InsertProfileResponse
	response.PreviousTXNID = cd.PreviousTXNID
	response.PreviousProfileID = cd.PreviousProfileID
	response.Identifiers = cd.Identifiers
	response.TxnType = cd.InsertType

	// save data
	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{publicKey},
		build.AutoSequence{horizon.DefaultTestNetClient},
		build.SetData("TransactionType", []byte(cd.InsertType)),
		build.SetData("PreviousTXNID", []byte(cd.PreviousTXNID)),
		build.SetData("PreviousProfileID", []byte(cd.PreviousProfileID)),
		build.SetData("Identifiers", []byte(cd.Identifiers)),
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
	response.Txn = resp.Hash

	return response

}
