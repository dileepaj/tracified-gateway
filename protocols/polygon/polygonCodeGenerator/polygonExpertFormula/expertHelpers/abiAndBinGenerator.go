package experthelpers

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/ethereum/deploy"
	"github.com/dileepaj/tracified-gateway/utilities"
)

func AbiAndBinGenerator(contractName string, reqType string, formulaId string, transactionUuid string, formulaObj model.EthereumExpertFormula) (error, string, string) {
	logger := utilities.NewCustomLogger()
	object := dao.Connection{}
	//call the ABI generator
	abiString, errWhenGeneratingAbiString := deploy.GenerateABI(contractName, reqType)
	if errWhenGeneratingAbiString != nil {
		logger.LogWriter("Error when generating the ABI file : "+errWhenGeneratingAbiString.Error(), constants.ERROR)
		return errors.New(errWhenGeneratingAbiString.Error()), "", ""
	}

	//update the collection on abi generation
	formulaObj.ABIstring = abiString
	formulaObj.ActualStatus = 107 // GENERATING_ABI_COMPLETED
	errWhenUpdatingAfterAbiGeneration := object.UpdateSelectedPolygonFormulaFields(formulaId, transactionUuid, formulaObj)
	if errWhenUpdatingAfterAbiGeneration != nil {
		logger.LogWriter("Error when updating polygon collection after the ABI generation : "+errWhenUpdatingAfterAbiGeneration.Error(), constants.ERROR)
		return errors.New(errWhenUpdatingAfterAbiGeneration.Error()), "", ""
	} else {
		logger.LogWriter("ABI inserted to collection", constants.INFO)
	}

	//call the BIN generator
	binString, errWhenGeneratingBinString := deploy.GenerateBIN(contractName, reqType)
	if errWhenGeneratingBinString != nil {
		logger.LogWriter("Error when generating the BIN file : "+errWhenGeneratingBinString.Error(), constants.ERROR)
		return errors.New(errWhenGeneratingBinString.Error()), "", ""
	}

	formulaObj.BINstring = binString
	formulaObj.ActualStatus = 109 // GENERATING_BIN_COMPLETED
	errWhenUpdatingAfterBinGeneration := object.UpdateSelectedPolygonFormulaFields(formulaId, transactionUuid, formulaObj)
	if errWhenUpdatingAfterBinGeneration != nil {
		logger.LogWriter("Error when updating polygon collection after the BIN generation : "+errWhenUpdatingAfterBinGeneration.Error(), constants.ERROR)
		return errors.New(errWhenUpdatingAfterBinGeneration.Error()), "", ""
	} else {
		logger.LogWriter("BIN inserted to the collection", constants.INFO)
	}

	return nil, abiString, binString
}
