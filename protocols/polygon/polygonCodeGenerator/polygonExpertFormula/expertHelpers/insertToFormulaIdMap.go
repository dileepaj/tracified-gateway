package experthelpers

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/utilities"
)

//Get the next sequence for the formula ID map

//1 - Ethereum, 2 - Polygon
func InsertToFormulaIdMap(expertFormulaID string, blockchain int) error {
	object := dao.Connection{}
	var data model.Counters
	var errInGettingNextMapSequence error
	logger := utilities.NewCustomLogger()

	if blockchain == 1 {
		//get the sequence from ethereum formula ID map
		data, errInGettingNextMapSequence = object.GetNextSequenceValue("ETHFORMULAID")
	} else if blockchain == 2 {
		//get the sequence from polygon formula ID
		data, errInGettingNextMapSequence = object.GetNextSequenceValue("POLYGONFORMULAID")
	} else {
		logger.LogWriter("Invalid blockchain type", constants.INFO)
		return errors.New("Invalid blockchain type")
	}

	if errInGettingNextMapSequence != nil {
		logger.LogWriter("Error while getting next sequence value : "+errInGettingNextMapSequence.Error(), constants.ERROR)
		return errors.New("Error while getting next sequence value : " + errInGettingNextMapSequence.Error())
	}

	formulaIdMap := model.EthFormulaIDMap{
		FormulaID: expertFormulaID,
		MapID:     data.SequenceValue,
	}

	var errorWhenInsertingToFormulaIDMap error
	if blockchain == 1 {
		errorWhenInsertingToFormulaIDMap = object.InsertEthFormulaIDMap(formulaIdMap)
	} else if blockchain == 2 {
		errorWhenInsertingToFormulaIDMap = object.InsertPolygonFormulaIDMap(formulaIdMap)
	} else {
		logger.LogWriter("Invalid blockchain type", constants.INFO)
		return errors.New("Invalid blockchain type")
	}

	if errorWhenInsertingToFormulaIDMap != nil {
		logger.LogWriter("Unable to connect to gateway datastore : "+errorWhenInsertingToFormulaIDMap.Error(), constants.ERROR)
		return errors.New("Unable to connect to gateway datastore : " + errorWhenInsertingToFormulaIDMap.Error())
	}

	return nil

}
