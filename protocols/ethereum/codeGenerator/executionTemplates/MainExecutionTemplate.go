package executionTemplates

import "github.com/dileepaj/tracified-gateway/model"

func ExecutionTemplateDivider(executionTemplate model.ExecutionTemplate) (string, error) {
	var strTemplate string

	if executionTemplate.Lst_Commands != nil {
		strTemplate, _ = Template1Builder(executionTemplate)
	} else {
		strTemplate, _ = Template2Builder(executionTemplate)
	}

	return strTemplate, nil
}