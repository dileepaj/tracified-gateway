package massbalance

import (
	"log"
	"strconv"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

func SetConversion(sellerSourceAccount model.SourceAccount, buyerSourceAccount model.SourceAccount, manageSellOffer model.ManageOffer, manageBuyOffer model.ManageOffer) (string, error) {
	client := horizonclient.DefaultTestNetClient

	buyingasset, creditAsseterr := txnbuild.CreditAsset{Code: manageBuyOffer.TokenName, Issuer: manageBuyOffer.TokenIssuerAccount}.ToChangeTrustAsset()

	if creditAsseterr != nil {
		log.Println("Failed to create credit asset: ", creditAsseterr.Error())
		return "", creditAsseterr
	}

	changeTrustOp := txnbuild.ChangeTrust{
		Line:          buyingasset,
		Limit:         strconv.Itoa(manageSellOffer.Amount * manageSellOffer.UnitPrice),
		SourceAccount: sellerSourceAccount.Source,
	}

	sellOp := txnbuild.ManageSellOffer{
		Selling: txnbuild.CreditAsset{Code: manageSellOffer.TokenName,
			Issuer: manageSellOffer.TokenIssuerAccount},
		Buying: txnbuild.CreditAsset{Code: manageBuyOffer.TokenName,
			Issuer: manageBuyOffer.TokenIssuerAccount},
		Amount: strconv.Itoa(manageSellOffer.Amount),
		Price: xdr.Price{
			N: xdr.Int32(manageSellOffer.UnitPrice),
			D: xdr.Int32(1),
		},
		OfferID:       0,
		SourceAccount: sellerSourceAccount.Source,
	}

	accountRequest := horizonclient.AccountRequest{AccountID: sellerSourceAccount.Source}
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

	senderSK := sellerSourceAccount.Sign
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

func ConvertBatches(sellerSourceAccount model.SourceAccount, buyerSourceAccount model.SourceAccount, manageSellOffer model.ManageOffer, manageBuyOffer model.ManageOffer) (string, error) {
	client := horizonclient.DefaultTestNetClient

	sellingAsset, creditAsseterr := txnbuild.CreditAsset{Code: manageSellOffer.TokenName, Issuer: manageSellOffer.TokenIssuerAccount}.ToChangeTrustAsset()

	if creditAsseterr != nil {
		log.Println("Failed to create credit asset: ", creditAsseterr.Error())
		return "", creditAsseterr
	}

	changeTrustOp := txnbuild.ChangeTrust{
		Line:          sellingAsset,
		Limit:         strconv.Itoa(manageBuyOffer.Amount * manageBuyOffer.UnitPrice),
		SourceAccount: buyerSourceAccount.Source,
	}

	BuyOp := txnbuild.ManageBuyOffer{
		Selling: txnbuild.CreditAsset{Code: manageBuyOffer.TokenName,
			Issuer: manageBuyOffer.TokenIssuerAccount},
		Buying: txnbuild.CreditAsset{Code: manageSellOffer.TokenName,
			Issuer: manageSellOffer.TokenIssuerAccount},
		Amount: strconv.Itoa(manageBuyOffer.Amount),
		Price: xdr.Price{
			N: xdr.Int32(manageBuyOffer.UnitPrice),
			D: xdr.Int32(1),
		},
		OfferID:       0,
		SourceAccount: buyerSourceAccount.Source,
	}

	accountRequest := horizonclient.AccountRequest{AccountID: buyerSourceAccount.Source}
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

	senderSK := buyerSourceAccount.Sign
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
