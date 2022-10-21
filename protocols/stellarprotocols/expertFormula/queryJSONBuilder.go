package expertformula

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dileepaj/tracified-gateway/model"
)

func BuildJSONStructure(executionTemplateString string) (model.ExecutionTemplate, error) {

	var executionTemplate model.ExecutionTemplate
	errWhenUnmarshelling := json.Unmarshal([]byte(executionTemplateString), &executionTemplate)
	if errWhenUnmarshelling != nil {
		return model.ExecutionTemplate{}, errors.New("error when unmarshelling string to JSON execution template object " + errWhenUnmarshelling.Error())
	}
	// logrus.Info("Exection Template created : ", "%+v\n", executionTemplate)
	fmt.Printf("%+v\n", executionTemplate)

	return executionTemplate, nil

}
