package experthelpers

import (
	"encoding/base64"
	"errors"

	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/utilities"
)

func ContractTemplateBuilder(formulaObj model.EthereumExpertFormula, license string, pragmaLine string, importCalculationSol string, contractStart, contractBody string, contractEnd string) (string, error) {
	object := dao.Connection{}
	logger := utilities.NewCustomLogger()
	template := license + "\n\n" + pragmaLine + "\n\n" + importCalculationSol + "\n\n" + contractStart + "\n\t" + contractBody + "\n" + contractEnd
	b64Template := base64.StdEncoding.EncodeToString([]byte(template))

	formulaObj.TemplateString = b64Template
	formulaObj.ActualStatus = 103 // SMART_CONTRACT_GENERATION_COMPLETED

	errWhenUpdatingFormulaWithTemplateString := object.UpdateSelectedPolygonFormulaFields(formulaObj.FormulaID, formulaObj.TransactionUUID, formulaObj)
	if errWhenUpdatingFormulaWithTemplateString != nil {
		logger.LogWriter("Error while updating the polygon collection after the contract generation "+errWhenUpdatingFormulaWithTemplateString.Error(), constants.ERROR)
		return "", errors.New("Error while updating the polygon collection after the contract generation " + errWhenUpdatingFormulaWithTemplateString.Error())
	}

	return template, nil
}
