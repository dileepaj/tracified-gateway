package pools

import (
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/clients/horizonclient"
	sdk "github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	//"github.com/stellar/go/txnbuild"
)

var (
	sponsorPK = commons.GoDotEnvVariable("SPONSORPK")
	sponsorSK = commons.GoDotEnvVariable("SPONSORSeed")
)

var netClient = sdk.DefaultTestNetClient

// CreateSponseredAccount() retur the new stellar account ket pair (created 0 lumen account )
func CreateSponseredAccount(batchAccount model.CoinAccount) (string, string, error) {
	// create keypair
	pair, err := keypair.Random()
	if err != nil {
		logrus.Error("", err)
		return "", "", err
	}

	logrus.Info("Batch Account PK ", pair.Address())
	// encrypt key pair and add to DB
	encSK := commons.Encrypt(pair.Seed())
	// logrus.Info("keys ----------------------+++++++-------------- ",encPK,"   ",encSK)

	logrus.Info("Encrypted PK ", pair.Address())
	// logrus.Info("Encrypted SK ", encSK)

	batchAccount.CoinAccountPK = pair.Address()
	batchAccount.CoinAccountSK = encSK

	object := dao.Connection{}
	errResult := object.InsertAccount(batchAccount)
	if errResult != nil {
		logrus.Info("Error when inserting batch acccount to DB " + errResult.Error())
		return "", "", errResult
	}
	address := pair.Address()
	generatedAccount, err := keypair.ParseFull(pair.Seed())
	if err != nil {
		logrus.Error(err)
		return "", "", err
	}
	request := horizonclient.AccountRequest{AccountID: sponsorPK}
	sponsorAccount, err := commons.GetHorizonNetwork().AccountDetail(request)
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

	signedTx, err := tx.Sign(commons.GetStellarNetwork(), sponsor, generatedAccount)
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
