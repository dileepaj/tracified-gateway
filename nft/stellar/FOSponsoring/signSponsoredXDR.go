package fosponsoring

import (
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

func BuildSignedSponsoredXDR(payload model.TransactionData) (string, error) {
	// Replace 'transactionXDR'  actual XDR
	transactionXDR := payload.XDR

	object := dao.Connection{}
	//getting the issuer secret key
	data, errkey := object.GetIssuerSK(payload.AccountIssuer).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errkey != nil {
		logrus.Error(errkey)
		return "", errkey
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
			return "", err1
		}

		txe, val := transaction.Transaction()
		logrus.Info("The value of the transaction is: ", val)

		additionalSigner, errpair := keypair.Parse(myStellarSeed) //decryptSK
		if errpair != nil {
			return "", errpair
		}

		hashXDR, errhash := txe.Hash(commons.GetStellarNetwork())
		if errhash != nil {
			return "", errhash
		}

		signer, errsigner := additionalSigner.SignDecorated(hashXDR[:])
		if errsigner != nil {
			return "", errsigner
		}

		hint := additionalSigner.Hint()

		decoratedSignature := xdr.DecoratedSignature{
			Signature: signer.Signature,
			Hint:      hint,
		}

		txesignex, errx := txe.AddSignatureDecorated(decoratedSignature)
		if errx != nil {
			return "", errx
		}

		bs64, err64 := txesignex.Base64()
		if err64 != nil {
			return "", err64
		}
		return bs64, nil
	}

}
