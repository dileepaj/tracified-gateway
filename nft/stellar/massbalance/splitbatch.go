package massbalance

import (
	"log"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func Split(destination string, sign string, amount string, nftname string, sender string, issuer string, limit string) (string, error) {
	client := horizonclient.DefaultTestNetClient
	log.Println("data in stellar service: ", destination, amount, nftname, sender, issuer, limit)
	asset, err := txnbuild.CreditAsset{Code: nftname, Issuer: issuer}.ToChangeTrustAsset()
	if err != nil {
		log.Fatal("Error on asset", err)
	}
	log.Println("test1------------------------------")
	changeTrustOp := txnbuild.ChangeTrust{
		Line:          asset,
		Limit:         limit,
		SourceAccount: destination,
	}

	accountRequest := horizonclient.AccountRequest{AccountID: destination}
	sourceAccount, err := client.AccountDetail(accountRequest)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal("Error while trying to build tranaction: ", err)
	}

	// issuerSK := "SBZAXSIYUW4YENKKKFTXCJUDYVOBTRCVLF6DEDL2YT34IINN4OVWXC6P"
	// issuerKeypair, _ := keypair.ParseFull(issuerSK)
	// log.Println("-----------keypairs 111-------------  ", issuerKeypair)

	destSK := sign
	destKeypair, _ := keypair.ParseFull(destSK)
	log.Println("-----------keypairs 222-------------  ", destKeypair)

	txe64, err := tx.Sign(network.TestNetworkPassphrase, destKeypair)
	if err != nil {
		hError := err.(*horizonclient.Error)
		log.Fatal("Error when submitting thetransaction : ", hError)
	}
	log.Println("-----------keypairs 3333-------------  ", txe64)
	respn, err := commons.GetHorizonClient().SubmitTransaction(txe64)
	if err != nil {
		log.Fatal("Error submitting transaction:", err)
		panic(err)
	}
	log.Println("txxxxn---------------------", respn.Hash)
	return respn.Hash, nil
}

func SplitPayment(destination string, sign string, amount string, nftname string, sender string, issuer string, limit string) (string, error) {
	client := horizonclient.DefaultTestNetClient
	log.Println("data in stellar service transferrrrr: ", destination, amount, nftname, sender, issuer, limit)

	paymentOp := txnbuild.Payment{
		Destination: destination,
		Amount:      amount,
		Asset: txnbuild.CreditAsset{Code: nftname,
			Issuer: issuer},
		SourceAccount: sender,
	}
	log.Println("test3------------------------------")

	accountRequest := horizonclient.AccountRequest{AccountID: sender}
	sourceAccount, err := client.AccountDetail(accountRequest)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal("Error while trying to build tranaction: ", err)
	}

	senderSK := "SBZAXSIYUW4YENKKKFTXCJUDYVOBTRCVLF6DEDL2YT34IINN4OVWXC6P"
	senderKeypair, _ := keypair.ParseFull(senderSK)
	log.Println("-----------keypairs 111-------------  ", senderKeypair)

	txe64, err := tx.Sign(network.TestNetworkPassphrase, senderKeypair)
	if err != nil {
		hError := err.(*horizonclient.Error)
		log.Fatal("Error when submitting thetransaction : ", hError)
	}

	respn, err := commons.GetHorizonClient().SubmitTransaction(txe64)
	if err != nil {
		log.Fatal("Error submitting transaction:", err)
		panic(err)
	}
	log.Println("txxxxn---------------------", respn.Hash)
	return respn.Hash, nil
}
