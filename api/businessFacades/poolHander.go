package businessFacades

import (
	"encoding/json"
	"net/http"
	"strings"

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
	var coinValidationPassed bool
	var newBatchConvertCoinObj model.BatchCoinConvert
	err := json.NewDecoder(r.Body).Decode(&newBatchConvertCoinObj)
	if err != nil {
		logrus.Info(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while decoding the body")
		return
	}

	err0 := validations.ValidateBatchCoinConvert(newBatchConvertCoinObj)
	if err0 != nil {
		logrus.Error(err0)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Request body is invalid")
		return
	} else {
		object := dao.Connection{}
		data, _ := object.GetLiquidityPool(newBatchConvertCoinObj.EquationID, newBatchConvertCoinObj.ProductName, newBatchConvertCoinObj.TenantID, newBatchConvertCoinObj.FormulaType).Then(func(data interface{}) interface{} {
			return data
		}).Await()
		if data == nil {
			logrus.Error("Can not find the pool from DB")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Can not find the pool from DB")
			return
		}
		for _, coin := range data.(model.BuildPoolResponse).CoinMap {
			if coin.CoinName == strings.ToUpper(newBatchConvertCoinObj.MetricCoin.CoinName) {
				newBatchConvertCoinObj.MetricCoin.GeneratedName = coin.GeneratedName
				break
			}
		}
		// Check if the coin name is more than 4 characters
		if len(newBatchConvertCoinObj.MetricCoin.CoinName) == 4 || len(newBatchConvertCoinObj.MetricCoin.GeneratedName) == 12 {
			// loop through the user inputs
			userInputs := newBatchConvertCoinObj.UserInputs
			for i := 0; i < len(userInputs); i++ {
				logrus.Info("Coin name length ", len(userInputs[i].CoinName))
				// check the string length
				if len(userInputs[i].CoinName) != 4 {
					coinValidationPassed = false
					logrus.Error("Coin name character limit should be 4")
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode("Coin name character limit should be 4")
					return
				} else {
					coinValidationPassed = true
				}
			}
		} else {
			logrus.Error("Metric coin name or generated name length error")
			coinValidationPassed = false
		}

		// check if the coin validations are passed
		if coinValidationPassed {
			// check if the formula type is batch or artifact
			if newBatchConvertCoinObj.FormulaType == "BATCH" {
				// execute the rest
				queue := model.SendToQueue{
					Type:        "COINCONVERT",
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
			} else {
				logrus.Error("Invalid formula type")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("Invalid formula type")
				return
			}
		} else {
			logrus.Error("Coin name character limit should be 4")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Coin name character limit should be 4")
			return
		}
	}
}

func ArtifactConvertCoin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var coinValidationPassed bool
	var convertValidate model.ArtifactCoinConvert
	err := json.NewDecoder(r.Body).Decode(&convertValidate)
	if err != nil {
		logrus.Info(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while decoding the body")
		return
	}

	err0 := validations.ValidateArtifactCoinConvert(convertValidate)
	if err0 != nil {
		logrus.Error(err0)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Request body is invalid")
		return
	} else {

		newBatchConvertCoinObj := model.BatchCoinConvert{
			EquationID:      convertValidate.EquationID,
			TenantID:        convertValidate.TenantID,
			ProductName:     convertValidate.ProductName,
			ProductID:       convertValidate.ProductID,
			FormulaType:     convertValidate.FormulaType,
			FormulaTypeID:   convertValidate.FormulaTypeID,
			FormulaTypeName: convertValidate.FormulaTypeID,
			StageId:         convertValidate.StageId,
			MetricCoin:      convertValidate.MetricCoin,
			UserInputs:      convertValidate.UserInputs,
		}
		object := dao.Connection{}
		data, _ := object.GetLiquidityPoolForArtifact(newBatchConvertCoinObj.EquationID, newBatchConvertCoinObj.TenantID, newBatchConvertCoinObj.FormulaType).Then(func(data interface{}) interface{} {
			return data
		}).Await()
		if data == nil {
			logrus.Error("Can not find the pool from DB")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Can not find the pool from DB")
			return
		}
		for _, coin := range data.(model.BuildPoolResponse).CoinMap {
			if coin.CoinName == strings.ToUpper(newBatchConvertCoinObj.MetricCoin.CoinName) {
				newBatchConvertCoinObj.MetricCoin.GeneratedName = coin.GeneratedName
				break
			}
		}
		// Check if the coin name is more than 4 characters
		if len(newBatchConvertCoinObj.MetricCoin.CoinName) == 4 || len(newBatchConvertCoinObj.MetricCoin.GeneratedName) == 12 {
			// loop through the user inputs
			userInputs := newBatchConvertCoinObj.UserInputs
			for i := 0; i < len(userInputs); i++ {
				logrus.Info("Coin name length ", len(userInputs[i].CoinName))
				// check the string length
				if len(userInputs[i].CoinName) != 4 {
					coinValidationPassed = false
					logrus.Error("Coin name character limit should be 4")
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode("Coin name character limit should be 4")
					return
				} else {
					coinValidationPassed = true
				}
			}
		} else {
			logrus.Error("Metric coin name or generated name length error")
			coinValidationPassed = false
		}

		// check if the coin validations are passed
		if coinValidationPassed {
			// check if the formula type is batch or artifact
			if newBatchConvertCoinObj.FormulaType == "ARTIFACT" {
				// execute the rest
				queue := model.SendToQueue{
					Type:        "COINCONVERT",
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
			} else {
				logrus.Error("Invalid formula type")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("Invalid formula type")
				return
			}
		} else {
			logrus.Error("Coin name character limit should be 4")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Coin name character limit should be 4")
			return
		}
	}
}

func CreatePoolForBatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var equationJsonObj model.CreatePool
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
		json.NewEncoder(w).Encode("Request body is invalid")
		return
	} else {
		// Check if the coin name's character equalto 4
		if len(equationJsonObj.MetricCoin.CoinName) != 4 {
			logrus.Error("Coin name character limit should be 4")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Metric Coin name character limit should be 4")
			return
		}
		// check if the equation type is either batch or artifact
		if equationJsonObj.FormulaType == "BATCH" {

			object := dao.Connection{}
			data, _ := object.GetLiquidityPool(equationJsonObj.EquationID, equationJsonObj.ProductName, equationJsonObj.TenantID, equationJsonObj.FormulaType).Then(func(data interface{}) interface{} {
				return data
			}).Await()
			if data != nil {
				logrus.Error("GetLiquidityPool did not empty, Pool already created")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(err.Error())
				return
			}

			// reformate the equation json
			equationJson, coinMap, err := pools.RemoveDivisionAndOperator(equationJsonObj)
			if err != nil {
				logrus.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(err.Error())
				return
			}
			logrus.Info("CoinMap ", coinMap)
			// 	// build the pool creation json
			poolCreationJSON, err := pools.BuildPoolCreationJSON(equationJson)
			if err != nil {
				logrus.Error(poolCreationJSON, err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(err.Error())
				return
			}
			logrus.Info("PoolCreationJSON ", poolCreationJSON)

			equationDetails := model.CreatePool{
				EquationID:  equationJsonObj.EquationID,
				ProductName: equationJsonObj.ProductName,
				ProductID:   equationJsonObj.ProductID,
				TenantID:    equationJsonObj.TenantID,
				FormulaType: equationJson.FormulaType,
				EquatinStringFormate: equationJsonObj.EquatinStringFormate,
				SimpleifedEquation: equationJsonObj.SimpleifedEquation,
			}

			queue := model.SendToQueue{
				Type:              "POOL",
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
			json.NewEncoder(w).Encode("Invalid equation type")
			return
		}
	}
}

func CreatePoolForArtifact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var equationJsonObjValidate model.CreatePoolForArtifact
	// var equationJsonObj model.CreatePool
	errValidate := json.NewDecoder(r.Body).Decode(&equationJsonObjValidate)
	if errValidate != nil {
		logrus.Error(errValidate)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error1 while decoding the body")
		return
	}

	err1 := validations.ValidateCreatePoolForArtifact(equationJsonObjValidate)
	if err1 != nil {
		logrus.Error(err1)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Request body is invalid")
		return
	} else {
		equationJsonObj := model.CreatePool{
			EquationID:           equationJsonObjValidate.EquationID,
			ProductName:          equationJsonObjValidate.ProductName,
			ProductID:            equationJsonObjValidate.ProductID,
			TenantID:             equationJsonObjValidate.TenantID,
			ActivityId:           equationJsonObjValidate.ActivityId,
			FormulaType:          equationJsonObjValidate.FormulaType,
			EquatinStringFormate: equationJsonObjValidate.EquatinStringFormate,
			SimpleifedEquation:   equationJsonObjValidate.SimpleifedEquation,
			MetricCoin:           equationJsonObjValidate.MetricCoin,
			EquationSubPortion:   equationJsonObjValidate.EquationSubPortion,
		}
		// Check if the coin name's character equalto 4
		if len(equationJsonObj.MetricCoin.CoinName) != 4 {
			logrus.Error("Coin name character limit should be 4")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Metric Coin name character limit should be 4")
			return
		}
		// check if the equation type is either batch or artifact
		if equationJsonObj.FormulaType == "ARTIFACT" {

			object := dao.Connection{}
			data, _ := object.GetLiquidityPoolForArtifact(equationJsonObj.EquationID, equationJsonObj.TenantID, equationJsonObj.FormulaType).Then(func(data interface{}) interface{} {
				return data
			}).Await()
			if data != nil {
				logrus.Error("GetLiquidityPool did not empty, Pool already created")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("GetLiquidityPool did not empty, Pool already created")
				return
			}

			// reformate the equation json
			equationJson, coinMap, err := pools.RemoveDivisionAndOperator(equationJsonObj)
			if err != nil {
				logrus.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(err.Error())
				return
			}
			logrus.Info("CoinMap ", coinMap)
			// 	// build the pool creation json
			poolCreationJSON, err := pools.BuildPoolCreationJSON(equationJson)
			if err != nil {
				logrus.Error(poolCreationJSON, err)
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(err.Error())
				return
			}
			logrus.Info("PoolCreationJSON ", poolCreationJSON)

			equationDetails := model.CreatePool{
				EquationID: equationJsonObj.EquationID,
				TenantID:   equationJsonObj.TenantID,
				ProductName: equationJson.ProductName,
				ProductID: equationJson.ProductID,
				FormulaType: equationJson.FormulaType,
			}

			queue := model.SendToQueue{
				Type:              "POOL",
				EqationJson:       equationDetails,
				CoinMap:           coinMap,
				PoolCreationArray: poolCreationJSON,
			}
			// sent data to mgs amq queue
			logrus.Info("Sent..", queue)
			services.SendToQueue(queue)

			log.Println("Pool added to the Queue   (Artifact)")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode("Pool added to the Queue  (Artifact)")
			return

		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Invalid equation type")
			return
		}
	}
}

func CalculateEquationForBatch(w http.ResponseWriter, r *http.Request) {
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
	data, _ := object.GetBatchSpecificAccount(calculateEquationObj.FormulaType, calculateEquationObj.FormulaTypeName,
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
	dbData := data.(model.CoinAccount)
	coinBalance, err := pools.CalculateCoin(dbData)
	if err != nil {
		logrus.Info("Can not find the assert in account")
		w.WriteHeader(http.StatusNoContent)
		result := "Can not find the coin in account"
		json.NewEncoder(w).Encode(result)
		return
	}
	equationResponse = model.EquationResultForBatch{
		EquationID:      dbData.EquationID,
		TenantID:        dbData.TenantID,
		ProductName:     dbData.ProductName,
		FormulaTypeName: dbData.FormulaTypeName,
		FormulaType:     dbData.FormulaType,
		StageId:         dbData.StageID,
		MetricCoin:      dbData.MetricCoin,
		Account:    dbData.CoinAccountPK,
		EquationResult:  coinBalance,
	}

	log.Println("Equation result")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(equationResponse)
	return
}
