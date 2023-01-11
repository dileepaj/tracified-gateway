package ethereuemmetricbind

import (
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/ethereum"
)

func WriteAddPivotField(element model.MetricDataBindActivityRequest) (string, error) {
	str := "" // to store the method string

	// loop through all the pivot fields and add the method calls
	for _, pivotField := range element.MetricFormula.PivotFields {
		// add the pivot details to the string
		str += "\t\t" + `allPivotFields.push(PivotField("` +
			pivotField.Name + `", "` +
			ethereum.StringToHexString(pivotField.Key) + `", "` +
			pivotField.Field + `", "` +
			pivotField.Condition + `", "` +
			pivotField.Value + `", "` +
			pivotField.ArtifactTemplateId + `", "` +
			pivotField.ArtifactDataId + `", "` +
			element.MetricFormula.MetricExpertFormula.ID + `"));` + "\n"
	}

	return str, nil
}
