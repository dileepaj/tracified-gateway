package expertformula

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dileepaj/tracified-gateway/model"
	fclqueryexecuter "github.com/dileepaj/tracified-gateway/protocols/stellarprotocols/FCLQueryExecuter"
	fcl "github.com/shanukabps/FCL-GO_test"
	"github.com/sirupsen/logrus"
)

/**
 * Get the execution template as a string by passing the query string to the FCL query executer
 * and convert the response to a JSON and store it as a ExecutionTemplate struct
 * @param query - string
 * @return ExecutionTemplate - model.ExecutionTemplate
 */

func BuildJSONStructure(executionTemplateString string) (model.ExecutionTemplate, error) {
	var result string = fcl.NewFCLWrapper().GetExecutionTemplateJSONString("./Defs.txt", "$WATER.Multiply($WATER_TO_ELECTRICITY_UNIT).Multiply($ELECTRICITY_UNIT_TO_CARBON_EMISSION)")
	fmt.Println(result)

	var executionTemplate model.ExecutionTemplate
	errWhenUnmarshelling := json.Unmarshal([]byte(executionTemplateString), &executionTemplate)
	if errWhenUnmarshelling != nil {
		return model.ExecutionTemplate{}, errors.New("error when unmarshelling string to JSON execution template object " + errWhenUnmarshelling.Error())
	}
	// logrus.Info("Exection Template created : ", "%+v\n", executionTemplate)
	fmt.Printf("%+v\n", executionTemplate)

	return executionTemplate, nil
}

func BuildExecutionTemplateByQuery(query string) (model.ExecutionTemplate, error) {
	var executionTemplate model.ExecutionTemplate
	executionTemplateString, err := fclqueryexecuter.FCLQueryToExecutionTempalteJsonString(query)
	if err != nil {
		return model.ExecutionTemplate{}, errors.New("error when getting execution template string from FCL query executer(queryJSONBuilder) " + err.Error())
	}
	errWhenUnmarshelling := json.Unmarshal([]byte(executionTemplateString), &executionTemplate)
	if errWhenUnmarshelling != nil {
		return model.ExecutionTemplate{}, errors.New("error when unmarshelling string to JSON execution template object(queryJSONBuilder) " + errWhenUnmarshelling.Error())
	}
	logrus.Printf("%+v\n", executionTemplate)
	if executionTemplate.Error != "" {
		return model.ExecutionTemplate{}, errors.New("error when getting execution template string from FCL query executer(queryJSONBuilder)" + executionTemplate.Error)
	}
	return executionTemplate, nil
}
