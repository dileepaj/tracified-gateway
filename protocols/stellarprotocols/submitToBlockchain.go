package stellarprotocols

import (
	"net/http"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
)

type StellarTrasaction struct {
	PublicKey  string
	Operations []txnbuild.Operation
	Memo       string
	Type string
	TransactionHash string
	Sequence        int64
	XDR             string
	FormulaId       string
	MetricId        string
	User            model.User
}

/*
des - common method to send a transaction to the blockchain
*/

func (transaction StellarTrasaction) SubmitToStellerBlockchain() (error, int, string, string) {
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
		BaseFee:              txnbuild.MinBaseFee,
		Memo:                 txnbuild.MemoText(transaction.Memo),
		Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
	})
	if err != nil {
		logrus.Println("Error while buliding XDR " + err.Error())
	}
	// SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
	GatewayTXE, err := tx.Sign(commons.GetStellarNetwork(), tracifiedAccount)
	if err != nil {
		logrus.Error("Error while signing the XDR by secretKey  ", err)
		return err, http.StatusInternalServerError, "", GatewayTXE.ToXDR().GoString()
	}
	// CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
	resp, err := client.SubmitTransaction(GatewayTXE)
	if err != nil {
		logrus.Error("XDR submitting issue  ", err)
		return err, http.StatusInternalServerError, "", GatewayTXE.ToXDR().GoString()
	}
	return nil, 200, resp.Hash, GatewayTXE.ToXDR().GoString()
}
