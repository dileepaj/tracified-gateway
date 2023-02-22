package gasPriceServices

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net/http"
	"strconv"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

// getting minimum gas price from the latest transactions and get the amount 10% less than that (using infura, etherscan and go-ethereum)

type Response struct {
	Status  string
	Message string
	Result  []struct {
		BlockNumber     string
		TimeStamp       string
		Hash            string
		From            string
		To              string
		Value           string
		ContractAddress string
		Input           string
		Type            string
		Gas             string
		GasUsed         string
		TraceId         string
		IsError         string
		ErrCode         string
	}
}

func GetMinGasPrice() (*big.Int, error) {
	client1, errDiallingClient := ethclient.Dial(commons.GoDotEnvVariable("ETHEREUMMAINNETLINK"))
	if errDiallingClient != nil {
		logrus.Error(errDiallingClient)
		return nil, errors.New(errDiallingClient.Error())
	}

	offSet := 100
	urlToGetTransactions := "https://api.etherscan.io/api?apikey=" + commons.GoDotEnvVariable("ETHERSCANAPIKEY") + "&module=account&action=txlistinternal&startblock=latest&sort=desc&offset=" + strconv.Itoa(offSet) + "&page=1&endblock=latest+1"
	method := "GET"
	httpClient := &http.Client{}
	request, errInNewRequest := http.NewRequest(method, urlToGetTransactions, nil)
	if errInNewRequest != nil {
		logrus.Error("Error in creating new request: " + errInNewRequest.Error())
		return nil, errors.New(errInNewRequest.Error())
	}

	response, errInDo := httpClient.Do(request)
	if errInDo != nil {
		logrus.Error("Error in getting response: " + errInDo.Error())
		return nil, errors.New(errInDo.Error())
	}

	defer response.Body.Close()

	body, errInReadAll := io.ReadAll(response.Body)
	if errInReadAll != nil {
		logrus.Error("Error in reading response body: " + errInReadAll.Error())
		return nil, errors.New(errInReadAll.Error())
	}

	var response1 Response
	errInUnmarshal := json.Unmarshal(body, &response1)
	if errInUnmarshal != nil {
		logrus.Error("Error in unmarshalling: " + errInUnmarshal.Error())
		return nil, errors.New(errInUnmarshal.Error())
	}

	min := new(big.Int)
	if len(response1.Result) > 0 {
		// getting all the hashes to an array
		var hashes []string
		for result1 := range response1.Result {
			if response1.Result[result1].IsError == "0" {
				hashes = append(hashes, response1.Result[result1].Hash)
			}
		}

		// remove duplicates from hashes
		uniqueHashes := getUniqueStringsInAnArray(hashes)
		logrus.Info("No of Unique Hashes : ", len(uniqueHashes))

		tx1, _, errInGettingTransaction1 := client1.TransactionByHash(context.Background(), common.HexToHash(response1.Result[0].Hash))
		if errInGettingTransaction1 != nil {
			logrus.Error("Error in getting transaction: " + errInGettingTransaction1.Error())
			return nil, errors.New(errInGettingTransaction1.Error())
		}

		min = tx1.GasPrice()
		for _, hash := range uniqueHashes {
			tx2, _, errInGettingTransaction := client1.TransactionByHash(context.Background(), common.HexToHash(hash))
			if errInGettingTransaction != nil {
				logrus.Error("Error in getting transaction: " + errInGettingTransaction.Error())
				return nil, errors.New(errInGettingTransaction.Error())
			}

			if tx2.GasPrice().Cmp(min) < 0 {
				min = tx2.GasPrice()
			}

		}
	} else {
		logrus.Error("No transactions found")
		logrus.Info("Using the lowest gas price from the network")
		lowestPrice, errorInGettingLowestPrice := GetCurrentGasPrice()
		if errorInGettingLowestPrice != nil {
			logrus.Error(errorInGettingLowestPrice)
			return nil, errors.New(errorInGettingLowestPrice.Error())
		}
		min = big.NewInt(int64(lowestPrice))
	}
	logrus.Info("Initial gas price : ", min)
	// get the less than 10% value as the minimum gas price
	min = new(big.Int).Sub(min, new(big.Int).Div(min, big.NewInt(10)))

	return min, nil
}

func getUniqueStringsInAnArray(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
