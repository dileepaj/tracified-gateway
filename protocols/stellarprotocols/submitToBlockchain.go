package stellarprotocols

import (
	"net/http"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
)

type StellarTransaction struct {
	Operations []txnbuild.Operation
	Memo       string
	Type       string
}

/*
des - common method to send a transaction to the blockchain
*/

func (transaction StellarTransaction) SubmitToStellarBlockchain() (error, int, string, int64, string, string) {
	// load account
	publicKey := commons.GoDotEnvVariable("SOCILAIMPACTPUBLICKKEY")
	secretKey := commons.GoDotEnvVariable("SOCILAIMPACTSEED")
	tracifiedAccount, err := keypair.ParseFull(secretKey)
	client := commons.GetHorizonClient()
	pubaccountRequest := horizonclient.AccountRequest{AccountID: publicKey}
	pubaccount, err := client.AccountDetail(pubaccountRequest)
	// BUILD THE GATEWAY XDR
	tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount:        &pubaccount,
		IncrementSequenceNum: true,
		Operations:           transaction.Operations,
		BaseFee:              constants.MinBaseFee,
		Memo:                 txnbuild.MemoText(transaction.Memo),
		Preconditions:        txnbuild.Preconditions{TimeBounds:constants.TransactionTimeOut},
	})
	if err != nil {
		logrus.Println("Error while building XDR " + err.Error())
		return err, http.StatusInternalServerError, "", 0, "", publicKey
	}
	// SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
	gatewayTXE, err := tx.Sign(commons.GetStellarNetwork(), tracifiedAccount)
	xdrBase64, err := gatewayTXE.Base64()
	if err != nil {
		logrus.Error("Error while signing the XDR by secretKey  ", err)
		return err, http.StatusInternalServerError, "", 0, xdrBase64, publicKey
	}
	// CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
	resp, err := client.SubmitTransaction(gatewayTXE)
	if err != nil {
		logrus.Error("XDR submitting issue  ", err)
		return err, http.StatusInternalServerError, "", 0, xdrBase64, publicKey
	}
	return nil, 200, resp.Hash, resp.AccountSequence, xdrBase64, publicKey
}
