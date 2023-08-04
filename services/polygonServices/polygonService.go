package polygonservices

import (
	"math/big"

	"github.com/dileepaj/tracified-gateway/model"
	polygoncontractdeployer "github.com/dileepaj/tracified-gateway/services/polygonServices/polygonContractDeployer"
)

type PolygonDeploymentStrategy interface {
	PolygonAbstractContractDeployment() (string, string, string, *big.Int, *big.Int, int, error)
}

type PolygonContractDeployerContext struct {
	contractDeploymentStrategy PolygonDeploymentStrategy
}

func (context *PolygonContractDeployerContext) SetContractDeploymentStrategyForPolygon(strategy PolygonDeploymentStrategy) {
	context.contractDeploymentStrategy = strategy
}

// creating a method to execute the strategy
func (context *PolygonContractDeployerContext) PolygonExecuteContractDeployment() (string, string, string, *big.Int, *big.Int, int, error) {
	return context.contractDeploymentStrategy.PolygonAbstractContractDeployment()
}

type PolygonAbstractContractDeployment struct {
	ABI          string
	BIN          string
	Identifier   string
	ContractType string
	OtherParams  []any
}

func (contractObject *PolygonAbstractContractDeployment) PolygonAbstractContractDeployment() (string, string, string, *big.Int, *big.Int, int, error) {
	//call the deployer method that is able to send the transaction to the blockchain with multiple try outs on failures
	address, hash, transactionCost, erInContractDeployment := polygoncontractdeployer.PolygonContractDeployer(contractObject.BIN, contractObject.ABI, contractObject.Identifier, contractObject.ContractType, contractObject.OtherParams)
	return address, hash, transactionCost, big.NewInt(0), big.NewInt(0), 0, erInContractDeployment
}

// struct and implementation for re-deployment strategy
type PolygonAbstractContractRedeployment struct {
	PendingContract model.PendingContracts
}

func (contractObject *PolygonAbstractContractRedeployment) PolygonAbstractContractDeployment() (string, string, string, *big.Int, *big.Int, int, error) {
	//call the deployer method that is able to send the transaction to the blockchain with multiple try outs on failures
	address, hash, transactionCost, nonce, gasPrice, gasLimit, errInContractDeployment := polygoncontractdeployer.RedeployFailedContractsOnPolygon(contractObject.PendingContract)
	return address, hash, transactionCost, nonce, gasPrice, gasLimit, errInContractDeployment
}
