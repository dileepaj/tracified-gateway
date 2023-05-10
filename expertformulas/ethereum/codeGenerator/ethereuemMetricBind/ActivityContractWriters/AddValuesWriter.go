package ActivityContractWriters

import (
	"encoding/base64"
	"errors"
	"strconv"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

// For writing the code for adding values to the values array (inside the addDetails method)

func AddValuesWriter(elements model.MetricDataBindActivityRequest) (string, error) {
	addValueArrayString := ""
	bindDataType := 0
	object := dao.Connection{}

	for _, value := range elements.MetricFormula.Formula {
		// Get the value name from the DB
		valueName := ""
		variableDefMap, errWhenGettingVariableData := object.EthereumGetValueMapID(value.ID, elements.MetricFormula.MetricExpertFormula.ID).Then(func(data interface{}) interface{} {
			return data
		}).Await()
		if errWhenGettingVariableData != nil {
			logrus.Error("Unable to connect to gateway datastore ", errWhenGettingVariableData)
		}
		if variableDefMap == nil {
			logrus.Error("Requested variable " + value.ID + " does not exists in the gateway DB")
			return "", errors.New("Requested variable " + value.ID + " does not exists in the gateway DB")
		} else {
			valueMapData := variableDefMap.(model.ValueIDMap)
			valueName = valueMapData.ValueName
		}

		// Finding the binding type (1 - stage data, 2 - master data)
		if value.ArtifactTemplateID != "" {
			bindDataType = 2
		} else {
			bindDataType = 1
		}

		// Get the primary key row ID
		primaryKeyRowID := ""
		if len(elements.MetricFormula.PivotFields) > 0 {
			for _, pivot := range elements.MetricFormula.PivotFields {
				if value.ArtifactTemplateID == pivot.ArtifactTemplateId && value.ArtifactTemplate.FieldName == pivot.Field && pivot.Condition == "EQUAL" {
					primaryKeyRowID = pivot.ArtifactDataId
				}
			}
		}

		// Code for adding the value to the array
		valueAdder := "\t\t" + `allValues.push(Value(`
		valueAdder += `"` + value.ID + `",`
		valueAdder += `"` + valueName + `",`
		valueAdder += `"` + elements.WorkflowID + `",`
		valueAdder += `"` + elements.Stage.StageID + `",`
		valueAdder += `"` + base64.StdEncoding.EncodeToString([]byte(elements.Stage.Name)) + `",`
		valueAdder += `"` + base64.StdEncoding.EncodeToString([]byte(value.Key)) + `",`
		valueAdder += `"` + strconv.Itoa(value.Type) + `",`
		valueAdder += strconv.Itoa(bindDataType) + `,`
		valueAdder += `"` + primaryKeyRowID + `",`
		valueAdder += `"` + value.ArtifactTemplateID + `",`
		valueAdder += `"` + base64.StdEncoding.EncodeToString([]byte(value.ArtifactTemplate.Name)) + `",`
		valueAdder += `"` + base64.StdEncoding.EncodeToString([]byte(value.Field)) + `",`
		valueAdder += `"` + base64.StdEncoding.EncodeToString([]byte(value.ArtifactTemplate.FieldName)) + `"`
		valueAdder += `));` + "\n"
		addValueArrayString += valueAdder
	}

	addValuesCodeComment := "\t\t" + "// Add values to the values array" + "\n"
	addValuesCode := addValuesCodeComment + addValueArrayString + "\n"

	return addValuesCode, nil
}
