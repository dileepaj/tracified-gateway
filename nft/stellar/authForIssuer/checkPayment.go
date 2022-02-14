package authForIssuer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
	//"github.com/dileepaj/tracified-gateway/model"
	//"github.com/stellar/go/xdr"
)

func CheckPayment(hash string) bool {
	log.Println("Inside Check Payment!")
	fmt.Println(hash)

	var result bool

	//decode the hash to get memo
	result1, err := http.Get(commons.GetHorizonClient().URL + "/transactions/" + hash)
	if err != nil {
		log.Println("Error while Decoding" + err.Error())
	}
	data, _ := ioutil.ReadAll(result1.Body)
	if result1.StatusCode == 200 {
		var txn model.StellarTransaction
		error1 := json.Unmarshal(data, &txn)
		if error1 != nil {
			log.Println("Error while unmarshalling")
		} else {
			if *&txn.Memo == "Payment has been made!" {
				result = true
			} else {
				result = false
			}
		}
		fmt.Println(result)
		return result

	}
	logrus.Error("Transaction is not valid")
	return false
}
