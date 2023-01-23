package executionTemplates

import "github.com/dileepaj/tracified-gateway/model"

/**
 * @return start variable declarations, setter list, built equation, error
 * @param executionTemplate
 */
func ExecutionTemplateDivider(executionTemplate model.ExecutionTemplate) (string, error) {
	var strTemplate string						// the final equation

	// check whether the execution template has list of commands or not and call the relevant template builder
	if executionTemplate.Lst_Commands != nil {
		templateString, _ := Template1Builder(executionTemplate)
		strTemplate = templateString
	} else {
		templateString, _ := Template2Builder(executionTemplate)
		strTemplate = templateString
	}

	return strTemplate, nil
}