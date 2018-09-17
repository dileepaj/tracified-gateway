package stellarExecuter

import (
	"log"
	"main/api/apiModel"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
)

type ConcreteCoCLinkage struct {
	// *builder.AbstractTDPInsert
	// Code       string
	// Amount     string
	// IssuerKey  string
	// Reciverkey string
	// Sender     string
	ChangeOfCustodyLink apiModel.ChangeOfCustodyLink
	ProfileId           string
}

func (cd *ConcreteCoCLinkage) CoCLinkage() string {

	// Second, the issuing account actually sends a payment using the asset
	//RA sign
	// signerSeed := cd.coc.Sender
	// recipientSeed := reciverkey

	// // Keys for accounts to issue and receive the new asset
	signerSeed, err := keypair.Parse(cd.ChangeOfCustodyLink.SignerKey)
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
		build.SetData("Transaction Type", []byte(cd.ChangeOfCustodyLink.Type)),
		build.SetData("PreviousTXNID", []byte(cd.ProfileId)),
		build.SetData("ProfileID", []byte(cd.ProfileId)),
		build.SetData("COCTxnID", []byte(cd.ChangeOfCustodyLink.COCTxn)),
	)

	if err != nil {
		log.Fatal(err)
	}
	paymentTxe, err := paymentTx.Sign(cd.ChangeOfCustodyLink.SignerKey)
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
