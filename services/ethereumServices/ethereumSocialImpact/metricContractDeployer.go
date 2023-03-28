package ethereumsocialimpact

import (
	"time"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	ethereumservices "github.com/dileepaj/tracified-gateway/services/ethereumServices"
	"github.com/dileepaj/tracified-gateway/services/ethereumServices/dbCollectionHandler"
	"github.com/sirupsen/logrus"
)

func DeployMetricContract(ethMetricObj model.EthereumMetricBind) error {
	object := dao.Connection{}
	logrus.Info("Trying to deploy contract for metric : ", ethMetricObj.MetricID)
	ethMetricObj.ActualStatus = 110 // DEPLOYMENT_STARTED
	errorWhenUpdatingStatus := object.UpdateSelectedEthMetricFields(ethMetricObj.MetricID, ethMetricObj.TransactionUUID, ethMetricObj)
	if errorWhenUpdatingStatus != nil {
		logrus.Error("Error when updating the status of the metric : ", ethMetricObj.MetricID, " Error : ", errorWhenUpdatingStatus)
		return errorWhenUpdatingStatus
	}
	// use deployer strategy
	metricDeployer := &ethereumservices.ContractDeployerContext{}
	metricDeployer.SetContractDeploymentStrategy(&ethereumservices.AbstractContractDeployment{
		ABI: 	  ethMetricObj.ABIstring,
		BIN: 	  ethMetricObj.BINstring,
		Identifier: ethMetricObj.TransactionUUID,
		ContractType: "ETHMETRICBIND",
		OtherParams: []any{ethMetricObj},
	})

	// call the deploy method
	address, txnHash, deploymentCost, _, _, _, errWhenDeploying := metricDeployer.ExecuteContractDeployment()
	ethMetricObj.Timestamp = time.Now().UTC().String()
	ethMetricObj.ContractAddress = address
	ethMetricObj.TransactionHash = txnHash
	ethMetricObj.TransactionCost = deploymentCost
	ethMetricObj.Status = "PENDING"
	if errWhenDeploying != nil {
		ethMetricObj.ErrorMessage = errWhenDeploying.Error()
		ethMetricObj.Status = "FAILED"
		ethMetricObj.ActualStatus = 111 // DEPLOYMENT_FAILED
		errorWhenUpdatingStatus1 := object.UpdateSelectedEthMetricFields(ethMetricObj.MetricID, ethMetricObj.TransactionUUID, ethMetricObj)
		if errorWhenUpdatingStatus1 != nil {
			logrus.Error("Error when updating the status of the metric : ", ethMetricObj.MetricID, " Error : ", errorWhenUpdatingStatus1)
			return errorWhenUpdatingStatus1
		}
		logrus.Error("Error when deploying contract for metric : ", ethMetricObj.MetricID, " Error : ", errWhenDeploying)

		pendingContract := model.PendingContracts{
			ContractAddress: ethMetricObj.ContractAddress,
			ContractType:    "ETHMETRICBIND",
			TransactionHash: ethMetricObj.TransactionHash,
			Identifier:      ethMetricObj.TransactionUUID,
			Status:          "FAILED",
			ErrorMessage:    ethMetricObj.ErrorMessage,
		}
		errWheninvalidatingMetric := dbCollectionHandler.InvalidateMetric(pendingContract, "FAILED", ethMetricObj.ErrorMessage)
		if errWheninvalidatingMetric != nil {
			logrus.Error("Error when invalidating the metric : " + ethMetricObj.MetricID)
		}
		logrus.Info("Metric update called with FAILED status. Type: " + ethMetricObj.Type)
		logrus.Info("Contract deployment unsuccessful")
		return errWhenDeploying
	} else {
		ethMetricObj.ActualStatus = 112 // DEPLOYMENT_TRANSACTION_PENDING
		//if deploy method is success update the status into success
		errWhenUpdatingStatus := object.UpdateEthereumMetricStatus(ethMetricObj.MetricID, ethMetricObj.TransactionUUID, ethMetricObj)
		if errWhenUpdatingStatus != nil {
			logrus.Error("Error when updating the status of metric status for Eth , formula ID " + ethMetricObj.MetricID)
		}
		logrus.Info("Metric update called with status " + ethMetricObj.Status + ". Type: " + ethMetricObj.Type)
		logrus.Info("-------------------------------------------------------------------------------------------------------------------------------------")
		logrus.Info("Deployed expert metric bind smart contract to blockchain")
		logrus.Info("Contract address : " + address)
		logrus.Info("Transaction hash : " + txnHash)
		logrus.Info("-------------------------------------------------------------------------------------------------------------------------------------")
	}

	return nil
}