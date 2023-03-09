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
	
	return model.EthereumMetricBind{}, nil
}