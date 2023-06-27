package stellar

import (
	"log"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
)

func SponsorAccount(userPK string) (string, error) {
	logger := utilities.NewCustomLogger()
	client := commons.GetHorizonClient()
	beginSponsorship := txnbuild.BeginSponsoringFutureReserves{
		SponsoredID:   userPK,
		SourceAccount: commons.GoDotEnvVariable("SPONSORERPK"),
	}

	createAccount := txnbuild.CreateAccount{
		Destination:   userPK,
		Amount:        "10",
		SourceAccount: commons.GoDotEnvVariable("SPONSORERPK"),
	}

	endSponsorship := txnbuild.EndSponsoringFutureReserves{
		SourceAccount: userPK,
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
			Operations:           []txnbuild.Operation{&beginSponsorship, &createAccount, &endSponsorship},
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

	txe64, err := tx.Sign(commons.GetStellarNetwork(), sponsorerKeypair)
	if err != nil {
		hError := err.(*horizonclient.Error)
		log.Fatal("Error when submitting the transaction : ", hError)
	}

	txe, err := txe64.Base64()
	if err != nil {
		logger.LogWriter("Error when converting to b64 the transaction :  : "+err.Error(), constants.ERROR)
		return txe, err
	}

	return txe, nil
}
