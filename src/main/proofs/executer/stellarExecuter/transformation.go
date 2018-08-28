package stellarExecuter

import (
	"fmt"
	"log"
	"main/api/apiModel"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
)

// func TransformMerge(code1 string, limit1 string, code2 string, limit2 string, code3 string, limit3 string, code4 string, limit4 string, reciver1 string, reciver2 string) string {

// 	// router := routes.NewRouter()

// 	// log.Fatal(http.ListenAndServe(":8030", router))

// 	//RA
// 	issuerSeed := "GASTEFX4WMC7PN3WIJTHYDJHR3D4FXVTKR7JWMBL4OUYEMPQDNPGNOAG"

// 	// recipientSeed := "SCSLMLE2GMIWCIF6GX2KWS64GJUSUGUF3NKKKP7WMRWZQJEOLVWNTJY7" //factory
// 	// recipientSeed2 := "SD7JWOWBL4Y777WHHTTWTYCFJME7CEOUYS2OZEJTNN5OHN3BYQOBLXLV" //reg

// 	// issuer, err := keypair.Parse(issuerSeed)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	recipient, err := keypair.Parse(reciver1)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// recipient2, err := keypair.Parse(recipientSeed2)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// Second, the issuing account actually sends a payment using the asset
// 	paymentTx, err := build.Transaction(
// 		build.SourceAccount{recipient.Address()},
// 		build.TestNetwork,
// 		build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
// 		build.Payment(
// 			build.Destination{AddressOrSeed: reciver2},
// 			build.CreditAmount{code1, issuerSeed, limit1},
// 		), build.Payment(
// 			build.Destination{AddressOrSeed: reciver2},
// 			build.CreditAmount{code2, issuerSeed, limit2},
// 		), build.Payment(
// 			build.Destination{AddressOrSeed: reciver2},
// 			build.CreditAmount{code3, issuerSeed, limit3},
// 		),

// 		build.Payment(
// 			build.SourceAccount{reciver2},
// 			build.Destination{AddressOrSeed: recipient.Address()},
// 			build.CreditAmount{code4, issuerSeed, limit4},
// 		),
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	paymentTxe, err := paymentTx.Sign(reciver1)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	paymentTxeB64, err := paymentTxe.Base64()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return paymentTxeB64
// 	// fmt.Printf("tx base64: %s", paymentTxeB64)
// }

type ConcreteTransform struct {
	AssestTransfer apiModel.AssestTransfer
	// Code1  string
	// Limit1 string
	// Code2  string
	// Limit2 string
	// Code3  string
	// Limit3 string
	// Code4  string
	// Limit4 string

	// Reciver1 string
	// Reciver2 string
}

func (cd *ConcreteTransform) TransformMerge() string {

	// router := routes.NewRouter()

	// log.Fatal(http.ListenAndServe(":8030", router))

	//issuer
	issuerSeed := cd.AssestTransfer.Issuer
	Reciver2 := cd.AssestTransfer.Reciver
	recipientSeed := cd.AssestTransfer.Sender //factory
	// recipientSeed2 := reciver2 //reg

	// issuer, err := keypair.Parse(issuerSeed)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	recipient, err := keypair.Parse(recipientSeed)
	if err != nil {
		log.Fatal(err)
	}
	// recipient2, err := keypair.Parse(recipientSeed2)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// Second, the issuing account actually sends a payment using the asset

	muts := []build.TransactionMutator{
		build.SourceAccount{recipient.Address()},
		build.TestNetwork,
		build.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
	}
	arrSize := len(cd.AssestTransfer.Asset)
	fmt.Println(arrSize)
	for i := 0; i < arrSize; i++ {
		if i < arrSize-1 {
			ops := []build.TransactionMutator{
				build.Payment(
					build.SourceAccount{recipient.Address()},
					build.Destination{AddressOrSeed: Reciver2},
					build.CreditAmount{cd.AssestTransfer.Asset[i].Code, issuerSeed, cd.AssestTransfer.Asset[i].Limit},
				),
			}
			muts = append(muts, ops...)
		} else {
			ops := []build.TransactionMutator{
				build.Payment(
					build.SourceAccount{Reciver2},
					build.Destination{AddressOrSeed: recipient.Address()},
					build.CreditAmount{cd.AssestTransfer.Asset[i].Code, issuerSeed, cd.AssestTransfer.Asset[i].Limit},
				),
			}
			muts = append(muts, ops...)
		}

	}

	tx, err := build.Transaction(muts...)

	if err != nil {
		log.Fatal(err)
	}

	paymentTxe, err := tx.Sign(recipientSeed)

	if err != nil {
		log.Fatal(err)
	}
	paymentTxeB64, err := paymentTxe.Base64()
	if err != nil {
		log.Fatal(err)
	}
	return paymentTxeB64
	// fmt.Printf("tx base64: %s", paymentTxeB64)
}
