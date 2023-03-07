package ethereumservices

import (
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/services/ethereumServices/dbCollectionHandler"
	"github.com/sirupsen/logrus"
)

type AbstractCollectionUpdate struct {
	PendingContract model.PendingContracts
	Status          string
	Type 		  	string
}

func (collectionUpdateObject *AbstractCollectionUpdate) AbstractCollectionUpdater() (error) {
	// Checking the type and calling the relevant method according to the availability of parameters
	if collectionUpdateObject.Type == "SOCIALIMPACT" {
		errorInCallingStatusUpdater := dbCollectionHandler.UpdateCollectionsWithNewStatus(collectionUpdateObject.PendingContract, collectionUpdateObject.Status)
		if errorInCallingStatusUpdater != nil {
			logrus.Info("Error when calling the status updater. Error: ", errorInCallingStatusUpdater)
			return errorInCallingStatusUpdater
		}
	}

	return nil
}
