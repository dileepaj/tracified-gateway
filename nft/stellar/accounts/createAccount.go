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
func CreateIssuerAccount() (string,string,error) {
	var NFTAccountKeyEncodedPassword string =commons.GoDotEnvVariable("NFTAccountKeyEncodedPassword")
	//generate new issuer keypair
	pair, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	//NFT issuer keys
	nftIssuerPK := pair.Address()
	nftIssuerSK := pair.Seed()
	//Account Creater keys(Tracified main account)
	issuerPK := commons.GoDotEnvVariable("NFTSTELLARISSUERPUBLICKEYK")
	issuerSK := commons.GoDotEnvVariable("NFTSTELLARISSUERSECRETKEY")
	transactionNft, err := build.Transaction(
		commons.GetHorizonNetwork(),
		build.SourceAccount{AddressOrSeed: issuerPK},
		build.AutoSequence{SequenceProvider: commons.GetHorizonClient()},
		build.CreateAccount(
			build.Destination{AddressOrSeed: nftIssuerPK},
			build.NativeAmount{
				Amount: "2.500000"},
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
	//encrypt the issuer secret key
	encryptedSK,err :=Encrypt(nftIssuerSK,NFTAccountKeyEncodedPassword)
	 if err!=nil{
		 panic(err)
	 }
	log.Println("Transaction Hash for new Account creation: ", respn.Hash)
	return nftIssuerPK,encryptedSK,err
}
