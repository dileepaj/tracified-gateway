package deprecatedStellarExecuter

import (
	"fmt"
	"net/http"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
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
	tracifiedAccount, err := keypair.ParseFull(secretKey)
	if err != nil {
		logrus.Error(err)
	}
	var response model.InsertProfileResponse
	response.PreviousTXNID = cd.InsertProfileStruct.PreviousTXNID
	response.PreviousProfileID = cd.InsertProfileStruct.PreviousProfileID
	response.Identifiers = cd.InsertProfileStruct.Identifier
	response.TxnType = cd.InsertProfileStruct.Type

	// netClient := commons.GetHorizonClient()
	// accountRequest := horizonclient.AccountRequest{AccountID: publicKey}
	// account, err := netClient.AccountDetail(accountRequest)
	kp,_ := keypair.Parse(publicKey)
	client := horizonclient.DefaultTestNetClient
	accountRequest := horizonclient.AccountRequest{AccountID: kp.Address()}
	account, err := client.AccountDetail(accountRequest)
	if err != nil {
		// log.Fatal(err)
	}

	typeTXNBuilder := txnbuild.ManageData{Name: "Transaction Type", Value: []byte(cd.InsertProfileStruct.Type)}
	previousTXNBuilder := txnbuild.ManageData{Name: "PreviousTXNID", Value: []byte(cd.InsertProfileStruct.PreviousTXNID)}
	profileIDTXNBuilder := txnbuild.ManageData{Name: "PreviousProfileID", Value: []byte(cd.InsertProfileStruct.PreviousTXNID)}
	identifierTXNBuilder := txnbuild.ManageData{Name: "Identifier", Value: []byte(cd.InsertProfileStruct.Identifier)}
	// BUILD THE GATEWAY XDR
	tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount:        &account,
		IncrementSequenceNum: true,
		Operations:           []txnbuild.Operation{&previousTXNBuilder, &typeTXNBuilder, &identifierTXNBuilder, &profileIDTXNBuilder},
		BaseFee:              constants.MinBaseFee,
		Memo:                 nil,
		Preconditions:        txnbuild.Preconditions{TimeBounds: constants.TransactionTimeOut},
	})
	// save data
	// tx, err := build.Transaction(
	// 	commons.GetHorizonNetwork(),
	// 	build.SourceAccount{publicKey},
	// 	build.AutoSequence{commons.GetHorizonClient()},
	// 	build.SetData("Transaction Type", []byte(cd.InsertProfileStruct.Type)),
	// 	build.SetData("PreviousTXNID", []byte(cd.InsertProfileStruct.PreviousTXNID)),
	// 	build.SetData("PreviousProfileID", []byte(cd.InsertProfileStruct.PreviousProfileID)),
	// 	build.SetData("Identifiers", []byte(cd.InsertProfileStruct.Identifier)),
	// )
	if err != nil {
		// panic(err)
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "The HTTP request failed for InsertProfile "
		return response
	}

	// Sign the transaction to prove you are actually the person sending it.
	txe, err := tx.Sign(commons.GetStellarNetwork(), tracifiedAccount)
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
	resp, err := commons.GetHorizonClient().SubmitTransactionXDR(txeB64)
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
