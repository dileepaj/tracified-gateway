package ethereuemmetricbind

import (
	"strconv"

	"github.com/dileepaj/tracified-gateway/model"
)

/*
 * This function is used to write the addDetails() method implementation
 */

func WriteAddDetailsFunction(element model.MetricDataBindingRequest) ([]string, []string, string, error) {
	functionStr := `` 	// to store the method string
	formulaCount := 0 	// to keep track of the formula count
	formulaIDs := []string{} // to store the formula IDs
	valueIDs := []string{} // to store the value IDs

	functionStr += "\t" + `// function to add details` + "\n"		// adding comment for the method
	functionStr += "\t" + `function addDetails() public {` + "\n"	// adding method declaration start

	// loop through all the activities and get the method calls
	for _, activity := range element.Metric.Activities {
		// add the formula ID to the array
		formulaIDs = append(formulaIDs, activity.MetricFormula.MetricExpertFormula.ID)
		formulaCount++
		// get the method call string for the formula and the value IDs
		valueIdList, addFormulaStr, errInGettingFormulaString := WriteAddFormula(activity, formulaCount)
		if errInGettingFormulaString != nil {
			return []string{}, []string{}, ``, errInGettingFormulaString
		}
		// add the codes for adding formula if the addFormulaStr is not empty
		if addFormulaStr != "" {
			formulaComment := "\t\t// add formula " + strconv.Itoa(formulaCount) + "\n"
			functionStr += formulaComment + addFormulaStr + "\n"
		}

		// get the string for adding the pivot fields
		addPivotFieldsStr, errInGettingPivotFields := WriteAddPivotField(activity)
		if errInGettingPivotFields != nil {
			return []string{}, []string{}, ``, errInGettingPivotFields
		}
		// add the codes for adding pivot fields if the addPivotFieldsStr is not empty
		if addPivotFieldsStr != "" {
			addPivotFieldComment := "\t\t// add pivot fields related to the formula " + strconv.Itoa(formulaCount) + " to the array" + "\n"
			functionStr += addPivotFieldComment + addPivotFieldsStr + "\n"
		}

		valueIDs = append(valueIDs, valueIdList...)
	}

	functionStr += "\t" + `}` + "\n\n"

	return formulaIDs, valueIDs, functionStr, nil
}