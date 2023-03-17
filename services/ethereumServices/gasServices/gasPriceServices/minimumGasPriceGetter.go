package gasPriceServices

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

type BlockNumberResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

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
	client, errDiallingClient := ethclient.Dial(commons.GoDotEnvVariable("ETHEREUMMAINNETLINK"))
	if errDiallingClient != nil {
		logrus.Error("Error when dialing client : " + errDiallingClient.Error())
		return nil, errors.New(errDiallingClient.Error())
	}

	//call the block number endpoint
	blockUrl := `https://api.etherscan.io/api?module=block&action=getblocknobytime&timestamp=` + strconv.FormatInt(time.Now().Unix(), 10) + `&closest=before&apikey=AER6M2C3436231IGT7SV7JZ2URFYFX7MZ1`

	blockNoRes, errWhenGettingBlockNumber := http.Get(blockUrl)
	if errWhenGettingBlockNumber != nil {
		logrus.Error("Error when getting block number : " + errWhenGettingBlockNumber.Error())
		return nil, errors.New(errWhenGettingBlockNumber.Error())
	}
	defer blockNoRes.Body.Close()
	body, errWhenReadingTheBlockNo := ioutil.ReadAll(blockNoRes.Body)
	if errWhenReadingTheBlockNo != nil {
		logrus.Error("Error when reading the block number : " + errWhenReadingTheBlockNo.Error())
		return nil, errors.New(errWhenReadingTheBlockNo.Error())
	}

	var blockNoJsonResponse BlockNumberResponse
	errWhenUnmarsallingBlock := json.Unmarshal(body, &blockNoJsonResponse)
	if errWhenUnmarsallingBlock != nil {
		logrus.Error("Error when unmarshalling block : " + errWhenUnmarsallingBlock.Error())
		return nil, errors.New(errWhenUnmarsallingBlock.Error())
	}

	var i big.Int
	_, suc := i.SetString(blockNoJsonResponse.Result, 10)
	if !suc {
		logrus.Error("Error when converting block number to big int")
		return nil, errors.New("Error when converting block number to big int")
	}

	//load th block
	block, err := client.BlockByNumber(context.Background(), &i)
	if err != nil {
		logrus.Error("Error when loading the block : " + err.Error())
		return nil, errors.New(err.Error())
	}

	count := 0

	min := block.Transactions()[0].GasPrice()
	//query block transactions
	for _, tx := range block.Transactions() {
		//pick only the normal transaction by checking if the "To" is nil
		if tx.To() != nil {
			if tx.GasPrice().Cmp(min) < 0 {
				min = tx.GasPrice()
			}
			count++
		}
	}

	logrus.Info("No of unique transactions considered to get the minimum gas price : " + strconv.Itoa(count))

	return min, nil

}
