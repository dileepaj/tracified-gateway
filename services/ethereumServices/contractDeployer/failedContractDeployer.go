package contractdeployer

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/services/ethereumServices/dbCollectionHandler"
	"github.com/dileepaj/tracified-gateway/services/ethereumServices/gasServices/gasPriceServices"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

func RedeployFailedContracts(failedContract model.PendingContracts) (string, string, string, *big.Int, *big.Int, int, error) {
	logrus.Info("----Redeploying failed transaction-------------")

	transactionHash := ""
	contractAddress := ""
	transactionCost := ""
	var nonce uint64
	var deploymentError string
	var errWhenGettingNonce error
	var isFailed = true
	var predictedGasPrice = new(big.Int)
	object := dao.Connection{}
	var gasLimit int

	//get the ABI and BIN
	abiString, binString, errWhenGettingABIandBIN := dbCollectionHandler.GetAbiAndBin(failedContract.ContractType, failedContract.Identifier)
	if errWhenGettingABIandBIN != nil {
		logrus.Error("Error when getting the ABI and BIN : " + errWhenGettingABIandBIN.Error())
		return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when getting ABI and BIN : " + errWhenGettingABIandBIN.Error())
	}

	//Dial infura client
	client, errWhenDialingEthClient := ethclient.Dial(commons.GoDotEnvVariable("ETHEREUMTESTNETLINK"))
	if errWhenDialingEthClient != nil {
		logrus.Error("Error when dialing the eth client : " + errWhenDialingEthClient.Error())
		return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when dialing eth client , ERROR : " + errWhenDialingEthClient.Error())
	}

	//load ECDSA private key
	privateKey, errWhenGettingECDSAKey := crypto.HexToECDSA(commons.GoDotEnvVariable("ETHEREUMSECKEY"))
	if errWhenGettingECDSAKey != nil {
		logrus.Error("Error when getting ECDSA key " + errWhenGettingECDSAKey.Error())
		return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when getting ECDSA key , ERROR : " + errWhenGettingECDSAKey.Error())
	}

	//get the public key
	publicKey := privateKey.Public()
	//get public key ECDSA
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		logrus.Error("Cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	//assign metadata for the contract
	var BuildData = &bind.MetaData{
		ABI: abiString,
		Bin: binString,
	}

	//var ContractABI = BuildData.ABI
	var ContractBIN = BuildData.Bin

	parsed, errWhenGettingABI := BuildData.GetAbi()
	if errWhenGettingABI != nil {
		logrus.Error("Error when getting abi from passed ABI string " + errWhenGettingABI.Error())
		return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when getting abi from passed ABI string , ERROR : " + errWhenGettingABI.Error())
	}

	if parsed == nil {
		logrus.Error("GetABI returned nil")
		return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when getting ABI string , ERROR : GetAbi() returned nil")
	}

	//create the keyed transactor
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Value = big.NewInt(0) // in wei

	tryoutCap, errInTryConvert := strconv.Atoi(commons.GoDotEnvVariable("CONTRACTDEPLOYLIMIT"))
	if errInTryConvert != nil {
		logrus.Error("Error when converting the tryout limit , ERROR : " + errInTryConvert.Error())
		return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when converting the tryout limit , ERROR : " + errInTryConvert.Error())
	}

	gasPriceCap, errInGasPriceCapConcert := strconv.Atoi(commons.GoDotEnvVariable("GASPRICECAP"))
	if errInGasPriceCapConcert != nil {
		logrus.Error("Error when converting the gas price cap , ERROR : " + errInGasPriceCapConcert.Error())
		return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when converting the gas price cap , ERROR : " + errInGasPriceCapConcert.Error())
	}

	gasLimitCap, errInGasLimitCapConcert := strconv.Atoi(commons.GoDotEnvVariable("GASLIMITCAP"))
	if errInGasLimitCapConcert != nil {
		logrus.Error("Error when converting the gas limit cap , ERROR : " + errInGasLimitCapConcert.Error())
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
			logrus.Info("Deploying the contract for the ", i+1, " th time")
			if i == 0 {
				//get the initially corrected values
				auth.GasLimit = uint64(gasLimit)
				nonce, errWhenGettingNonce = client.PendingNonceAt(context.Background(), fromAddress)
				if errWhenGettingNonce != nil {
					logrus.Error("Error when getting nonce " + errWhenGettingNonce.Error())
					return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when getting nonce , ERROR : " + errWhenGettingNonce.Error())
				}
				var errWhenGettingGasPrice error
				predictedGasPrice, errWhenGettingGasPrice = gasPriceServices.GetMinGasPrice()
				if errWhenGettingGasPrice != nil {
					logrus.Error("Error when getting gas price " + errWhenGettingGasPrice.Error())
					return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when getting gas price, ERROR : " + errWhenGettingGasPrice.Error())
				}
				if predictedGasPrice.Cmp(big.NewInt(0)) == 0 {
					logrus.Error("Error when getting gas price , gas price is zero")
					return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when getting gas price , gas price is zero")
				}
			} else {
				//check the error
				if deploymentError == "nonce too low" {
					//pick up the latest the nonce available
					nonce, errWhenGettingNonce = client.PendingNonceAt(context.Background(), fromAddress)
					if errWhenGettingNonce != nil {
						logrus.Error("Error when getting nonce " + errWhenGettingNonce.Error())
						return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when getting nonce , ERROR : " + errWhenGettingNonce.Error())
					}

				} else if deploymentError == "intrinsic gas too low" {
					//increase gas limit by 10%
					gasLimit = gasLimit + int(gasLimit*10/100)
				} else if deploymentError == "insufficient funds for gas * price + value" {
					//send email to increase the account balance
					errorInSendingEmail := RequestFunds()
					if errorInSendingEmail != nil {
						logrus.Error("Error when sending email " + errorInSendingEmail.Error())
						return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Error when sending email , ERROR : " + errorInSendingEmail.Error())
					}
					return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Gateway Ethereum account funds are not enough")

				}
			}
			//check the gas limit cap and gas price cap
			if gasLimit > gasLimitCap || predictedGasPrice.Cmp(big.NewInt(int64(gasPriceCap))) == 1 {
				logrus.Error("Gas values are passing specified thresholds")
				return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, errors.New("Gas values are passing specified thresholds")
			}

			logrus.Info("Predicted gas limit : ", gasLimit)
			logrus.Info("Predicted gas price : ", predictedGasPrice)
			logrus.Info("Current nonce : ", nonce)

			auth.GasLimit = uint64(gasLimit) // in units
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
					Status:          "PENDING",
					CurrentIndex:    0,
					ErrorMessage:    "",
					ContractType:    failedContract.ContractType,
					Identifier:      failedContract.Identifier,
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
				costInWei := new(big.Int).Mul(big.NewInt(int64(gasLimit)), predictedGasPrice)
				cost := new(big.Float).Quo(new(big.Float).SetInt(costInWei), big.NewFloat(math.Pow10(18)))
				transactionCost = fmt.Sprintf("%g", cost) + " ETH"
			}
		}
	}
	if !isFailed {
		return contractAddress, transactionHash, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, nil
	}

	return transactionHash, contractAddress, transactionCost, big.NewInt(int64(nonce)), predictedGasPrice, gasLimit, nil
}
