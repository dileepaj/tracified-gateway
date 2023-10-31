package experthelpers

func WriteCalculationGetterCode(executionTemplateString string) string {
	commentForExecutor := `// method to get the result of the calculation`
	startOfTheExecutor := `function executeCalculation() public {`
	endOfTheExecutor := "\n\t" + `}`

	executorBody := "\n\t\t" + `result.value` + " = " + executionTemplateString + ";" + "\n\t\t"
	executorBody = executorBody + `result.exponent = calculations.GetExponent();` + "\n\t\t"
	contractBody := "\n\n\t" + commentForExecutor + "\n\t" + startOfTheExecutor + executorBody + endOfTheExecutor + "\n"

	return contractBody
}
