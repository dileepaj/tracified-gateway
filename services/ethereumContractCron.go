package services

import (
	"context"
	"strconv"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/ethereum/deploy"
	ethereumservices "github.com/dileepaj/tracified-gateway/services/ethereumServices"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/sirupsen/logrus"
)

func CheckContractStatus() {
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
		pendingCap, errWhenConvertingPendingCap := strconv.Atoi(commons.GoDotEnvVariable("PENDINGTHRESHOLD"))
		if errWhenConvertingPendingCap != nil {
			logrus.Error("Error when converting the pending cap : " + errWhenConvertingPendingCap.Error())
			return nil
		}
		for i := 0; i < len(result); i++ {
			pendingHash := result[i].TransactionHash
			//check the pending threshold
			if result[i].CurrentIndex == pendingCap {
				logrus.Error(pendingHash + " hash pending checking capacity met, invalidating transaction")
				updateCancel := model.PendingContracts{
					TransactionHash: result[i].TransactionHash,
					ContractAddress: result[i].ContractAddress,
					Status:          "Cancelled",
					CurrentIndex:    result[i].CurrentIndex + 1,
					ErrorMessage:    result[i].ErrorMessage,
					ContractType:    result[i].ContractType,
					Identifier:      result[i].Identifier,
					Nonce:           result[i].Nonce,
					GasPrice:        result[i].GasPrice,
					GasLimit:        result[i].GasLimit,
				}
				errWhenUpdatingStatus := object.UpdateEthereumPendingContract(result[i].TransactionHash, result[i].ContractAddress, result[i].Identifier, updateCancel)
				if errWhenUpdatingStatus != nil {
					logrus.Error("Error when updating status of the transaction : " + errWhenUpdatingStatus.Error())
					continue
				}
				continue
			} else {
				//check the transaction status
				transactionReceipt, errWhenTakingTheReceipt := ethClient.TransactionReceipt(context.Background(), common.HexToHash(pendingHash))
				if errWhenTakingTheReceipt != nil {
					//check the error
					if errWhenTakingTheReceipt.Error() == "not found" {
						//transaction is still pending
						logrus.Info(pendingHash + " transaction is still in pending state.")
						updatePending := model.PendingContracts{
							TransactionHash: result[i].TransactionHash,
							ContractAddress: result[i].ContractAddress,
							Status:          "Pending",
							CurrentIndex:    result[i].CurrentIndex + 1,
							ErrorMessage:    result[i].ErrorMessage,
							ContractType:    result[i].ContractType,
							Identifier:      result[i].Identifier,
							Nonce:           result[i].Nonce,
							GasPrice:        result[i].GasPrice,
							GasLimit:        result[i].GasLimit,
						}

						errWhenUpdatingStatus := object.UpdateEthereumPendingContract(result[i].TransactionHash, result[i].ContractAddress, result[i].Identifier, updatePending)
						if errWhenUpdatingStatus != nil {
							logrus.Error("Error when updating status of the transaction : " + errWhenUpdatingStatus.Error())
							continue
						}
					} else {
						logrus.Error("Error when calling the transaction receipt : " + errWhenTakingTheReceipt.Error())
						continue
					}

				} else if transactionReceipt.Status == 1 {
					//Transaction is successful
					// update both collections
					abstractObj := ethereumservices.AbstractCollectionUpdate{
						PendingContract: result[i],
						Status:          "SUCCESS",
						Type: 			 "SOCIALIMPACT",
					}
					errInUpdatingDBForSuccessfulTransactions := abstractObj.AbstractCollectionUpdater()
					if errInUpdatingDBForSuccessfulTransactions != nil {
						logrus.Error("Error when updating the database for successful transactions : " + errInUpdatingDBForSuccessfulTransactions.Error())
						continue
					}
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
