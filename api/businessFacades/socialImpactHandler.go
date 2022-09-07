package businessFacades

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

// BuildSocialImpactFormula --> this method take the farmulaJSOn from backend and map the equationID to 4 byte value and store it DB,
// generate  executionTemplete object via passing FCL query to FCL library,
// base on the executionTemplete objcet store the formula in stellar blckchin
func BuildSocialImpactFormula(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var formulaJSON model.BuildFormula

	err := json.NewDecoder(r.Body).Decode(&formulaJSON)
	if err != nil {
		logrus.Error(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while decoding the body")
		return
	}
	url := "https://horizon-testnet.stellar.org/transactions/db3d8cb1bfd1f393f0691c3efaf1101336f26309f997ef48c4164f8b8c6f18ce/operations"
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
	keys := make([]PublicKeyPOCOC, 0)
	err = json.Unmarshal(keysBody, &keys)
	if err != nil {
		logrus.Error("Unable to unmarshal keys data")
	}
	fmt.Println("--raw1----- ", keys)
	for i := range keys {
		fmt.Println(keys[i])
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(formulaJSON)
	return
}
