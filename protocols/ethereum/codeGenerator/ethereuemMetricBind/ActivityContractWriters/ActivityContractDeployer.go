package ActivityContractWriters

import (
	"encoding/base64"
	"os"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/ethereum/deploy"
	ethereumsocialimpact "github.com/dileepaj/tracified-gateway/services/ethereumServices/ethereumSocialImpact"
	"github.com/sirupsen/logrus"
)

func ActivityContractDeployer(metricMapID string, formulaMapID string, metricID string, element model.MetricDataBindActivityRequest, metricName string, metricElement model.MetricReq, userElement model.User, ethMetricActivityObj model.EthereumMetricBind) error {
	object := dao.Connection{}
	var pivotCode model.EthGeneralPivotField
	var errWhenGettingPivotCodes error
	reqType := "METRIC"

	//call the general code writer
	generalCodes, errWhenGettingGeneralCodes := WriteGeneralCode(metricMapID, formulaMapID)
	if errWhenGettingGeneralCodes != nil {
		ethMetricActivityObj.ActualStatus = 102 	// SMART_CONTRACT_GENERATION_FAILED
		errorWhenUpdatingStatus := object.UpdateSelectedEthMetricFields(metricID, ethMetricActivityObj.TransactionUUID, ethMetricActivityObj)
		if errorWhenUpdatingStatus != nil {
			logrus.Error("Error when updating status for ethereum metric : ", errorWhenUpdatingStatus)
		}
		logrus.Error("Error when generating general code for activity contract : ", errWhenGettingGeneralCodes)
		return errWhenGettingGeneralCodes
	}

	//check if the pivot is not empty and then call the pivot commands
	if len(element.MetricFormula.PivotFields) > 0 {
		pivotCode, errWhenGettingPivotCodes = WritePivotCommonCode()
		if errWhenGettingPivotCodes != nil {
			ethMetricActivityObj.ActualStatus = 102 	// SMART_CONTRACT_GENERATION_FAILED
			errorWhenUpdatingStatus := object.UpdateSelectedEthMetricFields(metricID, ethMetricActivityObj.TransactionUUID, ethMetricActivityObj)
			if errorWhenUpdatingStatus != nil {
				logrus.Error("Error when updating status for ethereum metric : ", errorWhenUpdatingStatus)
			}
			logrus.Error("Error when generating pivot code for activity contract : ", errWhenGettingPivotCodes)
			return errWhenGettingPivotCodes
		}
	}

	//call previous address code writer
	previousCode, errWhenGettingPreviousCodes := WritePreviousCommonCode(metricID)
	if errWhenGettingPreviousCodes != nil {
		ethMetricActivityObj.ActualStatus = 102 	// SMART_CONTRACT_GENERATION_FAILED
		errorWhenUpdatingStatus := object.UpdateSelectedEthMetricFields(metricID, ethMetricActivityObj.TransactionUUID, ethMetricActivityObj)
		if errorWhenUpdatingStatus != nil {
			logrus.Error("Error when updating status for ethereum metric : ", errorWhenUpdatingStatus)
		}
		logrus.Error("Error when generating previous code for activity contract : ", errWhenGettingPreviousCodes)
		return errWhenGettingPreviousCodes
	}

	//call formula deceleration
	formulaDeceleration, errWhenGettingFormulaDeceleration := GetFormulaDefinitionCode(element)
	if errWhenGettingFormulaDeceleration != nil {
		ethMetricActivityObj.ActualStatus = 102 	// SMART_CONTRACT_GENERATION_FAILED
		errorWhenUpdatingStatus := object.UpdateSelectedEthMetricFields(metricID, ethMetricActivityObj.TransactionUUID, ethMetricActivityObj)
		if errorWhenUpdatingStatus != nil {
			logrus.Error("Error when updating status for ethereum metric : ", errorWhenUpdatingStatus)
		}
		logrus.Error("Error when generating formula deceleration for activity contract : ", errWhenGettingFormulaDeceleration)
		return errWhenGettingFormulaDeceleration
	}

	addDetailsFunction, errWhenGettingAddDetailsFunction := AddDetailsMethodWriter(element)
	if errWhenGettingAddDetailsFunction != nil {
		ethMetricActivityObj.ActualStatus = 102 	// SMART_CONTRACT_GENERATION_FAILED
		errorWhenUpdatingStatus := object.UpdateSelectedEthMetricFields(metricID, ethMetricActivityObj.TransactionUUID, ethMetricActivityObj)
		if errorWhenUpdatingStatus != nil {
			logrus.Error("Error when updating status for ethereum metric : ", errorWhenUpdatingStatus)
		}
		logrus.Error("Error when generating add details function for activity contract : ", errWhenGettingAddDetailsFunction)
		return errWhenGettingAddDetailsFunction
	}

	//!File structure should be a follows
	//License
	//Pragma line
	//Contract start
	//Structures
	//Arrays
	//Previous address declaration
	//Formula deceleration
	//Add detail function
	//previous address getter
	//formula getter
	//value getter
	//pivot getter

	template := generalCodes.License + generalCodes.PragmaLine + generalCodes.ContractStart + generalCodes.FormulaStructure + generalCodes.ValueStructure + pivotCode.PivotStructure + generalCodes.ValueArray + pivotCode.PivotArray + previousCode.Setter + formulaDeceleration + addDetailsFunction + previousCode.Getter + generalCodes.FormulaGetter + generalCodes.ValueGetter + pivotCode.PivotGetter + generalCodes.ContractEnd
	ethMetricActivityObj.ActualStatus = 103 	// SMART_CONTRACT_GENERATION_COMPLETED
	errorWhenUpdatingStatus1 := object.UpdateSelectedEthMetricFields(metricID, ethMetricActivityObj.TransactionUUID, ethMetricActivityObj)
	if errorWhenUpdatingStatus1 != nil {
		logrus.Error("Error when updating status for ethereum metric : ", errorWhenUpdatingStatus1)
	}
	//generate the solidity file in the specified location
	contractName := generalCodes.ContractName
	fo, errInOutput := os.Create(commons.GoDotEnvVariable("METRICCONTRACTLOCATION") + "/" + contractName + `.sol`)
	if errInOutput != nil {
		ethMetricActivityObj.ActualStatus = 104	// WRITING_CONTRACT_TO_FILE_FAILED
		errorWhenUpdatingStatus := object.UpdateSelectedEthMetricFields(metricID, ethMetricActivityObj.TransactionUUID, ethMetricActivityObj)
		if errorWhenUpdatingStatus != nil {
			logrus.Error("Error when updating status for ethereum metric : ", errorWhenUpdatingStatus)
		}
		logrus.Error("Error when generating metadata contract file: ", errInOutput)
		return errInOutput
	}

	//write into the file
	defer fo.Close()
	_, errInWritingOutput := fo.Write([]byte(template))
	if errInWritingOutput != nil {
		ethMetricActivityObj.ActualStatus = 104	// WRITING_CONTRACT_TO_FILE_FAILED
		errorWhenUpdatingStatus := object.UpdateSelectedEthMetricFields(metricID, ethMetricActivityObj.TransactionUUID, ethMetricActivityObj)
		if errorWhenUpdatingStatus != nil {
			logrus.Error("Error when updating status for ethereum metric : ", errorWhenUpdatingStatus)
		}
		logrus.Error("Error when writing into the metadata contract file: ", errInWritingOutput)
		return errInWritingOutput
	}

	ethMetricActivityObj.ActualStatus = 105	// WRITING_CONTRACT_TO_FILE_COMPLETED
	errorWhenUpdatingStatus2 := object.UpdateSelectedEthMetricFields(metricID, ethMetricActivityObj.TransactionUUID, ethMetricActivityObj)
	if errorWhenUpdatingStatus2 != nil {
		logrus.Error("Error when updating status for ethereum metric : ", errorWhenUpdatingStatus2)
	}

	//generate ABI
	abiString, errWhenGeneratingABI := deploy.GenerateABI(contractName, reqType)
	if errWhenGeneratingABI != nil {
		ethMetricActivityObj.ActualStatus = 106	// GENERATING_ABI_FAILED
		errorWhenUpdatingStatus := object.UpdateSelectedEthMetricFields(metricID, ethMetricActivityObj.TransactionUUID, ethMetricActivityObj)
		if errorWhenUpdatingStatus != nil {
			logrus.Error("Error when updating status for ethereum metric : ", errorWhenUpdatingStatus)
		}
		logrus.Error("Error when generating ABI for metric metadata contract : ", errWhenGeneratingABI)
		return errWhenGeneratingABI
	}

	ethMetricActivityObj.ActualStatus = 107	// GENERATING_ABI_COMPLETED
	errorWhenUpdatingStatus3 := object.UpdateSelectedEthMetricFields(metricID, ethMetricActivityObj.TransactionUUID, ethMetricActivityObj)
	if errorWhenUpdatingStatus3 != nil {
		logrus.Error("Error when updating status for ethereum metric : ", errorWhenUpdatingStatus3)
	}

	//generate BIN
	binString, errWhenGeneratingBIN := deploy.GenerateBIN(contractName, reqType)
	if errWhenGeneratingBIN != nil {
		ethMetricActivityObj.ActualStatus = 108	// GENERATING_BIN_FAILED
		errorWhenUpdatingStatus := object.UpdateSelectedEthMetricFields(metricID, ethMetricActivityObj.TransactionUUID, ethMetricActivityObj)
		if errorWhenUpdatingStatus != nil {
			logrus.Error("Error when updating status for ethereum metric : ", errorWhenUpdatingStatus)
		}
		logrus.Error("Error when generating BIN for metric metadata contract : ", errWhenGeneratingBIN)
		return errWhenGeneratingBIN
	}

	ethMetricActivityObj.ActualStatus = 109	// GENERATING_BIN_COMPLETED
	errorWhenUpdatingStatus4 := object.UpdateSelectedEthMetricFields(metricID, ethMetricActivityObj.TransactionUUID, ethMetricActivityObj)
	if errorWhenUpdatingStatus4 != nil {
		logrus.Error("Error when updating status for ethereum metric : ", errorWhenUpdatingStatus4)
	}

	templateB64 := base64.StdEncoding.EncodeToString([]byte(template))

	ethMetricActivityObj.ContractName = contractName
	ethMetricActivityObj.TemplateString = templateB64
	ethMetricActivityObj.BINstring = binString
	ethMetricActivityObj.ABIstring = abiString

	errWhenUpdatingMetricDetails := object.UpdateEthereumMetricStatus(ethMetricActivityObj.MetricID, ethMetricActivityObj.TransactionUUID, ethMetricActivityObj)
	if errWhenUpdatingMetricDetails != nil {
		logrus.Error("Error when updating the metric metadata contract details : ", errWhenUpdatingMetricDetails)
		return errWhenUpdatingMetricDetails
	}

	errWhenDeploying := ethereumsocialimpact.DeployMetricContract(ethMetricActivityObj)
	if errWhenDeploying != nil {
		logrus.Error("Error when sending to the metric activity contract to deployer : ", errWhenDeploying)
		return errWhenDeploying
	}

	return nil
}
