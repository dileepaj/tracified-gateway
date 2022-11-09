package deprecatedStellarExecuter

import (
	"net/http"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
)

type ConcreteChangeOfCustody struct {
	// *builder.AbstractTDPInsert
	// Code       string
	// Amount     string
	// IssuerKey  string
	// Reciverkey string
	// Sender     string
	COC apiModel.ChangeOfCustody
	// ProfileId string
}

func (cd *ConcreteChangeOfCustody) ChangeOfCustody() model.COCResponse {
	// Second, the issuing account actually sends a payment using the asset
	// RA sign
	// signerSeed := cd.coc.Sender
	// recipientSeed := reciverkey
	var response model.COCResponse
	// // Keys for accounts to issue and receive the new asset
	signerSeed, err := keypair.Parse(cd.COC.Sender)
	if err != nil {
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "Account not found"
		return response
	}
	// recipient, err := keypair.Parse(recipientSeed)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	netClient := commons.GetHorizonClient()
	accountRequest := horizonclient.AccountRequest{AccountID: signerSeed.Address()}
	account, err := netClient.AccountDetail(accountRequest)
	if err != nil {
		// log.Fatal(err)
	}
	asset := txnbuild.CreditAsset{Code: cd.COC.Code, Issuer: cd.COC.IssuerKey}

	typeTXNBuilder := txnbuild.ManageData{Name: "Transaction Type", Value: []byte(cd.COC.Type)}
	CertTypeTXNBuilder := txnbuild.ManageData{Name: "PreviousTXNID", Value: []byte(cd.COC.PreviousTXNID)}
	ProfileIDTXNBuilder := txnbuild.ManageData{Name: "ProfileID", Value: []byte(cd.COC.PreviousProfileID)}
	IdentifierTXNBuilder := txnbuild.ManageData{Name: "Identifier", Value: []byte(cd.COC.PreviousProfileID)}
	paymentTXNBuilder := txnbuild.Payment{
		Destination:   signerSeed.Address(),
		Amount:        cd.COC.Amount,
		Asset:         asset,
		SourceAccount: cd.COC.IssuerKey,
	}

	// BUILD THE GATEWAY XDR
	paymentTx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount:        &account,
		IncrementSequenceNum: true,
		Operations:           []txnbuild.Operation{&typeTXNBuilder, &CertTypeTXNBuilder, &ProfileIDTXNBuilder, &IdentifierTXNBuilder, &paymentTXNBuilder},
		BaseFee:              txnbuild.MinBaseFee,
		Memo:                 nil,
		Preconditions:        txnbuild.Preconditions{},
	})
	// paymentTx, err := build.Transaction(
	// 	build.SourceAccount{cd.COC.Reciverkey},
	// 	commons.GetHorizonNetwork(),
	// 	build.AutoSequence{SequenceProvider: commons.GetHorizonClient()},
	// 	build.SetData("Transaction Type", []byte(cd.COC.Type)),
	// 	build.SetData("PreviousTXNID", []byte(cd.COC.PreviousTXNID)),
	// 	build.SetData("ProfileID", []byte(cd.COC.PreviousProfileID)),
	// 	build.SetData("Identifier", []byte(cd.COC.PreviousProfileID)),
	// 	build.Payment(
	// 		build.SourceAccount{signerSeed.Address()},
	// 		build.Destination{AddressOrSeed: cd.COC.Reciverkey},
	// 		build.CreditAmount{cd.COC.Code, cd.COC.IssuerKey, cd.COC.Amount},
	// 	),
	// )
	if err != nil {
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "Transaction Build Issues"
		return response
	}
	paymentTxe, err := paymentTx.Sign(cd.COC.Sender)
	if err != nil {
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "Transaction 	Singed Issues"
		return response
	}
	paymentTxeB64, err := paymentTxe.Base64()
	if err != nil {
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "Transaction Converter Issues"
		return response
	}
	// resp, err := commons.GetHorizonClient().SubmitTransaction(paymentTxeB64)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	response.Amount = cd.COC.Amount
	response.Code = cd.COC.Code
	response.From = signerSeed.Address()
	response.To = cd.COC.Reciverkey
	response.PreviousProfileID = cd.COC.PreviousProfileID
	response.PreviousTXNID = cd.COC.PreviousTXNID
	response.TxnXDR = paymentTxeB64
	response.Error.Code = http.StatusOK
	response.Error.Message = "Success"
	return response
	// fmt.Println("Successful Transaction:")
	// fmt.Println("Ledger:", resp.Ledger)
	// fmt.Println("Hash:", resp.Hash)
}
