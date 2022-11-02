package executionTemplates

import "github.com/dileepaj/tracified-gateway/model"

/**
 * @return start variable declarations, setter list, built equation, error
 * @param executionTemplate
 */
func ExecutionTemplateDivider(executionTemplate model.ExecutionTemplate) ([]string, []string, string, error) {
	var strTemplate string						// the final equation
	var startVariableDeclarations []string		// the list of start variable declarations
	var setterList []string						// the list of setters for the start variables

	// check whether the execution template has list of commands or not and call the relevant template builder
	if executionTemplate.Lst_Commands != nil {
		startVariables, setters, templateString, _ := Template1Builder(executionTemplate)
		strTemplate = templateString
		startVariableDeclarations = append(startVariableDeclarations, startVariables...)
		setterList = append(setterList, setters...)
	} else {
		startVariable, setter, templateString, _ := Template2Builder(executionTemplate)
		strTemplate = templateString
		startVariableDeclarations = append(startVariableDeclarations, startVariable)
		setterList = append(setterList, setter)
	}

	return startVariableDeclarations, setterList, strTemplate, nil
}