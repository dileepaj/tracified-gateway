package codeGenerator

import (
	"github.com/dileepaj/tracified-gateway/model"
)

/*
	Generate the general code snippets common to all the formula build smart contracts
	Building:
		Header code snippets
		Result variable
		Metadata structure
		Value structure
		Variable structure
		Semantic constant structure
		Referred constant structure
		Metadata declaration
*/
func WriteGeneralCodeSnippets(element model.FormulaBuildingRequest, contractName string) (model.ContractGeneral, error) {
	//Meta data structure
	metaDataStructComment := "\t" + `//Metadata structure` + "\n"
	metaDataStructHead := "\t" + `struct Metadata {` + "\n"
	metaDataFormulaID := "\t\t" + `string formulaID; //initialize at start` + "\n"
	metaDataFormulaName := "\t\t" + `string formulaName; //initialize at start` + "\n"
	metaDataExpertPK := "\t\t" + `string expertPK; //initialize at start` + "\n"
	metaDataStructEnd := "\t" + `}` + "\n"
	metaDataStructStr := metaDataStructComment + metaDataStructHead + metaDataFormulaID + metaDataFormulaName + metaDataExpertPK + metaDataStructEnd

	//Value data structure
	valueDataStructComment := "\t" + `//Parent value structure` + "\n"
	valueDataStructHead := "\t" + `struct Value {` + "\n"
	valueType := "\t\t" + `string valueType; //initialize at start` + "\n"
	valueID := "\t\t" + `bytes32 valueID; //initialize at start` + "\n"
	valueName := "\t\t" + `string valueName; //initialize at start` + "\n"
	valueDef := "\t\t" + `int256 value; //initialize at start, added using setter` + "\n"
	valueDef = valueDef + "\t\t" + `int256 exponent; //initialize at start, added using setter` + "\n"
	valueDescription := "\t\t" + `string description; //initialize at start` + "\n"
	valueDataStructEnd := "\t" + `}` + "\n"
	valueDataStructStr := valueDataStructComment + valueDataStructHead + valueType + valueID + valueName + valueDef + valueDescription + valueDataStructEnd

	//Variable data structure
	variableStructComment := "\t" + `//Variable structure, child of Value` + "\n"
	variableStructHead := "\t" + `struct Variable {` + "\n"
	variableValue := "\t\t" + `Value value; //initialize at start` + "\n"
	variableUnit := "\t\t" + `bytes32 unit; //initialize at start` + "\n"
	variablePrecision := "\t\t" + `bytes32 precision; //initialize at start` + "\n"
	variableStructEnd := "\t" + `}` + "\n"
	variableStructStr := variableStructComment + variableStructHead + variableValue + variableUnit + variablePrecision + variableStructEnd

	//Semantic data structure
	semanticStructComment := "\t" + `//Semantic constant structure, child of Value` + "\n"
	semanticStructHead := "\t" + `struct SemanticConstant {` + "\n"
	semanticValue := "\t\t" + `Value value; //initialize at start` + "\n"
	semanticStructEnd := "\t" + `}` + "\n"
	semanticStructStr := semanticStructComment + semanticStructHead + semanticValue + semanticStructEnd

	//Referred data structure
	referredStructComment := "\t" + `//Referred constant structure, child of Value` + "\n"
	referredStructHead := "\t" + `struct ReferredConstant {` + "\n"
	referredValue := "\t\t" + `Value value; //initialize at start` + "\n"
	referredUnit := "\t\t" + `bytes32 unit; //initialize at start` + "\n"
	referredRefURL := "\t\t" + `string refUrl; //initialize at start` + "\n"
	referredStructEnd := "\t" + `}` + "\n"
	referredStructStr := referredStructComment + referredStructHead + referredValue + referredUnit + referredRefURL + referredStructEnd

	//Metadata declaration
	metaDataInitComment := "\t" + `//Metadata declaration` + "\n"
	metaDataInit := "\t" + `Metadata metadata = Metadata("` + element.MetricExpertFormula.ID + `","` + element.MetricExpertFormula.Name + `","` + element.User.Publickey + `");` + "\n"

	//Result structure
	resultVariable := "\n\t" + `// Result structure` +
		"\n\t" + `struct Result {` +
		"\n\t\t" + `int256 value;` +
		"\n\t\t" + `int256 exponent;` +
		"\n\t" + `}` + "\n"

	// Result initialization
	resultDeclaration := "\n\t" + `// Result initialization` + "\n\t" + `Result result = Result(0, 0);` + "\n"

	// Calculations object declaration
	calculationObject := "\n\t" + `// Calculation object creation` + "\n\t" + `Calculations calculations = new Calculations();` + "\n"

	generalBuilder := model.ContractGeneral{
		License:                   `// SPDX-License-Identifier: MIT`,
		PragmaLine:                `pragma solidity ^0.8.7;`,
		ImportCalculationsSol:     `import './Calculations.sol';`,
		ContractStart:             `contract ` + contractName + ` {`,
		ResultVariable:            resultVariable,
		MetaDataStructure:         metaDataStructStr,
		ValueDataStructure:        valueDataStructStr,
		VariableStructure:         variableStructStr,
		SemanticConstantStructure: semanticStructStr,
		ReferredConstant:          referredStructStr,
		MetadataDeclaration:       metaDataInitComment + metaDataInit,
		ResultDeclaration:         resultDeclaration,
		CalculationObject:         calculationObject,
		ContractEnd:               `}`,
	}

	return generalBuilder, nil
}
