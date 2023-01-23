package ethereuemmetricbind

import (
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/ethereum"
)

/*
 * This function is used to write the addFormula() method call
 */

func WriteAddFormula(activity model.MetricDataBindActivityRequest, formulaCount int) ([]string, string, error) {
	addFormulaStr := ``
	object := dao.Connection{}

	// get the contract address of the formula from DB
	contract := ""
	formulaDet, errInGettingFormulaDet := object.GetEthFormulaStatus(activity.MetricFormula.MetricExpertFormula.ID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errInGettingFormulaDet != nil {
		return []string{}, ``, errInGettingFormulaDet
	}
	formulaDetData := formulaDet.(model.EthereumExpertFormula)
	contract = formulaDetData.ContractAddress

	valueIDs := ""
	valueCount := 0
	valueIDArray := []string{} // to store the value IDs
	addValueStrings := ``      // to get all the addValues() method strings of the formula

	// loop through all the values and get the method call string
	for _, value := range activity.MetricFormula.Formula {
		valueCount++
		addValueStr, errInGettingValueDetails := WriteAddValue(activity.MetricFormula.MetricExpertFormula.ID, value, valueCount, activity.Stage.StageID, activity.Stage.Name, activity.WorkflowID, activity.MetricFormula.PivotFields)
		if errInGettingValueDetails != nil {
			return []string{}, ``, errInGettingValueDetails
		}
		valueIDArray = append(valueIDArray, value.ID)

		// add the codes for adding values if the addValueStr is not empty
		if addValueStr != `` {
			valueComment := "\t\t// add value " + strconv.Itoa(valueCount) + " for formula " + strconv.Itoa(formulaCount) + "\n"
			addValueStrings += valueComment + addValueStr
		}
	}
	// get the value IDs a string
	valueIDs = strings.Join(valueIDArray, ", ")

	// add the addFormula() method call string
	addFormulaStr += "\t\tallFormulas.push(Formula(" + `"` +
		activity.MetricFormula.MetricExpertFormula.ID + `", "` +
		contract + `", ` +
		strconv.Itoa(len(activity.MetricFormula.Formula)) + `, "` +
		activity.ID + `", "` +
		ethereum.StringToHexString(activity.Name) + `", "` +
		valueIDs + `"));` + "\n"

	addFormulaStr += addValueStrings

	return valueIDArray, addFormulaStr, nil
}
