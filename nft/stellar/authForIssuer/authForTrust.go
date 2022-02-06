package authForIssuer

import (
	"log"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/stellar/go/build"
	// "github.com/stellar/go/keypair"
	// "honnef.co/go/tools/analysis/code"
)

func AuthTrust(CurrentIssuerPK string, trustor string, code string) error {
	txn, err := build.Transaction(
		commons.GetHorizonNetwork(),
		build.SourceAccount{AddressOrSeed: CurrentIssuerPK},
		build.AutoSequence{SequenceProvider: commons.GetHorizonClient()},
		// build.SetFlag(
		// 	build.Asset{
		// 		Code:   code,
		// 		Issuer: CurrentIssuerPK,
		// 		Native: false,
		// 	},
		// 	build.Trustor{Address: trustor},
		// 	build.Authorize{Value: true},
		// ),
	)

	if err != nil {
		log.Fatal(err)

	}
	signTxn, err := txn.Sign(commons.GoDotEnvVariable("NFTSTELLARISSUERSECRETKEY"))
	if err != nil {
		log.Fatal("Error when submitting the transaction : ", " hError")

	}
	encodedTxn, err := signTxn.Base64()
	if err != nil {
		log.Fatal(err)
	}
	//submit transaction
	respn, err := commons.GetHorizonClient().SubmitTransaction(encodedTxn)
	if err != nil {
		log.Fatal("Error submitting transaction:", err)

	}
	log.Println("Hash for locking is", respn.Hash)
	return nil

}
