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

	// PivotField structure
	pivotFieldDataStructComment := "\t" + `// PivotField structure` + "\n"
	pivotFieldDataStructHead := "\t" + `struct PivotField {` + "\n"
	pivotFieldDataStructName := "\t\t" + `string name;` + "\n"
	pivotFieldDataStructKey := "\t\t" + `string key; // converted value to bytes` + "\n"
	pivotFieldDataStructField := "\t\t" + `string field;` + "\n"
	pivotFieldDataStructCondition := "\t\t" + `string condition;` + "\n"
	pivotFieldDataStructValue := "\t\t" + `string value;` + "\n"
	pivotFieldDataStructArtifactTemplateID := "\t\t" + `string artifactTemplateID;` + "\n"
	pivotFieldDataStructArtifactDataID := "\t\t" + `string artifactDataID;` + "\n"
	pivotFieldDataStructFormulaID := "\t\t" + `string formulaID;` + "\n"
	pivotFieldDataStructEnd := "\t" + `}` + "\n"
	pivotFieldDataStructStr := pivotFieldDataStructComment + pivotFieldDataStructHead + pivotFieldDataStructName + pivotFieldDataStructKey + pivotFieldDataStructField + pivotFieldDataStructCondition + pivotFieldDataStructValue + pivotFieldDataStructArtifactTemplateID + pivotFieldDataStructArtifactDataID + pivotFieldDataStructFormulaID + pivotFieldDataStructEnd

	// Metadata declaration
	metaDataInitComment := "\t" + `// Metadata declaration` + "\n"
	metaDataInit := "\t" + `Metadata metadata = Metadata("` + element.Metric.ID + `", "` + element.Metric.Name + `", "` + element.User.TenantID + `", ` + strconv.Itoa(len(element.Metric.MetricActivities)) + `, "` + element.User.Publickey + `");` + "\n"
	metaDataDeclaration := metaDataInitComment + metaDataInit

	// Array to store all the values
	valueMapComment := "\t" + `// Array to store all the values` + "\n"
	valueMapHead := "\t" + `Value[] private allValues;` + "\n"
	valueList := valueMapComment + valueMapHead

	// Array to store all the formulas
	formulaMapComment := "\t" + `// Array to store all the formulas` + "\n"
	formulaMapHead := "\t" + `Formula[] private allFormulas;` + "\n"
	formulaList := formulaMapComment + formulaMapHead

	// Array to store all the pivot fields
	pivotFieldMapComment := "\t" + `// Array to store all the pivot fields` + "\n"
	pivotFieldMapHead := "\t" + `PivotField[] private allPivotFields;` + "\n"
	pivotFieldsList := pivotFieldMapComment + pivotFieldMapHead

	// getFormulaDetails function to get the formula details
	getFormulaFunctionComment := "\t" + `// Getter for formulas` + "\n"
	getFormulaFunctionHead := "\t" + `function getFormulaDetails() public view returns (Formula[] memory) {` + "\n"
	getFormulaFunctionBody := "\t\t" + `Formula[] memory formulas = allFormulas;` + "\n" + "" + "\t\t" + `return formulas;` + "\n"
	getFormulaFunctionEnd := "\t" + `}` + "\n"
	getFormulaFunction := getFormulaFunctionComment + getFormulaFunctionHead + getFormulaFunctionBody + getFormulaFunctionEnd

	// getValueDetails function to get the value details
	getValueFunctionComment := "\t" + `// Getter for values` + "\n"
	getValueFunctionHead := "\t" + `function getValueDetails() public view returns (Value[] memory) {` + "\n"
	getValueFunctionBody := "\t\t" + `Value[] memory values = allValues;` + "\n" + "" + "\t\t" + `return values;` + "\n"
	getValueFunctionEnd := "\t" + `}` + "\n"
	getValueFunction := getValueFunctionComment + getValueFunctionHead + getValueFunctionBody + getValueFunctionEnd

	// getPivotFieldDetails function to get the pivot fields
	getPivotFieldFunctionComment := "\t" + `// Getter for pivot fields` + "\n"
	getPivotFieldFunctionHead := "\t" + `function getPivotFieldDetails() public view returns (PivotField[] memory) {` + "\n"
	getPivotFieldFunctionBody := "\t\t" + `PivotField[] memory pivotFields = allPivotFields;` + "\n" + "" + "\t\t" + `return pivotFields;` + "\n"
	getPivotFieldFunctionEnd := "\t" + `}` + "\n"
	getPivotFieldFunction := getPivotFieldFunctionComment + getPivotFieldFunctionHead + getPivotFieldFunctionBody + getPivotFieldFunctionEnd

	generalBuilder := model.MetricContractGeneral{
		License:                   `// SPDX-License-Identifier: MIT` + "\n\n",
		PragmaLine:                `pragma solidity ^0.8.7;` + "\n\n",
		ContractStart:             `contract ` + contractName + ` { ` + "\n\n",
		MetaDataStructure:         metaDataStructStr + "\n",
		FormulaStructure:          formulaStructStr + "\n",
		ValueDataStructure:        valueDataStructStr + "\n",
		PivotFieldStructure:       pivotFieldDataStructStr + "\n",
		MetadataDeclaration:       metaDataDeclaration + "\n",
		ValueList:                 valueList + "\n",
		FormulaList:               formulaList + "\n",
		PivotFieldList:            pivotFieldsList + "\n",
		GetFormulaDetailsFunction: getFormulaFunction + "\n",
		GetValueDetailsFunction:   getValueFunction + "\n",
		GetPivotFieldDetails:      getPivotFieldFunction + "\n",
		ContractEnd:               `}`,
	}

	return generalBuilder, nil
}
