package pools

import (
	"log"

	sdk "github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	//"github.com/stellar/go/txnbuild"
)

var (
	accountCreatorPK = "GBSRUSQCEKMIJPPIXFVWTP2EAOU3QHWDVBIKAWBLUEK7VEUTZJK3OXLL"
	accountCreatorSK = "SCYSUJCORHLZ3MJPXI7YCW7K2ASN55UURXWUHTJ4D3LD3UFOC7IXRRK3"
)

var netClient = sdk.DefaultTestNetClient

func createSponseredAccount()(string, string, error){

	//create keypair
	pair, err := keypair.Random()
    if err != nil {
        log.Fatal(err)
    }

	log.Println(pair.Seed())
	log.Println(pair.Address())

	//begin sponser account
	// creatorAccount, err := netClient.AccountDetail(sdk.AccountRequest{AccountID: accountCreatorPK})
	// if err != nil {
	// 	return "" , "", err
	// }

	// creator, err := keypair.ParseFull(accountCreatorSK)
	// if err != nil {
	// 	return "", "",err
	// }

	// tx, err := txnbuild.NewTransaction(
	// 	txnbuild.TransactionParams{
	// 		SourceAccount: &creatorAccount,
	// 		IncrementSequenceNum: true,
	// 		Operations: []txnbuild.Operation{&txnbuild.BeginSponsoringFutureReserves{}},
	// 		BaseFee: txnbuild.MinBaseFee,
	// 		Memo: nil,
	// 		Preconditions:txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
	// 	}
	// )

	return pair.Address(), pair.Seed(), nil
}