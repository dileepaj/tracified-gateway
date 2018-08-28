package stellarExecuter

import (
	"log"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
)

type ConcreteSendAssest struct {
	// *builder.AbstractTDPInsert
	Code       string
	Amount     string
	Issuerkey  string
	Reciverkey string
	Signer     string
}

func (cd *ConcreteSendAssest) SendAssest() string {

	// Second, the issuing account actually sends a payment using the asset
	//RA sign
	signerSeed := cd.Signer
	// recipientSeed := cd.Reciverkey

	// // Keys for accounts to issue and receive the new asset
	issuer, err := keypair.Parse(signerSeed)
	if err != nil {
		log.Fatal(err)
	}
	// recipient, err := keypair.Parse(recipientSeed)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	paymentTx, err := build.Transaction(
		build.SourceAccount{issuer.Address()},
		build.TestNetwork,
		build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
		build.Payment(
			build.Destination{AddressOrSeed: cd.Reciverkey},
			build.CreditAmount{cd.Code, cd.Issuerkey, cd.Amount},
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
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(paymentTxeB64)
	if err != nil {
		log.Fatal(err)
	}
	return resp.Hash
	// fmt.Println("Successful Transaction:")
	// fmt.Println("Ledger:", resp.Ledger)
	// fmt.Println("Hash:", resp.Hash)
}
