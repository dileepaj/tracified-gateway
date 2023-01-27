package ethereuemmetricbind

import (
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/vendor/github.com/sirupsen/logrus"
)

// get the status of the metric smart contract from the DB by metric ID

func GetMetricSmartContractStatus(metricId string, contracttype string) (string, model.EthereumMetricBind, error) {
	object := dao.Connection{}
	var status string
	var metricObject model.EthereumMetricBind

	metricDetails, errWhenGettingMetricStatus := object.GetEthMetricStatus(metricId).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errWhenGettingMetricStatus != nil {
		logrus.Error("An error occurred when getting metric status ", errWhenGettingMetricStatus)
		return "", model.EthereumMetricBind{}, errWhenGettingMetricStatus
	}
	if metricDetails == nil {
		status = ""
	} else {
		metricObject = metricDetails.(model.EthereumMetricBind)
		status = metricObject.Status
	}

	return status, metricObject, nil
}