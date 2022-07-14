package pools

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	log "github.com/sirupsen/logrus"
)

func CalculateCoin(batchAccount model.CoinAccount) (string, error) {
	url := commons.GetHorizonClient().HorizonURL + "accounts/" + batchAccount.CoinAccountPK
	result, err := http.Get(url)
	if err != nil {
		log.Error("Unable to reach Stellar network", url)
		return "", err
	}
	if result.StatusCode != 200 {
		return "", errors.New(result.Status + " The request you sent to pool assert convertion was invalid in some way" + " " + url)
	}
	defer result.Body.Close()
	assertInfo, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Error(err)
		return "", err
	}

	var raw map[string]interface{}
	var raw1 []interface{}

	json.Unmarshal(assertInfo, &raw)

	out1, _ := json.Marshal(raw["balances"])
	json.Unmarshal(out1, &raw1)

	// retrive the coin converion paths and push it to array
	for i := range raw1 {
		accoutBalance := raw1[i].(map[string]interface{})
		assetCode := fmt.Sprintf("%v", accoutBalance["asset_code"])
		balance := fmt.Sprintf("%v", accoutBalance["balance"])
		if batchAccount.MetricCoin.CoinName == assetCode {
			return balance, nil
		}
	}
	return "", errors.New("Can not find the Coin")
}
