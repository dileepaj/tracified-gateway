package activitywriters

import (
	"encoding/base64"
	"os"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/ethereum/deploy"
	"github.com/dileepaj/tracified-gateway/services"
	"github.com/sirupsen/logrus"
)

func ActivityContractDeployer(metricMapID string, formulaMapID string, metricID string, element model.MetricDataBindActivityRequest, metricName string, metricElement model.MetricReq, userElement model.User, ethMetricActivityObj model.EthereumMetricBind) error {

	var pivotCode model.EthGeneralPivotField
	var errWhenGettingPivotCodes error
	reqType := "METRIC"

	//call the general code writer
	generalCodes, errWhenGettingGeneralCodes := WriteGeneralCode(metricMapID, formulaMapID)
	if errWhenGettingGeneralCodes != nil {
		logrus.Error("Error when generating general code for activity contract : ", errWhenGettingGeneralCodes)
		return errWhenGettingGeneralCodes
	}

	//check if the pivot is not empty and then call the pivot commands
	if len(element.MetricFormula.PivotFields) > 0 {
		pivotCode, errWhenGettingPivotCodes = WritePivotCommonCode()
		if errWhenGettingPivotCodes != nil {
			logrus.Error("Error when generating pivot code for activity contract : ", errWhenGettingPivotCodes)
			return errWhenGettingPivotCodes
		}
	}

	//call previous address code writer
	previousCode, errWhenGettingPreviousCodes := WritePreviousCommonCode(metricID)
	if errWhenGettingPreviousCodes != nil {
		logrus.Error("Error when generating previous code for activity contract : ", errWhenGettingPreviousCodes)
		return errWhenGettingPreviousCodes
	}

	//call formula deceleration
	formulaDeceleration, errWhenGettingFormulaDeceleration := GetFormulaDefinitionCode(element)
	if errWhenGettingFormulaDeceleration != nil {
		logrus.Error("Error when generating formula deceleration for activity contract : ", errWhenGettingFormulaDeceleration)
		return errWhenGettingFormulaDeceleration
	}

	//TODO: call add detail function
	addDetailsFunction, errWhenGettingAddDetailsFunction := AddDetailsMethodWriter(element)
	if errWhenGettingAddDetailsFunction != nil {
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

	//generate the solidity file in the specified location
	contractName := generalCodes.ContractName
	fo, errInOutput := os.Create(commons.GoDotEnvVariable("METRICCONTRACTLOCATION") + "/" + contractName + `.sol`)
	if errInOutput != nil {
		logrus.Error("Error when generating metadata contract file: ", errInOutput)
		return errInOutput
	}

	//write into the file
	defer fo.Close()
	_, errInWritingOutput := fo.Write([]byte(template))
	if errInWritingOutput != nil {
		logrus.Error("Error when writing into the metadata contract file: ", errInWritingOutput)
		return errInWritingOutput
	}

	//generate ABI
	abiString, errWhenGeneratingABI := deploy.GenerateABI(contractName, reqType)
	if errWhenGeneratingABI != nil {
		logrus.Error("Error when generating ABI for metric metadata contract : ", errWhenGeneratingABI)
		return errWhenGeneratingABI
	}

	//generate BIN
	binString, errWhenGeneratingBIN := deploy.GenerateBIN(contractName, reqType)
	if errWhenGeneratingBIN != nil {
		logrus.Error("Error when generating BIN for metric metadata contract : ", errWhenGeneratingBIN)
		return errWhenGeneratingBIN
	}

	templateB64 := base64.StdEncoding.EncodeToString([]byte(template))

	ethMetricActivityObj.ContractName = contractName
	ethMetricActivityObj.TemplateString = templateB64
	ethMetricActivityObj.BINstring = binString
	ethMetricActivityObj.ABIstring = abiString

	buildQueueObject := model.SendToQueue{
		EthereumMetricBind: ethMetricActivityObj,
		Type:               "ETHMETRICBIND",
		User:               userElement,
		Status:             "QUEUE",
	}

	errWhenSendingToQueue := services.SendToQueue(buildQueueObject)
	if errWhenSendingToQueue != nil {
		logrus.Error("Error when sending to the metric activity contract to queue : ", errWhenSendingToQueue)
		return errWhenSendingToQueue
	}

	return nil
}
