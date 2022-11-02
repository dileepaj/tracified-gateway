package executionTemplates

import (
	"github.com/dileepaj/tracified-gateway/model"
)

/**
 * @return start variable declaration, the setter for the start variable, variable name, error
 * @param executionTemplate
 */
func Template2Builder(executionTemplate model.ExecutionTemplate) (string, string, string, error) {
	var varDeclaration string		// the start variable declaration
	var setter string				// the setter for the start variable

	// get the generated solidity code for start variable and append it to the strTemplate
	varDeclaration = `uint public ` + executionTemplate.S_StartVarName + `;`	
	// get the generated solidity code for setter and append it to the setterList
	setter = `function set` + executionTemplate.S_StartVarName + `(uint _` + executionTemplate.S_StartVarName + `) public {` 
	setter = setter + "\n\t\t" + executionTemplate.S_StartVarName + ` = _` + executionTemplate.S_StartVarName + `;` + "\n\t" + `}`

	return varDeclaration, setter, executionTemplate.S_StartVarName, nil
}