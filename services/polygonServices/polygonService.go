package polygonservices

import "math/big"

type PolygonDeploymentStrategy interface {
}

type PolygonAbstractContractDeployment struct {
	ABI          string
	BIN          string
	Identifier   string
	ContractType string
	OtherParams  []any
}

type PolygonContractDeployerContext struct {
	contractDeploymentStrategy PolygonDeploymentStrategy
}

func (context *PolygonContractDeployerContext) SetContractDeploymentStrategy(strategy PolygonDeploymentStrategy) {
	context.contractDeploymentStrategy = strategy
}

// creating a method to execute the strategy
func (context *PolygonAbstractContractDeployment) ExecuteContractDeployment() (string, string, string, *big.Int, *big.Int, int, error) {
	return context.PolygonAbstractContractDeployer()
}

func (contractObject *PolygonAbstractContractDeployment) PolygonAbstractContractDeployer() (string, string, string, *big.Int, *big.Int, int, error) {
	//call the deployer method that is able to send the transaction to the blockchain with multiple try outs on failures

	return "", "", "", big.NewInt(0), big.NewInt(0), 0, nil
}
