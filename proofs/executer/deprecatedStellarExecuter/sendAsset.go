package deprecatedStellarExecuter

import (
	"net/http"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
)

type ConcreteSendAssest struct {
	// *builder.AbstractTDPInsert
	// Code       string
	// Amount     string
	// Issuerkey  string
	// Reciverkey string
	// Signer     string
	Assest apiModel.SendAssest
}

func (cd *ConcreteSendAssest) SendAsset() model.SendAssetResponse {

	// Second, the issuing account actually sends a payment using the asset
	//RA sign
	// signerSeed :=
	// recipientSeed := cd.Reciverkey
	var response model.SendAssetResponse
	// // Keys for accounts to issue and receive the new asset

	signerSeed, err := keypair.Parse(cd.Assest.Signer)
	if err != nil {
		// log.Fatal(err)
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "Account not found"
		return response
	}

	// recipient, err := keypair.Parse(recipientSeed)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// netClient := commons.GetHorizonClient()
	// accountRequest := horizonclient.AccountRequest{AccountID: cd.Assest.Issuerkey}
	// account, err := netClient.AccountDetail(accountRequest)
	kp,_ := keypair.Parse(cd.Assest.Issuerkey)
	client := horizonclient.DefaultTestNetClient
	accountRequest := horizonclient.AccountRequest{AccountID: kp.Address()}
	account, err := client.AccountDetail(accountRequest)

	AssestSigner, err := keypair.ParseFull(cd.Assest.Signer)
	if err != nil {
		logrus.Error(err)
	}
	asset := txnbuild.CreditAsset{Code:cd.Assest.Code, Issuer:  cd.Assest.Issuerkey}

	transactionTypeTXNBuilder := txnbuild.ManageData{Name: "Transaction Type", Value: []byte(cd.Assest.Type)}
	previousTXNIDTXNBuilder := txnbuild.ManageData{Name: "PreviousTXNID", Value:  []byte(cd.Assest.Type)}
	profileIDTXNBuilder := txnbuild.ManageData{Name: "ProfileID", Value: []byte(cd.Assest.ProfileID)}
	paymentTXNBuilder := txnbuild.Payment{ Asset:asset,Destination: cd.Assest.Reciverkey, SourceAccount: cd.Assest.Issuerkey,Amount: cd.Assest.Amount}


	// BUILD THE GATEWAY XDR
	paymentTx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount:        &account,
		IncrementSequenceNum: true,
		Operations:           []txnbuild.Operation{&transactionTypeTXNBuilder, &previousTXNIDTXNBuilder, &profileIDTXNBuilder, &paymentTXNBuilder},
		BaseFee:              txnbuild.MinBaseFee,
		Memo:                 nil,
		Preconditions:        txnbuild.Preconditions{},
	})
	
	//!------------------------------------------ not added payment
	
	//!------------------------------------------ not added payment
	
	//!------------------------------------------ not added payment
	// paymentTx, err := build.Transaction(
	// 	build.SourceAccount{cd.Assest.Issuerkey},
	// 	commons.GetHorizonNetwork(),
	// 	build.AutoSequence{SequenceProvider: commons.GetHorizonClient()},
	// 	build.SetData("Transaction Type", []byte(cd.Assest.Type)),
	// 	build.SetData("PreviousTXNID", []byte(cd.Assest.PreviousTXNID)),
	// 	build.SetData("ProfileID", []byte(cd.Assest.ProfileID)),
	// 	build.Payment(
	// 		build.Destination{AddressOrSeed: cd.Assest.Reciverkey},
	// 		build.CreditAmount{cd.Assest.Code, cd.Assest.Issuerkey, cd.Assest.Amount},
	// 	),
	// )
	if err != nil {
		// log.Fatal(err)
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "Transaction Build Issues"
		return response
	}
	paymentTxe, err := paymentTx.Sign(commons.GetStellarNetwork(),AssestSigner)
	if err != nil {
		// log.Fatal(err)
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "Transaction Signed Issues"
		return response
	}
	paymentTxeB64, err := paymentTxe.Base64()
	if err != nil {
		// log.Fatal(err)
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "Transaction Converter Issues"
		return response
	}
	resp, err := commons.GetHorizonClient().SubmitTransactionXDR(paymentTxeB64)
	if err != nil {
		// log.Fatal(err)
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "Transaction Submit Issues"
		return response
	}

	response.Amount = cd.Assest.Amount
	response.Code = cd.Assest.Code
	response.From = signerSeed.Address()
	response.To = cd.Assest.Reciverkey
	response.PreviousProfileID = cd.Assest.ProfileID
	response.PreviousTXNID = cd.Assest.PreviousTXNID
	response.Txn = resp.Hash
	response.Error.Code = http.StatusOK
	response.Error.Message = "Success"

	return response
	// fmt.Println("Successful Transaction:")
	// fmt.Println("Ledger:", resp.Ledger)
	// fmt.Println("Hash:", resp.Hash)
}
