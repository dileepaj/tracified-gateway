package activitywriters

import (
	"github.com/dileepaj/tracified-gateway/model"
)

// For generating the solidity codes for the common struct

func WriteGeneralCode(metricID string, formulaID string) (model.ActivityContractGeneral, error) {

	// Contract start
	contractName := "Metric_" + metricID + "_" + formulaID
	contractStart := `contract ` + contractName + ` {` + "\n\n"

	// Formula Structure
	formulaStructComment := `// Formula structure` + "\n"
	formulaStructStart := `struct Formula {` + "\n"
	actualFormulaID := `string formulaID;		// actual formula ID` + "\n"
	contractAddress := `string contractAddress;` + "\n"
	noOfValues := `uint noOfValues;` + "\n"
	activityID := `string activityID;` + "\n"
	activityName := `string activityName;	// converted value to base64` + "\n"
	valueIDs := `string valueIDs;` + "\n"
	formulaStructEnd := `}` + "\n\n"
	formulaStruct := formulaStructComment + formulaStructStart + actualFormulaID + contractAddress + noOfValues + activityID + activityName + valueIDs + formulaStructEnd

	// Value Structure
	valueStructComment := `// Value structure` + "\n"
	valueStructStart := `struct Value {` + "\n"
	valueID := `string valueID;` + "\n"
	valueName := `string valueName;` + "\n"
	workflowID := `string workflowID;` + "\n"
	stageID := `string stageID;` + "\n"
	stageName := `string stageName;		// converted value to base64` + "\n"
	keyName := `string keyName;		// converted value to base64` + "\n"
	tdpType := `string tdpType;` + "\n"
	bindingType := `int bindingType;` + "\n"
	artifactID := `string artifactID;` + "\n"
	primaryKeyRowID := `string primaryKeyRowID;` + "\n"
	artifactTemplateName := `string artifactTemplateName;	// converted value to base64` + "\n"
	fieldKey := `string fieldKey;		// converted value to base64` + "\n"
	fieldName := `string fieldName;		// converted value to base64` + "\n"
	valueStructEnd := `}` + "\n\n"
	valueStruct := valueStructComment + valueStructStart + valueID + valueName + workflowID + stageID + stageName + keyName + tdpType + bindingType + artifactID + primaryKeyRowID + artifactTemplateName + fieldKey + fieldName + valueStructEnd

	// Value array
	valueArrayComment := `// Array to store all the values` + "\n"
	valueArrayDeclaration := `Value[] private allValues;` + "\n\n"
	valueArray := valueArrayComment + valueArrayDeclaration

	// Formula getter
	formulaGetterComment := `// Getter for formula` + "\n"
	formulaGetter := `function getFormula() public view returns (Formula memory) {` + "\n"
	formulaGetterReturn := `return formula;` + "\n"
	formulaGetterEnd := `}` + "\n\n"
	formulaGetterCode := formulaGetterComment + formulaGetter + formulaGetterReturn + formulaGetterEnd

	// Value getter
	valueGetterComment := `// Getter for retrieving all values` + "\n"
	valueGetter := `function getValues() public view returns (Value[] memory) {` + "\n"
	valueGetterReturn := `return allValues;` + "\n"
	valueGetterEnd := `}` + "\n\n"
	valueGetterCode := valueGetterComment + valueGetter + valueGetterReturn + valueGetterEnd

	activity := model.ActivityContractGeneral{
		License:          `// SPDX-License-Identifier: MIT` + "\n\n",
		PragmaLine:       `pragma solidity ^0.8.7;` + "\n\n",
		ContractStart:    contractStart,
		FormulaStructure: formulaStruct,
		ValueStructure:   valueStruct,
		ValueArray:       valueArray,
		FormulaGetter:    formulaGetterCode,
		ValueGetter:      valueGetterCode,
		ContractEnd:      `}`,
	}

	return activity, nil
}
