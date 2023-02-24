package pendingTransactionHandler

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// Check if the transaction is in blockchain or not
func CheckTransaction(tx string) (bool, error) {
	isInBlockchain := false

	url := commons.GoDotEnvVariable("ETHERSCANAPISEPOLIATESTNET") + "api?module=account&action=txlist&address=" + commons.GoDotEnvVariable("ETHEREUMPUBKEY") + "&startblock=0&endblock=latest+1&page=1&offset=1&sort=desc&apikey=" + commons.GoDotEnvVariable("ETHERSCANAPIKEY")
	logrus.Info("Getting the latest transaction from the address : " + url)

	result, errWhenCallingUrl := http.Get(url)
	if errWhenCallingUrl != nil {
		logrus.Error("Error when calling the url to get the last transaction : ", errWhenCallingUrl.Error())
		return false, errors.New("Error when calling the url to get the last transaction  : " + errWhenCallingUrl.Error())
	}
	
	if result.StatusCode != 200 {
		logrus.Error("Getting last transaction response status code : ", result.StatusCode)
		return false, errors.New("Getting last transaction response status code : " + strconv.Itoa(result.StatusCode))
	}

	data, errWhenReadingTheResult := io.ReadAll(result.Body)
	if errWhenReadingTheResult != nil {
		logrus.Error("Error when reading the response : ", errWhenReadingTheResult.Error())
		return false, errors.New("Error when reading the response : " + errWhenReadingTheResult.Error())
	}
	defer result.Body.Close()

	value := gjson.GetBytes(data, "result.0.hash")
	logrus.Info("Last transaction hash in the account : ", value.String())

	if value.String() == tx {
		isInBlockchain = true
	}

	return isInBlockchain, nil
}