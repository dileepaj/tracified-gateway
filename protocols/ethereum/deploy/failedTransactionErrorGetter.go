package deploy

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/sirupsen/logrus"
)

type TransactionStatus struct {
	Status  string
	Message string
	Result  struct {
		IsError        string
		ErrDescription string
	}
}

func GetErrorOfFailedTransaction(transactionHash string) (string, error) {
	url := commons.GoDotEnvVariable("ETHERSCANAPISEPOLIATESTNET") + "api?module=transaction&action=getstatus&apikey=" + commons.GoDotEnvVariable("ETHERSCANAPIKEY") + "&txhash=" + transactionHash
	// url := "https://api-goerli.etherscan.io/api?module=transaction&action=getstatus&apikey=AER6M2C3436231IGT7SV7JZ2URFYFX7MZ1&txhash=" + transactionHash
	logrus.Info("Calling the transaction status getter at : " + url)
	method := "GET"

	logrus.Info("---------------Waiting 30 seconds for updating the transaction status...")
	time.Sleep(30 * time.Second)

	client := &http.Client{}
	req, errWhenCallingTheNewRequest := http.NewRequest(method, url, nil)
	if errWhenCallingTheNewRequest != nil {
		logrus.Error("Error when calling the new request : ", errWhenCallingTheNewRequest.Error())
		return "", errors.New("Error when calling the new request : " + errWhenCallingTheNewRequest.Error())
	}

	res, errWhenCallingUrl := client.Do(req)
	if errWhenCallingUrl != nil {
		logrus.Error("Error when calling the transaction url : ", errWhenCallingUrl.Error())
		return "", errors.New("Error when calling the transaction url : " + errWhenCallingUrl.Error())
	}
	defer res.Body.Close()

	body, errWhenReadingTheResponse := ioutil.ReadAll(res.Body)
	if errWhenReadingTheResponse != nil {
		logrus.Error("Error when reading the responce : ", errWhenReadingTheResponse.Error())
		return "", errors.New("Error when reading the responce : " + errWhenReadingTheResponse.Error())
	}

	var transactionStatus TransactionStatus
	errorWhenUnmarshalling := json.Unmarshal(body, &transactionStatus)
	if errorWhenUnmarshalling != nil {
		logrus.Error("Error when unmarshalling the response : ", errorWhenUnmarshalling.Error())
		return "", errors.New("Error when unmarshalling the response : " + errorWhenUnmarshalling.Error())
	}

	return transactionStatus.Result.ErrDescription, nil
}
