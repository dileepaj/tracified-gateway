package stellar

import (

	//"runtime/trace"
	"errors"
	"fmt"
	"log"

	//"strings"

	//"github.com/segmentio/ksuid"
	//"github.com/rs/xid"
	//"github.com/sony/sonyflake"
	//"github.com/lithammer/shortuuid"
	//"github.com/google/uuid"

	//"github.com/stellar/go/clients/horizonclient"
	//"github.com/stellar/go/txnbuild"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/stellar/go/build"
	//"github.com/stellar/go/support/env"
	//"github.com/stellar/go/xdr"
	//"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
)

func IssueNft(distributerPK string, assetcode string,TDPtxnhas string) (string , string, string,error) {
	fmt.Println("---------------Issuing NFT-------------")
	fmt.Println("---------------Dis key : ", distributerPK)
	fmt.Println("---------------Asset : ", assetcode)
	fmt.Println("---------------TDP hash : ", TDPtxnhas)

	//Issuer
	issuerPK := commons.GoDotEnvVariable("NFTSTELLARISSUERPUBLICKEYK")
	issuerSK := commons.GoDotEnvVariable("NFTSTELLARISSUERSECRETKEY")
	
	// log.Println("Transaction Hash for trust line: ", resp.Hash)
	fmt.Println("---------------issuerSK : ", issuerPK)
	fmt.Println("---------------issuerSK : ", issuerSK)
	var b = []byte("POE proofbot URL")

	txn, err := build.Transaction(
		commons.GetHorizonNetwork(),
		build.SourceAccount{AddressOrSeed: issuerSK},
		build.AutoSequence{SequenceProvider: commons.GetHorizonClient()},
		build.Payment(
			build.Destination{AddressOrSeed: distributerPK},
			build.CreditAmount{
				Code:   assetcode,
				Issuer: issuerPK,
				Amount: "1",
			},
		),
		build.SetData(assetcode, b),
		build.SetOptions(
			build.HomeDomain("https://tracified.com")),

		//Network:              network.TestNetworkPassphrase,
		// Use a real timeout in production!

	)
	if err != nil {
		//" hError" := err.(*horizon.Error)
		fmt.Println("error1",	"errr")
		log.Fatal(err)
		return "","","",errors.New("errr")	
	}

	txen64, err := txn.Sign(issuerSK)
	if err != nil {	
		//" hError" := err.(*horizon.Error)
		log.Fatal("Error when submitting the transaction : "," hError")
		return "","","",errors.New("errr")
	}

	txen, err := txen64.Base64()
	if err != nil {
		fmt.Println("error base64")
		//" hError" := err.(*horizon.Error)
		return "","","",errors.New("errr")
	}
	
	//submit transaction
	respn, err := commons.GetHorizonClient().SubmitTransaction(txen)
	if err != nil {
		fmt.Println("error---------------",err)
		//" hError" := err.(*horizon.Error)
		log.Fatal("Error submitting transaction:"," hError")
		return "","","",errors.New("errr")
	}
	fmt.Println("rsponse---------------------------------",respn.Result)
	log.Println("Transaction Hash for NFT: ", respn.Hash)

	return respn.Hash, issuerPK, string(b),nil
}