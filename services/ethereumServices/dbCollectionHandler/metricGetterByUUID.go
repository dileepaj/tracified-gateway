package dbCollectionHandler

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

func GetEthMetricByUUID(identifier string) (model.EthereumMetricBind, error) {
	object := dao.Connection{}

	metricBind, errWhenRetrievingTheMetric := object.GetEthMetricByUUID(identifier).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errWhenRetrievingTheMetric != nil {
		logrus.Error("error when retrieving the metric bind for the given uuid: " + errWhenRetrievingTheMetric.Error())
		return model.EthereumMetricBind{}, errors.New("error when retrieving the metric bind for the given uuid: " + errWhenRetrievingTheMetric.Error())
	}
	if metricBind != nil {
		metric := metricBind.(model.EthereumMetricBind)
		return metric, nil
	} else {
		logrus.Error("no metric bind found for the given uuid")
		return model.EthereumMetricBind{}, errors.New("no metric bind found for the given uuid")
	}
}

func GetEthFormulaByUUID(identifier string, collection string) (model.EthereumExpertFormula, error) {
	object := dao.Connection{}

	formula, errWhenRetrievingTheFormula := object.GetEthFormulaBinAndAbiByIdentifier(identifier, collection).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errWhenRetrievingTheFormula != nil {
		logrus.Error("error when retrieving the formula for the given uuid: " + errWhenRetrievingTheFormula.Error())
		return model.EthereumExpertFormula{}, errors.New("error when retrieving the formula for the given uuid: " + errWhenRetrievingTheFormula.Error())
	}
	if formula != nil {
		formulaObj := formula.(model.EthereumExpertFormula)
		return formulaObj, nil
	} else {
		logrus.Error("no formula found for the given uuid")
		return model.EthereumExpertFormula{}, errors.New("no formula found for the given uuid")
	}
}
