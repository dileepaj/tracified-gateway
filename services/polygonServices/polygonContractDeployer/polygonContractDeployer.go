package polygoncontractdeployer

import (
	"errors"
	"math/big"
	"strconv"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	generalservices "github.com/dileepaj/tracified-gateway/services/ethereumServices/generalServices"
	polygongasservice "github.com/dileepaj/tracified-gateway/services/polygonServices/polygonGasService"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/dileepaj/tracified-gateway/vendor/github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func PolygonContractDeployer(bin string, abi string, contractIdentifier string, contractType string, otherParams []any) (string, string, string, error) {
	contractAddress := ""
	transactionHash := ""
	transactionCost := ""
	var isFailed = true
	var predictedGasLimit int
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

	tryoutCap, errInTryConvert := strconv.Atoi(commons.GoDotEnvVariable("CONTRACTDEPLOYLIMIT"))
	if errInTryConvert != nil {
		logger.LogWriter("Error when converting the tryout limit , ERROR : "+errInTryConvert.Error(), constants.ERROR)
		return contractAddress, transactionHash, transactionCost, errors.New("Error when converting the tryout limit , ERROR : " + errInTryConvert.Error())
	}

	gasLimitCap, errInGasLimitCapConcert := strconv.Atoi(commons.GoDotEnvVariable("GASLIMITCAP"))
	if errInGasLimitCapConcert != nil {
		logger.LogWriter("Error when converting the gas limit cap , ERROR : "+errInGasLimitCapConcert.Error(), constants.ERROR)
		return contractAddress, transactionHash, transactionCost, errors.New("Error when converting the gas limit cap , ERROR : " + errInGasLimitCapConcert.Error())
	}

	gasPriceCap, errInGasPriceCapConcert := strconv.Atoi(commons.GoDotEnvVariable("GASPRICECAP"))
	if errInGasPriceCapConcert != nil {
		logger.LogWriter("Error when converting the gas price cap , ERROR : "+errInGasPriceCapConcert.Error(), constants.ERROR)
		return contractAddress, transactionHash, transactionCost, errors.New("Error when converting the gas price cap , ERROR : " + errInGasPriceCapConcert.Error())
	}

	logger.LogWriter(tryoutCap, constants.INFO)
	logger.LogWriter(gasLimitCap, constants.INFO)
	logger.LogWriter(gasPriceCap, constants.INFO)

	for i := 0; i < tryoutCap; i++ {
		if !isFailed {
			return contractAddress, transactionHash, transactionCost, nil
		} else {
			logger.LogWriter("Deploying the contract for the "+strconv.FormatInt(int64(i+1), 10)+" th time", constants.INFO)
			if i == 0 {
				//gas limit estimation
				gasLimit, errWhenGettingGasLimit := polygongasservice.EstimateGasLimitForPolygon(commons.GoDotEnvVariable("ETHEREUMPUBKEY"), "", "", "", "", bin)
				if errWhenGettingGasLimit != nil {
					logger.LogWriter("Error when getting gas limit "+errWhenGettingGasLimit.Error(), constants.ERROR)
					return contractAddress, transactionHash, transactionCost, errors.New("Error when getting gas limit, ERROR : " + errWhenGettingGasLimit.Error())
				}
				predictedGasLimit = int(gasLimit)
				auth.GasLimit = uint64(predictedGasLimit) // in units

				//gas price estimation
			} else {

			}
		}
	}

	return contractAddress, transactionHash, transactionCost, nil
}
