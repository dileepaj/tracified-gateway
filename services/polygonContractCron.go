package services

import (
	"strconv"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/services/ethereumServices/dbCollectionHandler"
	"github.com/dileepaj/tracified-gateway/services/ethereumServices/pendingTransactionHandler"
	transactionrecipthandler "github.com/dileepaj/tracified-gateway/services/polygonServices/transactionReciptHandler"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/ethereum/go-ethereum/ethclient"
)

func CheckPolygonContractStatus() {
	object := dao.Connection{}
	var formulaObj model.EthereumExpertFormula
	logger := utilities.NewCustomLogger()
	p := object.GetPendingContractsByStatus(117, "PolygonPendingTransactions")
	p.Then(func(data interface{}) interface{} {
		result := data.([]model.PendingContracts)
		ethClient, errWHenDialingEthClient := ethclient.Dial(commons.GoDotEnvVariable("POLYGONALCHEMYTESTNETLINK") + commons.GoDotEnvVariable("POLYGONALCHEMYAPIKEY"))
		if errWHenDialingEthClient != nil {
			logger.LogWriter("Error when calling th Ethereum client on Cron job, Error : "+errWHenDialingEthClient.Error(), constants.ERROR)
			return nil
		}
		pendingCap, errWhenConvertingPendingCap := strconv.Atoi(commons.GoDotEnvVariable("PENDINGTHRESHOLD"))
		if errWhenConvertingPendingCap != nil {
			logger.LogWriter("Error when converting the pending cap : "+errWhenConvertingPendingCap.Error(), constants.ERROR)
			return nil
		}
		for i := 0; i < len(result); i++ {
			if result[i].ContractType == "POLYGONEXPERTFORMULA" {
				// get the formula by uuid
				formulaDetails, errorInGettingFormulaDetails := dbCollectionHandler.GetEthFormulaByUUID(result[i].Identifier, "PolygonExpertFormula")
				if errorInGettingFormulaDetails != nil {
					logger.LogWriter("Error when getting the formula details : "+errorInGettingFormulaDetails.Error(), constants.ERROR)
					continue
				}

				formulaObj = formulaDetails
			}
			//TODO - handle metric bind
			pendingHash := result[i].TransactionHash
			if result[i].CurrentIndex == pendingCap {
				logger.LogWriter(pendingHash+" hash pending checking capacity met, invalidating transaction", constants.ERROR)
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
				if result[i].ContractType == "POLYGONEXPERTFORMULA" {
					formulaObj.ActualStatus = 115 // DEPLOYMENT_TRANSACTION_CANCELLED
					errWhenUpdatingActualStatus := object.UpdateSelectedPolygonFormulaFields(formulaObj.FormulaID, formulaObj.TransactionUUID, formulaObj)
					if errWhenUpdatingActualStatus != nil {
						logger.LogWriter("Error when updating the actual status of the formula : "+errWhenUpdatingActualStatus.Error(), constants.ERROR)
						continue
					}
				}
				//TODO - Handle metric bind
				if result[i].ContractType == "POLYGONEXPERTFORMULA" {
					errWhenUpdatingStatus := dbCollectionHandler.UpdateCollectionsWithNewStatusForPolygon(updateCancel, 120)
					if errWhenUpdatingStatus != nil {
						logger.LogWriter("Error when updating status of the transaction : "+errWhenUpdatingStatus.Error(), constants.ERROR)
						continue
					}
				}
				//TODOD - Handle metric bind
				continue
			} else {
				//Get transaction receipt
				transactionReceipt, errWhenGettingTheTransactionReceipt := transactionrecipthandler.GetTransactionReceiptForPolygon(pendingHash)
				if errWhenGettingTheTransactionReceipt != nil {
					//check the error
					if errWhenGettingTheTransactionReceipt.Error() == "not found" {
						logger.LogWriter(pendingHash+" transaction is still at the pending state", constants.INFO)
					} else {
						logger.LogWriter("Error when getting the transaction receipt : "+errWhenGettingTheTransactionReceipt.Error(), constants.ERROR)
					}
					continue
				} else if transactionReceipt.Status == "0x1" {
					//Transaction is successfull
					result[i].Status = 118 // SUCCESS

					errWhenUpdatingColletion := dbCollectionHandler.UpdateCollectionsWithNewStatusForPolygon(result[i], 118)
					if errWhenUpdatingColletion != nil {
						logger.LogWriter("Error when updating collection : "+errWhenUpdatingColletion.Error(), constants.ERROR)
						continue
					}

					//updating the actual status in the database
					if result[i].ContractType == "POLYGONEXPERTFORMULA" {
						formulaObj.Status = 118       // SUCCESS
						formulaObj.ActualStatus = 113 // DEPLOYMENT_TRANSACTION_SUCCESS
						errWhenUpdatingActualStatus := object.UpdateSelectedPolygonFormulaFields(formulaObj.FormulaID, formulaObj.TransactionUUID, formulaObj)
						if errWhenUpdatingActualStatus != nil {
							logger.LogWriter("Error when updating the actual status of the formula : "+errWhenUpdatingActualStatus.Error(), constants.ERROR)
							continue
						}
					}
					//TODO-handle for metric bind
					continue

				} else if transactionReceipt.Status == "0x0" {
					//Transaction is failed
					//Get the error for the transaction failure
					errorOccurred, errWhenGettingTheTransaction := pendingTransactionHandler.GetErrorOfFailedTransaction(pendingHash, 2)
					if errWhenGettingTheTransaction != nil {
						logger.LogWriter("Error when getting the transaction error : "+errWhenGettingTheTransaction.Error(), constants.ERROR)
						continue
					}
					logger.LogWriter(pendingHash+" hash failed due to the error : "+errorOccurred, constants.INFO)

					errorMessage := model.EthErrorMessage{
						TransactionHash: pendingHash,
						ErrorMessage:    errorOccurred,
					}

					errWhenInsertingErrorMessage := object.InsertPolygonErrorMessage(errorMessage)
					if errWhenInsertingErrorMessage != nil {
						logger.LogWriter("Error when inserting the error message : "+errWhenInsertingErrorMessage.Error(), constants.ERROR)
					}

					result[i].Status = 119
					result[i].ErrorMessage = errorOccurred

					//TODO -Handle metric bind

				}
				logger.LogWriter(ethClient, constants.INFO)
			}

		}

		return nil
	}).Catch(func(error error) error {
		logger.LogWriter("Error @CheckContractStatus "+error.Error(), constants.DEBUG)
		return error
	})
	p.Await()
}
