package accounts

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
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
	//"github.com/stellar/go/support/env"
	//"github.com/stellar/go/xdr"
	//"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
)

func CreateIssuerAccount() (string, string) {
	fmt.Println("-------------Create account inside--------------------")
	pair, err := keypair.Random()
	fmt.Println("------------------------Random called------------------")
	if err != nil {
		fmt.Println("--------------------------------", err)
		log.Fatal(err)
	}

	nftIssuerPK := pair.Address()
	nftIssuerSK := pair.Seed()

	//Issuer
	issuerPK := commons.GoDotEnvVariable("NFTSTELLARISSUERPUBLICKEYK")
	issuerSK := commons.GoDotEnvVariable("NFTSTELLARISSUERSECRETKEY")
	//issuerPK := commons.GoDotEnvVariable("GCNR45BE32EEGYBNEVSYMHLMUMBKGIO6AY5RFUAZF745FIA5IIBNIUGD")
	//issuerSK := commons.GoDotEnvVariable("SB5BST4V2RC3K5QTYJRDGNSVFA2QHJQDPLR4YJLZND5VQ73XNCNZLWG2")

	transactionNft, errpk := build.Transaction(
		commons.GetHorizonNetwork(),
		build.SourceAccount{AddressOrSeed: issuerPK},
		build.AutoSequence{SequenceProvider: commons.GetHorizonClient()},
		build.CreateAccount(
			build.Destination{AddressOrSeed: nftIssuerPK},
			build.NativeAmount{
				Amount: "1"},
		),
	)

	if errpk != nil {
		//" hError" := err.(*horizon.Error)
		fmt.Println("error1", "errr")
		log.Fatal(errpk)

	}

	txen64, errr := transactionNft.Sign(issuerSK)
	if errr != nil {
		//" hError" := err.(*horizon.Error)
		log.Fatal("Error when submitting the transaction : ", " hError")

	}

	txen, err1 := txen64.Base64()
	if err1 != nil {
		fmt.Println("error base64")
		//" hError" := err.(*horizon.Error)

	}

	//submit transaction
	respn, err2 := commons.GetHorizonClient().SubmitTransaction(txen)
	if err2 != nil {
		fmt.Println("error---------------", err)
		//" hError" := err.(*horizon.Error)
		log.Fatal("Error submitting transaction:", " hError")

	}

	fmt.Println("rsponse---------------------------------", respn.Result)
	log.Println("Transaction Hash for new Account creation: ", respn.Hash)

	return nftIssuerPK, nftIssuerSK
}
