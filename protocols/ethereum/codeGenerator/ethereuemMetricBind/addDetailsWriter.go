package ethereuemmetricbind

import (
	"strconv"

	"github.com/dileepaj/tracified-gateway/model"
)

func WriteAddDetailsFunction(element model.MetricDataBindingRequest) (string, error) {
	functionStr := ``
	formulaCount := 0

	functionStr += "\t" + `function addDetails() public {` + "\n\n"

	// loop through all the activities and get the method calls
	for _, formula := range element.Metric.Activities {
		formulaCount++
		formulaComment := "\t\t // add formula " + strconv.Itoa(formulaCount) + "\n"
		addFormulaStr, errInGettingFormulaString := WriteAddFormula(formula, formulaCount)
		if errInGettingFormulaString != nil {
			return ``, errInGettingFormulaString
		}

		functionStr += formulaComment + addFormulaStr + "\n"
	}

	functionStr += "\t" + `}` + "\n"	

	return functionStr, nil
}