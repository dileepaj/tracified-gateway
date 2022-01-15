package stellar

import (

	//"runtime/trace"
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

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	//"github.com/stellar/go/support/env"
	//"github.com/stellar/go/xdr"
	//"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
)

func IssueNft(distributerPK string, assetcode string,TDPtxnhas string, XDR string) (string , string, string) {
	fmt.Println("---------------Issuing NFT-------------")
	fmt.Println("---------------Dis key : ", distributerPK)
	fmt.Println("---------------Asset : ", assetcode)
	fmt.Println("---------------TDP hash : ", TDPtxnhas)
	fmt.Println("---------------XDR : ", XDR)
	
	//Issuer
	issuerPK := "GC6SZI57VRGFULGMBEJGNMPRMDWEJYNL647CIT7P2G2QKNLUHTTOVFO3"
	issuerSK := "SCHSOQDUY2BFFAKO3SPK6WEX4QRTRFUQ7KCI4T6VU4UUGBSLEQYFZCK3"
	
 	//submit transaction for trustline
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(XDR)
	if err != nil {
		hError := err.(*horizon.Error)
		log.Fatal("Error submitting transaction:", hError)
	}

	log.Println("Transaction Hash for trust line: ", resp.Hash)

	
	var b = []byte("POE proofbot URL")

	txn, err := build.Transaction(

		build.SourceAccount{AddressOrSeed: issuerSK},
		build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
		build.TestNetwork,
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
		log.Fatal(err)
	}

	

	txen64, err := txn.Sign(issuerSK)
	if err != nil {
		hError := err.(*horizon.Error)
		log.Fatal("Error when submitting the transaction : ", hError)
	}

	txen, err := txen64.Base64()
	if err != nil {
		panic(err)
	}

	//submit transaction
	respn, err := horizon.DefaultTestNetClient.SubmitTransaction(txen)
	if err != nil {
		hError := err.(*horizon.Error)
		log.Fatal("Error submitting transaction:", hError)
	}

	log.Println("Transaction Hash for NFT: ", respn.Hash)

	return respn.Hash, issuerPK, string(b)
}