package deprecatedStellarExecuter

import (
	"log"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
)

type ConcreteTrustline struct {
	TrustlineStruct apiModel.TrustlineStruct
	// Code      string
	// Limit     string
	// Issuerkey string
	// Signerkey string
}

func (cd *ConcreteTrustline) CreateTrustline() string {
	// router := routes.NewRouter()

	// log.Fatal(http.ListenAndServe(":8030", router))

	// issuerSeed := "SDIYQG2GHVQM3ALVIG3BOJSFNZUOUPRGV3YWWNA7UER5QVCYFLCTQDUT"
	// Reg own account
	recipientSeed := cd.TrustlineStruct.Signerkey

	// Keys for accounts to issue and receive the new asset
	_, err := keypair.Parse(cd.TrustlineStruct.Issuerkey)
	if err != nil {
		log.Fatal(err)
	}
	recipient, err := keypair.Parse(recipientSeed)
	if err != nil {
		log.Fatal(err)
	}

	netClient := commons.GetHorizonClient()
	accountRequest := horizonclient.AccountRequest{AccountID: recipient.Address()}
	account, err := netClient.AccountDetail(accountRequest)
	if err != nil {
		// log.Fatal(err)
	}
	// !Create an object to represent the new asset
	//assetName := txnbuild.CreditAsset{Code: cd.TrustlineStruct.Code, Issuer: issuer.Address()}

	// First, the receiving (distribution) account must trust the asset from the
	// issuer.
	trustTx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &account,
			IncrementSequenceNum: true,
			// !Operations: []txnbuild.TrustLineAsset{&txnbuild.ChangeTrust{
			// 	Line:          assetName,
			// 	Limit:         cd.TrustlineStruct.Limit,
			// 	SourceAccount: "",
			// }},
			BaseFee:       constants.MinBaseFee,
			Memo:          nil,
			Preconditions: txnbuild.Preconditions{TimeBounds: constants.TransactionTimeOut},
		},
	)
	// // Create an object to represent the new asset
	// Asset := build.CreditAsset(cd.TrustlineStruct.Code, issuer.Address())
	// // First, the receiving account must trust the asset
	// trustTx, err := build.Transaction(
	// 	build.SourceAccount{recipient.Address()},
	// 	build.AutoSequence{SequenceProvider: commons.GetHorizonClient()},
	// 	commons.GetHorizonNetwork(),
	// 	build.Trust(Asset.Code, Asset.Issuer, build.Limit(cd.TrustlineStruct.Limit)),
	// )
	if err != nil {
		log.Fatal(err)
	}
	trustTxe, err := trustTx.Sign(recipientSeed)
	if err != nil {
		log.Fatal(err)
	}
	trustTxeB64, err := trustTxe.Base64()
	if err != nil {
		log.Fatal(err)
	}
	resp, err := commons.GetHorizonClient().SubmitTransactionXDR(trustTxeB64)
	if err != nil {
		log.Fatal(err)
	}

	return resp.Hash
	// fmt.Println("Successful Transaction:")
	// fmt.Println("Ledger:", resp.Ledger)
	// fmt.Println("Hash:", resp.Hash)
}
