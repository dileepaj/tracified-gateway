package gasservices

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/sirupsen/logrus"
)

//from, to, gas, gasPrice, maxPriorityFeePerGas, maxFeePerGas, value, data
// EstimateGasLimit estimates the gas limit for a transaction
// It takes parameters: from[optional], to, gas[optional], gasPrice[optional], maxPriorityFeePerGas[optional], maxFeePerGas[optional], value[optional], data[optional]
// to is the address the transaction is directed to or null if deploying a contract
// data is the bin of the contract to deploy

func EstimateGasLimit(from string, to string, gas string, gasPrice string, maxPriorityFeePerGas string, maxFeePerGas string, value string, data string) (uint64, error) {
	url := commons.GoDotEnvVariable("ETHEREUMTESTNETLINK")
	method := "POST"

	// if any of the parameters are empty, set them to null
	if from == "" {
		from = "null"
	}
	if to == "" {
		to = "null"
	}
	if gas == "" {
		gas = "null"
	}
	if gasPrice == "" {
		gasPrice = "null"
	}
	if maxPriorityFeePerGas == "" {
		maxPriorityFeePerGas = "null"
	}
	if maxFeePerGas == "" {
		maxFeePerGas = "null"
	}
	if value == "" {
		value = "null"
	}
	if data == "" {
		data = "null"
	}

	payload := strings.NewReader(`{` +
		`"jsonrpc": "2.0",` +
		`"method": "eth_estimateGas",` +
		`"params": [` +
		`{` +
		`"from": "` + from + `",` +
		`"to": ` + to + `,` +
		`"gas": "`+ gas + `",` +
		`"gasPrice": "`+ gasPrice +`",` +
		`"maxPriorityFeePerGas": "`+ maxPriorityFeePerGas +`",` +
		`"maxFeePerGas": "`+ maxFeePerGas +`",` +
		`"value": "` + value + `",` +
		`"data": "` + data + `"` +
		`}` +
		`],` +
		`"id": 1` +
		`}`)

	client := &http.Client{}
	req, errNewRequest := http.NewRequest(method, url, payload)
	if errNewRequest != nil {
		logrus.Error("Error in creating new request: " + errNewRequest.Error())
		return 0, errNewRequest
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		logrus.Error("Error in sending request: " + err.Error())
		return 0, err
	}
	defer res.Body.Close()

	body, errInReadAll := ioutil.ReadAll(res.Body)
	if errInReadAll != nil {
		logrus.Error("Error in reading response body: " + errInReadAll.Error())
		return 0, errInReadAll
	}

	// struct to unmarshal the response
	type GasLimitResponse struct {
		Jsonrpc string
		Id int
		Result  string
	}

	var gasLimitResponse GasLimitResponse

	// unmarshal the response
	errorInUnmarshal := json.Unmarshal(body, &gasLimitResponse)
	if errorInUnmarshal != nil {
		logrus.Error("Error in un-marshalling response body: " + errorInUnmarshal.Error())
		return 0, errorInUnmarshal
	}

	// remove 0x from the hex string
	hexString := strings.Replace(gasLimitResponse.Result, "0x", "", 1)

	// convert hex to decimal
	decimalValue, errInConversion := strconv.ParseInt(hexString, 16, 64)
	if errInConversion != nil {
		logrus.Error("Error in converting gas limit hex to decimal: " + err.Error())
		return 0, errInConversion
	}

	// add 10% to the gas limit
	safeValue := decimalValue + (decimalValue * 10 / 100)
	logrus.Info("Estimated gas limit: ", decimalValue)
	logrus.Info("Estimated safe gas limit: ", safeValue)

	return uint64(safeValue), nil
}