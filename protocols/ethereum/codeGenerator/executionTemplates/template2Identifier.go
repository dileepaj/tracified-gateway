package executionTemplates

import (
	"github.com/dileepaj/tracified-gateway/model"
)

/**
 * @return start variable declaration, the setter for the start variable, variable name, error
 * @param executionTemplate
 */
func Template2Builder(executionTemplate model.ExecutionTemplate) (string, error) {

	return executionTemplate.S_StartVarName, nil
}