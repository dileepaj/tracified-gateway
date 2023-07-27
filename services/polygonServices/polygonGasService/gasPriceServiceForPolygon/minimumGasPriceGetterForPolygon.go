package gaspriceserviceforpolygon

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
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/services/ethereumServices/gasServices/gasPriceServices"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/ethereum/go-ethereum/ethclient"
)

func MinimumGasPriceGetterForPolygon() (*big.Int, error) {
	logger := utilities.NewCustomLogger()

	client, errWhenDialingClient := ethclient.Dial(commons.GoDotEnvVariable("POLYGONALCHEMYMAINNETLINK"))
	if errWhenDialingClient != nil {
		logger.LogWriter("Error when dialing client : "+errWhenDialingClient.Error(), constants.ERROR)
		return nil, errors.New(errWhenDialingClient.Error())
	}

	blockUrl := commons.GoDotEnvVariable("POLYGONSCANTESTNETAPI") + `?module=block&action=getblocknobytime&timestamp=` + strconv.FormatInt(time.Now().Unix(), 10) + `&closest=before&apikey=` + commons.GoDotEnvVariable("16RXSK5PUS2S428HSFBXP4C1GV81EK3J4F")

	blockNoRes, errWhenGettingBlockNo := http.Get(blockUrl)
	if errWhenGettingBlockNo != nil {
		logger.LogWriter("Error when getting block No :"+errWhenGettingBlockNo.Error(), constants.ERROR)
		return nil, errors.New("Error when getting block No : " + errWhenGettingBlockNo.Error())
	}
	defer blockNoRes.Body.Close()

	body, errWhenReadingTheBlockNo := ioutil.ReadAll(blockNoRes.Body)
	if errWhenReadingTheBlockNo != nil {
		logger.LogWriter("Error when reading the block number : "+errWhenReadingTheBlockNo.Error(), constants.ERROR)
		return nil, errors.New(errWhenReadingTheBlockNo.Error())
	}

	var blockNoJsonResponse gasPriceServices.BlockNumberResponse
	errWhenUnmarsallingBlock := json.Unmarshal(body, &blockNoJsonResponse)
	if errWhenUnmarsallingBlock != nil {
		logger.LogWriter("Error when unmarshalling block : "+errWhenUnmarsallingBlock.Error(), constants.ERROR)
		return nil, errors.New(errWhenUnmarsallingBlock.Error())
	}

	var i big.Int
	_, suc := i.SetString(blockNoJsonResponse.Result, 10)
	if !suc {
		logger.LogWriter("Error when converting block number to big int", constants.ERROR)
		return nil, errors.New("Error when converting block number to big int")
	}

	block, err := client.BlockByNumber(context.Background(), &i)
	if err != nil {
		logger.LogWriter("Error when loading the block : "+err.Error(), constants.ERROR)
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

	logger.LogWriter("No of unique transactions considered to get the minimum gas price : "+strconv.Itoa(count), constants.INFO)

	return min, nil
}
