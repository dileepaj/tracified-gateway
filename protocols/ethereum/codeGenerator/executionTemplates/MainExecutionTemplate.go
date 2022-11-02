package executionTemplates

import "github.com/dileepaj/tracified-gateway/model"

func ExecutionTemplateDivider(executionTemplate model.ExecutionTemplate) ([]string, []string, string, error) {
	var strTemplate string
	var startVariableDeclarations []string
	var setterList []string

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