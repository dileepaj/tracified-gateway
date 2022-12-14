package deploy

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/go-errors/errors"
	"github.com/sirupsen/logrus"
)

type BalanceResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func CheckBalance() (float64, error) {
	var balanceResponse BalanceResponse

	//Create the URL to be called
	URLString := commons.GoDotEnvVariable("SEPOLIALINK")
	URLString += "module=account&action=balance&address=" + commons.GoDotEnvVariable("ETHEREUMPUBKEY") + "&tag=latest&apikey=" + commons.GoDotEnvVariable("TRACIFIEDETHERSCANAPIKEY")

	client := &http.Client{}
	request, _ := http.NewRequest(http.MethodGet, URLString, nil)

	//Send the request and get the response
	response, err := client.Do(request)
	if err != nil {
		logrus.Error("Error in CheckBalance: ", err)
		return 0, err
	}
	defer response.Body.Close()

	//Decode the response
	json.NewDecoder(response.Body).Decode(&balanceResponse)

	if balanceResponse.Status == "1" {
		//Convert the balance to float64
		balanceInFloat, _ := strconv.ParseFloat(balanceResponse.Result, 64)
		//Convert the balance to Ether	(1 Ether = 1000000000000000000 Wei)
		convertedBalance := balanceInFloat / 1000000000000000000
		return convertedBalance, nil
	} else {
		logrus.Error("Error in CheckBalance: ", balanceResponse.Message)
		return 0, errors.New("Error in when checking account balance")
	}
}