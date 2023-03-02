package dbCollectionHandler

import (
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

// Check the contracts type and update the database
func UpdateCollectionsWithNewStatus(pendingContract model.PendingContracts, status string) error {
	object := dao.Connection{}

	// Update main collections
	if pendingContract.ContractType == "ETHEXPERTFORMULA" {
		// unique identifier -> formula ID
		errWhenUpdatingTheFormula := object.UpdateEthFormulaStatusByUUID(pendingContract.Identifier, status, pendingContract.ErrorMessage)
		if errWhenUpdatingTheFormula != nil {
			return errWhenUpdatingTheFormula
		}
		logrus.Info("Updated the formula status: ", status)
	} else if pendingContract.ContractType == "ETHMETRICBIND" {
		// unique identifier -> uuid
		errWhenUpdatingTheMetric := object.UpdateEthMetricStatusByUUID(pendingContract.Identifier, status, pendingContract.ErrorMessage)
		if errWhenUpdatingTheMetric != nil {
			return errWhenUpdatingTheMetric
		}
		logrus.Info("Updated the metric status: ", status)
	}

	// update pending transactions collection
	pendingContract.Status = status
	errWhenUpdatingThePendingContracts := object.UpdateEthereumPendingContract(pendingContract.TransactionHash, pendingContract.ContractAddress, pendingContract.Identifier, pendingContract)
	if errWhenUpdatingThePendingContracts != nil {
		return errWhenUpdatingThePendingContracts
	}

	return nil
}