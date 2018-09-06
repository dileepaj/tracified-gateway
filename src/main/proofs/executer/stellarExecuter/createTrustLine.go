package stellarExecuter

import (
	"log"
	"main/api/apiModel"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
)

type ConcreteTrustline struct {
	TrustlineStruct apiModel.TrustlineStruct
	// Code      string
	// Limit     string
	// Issuerkey string
	// Signerkey string
}

func (cd *ConcreteTrustline) CreateTrustline() string {

	// router := routes.NewRouter()

	// log.Fatal(http.ListenAndServe(":8030", router))

	// issuerSeed := "SDIYQG2GHVQM3ALVIG3BOJSFNZUOUPRGV3YWWNA7UER5QVCYFLCTQDUT"
	//Reg own account
	recipientSeed := cd.TrustlineStruct.Signerkey

	// Keys for accounts to issue and receive the new asset
	issuer, err := keypair.Parse(cd.TrustlineStruct.Issuerkey)
	if err != nil {
		log.Fatal(err)
	}
	recipient, err := keypair.Parse(recipientSeed)
	if err != nil {
		log.Fatal(err)
	}

	// Create an object to represent the new asset
	Asset := build.CreditAsset(cd.TrustlineStruct.Code, issuer.Address())

	// First, the receiving account must trust the asset
	trustTx, err := build.Transaction(
		build.SourceAccount{recipient.Address()},
		build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
		build.TestNetwork,
		build.Trust(Asset.Code, Asset.Issuer, build.Limit(cd.TrustlineStruct.Limit)),
	)
	if err != nil {
		log.Fatal(err)
	}
	trustTxe, err := trustTx.Sign(recipientSeed)
	if err != nil {
		log.Fatal(err)
	}
	trustTxeB64, err := trustTxe.Base64()
	if err != nil {
		log.Fatal(err)
	}
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(trustTxeB64)
	if err != nil {
		log.Fatal(err)
	}

	return resp.Hash
	// fmt.Println("Successful Transaction:")
	// fmt.Println("Ledger:", resp.Ledger)
	// fmt.Println("Hash:", resp.Hash)
}
