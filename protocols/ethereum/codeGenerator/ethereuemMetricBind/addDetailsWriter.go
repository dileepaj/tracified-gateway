package ethereuemmetricbind

import (
	"strconv"

	"github.com/dileepaj/tracified-gateway/model"
)

/*
 * This function is used to write the addDetails() method implementation
 */

func WriteAddDetailsFunction(element model.MetricDataBindingRequest) (string, error) {
	functionStr := `` 	// to store the method string
	formulaCount := 0 	// to keep track of the formula count

	functionStr += "\t" + `// function to add details` + "\n"		// adding comment for the method
	functionStr += "\t" + `function addDetails() public {` + "\n"	// adding method declaration start

	// loop through all the activities and get the method calls
	for _, activity := range element.Metric.Activities {
		formulaCount++
		formulaComment := "\t\t// add formula " + strconv.Itoa(formulaCount) + "\n"
		// get the method call string for the formula
		addFormulaStr, errInGettingFormulaString := WriteAddFormula(activity, formulaCount)
		if errInGettingFormulaString != nil {
			return ``, errInGettingFormulaString
		}

		functionStr += formulaComment + addFormulaStr + "\n"
	}

	functionStr += "\t" + `}` + "\n\n"	

	return functionStr, nil
}