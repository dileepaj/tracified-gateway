package accounts

import (
	"log"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
)

/*CreateNFTIssuerAccount
@desc - Create new issuer account for each NFT, encrypt the secret key and send the PK and encrypted sk
@params - None
*/
func CreateIssuerAccount() (string, string, error) {
	log.Println("-------------------------------------creating an account-------------------")
	var NFTAccountKeyEncodedPassword string = commons.GoDotEnvVariable("NFTAccountKeyEncodedPassword")
	//generate new issuer keypair
	pair, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	log.Println("Key pair generated: ", pair)
	//NFT issuer keys
	nftIssuerPK := pair.Address()
	nftIssuerSK := pair.Seed()
	//Account Creater keys(Tracified main account)
	log.Println("nft issuer PK: ", nftIssuerPK)
	issuerPK := commons.GoDotEnvVariable("NFTSTELLARISSUERPUBLICKEYK")
	issuerSK := commons.GoDotEnvVariable("NFTSTELLARISSUERSECRETKEY")
	log.Println("Starting transaction build")
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
	log.Println("---------1-------------")
	if err != nil {
		log.Fatal("Error when build transaction : ", err)
		panic(err)
	}
	log.Println("---------2-------------")
	txen64, err := transactionNft.Sign(issuerSK)
	log.Println("---------3-------------")
	if err != nil {
		log.Fatal("Error when sign the transaction : ", err)
		panic(err)
	}
	log.Println("---------4-------------")
	txen, err := txen64.Base64()
	log.Println("---------5-------------")
	if err != nil {
		panic(err)
	}
	//submit transaction
	log.Println("---------6-------------")
	respn, err := commons.GetHorizonClient().SubmitTransaction(txen)
	log.Println("---------7-------------")
	if err != nil {
		log.Fatal("Error submitting transaction:", err)
		panic(err)
	}
	//encrypt the issuer secret key
	log.Println("---------8-------------")
	encryptedSK, err := Encrypt(nftIssuerSK, NFTAccountKeyEncodedPassword)
	log.Println("---------9-------------")
	if err != nil {
		panic(err)
	}
	log.Println("Transaction Hash for new Account creation: ", respn.Hash)
	return nftIssuerPK, encryptedSK, err
}
