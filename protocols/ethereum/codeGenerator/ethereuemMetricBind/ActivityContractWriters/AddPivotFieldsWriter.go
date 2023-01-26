package activitywriters

import "github.com/dileepaj/tracified-gateway/model"

// For writing the code for adding pivot fields to the values array (inside the addDetails method)/*

func AddPivotFieldsWriter(elements []model.PivotField, formulaID string) (string, error){

	addPivotFieldArrayString := ""
	for _, pivot := range elements {
		pivotField := "\t\t" + `values.push(PivotField(`
		pivotField += `"` + pivot.Name + `", `
		pivotField += `"` + pivot.Key + `", `
		pivotField += `"` + pivot.Field + `", `
		pivotField += `"` + pivot.Condition + `", `
		pivotField += `"` + pivot.Value + `", `
		pivotField += `"` + pivot.ArtifactTemplateId + `", `
		pivotField += `"` + pivot.ArtifactDataId + `", `
		pivotField += `"` + formulaID + `"` + `));` + "\n"
		addPivotFieldArrayString += pivotField
	}

	addPivotFieldsCodeComment := "\t\t" + "// Add pivot fields to the pivot fields array" + "\n"
	addPivotFieldsCode := addPivotFieldsCodeComment + addPivotFieldArrayString + "\n"

	return addPivotFieldsCode, nil
}