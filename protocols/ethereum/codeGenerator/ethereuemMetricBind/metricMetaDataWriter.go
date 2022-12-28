package ethereuemmetricbind

import (
	"github.com/dileepaj/tracified-gateway/model"

	"strconv"
)

/*
	Generate the general code snippets common to all the metric bind smart contracts
	Building:
		Header code snippets
		Metadata structure
		Formula structure
		Value structure
		Metadata declaration
		Map to store all the values
		Map to store all the formulas
		Array declaration to store all the formula details as a string
		AddValue function
		AddFormula function
		GetFormula function to get the formula details as array
*/

func WriteMetricGeneralCodeSnippets(element model.MetricDataBindingRequest, contractName string) (model.MetricContractGeneral, error) {
	// Metadata structure
	metaDataStructComment := "\t" + `// Metadata structure` + "\n"
	metaDataStructHead := "\t" + `struct Metadata {` + "\n"
	metaDataMetricID := "\t\t" + `string metricID; ` + "\n"
	metaDataMetricName := "\t\t" + `string metricName; // converted value to bytes` + "\n"
	metaDataTenantID := "\t\t" + `string tenantID;` + "\n"
	metaDataNoOfFormulas := "\t\t" + `uint noOfFormulas;` + "\n"
	metaDataTrustNetPK := "\t\t" + `string trustNetPK;` + "\n"
	metaDataStructEnd := "\t" + `}` + "\n"
	metaDataStructStr := metaDataStructComment + metaDataStructHead + metaDataMetricID + metaDataMetricName + metaDataTenantID + metaDataNoOfFormulas + metaDataTrustNetPK + metaDataStructEnd

	// Formula structure
	formulaStructComment := "\t" + `// Formula structure` + "\n"
	formulaStructHead := "\t" + `struct Formula {` + "\n"
	formulaStructActualID := "\t\t" + `string formulaID; // actual formula ID` + "\n"
	formulaStructContractAddress := "\t\t" + `string contractAddress;` + "\n"
	formulaStructNoOfValues := "\t\t" + `uint noOfValues;` + "\n"
	formulaStructActivityID := "\t\t" + `string activityID;` + "\n"
	formulaStructActivityName := "\t\t" + `string activityName; // converted value to bytes` + "\n"
	formulaStructValueIDs := "\t\t" + `string valueIDs;` + "\n"
	formulaStructEnd := "\t" + `}` + "\n"
	formulaStructStr := formulaStructComment + formulaStructHead + formulaStructActualID + formulaStructContractAddress + formulaStructNoOfValues + formulaStructActivityID + formulaStructActivityName + formulaStructValueIDs + formulaStructEnd

	// Value structure
	valueDataStructComment := "\t" + `// Value structure` + "\n"
	valueDataStructHead := "\t" + `struct Value {` + "\n"
	valueDataStructValueID := "\t\t" + `string valueID;` + "\n"
	valueDataStructValueName := "\t\t" + `string valueName;` + "\n"
	valueDataStructWorkflowID := "\t\t" + `string workflowID;` + "\n"
	valueDataStructStageID := "\t\t" + `string stageID;` + "\n"
	valueDataStructStageName := "\t\t" + `string stageName; // converted value to bytes` + "\n"
	valueDataStructKeyName := "\t\t" + `string keyName; // converted value to bytes` + "\n"
	valueDataStructTDPType := "\t\t" + `string tdpType;` + "\n"
	valueDataStructBindingType := "\t\t" + `int bindingType;` + "\n"
	valueDataStructArtifactID := "\t\t" + `string artifactID;` + "\n"
	valueDataStructPrimaryKeyRowID := "\t\t" + `string primaryKeyRowID;` + "\n"
	valueDataStructArtifactTemplateName := "\t\t" + `string artifactTemplateName; // converted value to bytes` + "\n"
	valueDataStructFieldKey := "\t\t" + `string fieldKey; // converted value to bytes` + "\n"
	valueDataStructFieldName := "\t\t" + `string fieldName; // converted value to bytes` + "\n"
	valueDataStructEnd := "\t" + `}` + "\n"
	valueDataStructStr := valueDataStructComment + valueDataStructHead + valueDataStructValueID + valueDataStructValueName + valueDataStructWorkflowID + valueDataStructStageID + valueDataStructStageName + valueDataStructKeyName + valueDataStructTDPType + valueDataStructBindingType + valueDataStructArtifactID + valueDataStructPrimaryKeyRowID + valueDataStructArtifactTemplateName + valueDataStructFieldKey + valueDataStructFieldName + valueDataStructEnd

	// Metadata declaration
	metaDataInitComment := "\t" + `// Metadata declaration` + "\n"
	metaDataInit := "\t" + `Metadata metadata = Metadata("` + element.Metric.ID + `", "` + element.Metric.Name + `", "` + element.User.TenantID + `", ` + strconv.Itoa(len(element.Metric.Activities)) + `, "` + element.User.Publickey + `");` + "\n"
	metaDataDeclaration := metaDataInitComment + metaDataInit

	// Map to store all the values
	valueMapComment := "\t" + `// Map to store all the values` + "\n"
	valueMapHead := "\t" + `mapping(string => Value) private allValues;` + "\n"
	valueMap := valueMapComment + valueMapHead

	// Map to store all the formulas
	formulaMapComment := "\t" + `// Map to store all the formulas` + "\n"
	formulaMapHead := "\t" + `mapping(string => Formula) private allFormulas;` + "\n"
	formulaMap := formulaMapComment + formulaMapHead

	// AddValue function
	addValueFunctionComment := "\t" + `// AddValue function` + "\n"
	addValueFunctionHead := "\t" + `function addValue(string memory _valueID, string memory _valueName, string memory _workflowID, string memory _stageID, string memory _stageName, string memory _keyName, string memory _TDPType, int _bindingType, string memory _artifactID, string memory _primaryKeyRowID, string memory _artifactTemplateName, string memory _fieldKey, string memory _fieldName) internal {` + "\n"
	addValueFunctionBodyComment := "\t\t" + `// Add the value to the map` + "\n"
	addValueFunctionBody := "\t\t" + `allValues[_valueID] = Value(_valueID, _valueName, _workflowID, _stageID, _stageName, _keyName, _TDPType, _bindingType, _artifactID, _primaryKeyRowID, _artifactTemplateName, _fieldKey, _fieldName);` + "\n"
	addValueFunctionEnd := "\t" + `}` + "\n"
	addValueFunction := addValueFunctionComment + addValueFunctionHead + addValueFunctionBodyComment + addValueFunctionBody + addValueFunctionEnd

	// AddFormula function
	addFormulaFunctionComment := "\t" + `// AddFormula function` + "\n"
	addFormulaFunctionHead := "\t" + `function addFormula(string memory _formulaID, string memory _contractAddress, uint256 _noOfVariables, string memory _activityID, string memory _activityName, string memory _valueList) internal {` + "\n"
	addFormulaBodyComment := "\t\t" + `// Add the formula to the map` + "\n"
	addFormulaBody := "\t\t" + `allFormulas[_formulaID] = Formula(_formulaID, _contractAddress, _noOfVariables, _activityID, _activityName, _valueList);` + "\n"
	addFormulaFunctionEnd := "\t" + `}` + "\n"
	addFormulaFunction := addFormulaFunctionComment + addFormulaFunctionHead + addFormulaBodyComment + addFormulaBody + addFormulaFunctionEnd

	// GetFormula function to get the formula details by ID
	getFormulaFunctionComment := "\t" + `// Getter to get the formula details by ID` + "\n"
	getFormulaFunctionHead := "\t" + `function getFormulaDetails(string memory _id) public view returns (Formula memory) {` + "\n"
	getFormulaFunctionBody := "\t\t" + `Formula memory formula = allFormulas[_id];` + "\n" + "" + "\t\t" + `return formula;` + "\n"
	getFormulaFunctionEnd := "\t" + `}` + "\n"
	getFormulaFunction := getFormulaFunctionComment + getFormulaFunctionHead + getFormulaFunctionBody + getFormulaFunctionEnd

	// GetValue function to get the value details by ID
	getValueFunctionComment := "\t" + `// Getter to get the value details by ID` + "\n"
	getValueFunctionHead := "\t" + `function getValueDetails(string memory _id) public view returns (Value memory) {` + "\n"
	getValueFunctionBody := "\t\t" + `Value memory value = allValues[_id];` + "\n" + "" + "\t\t" + `return value;` + "\n"
	getValueFunctionEnd := "\t" + `}` + "\n"
	getValueFunction := getValueFunctionComment + getValueFunctionHead + getValueFunctionBody + getValueFunctionEnd

	generalBuilder := model.MetricContractGeneral{
		License:                   `// SPDX-License-Identifier: MIT` + "\n\n",
		PragmaLine:                `pragma solidity ^0.8.7;` + "\n\n",
		ContractStart:             `contract ` + contractName + ` { ` + "\n\n",
		MetaDataStructure:         metaDataStructStr + "\n",
		FormulaStructure:          formulaStructStr + "\n",
		ValueDataStructure:        valueDataStructStr + "\n",
		MetadataDeclaration:       metaDataDeclaration + "\n",
		ValueMap:                  valueMap + "\n",
		FormulaMap:                formulaMap + "\n",
		AddValueFunction:          addValueFunction + "\n",
		AddFormulaFunction:        addFormulaFunction + "\n",
		GetFormulaDetailsFunction: getFormulaFunction + "\n",
		GetValueDetailsFunction:   getValueFunction,
		ContractEnd:               `}`,
	}

	return generalBuilder, nil
}
