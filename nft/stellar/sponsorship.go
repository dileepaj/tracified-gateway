package stellar

import (
	"log"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func SponsorCreateAccount(buyerPK string, nftname string, issuer string) (string, error) {
	client := horizonclient.DefaultTestNetClient
	beginSponsorship := txnbuild.BeginSponsoringFutureReserves{
		SponsoredID:   buyerPK,
		SourceAccount: commons.GoDotEnvVariable("SPONSORERPK"),
	}

	createAccount := txnbuild.CreateAccount{
		Destination:   buyerPK,
		Amount:        "0",
		SourceAccount: commons.GoDotEnvVariable("SPONSORERPK"),
	}

	endSponsorship := txnbuild.EndSponsoringFutureReserves{
		SourceAccount: buyerPK,
	}

	beginSponsorship1 := txnbuild.BeginSponsoringFutureReserves{
		SponsoredID:   buyerPK,
		SourceAccount: commons.GoDotEnvVariable("SPONSORERPK"),
	}

	asset, err := txnbuild.CreditAsset{Code: nftname, Issuer: issuer}.ToChangeTrustAsset()
	if err != nil {
		log.Fatal("Error on asset", err)
	}

	changeTrustOp := txnbuild.ChangeTrust{
		Line:          asset,
		Limit:         "1",
		SourceAccount: buyerPK,
	}

	endSponsorship1 := txnbuild.EndSponsoringFutureReserves{
		SourceAccount: buyerPK,
	}

	accountRequest := horizonclient.AccountRequest{AccountID: commons.GoDotEnvVariable("SPONSORERPK")}
	sourceAccount, err := client.AccountDetail(accountRequest)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&beginSponsorship, &createAccount, &endSponsorship, &beginSponsorship1, &changeTrustOp, &endSponsorship1},
			BaseFee:              txnbuild.MinBaseFee,
			Memo:                 nil,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
		},
	)
	if err != nil {
		log.Fatal("Error while trying to build tranaction: ", err)
	}

	sposorerSK := commons.GoDotEnvVariable("SPONSORERSK")
	sponsorerKeypair, _ := keypair.ParseFull(sposorerSK)

	txe64, err := tx.Sign(network.TestNetworkPassphrase, sponsorerKeypair)
	if err != nil {
		hError := err.(*horizonclient.Error)
		log.Fatal("Error when submitting the transaction : ", hError)
	}

	txe, err := txe64.Base64()
	if err != nil {
		panic(err)
	}

	return txe, nil
}
