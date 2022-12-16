package deploy

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math"
	"math/big"
	"reflect"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/ethereum/go-ethereum"
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
	client, errWhenDialingEthClinet := ethclient.Dial(commons.GoDotEnvVariable("SPOLIALINK"))
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

	//get the gas price
	gasPrice, errWhenGettingGasPrice := client.EstimateGas(context.Background(), ethereum.CallMsg{})
	if errWhenGettingGasPrice != nil {
		logrus.Error("Error when getting gas price " + errWhenGettingGasPrice.Error())
		return contractAddress, transactionHash, transactionCost, errors.New("Error when getting gas price , ERROR : " + errWhenGettingGasPrice.Error())
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
			logrus.Error("Transaction failed.")
			return contractAddress, transactionHash, transactionCost, errors.New("Transaction failed.")
		}
	}
	return contractAddress, transactionHash, transactionCost, nil
}

/*
	Call the deploy method in a dynamic manner
*/
func InvokeDeploy(funcName string, auth *bind.TransactOpts, backend bind.ContractBackend) (string, string, error) {
	type MethodCall struct{}
	var m MethodCall

	params := []reflect.Value{
		reflect.ValueOf(auth),
		reflect.ValueOf(backend),
	}
	result := reflect.ValueOf(&m).MethodByName(funcName).Call(params)
	address := result[0].Interface()
	hash := result[1].Interface()
	errDeploy := result[3].Interface()
	if errDeploy != nil {
		logrus.Error("Error in the deployment method ", errDeploy.(error))
		return "", "", errDeploy.(error)
	}

	return address.(string), hash.(string), nil
}
