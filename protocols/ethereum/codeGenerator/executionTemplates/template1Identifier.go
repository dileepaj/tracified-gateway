package executionTemplates

import (
	"github.com/dileepaj/tracified-gateway/model"
)

/**
 * @return start variable declarations, setter list, built equation, error
 * @param executionTemplate
 */
 func Template1Builder(executionTemplate model.ExecutionTemplate) ([]string, []string, string, error) {
	var startVariableDeclarations []string	// the list of start variable declarations
	var setterList []string					// the list of setters for the start variables
	var strTemplate string					// the final equation from the execution template

	// get the generated solidity code for start variable and append it to the strTemplate
	varDeclaration := `uint public ` + executionTemplate.S_StartVarName + `;`	
	startVariableDeclarations = append(startVariableDeclarations, varDeclaration)

	// get the generated solidity code for setter and append it to the setterList
	setter := `function set` + executionTemplate.S_StartVarName + `(uint _` + executionTemplate.S_StartVarName + `) public {` 
	setter = setter + "\n\t\t" + executionTemplate.S_StartVarName + ` = _` + executionTemplate.S_StartVarName + `;` + "\n\t" + `}`
	setterList = append(setterList, setter)
	strTemplate = `(` + executionTemplate.S_StartVarName

	// loop through the commands 
	for _, command := range executionTemplate.Lst_Commands {
		startVariables, setters, commandForSolidity, errInCommand := CommandBuilder(command)
		if errInCommand != nil {
			return nil, nil , "", errInCommand
		}
		// append the generated solidity code for each command to the strTemplate
		strTemplate = strTemplate + commandForSolidity
		// append the returned start variable list to the startVariableDeclarations
		startVariableDeclarations = append(startVariableDeclarations, startVariables...)
		// append the returned setter list to the setterList
		setterList = append(setterList, setters...)
	}
	strTemplate = strTemplate + `)`

	return startVariableDeclarations, setterList, strTemplate, nil
}