package ethereumsocialimpact

import (
	"errors"
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
	// call the deploy method
	metadataDeployObj := ethereumservices.AbstractContractDeployment{
		ABI:          ethMetricObj.ABIstring,
		BIN:          ethMetricObj.BINstring,
		Identifier:   ethMetricObj.TransactionUUID,
		ContractType: "ETHMETRICBIND",
	}
	address, txnHash, deploymentCost, errWhenDeploying := metadataDeployObj.AbstractContractDeployer()
	ethMetricObj.Timestamp = time.Now().String()
	ethMetricObj.ContractAddress = address
	ethMetricObj.TransactionHash = txnHash
	ethMetricObj.TransactionCost = deploymentCost
	ethMetricObj.Status = "PENDING"
	if errWhenDeploying != nil {
		ethMetricObj.ErrorMessage = errWhenDeploying.Error()
		ethMetricObj.Status = "FAILED"
		logrus.Error("Error when deploying contract for metric : ", ethMetricObj.MetricID, " Error : ", errWhenDeploying)

		errWhenUpdatingMetric := object.UpdateEthereumMetricStatus(ethMetricObj.MetricID, ethMetricObj.TransactionUUID, ethMetricObj)
		if errWhenUpdatingMetric != nil {
			logrus.Error("Error when updating the status of metric status for Eth , formula ID " + ethMetricObj.MetricID)
		}

		pendingContract := model.PendingContracts{
			ContractAddress: ethMetricObj.ContractAddress,
			ContractType:    "ETHMETRICBIND",
			TransactionHash: ethMetricObj.TransactionHash,
			Identifier:      ethMetricObj.TransactionUUID,
			Status:          "FAILED",
			ErrorMessage:    ethMetricObj.ErrorMessage,
		}
		errWheninvalidatingMetric := dbCollectionHandler.InvalidateMetric(pendingContract, ethMetricObj.Status, ethMetricObj.ErrorMessage)
		if errWheninvalidatingMetric != nil {
			logrus.Error("Error when invalidating the metric : " + ethMetricObj.MetricID)
		}
		logrus.Info("Metric update called with FAILED status. Type: " + ethMetricObj.Type)
		logrus.Info("Contract deployment unsuccessful")
		return errors.New("Error in metric contract deployer. Error : " + errWhenDeploying.Error())
	} else {
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