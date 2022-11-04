package executionTemplates

import (
	"github.com/dileepaj/tracified-gateway/model"
)

/**
 * @return start variable declarations, setter list, built equation, error
 * @param executionTemplate
 */
 func Template1Builder(executionTemplate model.ExecutionTemplate) (string, error) {
	var strTemplate string					// the final equation from the execution template

	strTemplate = `(` + executionTemplate.S_StartVarName + `.value.value`

	// loop through the commands 
	for _, command := range executionTemplate.Lst_Commands {
		commandForSolidity, errInCommand := CommandBuilder(command)
		if errInCommand != nil {
			return "", errInCommand
		}
		// append the generated solidity code for each command to the strTemplate
		strTemplate = strTemplate + commandForSolidity
	}
	strTemplate = strTemplate + `)`

	return strTemplate, nil
}