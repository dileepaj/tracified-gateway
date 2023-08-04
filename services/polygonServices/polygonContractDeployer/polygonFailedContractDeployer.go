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
	"github.com/dileepaj/tracified-gateway/services/ethereumServices/dbCollectionHandler"
	generalservices "github.com/dileepaj/tracified-gateway/services/ethereumServices/generalServices"
	gaspriceserviceforpolygon "github.com/dileepaj/tracified-gateway/services/polygonServices/polygonGasService/gasPriceServiceForPolygon"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
)

func RedeployFailedContractsOnPolygon(failedContract model.PendingContracts) (string, string, string, *big.Int, *big.Int, int, error) {
	logger := utilities.NewCustomLogger()
	object := dao.Connection{}

	logger.LogWriter("------------Redeploying failed transaction-----------", constants.INFO)

	transactionHash := ""
	contractAddress := ""
	transactionCost := ""
	var predictedGasPrice = new(big.Int)
	var gasLimit int
	var nonce uint64
	var errWhenGettingNonce error
	var isFailed = true
	var deploymentError string

	//get the ABI and BIN
	abiString, binString, errWhenGettingABIandBIN := dbCollectionHandler.GetAbiAndBin(failedContract.ContractType, failedContract.Identifier)
	if errWhenGettingABIandBIN != nil {
		logger.LogWriter("Error when getting the ABI and BIN : "+errWhenGettingABIandBIN.Error(), constants.ERROR)
		return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when getting ABI and BIN : " + errWhenGettingABIandBIN.Error())
	}

	//load client and the keys
	client, privateKey, fromAddress, errWhenLoadingClientAndKey := generalservices.LoadClientAndKey(2)
	if errWhenLoadingClientAndKey != nil {
		logger.LogWriter("Error when loading the client and the key : "+errWhenLoadingClientAndKey.Error(), constants.ERROR)
		return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when loading the client and the key : " + errWhenLoadingClientAndKey.Error())

	}

	ContractBIN, parsed, errWhenLoadingParsedABIAndBIN := generalservices.LoadContractBinAndParsedAbi(binString, abiString)
	if errWhenLoadingParsedABIAndBIN != nil {
		logger.LogWriter("Error when loading ContractBIN and Parsed ABI : "+errWhenLoadingParsedABIAndBIN.Error(), constants.ERROR)
		return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when loading ContractBIN and Parsed ABI : " + errWhenLoadingParsedABIAndBIN.Error())
	}

	if parsed == nil {
		logger.LogWriter("GetABI returned nil", constants.ERROR)
		return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when getting ABI string , ERROR : GetAbi() returned nil")
	}

	//create the keyed transactor
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Value = big.NewInt(0) // in wei

	tryoutCap, errInTryConvert := strconv.Atoi(commons.GoDotEnvVariable("CONTRACTDEPLOYLIMIT"))
	if errInTryConvert != nil {
		logger.LogWriter("Error when converting the tryout limit , ERROR : "+errInTryConvert.Error(), constants.ERROR)
		return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when converting the tryout limit , ERROR : " + errInTryConvert.Error())
	}

	gasPriceCap, errInGasPriceCapConcert := strconv.Atoi(commons.GoDotEnvVariable("GASPRICECAP"))
	if errInGasPriceCapConcert != nil {
		logger.LogWriter("Error when converting the gas price cap , ERROR : "+errInGasPriceCapConcert.Error(), constants.ERROR)
		return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when converting the gas price cap , ERROR : " + errInGasPriceCapConcert.Error())
	}

	gasLimitCap, errInGasLimitCapConcert := strconv.Atoi(commons.GoDotEnvVariable("GASLIMITCAP"))
	if errInGasLimitCapConcert != nil {
		logger.LogWriter("Error when converting the gas limit cap , ERROR : "+errInGasLimitCapConcert.Error(), constants.ERROR)
		return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when converting the gas limit cap , ERROR : " + errInGasLimitCapConcert.Error())
	}

	gasLimit = failedContract.GasLimit

	//check the error to be corrected
	if failedContract.ErrorMessage == "out of gas" || failedContract.ErrorMessage == "contract creation code storage out of gas" {
		gasLimit = gasLimit + int(gasLimit*10/100)
	}

	for i := 0; i < tryoutCap; i++ {
		if !isFailed {
			return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, nil
		} else {
			logger.LogWriter("Deploying the contract for the "+strconv.Itoa(i+1)+" th time", constants.INFO)
			if i == 0 {
				//get the initially corrected values
				auth.GasLimit = uint64(gasLimit)
				nonce, errWhenGettingNonce = client.PendingNonceAt(context.Background(), common.Address(fromAddress))
				if errWhenGettingNonce != nil {
					logger.LogWriter("Error when getting nonce "+errWhenGettingNonce.Error(), constants.ERROR)
					return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when getting nonce , ERROR : " + errWhenGettingNonce.Error())
				}
				var errWhenGettingGasPrice error
				predictedGasPrice, errWhenGettingGasPrice = gaspriceserviceforpolygon.MinimumGasPriceGetterForPolygon()
				if errWhenGettingGasPrice != nil {
					logger.LogWriter("Error when getting gas price "+errWhenGettingGasPrice.Error(), constants.ERROR)
					return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when getting gas price, ERROR : " + errWhenGettingGasPrice.Error())
				}
				if predictedGasPrice.Cmp(big.NewInt(0)) == 0 {
					logger.LogWriter("Error when getting gas price , gas price is zero", constants.ERROR)
					return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when getting gas price , gas price is zero")
				}
			} else {
				//check the error
				if deploymentError == "nonce too low" {
					//pick up the latest the nonce available
					nonce, errWhenGettingNonce = client.PendingNonceAt(context.Background(), common.Address(fromAddress))
					if errWhenGettingNonce != nil {
						logger.LogWriter("Error when getting nonce "+errWhenGettingNonce.Error(), constants.ERROR)
						return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when getting nonce , ERROR : " + errWhenGettingNonce.Error())
					}
				} else if deploymentError == "intrinsic gas too low" {
					//increase gas limit by 10%
					gasLimit = gasLimit + int(gasLimit*10/100)
				} else if deploymentError == "insufficient funds for gas * price + value" {
					//send email to increase the account balance
					errorInSendingEmail := contractdeployer.RequestFunds(2)
					if errorInSendingEmail != nil {
						logger.LogWriter("Error when sending email "+errorInSendingEmail.Error(), constants.ERROR)
						return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when sending email , ERROR : " + errorInSendingEmail.Error())
					}
					return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Gateway Ethereum account funds are not enough")
				}
			}
			if gasLimit > gasLimitCap {
				logger.LogWriter("Gas limit is passing the threshold", constants.ERROR)
				return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Gas limit is passing the threshold")
			}
			if predictedGasPrice.Cmp(big.NewInt(int64(gasPriceCap))) == 1 {
				logger.LogWriter("Gas price is passing the threshold", constants.ERROR)
				return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Gas price is passing the threshold")
			}

			logger.LogWriter("Predicted gas limit : "+strconv.Itoa(gasLimit), constants.INFO)
			logger.LogWriter("Predicted gas price : "+predictedGasPrice.String(), constants.INFO)
			logger.LogWriter("Current nonce : "+strconv.FormatUint(nonce, 10), constants.INFO)

			auth.GasLimit = uint64(gasLimit) // in units
			auth.Nonce = big.NewInt(int64(nonce))
			auth.GasPrice = predictedGasPrice

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
					Status:          117,
					CurrentIndex:    failedContract.CurrentIndex + 1,
					ErrorMessage:    "",
					ContractType:    failedContract.ContractType,
					Identifier:      failedContract.Identifier,
					Nonce:           auth.Nonce,
					GasLimit:        int(auth.GasLimit),
					GasPrice:        auth.GasPrice,
				}

				errInInsertingPendingTx := object.InsertPolygonPendingContract(pendingTransaction)
				if errInInsertingPendingTx != nil {
					logrus.Error("Error in inserting the pending transaction, ERROR : " + errInInsertingPendingTx.Error())
					isFailed = true
				} else {
					isFailed = false
				}

				// calculate the predicted transaction cost
				costInWei := new(big.Int).Mul(big.NewInt(int64(gasLimit)), predictedGasPrice)
				cost := new(big.Float).Quo(new(big.Float).SetInt(costInWei), big.NewFloat(math.Pow10(18)))
				transactionCost = fmt.Sprintf("%g", cost) + " MATIC"
			}
		}
	}
	if !isFailed {
		return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, nil
	}

	return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Threshold for contract redeployment exceeded")
}
