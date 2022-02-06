package authForIssuer

import (
	"log"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/nft/stellar/accounts"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/build"
)

func SetAuth(CurrentIssuerPK string) error {

	object := dao.Connection{}
	data, err := object.GetNFTIssuerSK(CurrentIssuerPK).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err != nil {
		log.Fatal(err)
		//return "", "", err
	} else {
		nftKeys := data.([]model.NFTKeys)
		//decrypt the secret key
		decrpytNftissuerSecretKey, err := accounts.Decrypt(nftKeys[0].SecretKey, commons.GoDotEnvVariable("NFTAccountKeyEncodedPassword"))
		if data == nil {
			logrus.Error("PublicKey is not found in gateway datastore")
			panic(data)
		}
		txn, err := build.Transaction(
			commons.GetHorizonNetwork(),
			build.SourceAccount{AddressOrSeed: CurrentIssuerPK},
			build.AutoSequence{SequenceProvider: commons.GetHorizonClient()},
			build.SetOptions(build.SetFlag(build.SetAuthRequired()),
				build.SetFlag(build.SetAuthRevocable()),
				build.Signer{
					Address: commons.GoDotEnvVariable("NFTSTELLARISSUERPUBLICKEYK"),
					Weight:  255,
				}),
		)

		if err != nil {
			log.Fatal(err)

		}
		signTxn, err := txn.Sign(decrpytNftissuerSecretKey)
		if err != nil {
			log.Fatal("Error when submitting the transaction : ", " hError")

		}
		encodedTxn, err := signTxn.Base64()
		if err != nil {
			log.Fatal(err)
		}
		//submit transaction
		respn, err := commons.GetHorizonClient().SubmitTransaction(encodedTxn)
		if err != nil {
			log.Fatal("Error submitting transaction:", err)

		}
		log.Println("Hash for locking is", respn.Hash)

	}
	return nil
}
