package activitywriters

import (
	"os"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

func ActivityContractDeployer(metricID string, element model.MetricDataBindActivityRequest) error {

	var pivotCode model.EthGeneralPivotField
	var errWhenGettingPivotCodes error
	//reqType := "METRIC"

	//call the general code writer
	generalCodes, errWhenGettingGeneralCodes := WriteGeneralCode(metricID, element.MetricFormula.ID)
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

	template := generalCodes.License + generalCodes.PragmaLine + generalCodes.ContractStart + generalCodes.FormulaStructure + generalCodes.ValueStructure + pivotCode.PivotStructure + generalCodes.ValueArray + pivotCode.PivotArray + previousCode.Setter + formulaDeceleration + previousCode.Getter + generalCodes.FormulaGetter + generalCodes.ValueGetter + pivotCode.PivotGetter + generalCodes.ContractEnd

	//generate the solidity file in the specified location
	contactName := generalCodes.ContractName
	fo, errInOutput := os.Create(commons.GoDotEnvVariable("METRICCONTRACTLOCATION") + "/" + contactName + `.sol`)
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

	return nil
}
