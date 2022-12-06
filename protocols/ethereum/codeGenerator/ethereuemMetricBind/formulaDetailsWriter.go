package ethereuemmetricbind

import (
	"strconv"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/ethereum"
)

/*
	TODO:
		* get the mapped formula ID from the DB
		* get the contract address from the DB
		* get the valueIDs as a string
		* check the request data with variables in the contract
*/

func WriteAddFormula(formula model.MetricDataBindActivityRequest, formulaCount int) (string, error) {
	addFormulaStr := ``

	//object := dao.Connection{}
	// get the mapped FormulaID
	//formulaMapID, err := object.GetEthFormulaMapID(formula.MetricFormula.ID).Then(func(data interface{}) interface{} {
	//	return data
	//}).Await()
	//if err != nil {
	//	return ``, err
	//}
	//formulaMapData := formulaMapID.(model.EthFormulaIDMap)
	//formulaMapIDInt := formulaMapData.MapID

	// get the contract address of the formula
	// TOdo: get the contract address from the DB
	contract := "0xC0f4DC75c610bC621CB27c6616a44013B6888DDc"

	// get the valueIDs as a string
	valueIDs := ""

	addFormulaStr += "\t\t addFormula(" + strconv.FormatUint(uint64(formulaCount), 10) + `, "` + formula.MetricFormula.ID + `", address(` + contract + `), ` + strconv.Itoa(len(formula.MetricFormula.Formula)) + `, "` + formula.ActivityFormulaDefinitionManageData.ActivityID + `", "` + ethereum.StringToHexString(formula.ActivityNameMangeData.ActivityName) + `", "` + ethereum.StringToHexString(valueIDs) + `", "` + ethereum.StringToHexString(formula.Stage.Name) + `", "`+ ethereum.StringToHexString(formula.ActivityFormulaDefinitionManageData.Key) + `");` + "\n"

	valueCount := 0
	for _, value := range formula.MetricFormula.Formula {
		valueCount++
		valueComment := "\t\t // add value " + strconv.Itoa(valueCount) + " for formula " + strconv.Itoa(formulaCount) + "\n"
		addValueStr, errInGettingValueDetails := WriteAddValue(value, valueCount)
		if errInGettingValueDetails != nil {
			return ``, errInGettingValueDetails
		}

		addFormulaStr += valueComment + addValueStr
	}

	return addFormulaStr, nil
}
