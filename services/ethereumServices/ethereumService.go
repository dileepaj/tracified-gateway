package ethereumservices

import (
	contractdeployer "github.com/dileepaj/tracified-gateway/services/ethereumServices/contractDeployer"
	"github.com/sirupsen/logrus"
)

type AbstractContractDeployment struct {
	ABI string
	BIN string
}

func (contractObject *AbstractContractDeployment) AbstractContractDeployer() {
	//call the deployer method that is able to send the transaction to the blockchain with multiple try outs on failures
	address, hash, transactionCost, errInContractDeploymet := contractdeployer.EthereumContractDeployerService(contractObject.BIN, contractObject.ABI)
	logrus.Info(address, hash, transactionCost, errInContractDeploymet)

}
