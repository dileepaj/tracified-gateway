package stellarprotocols

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

/*
des - calling the stellar endpoint to get the transaction details using the txn hash
*/

func AccecingManageData(hash string) {
	//! acessing the manage data in blockchain and convert byte to bit string
	url := "https://horizon-testnet.stellar.org/transactions/" + hash + "/operations"
	result, err := http.Get(url)
	if err != nil {
		logrus.Error("Unable to reach Stellar network", url)
	}
	if result.StatusCode != 200 {
		logrus.Error(result)
	}
	defer result.Body.Close()
	assertInfo, err := ioutil.ReadAll(result.Body)
	if err != nil {
		logrus.Error(err)
	}

	var raw map[string]interface{}
	var raw1 map[string]interface{}

	json.Unmarshal(assertInfo, &raw)

	out, err := json.Marshal(raw["_embedded"])
	if err != nil {
		logrus.Error("Unable to marshal embedded")
	}
	err = json.Unmarshal(out, &raw1)
	if err != nil {
		logrus.Error("Unable to unmarshal  json.Unmarshal(out, &raw1)")
	}
	out1, _ := json.Marshal(raw1["records"])
	json.Unmarshal(out1, &raw1)

	keysBody := out1
	keys := make([]model.ManageData, 0)
	err = json.Unmarshal(keysBody, &keys)
	if err != nil {
		logrus.Error("Unable to unmarshal keys data")
	}
	for i := range keys {
		// fmt.Println(keys[i])
		acceptTxn_byteData, err := base64.StdEncoding.DecodeString(keys[i].Value)
		if err != nil {
			logrus.Error(" Decodring issue Value ", keys[i].Value)
		}
		fmt.Println("-----value--------------- ", acceptTxn_byteData)
		acceptTxn := string(acceptTxn_byteData)
		fmt.Println("byte to string value  ", acceptTxn)

		byteVal := acceptTxn[2:10]
		getInt, err := ByteStingToInteger(byteVal)
		if err != nil {
			logrus.Error(err)
		}
		fmt.Println(byteVal, " converted int ", getInt)
	}
}
