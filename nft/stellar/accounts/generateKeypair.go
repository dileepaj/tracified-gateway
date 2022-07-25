package accounts

import (
	"log"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func FundAccount(buyerPK string) (string, error) {

	request := horizonclient.AccountRequest{AccountID: commons.GoDotEnvVariable("NFTSTELLARISSUERPUBLICKEYK")}
	issuerAccount, err := commons.GetHorizonNetwork().AccountDetail(request)
	if err != nil {
		return "", err
	}
	issuerSign, err := keypair.ParseFull(commons.GoDotEnvVariable("NFTSTELLARISSUERSECRETKEY"))
	if err != nil {
		return "", err
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &issuerAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&txnbuild.Payment{Destination: buyerPK, Asset:txnbuild.NativeAsset{}, Amount: "1"}},
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
		log.Fatal("Error submitting transaction:", err)
		panic(err)
	}
	return respn.Hash, nil
}
