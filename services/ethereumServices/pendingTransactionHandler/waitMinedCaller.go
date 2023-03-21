package pendingTransactionHandler

import (
	"context"
	"errors"
	"fmt"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

// returns status, usedGasAmount, error
func CallWaitMined(client *ethclient.Client, tx *types.Transaction) (int64, uint64, error) {

	object := dao.Connection{}

	receipt, errInGettingReceipt := bind.WaitMined(context.Background(), client, tx)
	if errInGettingReceipt != nil {
		logrus.Error("Error in getting receipt: Error: " + errInGettingReceipt.Error())
		return -1, 0, errors.New("Error in getting receipt: Error: " + errInGettingReceipt.Error())
	} else {
		if receipt.Status == 0 {
			errorMessageFromStatus, errorInCallingTransactionStatus := GetErrorOfFailedTransaction(tx.Hash().Hex())
			if errorInCallingTransactionStatus != nil {
				logrus.Error("Transaction failed.")
				logrus.Error("Error when getting the error for the transaction failure: Error: " + errorInCallingTransactionStatus.Error())
				return int64(receipt.Status), receipt.GasUsed, errors.New("Transaction failed.")
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
			return int64(receipt.Status), receipt.GasUsed, nil
		} else {
			logrus.Error("Invalid receipt status for 'WaitMined', Status : ", receipt.Status)
			return int64(receipt.Status), receipt.GasUsed, errors.New("Invalid receipt status for 'WaitMined', Status : " + fmt.Sprint(receipt.Status))
		}
		logrus.Info("Status of receipt : ", receipt.Status)
	}
	return int64(receipt.Status), receipt.GasUsed, nil
}