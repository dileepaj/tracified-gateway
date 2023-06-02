package ethereumservices

import (
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/services/ethereumServices/dbCollectionHandler"
	"github.com/sirupsen/logrus"
)

// define the interface
type CollectionServiceStrategy interface {
	AbstractCollectionService() (error)
}

// create context class
type CollectionServiceContext struct {
	collectionServiceStrategy CollectionServiceStrategy
}

// create a method to set the strategy
func (context *CollectionServiceContext) SetCollectionServiceStrategy(strategy CollectionServiceStrategy) {
	context.collectionServiceStrategy = strategy
}

// create a method to execute the strategy
func (context *CollectionServiceContext) ExecuteCollectionService() (error) {
	return context.collectionServiceStrategy.AbstractCollectionService()
}

// struct and function for social impact main collection and pending contract collection updates
type SocialImpactMainCollectionUpdate struct {
	PendingContract model.PendingContracts
	Status          int
}
func (collectionUpdateObject *SocialImpactMainCollectionUpdate) AbstractCollectionService() (error) {
	// update the collection with the new status
	errorInCallingCollectionService := dbCollectionHandler.UpdateCollectionsWithNewStatus(collectionUpdateObject.PendingContract, collectionUpdateObject.Status)
	if errorInCallingCollectionService != nil {
		logrus.Info("Error when calling the status updater. Error: ", errorInCallingCollectionService)
		return errorInCallingCollectionService
	}

	return nil
}
