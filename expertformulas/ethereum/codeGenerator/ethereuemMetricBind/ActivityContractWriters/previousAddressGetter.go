package ActivityContractWriters

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

// to retrieve the previous contract address from the DB

func getPreviousContractAddress(metricID string) (string, error) {
	contractAddress := ""

	object := dao.Connection{}
	contract, err := object.GetEthereumMetricLatestContract(metricID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err != nil {
		logrus.Error("Unable to connect to gateway datastore ", err)
		return "", errors.New("Requested latest contract address for metric " + metricID + " does not exists in the gateway DB")
	}
	if contract == nil {
		logrus.Error("Requested latest contract address for metric " + metricID + " does not exists in the gateway DB")
		return "", errors.New("requested latest contract address for metric " + metricID + " does not exists in the gateway DB")
	} else {
		latestContract := contract.(model.MetricLatestContract)
		contractAddress = latestContract.ContractAddress
	}

	return contractAddress, nil

}
