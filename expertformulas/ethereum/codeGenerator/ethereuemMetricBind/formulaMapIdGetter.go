package ethereuemmetricbind

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

// to get the metric map id from the DB

func GetFormulaMapId(formulaID string) (uint64, error) {
	var formulaMapID uint64
	object := dao.Connection{}

	formulaIDMap, errInFormulaIDMap := object.GetEthFormulaMapID(formulaID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errInFormulaIDMap != nil {
		logrus.Info("Error when retrieving formula id from DB. Error: " + errInFormulaIDMap.Error())
		return formulaMapID, errors.New("Error when retrieving formula id from DB. Error: " + errInFormulaIDMap.Error())
	}
	if formulaIDMap == nil {
		logrus.Error("Requested map id for formula " + formulaID + " does not exists in the gateway DB")
		return formulaMapID, errors.New("Requested map id for formula " + formulaID + " does not exists in the gateway DB")
	} else {
		logrus.Info("Formula ID " + formulaID + " is recorded in the DB")
		formulaIDMapData := formulaIDMap.(model.EthFormulaIDMap)
		formulaMapID = formulaIDMapData.MapID
	}
	return formulaMapID, nil
}
