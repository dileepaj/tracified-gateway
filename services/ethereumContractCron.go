package services

import (
	"context"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/ethereum/deploy"
	"github.com/dileepaj/tracified-gateway/vendor/github.com/ethereum/go-ethereum/ethclient"
	"github.com/dileepaj/tracified-gateway/vendor/github.com/ethereum/go-ethereum/log"
	"github.com/dileepaj/tracified-gateway/vendor/github.com/sirupsen/logrus"
	"github.com/ethereum/go-ethereum/common"
)

func CheckContractStatus() {
	//TODO loop through the transactions and call the ethereum endpoint to check the transaction status
	//TODO check the status
	/*
		TODO
			if pending - Update the index
			if success - update in the DB collection as completed
			if failed - log the error and mark the status as failed
	*/

	if commons.GoDotEnvVariable("LOGSTYPE") == "DEBUG" {
		log.Debug("---------------------------------------- Check pending Ethereum contracts -----------------------")
	}

	object := dao.Connection{}
	//Get the transactions with the pending status from the Database
	p := object.GetPendingContractsByStatus("Pending")
	p.Then(func(data interface{}) interface{} {
		result := data.([]model.PendingContracts)
		ethClient, errWHenDialingEthClient := ethclient.Dial(commons.GoDotEnvVariable("ETHEREUMTESTNETLINK"))
		if errWHenDialingEthClient != nil {
			logrus.Error("Error when calling th Ethereum client on Cron job, Error : " + errWHenDialingEthClient.Error())
			return nil
		}
		for i := 0; i < len(result); i++ {
			//check the transaction status
			pendingHash := result[i].TransactionHash
			transactionReceipt, errWhenTakingTheReceipt := ethClient.TransactionReceipt(context.Background(), common.HexToHash(pendingHash))
			if errWhenTakingTheReceipt != nil {
				//Pending transaction
			} else if transactionReceipt.Status == 1 {
				//Transaction is successful
			} else if transactionReceipt.Status == 0 {
				//Transaction failed
				//Get the error for the transaction
				errorOccurred, errWhenGettingTheTransactionError := deploy.GetErrorOfFailedTransaction(pendingHash)
				if errWhenGettingTheTransactionError != nil {
					logrus.Error("Error when getting the transaction error : " + errWhenGettingTheTransactionError.Error())
					continue
				}
				logrus.Info(pendingHash + " hash failed due to the error : " + errorOccurred)

				//call the failed contact redeployer

			} else {
				logrus.Error("Invalid transaction receipt status for transaction hash : ", pendingHash)
				continue
			}
		}
		return nil
	}).Catch(func(error error) error {
		if commons.GoDotEnvVariable("LOGSTYPE") == "DEBUG" {
			log.Error("Error @CheckContractStatus " + error.Error())
		}
		return error
	})
	p.Await()
}
