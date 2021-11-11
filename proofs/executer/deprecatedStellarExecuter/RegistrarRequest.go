package deprecatedStellarExecuter

import (
	"fmt"
	"log"
	"strconv"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/commons"

	"github.com/stellar/go/build"
	"github.com/stellar/go/keypair"
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

	signerAcc, err := keypair.Parse(cd.AppointRegistrar.AccountKey)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure destination account exists
	if _, err := commons.GetHorizonClient().LoadAccount(Registrar); err != nil {
		panic(err)
	}

	// passphrase := network.TestNetworkPassphrase
	weight32, err := strconv.ParseUint(cd.AppointRegistrar.Weight, 10, 64)
	weight := uint32(weight32)
	if err != nil {
		fmt.Println(err)
	}

	low32, err := strconv.ParseUint(cd.AppointRegistrar.Low, 10, 64)
	low := uint32(low32)
	if err != nil {
		fmt.Println(err)
	}
	medium32, err := strconv.ParseUint(cd.AppointRegistrar.Medium, 10, 64)
	medium := uint32(medium32)
	if err != nil {
		fmt.Println(err)
	}
	high32, err := strconv.ParseUint(cd.AppointRegistrar.High, 10, 64)
	high := uint32(high32)
	if err != nil {
		fmt.Println(err)
	}

	tx, err := build.Transaction(
		commons.GetHorizonNetwork(),
		build.SourceAccount{signerAcc.Address()},
		build.AutoSequence{commons.GetHorizonClient()},
		build.AddSigner(Registrar, weight),
		build.SetHighThreshold(high),
		build.SetLowThreshold(low),
		build.SetMediumThreshold(medium),
	)

	if err != nil {
		panic(err)
	}

	// Sign the transaction to prove you are actually the person sending it.
	txe, err := tx.Sign(cd.AppointRegistrar.AccountKey)
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
