package ActivityContractWriters

import "github.com/dileepaj/tracified-gateway/model"

//	For generating the solidity codes for the pivot struct, array, and getter
//	This method will be called only when the pivot field array is not empty

func WritePivotCommonCode() (model.EthGeneralPivotField, error) {

	// pivot struct
	structComment := "\t" + "// PivotField structure"
	structStart := "\t" + "struct PivotField {"
	name := "\t\t" + "string name;"
	key := "\t\t" + "string key;		// converted value to base64"
	field := "\t\t" + "string field;"
	condition := "\t\t" + "string condition;"
	value := "\t\t" + "string value;"
	artifactTemplateID := "\t\t" + "string artifactTemplateID;"
	artifactDataID := "\t\t" + "string artifactDataID;"
	formulaID := "\t\t" + "string formulaID;"
	structEnd := "\t" + "}"
	pivotStruct := structComment + "\n" + structStart + "\n" + name + "\n" + key + "\n" + field + "\n" + condition + "\n" + value + "\n" + artifactTemplateID + "\n" + artifactDataID + "\n" + formulaID + "\n" + structEnd + "\n\n"

	// pivot array
	arrayComment := "\t" + "// Array to store all the pivot fields of the formula"
	arrayDeclaration := "\t" + "PivotField[] private allPivotFields;"
	pivotArray := arrayComment + "\n" + arrayDeclaration + "\n\n"

	// pivot getter
	getterComment := "\t" + "// Getter for retrieving all pivot fields"
	getter := "\t" + "function getPivotFields() public view returns (PivotField[] memory) {"
	getterReturn := "\t\t" + "return allPivotFields;"
	getterEnd := "\t" + "}"
	pivotGetter := getterComment + "\n" + getter + "\n" + getterReturn + "\n" + getterEnd + "\n\n"

	pivotFieldCommonCode := model.EthGeneralPivotField{
		PivotStructure: pivotStruct,
		PivotArray:     pivotArray,
		PivotGetter:    pivotGetter,
	}

	return pivotFieldCommonCode, nil
}
