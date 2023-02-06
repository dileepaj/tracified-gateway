package massbalance

import (
	"log"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func Merge(sender string, sign string, amount string, nftname string, destination string, issuer string, limit string) (string, error) {
	client := horizonclient.DefaultTestNetClient

	asset, err := txnbuild.CreditAsset{Code: nftname, Issuer: issuer}.ToChangeTrustAsset()
	if err != nil {
		log.Println("Error on asset", err)
	}

	changeTrustOp := txnbuild.ChangeTrust{
		Line:          asset,
		Limit:         limit,
		SourceAccount: destination,
	}

	accountRequest := horizonclient.AccountRequest{AccountID: destination}
	sourceAccount, err := client.AccountDetail(accountRequest)
	if err != nil {
		log.Println(err)
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&changeTrustOp},
			BaseFee:              txnbuild.MinBaseFee,
			Memo:                 nil,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
		},
	)
	if err != nil {
		log.Println("Error while trying to build tranaction: ", err)
	}

	senderSK := "SDTJNIJWXBCBB4RAQD6NIBW6CXPDSK3U2WNMAMUPLLKKC5DIGTTYQ6DE"
	senderKeypair, _ := keypair.ParseFull(senderSK)

	txe64, err := tx.Sign(network.TestNetworkPassphrase, senderKeypair)
	if err != nil {
		hError := err.(*horizonclient.Error)
		log.Println("Error when submitting the transaction : ", hError)
	}

	respn, err := commons.GetHorizonClient().SubmitTransaction(txe64)
	if err != nil {
		log.Println("Error submitting transaction:", err)
		panic(err)
	}
	log.Println("txxxxn trust---------------------", respn.Hash)
	return respn.Hash, nil
}

func TransferMerge(sender string, sign string, amount string, nftname string, destination string, issuer string, limit string) (string, error) {
	client := horizonclient.DefaultTestNetClient

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
		log.Println(err)
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&paymentOp},
			BaseFee:              txnbuild.MinBaseFee,
			Memo:                 nil,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
		},
	)
	if err != nil {
		log.Println("Error while trying to build tranaction: ", err)
	}

	senderSK := sign
	senderKeypair, _ := keypair.ParseFull(senderSK)

	txe64, err := tx.Sign(network.TestNetworkPassphrase, senderKeypair)
	if err != nil {
		hError := err.(*horizonclient.Error)
		log.Println("Error when submitting the transaction : ", hError)
	}

	respn, err := commons.GetHorizonClient().SubmitTransaction(txe64)
	if err != nil {
		log.Println("Error submitting transaction:", err)
		panic(err)
	}
	log.Println("txxxxn- transfer--------------------", respn.Hash)
	return respn.Hash, nil

}
