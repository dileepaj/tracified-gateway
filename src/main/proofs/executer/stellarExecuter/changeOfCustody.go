package stellarExecuter

import (
	"log"
	"main/api/apiModel"

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
	COC       apiModel.ChangeOfCustody
	ProfileId string
}

func (cd *ConcreteChangeOfCustody) ChangeOfCustody() string {

	// Second, the issuing account actually sends a payment using the asset
	//RA sign
	// signerSeed := cd.coc.Sender
	// recipientSeed := reciverkey

	// // Keys for accounts to issue and receive the new asset
	signerSeed, err := keypair.Parse(cd.COC.Sender)
	if err != nil {
		log.Fatal(err)
	}
	// recipient, err := keypair.Parse(recipientSeed)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	paymentTx, err := build.Transaction(
		build.SourceAccount{cd.COC.Reciverkey},
		build.TestNetwork,
		build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
		build.Payment(
			build.SourceAccount{signerSeed.Address()},
			build.Destination{AddressOrSeed: cd.COC.Reciverkey},
			build.CreditAmount{cd.COC.Code, cd.COC.IssuerKey, cd.COC.Amount},
		),
		build.SetData("PreviousTXNID", []byte(cd.COC.PreviousTXNID)),
		build.SetData("ProfileID", []byte(cd.ProfileId)),
	)

	if err != nil {
		log.Fatal(err)
	}
	paymentTxe, err := paymentTx.Sign(cd.COC.Sender)
	if err != nil {
		log.Fatal(err)
	}
	paymentTxeB64, err := paymentTxe.Base64()
	if err != nil {
		log.Fatal(err)
	}
	// resp, err := horizon.DefaultTestNetClient.SubmitTransaction(paymentTxeB64)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	return paymentTxeB64
	// fmt.Println("Successful Transaction:")
	// fmt.Println("Ledger:", resp.Ledger)
	// fmt.Println("Hash:", resp.Hash)
}
