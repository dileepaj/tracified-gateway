package businessFacades

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/pools"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func BatchConvertCoin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newBatchConvertCoinObj model.BatchCoinConvert
	var batchAccountPK string
	var batchAccountSK string
	var coinConversions []model.BuildPathPayment
	err := json.NewDecoder(r.Body).Decode(&newBatchConvertCoinObj)
	if err != nil {
		logrus.Info(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error while decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		return
	}

	// check if there is an account in the DB for the batchID and get the account
	object := dao.Connection{}
	data, _ := object.GetBatchSpecificAccount(newBatchConvertCoinObj.BatchID, newBatchConvertCoinObj.EquationID, newBatchConvertCoinObj.ProductName).Then(func(data interface{}) interface{} {
		return data
	}).Await()

	if data == nil {
		// if not create the sponsering account
		batchPK, batchSK, err := pools.CreateSponseredAccount()
		batchAccountPK = batchPK
		batchAccountSK = batchSK
		logrus.Info(batchAccountPK)
		logrus.Info(batchAccountSK)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			result := "Can not Create Batch Account " + err.Error()
			json.NewEncoder(w).Encode(result)
		}
		// add account to the DB
		batchAccount := model.BatchAccount{
			BatchID:        newBatchConvertCoinObj.BatchID,
			BatchName:      newBatchConvertCoinObj.BatchName,
			TenentID:       newBatchConvertCoinObj.TenantId,
			ProductName:    newBatchConvertCoinObj.ProductName,
			EquationID:     newBatchConvertCoinObj.EquationID,
			StageID:        newBatchConvertCoinObj.StageId,
			BatchAccountPK: batchAccountPK,
			BatchAccountSK: batchAccountSK,
		}
		object := dao.Connection{}
		errResult := object.InsertBatchAccount(batchAccount)
		if errResult != nil {
			logrus.Info("Error when inserting batch acccount to DB " + errResult.Error())
			w.WriteHeader(http.StatusInternalServerError)
			result := "Error when inserting batch acccount to DB"
			json.NewEncoder(w).Encode(result)
			return
		}
	} else {
		// if there is an account go to path payments directly
		batchAccountPK = (data.(model.BatchAccount)).BatchAccountPK
		logrus.Info(batchAccountPK)

		batchAccountSK = (data.(model.BatchAccount)).BatchAccountSK
		logrus.Info(batchAccountSK)
	}

	//CoinConvertionJson return CoinConvertionJson that used to do a coin convert via pools
	pathpayments, err := pools.CoinConvertionJson(newBatchConvertCoinObj, batchAccountPK, batchAccountSK)
	if err != nil {
		logrus.Error("Can not create Path Payment Json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("Can not create Path Payment Json")
		return
	}

	for _, pathPayment := range pathpayments {
		coinConversion, err := pools.CoinConvert(pathPayment)
		if err != nil {
			logrus.Error("Coin converion issue ", err)
		} else {
			coinConversions = append(coinConversions, coinConversion)
		}
	}
	if len(coinConversions) <= 0 {
		logrus.Info("Can not convert any Coin")
		w.WriteHeader(http.StatusInternalServerError)
		result := "Empty coin convertion"
		json.NewEncoder(w).Encode(result)
		return
	}
	// build response with all coin details
	buildCoinConvertionResponse := model.BuildPathPaymentJSon{
		CoinConertions: coinConversions,
		ProductId:      newBatchConvertCoinObj.ProductID,
		ProductIdName:  newBatchConvertCoinObj.ProductName,
		EquationId:     newBatchConvertCoinObj.EquationID,
		TenantId:       newBatchConvertCoinObj.TenantId,
	}

	// TODO each coinconversion details shoud be saved with timestamp 

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(buildCoinConvertionResponse)
	return
}

func CreatePool(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var equationJsonObj model.CreatePool

	err := json.NewDecoder(r.Body).Decode(&equationJsonObj)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while decoding the body")
		return
	}

	// reformate the equation json
	equationJson, err := pools.RemoveDivisionAndOperator(equationJsonObj)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	// build the pool creation json
	poolCreationJSON, err := pools.BuildPoolCreationJSON(equationJson)
	if err != nil {
		logrus.Error(poolCreationJSON, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	// create the pools
	cratedPools, err := pools.CreatePoolsUsingJson(poolCreationJSON)
	if err != nil {
		logrus.Error(cratedPools, err)
		json.NewEncoder(w).Encode(err)
		return
	}

	response := model.BuildPoolResponse{
		MetricId:   equationJsonObj.MetricID,
		EquationId: equationJsonObj.EquationID,
		TenantId:   equationJsonObj.TenantID,
		BuildPools: cratedPools,
	}

	// insert the pool to the DB
	object := dao.Connection{}
	err1 := object.InsertPool(response)
	if err1 != nil {
		log.Println("Error when inserting pool to DB " + err.Error())
	} else {
		log.Println("Pool added to the DB")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}
}
