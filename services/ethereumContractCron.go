package services

import (
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/vendor/github.com/ethereum/go-ethereum/log"
)

func CheckContractStatus() {
	//TODO call the DB and filter out the transaction with pending status
	//TODO loop through the transactions and call the ethereum endpoint to check the transaction status
	//TODO check the status
	/*
		TODO
			if pending - Update the index
			if success - update in the DB collection as completed
			if failed - log the error and mark the status as failed
	*/

	if commons.GoDotEnvVariable("LOGSTYPE") == "DEBUG" {
		log.Debug("---------------------------------------- Check pending Ethereum contracts -----------------------")
	}

	object := dao.Connection{}
	p := object.GetPendingContractsByStatus("Pending")
	p.Then(func(data interface{}) interface{} {
		result := data.([]model.PendingContracts)
		for i := 0; i < len(result); i++ {

		}
		return nil
	}).Catch(func(error error) error {
		if commons.GoDotEnvVariable("LOGSTYPE") == "DEBUG" {
			log.Error("Error @CheckContractStatus " + error.Error())
		}
		return error
	})
	p.Await()
}
