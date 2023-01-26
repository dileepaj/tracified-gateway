package activitywriters

import "github.com/dileepaj/tracified-gateway/model"

//	For generating the solidity codes for the pivot struct, array, and getter
//	This method will be called only when the pivot field array is not empty

func WritePivotCommonCode() (model.EthGeneralPivotField, error) {

	// pivot struct
	structComment := "\t" + "// Pivot structure"
	structStart := "\t" + "struct Pivot {"
	name := "\t\t" + "string name;"
	key := "\t\t" + "string key;		// converted value to base64"
	value := "\t\t" + "string value;"
	structEnd := "\t" + "}"
	pivotStruct := structComment + "\n" + structStart + "\n" + name + "\n" + key + "\n" + value + "\n" + structEnd + "\n\n"

	// pivot array
	arrayComment := "\t" + "// Array to store all the pivot fields of the formula"
	arrayDeclaration := "\t" + "Pivot[] private allPivotFields;"
	pivotArray := arrayComment + "\n" + arrayDeclaration + "\n\n"

	// pivot getter
	getterComment := "\t" + "// Getter for retrieving all pivot fields"
	getter := "\t" + "function getPivotFields() public view returns (Pivot[] memory) {"
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