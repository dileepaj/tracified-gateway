package activitywriters

import (
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/model"
)

// Get the contract address of the formula
// Get the number of values in the formula
// Get the value IDs of the formula
// Convert the activity name to base64
// Create the formula definition code and return it

func GetFormulaDefinitionCode(element model.MetricDataBindActivityRequest) (string, error) {
	// Get the contract address of the formula
	contractAddress, errInGettingContractAddress := GetFormulaContractAddress(element.MetricFormula.MetricExpertFormula.ID)
	if errInGettingContractAddress != nil {
		return "", errInGettingContractAddress
	}

	// Getting value IDs string
	valueIDs := []string{}
	for _, value := range element.MetricFormula.Formula {
		valueIDs = append(valueIDs, value.ID)
	}
	valueIDsString := strings.Join(valueIDs, ", ")

	formulaDefinitionStart := "\t" + `Formula private formula = Formula(`
	formulaIDCode := `"` + element.MetricFormula.MetricExpertFormula.ID + `", `
	contractAddressCode := `"` + contractAddress + `", `
	noOfValuesCode := strconv.Itoa(len(element.MetricFormula.Formula)) + `, `
	activityIDCode := `"` + element.ID + `", `
	activityNameCode := `"` + base64.StdEncoding.EncodeToString([]byte(element.Name)) + `", `
	valueIDsCode := `"` + valueIDsString + `"`
	formulaDefinitionEnd := `);`
	formulaDefinitionCode := formulaDefinitionStart + formulaIDCode + contractAddressCode + noOfValuesCode + activityIDCode + activityNameCode + valueIDsCode + formulaDefinitionEnd

	return formulaDefinitionCode, nil
}
