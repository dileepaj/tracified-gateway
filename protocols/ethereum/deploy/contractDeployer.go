package deploy

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

/*
	Deploy smart contract in Ethereum
*/
func DeployContract(abi string, bin string) (string, string, string, error) {
	contractAddress := ""
	transactionHash := ""
	transactionCost := ""

	//Dial infura client
	client, errWhenDialingEthClinet := ethclient.Dial(commons.GoDotEnvVariable("ETHEREUMTESTNETLINK"))
	if errWhenDialingEthClinet != nil {
		logrus.Error("Error when dialing the eth client " + errWhenDialingEthClinet.Error())
		return contractAddress, transactionHash, transactionCost, errors.New("Error when dialing eth client , ERROR : " + errWhenDialingEthClinet.Error())
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

	//get nonce
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, errWhenGettingNonce := client.PendingNonceAt(context.Background(), fromAddress)
	if errWhenGettingNonce != nil {
		logrus.Error("Error when getting nonce " + errWhenGettingNonce.Error())
		return contractAddress, transactionHash, transactionCost, errors.New("Error when getting nonce , ERROR : " + errWhenGettingNonce.Error())
	}

	gasPrice, errWhenGettingGasPrice := GetCurrentGasPrice()
	if errWhenGettingGasPrice != nil {
		logrus.Error("Error when getting gas price " + errWhenGettingGasPrice.Error())
		return contractAddress, transactionHash, transactionCost, errors.New("Error when getting gas price, ERROR : " + errWhenGettingGasPrice.Error())
	}

	//create the keyed transactor
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = big.NewInt(int64(gasPrice))

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
		logrus.Info("GetABI returned nil")
		return contractAddress, transactionHash, transactionCost, errors.New("Error when getting ABI string , ERROR : GetAbi() returned nil")
	}

	address, tx, contract, errWhenDeployingContract := bind.DeployContract(auth, *parsed, common.FromHex(ContractBIN), client)
	if errWhenDeployingContract != nil {
		logrus.Info("Error when deploying contract " + errWhenDeployingContract.Error())
		return contractAddress, transactionHash, transactionCost, errors.New("Error when deploying contract, ERROR : " + errWhenDeployingContract.Error())
	}

	contractAddress = address.Hex()
	transactionHash = tx.Hash().Hex()
	_ = contract

	logrus.Info("View contract at : https://sepolia.etherscan.io/address/", address.Hex())
	logrus.Info("View transaction at : https://sepolia.etherscan.io/tx/", tx.Hash().Hex())

	// Wait for the transaction to be mined and calculate the cost
	receipt, errInGettingReceipt := bind.WaitMined(context.Background(), client, tx)
	if errInGettingReceipt != nil {
		logrus.Error("Error in getting receipt: Error: " + errInGettingReceipt.Error())
	} else {
		costInWei := new(big.Int).Mul(big.NewInt(int64(receipt.GasUsed)), big.NewInt(int64(gasPrice)))
		cost := new(big.Float).Quo(new(big.Float).SetInt(costInWei), big.NewFloat(math.Pow10(18)))
		transactionCost = fmt.Sprintf("%g", cost) + " ETH"

		if receipt.Status == 0 {
			errorMessageFromStatus, errorInCallingTransactionStatus := GetErrorOfFailedTransaction(tx.Hash().Hex())
			if errorInCallingTransactionStatus != nil {
				logrus.Error("Transaction failed.")
				logrus.Error("Error when getting the error for the transaction failure: Error: " + errorInCallingTransactionStatus.Error())
				return contractAddress, transactionHash, transactionCost, errors.New("Transaction failed.")
			} else {
				logrus.Error("Transaction failed. Error: " + errorMessageFromStatus)
				return contractAddress, transactionHash, transactionCost, errors.New("Transaction failed. Error: " + errorMessageFromStatus)
			}
		}
	}
	return contractAddress, transactionHash, transactionCost, nil
}
