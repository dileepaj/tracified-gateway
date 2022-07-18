package deprecatedStellarExecuter

import (
	"log"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/commons"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
)

type ConcreteCoCLinkage struct {
	// *builder.AbstractTDPInsert
	// Code       string
	// Amount     string
	// IssuerKey  string
	// Reciverkey string
	// Sender     string
	ChangeOfCustodyLink apiModel.ChangeOfCustodyLink
	ProfileId           string
}

func (cd *ConcreteCoCLinkage) CoCLinkage() string {

	// Second, the issuing account actually sends a payment using the asset
	//RA sign
	// signerSeed := cd.coc.Sender
	// recipientSeed := reciverkey

	// // Keys for accounts to issue and receive the new asset
	signerSeed, err := keypair.Parse(cd.ChangeOfCustodyLink.SignerKey)
	if err != nil {
		log.Fatal(err)
	}
	// recipient, err := keypair.Parse(recipientSeed)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	netClient := commons.GetHorizonClient()
	accountRequest := horizonclient.AccountRequest{AccountID: signerSeed.Address()}
	account, err := netClient.AccountDetail(accountRequest)
	if err != nil {
		// log.Fatal(err)
	}

	typeTXNBuilder := txnbuild.ManageData{Name: "Transaction Type", Value: []byte(cd.ChangeOfCustodyLink.Type)}
	CertTypeTXNBuilder := txnbuild.ManageData{Name: "PreviousTXNID", Value: []byte(cd.ProfileId)}
	ProfileIDTXNBuilder := txnbuild.ManageData{Name: "ProfileID", Value: []byte(cd.ProfileId)}
	COCTxnIDTXNBuilder := txnbuild.ManageData{Name: "COCTxnID", Value: []byte(cd.ChangeOfCustodyLink.COCTxn)}

	// BUILD THE GATEWAY XDR
	paymentTx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount:        &account,
		IncrementSequenceNum: true,
		Operations:           []txnbuild.Operation{&typeTXNBuilder, &CertTypeTXNBuilder, &ProfileIDTXNBuilder,&COCTxnIDTXNBuilder},
		BaseFee:              txnbuild.MinBaseFee,
		Memo:                 nil,
		Preconditions:        txnbuild.Preconditions{},
	})
	// paymentTx, err := build.Transaction(
	// 	build.SourceAccount{signerSeed.Address()},
	// 	commons.GetHorizonNetwork(),
	// 	build.AutoSequence{SequenceProvider: commons.GetHorizonClient()},
	// 	build.SetData("Transaction Type", []byte(cd.ChangeOfCustodyLink.Type)),
	// 	build.SetData("PreviousTXNID", []byte(cd.ProfileId)),
	// 	build.SetData("ProfileID", []byte(cd.ProfileId)),
	// 	build.SetData("COCTxnID", []byte(cd.ChangeOfCustodyLink.COCTxn)),
	// )

	if err != nil {
		log.Fatal(err)
	}
	paymentTxe, err := paymentTx.Sign(cd.ChangeOfCustodyLink.SignerKey)
	if err != nil {
		log.Fatal(err)
	}
	paymentTxeB64, err := paymentTxe.Base64()
	if err != nil {
		log.Fatal(err)
	}

	resp, err := commons.GetHorizonClient().SubmitTransactionXDR(paymentTxeB64)
	if err != nil {
		log.Fatal(err)
	}
	return resp.Hash
	// fmt.Println("Successful Transaction:")
	// fmt.Println("Ledger:", resp.Ledger)
	// fmt.Println("Hash:", resp.Hash)
}
