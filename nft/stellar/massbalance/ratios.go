package massbalance

import (
	"log"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func checkRange(number int, min int, max int) uint8 {
	if number >= min && number <= max {
		return 1
	} else {
		return 0
	}
}
func GenerateManageDataforAccount(accountweight uint8) txnbuild.ManageData {
	//accessLevelStatus := txnbuild.ManageData()
	if accountweight == 1 {
		accessLevelStatus := txnbuild.ManageData(
			txnbuild.ManageData{
				Name:  "Open account",
				Value: []byte(string(accountweight)),
			},
		)
		return accessLevelStatus
	} else {
		accessLevelStatus := txnbuild.ManageData(
			txnbuild.ManageData{
				Name:  "Account locked ",
				Value: []byte(string(accountweight)),
			},
		)
		return accessLevelStatus
	}
}
func CheckAcountStatus(user model.AccountCredentials, sourceWeight int) (string, error) {

	manageData := txnbuild.ManageData(
		txnbuild.ManageData{
			Name:  "Account Weight  ",
			Value: []byte(string(rune(sourceWeight))),
		},
	)
	request := horizonclient.AccountRequest{AccountID: user.PublicKey}
	account, err := commons.GetHorizonClient().AccountDetail(request)
	if err != nil {
		log.Println("Error accoured when getting account : ", err.Error())
		return "", err
	}
	tx, txErr := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &account,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&manageData},
			BaseFee:              txnbuild.MinBaseFee,
			Memo:                 nil,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
		},
	)
	if txErr != nil {
		log.Println("Error creating transaction : ", txErr.Error())
		return "", txErr
	}

	senderKeypair, _ := keypair.ParseFull(user.SecretKey)

	txe64, signErr := tx.Sign(network.TestNetworkPassphrase, senderKeypair)
	if signErr != nil {
		log.Println("failed to sign : ", signErr.Error())
		return "", signErr
	}
	respn, reserr := commons.GetHorizonClient().SubmitTransaction(txe64)
	if reserr != nil {
		log.Println("Error submitting Manage data transaction:", reserr)
		return "locked", reserr
	}

	return respn.Hash, nil
}
func SetAccountLockLevel(productName string, userinput int, lowerlimit int, higherlimit int, singer model.AccountCredentials, user model.AccountCredentials) (string, string, error) {
	accountWeight := checkRange(userinput, lowerlimit, higherlimit)
	accountSetOptions := txnbuild.SetOptions{
		MasterWeight: txnbuild.NewThreshold(
			txnbuild.Threshold(accountWeight),
		),

		Signer: &txnbuild.Signer{
			Address: singer.PublicKey,
			Weight:  1,
		},
	}

	accountManageData := GenerateManageDataforAccount(accountWeight)
	request := horizonclient.AccountRequest{AccountID: user.PublicKey}
	account, err := commons.GetHorizonClient().AccountDetail(request)
	if err != nil {
		log.Println("Error accoured when getting account : ", err.Error())
		return "", "", err
	}

	tx, txErr := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &account,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&accountSetOptions, &accountManageData},
			BaseFee:              txnbuild.MinBaseFee,
			Memo:                 nil,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
		},
	)
	if txErr != nil {
		log.Println("Error creating transaction : ", txErr.Error())
		return "", "", txErr
	}

	senderKeypair, _ := keypair.ParseFull(user.SecretKey)

	txe64, signErr := tx.Sign(network.TestNetworkPassphrase, senderKeypair)
	if signErr != nil {
		log.Println("failed to sign : ", signErr.Error())
		return "", "", signErr
	}
	respn, reserr := commons.GetHorizonClient().SubmitTransaction(txe64)

	manageDataresp, manageDataErr := CheckAcountStatus(user, int(accountWeight))
	if manageDataErr != nil {
		log.Println("Manage data failed account locked")
		return respn.Hash, "locked", manageDataErr
	}

	if reserr != nil {
		log.Println("Error submitting transaction:", reserr)
		return respn.Hash, "locked", reserr
	}
	return respn.Hash, manageDataresp, nil
}

func UnlockAccount(singer model.AccountCredentials, user model.AccountCredentials) (string, error) {
	accountSetOptions := txnbuild.SetOptions{
		MasterWeight: txnbuild.NewThreshold(
			txnbuild.Threshold(1),
		),
	}
	request := horizonclient.AccountRequest{AccountID: user.PublicKey}
	account, err := commons.GetHorizonClient().AccountDetail(request)
	if err != nil {
		log.Println("Error accoured when getting account : ", err.Error())
		return "", err
	}

	tx, txErr := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &account,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&accountSetOptions},
			BaseFee:              txnbuild.MinBaseFee,
			Memo:                 nil,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
		},
	)
	if txErr != nil {
		log.Println("Error creating transaction : ", txErr.Error())
		return "", txErr
	}

	senderKeypair, _ := keypair.ParseFull(singer.SecretKey)

	txe64, signErr := tx.Sign(network.TestNetworkPassphrase, senderKeypair)
	if signErr != nil {
		log.Println("failed to sign : ", signErr.Error())
		return "", signErr
	}
	respn, reserr := commons.GetHorizonClient().SubmitTransaction(txe64)
	if reserr != nil {
		log.Println("Error submitting transaction:", reserr)
		return "account locked", reserr
	}
	return respn.Hash, nil
}
