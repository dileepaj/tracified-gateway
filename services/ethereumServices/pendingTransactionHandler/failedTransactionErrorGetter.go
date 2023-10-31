package pendingTransactionHandler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/utilities"
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

//1- Ethereum, 2-Polygon
func GetErrorOfFailedTransaction(transactionHash string, blockchainType int) (string, error) {
	var url string
	logger := utilities.NewCustomLogger()
	if blockchainType == 1 {
		url = commons.GoDotEnvVariable("ETHERSCANAPITESTNET") + "api?module=transaction&action=getstatus&apikey=" + commons.GoDotEnvVariable("ETHERSCANAPIKEY") + "&txhash=" + transactionHash
	} else if blockchainType == 2 {
		url = commons.GoDotEnvVariable("POLYGON_SCAN_API_LINK") + "api?module=transaction&action=getstatus&apikey=" + commons.GoDotEnvVariable("POLYGON_SCAN_API_KEY") + "&txhash=" + transactionHash
	} else {
		logger.LogWriter("Invalid blockchain type", constants.ERROR)
		return "", errors.New("Invalid blockchain type")
	}
	logrus.Info("Calling the transaction status getter at : " + url)
	method := "GET"

	client := &http.Client{}
	req, errWhenCallingTheNewRequest := http.NewRequest(method, url, nil)
	if errWhenCallingTheNewRequest != nil {
		logger.LogWriter("Error when calling the new request : "+errWhenCallingTheNewRequest.Error(), constants.ERROR)
		return "", errors.New("Error when calling the new request : " + errWhenCallingTheNewRequest.Error())
	}

	res, errWhenCallingUrl := client.Do(req)
	if errWhenCallingUrl != nil {
		logger.LogWriter("Error when calling the transaction url : "+errWhenCallingUrl.Error(), constants.ERROR)
		return "", errors.New("Error when calling the transaction url : " + errWhenCallingUrl.Error())
	}
	defer res.Body.Close()

	body, errWhenReadingTheResponse := ioutil.ReadAll(res.Body)
	if errWhenReadingTheResponse != nil {
		logger.LogWriter("Error when reading the responce : "+errWhenReadingTheResponse.Error(), constants.ERROR)
		return "", errors.New("Error when reading the responce : " + errWhenReadingTheResponse.Error())
	}

	var transactionStatus TransactionStatus
	errorWhenUnmarshalling := json.Unmarshal(body, &transactionStatus)
	if errorWhenUnmarshalling != nil {
		logger.LogWriter("Error when unmarshalling the response : "+errorWhenUnmarshalling.Error(), constants.ERROR)
		return "", errors.New("Error when unmarshalling the response : " + errorWhenUnmarshalling.Error())
	}

	return transactionStatus.Result.ErrDescription, nil
}
