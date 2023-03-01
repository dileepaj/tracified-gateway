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
	"github.com/dileepaj/tracified-gateway/protocols/ethereum/deploy"
	gasServices "github.com/dileepaj/tracified-gateway/services/ethereumServices/gasServices"
	"github.com/dileepaj/tracified-gateway/services/ethereumServices/gasServices/gasPriceServices"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

/*
Deploy smart contracts on to Ethereum with failure replacements
*/
func EthereumContractDeployerService(bin string, abi string) (string, string, string, error) {
	contractAddress := ""
	transactionHash := ""
	transactionCost := ""

	object := dao.Connection{}

	logrus.Info("Calling the deployer service.............")

	//Dial infura client
	client, errWhenDialingEthClient := ethclient.Dial(commons.GoDotEnvVariable("ETHEREUMTESTNETLINK"))
	if errWhenDialingEthClient != nil {
		logrus.Error("Error when dialing the eth client " + errWhenDialingEthClient.Error())
		return contractAddress, transactionHash, transactionCost, errors.New("Error when dialing eth client , ERROR : " + errWhenDialingEthClient.Error())
	}

	//load ECDSA private key
	privateKey, errWhenGettingECDSAKey := crypto.HexToECDSA(commons.GoDotEnvVariable("ETHEREUMSECKEY"))
	if errWhenGettingECDSAKey != nil {
		logrus.Error("Error when getting ECDSA key " + errWhenGettingECDSAKey.Error())
		return contractAddress, transactionHash, transactionCost, errors.New("Error when getting ECDSA key , ERROR : " + errWhenGettingECDSAKey.Error())
	}

	//get the public key
	publicKey := privateKey.Public()
	//get public key ECDSA
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		logrus.Error("Cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return contractAddress, transactionHash, transactionCost, errors.New("Cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	//assign metadata for the contract
	var BuildData = &bind.MetaData{
		ABI: abi,
		Bin: bin,
	}

	//var ContractABI = BuildData.ABI
	var ContractBIN = BuildData.Bin

	parsed, errWhenGettingABI := BuildData.GetAbi()
	if errWhenGettingABI != nil {
		logrus.Error("Error when getting abi from passed ABI string " + errWhenGettingABI.Error())
		return contractAddress, transactionHash, transactionCost, errors.New("Error when getting abi from passed ABI string , ERROR : " + errWhenGettingABI.Error())
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

	var isFailed = true
	var predictedGasLimit int
	var predictedGasPrice = new(big.Int)
	var deploymentError string
	var nonce uint64
	var errWhenGettingNonce error
	var initialNonce uint64
	var currentHash common.Hash
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
				nonce, errWhenGettingNonce = client.PendingNonceAt(context.Background(), fromAddress)
				if errWhenGettingNonce != nil {
					logrus.Error("Error when getting nonce " + errWhenGettingNonce.Error())
					return contractAddress, transactionHash, transactionCost, errors.New("Error when getting nonce , ERROR : " + errWhenGettingNonce.Error())
				}
				initialNonce = nonce
			} else {
				//check the error
				if deploymentError == "nonce too low" {
					//pick up the latest the nonce available
					nonce, errWhenGettingNonce = client.PendingNonceAt(context.Background(), fromAddress)
					if errWhenGettingNonce != nil {
						logrus.Error("Error when getting nonce " + errWhenGettingNonce.Error())
						return contractAddress, transactionHash, transactionCost, errors.New("Error when getting nonce , ERROR : " + errWhenGettingNonce.Error())
					}

					//check the previous nonce with the current nonce if the 2nd loop is in action to see whether the transaction is already on the blockchain
					if i == 1 {
						if initialNonce != nonce {
							//! if this is true that means this transaction is already happened in the blockchain, next check the success of failure status
							transactionReceipt, errWhenTakingTransactionForSameNonce := client.TransactionReceipt(context.Background(), currentHash)
							if errWhenTakingTransactionForSameNonce != nil {
								logrus.Info("Error when taking the transaction receipt for the previously pending transaction, Error : " + errWhenTakingTransactionForSameNonce.Error())
								return contractAddress, transactionHash, transactionCost, errors.New("Error when taking the transaction receipt for the previously pending transaction, Error : " + errWhenTakingTransactionForSameNonce.Error())
							}

							//checking the transaction receipt status
							if transactionReceipt.Status == 0 {
								//transaction has been failed
								//try in the next iteration
								isFailed = true

							} else if transactionReceipt.Status == 1 {
								//transaction is successful -> return with success message
								return contractAddress, transactionHash, transactionCost, nil

							} else {
								logrus.Error("Invalid transaction receipt status")
								return contractAddress, transactionHash, transactionCost, errors.New("Invalid transaction receipt status")
							}
						}
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

				} else if deploymentError == "replacement transaction underpriced" {
					//increase gas price by 10%
					predictedGasPrice = new(big.Int).Add(predictedGasPrice, new(big.Int).Div(predictedGasPrice, big.NewInt(10)))
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
			} else {
				contractAddress = address.Hex()
				transactionHash = tx.Hash().Hex()
				_ = contract

				currentHash = tx.Hash()

				logrus.Info("View contract at : https://sepolia.etherscan.io/address/", address.Hex())
				logrus.Info("View transaction at : https://sepolia.etherscan.io/tx/", tx.Hash().Hex())

				receipt, errInGettingReceipt := bind.WaitMined(context.Background(), client, tx)
				if errInGettingReceipt != nil {
					logrus.Error("Error in getting receipt: Error: " + errInGettingReceipt.Error())
					return contractAddress, transactionHash, transactionCost, errors.New("Error in getting receipt: Error: " + errInGettingReceipt.Error())
				} else {
					costInWei := new(big.Int).Mul(big.NewInt(int64(receipt.GasUsed)), predictedGasPrice)
					cost := new(big.Float).Quo(new(big.Float).SetInt(costInWei), big.NewFloat(math.Pow10(18)))
					transactionCost = fmt.Sprintf("%g", cost) + " ETH"

					if receipt.Status == 0 {
						isFailed = true
						errorMessageFromStatus, errorInCallingTransactionStatus := deploy.GetErrorOfFailedTransaction(tx.Hash().Hex())
						if errorInCallingTransactionStatus != nil {
							logrus.Error("Transaction failed.")
							logrus.Error("Error when getting the error for the transaction failure: Error: " + errorInCallingTransactionStatus.Error())
							return contractAddress, transactionHash, transactionCost, errors.New("Transaction failed.")
						} else {
							logrus.Error("Transaction failed. Error: " + errorMessageFromStatus)
							// inserting error message to the database
							errorMessage := model.EthErrorMessage{
								TransactionHash: tx.Hash().Hex(),
								ErrorMessage:    errorMessageFromStatus,
								Network:         "sepolia",
							}
							errInInsertingErrorMessage := object.InsertEthErrorMessage(errorMessage)
							if errInInsertingErrorMessage != nil {
								logrus.Error("Error in inserting the error message, ERROR : " + errInInsertingErrorMessage.Error())
							}
						}
					} else if receipt.Status == 1 {
						isFailed = false
					} else {
						logrus.Error("Invalid receipt status for 'WaitMined', Status : ", receipt.Status)
						return contractAddress, transactionHash, transactionCost, errors.New("Invalid receipt status for 'WaitMined', Status : " + fmt.Sprint(receipt.Status))
					}
					logrus.Info("Status of receipt : ", receipt.Status)
					logrus.Info(isFailed)
				}

			}

		}
	}
	if !isFailed {
		return contractAddress, transactionHash, transactionCost, nil
	}

	return contractAddress, transactionHash, transactionCost, errors.New("Threshold for contract redeployment exceeded")
}
