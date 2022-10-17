package massbalance

import (
	"log"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func Merge(sender string, amount string, nftname string, destination string, issuer string, limit string) (string, error) {
	client := horizonclient.DefaultTestNetClient

	asset, err := txnbuild.CreditAsset{Code: nftname, Issuer: issuer}.ToChangeTrustAsset()
	if err != nil {
		log.Fatal("Error on asset", err)
	}

	changeTrustOp := txnbuild.ChangeTrust{
		Line:          asset,
		Limit:         limit,
		SourceAccount: destination,
	}

	paymentOp := txnbuild.Payment{
		Destination: destination,
		Amount:      amount,
		Asset: txnbuild.CreditAsset{Code: nftname,
			Issuer: issuer},
		SourceAccount: sender,
	}

	accountRequest := horizonclient.AccountRequest{AccountID: sender}
	sourceAccount, err := client.AccountDetail(accountRequest)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&changeTrustOp, &paymentOp},
			BaseFee:              txnbuild.MinBaseFee,
			Memo:                 nil,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
		},
	)
	if err != nil {
		log.Fatal("Error while trying to build tranaction: ", err)
	}

	senderSK := ""
	senderKeypair, _ := keypair.ParseFull(senderSK)

	txe64, err := tx.Sign(network.TestNetworkPassphrase, senderKeypair)
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
