package polygoncontractdeployer

import (
	"errors"
	"math/big"

	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	generalservices "github.com/dileepaj/tracified-gateway/services/ethereumServices/generalServices"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/dileepaj/tracified-gateway/vendor/github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func PolygonContractDeployer(bin string, abi string, contractIdentifier string, contractType string, otherParams []any) (string, string, string, error) {
	contractAddress := ""
	transactionHash := ""
	transactionCost := ""
	logger := utilities.NewCustomLogger()
	object := dao.Connection{}
	var expertFormulaObj model.EthereumExpertFormula

	logger.LogWriter("Calling the polygon contract deployer service...........", constants.INFO)

	if contractType == "POLYGONEXPERTFORMULA" {
		expertFormulaObj = otherParams[0].(model.EthereumExpertFormula)
		expertFormulaObj.ActualStatus = 110 //DEPLOYMENT_STARTED
		errWhenUpdatingStatus := object.UpdateSelectedPolygonFormulaFields(expertFormulaObj.FormulaID, expertFormulaObj.TransactionUUID, expertFormulaObj)
		if errWhenUpdatingStatus != nil {
			logger.LogWriter("Error when updating polygon collection after the deployment started : "+errWhenUpdatingStatus.Error(), constants.ERROR)
			return contractAddress, transactionHash, transactionCost, errors.New("Error when updating polygon collection after the deployment started : " + errWhenUpdatingStatus.Error())
		}
	}

	//load client and key
	client, privateKey, fromAddress, errWhenLoadingClientAndKey := generalservices.LoadClientAndKey(2)
	if errWhenLoadingClientAndKey != nil {
		logger.LogWriter("Error when loading the client and the key : "+errWhenLoadingClientAndKey.Error(), constants.ERROR)
		return contractAddress, transactionHash, transactionCost, errors.New("Error when loading the client and the key : " + errWhenLoadingClientAndKey.Error())
	}

	logger.LogWriter(client, constants.INFO)
	logger.LogWriter(fromAddress, constants.INFO)

	ContractBIN, parsed, errWhenLoadingParsedABIAndBIN := generalservices.LoadContractBinAndParsedAbi(bin, abi)
	if errWhenLoadingParsedABIAndBIN != nil {
		logger.LogWriter("Error when loading ContractBIN and Parsed ABI : "+errWhenLoadingParsedABIAndBIN.Error(), constants.ERROR)
		return contractAddress, transactionHash, transactionCost, errors.New("Error when loading ContractBIN and Parsed ABI : " + errWhenLoadingParsedABIAndBIN.Error())
	}

	logger.LogWriter(ContractBIN, constants.INFO)

	if parsed == nil {
		logger.LogWriter("GetABI returned nil", constants.ERROR)
		return contractAddress, transactionHash, transactionCost, errors.New("Error when getting ABI string , ERROR : GetAbi() returned nil")
	}

	//create the keyed transactor
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Value = big.NewInt(0) // in wei

	return contractAddress, transactionHash, transactionCost, nil
}