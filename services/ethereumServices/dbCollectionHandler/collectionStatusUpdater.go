package dbCollectionHandler

import (
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

// Check the contracts type and update the database
func UpdateCollectionsWithNewStatus(pendingContract model.PendingContracts, status int) error {
	object := dao.Connection{}

	// Update main collections
	if pendingContract.ContractType == "ETHEXPERTFORMULA" {
		// unique identifier -> uuid
		errWhenUpdatingTheFormula := object.UpdateEthFormulaStatusByUUID(pendingContract.Identifier, status, pendingContract.ErrorMessage)
		if errWhenUpdatingTheFormula != nil {
			return errWhenUpdatingTheFormula
		}
		logrus.Info("Updated the formula("+pendingContract.TransactionHash+") status: ", status)
	} else if pendingContract.ContractType == "ETHMETRICBIND" {
		// unique identifier -> uuid
		errWhenUpdatingTheMetric := object.UpdateEthMetricStatusByUUID(pendingContract.Identifier, status, pendingContract.ErrorMessage)
		if errWhenUpdatingTheMetric != nil {
			return errWhenUpdatingTheMetric
		}
		logrus.Info("Updated the metric("+pendingContract.TransactionHash+") status: ", status)
	}

	// update pending transactions collection if the hash is present
	if pendingContract.TransactionHash != "" {
		pendingContract.Status = status
		errWhenUpdatingThePendingContracts := object.UpdateEthereumPendingContract(pendingContract.TransactionHash, pendingContract.ContractAddress, pendingContract.Identifier, pendingContract)
		if errWhenUpdatingThePendingContracts != nil {
			return errWhenUpdatingThePendingContracts
		}
		logrus.Info("Updated the pending contract collection for "+pendingContract.TransactionHash+" status: ", status)
	}

	return nil
}

func UpdateCollectionsWithNewStatusForPolygon(pendingContract model.PendingContracts, status int) error {
	object := dao.Connection{}

	// Update main collections
	if pendingContract.ContractType == "POLYGONEXPERTFORMULA" {
		// unique identifier -> uuid
		errWhenUpdatingTheFormula := object.UpdatePolygonFormulaStatusByUUID(pendingContract.Identifier, status, pendingContract.ErrorMessage)
		if errWhenUpdatingTheFormula != nil {
			return errWhenUpdatingTheFormula
		}
		logrus.Info("Updated the formula("+pendingContract.TransactionHash+") status: ", status)
	} else if pendingContract.ContractType == "POLYGONMETRICBIND" {
		//TODO - UPDATE POLYGON METRIC
	}

	// update pending transactions collection if the hash is present
	if pendingContract.TransactionHash != "" {
		pendingContract.Status = status
		errWhenUpdatingThePendingContracts := object.UpdatePolygonPendingContract(pendingContract.TransactionHash, pendingContract.ContractAddress, pendingContract.Identifier, pendingContract)
		if errWhenUpdatingThePendingContracts != nil {
			return errWhenUpdatingThePendingContracts
		}
		logrus.Info("Updated the pending contract collection for "+pendingContract.TransactionHash+" status: ", status)
	}

	return nil
}
