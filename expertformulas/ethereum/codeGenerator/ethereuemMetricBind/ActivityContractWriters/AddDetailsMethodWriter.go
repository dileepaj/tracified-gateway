package ActivityContractWriters

import "github.com/dileepaj/tracified-gateway/model"

// AddDetailsMethodWriter writes the AddDetails method of the activity contract
// Get codes for adding values to the values array
// Check whether the pivot field is empty or not and add the pivot field to the values array

func AddDetailsMethodWriter(element model.MetricDataBindActivityRequest) (string, error) {

	// Get codes for adding values to the values array
	addValuesCode, errWhenGettingVariableData := AddValuesWriter(element)
	if errWhenGettingVariableData != nil {
		return "", errWhenGettingVariableData
	}

	// Check whether the pivot field is empty or not and add the pivot field to the values array
	addPivotFieldsCode := ""
	if len(element.MetricFormula.PivotFields) > 0 {
		addPivotFields, errWhenGettingPivotFields := AddPivotFieldsWriter(element.MetricFormula.PivotFields, element.MetricFormula.MetricExpertFormula.ID)
		if errWhenGettingPivotFields != nil {
			return "", errWhenGettingPivotFields
		}
		addPivotFieldsCode += addPivotFields
	}

	addDetailsCodeComment := "\t" + "// Method to add all the values and pivot fields relevant to the formula" + "\n"
	functionStart := "\t" + `function addDetails() public {` + "\n"
	functionEnd := "\t" + `}` + "\n\n"
	addDetailsCode := addDetailsCodeComment + functionStart + addValuesCode + addPivotFieldsCode + functionEnd
	
	return addDetailsCode, nil
}