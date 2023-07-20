package polygoncontractdeployer

import (
	"context"
	"errors"
	"fmt"
	"math"
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

	//TODO-Create metric bind object

	//load client and key
	client, privateKey, fromAddress, errWhenLoadingClientAndKey := generalservices.LoadClientAndKey(2)
	if errWhenLoadingClientAndKey != nil {
		logger.LogWriter("Error when loading the client and the key : "+errWhenLoadingClientAndKey.Error(), constants.ERROR)
		return contractAddress, transactionHash, transactionCost, errors.New("Error when loading the client and the key : " + errWhenLoadingClientAndKey.Error())
	}
	ContractBIN, parsed, errWhenLoadingParsedABIAndBIN := generalservices.LoadContractBinAndParsedAbi(bin, abi)
	if errWhenLoadingParsedABIAndBIN != nil {
		logger.LogWriter("Error when loading ContractBIN and Parsed ABI : "+errWhenLoadingParsedABIAndBIN.Error(), constants.ERROR)
		return contractAddress, transactionHash, transactionCost, errors.New("Error when loading ContractBIN and Parsed ABI : " + errWhenLoadingParsedABIAndBIN.Error())
	}
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

			//call the bind deployer method
			address, tx, contract, errWhenDeployingContract := bind.DeployContract(auth, *parsed, common.FromHex(ContractBIN), client)
			if errWhenDeployingContract != nil {
				logger.LogWriter("Error when deploying contract "+errWhenDeployingContract.Error(), constants.ERROR)
				isFailed = true
				deploymentError = errWhenDeployingContract.Error()
				// inserting error message to the database
				errorMessage := model.EthErrorMessage{
					TransactionHash: "",
					ErrorMessage:    deploymentError,
					Network:         "polygon",
				}
				errInInsertingErrorMessage := object.InsertPolygonErrorMessage(errorMessage)
				if errInInsertingErrorMessage != nil {
					logger.LogWriter("Error in inserting the error message, ERROR : "+errInInsertingErrorMessage.Error(), constants.ERROR)
				}

				//TODO-handle metric bind request
			} else {
				contractAddress = address.Hex()
				transactionHash = tx.Hash().Hex()
				_ = contract

				env := commons.GoDotEnvVariable("ENVIRONMENT")
				if env == "QA" {
					logger.LogWriter("View contract at : https://mumbai.polygonscan.com/address/"+address.Hex(), constants.INFO)
					logger.LogWriter("View transaction at : https://mumbai.polygonscan.com/tx/"+tx.Hash().Hex(), constants.INFO)
				} else if env == "PRODUCTION" {
					logger.LogWriter("View contract at : https://polygonscan.com/tx/"+address.Hex(), constants.INFO)
					logger.LogWriter("View transaction at : https://polygonscan.com/address/"+tx.Hash().Hex(), constants.INFO)
				} else if env == "STAGING" {
					logger.LogWriter("View contract at : https://mumbai.polygonscan.com/address/"+address.Hex(), constants.INFO)
					logger.LogWriter("View transaction at : https://mumbai.polygonscan.com/tx/"+tx.Hash().Hex(), constants.INFO)
				}

				// Insert the pending transaction to the database
				pendingTransaction := model.PendingContracts{
					TransactionHash: tx.Hash().Hex(),
					ContractAddress: address.Hex(),
					Status:          117, //PENDING
					CurrentIndex:    0,
					ErrorMessage:    "",
					ContractType:    contractType,
					Identifier:      contractIdentifier,
					Nonce:           auth.Nonce,
					GasLimit:        int(auth.GasLimit),
					GasPrice:        auth.GasPrice,
				}
				errInInsertingPendingTx := object.InsertPolygonPendingContract(pendingTransaction)
				if errInInsertingPendingTx != nil {
					logger.LogWriter("Error in inserting the pending transaction, ERROR : "+errInInsertingPendingTx.Error(), constants.ERROR)
					isFailed = true
				} else {
					isFailed = false
				}

				// calculate the predicted transaction cost
				costInWei := new(big.Int).Mul(big.NewInt(int64(predictedGasLimit)), predictedGasPrice)
				cost := new(big.Float).Quo(new(big.Float).SetInt(costInWei), big.NewFloat(math.Pow10(18)))
				transactionCost = fmt.Sprintf("%g", cost) + " MATIC"

				//TODO - handle insert and update for metric binding
			}
		}
	}
	if !isFailed {
		return contractAddress, transactionHash, transactionCost, nil

	}

	return contractAddress, transactionHash, transactionCost, errors.New("Threshold for contract redeployment exceeded")
}
