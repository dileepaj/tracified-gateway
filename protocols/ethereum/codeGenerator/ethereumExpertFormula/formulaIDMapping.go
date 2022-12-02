package codeGenerator

import (
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

func MapFormulaID(formulaID string, status string) (uint64, error) {
	var formulaMapID uint64
	object := dao.Connection{}

	// check if the formula ID is already mapped
	formulaMap, err := object.GetEthFormulaMapID(formulaID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err != nil {
		logrus.Info(err)
	}
	// if the formula ID is already mapped, return the mapped ID
	// else map the formula ID and return the mapped ID (if only the status is empty)
	if formulaMap != nil {
		logrus.Info("Formula ID " + formulaID + " is recorded in the DB")
		formulaMapData := formulaMap.(model.EthFormulaIDMap)
		formulaMapID = formulaMapData.MapID
	} else if status == "" {
		data, errWhenMapping := object.GetNextSequenceValue("ETHFORMULAID")
		if errWhenMapping != nil {
			return data.SequenceValue, errWhenMapping
		}

		formulaMapID = data.SequenceValue
		formulaIDMap := model.EthFormulaIDMap{
			FormulaID: formulaID,
			MapID:     formulaMapID,
		}
		// map the formulaID with incrementing Integer put those object to blockchain
		err1 := object.InsertEthFormulaIDMap(formulaIDMap)
		if err1 != nil {
			logrus.Error("Inserting formula to the expert formula id map was failed " + err1.Error())
			return 0, err1
		}
	}

	return formulaMapID, nil
}
