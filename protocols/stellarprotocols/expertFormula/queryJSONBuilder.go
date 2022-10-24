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
		return model.ExecutionTemplate{}, err
	}
	errWhenUnmarshelling := json.Unmarshal([]byte(executionTemplateString), &executionTemplate)
	if errWhenUnmarshelling != nil {
		return model.ExecutionTemplate{}, errors.New("error when unmarshelling string to JSON execution template object " + errWhenUnmarshelling.Error())
	}
	logrus.Printf("%+v\n", executionTemplate)
	if executionTemplate.Error != "" {
		return model.ExecutionTemplate{}, errors.New(executionTemplate.Error)
	}
	return executionTemplate, nil
}
