package stellar

import (
	"log"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/nft/stellar/accounts"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/build"
)

func IssueNft(CurrentIssuerPK string, distributerPK string, assetcode string, TDPtxnhas string) (string, string, string, error) {
	object := dao.Connection{}
	data, err := object.GetNFTIssuerSK(CurrentIssuerPK).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err != nil {
		log.Fatal(err)
		return "", "", "", err
	}else{
		nftKeys :=data.([]model.NFTKeys)
		decrpytNftissuerSecretKey,err:=accounts.Decrypt(nftKeys[0].SecretKey,commons.GoDotEnvVariable("NFTAccountKeyEncodedPassword"))
		if data == nil {
			logrus.Error("PublicKey is not found in gateway datastore")
			panic(data)
		}
		var nftConten = []byte("POE proofbot URL")
		txn, err := build.Transaction(
			commons.GetHorizonNetwork(),
			build.SourceAccount{AddressOrSeed: CurrentIssuerPK},
			build.AutoSequence{SequenceProvider: commons.GetHorizonClient()},
			build.Payment(
				build.Destination{AddressOrSeed: distributerPK},
				build.CreditAmount{
					Code:   assetcode,
					Issuer: CurrentIssuerPK,
					Amount: "1",
				},
			),
			build.SetData(assetcode, nftConten),
			build.SetOptions(build.HomeDomain("https://tracified.com")),)
		if err != nil {
			log.Fatal(err)
			return "", "", "", err
		}
		signTxn, err := txn.Sign(decrpytNftissuerSecretKey)
		if err != nil {
			log.Fatal("Error when submitting the transaction : ", " hError")
			return "", "", "", err
		}
		encodedTxn, err := signTxn.Base64()
		if err != nil {
			return "", "", "", err
		}
		//submit transaction
		respn, err := commons.GetHorizonClient().SubmitTransaction(encodedTxn)
		if err != nil {
			log.Fatal("Error submitting transaction:",err)
			return "", "", "", err
		}
		return respn.Hash, CurrentIssuerPK, string(nftConten), nil
	}

}
