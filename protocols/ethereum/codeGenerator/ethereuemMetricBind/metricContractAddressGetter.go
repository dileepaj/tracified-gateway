package ethereuemmetricbind

import (
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/vendor/github.com/sirupsen/logrus"
)

// get the status of the metric smart contract from the DB by metric ID

func GetMetricSmartContractStatus(metricId string, contractType string) (string, error) {
	object := dao.Connection{}
	var status string

	metricDetails, errWhenGettingMetricStatus := object.GetEthMetricByMetricIdAndType(metricId, contractType).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errWhenGettingMetricStatus != nil {
		logrus.Error("An error occurred when getting metric status ", errWhenGettingMetricStatus)
		return "", errWhenGettingMetricStatus
	}
	if metricDetails == nil {
		status = ""
	} else {
		status = metricDetails.(model.EthereumMetricBind).Status
	}

	return status, nil
}