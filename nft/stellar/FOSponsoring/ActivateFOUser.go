package fosponsoring

import (
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

func ActivateFOUser(xdrs string) (string, error) {

	transactionx, err1 := txnbuild.TransactionFromXDR(xdrs)
	if err1 != nil {
		return "", err1
	}

	txes, vals := transactionx.Transaction()
	logrus.Info("value to show the GT can be packed is ", vals)

	additionalSigners, err2 := keypair.Parse(commons.GoDotEnvVariable("SPONSORERSK")) //decryptSK
	if err2 != nil {
		logrus.Error(err2)
	}

	hashXDRs, err3 := txes.Hash(network.TestNetworkPassphrase)
	if err3 != nil {
		logrus.Error(err3)
	}

	signers, err4 := additionalSigners.SignDecorated(hashXDRs[:])
	if err4 != nil {
		logrus.Error(err4)
	}

	hints := additionalSigners.Hint()

	decoratedSignatures := xdr.DecoratedSignature{
		Signature: signers.Signature,
		Hint:      hints,
	}

	txesignexs, err5 := txes.AddSignatureDecorated(decoratedSignatures)
	if err5 != nil {
		return "", err5
	}

	txesignexs.ToXDR()
	bs64xdrs, errsignex := txesignexs.Base64()
	if errsignex != nil {
		return "", errsignex
	}
	logrus.Info("xdr signed base 64: ", bs64xdrs)

	respns, errsubmitting := commons.GetHorizonClient().SubmitTransaction(txesignexs)
	if errsubmitting != nil {
		logrus.Error("Error submitting transaction:", errsubmitting)
	}
	return respns.Hash, nil
}
