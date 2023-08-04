package polygongasservice

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/utilities"
)

func EstimateGasLimitForPolygon(from string, to string, gas string, gasPrice string, value string, data string) (uint64, error) {
	url := commons.GoDotEnvVariable("POLYGONALCHEMYTESTNETLINK") + commons.GoDotEnvVariable("POLYGONALCHEMYAPIKEY")
	method := "POST"
	logger := utilities.NewCustomLogger()

	var fromString, toString, gasString, gasPriceString, valueString, dataString string

	//clear empty parameters
	if from == "" {
		fromString = `"from": null,`
	} else {
		fromString = `"from": "` + from + `",`
	}

	if to == "" {
		toString = `"to": null,`
	} else {
		toString = `"to": "` + to + `",`
	}

	if gas == "" {
		gasString = `"gas": null,`
	} else {
		gasString = `"gas": "` + gas + `",`
	}

	if gasPrice == "" {
		gasPriceString = `"gasPrice": null,`
	} else {
		gasPriceString = `"gasPrice": "` + gasPrice + `",`
	}

	if value == "" {
		valueString = `"value": null,`
	} else {
		valueString = `"value": "` + value + `",`
	}

	if data == "" {
		dataString = `"data": null`
	} else {
		dataString = `"data": "` + data + `"`
	}

	payload := strings.NewReader(`{` +
		`"jsonrpc": "2.0",` +
		`"method": "eth_estimateGas",` +
		`"params": [` +
		`{` +
		fromString +
		toString +
		gasString +
		gasPriceString +
		valueString +
		dataString +
		`}` +
		`],` +
		`"id": 1` +
		`}`)

	client := &http.Client{}
	req, errNewRequest := http.NewRequest(method, url, payload)
	if errNewRequest != nil {
		logger.LogWriter("Error in creating new request: "+errNewRequest.Error(), constants.ERROR)
		return 0, errNewRequest
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		logger.LogWriter("Error in sending request: "+err.Error(), constants.ERROR)
		return 0, err
	}
	defer res.Body.Close()

	body, errInReadAll := ioutil.ReadAll(res.Body)
	if errInReadAll != nil {
		logger.LogWriter("Error in reading response body: "+errInReadAll.Error(), constants.ERROR)
		return 0, errInReadAll
	}

	// struct to unmarshal the response
	type GasLimitResponse struct {
		Id      int
		Jsonrpc string
		Result  string
	}

	var gasLimitResponse GasLimitResponse

	// unmarshal the response
	errorInUnmarshal := json.Unmarshal(body, &gasLimitResponse)
	if errorInUnmarshal != nil {
		logger.LogWriter("Error in un-marshalling response body: "+errorInUnmarshal.Error(), constants.ERROR)
		return 0, errorInUnmarshal
	}

	// remove 0x from the hex string
	hexString := strings.Replace(gasLimitResponse.Result, "0x", "", 1)

	// convert hex to decimal
	decimalValue, errInConversion := strconv.ParseInt(hexString, 16, 64)
	if errInConversion != nil {
		logger.LogWriter("Error in converting gas limit hex to decimal: "+errInConversion.Error(), constants.ERROR)
		return 0, errInConversion
	}

	// add 10% to the gas limit
	safeValue := decimalValue + (decimalValue * 10 / 100)
	logger.LogWriter("Estimated gas limit: "+strconv.FormatInt(decimalValue, 10), constants.INFO)
	logger.LogWriter("Estimated safe gas limit: "+strconv.FormatInt(safeValue, 10), constants.INFO)

	return uint64(safeValue), nil
}
