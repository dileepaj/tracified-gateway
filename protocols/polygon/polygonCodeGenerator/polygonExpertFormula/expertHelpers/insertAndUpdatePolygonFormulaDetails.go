package experthelpers

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/utilities"
)

func InsertAndUpdateExpertFormulaDetailsToPolygonCollections(deployStatus int, status int, errorMsg string, actualStatus int, formulaObj model.EthereumExpertFormula, formulaId string, transactionUUID string) error {
	logger := utilities.NewCustomLogger()
	object := dao.Connection{}

	formulaObj.Status = status
	formulaObj.ErrorMessage = errorMsg
	formulaObj.ActualStatus = actualStatus

	if deployStatus == 0 {
		//Insert to collection
		errWhenInsertingFormulaDetails := object.InsertToPolygonFormulaDetails(formulaObj)
		if errWhenInsertingFormulaDetails != nil {
			logger.LogWriter("Error while inserting formula details to Polygon collection : "+errWhenInsertingFormulaDetails.Error(), constants.ERROR)
			return errors.New("Error while inserting formula details to Polygon collection : " + errWhenInsertingFormulaDetails.Error())
		}
	} else if deployStatus == 119 {
		//Update the collection
		errWhenUpdatingFormulaDetails := object.UpdatePolygonFormulaStatus(formulaId, transactionUUID, formulaObj)
		if errWhenUpdatingFormulaDetails != nil {
			logger.LogWriter("Error while updating the formula details in Polygon collection : "+errWhenUpdatingFormulaDetails.Error(), constants.ERROR)
			return errors.New("Error while updating the formula details in Polygon collection : " + errWhenUpdatingFormulaDetails.Error())
		}
	} else {
		logger.LogWriter("Invalid deployment status recorded for the contract ", constants.ERROR)
		return errors.New("Invalid deployment status recorded for the contract")
	}

	return nil
}
