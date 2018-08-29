package stellarExecuter

import (
	"log"
	"main/api/apiModel"

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

func (cd *ConcreteSendAssest) SendAsset() string {

	// Second, the issuing account actually sends a payment using the asset
	//RA sign
	// signerSeed :=
	// recipientSeed := cd.Reciverkey

	// // Keys for accounts to issue and receive the new asset
	signerSeed, err := keypair.Parse(cd.Assest.Signer)
	if err != nil {
		log.Fatal(err)
	}
	// recipient, err := keypair.Parse(recipientSeed)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	paymentTx, err := build.Transaction(
		build.SourceAccount{signerSeed.Address()},
		build.TestNetwork,
		build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
		build.Payment(
			build.Destination{AddressOrSeed: cd.Assest.Reciverkey},
			build.CreditAmount{cd.Assest.Code, cd.Assest.Issuerkey, cd.Assest.Amount},
		),
		build.SetData("PreviousTXNID", []byte(cd.Assest.PreviousTXNID)),
		build.SetData("ProfileID", []byte(cd.Assest.ProfileID)),
	)
	if err != nil {
		log.Fatal(err)
	}
	paymentTxe, err := paymentTx.Sign(cd.Assest.Signer)
	if err != nil {
		log.Fatal(err)
	}
	paymentTxeB64, err := paymentTxe.Base64()
	if err != nil {
		log.Fatal(err)
	}
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(paymentTxeB64)
	if err != nil {
		log.Fatal(err)
	}
	return resp.Hash
	// fmt.Println("Successful Transaction:")
	// fmt.Println("Ledger:", resp.Ledger)
	// fmt.Println("Hash:", resp.Hash)
}
