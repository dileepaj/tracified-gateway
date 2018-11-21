package stellarExecuter

import (
	"github.com/tracified-gateway/api/apiModel"
	"github.com/tracified-gateway/model"
	"net/http"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
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

	paymentTx, err := build.Transaction(
		build.SourceAccount{cd.Assest.Issuerkey},
		build.TestNetwork,
		build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
		build.SetData("Transaction Type", []byte(cd.Assest.Type)),
		build.SetData("PreviousTXNID", []byte(cd.Assest.PreviousTXNID)),
		build.SetData("ProfileID", []byte(cd.Assest.ProfileID)),
		build.Payment(
			build.Destination{AddressOrSeed: cd.Assest.Reciverkey},
			build.CreditAmount{cd.Assest.Code, cd.Assest.Issuerkey, cd.Assest.Amount},
		),
	)
	if err != nil {
		// log.Fatal(err)
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "Transaction Build Issues"
		return response
	}
	paymentTxe, err := paymentTx.Sign(cd.Assest.Signer)
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
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(paymentTxeB64)
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
