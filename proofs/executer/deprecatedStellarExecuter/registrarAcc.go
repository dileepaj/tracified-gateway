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

type ConcreteRegistrarAcc struct {
	// *builder.AbstractTDPInsert
	// Publickey  string
	// Publickey1 string
	// Publickey2 string
	// Publickey3 string
	// SignerKey  string
	// Weight     uint32
	// Low        uint32
	// Medium     uint32
	// High       uint32
	RegistrarAccount apiModel.RegistrarAccount
}

func (cd *ConcreteRegistrarAcc) SetupAccount() string {

	// router := routes.NewRouter()

	// log.Fatal(http.ListenAndServe(":8030", router))

	//RA account
	// source := cd.RegistrarAccount.SignerKey
	destination := cd.RegistrarAccount.SignerKeys[0].Publickey
	registrar, err := keypair.Parse(cd.RegistrarAccount.SignerKey)
	if err != nil {
		log.Fatal(err)
		return "Account not found"
	}
	// Make sure destination account exists
	if _, err := commons.GetHorizonClient().LoadAccount(destination); err != nil {
		panic(err)
		return "Account not found"
	}

	// passphrase := network.TestNetworkPassphrase

	low32, err := strconv.ParseUint(cd.RegistrarAccount.Low, 10, 64)
	low := uint32(low32)
	if err != nil {
		fmt.Println(err)
		return "value convertion error"
	}
	medium32, err := strconv.ParseUint(cd.RegistrarAccount.Medium, 10, 64)
	medium := uint32(medium32)
	if err != nil {
		fmt.Println(err)
		return "value convertion error"
	}
	high32, err := strconv.ParseUint(cd.RegistrarAccount.High, 10, 64)
	high := uint32(high32)
	if err != nil {
		fmt.Println(err)
		return "value convertion error"
	}

	// weight11, err := strconv.ParseUint(cd.RegistrarAccount.SignerKeys[0].Weight, 10, 64)
	// weight1 := uint32(weight11)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// weight22, err := strconv.ParseUint(cd.RegistrarAccount.SignerKeys[1].Weight, 10, 64)
	// weight2 := uint32(weight22)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// weight33, err := strconv.ParseUint(cd.RegistrarAccount.SignerKeys[2].Weight, 10, 64)
	// weight3 := uint32(weight33)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// weight44, err := strconv.ParseUint(cd.RegistrarAccount.SignerKeys[3].Weight, 10, 64)
	// weight4 := uint32(weight44)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	muts := []build.TransactionMutator{
		build.SourceAccount{registrar.Address()},
		commons.GetHorizonNetwork(),
		build.AutoSequence{SequenceProvider: commons.GetHorizonClient()},
	}
	arrSize := len(cd.RegistrarAccount.SignerKeys)
	fmt.Println(arrSize)
	for i := 0; i < arrSize; i++ {
		weight, err := strconv.ParseUint(cd.RegistrarAccount.SignerKeys[i].Weight, 10, 64)
		weight64 := uint32(weight)
		if err != nil {
			fmt.Println(err)
			return "value convertion error"
		}
		ops := []build.TransactionMutator{

			build.AddSigner(cd.RegistrarAccount.SignerKeys[i].Publickey, weight64),
		}
		muts = append(muts, ops...)

	}
	opsLow := []build.TransactionMutator{

		build.SetLowThreshold(low),
	}
	muts = append(muts, opsLow...)

	opsMedium := []build.TransactionMutator{

		build.SetMediumThreshold(medium),
	}
	muts = append(muts, opsMedium...)

	opsHigh := []build.TransactionMutator{

		build.SetHighThreshold(high),
	}
	muts = append(muts, opsHigh...)

	// tx, err := build.Transaction(
	// 	commons.GetHorizonNetwork(),
	// 	build.SourceAccount{registrar.Address()},
	// 	build.AutoSequence{commons.GetHorizonClient()},
	// 	build.AddSigner(cd.RegistrarAccount.SignerKeys[0].Publickey, weight1), //TA
	// 	build.AddSigner(cd.RegistrarAccount.SignerKeys[1].Publickey, weight2), //QC
	// 	build.AddSigner(cd.RegistrarAccount.SignerKeys[2].Publickey, weight3), //QC
	// 	build.AddSigner(cd.RegistrarAccount.SignerKeys[3].Publickey, weight4), //QC
	// 	build.SetLowThreshold(low),
	// 	build.SetMediumThreshold(medium),
	// 	build.SetHighThreshold(high),
	// 	// build.MasterWeight(1),
	// )
	tx, err := build.Transaction(muts...)

	if err != nil {
		log.Fatal(err)
		return "value convertion error"
	}

	txe, err := tx.Sign(cd.RegistrarAccount.SignerKey)
	// if err != nil {
	// 	panic(err)
	// }

	// // Sign the transaction to prove you are actually the person sending it.
	// txe, err := tx.Sign(cd.RegistrarAccount.SignerKey)
	if err != nil {
		panic(err)
	}

	txeB64, err := txe.Base64()
	if err != nil {
		panic(err)
	}

	// And finally, send it off to Stellar!
	resp, err := commons.GetHorizonClient().SubmitTransaction(txeB64)
	if err != nil {
		panic(err)
	}
	return resp.Hash
	// fmt.Println("Successful Transaction:")
	// fmt.Println("Ledger:", resp.Ledger)
	// fmt.Println("Hash:", resp.Hash)
}
