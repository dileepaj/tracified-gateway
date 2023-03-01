package ethereumservices

import (
	contractdeployer "github.com/dileepaj/tracified-gateway/services/ethereumServices/contractDeployer"
	gettercaller "github.com/dileepaj/tracified-gateway/services/ethereumServices/getterCaller"
	settercaller "github.com/dileepaj/tracified-gateway/services/ethereumServices/setterCaller"
	"github.com/sirupsen/logrus"
)

type AbstractContractDeployment struct {
	ABI string
	BIN string
	Identifier string
	ContractType string
}

type AbstractMethodExecution struct {
	IsWithParams    bool
	ABI             string
	BIN             string
	Method          string
	Values          []any //this should contain the parameters according to the correct order
	ContractAddress string
}

type AbstractGetters struct {
	IsWithFilter    bool
	ABI             string
	BIN             string
	Method          string
	Params          []any //this should contain the parameters according to the correct order
	ContractAddress string
}

func (contractObject *AbstractContractDeployment) AbstractContractDeployer() (string, string, string, error) {
	//call the deployer method that is able to send the transaction to the blockchain with multiple try outs on failures
	address, hash, transactionCost, errInContractDeployment := contractdeployer.EthereumContractDeployerService(contractObject.BIN, contractObject.ABI, contractObject.Identifier, contractObject.ContractType)
	return address, hash, transactionCost, errInContractDeployment
}

func (MethodExecuteObject *AbstractMethodExecution) AbstractMethodExecuter() {
	//Checking the type and calling the relevant method according to the availability of parameters
	if MethodExecuteObject.IsWithParams {
		//call the method execution method that supports method calling with parameters
		hash, errorWhenCallingSetter := settercaller.EthereumSetterCallerWithParams(MethodExecuteObject.Method, MethodExecuteObject.ABI, MethodExecuteObject.BIN, MethodExecuteObject.Values, MethodExecuteObject.ContractAddress)
		logrus.Info(hash, errorWhenCallingSetter)

	} else {
		//call the method execution method that does not support parameters
		hash, errorWhenCallingSetter := settercaller.EthereumSetterCaller(MethodExecuteObject.Method, MethodExecuteObject.ABI, MethodExecuteObject.BIN, MethodExecuteObject.ContractAddress)
		logrus.Info(hash, errorWhenCallingSetter)
	}
}

func (GetterObject *AbstractGetters) AbstractGetterExecutor() {
	//Checking the type and calling the relevant method according to the availability of query parameters
	if GetterObject.IsWithFilter {
		//call the method execution method that supports method calling with parameters
		returnedValue, errorWhenCallingGetter := gettercaller.EthereumGetterCallerWithFilter(GetterObject.Method, GetterObject.ABI, GetterObject.BIN, GetterObject.Params, GetterObject.ContractAddress)
		logrus.Info(returnedValue, errorWhenCallingGetter)
	} else {
		//call the method execution method that does not support parameters
		returnedValue, errorWhenCallingGetter := gettercaller.EthereumGetterCaller(GetterObject.Method, GetterObject.ABI, GetterObject.BIN, GetterObject.ContractAddress)
		logrus.Info(returnedValue, errorWhenCallingGetter)
	}
}
