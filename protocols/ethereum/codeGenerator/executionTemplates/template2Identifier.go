package executionTemplates

import (
	"github.com/dileepaj/tracified-gateway/model"
)

// returns -> start variable declaration, the setter for the start variable, variable name, error

func Template2Builder(executionTemplate model.ExecutionTemplate) (string, string, string, error) {
	var varDeclaration string
	var setter string
	var strTemplate string

	// get the generated solidity code for start variable and append it to the strTemplate
	varDeclaration = `uint public ` + executionTemplate.S_StartVarName + `;`	
	setter = `function set` + executionTemplate.S_StartVarName + `(uint _` + executionTemplate.S_StartVarName + `) public {` + executionTemplate.S_StartVarName + ` = _` + executionTemplate.S_StartVarName + `;}`
	strTemplate = executionTemplate.S_StartVarName

	return varDeclaration, setter, strTemplate, nil
}