package stellar

import (

	//"runtime/trace"

	"fmt"
	"log"

	//"github.com/stellar/go/clients/horizonclient"

	//"github.com/stellar/go/txnbuild"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

func IssueNft() {
	fmt.Println("---------------------------------------------------------------------this is sooooooooooooooooooooooooooooooooooooooooo ---------------------------------")
	/*IssuerPublicKey := "GCPEHZQOQ7BAY4XRCUHENR6FCYT2QANN5GZVS7DRWRNFDUDEQSHSWAT6"
	IssuerSecretKey := "SDAUCB5BXAL25LALK3WFJ4ZHOBPIPVWVHQ5WRKTHOLHW3CCFBUG7EDID"
	IssuerKeypair, _ := keypair.Parse(IssuerSecretKey)

	//Public Key	GAQIWUSQVN4X74KOJKKKRTWMW3LJCN2KYFFSKYFC7IHC25HUYI3VDHAN
	//Secret Key	SBXVSR2FL22OEWGCICTNIG47BCPBYWTLR2H6Z2RSRPZU75Z4YPMS4OXO
	// Distributor Credentials - replace with your respective keys
	DistributorPublicKey := "GCIBSEDAHFAHF45YJNQTNR4VB3C3YDSJQK4JI5DIM5N5AA65ZXTYAD3F"
	DistributorSecretKey := "SDLYEW4ABCH5WM5KPMIBFO7YMVVTQAW57ZFGBWN7UCZKROAKZSN2IKVV"
	DistributorKeypair, _ := keypair.Parse(DistributorSecretKey)*/
	seed := "SASBKAMDPK67S2WL7NXQL5QDAYDWEVITLYCHUC3B6UBCULNQISHGYN4V"
	seedI := "SD374J7BKTZ4AKONZYB3EF5DWQ3ZNI2ROYVFZSYUKZKQMLER7ZIZDACS"
	//client := horizon.DefaultTestNetClient

	//Get information on distributor
	/*accountRequest := horizonclient.AccountRequest{AccountID: seed}
	sourceAccount, err := client.AccountDetail(accountRequest)
	if err != nil {
		log.Fatal(err)
	}*/

	//Get information on Issuer
	/*accountRequestn := horizonclient.AccountRequest{AccountID: seedI}
	sourceAccountn, err := client.AccountDetail(accountRequestn)
	if err != nil {
		log.Fatal(err)
	}*/

	//asset := txnbuild.NativeAsset{}
	/*asset, err := build.CreditAsset{Code: "StellarNFT", Issuer: IssuerPublicKey}.ToChangeTrustAsset()
	if err != nil {
		log.Fatal("Error on asset", err)
	}

	changeTrustOp := build.ChangeTrust{
		Line:          asset,
		Limit:         "1",
		SourceAccount: DistributorPublicKey,
	}*/

	tx, err := build.Transaction(

		build.SourceAccount{AddressOrSeed: seed},
		build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
		build.TestNetwork,
		build.Trust("StellarNFT", "GDJ6A3ZNO4T4NOFR4AE6SPH6OW6DDLNALAHPDTJSVV2QUYONOECKKWXJ", build.Limit("1")),
		//Network:              network.TestNetworkPassphrase,
		// Use a real timeout in production!

	)
	if err != nil {
		log.Fatal(err)
	}

	//sign transaction

	//convert to base64
	//txe, err := tx.BuildSignEncode(DistributorKeypair)
	//log.Println("Transaction base64: ", txe)

	txe64, err := tx.Sign(seed)
	if err != nil {
		hError := err.(*horizon.Error)
		log.Fatal("Error when submitting the transaction : ", hError)
	}

	txe, err := txe64.Base64()
	if err != nil {
		panic(err)
	}

	//fmt.Println("Signed XDR is: ", txe)

	//submit transaction
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(txe)
	if err != nil {
		hError := err.(*horizon.Error)
		log.Fatal("Error submitting transaction:", hError)
	}

	log.Println("Transaction Hash: ", resp.Hash)

	/*manageDataOperation := build.ManageData{
		Name:  "nftData",
		Value: []byte("http://localhost:8000/tdp/proofs/618e8cfa39de283083f97743qs"),
	}

	//
	paymentOp := build.Payment{
		Destination: DistributorPublicKey,
		Amount:      "1",
		Asset: txnbuild.CreditAsset{
			Code:   "StellarNFT",
			Issuer: IssuerPublicKey,
		},
		//Asset:         txnbuild.NativeAsset{},
		SourceAccount: IssuerPublicKey,
	}

	setOptionsOp := build.SetOptions{
		MasterWeight: build.NewThreshold(0),
		HomeDomain:   build.NewHomeDomain("https://tracified.com/"),
	}*/

	//set master weight to 0 to lock account
	//setOptionsOp := txnbuild.SetOptions{MasterWeight: txnbuild.NewThreshold(0)}
	var b = []byte("Hello, goodbye, etc!")
	//submit transaction

	txn, err := build.Transaction(

		build.SourceAccount{AddressOrSeed: seedI},
		build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
		build.TestNetwork,
		build.Payment(
			build.Destination{AddressOrSeed: "GDZ6KUMEJGCLESNKAHWBC2GCABKXXZ4IYPEGPXBV6ZBAJQOL2MLCJE5J"},
			build.CreditAmount{
				Code:   "StellarNFT",
				Issuer: "GDJ6A3ZNO4T4NOFR4AE6SPH6OW6DDLNALAHPDTJSVV2QUYONOECKKWXJ",
				Amount: "1",
			},
		),
		build.SetData("NFTTest", b),
		build.SetOptions(
			build.HomeDomain("Hello")),

		//Network:              network.TestNetworkPassphrase,
		// Use a real timeout in production!

	)
	if err != nil {
		log.Fatal(err)
	}

	//sign transaction

	//convert to base64
	//txe, err := tx.BuildSignEncode(DistributorKeypair)
	//log.Println("Transaction base64: ", txe)

	txen64, err := txn.Sign(seedI)
	if err != nil {
		hError := err.(*horizon.Error)
		log.Fatal("Error when submitting the transaction : ", hError)
	}

	txen, err := txen64.Base64()
	if err != nil {
		panic(err)
	}

	//fmt.Println("Signed XDR is: ", txe)

	//submit transaction
	respn, err := horizon.DefaultTestNetClient.SubmitTransaction(txen)
	if err != nil {
		hError := err.(*horizon.Error)
		log.Fatal("Error submitting transaction:", hError)
	}

	log.Println("Transaction Hash: ", respn.Hash)

	/*setOptionsOp := txnbuild.SetOptions{
		MasterWeight: txnbuild.NewThreshold(0),
		HomeDomain:   txnbuild.NewHomeDomain("https://tracified.com/"),
	}

	txnn, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&setOptionsOp},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(),
			//Network:              network.TestNetworkPassphrase,
			// Use a real timeout in production!
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	//sign transaction

	//convert to base64
	//txe, err := tx.BuildSignEncode(DistributorKeypair)
	//log.Println("Transaction base64: ", txe)

	txenn64, err := txnn.Sign(network.TestNetworkPassphrase, DistributorKeypair)
	if err != nil {
		hError := err.(*horizonclient.Error)
		log.Fatal("Error when submitting the transaction : ", hError)
	}

	txenn, err := txenn64.Base64()
	if err != nil {
		panic(err)
	}

	//fmt.Println("Signed XDR is: ", txe)

	//submit transaction
	respnn, err := client.SubmitTransactionXDR(txenn)
	if err != nil {
		hError := err.(*horizonclient.Error)
		log.Fatal("Error submitting transaction:", hError)
	}

	log.Println("Transaction Hash: ", respnn.Hash)*/

}
