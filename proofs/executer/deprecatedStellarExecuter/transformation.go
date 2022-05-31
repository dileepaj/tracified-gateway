package deprecatedStellarExecuter

import (
	"fmt"
	"log"

	"github.com/dileepaj/tracified-gateway/api/apiModel"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
)

func TransformMerge(code1 string, limit1 string, code2 string, limit2 string, code3 string, limit3 string, code4 string, limit4 string, reciver1 string, reciver2 string) string {
	// router := routes.NewRouter()

	// log.Fatal(http.ListenAndServe(":8030", router))

	// RA
	issuerSeed := "GASTEFX4WMC7PN3WIJTHYDJHR3D4FXVTKR7JWMBL4OUYEMPQDNPGNOAG"

	// recipientSeed := "SCSLMLE2GMIWCIF6GX2KWS64GJUSUGUF3NKKKP7WMRWZQJEOLVWNTJY7" //factory
	// recipientSeed2 := "SD7JWOWBL4Y777WHHTTWTYCFJME7CEOUYS2OZEJTNN5OHN3BYQOBLXLV" //reg

	// issuer, err := keypair.Parse(issuerSeed)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	recipient, err := keypair.Parse(reciver1)
	if err != nil {
		log.Fatal(err)
	}
	// recipient2, err := keypair.Parse(recipientSeed2)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	netClient := commons.GetHorizonClient()
	accountRequest := horizonclient.AccountRequest{AccountID: recipient.Address()}
	account, err := netClient.AccountDetail(accountRequest)
	if err != nil {
		// log.Fatal(err)
	}

	asset1 := txnbuild.CreditAsset{Code: code1, Issuer: issuerSeed}
	paymentTXNBuilder1 := txnbuild.Payment{Asset: asset1, Destination: reciver2, SourceAccount: issuerSeed, Amount: limit1}

	asset2 := txnbuild.CreditAsset{Code: code2, Issuer: issuerSeed}
	paymentTXNBuilder2 := txnbuild.Payment{Asset: asset2, Destination: reciver2, SourceAccount: issuerSeed, Amount: limit2}

	asset3 := txnbuild.CreditAsset{Code: code3, Issuer: issuerSeed}
	paymentTXNBuilder3 := txnbuild.Payment{Asset: asset3, Destination: reciver2, SourceAccount: issuerSeed, Amount: limit3}

	asset4 := txnbuild.CreditAsset{Code: code4, Issuer: issuerSeed}
	paymentTXNBuilder4 := txnbuild.Payment{Asset: asset4, Destination: reciver2, SourceAccount: issuerSeed, Amount: limit4}
	// BUILD THE GATEWAY XDR
	paymentTx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount:        &account,
		IncrementSequenceNum: true,
		Operations:           []txnbuild.Operation{&paymentTXNBuilder1, &paymentTXNBuilder2, &paymentTXNBuilder3, &paymentTXNBuilder4},
		Memo:                 nil,
		Preconditions:        txnbuild.Preconditions{},
	})
	// // Second, the issuing account actually sends a payment using the asset
	// paymentTx, err := build.Transaction(
	// 	build.SourceAccount{recipient.Address()},
	// 	commons.GetHorizonNetwork(),
	// 	build.AutoSequence{SequenceProvider: commons.GetHorizonClient()},
	// 	build.Payment(
	// 		build.Destination{AddressOrSeed: reciver2},
	// 		build.CreditAmount{code1, issuerSeed, limit1},
	// 	), build.Payment(
	// 		build.Destination{AddressOrSeed: reciver2},
	// 		build.CreditAmount{code2, issuerSeed, limit2},
	// 	), build.Payment(
	// 		build.Destination{AddressOrSeed: reciver2},
	// 		build.CreditAmount{code3, issuerSeed, limit3},
	// 	),
	// 	build.Payment(
	// 		build.SourceAccount{reciver2},
	// 		build.Destination{AddressOrSeed: recipient.Address()},
	// 		build.CreditAmount{code4, issuerSeed, limit4},
	// 	),
	// )
	if err != nil {
		log.Fatal(err)
	}
	paymentTxe, err := paymentTx.Sign(reciver1)
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

type ConcreteTransform struct {
	AssetTransfer apiModel.AssetTransfer
	ProfileID     string
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

	// issuer
	issuerSeed := cd.AssetTransfer.Issuer
	Reciver2 := cd.AssetTransfer.Reciver
	recipientSeed := cd.AssetTransfer.Sender // factory
	// recipientSeed2 := reciver2 //reg

	// issuer, err := keypair.Parse(issuerSeed)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	recipient, err := keypair.Parse(recipientSeed)
	if err != nil {
		// log.Fatal(err)
		return "Account not found"
	}
	// recipient2, err := keypair.Parse(recipientSeed2)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// Second, the issuing account actually sends a payment using the asset


	netClient := commons.GetHorizonClient()
	accountRequest := horizonclient.AccountRequest{AccountID: recipient.Address()}
	account, err := netClient.AccountDetail(accountRequest)
	if err != nil {
		// log.Fatal(err)
	}


	// BUILD THE GATEWAY XDR
	muts:= []txnbuild.Operation{}

	opsType := txnbuild.ManageData{Name:"Transaction Type",Value: []byte(cd.AssetTransfer.Type)}
	muts = append(muts, &opsType)
	opsTxn  := txnbuild.ManageData{Name:"PreviousTXNID",Value: []byte(cd.AssetTransfer.Type)}
	muts = append(muts, &opsTxn)
	opsProfile  := txnbuild.ManageData{Name:"ProfileID",Value: []byte(cd.ProfileID)}
	muts = append(muts, &opsProfile)

	arrSize := len(cd.AssetTransfer.Asset)
	fmt.Println(arrSize)

	for i := 0; i < arrSize; i++ {
		if i < arrSize-1 {
			asset1 := txnbuild.CreditAsset{Code: cd.AssetTransfer.Asset[i].Code, Issuer: issuerSeed}
			ops := txnbuild.Payment{Asset: asset1, Destination: Reciver2, SourceAccount: issuerSeed, Amount:cd.AssetTransfer.Asset[i].Limit}
			muts = append(muts, &ops)
		} else {
			asset1 := txnbuild.CreditAsset{Code: cd.AssetTransfer.Asset[i].Code, Issuer: issuerSeed}
			ops := txnbuild.Payment{Asset: asset1, Destination:  recipient.Address(), SourceAccount: issuerSeed, Amount:cd.AssetTransfer.Asset[i].Limit}
			muts = append(muts, &ops)
		}
	}

	// muts := []build.TransactionMutator{
	// 	build.SourceAccount{recipient.Address()},
	// 	commons.GetHorizonNetwork(),
	// 	build.AutoSequence{SequenceProvider: commons.GetHorizonClient()},
	// }
	// opsType := []build.TransactionMutator{
	// 	build.SetData("Transaction Type", []byte(cd.AssetTransfer.Type)),
	// }
	// muts = append(muts, opsType...)
	// opsTxn := []build.TransactionMutator{
	// 	build.SetData("PreviousTXNID", []byte(cd.ProfileID)),
	// }
	// muts = append(muts, opsTxn...)
	// opsProfile := []build.TransactionMutator{
	// 	build.SetData("ProfileID", []byte(cd.ProfileID)),
	// }
	// muts = append(muts, opsProfile...)

	// arrSize := len(cd.AssetTransfer.Asset)
	// fmt.Println(arrSize)
	// for i := 0; i < arrSize; i++ {
	// 	if i < arrSize-1 {
	// 		ops := []build.TransactionMutator{
	// 			build.Payment(
	// 				build.SourceAccount{recipient.Address()},
	// 				build.Destination{AddressOrSeed: Reciver2},
	// 				build.CreditAmount{cd.AssetTransfer.Asset[i].Code, issuerSeed, cd.AssetTransfer.Asset[i].Limit},
	// 			),
	// 		}
	// 		muts = append(muts, ops...)
	// 	} else {
	// 		ops := []build.TransactionMutator{
	// 			build.Payment(
	// 				build.SourceAccount{Reciver2},
	// 				build.Destination{AddressOrSeed: recipient.Address()},
	// 				build.CreditAmount{cd.AssetTransfer.Asset[i].Code, issuerSeed, cd.AssetTransfer.Asset[i].Limit},
	// 			),
	// 		}
	// 		muts = append(muts, ops...)
	// 	}
	// }

	tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount:        &account,
		IncrementSequenceNum: true,
		Operations:           muts,
		Memo:                 nil,
		Preconditions:        txnbuild.Preconditions{},
	})

	if err != nil {
		log.Fatal(err)
		return "Transaction Build issue"
	}

	paymentTxe, err := tx.Sign(recipientSeed)
	if err != nil {
		log.Fatal(err)
		return "Transaction Signed issue"
	}
	paymentTxeB64, err := paymentTxe.Base64()
	if err != nil {

		log.Fatal(err)
		return "Transaction failed"
	}
	return paymentTxeB64
	// fmt.Printf("tx base64: %s", paymentTxeB64)
}
