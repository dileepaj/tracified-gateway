package executionTemplates

import (
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/ethereum/components"
)

func Template1Builder(executionTemplate model.ExecutionTemplate) (string, error) {
	var strTemplate string

	// get the generated solidity code for start variable and append it to the strTemplate
	startVariableForSolidity, errorInStartVariable := components.GenerateStartVariable(executionTemplate.S_StartVarName)
	if errorInStartVariable != nil {
		return "", errorInStartVariable
	}
	strTemplate = startVariableForSolidity

	// loop through the commands and get the generated solidity code for each command and append it to the strTemplate
	for _, command := range executionTemplate.Lst_Commands {
		commandForSolidity, errInCommand := CommandBuilder(command)
		if errInCommand != nil {
			return "", errInCommand
		}
		strTemplate = strTemplate + commandForSolidity
	}

	return strTemplate, nil
}