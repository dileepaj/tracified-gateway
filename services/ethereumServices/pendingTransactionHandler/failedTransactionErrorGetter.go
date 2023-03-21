package pendingTransactionHandler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

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
	url := commons.GoDotEnvVariable("ETHERSCANAPITESTNET") + "api?module=transaction&action=getstatus&apikey=" + commons.GoDotEnvVariable("ETHERSCANAPIKEY") + "&txhash=" + transactionHash
	logrus.Info("Calling the transaction status getter at : " + url)
	method := "GET"

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