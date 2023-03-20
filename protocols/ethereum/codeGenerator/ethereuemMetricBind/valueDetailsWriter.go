package ethereuemmetricbind

import (
	"errors"
	"strconv"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/ethereum"
	"github.com/sirupsen/logrus"
)

/*
 * This function is used to write the addValue function calls
 */

func WriteAddValue(formulaId string, value model.FormulaDetails, valueCount int, stageID string, stageName string, workflowID string, pivotFields []model.PivotField) (string, error) {
	addValueStr := ``
	bindDataType := 0
	object := dao.Connection{}

	// 1 - stage data, 2 - master data
	if value.ArtifactTemplateID != "" {
		bindDataType = 2
	} else {
		bindDataType = 1
	}

	// get the value name from DB
	valueName := ""
	variableDefMap, errWhenGettingVariableData := object.EthereumGetValueMapID(value.ID, formulaId).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errWhenGettingVariableData != nil {
		logrus.Error("Unable to connect to gateway datastore ", errWhenGettingVariableData)
	}
	if variableDefMap == nil {
		logrus.Error("Formula Id ", formulaId)
		logrus.Error("Requested variable " + value.ID + " does not exists in the gateway DB")
		return "", errors.New("Requested variable " + value.ID + " does not exists in the gateway DB")
	} else {
		valueMapData := variableDefMap.(model.ValueIDMap)
		valueName = valueMapData.ValueName
	}

	// get the primary key row ID
	primaryKeyRowID := ""
	if len(pivotFields) > 0 {
		for _, pivot := range pivotFields {
			if value.ArtifactTemplateID == pivot.ArtifactTemplateId && value.ArtifactTemplate.FieldName == pivot.Field && pivot.Condition == "EQUAL" {
				primaryKeyRowID = pivot.ArtifactDataId
			}
		}
	}

	// add the addValue function call string
	addValueStr += "\t\tallValues.push(Value(" + `"` + value.ID + `", "` +
		valueName + `", "` +
		workflowID + `", "` +
		stageID + `", "` +
		ethereum.StringToHexString(stageName) + `", "` +
		ethereum.StringToHexString(value.Key) + `", "` +
		strconv.Itoa(value.Type) + `", ` +
		strconv.Itoa(bindDataType) + `, "` +
		value.ArtifactTemplate.ID + `", "` +
		primaryKeyRowID + `", "` +
		ethereum.StringToHexString(value.ArtifactTemplate.Name) + `", "` +
		ethereum.StringToHexString(value.Field) + `", "` +
		ethereum.StringToHexString(value.ArtifactTemplate.FieldName) + `"));` + "\n"

	return addValueStr, nil
}
