package gasPriceServices

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

// get the lowest gas price using an external api

func GetCurrentGasPrice() (int, error) {
	url := commons.GoDotEnvVariable("ETHERSCANAPI") + "api?module=gastracker&action=gasoracle&apikey=" + commons.GoDotEnvVariable("ETHERSCANAPIKEY")
	logrus.Info("Calling gas tracker to get the optimum gas price at : " + url)
	result, errWhenCallingUrl := http.Get(url)
	if errWhenCallingUrl != nil {
		logrus.Error("Error when calling the gas tracker url : ", errWhenCallingUrl.Error())
		return 0, errors.New("Error when calling the gas tracker url : " + errWhenCallingUrl.Error())
	}

	if result.StatusCode != 200 {
		logrus.Error("Gas tracker response status code : ", result.StatusCode)
		return 0, errors.New("Gas tracker response status code : " + strconv.Itoa(result.StatusCode))
	}

	data, errWhenReadingTheResult := io.ReadAll(result.Body)
	if errWhenReadingTheResult != nil {
		logrus.Error("Error when reading the gas tracker response : ", errWhenReadingTheResult.Error())
		return 0, errors.New("Error when reading the gas tracker response : " + errWhenReadingTheResult.Error())
	}
	defer result.Body.Close()

	value := gjson.GetBytes(data, "result.FastGasPrice")
	logrus.Info("Gas price in Gwei : ", value.String())

	//convert value to int
	intValue, errWhenConvertingToInt := strconv.Atoi(value.String())
	if errWhenConvertingToInt != nil {
		logrus.Error("Error when converting to int : ", errWhenConvertingToInt.Error())
		return 0, errors.New("Error when converting to int : " + errWhenConvertingToInt.Error())
	}

	finalWeiAmount := intValue * 1000000000
	logrus.Info("Gas price in Wei : ", finalWeiAmount)

	return finalWeiAmount, nil

}
