package services

import (
	"context"
	"strconv"
	"time"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	ethereumservices "github.com/dileepaj/tracified-gateway/services/ethereumServices"
	"github.com/dileepaj/tracified-gateway/services/ethereumServices/dbCollectionHandler"
	"github.com/dileepaj/tracified-gateway/services/ethereumServices/pendingTransactionHandler"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/sirupsen/logrus"
)

func CheckContractStatus() {

	log.Debug("---------------------------------------- Check pending Ethereum contracts -----------------------")

	object := dao.Connection{}
	//Get the transactions with the pending status from the Database
	p := object.GetPendingContractsByStatus("PENDING")
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
			if result[i].ContractType == "ETHMETRICBIND" {
				// get the metric by uuid
				metricDetails, errorInGettingMetricDetails := dbCollectionHandler.GetEthMetricByUUID(result[i].Identifier)
				if errorInGettingMetricDetails != nil {
					logrus.Error("Error when getting the metric details : " + errorInGettingMetricDetails.Error())
					continue
				}
				//check the time difference between the current time and the time of the transaction and if it is less than 10 minutes, skip the transaction  
				givenTimestamp := metricDetails.Timestamp
				layout := "2006-01-02 15:04:05.9999999 -0700"
				// truncate the timestamp to the layout
				truncatedTime, errInParsingAndTruncating := time.Parse(layout, givenTimestamp[:len(layout)])
				if errInParsingAndTruncating != nil {
					logrus.Error("Error when parsing and truncating the time : " + errInParsingAndTruncating.Error())
					continue
				}
				// get the current time
				currentTime := time.Now()
				// get the difference between the current time and the time of the transaction
				timeDifference := currentTime.Sub(truncatedTime)
				if timeDifference.Abs().Minutes() < 10 {
					logrus.Info("Transaction " + result[i].TransactionHash + " is less than 10 minutes old, skipping the transaction")
					continue
				} 
			}

			pendingHash := result[i].TransactionHash
			//check the pending threshold
			if result[i].CurrentIndex == pendingCap {
				logrus.Error(pendingHash + " hash pending checking capacity met, invalidating transaction")
				updateCancel := model.PendingContracts{
					TransactionHash: result[i].TransactionHash,
					ContractAddress: result[i].ContractAddress,
					Status:          "CANCELLED",
					CurrentIndex:    result[i].CurrentIndex + 1,
					ErrorMessage:    "Pending checking capacity met",
					ContractType:    result[i].ContractType,
					Identifier:      result[i].Identifier,
					Nonce:           result[i].Nonce,
					GasPrice:        result[i].GasPrice,
					GasLimit:        result[i].GasLimit,
				}
				updateCancel.ErrorMessage = "Transaction pending checking capacity met"
				if result[i].ContractType == "ETHMETRICBIND" {
					errorWhenInvalidatingTransactions := dbCollectionHandler.InvalidateMetric(updateCancel, updateCancel.Status, updateCancel.ErrorMessage)
					if errorWhenInvalidatingTransactions != nil {
						logrus.Error("Error when invalidating the transaction : " + errorWhenInvalidatingTransactions.Error())
						continue
					}
				} else if result[i].ContractType == "ETHEXPERTFORMULA" {
					errWhenUpdatingStatus := dbCollectionHandler.UpdateCollectionsWithNewStatus(updateCancel, "CANCELLED")
					if errWhenUpdatingStatus != nil {
						logrus.Error("Error when updating status of the transaction : " + errWhenUpdatingStatus.Error())
						continue
					}
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
					} else {
						logrus.Error("Error when calling the transaction receipt : " + errWhenTakingTheReceipt.Error())
					}
					continue

				} else if transactionReceipt.Status == 1 {
					//Transaction is successful
					// update both collections
					result[i].Status = "SUCCESS"
					abstractObj := ethereumservices.AbstractCollectionUpdate{
						PendingContract: result[i],
						Status:          "SUCCESS",
						Type:            "SOCIALIMPACT",
					}
					errInUpdatingDBForSuccessfulTransactions := abstractObj.AbstractCollectionUpdater()
					if errInUpdatingDBForSuccessfulTransactions != nil {
						logrus.Error("Error when updating the database for successful transactions : " + errInUpdatingDBForSuccessfulTransactions.Error())
						continue
					}
					continue
				} else if transactionReceipt.Status == 0 {
					//Transaction failed
					//Get the error for the transaction
					errorOccurred, errWhenGettingTheTransactionError := pendingTransactionHandler.GetErrorOfFailedTransaction(pendingHash)
					if errWhenGettingTheTransactionError != nil {
						logrus.Error("Error when getting the transaction error : " + errWhenGettingTheTransactionError.Error())
						continue
					}
					logrus.Info(pendingHash + " hash failed due to the error : " + errorOccurred)

					// update error messages collection
					errorMessage := model.EthErrorMessage{
						TransactionHash: pendingHash,
						ErrorMessage:    errorOccurred,
					}

					errorWhenInsertingErrorMessage := object.InsertEthErrorMessage(errorMessage)
					if errorWhenInsertingErrorMessage != nil {
						logrus.Error("Error when inserting the error message : " + errorWhenInsertingErrorMessage.Error())
					}

					result[i].Status = "FAILED"
					result[i].ErrorMessage = errorOccurred

					if result[i].ContractType == "ETHMETRICBIND" {
						errWhenUpdatingCollection := dbCollectionHandler.InvalidateMetric(result[i], "FAILED", result[i].ErrorMessage)
						if errWhenUpdatingCollection != nil {
							logrus.Error("Error when updating the collection : " + errWhenUpdatingCollection.Error())
							continue
						}
					} else if result[i].ContractType == "ETHEXPERTFORMULA" {
						errWhenUpdatingCollection := dbCollectionHandler.UpdateCollectionsWithNewStatus(result[i], "FAILED")
						if errWhenUpdatingCollection != nil {
							logrus.Error("Error when updating the collection : " + errWhenUpdatingCollection.Error())
							continue
						}
					}

					//call the failed contact redeployer
					if result[i].ContractType == "ETHEXPERTFORMULA" {

						// use deployment strategy
						deployer := &ethereumservices.ContractDeployerContext{}
						deployer.SetContractDeploymentStrategy(&ethereumservices.AbstractContractRedeployment{
							PendingContract: model.PendingContracts{
								TransactionHash: "",
								ContractAddress: "",
								Status:          "",
								CurrentIndex:    result[i].CurrentIndex,
								ErrorMessage:    result[i].ErrorMessage,
								ContractType:    result[i].ContractType,
								Identifier:      result[i].Identifier,
								Nonce:    result[i].Nonce,
								GasPrice: result[i].GasPrice,
								GasLimit: result[i].GasLimit,
							},
						})
						contractAddress, transactionHash, _, nonce, gasPrice, gasLimit, errWhenRedeploying := deployer.ExecuteContractDeployment()
						if errWhenRedeploying != nil {
							logrus.Error("Error when redeploying the failed transaction : " + errWhenRedeploying.Error())
							//update collection
							updatePending := model.PendingContracts{
								TransactionHash: transactionHash,
								ContractAddress: contractAddress,
								Status:          "CANCELLED",
								CurrentIndex:    result[i].CurrentIndex + 1,
								ErrorMessage:    errWhenRedeploying.Error(),
								ContractType:    result[i].ContractType,
								Identifier:      result[i].Identifier,
								Nonce:           nonce,
								GasPrice:        gasPrice,
								GasLimit:        gasLimit,
							}
							errWhenUpdatingStatus := dbCollectionHandler.UpdateCollectionsWithNewStatus(updatePending, "CANCELLED")
							if errWhenUpdatingStatus != nil {
								logrus.Error("Error when updating status of the transaction : " + errWhenUpdatingStatus.Error())
								continue
							}

							continue
						}
					}

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
