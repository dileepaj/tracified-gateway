package fosponsoring

import (
	"fmt"
	"log"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

func BuildSignedSponsoredXDR(payload model.TransactionData) (string, error) {
	// Replace 'transactionXDR'  actual XDR
	transactionXDR := payload.XDR

	object := dao.Connection{}
	//getting the issuer secret key
	data, err := object.GetIssuerSK(payload.AccountIssuer).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err != nil {
		log.Println(err)
		return "", err
	} else {
		Keys := data.([]model.TransactionDataKeys)
		//decrypt the secret key
		decrpytNftissuerSecretKey := commons.Decrypt(Keys[0].AccountIssuerSK)
		if data == nil {
			logrus.Error("PublicKey is not found in gateway datastore")
		}
		myStellarSeed := decrpytNftissuerSecretKey

		// Parse the transaction XDR
		transaction, err1 := txnbuild.TransactionFromXDR(transactionXDR)
		if err1 != nil {
			fmt.Println("Error parsing transaction XDR:", err1)
			return "", err1
		}

		txe, val := transaction.Transaction()
		fmt.Println("value to show the GT can be packed is ", val)
		fmt.Println("New Transaction is ", txe)

		///////////////////////////
		additionalSigner, err := keypair.Parse(myStellarSeed) //decryptSK
		if err != nil {
			fmt.Println("------------rrrrrrrrr------------", err)
		}
		fmt.Println("key pait signer------------- : ", additionalSigner)

		hashXDR, err := txe.Hash(network.TestNetworkPassphrase)
		if err != nil {
			panic(err)
		}
		fmt.Println("hash is ", hashXDR)

		signer, err := additionalSigner.SignDecorated(hashXDR[:])
		if err != nil {
			panic(err)
		}
		fmt.Println("signer is ", signer)

		hint := additionalSigner.Hint()
		fmt.Println("hint ", hint)

		decoratedSignature := xdr.DecoratedSignature{
			Signature: signer.Signature,
			Hint:      hint,
		}

		txesignex, errx := txe.AddSignatureDecorated(decoratedSignature)
		if errx != nil {
			fmt.Println("Error parsing Stellar secret seed:", errx)
			return "", errx
		}
		fmt.Println("sign txe : ", txesignex)
		////////////////////////////////////////////////

		// txesigned, err3 := txe.SignWithKeyString(myStellarSeed)
		// if err3 != nil {
		// 	fmt.Println("Error parsing Stellar secret seed:", err3)
		// 	return
		// }

		// fmt.Println("signed transaction ", txesigned)
		xdrsignex := txesignex.ToXDR()
		bs64xdr, errsignex := txesignex.Base64()
		if errsignex != nil {
			fmt.Println("Error parsing Stellar base64:", errsignex)
			return "", errsignex
		}
		fmt.Println("xdr signed: ", xdrsignex)
		fmt.Println("xdr signed base 64: ", bs64xdr)

		// xdrs := txe.ToXDR()
		// bs64, err64 := txe.Base64()
		// if err64 != nil {
		// 	fmt.Println("Error parsing Stellar base64:", err64)
		// 	return
		// }

		// fmt.Println("signs are they present? : ", txe.Signatures())

		// fmt.Println("sign txe bs64: ", bs64)

		// fmt.Println("new xdr is : ", xdrs)
		// xdrstring := xdrs.GoString()
		// fmt.Println("new xdr string is : ", xdrstring)
		// fmt.Println("sign : ", xdrs)
		return bs64xdr, nil
	}

}
