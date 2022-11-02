package executionTemplates

import (
	"github.com/dileepaj/tracified-gateway/model"
)

// returns -> start variable declarations, setter list, built equation, error
func Template1Builder(executionTemplate model.ExecutionTemplate) ([]string, []string, string, error) {
	var startVariableDeclarations []string
	var setterList []string
	var strTemplate string

	// get the generated solidity code for start variable and append it to the strTemplate
	varDeclaration := `uint public ` + executionTemplate.S_StartVarName + `;`	
	startVariableDeclarations = append(startVariableDeclarations, varDeclaration)
	setter := `function set` + executionTemplate.S_StartVarName + `(uint _` + executionTemplate.S_StartVarName + `) public {` + executionTemplate.S_StartVarName + ` = _` + executionTemplate.S_StartVarName + `;}`
	setterList = append(setterList, setter)
	strTemplate = executionTemplate.S_StartVarName

	// loop through the commands and get the generated solidity code for each command and append it to the strTemplate
	for _, command := range executionTemplate.Lst_Commands {
		startVariables, setters, commandForSolidity, errInCommand := CommandBuilder(command)
		if errInCommand != nil {
			return nil, nil , "", errInCommand
		}
		strTemplate = strTemplate + commandForSolidity
		startVariableDeclarations = append(startVariableDeclarations, startVariables...)
		setterList = append(setterList, setters...)
	}

	return startVariableDeclarations, setterList, strTemplate, nil
}