package stellarExecuter

import (
	"main/api/apiModel"
	"main/model"
	"net/http"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
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
	//RA sign
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

	paymentTx, err := build.Transaction(
		build.SourceAccount{cd.COC.Reciverkey},
		build.TestNetwork,
		build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
		build.SetData("Transaction Type", []byte(cd.COC.Type)),
		build.SetData("PreviousTXNID", []byte(cd.COC.PreviousTXNID)),
		build.SetData("ProfileID", []byte(cd.COC.PreviousProfileID)),
		build.SetData("Identifier", []byte(cd.COC.PreviousProfileID)),
		build.Payment(
			build.SourceAccount{signerSeed.Address()},
			build.Destination{AddressOrSeed: cd.COC.Reciverkey},
			build.CreditAmount{cd.COC.Code, cd.COC.IssuerKey, cd.COC.Amount},
		),
	)

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
	// resp, err := horizon.DefaultTestNetClient.SubmitTransaction(paymentTxeB64)
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
