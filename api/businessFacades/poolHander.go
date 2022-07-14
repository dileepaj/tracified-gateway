package businessFacades

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/pools"
	"github.com/dileepaj/tracified-gateway/services"
	"github.com/dileepaj/tracified-gateway/validations"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func BatchConvertCoin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newBatchConvertCoinObj model.BatchCoinConvert
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

	err0 := validations.ValidateBatchCoinConvert(newBatchConvertCoinObj)
	if err0 != nil {
		logrus.Error(err0)
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Request body is invalid",
		}
		json.NewEncoder(w).Encode(result)
		return
	} else {

		queue := model.SendToQueue{
			Type:        "CionConvert",
			CoinConvert: newBatchConvertCoinObj,
		}

		logrus.Info("Sent..", queue)
		// sent data to mgs amq queue
		services.SendToQueue(queue)

		log.Println("Coin conversion details added to the DB")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Path payment added to queue")
		return
	}
}

func CreatePool(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var equationJsonObj model.CreatePool
	var coinValidationPassed bool
	err := json.NewDecoder(r.Body).Decode(&equationJsonObj)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while decoding the body")
		return
	}

	err1 := validations.ValidateCreatePool(equationJsonObj)
	if err1 != nil {
		logrus.Error(err1)
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Request body is invalid",
		}
		json.NewEncoder(w).Encode(result)
		return
	} else {
		//Check if the coin name is more than 4 characters
		logrus.Info("Coin name length ", len(equationJsonObj.MetricCoin.CoinName))
		if len(equationJsonObj.MetricCoin.CoinName) <= 4 {
			//loop through the equation sub portions to check the coin name
			subportion := equationJsonObj.EquationSubPortion
			for i := 0; i < len(subportion); i++ {
				for j := 0; j < len(subportion[i].FieldAndCoin); j++ {
					logrus.Info("Coin name length ", len(subportion[i].FieldAndCoin[j].CoinName))
					//check the coin name length
					if len(subportion[i].FieldAndCoin[j].CoinName) > 4 {
						coinValidationPassed = false
						logrus.Error("Coin name exceeded character limit of 4")
						w.WriteHeader(http.StatusBadRequest)
						result := apiModel.SubmitXDRSuccess{
							Status: "Coin name exceeded character limit of 4",
						}
						json.NewEncoder(w).Encode(result)
						return
					} else {
						coinValidationPassed = true
					}
				}
			}
		} else {
			coinValidationPassed = false
		}
		//chceck if the coin validations are passed
		if coinValidationPassed {
			logrus.Error("Coin name validations passed")

			//check if the equation type is either batch or artifact
			if equationJsonObj.FormulaType == "Batch" || equationJsonObj.FormulaType == "Artifact" {

				object := dao.Connection{}
				data, _ := object.GetLiquidityPool(equationJsonObj.EquationID, equationJsonObj.ProductName, equationJsonObj.TenantID).Then(func(data interface{}) interface{} {
					return data
				}).Await()
				if data != nil {
					logrus.Error("GetLiquidityPool did not empty, Pool already created")
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(err)
					return
				}

				// reformate the equation json
				equationJson, coinMap, err := pools.RemoveDivisionAndOperator(equationJsonObj)
				if err != nil {
					logrus.Error(err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(err)
					return
				}
				logrus.Info("CoinMap ", coinMap)
				// 	// build the pool creation json
				poolCreationJSON, err := pools.BuildPoolCreationJSON(equationJson)
				if err != nil {
					logrus.Error(poolCreationJSON, err)
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(err)
					return
				}
				logrus.Info("PoolCreationJSON ", poolCreationJSON)

				equationDetails := model.CreatePool{
					EquationID:  equationJsonObj.EquationID,
					ProductName: equationJsonObj.ProductName,
					ProductID:   equationJsonObj.ProductID,
					TenantID:    equationJsonObj.TenantID,
				}

				queue := model.SendToQueue{
					Type:              "Pool",
					EqationJson:       equationDetails,
					CoinMap:           coinMap,
					PoolCreationArray: poolCreationJSON,
				}
				// sent data to mgs amq queue
				logrus.Info("Sent..", queue)
				services.SendToQueue(queue)

				log.Println("Pool added to the Queue")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode("Pool added to the Queue")
				return

			} else {
				w.WriteHeader(http.StatusBadRequest)
				result := apiModel.SubmitXDRSuccess{
					Status: "Invalid equation type",
				}
				json.NewEncoder(w).Encode(result)
				return
			}

		} else {
			logrus.Error("Coin name exceeded character limit of 4")
			w.WriteHeader(http.StatusBadRequest)
			result := apiModel.SubmitXDRSuccess{
				Status: "Coin name exceeded character limit of 4",
			}
			json.NewEncoder(w).Encode(result)
			return
		}
	}
}

func CacluateEquationForBatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var calculateEquationObj model.CalculateEquationForBatch
	var equationResponse model.EquationResultForBatch

	err := json.NewDecoder(r.Body).Decode(&calculateEquationObj)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while decoding the body")
		return
	}
	object := dao.Connection{}
	data, _ := object.GetBatchSpecificAccount(calculateEquationObj.BatchID,
		calculateEquationObj.EquationID, calculateEquationObj.ProductName,
		calculateEquationObj.TenantID).Then(func(data interface{}) interface{} {
		return data
	}).Await()

	if data == nil {
		logrus.Info("Can not find the Batch account")
		w.WriteHeader(http.StatusNoContent)
		result := "Can not find the Batch account"
		json.NewEncoder(w).Encode(result)
		return
	}
	dbData := data.(model.BatchAccount)
	coinBalance, err := pools.CalculateCoin(dbData)
	if err != nil {
		logrus.Info("Can not find the assert in account")
		w.WriteHeader(http.StatusNoContent)
		result := "Can not find the coin in account"
		json.NewEncoder(w).Encode(result)
		return
	}
	equationResponse = model.EquationResultForBatch{
		TenantID:       dbData.TenantID,
		ProductName:    dbData.ProductName,
		ProductID:      dbData.ProductID,
		BatchID:        dbData.BatchID,
		BatchName:      dbData.BatchName,
		StageId:        dbData.StageID,
		EquationID:     dbData.EquationID,
		MetrixType:     calculateEquationObj.MetrixType,
		MetricCoin:     dbData.MetricCoin,
		BatchAccount:   dbData.BatchAccountPK,
		EquationResult: coinBalance,
	}

	log.Println("Equation result")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(equationResponse)
	return
}
