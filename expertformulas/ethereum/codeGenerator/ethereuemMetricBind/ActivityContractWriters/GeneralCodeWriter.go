package ActivityContractWriters

import (
	"github.com/dileepaj/tracified-gateway/model"
)

// For generating the solidity codes for the common struct, arrays, and getters

func WriteGeneralCode(metricMapID string, formulaMapID string) (model.ActivityContractGeneral, error) {

	// Contract start
	contractName := "Metric_" + metricMapID + "_Formula_" + formulaMapID
	contractStart := `contract ` + contractName + ` {` + "\n\n"

	// Formula Structure
	formulaStructComment := "\t" + `// Formula structure` + "\n"
	formulaStructStart := "\t" + `struct Formula {` + "\n"
	actualFormulaID := "\t\t" + `string formulaID;		// actual formula ID` + "\n"
	contractAddress := "\t\t" + `string contractAddress;` + "\n"
	noOfValues := "\t\t" + `uint noOfValues;` + "\n"
	activityID := "\t\t" + `string activityID;` + "\n"
	activityName := "\t\t" + `string activityName;	// converted value to base64` + "\n"
	valueIDs := "\t\t" + `string valueIDs;` + "\n"
	formulaStructEnd := "\t" + `}` + "\n\n"
	formulaStruct := formulaStructComment + formulaStructStart + actualFormulaID + contractAddress + noOfValues + activityID + activityName + valueIDs + formulaStructEnd

	// Value Structure
	valueStructComment := "\t" + `// Value structure` + "\n"
	valueStructStart := "\t" + `struct Value {` + "\n"
	valueID := "\t\t" + `string valueID;` + "\n"
	valueName := "\t\t" + `string valueName;` + "\n"
	workflowID := "\t\t" + `string workflowID;` + "\n"
	stageID := "\t\t" + `string stageID;` + "\n"
	stageName := "\t\t" + `string stageName;		// converted value to base64` + "\n"
	keyName := "\t\t" + `string keyName;		// converted value to base64` + "\n"
	tdpType := "\t\t" + `string tdpType;` + "\n"
	bindingType := "\t\t" + `int bindingType;` + "\n"
	primaryKeyRowID := "\t\t" + `string primaryKeyRowID;` + "\n"
	artifactTemplateID := "\t\t" + `string artifactTemplateID;` + "\n"
	artifactTemplateName := "\t\t" + `string artifactTemplateName;	// converted value to base64` + "\n"
	fieldKey := "\t\t" + `string fieldKey;		// converted value to base64` + "\n"
	fieldName := "\t\t" + `string fieldName;		// converted value to base64` + "\n"
	valueStructEnd := "\t" + `}` + "\n\n"
	valueStruct := valueStructComment + valueStructStart + valueID + valueName + workflowID + stageID + stageName + keyName + tdpType + bindingType + primaryKeyRowID + artifactTemplateID + artifactTemplateName + fieldKey + fieldName + valueStructEnd

	// Value array
	valueArrayComment := "\t" + `// Array to store all the values` + "\n"
	valueArrayDeclaration := "\t" + `Value[] private allValues;` + "\n\n"
	valueArray := valueArrayComment + valueArrayDeclaration

	// Formula getter
	formulaGetterComment := "\t" + `// Getter for formula` + "\n"
	formulaGetter := "\t" + `function getFormula() public view returns (Formula memory) {` + "\n"
	formulaGetterReturn := "\t\t" + `return formula;` + "\n"
	formulaGetterEnd := "\t" + `}` + "\n\n"
	formulaGetterCode := formulaGetterComment + formulaGetter + formulaGetterReturn + formulaGetterEnd

	// Value getter
	valueGetterComment := "\t" + `// Getter for retrieving all values` + "\n"
	valueGetter := "\t" + `function getValues() public view returns (Value[] memory) {` + "\n"
	valueGetterReturn := "\t\t" + `return allValues;` + "\n"
	valueGetterEnd := "\t" + `}` + "\n\n"
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
		ContractName:     contractName,
	}

	return activity, nil
}
