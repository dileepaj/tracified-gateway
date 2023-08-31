package fosponsoring

import (
	"log"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func FundAccount(FOIssuerPK string) (string, error) {

	request := horizonclient.AccountRequest{AccountID: commons.GoDotEnvVariable("SPONSORERPK")}
	issuerAccount, err := commons.GetHorizonClient().AccountDetail(request)
	if err != nil {
		return "", err
	}
	issuerSign, err := keypair.ParseFull(commons.GoDotEnvVariable("SPONSORERSK"))
	if err != nil {
		return "", err
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &issuerAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&txnbuild.Payment{Destination: FOIssuerPK, Asset: txnbuild.NativeAsset{}, Amount: "10"}},
			BaseFee:              txnbuild.MinBaseFee,
			Memo:                 nil,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
		},
	)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	signedTx, err := tx.Sign(network.TestNetworkPassphrase, issuerSign)
	if err != nil {
		return "", err
	}
	// submit transaction
	respn, err := commons.GetHorizonClient().SubmitTransaction(signedTx)
	if err != nil {
		logger := utilities.NewCustomLogger()
		logger.LogWriter("Error submitting transaction :"+err.Error(), constants.ERROR)
		return "", err
	}
	return respn.Hash, nil
}
