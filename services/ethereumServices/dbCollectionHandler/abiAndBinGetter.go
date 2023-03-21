package dbCollectionHandler

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
)

//Get ABI and BIN strings
func GetAbiAndBin(contractType string, identifier string) (string, string, error) {
	object := dao.Connection{}

	if contractType == "ETHEXPERTFORMULA" {
		expertObject, errInRetrieving := object.GetEthFormulaBinAndAbiByIdentifier(identifier).Then(func(data interface{}) interface{} {
			return data
		}).Await()
		if errInRetrieving != nil {
			return "", "", errInRetrieving
		}
		if expertObject != nil {
			expert := expertObject.(model.EthereumExpertFormula)
			return expert.ABIstring, expert.BINstring, nil
		} else {
			return "", "", errors.New("no expert formula found to get the abi and bin")
		}
	} else if contractType == "ETHMETRICBIND" {
		metricObj, errInRetrieving := object.GetEthMetricBinAndAbiByIdentifier(identifier).Then(func(data interface{}) interface{} {
			return data
		}).Await()
		if errInRetrieving != nil {
			return "", "", errInRetrieving
		}
		if metricObj != nil {
			metric := metricObj.(model.EthereumMetricBind)
			return metric.ABIstring, metric.BINstring, nil
		} else {
			return "", "", errors.New("no metric found to get the abi and bin")
		}
	} else {
		return "", "", errors.New("no contract type found to get the abi and bin")
	}
}
