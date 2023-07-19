package polygoncontractdeployer

import (
	"context"
	"errors"
	"math/big"
	"strconv"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	contractdeployer "github.com/dileepaj/tracified-gateway/services/ethereumServices/contractDeployer"
	generalservices "github.com/dileepaj/tracified-gateway/services/ethereumServices/generalServices"
	polygongasservice "github.com/dileepaj/tracified-gateway/services/polygonServices/polygonGasService"
	gaspriceserviceforpolygon "github.com/dileepaj/tracified-gateway/services/polygonServices/polygonGasService/gasPriceServiceForPolygon"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/dileepaj/tracified-gateway/vendor/github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func PolygonContractDeployer(bin string, abi string, contractIdentifier string, contractType string, otherParams []any) (string, string, string, error) {
	contractAddress := ""
	transactionHash := ""
	transactionCost := ""
	var isFailed = true
	var predictedGasLimit int
	var predictedGasPrice = new(big.Int)
	var nonce uint64
	var errWhenGettingNonce error
	var deploymentError string
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
				var errWhenGettingGasPrice error
				predictedGasPrice, errWhenGettingGasPrice = gaspriceserviceforpolygon.MinimumGasPriceGetterForPolygon()
				if errWhenGettingGasPrice != nil {
					logger.LogWriter("Error when getting gas price : "+errWhenGettingGasLimit.Error(), constants.ERROR)
					return contractAddress, transactionHash, transactionCost, errors.New("Error when getting gas price : " + errWhenGettingGasPrice.Error())
				}
				if predictedGasPrice.Cmp(big.NewInt(0)) == 0 {
					logger.LogWriter("Error when getting gas price , gas price is zero", constants.ERROR)
					return contractAddress, transactionHash, transactionCost, errors.New("Error when getting gas price , gas price is zero")
				}

				nonce, errWhenGettingNonce = client.PendingNonceAt(context.Background(), common.Address(fromAddress))
				if errWhenGettingNonce != nil {
					logger.LogWriter("Error when getting nonce "+errWhenGettingNonce.Error(), constants.ERROR)
					return contractAddress, transactionHash, transactionCost, errors.New("Error when getting nonce , ERROR : " + errWhenGettingNonce.Error())
				}
			} else {
				if deploymentError == "nonce too low" {
					nonce, errWhenGettingNonce = client.PendingNonceAt(context.Background(), common.Address(fromAddress))
					if errWhenGettingNonce != nil {
						logger.LogWriter("Error when getting nonce "+errWhenGettingNonce.Error(), constants.ERROR)
						return contractAddress, transactionHash, transactionCost, errors.New("Error when getting nonce , ERROR : " + errWhenGettingNonce.Error())
					}
				} else if deploymentError == "intrinsic gas too low" {
					//increase gas limit by 10%
					predictedGasLimit = predictedGasLimit + int(predictedGasLimit*10/100)
				} else if deploymentError == "insufficient funds for gas * price + value" {
					//send email to increase the account balance
					errorInSendingEmail := contractdeployer.RequestFunds(2)
					if errorInSendingEmail != nil {
						logger.LogWriter("Error when sending email "+errorInSendingEmail.Error(), constants.ERROR)
						return contractAddress, transactionHash, transactionCost, errors.New("Error when sending email , ERROR : " + errorInSendingEmail.Error())
					}
					return contractAddress, transactionHash, transactionCost, errors.New("Gateway Polygon account funds are not enough")
				}
			}

			if predictedGasLimit > gasLimitCap || predictedGasPrice.Cmp(big.NewInt(int64(gasPriceCap))) == 1 {
				logger.LogWriter("Gas values are passing specified thresholds", constants.ERROR)
				return contractAddress, transactionHash, transactionCost, errors.New("Gas values are passing specified thresholds")
			}

			logger.LogWriter("Predicted gas limit : "+strconv.FormatInt(int64(predictedGasLimit), 10), constants.INFO)
			logger.LogWriter("Predicted gas price : "+predictedGasPrice.String(), constants.INFO)
			logger.LogWriter("Current nonce : "+strconv.FormatUint(nonce, 10), constants.INFO)

			auth.GasLimit = uint64(predictedGasLimit)
			auth.Nonce = big.NewInt(int64(nonce))
			auth.GasPrice = predictedGasPrice
		}
	}

	return contractAddress, transactionHash, transactionCost, nil
}
