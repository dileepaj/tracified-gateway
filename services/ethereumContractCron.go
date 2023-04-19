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
	log.Debug("Ethereum cron job started")
	cronJobStartTime := time.Now().String()

	object := dao.Connection{}
	var formulaObj model.EthereumExpertFormula
	var metricObj model.EthereumMetricBind
	//Get the transactions with the pending status from the Database
	p := object.GetPendingContractsByStatus(117) // 117 is the status code for PENDING transactions
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
			if result[i].ContractType == "ETHEXPERTFORMULA" {
				// get the formula by uuid
				formulaDetails, errorInGettingFormulaDetails := dbCollectionHandler.GetEthFormulaByUUID(result[i].Identifier)
				if errorInGettingFormulaDetails != nil {
					logrus.Error("Error when getting the formula details : " + errorInGettingFormulaDetails.Error())
					continue
				} 

				formulaObj = formulaDetails
			} else if result[i].ContractType == "ETHMETRICBIND" {
				// get the metric by uuid
				metricDetails, errorInGettingMetricDetails := dbCollectionHandler.GetEthMetricByUUID(result[i].Identifier)
				if errorInGettingMetricDetails != nil {
					logrus.Error("Error when getting the metric details : " + errorInGettingMetricDetails.Error())
					continue
				}
				metricObj = metricDetails
				//check the time difference between the current time and the time of the transaction and if it is less than 10 minutes, skip the transaction  
				givenTimestamp := metricDetails.Timestamp
				layout := "2006-01-02 15:04:05"
				// truncate the timestamp to the layout
				truncatedTime, errInParsingAndTruncating := time.Parse(layout, givenTimestamp[:len(layout)])
				if errInParsingAndTruncating != nil {
					logrus.Error("Error when parsing and truncating the time : " + errInParsingAndTruncating.Error())
					continue
				}
				// get the current time
				currentTime := time.Now().UTC()
				// get the difference between the current time and the time of the transaction
				timeDifference := currentTime.Sub(truncatedTime)
				if timeDifference.Round(time.Second).Minutes() < 10 {
					// logrus.Info("Transaction " + result[i].TransactionHash + " is less than 10 minutes old, skipping the transaction")
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
					Status:          120,
					CurrentIndex:    result[i].CurrentIndex + 1,
					ErrorMessage:    "Pending checking capacity met",
					ContractType:    result[i].ContractType,
					Identifier:      result[i].Identifier,
					Nonce:           result[i].Nonce,
					GasPrice:        result[i].GasPrice,
					GasLimit:        result[i].GasLimit,
				}
				// updating actual status in the database
				if result[i].ContractType == "ETHEXPERTFORMULA" {
					formulaObj.ActualStatus = 115 	// DEPLOYMENT_TRANSACTION_CANCELLED
					errWhenUpdatingActualStatus := object.UpdateSelectedEthFormulaFields(formulaObj.FormulaID, formulaObj.TransactionUUID, formulaObj)
					if errWhenUpdatingActualStatus != nil {
						logrus.Error("Error when updating the actual status of the formula : " + errWhenUpdatingActualStatus.Error())
						continue
					}
				} else if result[i].ContractType == "ETHMETRICBIND" {
					metricObj.ActualStatus = 115 	// DEPLOYMENT_TRANSACTION_CANCELLED
					errWhenUpdatingActualStatus := object.UpdateSelectedEthMetricFields(metricObj.MetricID, metricObj.TransactionUUID, metricObj)
					if errWhenUpdatingActualStatus != nil {
						logrus.Error("Error when updating the actual status of the metric : " + errWhenUpdatingActualStatus.Error())
						continue
					}
				}
				
				updateCancel.ErrorMessage = "Transaction pending checking capacity met"
				if result[i].ContractType == "ETHMETRICBIND" {
					errorWhenInvalidatingTransactions := dbCollectionHandler.InvalidateMetric(updateCancel, updateCancel.Status, updateCancel.ErrorMessage)
					if errorWhenInvalidatingTransactions != nil {
						logrus.Error("Error when invalidating the transaction : " + errorWhenInvalidatingTransactions.Error())
						continue
					}
				} else if result[i].ContractType == "ETHEXPERTFORMULA" {
					errWhenUpdatingStatus := dbCollectionHandler.UpdateCollectionsWithNewStatus(updateCancel, 120)
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
					result[i].Status = 118	// SUCCESS
					// use collection update strategy
					collectionUpdater := &ethereumservices.CollectionServiceContext{}
					collectionUpdater.SetCollectionServiceStrategy(&ethereumservices.SocialImpactMainCollectionUpdate{
						PendingContract: result[i],
						Status:       118,	// SUCCESS
					})
					errInUpdatingDBForSuccessfulTransactions := collectionUpdater.ExecuteCollectionService()
					if errInUpdatingDBForSuccessfulTransactions != nil {
						logrus.Error("Error when updating the database for successful transactions : " + errInUpdatingDBForSuccessfulTransactions.Error())
						continue
					}

					// updating actual status in the database
					if result[i].ContractType == "ETHEXPERTFORMULA" {
						formulaObj.Status = 118 // SUCCESS
						formulaObj.ActualStatus = 113	// DEPLOYMENT_TRANSACTION_SUCCESS
						errWhenUpdatingActualStatus := object.UpdateSelectedEthFormulaFields(formulaObj.FormulaID, formulaObj.TransactionUUID, formulaObj)
						if errWhenUpdatingActualStatus != nil {
							logrus.Error("Error when updating the actual status of the formula : " + errWhenUpdatingActualStatus.Error())
							continue
						}
					} else if result[i].ContractType == "ETHMETRICBIND" {
						metricObj.Status = 118 // SUCCESS
						metricObj.ActualStatus = 113 	// DEPLOYMENT_TRANSACTION_SUCCESS
						errWhenUpdatingActualStatus := object.UpdateSelectedEthMetricFields(metricObj.MetricID, metricObj.TransactionUUID, metricObj)
						if errWhenUpdatingActualStatus != nil {
							logrus.Error("Error when updating the actual status of the metric : " + errWhenUpdatingActualStatus.Error())
							continue
						}
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

					result[i].Status = 119	// FAILED
					result[i].ErrorMessage = errorOccurred

					if result[i].ContractType == "ETHMETRICBIND" {
						errWhenUpdatingCollection := dbCollectionHandler.InvalidateMetric(result[i], 119, result[i].ErrorMessage)
						if errWhenUpdatingCollection != nil {
							logrus.Error("Error when updating the collection : " + errWhenUpdatingCollection.Error())
							continue
						}
						// updating status and actual status in the database
						metricObj.Status = 119 // FAILED
						metricObj.ActualStatus = 114 	// DEPLOYMENT_TRANSACTION_FAILED
						errWhenUpdatingActualStatus := object.UpdateSelectedEthMetricFields(metricObj.MetricID, metricObj.TransactionUUID, metricObj)
						if errWhenUpdatingActualStatus != nil {
							logrus.Error("Error when updating the actual status of the metric : " + errWhenUpdatingActualStatus.Error())
							continue
						}
					} else if result[i].ContractType == "ETHEXPERTFORMULA" {
						errWhenUpdatingCollection := dbCollectionHandler.UpdateCollectionsWithNewStatus(result[i], 119)
						if errWhenUpdatingCollection != nil {
							logrus.Error("Error when updating the collection : " + errWhenUpdatingCollection.Error())
							continue
						}
						// updating status and actual status in the database
						formulaObj.Status = 119 // FAILED
						formulaObj.ActualStatus = 114	// DEPLOYMENT_TRANSACTION_FAILED
						errWhenUpdatingActualStatus := object.UpdateSelectedEthFormulaFields(formulaObj.FormulaID, formulaObj.TransactionUUID, formulaObj)
						if errWhenUpdatingActualStatus != nil {
							logrus.Error("Error when updating the actual status of the formula : " + errWhenUpdatingActualStatus.Error())
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
								Status:          0,
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
								Status:          120,	// CANCELLED
								CurrentIndex:    result[i].CurrentIndex + 1,
								ErrorMessage:    errWhenRedeploying.Error(),
								ContractType:    result[i].ContractType,
								Identifier:      result[i].Identifier,
								Nonce:           nonce,
								GasPrice:        gasPrice,
								GasLimit:        gasLimit,
							}
							errWhenUpdatingStatus := dbCollectionHandler.UpdateCollectionsWithNewStatus(updatePending, 120)
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
		log.Debug("Ethereum cron job completed at " + time.Now().String() + "(for started time : " + cronJobStartTime + ")")
		return nil
	}).Catch(func(error error) error {
		if commons.GoDotEnvVariable("LOGSTYPE") == "DEBUG" {
			log.Error("Error @CheckContractStatus " + error.Error())
		}
		return error
	})
	p.Await()
}
