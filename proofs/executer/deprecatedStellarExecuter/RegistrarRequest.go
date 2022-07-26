package deprecatedStellarExecuter

import (
	"log"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/sirupsen/logrus"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
)

type ConcreteAppointReg struct {
	// *builder.AbstractTDPInsert
	// Publickey string
	// SignerKey string
	// Weight    uint32
	AppointRegistrar apiModel.AppointRegistrar
}

func (cd *ConcreteAppointReg) RegistrarRequest() string {
	// router := routes.NewRouter()

	// log.Fatal(http.ListenAndServe(":8030", router))

	//registrat own account privatekey
	// source := "SD7JWOWBL4Y777WHHTTWTYCFJME7CEOUYS2OZEJTNN5OHN3BYQOBLXLV"
	Registrar := cd.AppointRegistrar.Registrar
	//RA account

	// signerAcc, err := keypair.Parse(cd.AppointRegistrar.AccountKey)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// Make sure destination account exists
	// if _, err := commons.GetHorizonClient().LoadAccount(Registrar); err != nil {
	// 	panic(err)
	// }
		// Get information about the account we just created
		// netClient := commons.GetHorizonClient()
		// accountRequest := horizonclient.AccountRequest{AccountID: Registrar}
		// account, err := netClient.AccountDetail(accountRequest)
		kp,_ := keypair.Parse(Registrar)
		client := horizonclient.DefaultTestNetClient
		accountRequest := horizonclient.AccountRequest{AccountID: kp.Address()}
		account, err := client.AccountDetail(accountRequest)
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		
		appointRegistrar, err := keypair.ParseFull(cd.AppointRegistrar.AccountKey)
		if err != nil {
			logrus.Error(err)
		}
	// // passphrase := network.TestNetworkPassphrase
	// weight32, err := strconv.ParseUint(cd.AppointRegistrar.Weight, 10, 64)
	                                                                                                                                                                                                            
	// weight := uint8(weight32)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// low32, err := strconv.ParseUint(cd.AppointRegistrar.Low, 10, 64)
	// low := uint8(low32)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// medium32, err := strconv.ParseUint(cd.AppointRegistrar.Medium, 10, 64)
	// medium := uint8(medium32)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// high32, err := strconv.ParseUint(cd.AppointRegistrar.High, 10, 64)
	// high := uint8(high32)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// weightT,_ :=txnbuild.Threshold{weight}   
	//_ := txnbuild.Signer{Registrar, txnbuild.Threshold(weight)}
	//addSignerTXNBuilder := txnbuild.ManageData{Name:"Tyepe", Value:[]byte("G"+TxnBody.TxnType)}



				// BUILD THE GATEWAY XDR
				tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
					SourceAccount:        &account,
					IncrementSequenceNum: true,
					Operations:           []txnbuild.Operation{},
					BaseFee:              txnbuild.MinBaseFee,
					Memo:                 nil,
					Preconditions:        txnbuild.Preconditions{},
				})

	// tx, err := txnbuild.Transaction(
	// 	commons.GetHorizonNetwork(),
	// 	build.SourceAccount{signerAcc.Address()},
	// 	build.AutoSequence{commons.GetHorizonClient()},
	// 	build.AddSigner(Registrar, weight),
	// 	build.SetHighThreshold(high),
	// 	build.SetLowThreshold(low),
	// 	build.SetMediumThreshold(medium),
	// )

	if err != nil {
		panic(err)
	}

	// Sign the transaction to prove you are actually the person sending it.
	txe, err := tx.Sign(commons.GetStellarNetwork(),appointRegistrar)
	if err != nil {
		panic(err)
	}

	txeB64, err := txe.Base64()
	if err != nil {
		panic(err)
	}

	// // And finally, send it off to Stellar!
	// resp, err := commons.GetHorizonClient().SubmitTransaction(txeB64)
	// if err != nil {
	// 	panic(err)
	// }
	// txeB64, err := txe.Base64()

	if err != nil {
		panic(err)
	}
	return txeB64
	// fmt.Printf("tx base64: %s", txeB64)

}
