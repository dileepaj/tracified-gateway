package businessFacades

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	log "github.com/sirupsen/logrus"
)

func BatchConvertCoin(w http.ResponseWriter, r *http.Request){
	var newBatchConvertCoinObj model.BatchCoinConvert

	err := json.NewDecoder(r.Body).Decode(&newBatchConvertCoinObj)
	if err != nil{
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error while decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		return
	}

	//check if there is an account in the DB for the batchID and get the account
	object := dao.Connection{}
	data,_ := object.GetBatchSpecificAccount(newBatchConvertCoinObj.BatchID).Then(func(data interface{}) interface{}{
		return data
	}).Await()

	if data == nil{
		//if not create the sponsering account

		//add account to the DB
		batchAccount := model.BatchAccount{
			BatchID: newBatchConvertCoinObj.BatchID,
			TenentID: newBatchConvertCoinObj.TenantId,
			ProductName: newBatchConvertCoinObj.ProductName,
			EquationID: newBatchConvertCoinObj.EquationID,
			BatchName: newBatchConvertCoinObj.BatchName,
			BatchAccountPK: "NewAcc2PK",
			BatchAccountSK : "NewAcc2SK",
		}

		object := dao.Connection{}
		err := object.InsertBatchAccount(batchAccount)
		if err != nil{
			log.Println("Error when inserting batch acccount to DB " + err.Error())
		} else{
			log.Println("Batch account added to the DB")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			result := apiModel.SubmitXDRSuccess{
				Status: "Batch account added to the DB",
			}
			json.NewEncoder(w).Encode(result)
			return
		}

		//call path payment methods


	}else{
		//if there is an account go to path payments directly
		batchAccountPK := (data.(model.BatchAccount)).BatchAccountPK
		fmt.Println(batchAccountPK)

		batchAccountSK := (data.(model.BatchAccount)).BatchAccountSK
		fmt.Println(batchAccountSK)

		//call path payement methods
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Batch account is already created",
		}
		json.NewEncoder(w).Encode(result)
		return
	}
	
}

func CreatePool(w http.ResponseWriter, r *http.Request){
	var newCreatePoolObj model.CreatePool

	err := json.NewDecoder(r.Body).Decode(&newCreatePoolObj)
	if err != nil{
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error while decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		return
	}

	//create pool
	// var poolJson []model.BuildPool
	// pool1 := model.BuildPool{
	// 	Coin1:               "BTC",
	// 	DepositeAmountCoin1: "10000",
	// 	Coin2:               "ETH",
	// 	DepositeAmountCoin2: "20000",
	// 	Ratio:               2,
	// }
	// pool2 := model.BuildPool{
	// 	Coin1:               "ETH",
	// 	DepositeAmountCoin1: "10000",
	// 	Coin2:               "USDT",
	// 	DepositeAmountCoin2: "70000",
	// 	Ratio:               2,
	// }
	// poolJson = append(poolJson, pool1,pool2)
	// pools.CreatePoolsUsingJson(poolJson)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	result := apiModel.SubmitXDRSuccess{
		Status: "Create pool request came",
	}
	json.NewEncoder(w).Encode(result)
	return

}