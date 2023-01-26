package activitywriters

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

// Retrieve formula contract address from DB

func GetFormulaContractAddress(formulaID string) (string, error) {
	// Get the contract address of the formula
	object := dao.Connection{}

	contract := ``

	// get the contract address of the formula from DB
	formulaDet, errInGettingFormulaDet := object.GetEthFormulaStatus(formulaID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errInGettingFormulaDet != nil {
		logrus.Error("Unable to connect to gateway datastore to get the formula contract address. Error: ", errInGettingFormulaDet)
		return "", errors.New("Unable to connect to gateway datastore to get the formula contract address. Error: " + errInGettingFormulaDet.Error())
	}
	if formulaDet == nil {
		logrus.Error("Requested contract address for formula " + formulaID + " does not exists in the gateway DB")
		return "", errors.New("requested contract address for formula " + formulaID + " does not exists in the gateway DB")
	} else {
		formulaDetData := formulaDet.(model.EthereumExpertFormula)
		contract = formulaDetData.ContractAddress
	}
	

	return contract, nil
}