package accounts

import (
	"log"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

/*CreateNFTIssuerAccount
@desc - Create new issuer account for each NFT, encrypt the secret key and send the PK and encrypted sk
@params - None
*/
func CreateIssuerAccount() (string, string, error) {
	var NFTAccountKeyEncodedPassword string = commons.GoDotEnvVariable("NFTAccountKeyEncodedPassword")
	// generate new issuer keypair
	pair, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	// NFT issuer keys
	nftIssuerPK := pair.Address()
	nftIssuerSK := pair.Seed()
	// Account Creater keys(Tracified main account)
	issuerPK := commons.GoDotEnvVariable("NFTSTELLARISSUERPUBLICKEYK")
	issuerSK := commons.GoDotEnvVariable("NFTSTELLARISSUERSECRETKEY")
	request := horizonclient.AccountRequest{AccountID: issuerPK}
	issuerAccount, err := horizonclient.DefaultTestNetClient.AccountDetail(request)
	if err != nil {
		logrus.Error(err)
		return "", "", err
	}
	issuerSign, err := keypair.ParseFull(issuerSK)
	if err != nil {
		return "", "", err
	}
	CreateAccount := []txnbuild.Operation{
		&txnbuild.CreateAccount{
			Destination:   nftIssuerPK,
			Amount:        "0",
			SourceAccount: issuerPK,
		},
	}
	tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount:        &issuerAccount,
		IncrementSequenceNum: true,
		Operations:           CreateAccount,
		BaseFee:              txnbuild.MinBaseFee,
		Memo:                 nil,
		Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
	})
	if err != nil {
		log.Fatal("Error when build transaction : ", err)
		panic(err)
	}

	signedTx, err := tx.Sign(network.TestNetworkPassphrase, issuerSign)
	if err != nil {
		logrus.Error(err)
		return "", "", err
	}
	// submit transaction
	respn, err := commons.GetHorizonClient().SubmitTransaction(signedTx)
	if err != nil {
		log.Fatal("Error submitting transaction:", err)
		panic(err)
	}
	// encrypt the issuer secret key
	encryptedSK, err := Encrypt(nftIssuerSK, NFTAccountKeyEncodedPassword)
	if err != nil {
		panic(err)
	}
	log.Println("Transaction Hash for new Account creation: ", respn.Hash)
	return nftIssuerPK, encryptedSK, err
}