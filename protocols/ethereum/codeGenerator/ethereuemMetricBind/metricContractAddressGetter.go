package ethereuemmetricbind

import (
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

var object = dao.Connection{}
var status string
var metricObject model.EthereumMetricBind

// get the status of the metric smart contract from the DB by metric ID

func GetMetricSmartContractStatus(metricId string, contractType string) (int, model.EthereumMetricBind, error) {
	var metricObject model.EthereumMetricBind 
	status := 0
	metricDetails, errWhenGettingMetricStatus := object.GetEthMetricByMetricIdAndType(metricId, contractType).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errWhenGettingMetricStatus != nil {
		logrus.Error("An error occurred when getting metric status ", errWhenGettingMetricStatus)
		return status, model.EthereumMetricBind{}, errWhenGettingMetricStatus
	}
	if metricDetails != nil {
		metricObject = metricDetails.(model.EthereumMetricBind)
		status = metricObject.Status
	}

	return status, metricObject, nil
}

func GetMetricSmartContractStatusForFormula(metricId string, contractType string, formulaId string) (int, model.EthereumMetricBind, error) {
	var metricObject model.EthereumMetricBind 
	status := 0
	metricFormulaDetails, errWhenGettingMetricStatus := object.GetEthMetricStatusForFormula(metricId, contractType, formulaId).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errWhenGettingMetricStatus != nil {
		logrus.Error("An error occurred when getting metric status ", errWhenGettingMetricStatus)
		return status, model.EthereumMetricBind{}, errWhenGettingMetricStatus
	}
	if metricFormulaDetails != nil {
		metricObject = metricFormulaDetails.(model.EthereumMetricBind)
		status = metricObject.Status
	}

	return status, metricObject, nil

}
