package dbCollectionHandler

import (
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
)

// Check the contracts type and update the database
func UpdateCollectionsWithNewStatus(pendingContract model.PendingContracts, status string) error {
	object := dao.Connection{}

	// Update main collections
	if pendingContract.ContractType == "ETHEXPERTFORMULA" {
		// unique identifier -> formula ID
		errWhenUpdatingTheFormula := object.UpdateEthFormulaStatusByUUID(pendingContract.Identifier, status)
		if errWhenUpdatingTheFormula != nil {
			return errWhenUpdatingTheFormula
		}
	} else if pendingContract.ContractType == "ETHMETRICBIND" {
		// unique identifier -> uuid
		errWhenUpdatingTheMetric := object.UpdateEthMetricStatusByUUID(pendingContract.Identifier, status)
		if errWhenUpdatingTheMetric != nil {
			return errWhenUpdatingTheMetric
		}
	}

	// update pending transactions collection
	pendingContract.Status = status
	errWhenUpdatingThePendingContracts := object.UpdateEthereumPendingContract(pendingContract.TransactionHash, pendingContract.ContractAddress, pendingContract.Identifier, pendingContract)
	if errWhenUpdatingThePendingContracts != nil {
		return errWhenUpdatingThePendingContracts
	}

	return nil
}