package activitywriters

import "github.com/dileepaj/tracified-gateway/model"

//	For generating the solidity codes for the pivot struct, array, and getter
//	This method will be called only when the pivot field array is not empty

func WritePivotCommonCode() (model.EthGeneralPivotField, error) {

	// pivot struct
	structComment := "// Pivot structure"
	structStart := "struct Pivot {"
	name := "string name;"
	key := "string key;		// converted value to base64"
	value := "string value;"
	structEnd := "}"
	pivotStruct := structComment + "\n" + structStart + "\n" + name + "\n" + key + "\n" + value + "\n" + structEnd + "\n\n"

	// pivot array
	arrayComment := "// Array to store all the pivot fields of the formula"
	arrayDeclaration := "Pivot[] private allPivotFields;"
	pivotArray := arrayComment + "\n" + arrayDeclaration + "\n\n"

	// pivot getter
	getterComment := "// Getter for retrieving all pivot fields"
	getter := "function getPivotFields() public view returns (Pivot[] memory) {"
	getterReturn := "return allPivotFields;"
	getterEnd := "}"
	pivotGetter := getterComment + "\n" + getter + "\n" + getterReturn + "\n" + getterEnd + "\n\n"

	pivotFieldCommonCode := model.EthGeneralPivotField{
		PivotStructure: pivotStruct,
		PivotArray:     pivotArray,
		PivotGetter:    pivotGetter,
	}

	return pivotFieldCommonCode, nil
}