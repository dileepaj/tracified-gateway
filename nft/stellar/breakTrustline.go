package stellar

import (
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

func BreakTrustline(payload model.TransactionDataBreakTrustline) (string, error) {
	// Replace 'transactionXDR'  actual XDR
	transactionXDR := payload.XDR

	myStellarSeed := commons.GoDotEnvVariable("SPONSORERSK")

	// Parse the transaction XDR
	transaction, err1 := txnbuild.TransactionFromXDR(transactionXDR)
	if err1 != nil {
		return "", err1
	}

	txe, val := transaction.Transaction()
	logrus.Info("The value of the transaction is: ", val)

	additionalSigner, errpair := keypair.Parse(myStellarSeed) //decryptSK
	if errpair != nil {
		logrus.Error("Failed to parse into keypair", errpair)
	}

	hashXDR, errhash := txe.Hash(network.TestNetworkPassphrase)
	if errhash != nil {
		logrus.Error(errhash)
	}

	signer, errsigner := additionalSigner.SignDecorated(hashXDR[:])
	if errsigner != nil {
		logrus.Error(errsigner)
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

	txesignex.ToXDR()

	respn, errsubmitting := commons.GetHorizonClient().SubmitTransaction(txesignex)
	if errsubmitting != nil {
		logrus.Error("Error submitting transaction:", errsubmitting)
	}
	return respn.Hash, nil

}
