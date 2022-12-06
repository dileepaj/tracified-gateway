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

func WriteAddFormula(activity model.MetricDataBindActivityRequest, formulaCount int) (string, error) {
	addFormulaStr := ``
	object := dao.Connection{}

	// get the mapped FormulaID
	formulaMapID, err := object.GetEthFormulaMapID(activity.MetricFormula.MetricExpertFormula.ID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err != nil {
		return ``, err
	}
	formulaMapData := formulaMapID.(model.EthFormulaIDMap)
	formulaMapIDInt := formulaMapData.MapID

	// get the contract address of the formula from DB
	contract := ""
	formulaDet, errInGettingFormulaDet := object.GetEthFormulaStatus(activity.MetricFormula.MetricExpertFormula.ID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errInGettingFormulaDet != nil {
		return ``, errInGettingFormulaDet
	}
	formulaDetData := formulaDet.(model.EthereumExpertFormula)
	contract = formulaDetData.ContractAddress

	valueIDs := ""
	valueCount := 0
	valueIDArray := []string{}	// to store the value IDs
	addValueStrings := ``		// to get all the addValues() method strings of the formula

	// loop through all the values and get the method call string
	for _, value := range activity.MetricFormula.Formula {
		valueCount++
		valueComment := "\t\t// add value " + strconv.Itoa(valueCount) + " for formula " + strconv.Itoa(formulaCount) + "\n"
		addValueStr, errInGettingValueDetails := WriteAddValue(activity.MetricFormula.MetricExpertFormula.ID, value, valueCount, activity.Stage.StageID, activity.Stage.Name, activity.WorkflowID, activity.MetricFormula.PivotFields)
		if errInGettingValueDetails != nil {
			return ``, errInGettingValueDetails
		}
		valueIDArray = append(valueIDArray, value.ID)

		addValueStrings += valueComment + addValueStr
	}
	// get the value IDs a string
	valueIDs = strings.Join(valueIDArray, ", ")

	// add the addFormula() method call string
	addFormulaStr += "\t\taddFormula(" + 
								strconv.FormatUint(uint64(formulaMapIDInt), 10) + `, "` + 
								activity.MetricFormula.MetricExpertFormula.ID + `", address(` + 
								contract + `), ` + 
								strconv.Itoa(len(activity.MetricFormula.Formula)) + `, "` + 
								activity.ID + `", "` + 
								ethereum.StringToHexString(activity.Name) + `", "` + 
								valueIDs + `");` + "\n"

	addFormulaStr += "\t\t// add formula id and contract address to array" + "\n"
	addFormulaStr += "\t\tformulaDetails.push(" + `'{` + 
												`"formulaId": ` + strconv.FormatUint(uint64(formulaMapIDInt), 10) + 
												`, "contractAddress": "` + contract +  
												`"}'` + ");" + "\n"
	addFormulaStr += addValueStrings

	return addFormulaStr, nil
}
