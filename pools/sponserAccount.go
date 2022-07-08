package pools

import (
	"fmt"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/clients/horizonclient"
	sdk "github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	//"github.com/stellar/go/txnbuild"
)

var (
	sponsorPK = "GBSRUSQCEKMIJPPIXFVWTP2EAOU3QHWDVBIKAWBLUEK7VEUTZJK3OXLL"
	sponsorSK = "SCYSUJCORHLZ3MJPXI7YCW7K2ASN55UURXWUHTJ4D3LD3UFOC7IXRRK3"
)

var netClient = sdk.DefaultTestNetClient

// CreateSponseredAccount() retur the new stellar account ket pair (created 0 lumen account )
func CreateSponseredAccount(batchAccount model.BatchAccount) (string, string, error) {
	// create keypair
	pair, err := keypair.Random()
	if err != nil {
		logrus.Error("1", err)
		return "", "", err
	}

	logrus.Info(pair.Seed())
	logrus.Info(pair.Address())

	//encrypt key pair and add to DB
	encPK := commons.Encrypt(pair.Address())
	encSK := commons.Encrypt(pair.Seed())

	fmt.Println("Encrypted PK", encPK)
	fmt.Println("Encrypted SK", encSK)

	batchAccount.BatchAccountPK=encPK
	batchAccount.BatchAccountSK=encSK

	object := dao.Connection{}
	errResult := object.InsertBatchAccount(batchAccount)
	if errResult != nil {
		logrus.Info("Error when inserting batch acccount to DB " + errResult.Error())
		return "","",errResult
	}
	address := pair.Address()
	generatedAccount,err := keypair.ParseFull(pair.Seed())
	if err != nil {
		logrus.Error(err)
		return "", "", err
	}
	request := horizonclient.AccountRequest{AccountID: sponsorPK}
	sponsorAccount, err := horizonclient.DefaultTestNetClient.AccountDetail(request)
	if err != nil {
		logrus.Error(err)
		return "", "", err
	}

	sponsor, err := keypair.ParseFull(sponsorSK)
	if err != nil {
		logrus.Error(err)
		return "", "", err
	}
	// sponsering account and create accoun with 0 lumen
	CreateAccount := []txnbuild.Operation{
		&txnbuild.BeginSponsoringFutureReserves{
			SponsoredID:   address,
			SourceAccount: sponsorPK,
		},

		&txnbuild.CreateAccount{
			Destination:   address,
			Amount:        "0",
			SourceAccount: sponsorPK,
		},
		&txnbuild.EndSponsoringFutureReserves{
			SourceAccount: address,
		},
	}

	tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount:        &sponsorAccount,
		IncrementSequenceNum: true,
		Operations:           CreateAccount,
		BaseFee:              txnbuild.MinBaseFee,
		Memo:                 nil,
		Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
	})
	if err != nil {
		logrus.Error(err)
		return "", "", err
	}

	signedTx, err := tx.Sign(network.TestNetworkPassphrase, sponsor,generatedAccount)
	if err != nil {
		logrus.Error(err)
		return "", "", err
	}

	response, err := client.SubmitTransaction(signedTx)
	if err != nil {
		logrus.Error(err)
		return "", "", err
	} else {
		logrus.Info("Batch account created ", response.Hash)
		return pair.Address(), pair.Seed(), nil
	}
}
