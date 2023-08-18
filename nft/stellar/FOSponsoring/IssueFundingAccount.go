package fosponsoring

import (
	"log"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/nft/stellar/accounts"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
)

/*
CreateNFTIssuerAccount
@desc - Create new issuer account for each NFT, encrypt the secret key and send the PK and encrypted sk
@params - None
*/
func CreateIssuerAccountForFOUser() (string, string, []byte, error) {
	var AccountKeyEncodedPassword string = commons.GoDotEnvVariable("NFTAccountKeyEncodedPassword")
	// generate new issuer keypair
	pair, err := keypair.Random()
	if err != nil {
		log.Println(err)
	}
	// NFT issuer keys
	issuerPK := pair.Address()
	issuerSK := pair.Seed()

	//send to aws
	encSK := commons.Encrypt(issuerSK)
	// Account Creater keys(Tracified main account)
	mainIssuerPK := commons.GoDotEnvVariable("SPONSORERPK")
	mainIssuerSK := commons.GoDotEnvVariable("SPONSORERSK")
	request := horizonclient.AccountRequest{AccountID: mainIssuerPK}
	issuerAccount, err := commons.GetHorizonClient().AccountDetail(request)
	if err != nil {
		logrus.Error(err)
		return "", "", nil, err
	}
	issuerSign, err := keypair.ParseFull(mainIssuerSK)
	if err != nil {
		return "", "", nil, err
	}
	CreateAccount := []txnbuild.Operation{
		&txnbuild.CreateAccount{
			Destination:   issuerPK,
			Amount:        "10",
			SourceAccount: mainIssuerPK,
		},
	}
	tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount:        &issuerAccount,
		IncrementSequenceNum: true,
		Operations:           CreateAccount,
		BaseFee:              constants.MinBaseFee,
		Memo:                 nil,
		Preconditions:        txnbuild.Preconditions{TimeBounds: constants.TransactionTimeOut},
	})
	if err != nil {
		log.Println("Error when build transaction : ", err)
	}

	signedTx, err := tx.Sign(commons.GetStellarNetwork(), issuerSign)
	if err != nil {
		logrus.Error(err)
		return "", "", nil, err
	}
	// submit transaction
	respn, err := commons.GetHorizonClient().SubmitTransaction(signedTx)
	if err != nil {
		log.Println("Error submitting transaction:", err)
	}
	// encrypt the issuer secret key
	encryptedSK, err := accounts.Encrypt(issuerSK, AccountKeyEncodedPassword)
	if err != nil {
		log.Println(err)
	}
	log.Println("Transaction Hash for new Account creation: ", respn.Hash)
	return issuerPK, encryptedSK, encSK, err
}
