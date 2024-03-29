package contractdeployer

import (
	"context"

	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/services/ethereumServices/dbCollectionHandler"
	gasServices "github.com/dileepaj/tracified-gateway/services/ethereumServices/gasServices"
	"github.com/dileepaj/tracified-gateway/services/ethereumServices/gasServices/gasPriceServices"
	generalservices "github.com/dileepaj/tracified-gateway/services/ethereumServices/generalServices"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
)

/*
Deploy smart contracts on to Ethereum with failure replacements
*/
func EthereumContractDeployerService(bin string, abi string, contractIdentifier string, contractType string, otherParams []any) (string, string, string, error) {
	contractAddress := ""
	transactionHash := ""
	transactionCost := ""
	var isFailed = true
	var predictedGasLimit int
	var predictedGasPrice = new(big.Int)
	var deploymentError string
	var nonce uint64
	var errWhenGettingNonce error

	object := dao.Connection{}


	var expertFormulaObj model.EthereumExpertFormula
	var metricBindObj model.EthereumMetricBind

	logrus.Info("Calling the deployer service.............")

	if contractType == "ETHEXPERTFORMULA" {
		expertFormulaObj = otherParams[0].(model.EthereumExpertFormula)
		expertFormulaObj.ActualStatus = 110 // DEPLOYMENT_STARTED
		errorWhenUpdatingStatus := object.UpdateSelectedEthFormulaFields(expertFormulaObj.FormulaID, expertFormulaObj.TransactionUUID, expertFormulaObj)
		if errorWhenUpdatingStatus != nil {
			logrus.Error("Error when updating the status after deployment started : " + errorWhenUpdatingStatus.Error())
			return contractAddress, transactionHash, transactionCost, errors.New("Error when updating the status after deployment started : " + errorWhenUpdatingStatus.Error())
		}
	} else if contractType == "ETHMETRICBIND" {
		metricBindObj = otherParams[0].(model.EthereumMetricBind)
		metricBindObj.ActualStatus = 110 // DEPLOYMENT_STARTED
		errorWhenUpdatingStatus := object.UpdateSelectedEthMetricFields(metricBindObj.MetricID, metricBindObj.TransactionUUID, metricBindObj)
		if errorWhenUpdatingStatus != nil {
			logrus.Error("Error when updating the status after deployment started : " + errorWhenUpdatingStatus.Error())
			return contractAddress, transactionHash, transactionCost, errors.New("Error when updating the status after deployment started : " + errorWhenUpdatingStatus.Error())
		}
	}

	//load client and the keys
	client, privateKey, fromAddress, errWhenLoadingClientAndKey := generalservices.LoadClientAndKey()
	if errWhenLoadingClientAndKey != nil {
		logrus.Error("Error when loading the client and the key : " + errWhenLoadingClientAndKey.Error())
		return contractAddress, transactionHash, transactionCost, errors.New("Error when loading the client and the key : " + errWhenLoadingClientAndKey.Error())
	}

	ContractBIN, parsed, errWhenLoadingParsedABIAndBIN := generalservices.LoadContractBinAndParsedAbi(bin, abi)
	if errWhenLoadingParsedABIAndBIN != nil {
		logrus.Error("Error when loading ContractBIN and Parsed ABI : " + errWhenLoadingParsedABIAndBIN.Error())
		return contractAddress, transactionHash, transactionCost, errors.New("Error when loading ContractBIN and Parsed ABI : " + errWhenLoadingParsedABIAndBIN.Error())
	}

	if parsed == nil {
		logrus.Error("GetABI returned nil")
		return contractAddress, transactionHash, transactionCost, errors.New("Error when getting ABI string , ERROR : GetAbi() returned nil")
	}

	//create the keyed transactor
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Value = big.NewInt(0) // in wei

	tryoutCap, errInTryConvert := strconv.Atoi(commons.GoDotEnvVariable("CONTRACTDEPLOYLIMIT"))
	if errInTryConvert != nil {
		logrus.Error("Error when converting the tryout limit , ERROR : " + errInTryConvert.Error())
		return contractAddress, transactionHash, transactionCost, errors.New("Error when converting the tryout limit , ERROR : " + errInTryConvert.Error())
	}

	gasLimitCap, errInGasLimitCapConcert := strconv.Atoi(commons.GoDotEnvVariable("GASLIMITCAP"))
	if errInGasLimitCapConcert != nil {
		logrus.Error("Error when converting the gas limit cap , ERROR : " + errInGasLimitCapConcert.Error())
		return contractAddress, transactionHash, transactionCost, errors.New("Error when converting the gas limit cap , ERROR : " + errInGasLimitCapConcert.Error())
	}

	gasPriceCap, errInGasPriceCapConcert := strconv.Atoi(commons.GoDotEnvVariable("GASPRICECAP"))
	if errInGasPriceCapConcert != nil {
		logrus.Error("Error when converting the gas price cap , ERROR : " + errInGasPriceCapConcert.Error())
		return contractAddress, transactionHash, transactionCost, errors.New("Error when converting the gas price cap , ERROR : " + errInGasPriceCapConcert.Error())
	}

	for i := 0; i < tryoutCap; i++ {
		if !isFailed {
			return contractAddress, transactionHash, transactionCost, nil
		} else {
			logrus.Info("Deploying the contract for the ", i+1, " th time")
			//if the first iteration take the initial gas limit and gas price
			if i == 0 {
				//get the initial gas limit
				gasLimit, errInGettingGasLimit := gasServices.EstimateGasLimit(commons.GoDotEnvVariable("ETHEREUMPUBKEY"), "", "", "", "", "", "", bin)
				if errInGettingGasLimit != nil {
					logrus.Error("Error when getting gas limit " + errInGettingGasLimit.Error())
					return contractAddress, transactionHash, transactionCost, errors.New("Error when getting gas limit, ERROR : " + errInGettingGasLimit.Error())
				}
				predictedGasLimit = int(gasLimit)
				//get the initial gas price
				var errWhenGettingGasPrice error
				predictedGasPrice, errWhenGettingGasPrice = gasPriceServices.GetMinGasPrice()
				if errWhenGettingGasPrice != nil {
					logrus.Error("Error when getting gas price " + errWhenGettingGasPrice.Error())
					return contractAddress, transactionHash, transactionCost, errors.New("Error when getting gas price, ERROR : " + errWhenGettingGasPrice.Error())
				}
				if predictedGasPrice.Cmp(big.NewInt(0)) == 0 {
					logrus.Error("Error when getting gas price , gas price is zero")
					return contractAddress, transactionHash, transactionCost, errors.New("Error when getting gas price , gas price is zero")
				}

				auth.GasLimit = uint64(predictedGasLimit) // in units
				nonce, errWhenGettingNonce = client.PendingNonceAt(context.Background(), common.Address(fromAddress))
				if errWhenGettingNonce != nil {
					logrus.Error("Error when getting nonce " + errWhenGettingNonce.Error())
					return contractAddress, transactionHash, transactionCost, errors.New("Error when getting nonce , ERROR : " + errWhenGettingNonce.Error())
				}
			} else {
				//check the error
				if deploymentError == "nonce too low" {
					//pick up the latest the nonce available
					nonce, errWhenGettingNonce = client.PendingNonceAt(context.Background(), common.Address(fromAddress))
					if errWhenGettingNonce != nil {
						logrus.Error("Error when getting nonce " + errWhenGettingNonce.Error())
						return contractAddress, transactionHash, transactionCost, errors.New("Error when getting nonce , ERROR : " + errWhenGettingNonce.Error())
					}

				} else if deploymentError == "intrinsic gas too low" {
					//increase gas limit by 10%
					predictedGasLimit = predictedGasLimit + int(predictedGasLimit*10/100)
				} else if deploymentError == "insufficient funds for gas * price + value" {
					//send email to increase the account balance
					errorInSendingEmail := RequestFunds()
					if errorInSendingEmail != nil {
						logrus.Error("Error when sending email " + errorInSendingEmail.Error())
						return contractAddress, transactionHash, transactionCost, errors.New("Error when sending email , ERROR : " + errorInSendingEmail.Error())
					}
					return contractAddress, transactionHash, transactionCost, errors.New("Gateway Ethereum account funds are not enough")

				}
			}

			//check the gas limit cap and gas price cap
			if predictedGasLimit > gasLimitCap || predictedGasPrice.Cmp(big.NewInt(int64(gasPriceCap))) == 1 {
				logrus.Error("Gas values are passing specified thresholds")
				return contractAddress, transactionHash, transactionCost, errors.New("Gas values are passing specified thresholds")
			}

			logrus.Info("Predicted gas limit : ", predictedGasLimit)
			logrus.Info("Predicted gas price : ", predictedGasPrice)
			logrus.Info("Current nonce : ", nonce)

			auth.GasLimit = uint64(predictedGasLimit) // in units
			auth.Nonce = big.NewInt(int64(nonce))
			auth.GasPrice = predictedGasPrice

			//call the deployer method
			address, tx, contract, errWhenDeployingContract := bind.DeployContract(auth, *parsed, common.FromHex(ContractBIN), client)
			if errWhenDeployingContract != nil {
				logrus.Info("Error when deploying contract " + errWhenDeployingContract.Error())
				isFailed = true
				deploymentError = errWhenDeployingContract.Error()
				// inserting error message to the database
				errorMessage := model.EthErrorMessage{
					TransactionHash: "",
					ErrorMessage:    deploymentError,
					Network:         "sepolia",
				}
				errInInsertingErrorMessage := object.InsertEthErrorMessage(errorMessage)
				if errInInsertingErrorMessage != nil {
					logrus.Error("Error in inserting the error message, ERROR : " + errInInsertingErrorMessage.Error())
				}

				if contractType == "ETHMETRICBIND" {
					pendingContractObj := model.PendingContracts{
						TransactionHash: "",
						ContractAddress: "",
						Status:          119,	//FAILED
						CurrentIndex:    0,
						ErrorMessage:    deploymentError,
						ContractType:    contractType,
						Identifier:      contractIdentifier,
					}
					errorWhenInvalidating := dbCollectionHandler.InvalidateMetric(pendingContractObj, 119, deploymentError)
					if errorWhenInvalidating != nil {
						logrus.Error("Error when invalidating the metric, ERROR : " + errorWhenInvalidating.Error())
					}
				}
			} else {
				contractAddress = address.Hex()
				transactionHash = tx.Hash().Hex()
				_ = contract

				logrus.Info("View contract at : https://sepolia.etherscan.io/address/", address.Hex())
				logrus.Info("View transaction at : https://sepolia.etherscan.io/tx/", tx.Hash().Hex())

				// Insert the pending transaction to the database
				pendingTransaction := model.PendingContracts{
					TransactionHash: tx.Hash().Hex(),
					ContractAddress: address.Hex(),
					Status:          117,	//PENDING
					CurrentIndex:    0,
					ErrorMessage:    "",
					ContractType:    contractType,
					Identifier:      contractIdentifier,
					Nonce:           auth.Nonce,
					GasLimit:        int(auth.GasLimit),
					GasPrice:        auth.GasPrice,
				}
				errInInsertingPendingTx := object.InsertEthPendingContract(pendingTransaction)
				if errInInsertingPendingTx != nil {
					logrus.Error("Error in inserting the pending transaction, ERROR : " + errInInsertingPendingTx.Error())
					isFailed = true
				} else {
					isFailed = false
				}

				// calculate the predicted transaction cost
				costInWei := new(big.Int).Mul(big.NewInt(int64(predictedGasLimit)), predictedGasPrice)
				cost := new(big.Float).Quo(new(big.Float).SetInt(costInWei), big.NewFloat(math.Pow10(18)))
				transactionCost = fmt.Sprintf("%g", cost) + " ETH"

				// insert and update latest metric contract address
				if contractType == "ETHMETRICBIND" {
					// get the metric object from the database using uuid
					ethMetricObj, errWhenGettingMetric := dbCollectionHandler.GetEthMetricByUUID(contractIdentifier)
					if errWhenGettingMetric != nil {
						logrus.Error("Error when getting metric object from DB: ", errWhenGettingMetric)
						return contractAddress, transactionHash, transactionCost, errors.New("Error when getting metric object from DB for the latest contract update, ERROR : " + errWhenGettingMetric.Error())
					}

					insertObj := model.MetricLatestContract{
						MetricID:        ethMetricObj.MetricID,
						ContractAddress: address.Hex(),
						Type:            ethMetricObj.Type,
					}
					if ethMetricObj.Type == "METADATA" {
						//insert the latest contract address in DB
						errWhenInsertingToLatest := object.EthereumInsertToMetricLatestContract(insertObj)
						if errWhenInsertingToLatest != nil {
							logrus.Error("Error when inserting to latest contract to DB: ", errWhenInsertingToLatest)
						}
						logrus.Info("Added " + address.Hex() + " to latest contract collection")
					} else if ethMetricObj.Type == "ACTIVITY" {
						//update the latest contract address in DB
						errWhenUpdatingLatest := object.UpdateEthereumMetricLatestContract(ethMetricObj.MetricID, insertObj)
						if errWhenUpdatingLatest != nil {
							logrus.Errorf("Error when updating latest contract address in DB: ", errWhenUpdatingLatest)
						}
						logrus.Info("Updated " + address.Hex() + " as latest contract")
					}
				}
			}
		}
	}
	if !isFailed {
		return contractAddress, transactionHash, transactionCost, nil
	}

	return contractAddress, transactionHash, transactionCost, errors.New("Threshold for contract redeployment exceeded")
}
