package deprecatedStellarExecuter

import (
	// "encoding/base64"
	// "encoding/json"
	"fmt"
	"net/http"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
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

	publicKey := "GAEO4AVTWOD6YRC3WFYYXFR6EYYRD2MYKLBB6XTHC3YDUPIEXEIKD5C3"
	secretKey := "SBSEIZJJXYL6SIC5Y2RDYEQYSBBSRTPSAPGBQPKXGLHC5TZZBC3TSYLC"
	tracifiedAccount, err := keypair.ParseFull(secretKey)
	if err != nil {
		logrus.Error(err)
	}
	// netClient := commons.GetHorizonClient()
	// accountRequest := horizonclient.AccountRequest{AccountID: publicKey}
	// account, err := netClient.AccountDetail(accountRequest)

	kp,_ := keypair.Parse(publicKey)
	client := horizonclient.DefaultTestNetClient
	accountRequest := horizonclient.AccountRequest{AccountID: kp.Address()}
	account, err := client.AccountDetail(accountRequest)
	
	var response model.SplitProfileResponse
	transactionTypeTXNBuilder := txnbuild.ManageData{Name: "TransactionType", Value: []byte(cd.SplitProfileStruct.Type)}
	previousTXNIDTXNBuilder := txnbuild.ManageData{Name: "PreviousTXNID", Value: []byte(cd.SplitProfileStruct.PreviousTXNID)}
	profileIDTXNBuilder := txnbuild.ManageData{Name: "ProfileID", Value: []byte(cd.SplitProfileStruct.ProfileID)}
	identifierTXNBuilder := txnbuild.ManageData{Name: "Identifier", Value: []byte(cd.SplitProfileStruct.Identifier)}
	assetsTXNBuilder := txnbuild.ManageData{Name: "Assets", Value: []byte(cd.CurAssets)}
	codeTXNBuilder := txnbuild.ManageData{Name: "Code", Value: []byte(cd.SplitProfileStruct.Code)}

	// BUILD THE GATEWAY XDR
	tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount:        &account,
		IncrementSequenceNum: true,
		Operations:           []txnbuild.Operation{&transactionTypeTXNBuilder, &previousTXNIDTXNBuilder, &profileIDTXNBuilder, &identifierTXNBuilder, &assetsTXNBuilder, &codeTXNBuilder},
		BaseFee:              constants.MinBaseFee,
		Memo:                 nil,
		Preconditions:        txnbuild.Preconditions{TimeBounds: constants.TransactionTimeOut},
	})
	if err != nil {
		// panic(err)
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "The HTTP request failed for SplitProfile "
		fmt.Println(err)
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
	response.Txn = resp.Hash
	response.PreviousTXNID = cd.SplitProfileStruct.PreviousTXNID

	return response
}
