package deprecatedStellarExecuter

import (
	"fmt"
	"net/http"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/model"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
)

type ConcreteGenesis struct {
	InsertGenesisStruct apiModel.InsertGenesisStruct
	// Identifiers   string
	// InsertType    string
	// PreviousTXNID string
}

// var GenesisTxn string

func (cd *ConcreteGenesis) InsertGenesis() model.InsertGenesisResponse {
	publicKey := "GD3EEFYWEP2XLLHONN2TRTQV4H5GSXJGCSUXZJGXGNZT4EFACOXEVLDJ"
	secretKey := "SA46OTS655ZDALIAODVCBWLWBXZWO6VUS6TU4U4GAIUVCKS2SYPDS7N4"
	var response model.InsertGenesisResponse
	response.Identifiers = cd.InsertGenesisStruct.Identifier
	response.TxnType = cd.InsertGenesisStruct.Type

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

	typeTXNBuilder := txnbuild.ManageData{Name: "Transaction Type", Value: []byte(cd.InsertGenesisStruct.Type)}
	CertTypeTXNBuilder := txnbuild.ManageData{Name: "PreviousTXNID", Value: []byte("")}
	IdentifierTXNBuilder := txnbuild.ManageData{Name: "Identifier", Value: []byte(cd.InsertGenesisStruct.Identifier)}

	// BUILD THE GATEWAY XDR
	tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount:        &account,
		IncrementSequenceNum: true,
		Operations:           []txnbuild.Operation{&typeTXNBuilder, &CertTypeTXNBuilder, &IdentifierTXNBuilder},
		BaseFee:              constants.MinBaseFee,
		Memo:                 nil,
		Preconditions:        txnbuild.Preconditions{TimeBounds: constants.TransactionTimeOut},
	})
	// // save data
	// tx, err := build.Transaction(
	// 	commons.GetHorizonNetwork(),
	// 	build.SourceAccount{publicKey},
	// 	build.AutoSequence{commons.GetHorizonClient()},
	// 	build.SetData("Transaction Type", []byte(cd.InsertGenesisStruct.Type)),
	// 	build.SetData("PreviousTXNID", []byte("")),
	// 	build.SetData("Identifiers", []byte(cd.InsertGenesisStruct.Identifier)),
	// )
	if err != nil {
		// panic(err)
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "The HTTP request failed for InsertGenesis "
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
	response.GenesisTxn = resp.Hash

	// cd.PreviousTXNID = resp.Hash

	return response
}
