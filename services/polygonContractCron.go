package services

import (
	"strconv"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/services/ethereumServices/dbCollectionHandler"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
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
						logrus.Error("Error when updating status of the transaction : " + errWhenUpdatingStatus.Error())
						continue
					}
				}
				//TODOD - Handle metric bind
				continue
			} else {
				//Get transaction receipt
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
