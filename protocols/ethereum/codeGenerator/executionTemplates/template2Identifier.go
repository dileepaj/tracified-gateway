package executionTemplates

import (
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/ethereum/components"
)

func Template2Builder(executionTemplate model.ExecutionTemplate) (string, error) {
	var strTemplate string

	// get the generated solidity code for start variable and append it to the strTemplate
	startVariableForSolidity, errorInStartVariable := components.GenerateStartVariable(executionTemplate.S_StartVarName)
	if errorInStartVariable != nil {
		return "", errorInStartVariable
	}
	strTemplate = startVariableForSolidity

	return strTemplate, nil
}