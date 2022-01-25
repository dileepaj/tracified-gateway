package accounts

import (
	"log"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
)

/**
CreateIssuerAccount()
This Method Use for Create New Account for Issue NFT
Created New Account kyes added To GatewayDB NFTKyes collection
**/
func CreateIssuerAccount() (string,string,error) {
	var NFTAccountKeyEncodedPassword string =commons.GoDotEnvVariable("NFTAccountKeyEncodedPassword")
	pair, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	//NFT issuer keys
	nftIssuerPK := pair.Address()
	nftIssuerSK := pair.Seed()
	//Account Creater keys(Tracified)
	issuerPK := commons.GoDotEnvVariable("NFTSTELLARISSUERPUBLICKEYK")
	issuerSK := commons.GoDotEnvVariable("NFTSTELLARISSUERSECRETKEY")
	transactionNft, err := build.Transaction(
		commons.GetHorizonNetwork(),
		build.SourceAccount{AddressOrSeed: issuerPK},
		build.AutoSequence{SequenceProvider: commons.GetHorizonClient()},
		build.CreateAccount(
			build.Destination{AddressOrSeed: nftIssuerPK},
			build.NativeAmount{
				Amount: "2"},
		),
	)
	if err != nil {
		log.Fatal("Error when build transaction : ",err)
		panic(err)
	}
	txen64, err := transactionNft.Sign(issuerSK)
	if err != nil {
		log.Fatal("Error when sign the transaction : ",err)
		panic(err)
	}
	txen, err := txen64.Base64()
	if err != nil {
		panic(err)
	}
	//submit transaction
	respn, err := commons.GetHorizonClient().SubmitTransaction(txen)
	if err != nil {
		log.Fatal("Error submitting transaction:",err)
		panic(err)
	}
	encryptedSK,err :=Encrypt(nftIssuerSK,NFTAccountKeyEncodedPassword)
	 if err!=nil{
		 panic(err)
	 }
	log.Println("Transaction Hash for new Account creation: ", respn.Hash)
	return nftIssuerPK,encryptedSK,err
}
