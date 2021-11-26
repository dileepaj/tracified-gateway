package stellarRetriever

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
)

type ConcreteStellarTransaction struct {
	Txnhash string
}

func (stxn *ConcreteStellarTransaction) RetrieveTransaction() (*model.StellarTransaction, error) {
	result, err := http.Get(commons.GetHorizonClient().URL + "/transactions/" + stxn.Txnhash)
	if err != nil {
		return nil, err
	}
	data, _ := ioutil.ReadAll(result.Body)
	if result.StatusCode == 200 {
		var txn model.StellarTransaction
		error := json.Unmarshal(data, &txn)
		if error != nil {
			return nil, error
		}
		return &txn, nil
	}
	return nil, errors.New("Transaction is not valid")
}

func (stxn *ConcreteStellarTransaction) RetrieveOperations() (*model.StellarOperations, error) {
	result, err := http.Get(commons.GetHorizonClient().URL + "/transactions/" + stxn.Txnhash + "/operations")
	if err != nil {
		return nil, err
	}
	data, _ := ioutil.ReadAll(result.Body)
	if result.StatusCode == 200 {
		var oprn model.StellarOperations
		error := json.Unmarshal(data, &oprn)
		if error != nil {
			return nil, error
		}
		return &oprn, nil
	}
	return nil, errors.New("Transaction is not valid")
}