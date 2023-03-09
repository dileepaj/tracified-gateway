package dbCollectionHandler

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

func InvalidateMetric(pendingContract model.PendingContracts, status string, errorMessage string) error {
	object := dao.Connection{}

	metricID := ""
	// Retrieve the metric ID from the metric bind using UUID
	metricBind, errWhenRetrievingTheMetric := GetEthMetricByUUID(pendingContract.Identifier)
	if errWhenRetrievingTheMetric != nil {
		logrus.Error("error when retrieving the metric bind for the given uuid: " + errWhenRetrievingTheMetric.Error())
		return errors.New("error when retrieving the metric bind for the given uuid: " + errWhenRetrievingTheMetric.Error())
	}
	metricID = metricBind.MetricID
	

	// get the uuids for the metric and update the status
	metrics, errWhenRetrievingTheMetricUUIDs := object.GetEthMetricsByMetricID(metricID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errWhenRetrievingTheMetricUUIDs != nil {
		logrus.Error("error when retrieving the metric contracts for the given metric id: " + errWhenRetrievingTheMetricUUIDs.Error())
		return errors.New("error when retrieving the metric contracts for the given metric id: " + errWhenRetrievingTheMetricUUIDs.Error())
	}
	if metrics != nil {
		metricList := metrics.([]model.EthereumMetricBind)
		for _, metric := range metricList {

			if metric.TransactionUUID != pendingContract.Identifier {
				errorMessage = "One of the contracts under this metric ID is failed or invalid."
			}

			pendingContractNew := model.PendingContracts{}
			// get the pending contracts for the uuid if the hash is present
			if metric.TransactionHash != "" {
				contract, errWhenRetrievingThePendingContract := object.GetPendingContractByIdentifier(metric.TransactionUUID).Then(func(data interface{}) interface{} {
					return data
				}).Await()
				if errWhenRetrievingThePendingContract != nil {
					logrus.Error("error when retrieving the pending contract for the given uuid: " + errWhenRetrievingThePendingContract.Error())
					return errors.New("error when retrieving the pending contract for the given uuid: " + errWhenRetrievingThePendingContract.Error())
				}
				if contract != nil {
					pendingContractNew = contract.(model.PendingContracts)
					// update the status
					pendingContractNew.Status = status
					pendingContractNew.ErrorMessage = errorMessage
				}
			}

			errorWhenUpdatingCollections := UpdateCollectionsWithNewStatus(pendingContractNew, status)
			if errorWhenUpdatingCollections != nil {
				logrus.Error("error when updating the collections: " + errorWhenUpdatingCollections.Error())
				return errors.New("error when updating the collections: " + errorWhenUpdatingCollections.Error())
			}
		}

	} else {
		return errors.New("no metrics found for the given metric id")
	}

	logrus.Info("Invalidated all the contracts for the metric " + metricID + " with status: " + status)

	return nil
}
