package stellarExecuter

import (
	"log"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
)

type ConcreteChangeOfCustody struct {
	// *builder.AbstractTDPInsert
	Code       string
	Amount     string
	IssuerKey  string
	Reciverkey string
	Sender     string
}

func (cd *ConcreteChangeOfCustody) ChangeOfCustody() string {

	// Second, the issuing account actually sends a payment using the asset
	//RA sign
	signerSeed := cd.Sender
	// recipientSeed := reciverkey

	// // Keys for accounts to issue and receive the new asset
	sendserKey, err := keypair.Parse(signerSeed)
	if err != nil {
		log.Fatal(err)
	}
	// recipient, err := keypair.Parse(recipientSeed)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	paymentTx, err := build.Transaction(
		build.SourceAccount{cd.Reciverkey},
		build.TestNetwork,
		build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
		build.Payment(
			build.SourceAccount{sendserKey.Address()},
			build.Destination{AddressOrSeed: cd.Reciverkey},
			build.CreditAmount{cd.Code, cd.IssuerKey, cd.Amount},
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	paymentTxe, err := paymentTx.Sign(signerSeed)
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
