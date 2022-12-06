package ethereuemmetricbind

import (
	"strconv"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/ethereum"
)

/*
	TODO:
		* check the request data with variables in the contract
*/

func WriteAddValue(value model.FormulaDetails, valueCount int) (string, error) {
	addValueStr := ``

	addValueStr += "\t\t addValue(" + `"` + value.ID + `", "` + value.Description + `", "` + value.BindManageData.BindData.WorkflowID + `", "` + value.BindManageData.BindData.StageID + `", "` + value.BindManageData.Master.ArtifactTemplateName.ManageDataType + `", ` + strconv.Itoa(value.Type) + `, "` + value.BindManageData.Master.ArtifactID + `", "` + ethereum.StringToHexString(value.Key) + `", "` + ethereum.StringToHexString(value.ArtifactTemplateID) + `", "` + ethereum.StringToHexString(value.BindManageData.Master.ArtifactFieldKey.FieldKey) + `", "` + ethereum.StringToHexString(value.ArtifactTemplate.FieldName) + `");` + "\n"

	return addValueStr, nil
}