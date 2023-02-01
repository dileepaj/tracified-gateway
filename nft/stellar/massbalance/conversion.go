package massbalance

import (
	"log"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

func SetConversion(sender string, amount string, SellAsset string, BuyAsset string, SellIssuer string, BuyIssuer string, numerator int, denominator int) (string, error) {
	client := horizonclient.DefaultTestNetClient

	sellingasset, creditAsseterr := txnbuild.CreditAsset{Code: SellAsset, Issuer: SellIssuer}.ToChangeTrustAsset()

	if creditAsseterr != nil {
		log.Println("Failed to create credit asset: ", creditAsseterr.Error())
		return "", creditAsseterr
	}

	changeTrustOp := txnbuild.ChangeTrust{
		Line:          sellingasset,
		Limit:         amount,
		SourceAccount: SellIssuer,
	}
	sellOp := txnbuild.ManageSellOffer{
		Selling: txnbuild.CreditAsset{Code: SellAsset,
			Issuer: SellIssuer},
		Buying: txnbuild.CreditAsset{Code: BuyAsset,
			Issuer: BuyIssuer},
		Amount: amount,
		Price: xdr.Price{
			N: xdr.Int32(numerator),
			D: xdr.Int32(denominator),
		},
		OfferID:       0,
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
			Operations:           []txnbuild.Operation{&changeTrustOp, &sellOp},
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
		log.Fatal("Error when submitting thetransaction : ", hError)
	}

	txe, err := txe64.Base64()
	if err != nil {
		panic(err)
	}

	return txe, nil
}

func ConvertBatches(sender string, amount string, SellAsset string, BuyAsset string, SellIssuer string, BuyIssuer string, numerator int, denominator int) (string, error) {
	client := horizonclient.DefaultTestNetClient

	buyingAsset, creditAsseterr := txnbuild.CreditAsset{Code: BuyAsset, Issuer: BuyIssuer}.ToChangeTrustAsset()

	if creditAsseterr != nil {
		log.Println("Failed to create credit asset: ", creditAsseterr.Error())
		return "", creditAsseterr
	}

	changeTrustOp := txnbuild.ChangeTrust{
		Line:          buyingAsset,
		Limit:         amount,
		SourceAccount: BuyIssuer,
	}

	BuyOp := txnbuild.ManageBuyOffer{
		Selling: txnbuild.CreditAsset{Code: SellAsset,
			Issuer: SellIssuer},
		Buying: txnbuild.CreditAsset{Code: BuyAsset,
			Issuer: BuyIssuer},
		Amount: amount,
		Price: xdr.Price{
			N: xdr.Int32(numerator),
			D: xdr.Int32(denominator),
		},
		OfferID:       0,
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
			Operations:           []txnbuild.Operation{&changeTrustOp, &BuyOp},
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
		log.Fatal("Error when submitting thetransaction : ", hError)
	}

	txe, err := txe64.Base64()
	if err != nil {
		panic(err)
	}

	return txe, nil
}
