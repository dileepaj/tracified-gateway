package expertformula

import (
	"encoding/json"
	"errors"

	fclqueryexecuter "github.com/dileepaj/tracified-gateway/expertformulas/stellarprotocols/FCLQueryExecuter"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

/**
 * Get the execution template as a string by passing the query string to the FCL query executer
 * and convert the response to a JSON and store it as a ExecutionTemplate struct
 * @param query - string
 * @return ExecutionTemplate - model.ExecutionTemplate
 */

func BuildExecutionTemplateByQuery(query string) (model.ExecutionTemplate, error) {
	var executionTemplate model.ExecutionTemplate
	executionTemplateString, err := fclqueryexecuter.FCLQueryToExecutionTempalteJsonString(query)
	if err != nil {
		return model.ExecutionTemplate{}, errors.New("error when getting execution template string from FCL query executer " + err.Error())
	}
	errWhenUnmarshelling := json.Unmarshal([]byte(executionTemplateString), &executionTemplate)
	if errWhenUnmarshelling != nil {
		return model.ExecutionTemplate{}, errors.New("error when unmarshelling string to JSON execution template object " + errWhenUnmarshelling.Error())
	}
	logrus.Printf("%+v\n", executionTemplate)
	if executionTemplate.Error != "" {
		return model.ExecutionTemplate{}, errors.New("error when getting execution template string from FCL query executer" + executionTemplate.Error)
	}
	return executionTemplate, nil
}
