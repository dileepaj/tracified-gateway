package businessFacades

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/pools"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func BatchConvertCoin(w http.ResponseWriter, r *http.Request) {
	var newBatchConvertCoinObj model.BatchCoinConvert

	err := json.NewDecoder(r.Body).Decode(&newBatchConvertCoinObj)
	if err != nil {
		fmt.Println(err)
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
	data, _ := object.GetBatchSpecificAccount(newBatchConvertCoinObj.BatchID).Then(func(data interface{}) interface{} {
		return data
	}).Await()

	if data == nil {
		// if not create the sponsering account

		// add account to the DB
		batchAccount := model.BatchAccount{
			BatchID:        newBatchConvertCoinObj.BatchID,
			TenentID:       newBatchConvertCoinObj.TenantId,
			ProductName:    newBatchConvertCoinObj.ProductName,
			EquationID:     newBatchConvertCoinObj.EquationID,
			BatchName:      newBatchConvertCoinObj.BatchName,
			BatchAccountPK: "NewAcc2PK",
			BatchAccountSK: "NewAcc2SK",
		}

		object := dao.Connection{}
		err := object.InsertBatchAccount(batchAccount)
		if err != nil {
			log.Println("Error when inserting batch acccount to DB " + err.Error())
		} else {
			log.Println("Batch account added to the DB")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			result := apiModel.SubmitXDRSuccess{
				Status: "Batch account added to the DB",
			}
			json.NewEncoder(w).Encode(result)
			return
		}

		// call path payment methods

	} else {
		// if there is an account go to path payments directly
		batchAccountPK := (data.(model.BatchAccount)).BatchAccountPK
		fmt.Println(batchAccountPK)

		batchAccountSK := (data.(model.BatchAccount)).BatchAccountSK
		fmt.Println(batchAccountSK)

		// call path payement methods
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Batch account is already created",
		}
		json.NewEncoder(w).Encode(result)
		return
	}
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

	response:=model.BuildPoolResponse{
		MetricId:equationJsonObj.MetricID,
		EquationId: equationJsonObj.EquationID,
		TenantId:   equationJsonObj.TenantID,
		BuildPools: cratedPools,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	return
}
