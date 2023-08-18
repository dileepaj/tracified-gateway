package transactionrecipthandler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/utilities"
)

type Log struct {
	TransactionHash  string   `json:"transactionHash"`
	Address          string   `json:"address"`
	BlockHash        string   `json:"blockHash"`
	BlockNumber      string   `json:"blockNumber"`
	Data             string   `json:"data"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
	Topics           []string `json:"topics"`
	TransactionIndex string   `json:"transactionIndex"`
}

// Result represents the result section of the JSON data.
type Result struct {
	TransactionHash   string `json:"transactionHash"`
	BlockHash         string `json:"blockHash"`
	BlockNumber       string `json:"blockNumber"`
	Logs              []Log  `json:"logs"`
	ContractAddress   string `json:"contractAddress"`
	EffectiveGasPrice string `json:"effectiveGasPrice"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	From              string `json:"from"`
	GasUsed           string `json:"gasUsed"`
	LogsBloom         string `json:"logsBloom"`
	Status            string `json:"status"`
	To                string `json:"to"`
	TransactionIndex  string `json:"transactionIndex"`
	Type              string `json:"type"`
}

type PolygonTransactionReceiptResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  Result `json:"result"`
}

func GetTransactionReceiptForPolygon(transactionHash string) (Result, error) {
	url := commons.GoDotEnvVariable("POLYGONALCHEMYAPILINK") + commons.GoDotEnvVariable("POLYGONALCHEMYAPIKEY")
	method := "POST"
	logger := utilities.NewCustomLogger()

	payload := strings.NewReader("{\"params\":[\"" + transactionHash + "\"],\"id\":1,\"jsonrpc\":\"2.0\",\"method\":\"eth_getTransactionReceipt\"}")

	client := &http.Client{}
	req, errNewRequest := http.NewRequest(method, url, payload)
	if errNewRequest != nil {
		logger.LogWriter("Error in creating new request: "+errNewRequest.Error(), constants.ERROR)
		return Result{}, errNewRequest
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		logger.LogWriter("Error in sending request: "+err.Error(), constants.ERROR)
		return Result{}, err
	}
	defer res.Body.Close()

	body, errInReadAll := ioutil.ReadAll(res.Body)
	if errInReadAll != nil {
		logger.LogWriter("Error in reading response body: "+errInReadAll.Error(), constants.ERROR)
		return Result{}, errInReadAll
	}

	var responseBody PolygonTransactionReceiptResponse

	// unmarshal the response
	errorInUnmarshal := json.Unmarshal(body, &responseBody)
	if errorInUnmarshal != nil {
		logger.LogWriter("Error in un-marshalling response body: "+errorInUnmarshal.Error(), constants.ERROR)
		return Result{}, errorInUnmarshal
	}

	//check if the result is null
	if responseBody.Result.Status == "" {
		logger.LogWriter("The transaction is not yet added to blockchain", constants.ERROR)
		return Result{}, errors.New("not found")
	}

	return responseBody.Result, nil
}
