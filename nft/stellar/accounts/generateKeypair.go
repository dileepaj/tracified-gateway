package accounts

import (
	"log"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/stellar/go/build"
)

func FundAccount(buyerPK string) (string, error) {
	txn, err := build.Transaction(
		commons.GetHorizonNetwork(),
		build.SourceAccount{AddressOrSeed: commons.GoDotEnvVariable("NFTSTELLARISSUERPUBLICKEYK")},
		build.AutoSequence{SequenceProvider: commons.GetHorizonClient()},
		build.Payment(
			build.Destination{AddressOrSeed: buyerPK},
			build.CreditAmount{
				Code:   build.NativeAsset().Code,
				Issuer: commons.GoDotEnvVariable("NFTSTELLARISSUERPUBLICKEYK"),
				Amount: "1",
			},
		),
	)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	signTxn, err := txn.Sign(commons.GoDotEnvVariable("NFTSTELLARISSUERSECRETKEY"))
	if err != nil {
		log.Fatal("Error when submitting the transaction : ", " hError")
		return "", err
	}
	encodedTxn, err := signTxn.Base64()
	if err != nil {
		return "", err
	}
	//submit transaction
	respn, err := commons.GetHorizonClient().SubmitTransaction(encodedTxn)
	if err != nil {
		log.Fatal("Error submitting transaction:", err)
		return "", err
	}
	return respn.Hash, nil
}
