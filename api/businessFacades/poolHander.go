package businessFacades

import (
	"encoding/json"
	"fmt"
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

//! CoinName = first 4 chacter form coin name
//! fullCoinName = user inserted coin name

func BatchConvertCoin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var coinConvertBody model.CoinConvertBody

	err := json.NewDecoder(r.Body).Decode(&coinConvertBody)
	if err != nil {
		logrus.Info(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while decoding the body")
		return
	}

	err0 := validations.ValidateBatchCoinConvert(coinConvertBody)
	if err0 != nil {
		logrus.Error(err0)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Request body is invalid")
		return
	} else {
		coinConvertBody.Type = "BATCH"
		object := dao.Connection{}
		data, _ := object.GetLiquidityPoolByProductAndActivity(coinConvertBody.MetricFormulaId,
			coinConvertBody.TenantID, coinConvertBody.Type, coinConvertBody.MetricActivityId,
			coinConvertBody.Event.Details.StageID).Then(func(data interface{}) interface{} {
			return data
		}).Await()

		if data == nil {
			logrus.Error("Can not find the pool from DB")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Can not find the pool from DB")
			return
		}

		for _, coin := range data.(model.BuildPoolResponse).CoinMap {
			if coin.FullCoinName == coinConvertBody.Metric.Name {
				coinConvertBody.ID = coin.ID
				coinConvertBody.Metric.Description = coin.Description
				coinConvertBody.Metric.GeneratedName = coin.GeneratedName
				coinConvertBody.Metric.FullCoinName = coin.FullCoinName
				coinConvertBody.Metric.CoinName = coin.CoinName
				break
			}
		}

		queue := model.SendToQueue{
			Type:        "COINCONVERT",
			CoinConvert: coinConvertBody,
		}

		logrus.Info("Sent.. to rabbitmq")
		// sent data to mgs amq queue
		services.SendToQueue(queue)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Path payment added to queue")
		return
	}
}

func ArtifactConvertCoin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var coinConvertBody model.CoinConvertBody

	err := json.NewDecoder(r.Body).Decode(&coinConvertBody)
	if err != nil {
		logrus.Info(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while decoding the body")
		return
	}

	err0 := validations.ValidateBatchCoinConvert(coinConvertBody)
	if err0 != nil {
		logrus.Error(err0)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Request body is invalid")
		return
	} else {
		coinConvertBody.Type = "ARTIFACT"
		object := dao.Connection{}
		data, _ := object.GetLiquidityPoolByProductAndActivity(coinConvertBody.MetricFormulaId,
			coinConvertBody.TenantID, coinConvertBody.Type, coinConvertBody.MetricActivityId,
			coinConvertBody.Event.Details.StageID).Then(func(data interface{}) interface{} {
			return data
		}).Await()

		if data == nil {
			logrus.Error("Can not find the pool from DB")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Can not find the pool from DB")
			return
		}

		for _, coin := range data.(model.BuildPoolResponse).CoinMap {
			if coin.FullCoinName == coinConvertBody.Metric.Name {
				coinConvertBody.ID = coin.ID
				coinConvertBody.Metric.Description = coin.Description
				coinConvertBody.Metric.GeneratedName = coin.GeneratedName
				coinConvertBody.Metric.FullCoinName = coin.FullCoinName
				coinConvertBody.Metric.CoinName = coin.CoinName
				break
			}
		}

		queue := model.SendToQueue{
			Type:        "COINCONVERT",
			CoinConvert: coinConvertBody,
		}

		logrus.Info("Sent.. to rabbitmq")
		// sent data to mgs amq queue
		services.SendToQueue(queue)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Path payment added to queue")
		return
	}
}

func CreatePoolForBatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var equationJsonBody model.CreatePoolBody

	err := json.NewDecoder(r.Body).Decode(&equationJsonBody)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while decoding the body")
		return
	}

	err1 := validations.ValidateCreatePool(equationJsonBody)
	if err1 != nil {
		logrus.Error(err1)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Request body is invalid")
		return
	} else {

		var equationSubPortions []model.EquationSubPortion
		// build the 4 charater coin name by using equation jason
		for i := 0; i < len(equationJsonBody.Subsets); i++ {
			var fieldAndCoins []model.FieldAndCoin

			for j := 0; j < len(equationJsonBody.Subsets[i].SubSet); j++ {
				// if subset elemet's type eqal to "DATA" take first 4 charter from coin name
				if equationJsonBody.Subsets[i].SubSet[j].Type == "DATA" {
					fieldAndCoin := model.FieldAndCoin{
						ID:            equationJsonBody.Subsets[i].SubSet[j].ID,
						CoinName:      equationJsonBody.Subsets[i].SubSet[j].CoinName[0:4],
						GeneratedName: "",
						FullCoinName:  equationJsonBody.Subsets[i].SubSet[j].CoinName,
						Description:   equationJsonBody.Subsets[i].SubSet[j].Description,
						VariableType:  equationJsonBody.Subsets[i].SubSet[j].Type,
						Value:         "",
					}
					fieldAndCoins = append(fieldAndCoins, fieldAndCoin)
				} else if equationJsonBody.Subsets[i].SubSet[j].Type == "CONSTANT" {
					// if subset elemet's type eqal to "CONSTANT" take first 4 charter from coin name
					if equationJsonBody.Subsets[i].SubSet[j].CoinName != "" {
						fieldAndCoin := model.FieldAndCoin{
							ID:            equationJsonBody.Subsets[i].SubSet[j].ID,
							CoinName:      equationJsonBody.Subsets[i].SubSet[j].CoinName[0:4],
							GeneratedName: "",
							FullCoinName:  equationJsonBody.Subsets[i].SubSet[j].CoinName,
							Description:   equationJsonBody.Subsets[i].SubSet[j].Description,
							VariableType:  equationJsonBody.Subsets[i].SubSet[j].Type,
							Value:         equationJsonBody.Subsets[i].SubSet[j].Value,
						}
						fieldAndCoins = append(fieldAndCoins, fieldAndCoin)
					} else {
						// if subset elemet's type eqal to "CONSTANT" and did not send the coin name uild the coin name by using the CreateCoinnameUsingValue
						fieldAndCoin := model.FieldAndCoin{
							ID:            equationJsonBody.Subsets[i].SubSet[j].ID,
							CoinName:      pools.CreateCoinnameUsingValue(equationJsonBody.Subsets[i].SubSet[j].Value),
							GeneratedName: "",
							FullCoinName:  equationJsonBody.Subsets[i].SubSet[j].Value,
							Description:   equationJsonBody.Subsets[i].SubSet[j].Description,
							VariableType:  equationJsonBody.Subsets[i].SubSet[j].Type,
							Value:         equationJsonBody.Subsets[i].SubSet[j].Value,
						}
						fieldAndCoins = append(fieldAndCoins, fieldAndCoin)
					}
				} else {
					fieldAndCoin := model.FieldAndCoin{
						ID:            equationJsonBody.Subsets[i].SubSet[j].ID,
						CoinName:      "",
						GeneratedName: "",
						FullCoinName:  "",
						Description:   "",
						VariableType:  equationJsonBody.Subsets[i].SubSet[j].Type,
						Value:         equationJsonBody.Subsets[i].SubSet[j].Value,
					}
					fieldAndCoins = append(fieldAndCoins, fieldAndCoin)
				}
			}
			equationSubPortion := model.EquationSubPortion{
				Subportion:   equationJsonBody.Subsets[i].SubsetAsString,
				FieldAndCoin: fieldAndCoins,
			}
			equationSubPortions = append(equationSubPortions, equationSubPortion)
		}

		var activity []model.Activity
		activity = append(activity, model.Activity{
			ID:              equationJsonBody.Activity.ID,
			Name:            equationJsonBody.Activity.Name,
			ProductName:     equationJsonBody.ProductName,
			TracifiedItemId: equationJsonBody.Activity.TracifiedItemId,
			StageId:         equationJsonBody.Activity.StageId,
		})

		var products []model.Product
		products = append(products, model.Product{
			ProductName: equationJsonBody.ProductName,
			ProductID:   equationJsonBody.ProductID,
		})

		equationJsonObj := model.CreatePool{
			EquationID:           equationJsonBody.ID,
			Products:             products,
			TenantID:             equationJsonBody.TenantID,
			Activity:             activity,
			FormulaType:          "BATCH",
			EquatinStringFormate: equationJsonBody.FormulaAsString,
			SimpleifedEquation:   equationJsonBody.FormulaAsString,
			MetricCoin: model.MetricCoin{
				ID:           equationJsonBody.Metric.ID,
				Name:         equationJsonBody.Metric.Name,
				CoinName:     equationJsonBody.Metric.Name[0:4],
				FullCoinName: equationJsonBody.Metric.Name,
				Description:  equationJsonBody.Metric.Description,
			},
			EquationSubPortion: equationSubPortions,
			CreatedAt:          equationJsonBody.CreatedAt,
			UpdatedAt:          equationJsonBody.UpdatedAt,
			Active:             equationJsonBody.Active,
		}

		// Check if the coin name's character equalto 4
		if len(equationJsonObj.MetricCoin.CoinName) != 4 {
			logrus.Error("Coin name character limit should be 4")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Metric Coin name character limit should be 4")
			return
		}

		object := dao.Connection{}
		data, _ := object.GetLiquidityPool(equationJsonObj.EquationID, equationJsonObj.TenantID, equationJsonObj.FormulaType).Then(func(data interface{}) interface{} {
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
		equationDetails := model.CreatePool{
			EquationID:           equationJsonObj.EquationID,
			Products:             equationJson.Products,
			TenantID:             equationJsonObj.TenantID,
			Activity:             equationJson.Activity,
			FormulaType:          equationJson.FormulaType,
			EquatinStringFormate: equationJson.EquatinStringFormate,
			SimpleifedEquation:   equationJson.EquatinStringFormate,
			MetricCoin:           equationJson.MetricCoin,
			EquationSubPortion:   equationSubPortions,
			CreatedAt:            equationJson.CreatedAt,
			UpdatedAt:            equationJson.UpdatedAt,
			Active:               equationJson.Active,
		}
		queue := model.SendToQueue{
			Type:              "POOL",
			EqationJson:       equationDetails,
			CoinMap:           coinMap,
			PoolCreationArray: poolCreationJSON,
		}
		// sent data to mgs amq queue
		out, err := json.Marshal(queue)
		if err != nil {
			log.Error(err)
		}

		fmt.Println("SENT", string(out))
		services.SendToQueue(queue)

		log.Println("Pool added to the Queue")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Pool added to the Queue")
		return
	}
}

func CreatePoolForArtifact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var equationJsonBody model.CreatePoolBody
	// var equationJsonObj model.CreatePool
	errValidate := json.NewDecoder(r.Body).Decode(&equationJsonBody)
	if errValidate != nil {
		logrus.Error(errValidate)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error1 while decoding the body")
		return
	}

	err1 := validations.ValidateCreatePool(equationJsonBody)
	if err1 != nil {
		logrus.Error(err1)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Request body is invalid")
		return
	} else {
		var equationSubPortions []model.EquationSubPortion

		for i := 0; i < len(equationJsonBody.Subsets); i++ {
			var fieldAndCoins []model.FieldAndCoin

			for j := 0; j < len(equationJsonBody.Subsets[i].SubSet); j++ {
				if equationJsonBody.Subsets[i].SubSet[j].Type == "DATA" {
					fieldAndCoin := model.FieldAndCoin{
						ID:            equationJsonBody.Subsets[i].SubSet[j].ID,
						CoinName:      equationJsonBody.Subsets[i].SubSet[j].CoinName[0:4],
						GeneratedName: "",
						FullCoinName:  equationJsonBody.Subsets[i].SubSet[j].CoinName,
						Description:   equationJsonBody.Subsets[i].SubSet[j].Description,
						VariableType:  equationJsonBody.Subsets[i].SubSet[j].Type,
						Value:         "",
					}
					fieldAndCoins = append(fieldAndCoins, fieldAndCoin)
				} else if equationJsonBody.Subsets[i].SubSet[j].Type == "CONSTANT" {
					if equationJsonBody.Subsets[i].SubSet[j].CoinName != "" {
						fieldAndCoin := model.FieldAndCoin{
							ID:            equationJsonBody.Subsets[i].SubSet[j].ID,
							CoinName:      equationJsonBody.Subsets[i].SubSet[j].CoinName[0:4],
							GeneratedName: "",
							FullCoinName:  equationJsonBody.Subsets[i].SubSet[j].CoinName,
							Description:   equationJsonBody.Subsets[i].SubSet[j].Description,
							VariableType:  equationJsonBody.Subsets[i].SubSet[j].Type,
							Value:         equationJsonBody.Subsets[i].SubSet[j].Value,
						}
						fieldAndCoins = append(fieldAndCoins, fieldAndCoin)
					} else {
						fieldAndCoin := model.FieldAndCoin{
							ID:            equationJsonBody.Subsets[i].SubSet[j].ID,
							CoinName:      pools.CreateCoinnameUsingValue(equationJsonBody.Subsets[i].SubSet[j].Value),
							GeneratedName: "",
							FullCoinName:  equationJsonBody.Subsets[i].SubSet[j].Value,
							Description:   equationJsonBody.Subsets[i].SubSet[j].Description,
							VariableType:  equationJsonBody.Subsets[i].SubSet[j].Type,
							Value:         equationJsonBody.Subsets[i].SubSet[j].Value,
						}
						fieldAndCoins = append(fieldAndCoins, fieldAndCoin)
					}
				} else {
					fieldAndCoin := model.FieldAndCoin{
						ID:            equationJsonBody.Subsets[i].SubSet[j].ID,
						CoinName:      "",
						GeneratedName: "",
						FullCoinName:  "",
						Description:   "",
						VariableType:  equationJsonBody.Subsets[i].SubSet[j].Type,
						Value:         equationJsonBody.Subsets[i].SubSet[j].Value,
					}
					fieldAndCoins = append(fieldAndCoins, fieldAndCoin)
				}
			}
			equationSubPortion := model.EquationSubPortion{
				Subportion:   equationJsonBody.Subsets[i].SubsetAsString,
				FieldAndCoin: fieldAndCoins,
			}
			equationSubPortions = append(equationSubPortions, equationSubPortion)
		}

		var activity []model.Activity
		activity = append(activity, model.Activity{
			ID:              equationJsonBody.Activity.ID,
			Name:            equationJsonBody.Activity.Name,
			ProductName:     equationJsonBody.ProductName,
			TracifiedItemId: equationJsonBody.ProductID,
		})

		var products []model.Product
		products = append(products, model.Product{
			ProductName: equationJsonBody.ProductName,
			ProductID:   equationJsonBody.ProductID,
		})

		equationJsonObj := model.CreatePool{
			EquationID:           equationJsonBody.ID,
			Products:             products,
			TenantID:             equationJsonBody.TenantID,
			Activity:             activity,
			FormulaType:          "BATCH",
			EquatinStringFormate: equationJsonBody.FormulaAsString,
			SimpleifedEquation:   equationJsonBody.FormulaAsString,
			MetricCoin: model.MetricCoin{
				ID:           equationJsonBody.Metric.ID,
				CoinName:     equationJsonBody.Metric.Name[0:4],
				FullCoinName: equationJsonBody.Metric.Name,
				Description:  equationJsonBody.Metric.Description,
			},
			EquationSubPortion: equationSubPortions,
			CreatedAt:          equationJsonBody.CreatedAt,
			UpdatedAt:          equationJsonBody.UpdatedAt,
			Active:             equationJsonBody.Active,
		}
		// Check if the coin name's character equalto 4
		if len(equationJsonObj.MetricCoin.CoinName) != 4 {
			logrus.Error("Coin name character limit should be 4")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Metric Coin name character limit should be 4")
			return
		}
		// check if the equation type is either batch or artifact
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
			EquationID:           equationJsonObj.EquationID,
			Products:             equationJson.Products,
			TenantID:             equationJsonObj.TenantID,
			Activity:             equationJson.Activity,
			FormulaType:          equationJson.FormulaType,
			EquatinStringFormate: equationJson.EquatinStringFormate,
			SimpleifedEquation:   equationJson.EquatinStringFormate,
			MetricCoin:           equationJson.MetricCoin,
			EquationSubPortion:   equationSubPortions,
			CreatedAt:            equationJson.CreatedAt,
			UpdatedAt:            equationJson.UpdatedAt,
			Active:               equationJson.Active,
		}

		queue := model.SendToQueue{
			Type:              "POOL",
			EqationJson:       equationDetails,
			CoinMap:           coinMap,
			PoolCreationArray: poolCreationJSON,
		}
		// sent data to mgs amq queue
		services.SendToQueue(queue)
		// sent data to mgs amq queue
		out, err := json.Marshal(queue)
		if err != nil {
			log.Error(err)
		}

		fmt.Println("SENT", string(out))

		log.Println("Pool added to the Queue   (Artifact)")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Pool added to the Queue  (Artifact)")
		return
	}
}

func CalculateEquationForBatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var calculateEquationObj model.CalculateEquation
	var equationResponse model.EquationResultForBatch

	err := json.NewDecoder(r.Body).Decode(&calculateEquationObj)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while decoding the body")
		return
	}
	object := dao.Connection{}
	var data interface{}
	if calculateEquationObj.Type == "BATCH" {
		data, _ = object.GetBatchSpecificAccount(calculateEquationObj.Type, calculateEquationObj.BatchID,
			calculateEquationObj.MetricFormulaId, calculateEquationObj.TracifiedItemId,
			calculateEquationObj.TenantID, calculateEquationObj.StageID, calculateEquationObj.MetricActivityId).Then(func(data interface{}) interface{} {
			return data
		}).Await()
	} else if calculateEquationObj.Type == "ARTIFACT" {
		data, _ = object.GetBatchSpecificAccount(calculateEquationObj.Type, calculateEquationObj.BatchID,
			calculateEquationObj.MetricFormulaId, calculateEquationObj.TracifiedItemId,
			calculateEquationObj.TenantID, calculateEquationObj.StageID, calculateEquationObj.MetricActivityId).Then(func(data interface{}) interface{} {
			return data
		}).Await()
	} else {
		logrus.Info("Invalied equation Type")
		w.WriteHeader(http.StatusBadRequest)
		result := "Invalied equation Type"
		json.NewEncoder(w).Encode(result)
		return
	}
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
		MetricFormulaId:  dbData.MetricFormulaId,
		TenantID:         dbData.TenantID,
		MetricActivityId: dbData.MetricActivivtyId,
		Type:             dbData.Type,
		BatchID:          dbData.Event.Details.BatchID,
		ArtifactID:       dbData.Event.Details.ArtifactID,
		StageID:          dbData.Event.Details.StageID,
		TracifiedItemId:  dbData.Event.Details.TracifiedItemId,
		Metric:           dbData.Metric,
		Account:          dbData.CoinAccountPK,
		EquationResult:   coinBalance,
	}

	log.Println("Equation result")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(equationResponse)
	return
}

// get the path payment details
func GetPathPaymentDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	tenantid, error := r.URL.Query()["tenantid"]
	if !error || len(tenantid[0]) < 1 {
		logrus.Error("Url Parameter 'tenantid' is missing")
		w.WriteHeader(http.StatusNotFound)
		response := apiModel.SubmitXDRSuccess{
			Status: "Url Parameter 'tenantid' is missing",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	equationid, error := r.URL.Query()["equationid"]
	if !error || len(equationid[0]) < 1 {
		logrus.Error("Url Parameter 'equationid' is missing")
		w.WriteHeader(http.StatusNotFound)
		response := apiModel.SubmitXDRSuccess{
			Status: "Url Parameter 'equationid' is missing",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// getting query parameters
	formulatype, error := r.URL.Query()["formulatype"]
	if !error || len(formulatype[0]) < 1 {
		logrus.Error("Url Parameter 'formulatype' is missing")
		w.WriteHeader(http.StatusNotFound)
		response := apiModel.SubmitXDRSuccess{
			Status: "Url Parameter 'formulatype' is missing",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	productname, error := r.URL.Query()["productname"]
	if !error {
		logrus.Error("Url Parameter 'productname' is having an error")
		w.WriteHeader(http.StatusNotFound)
		response := apiModel.SubmitXDRSuccess{
			Status: "Url Parameter 'productname' is having an error",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	object := dao.Connection{}

	qData, err := object.GetCoinConversionDetails(formulatype[0], equationid[0], productname[0], tenantid[0]).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err != nil {
		logrus.Error("Unable to connect gateway datastore")
		w.WriteHeader(http.StatusNotFound)
		response := apiModel.SubmitXDRSuccess{
			Status: "Unable to connect gateway datastore",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	if qData == nil {
		logrus.Error("Coin convert details are not found in gateway datastore")
		w.WriteHeader(http.StatusNoContent)
		response := apiModel.SubmitXDRSuccess{
			Status: "Coin convert details are not found in gateway datastore",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(qData)
}

func UpdateEquation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	object := dao.Connection{}
	var addStageAndProduct model.AddStageAndProduct

	err := json.NewDecoder(r.Body).Decode(&addStageAndProduct)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while decoding the body")
		return
	}
	data, _ := object.GetLiquidityPool(addStageAndProduct.EquationID, addStageAndProduct.TenantID, addStageAndProduct.Type).Then(func(data interface{}) interface{} {
		return data
	}).Await()

	if data == nil {
		logrus.Info("Can not find the Batch account")
		w.WriteHeader(http.StatusNoContent)
		result := "Can not find the Batch account"
		json.NewEncoder(w).Encode(result)
		return
	}

	dbData := data.(model.BuildPoolResponse)
	products := dbData.Products
	activity := dbData.Activity

	productExist := false
	activityExist := false
	for _, element := range dbData.Products {
		if addStageAndProduct.ProductID != "" && element.ProductID == addStageAndProduct.ProductID {
			productExist = true
		}
	}
	if !productExist || len(dbData.Products) == 0 {
		products = append(dbData.Products, model.Product{
			ProductName: addStageAndProduct.ProductName,
			ProductID:   addStageAndProduct.ProductID,
		})
	}

	for i, element := range dbData.Activity {
		if addStageAndProduct.Activity.ID != "" && element.ID == addStageAndProduct.Activity.ID {
			activityExist = true
			if dbData.Activity[i].TracifiedItemId == "" {
				dbData.Activity[i].TracifiedItemId = addStageAndProduct.ProductID
				dbData.Activity[i].ProductName = addStageAndProduct.ProductName
			} else if addStageAndProduct.ProductID != "" && element.TracifiedItemId != addStageAndProduct.ProductID {
				activity = append(activity, model.Activity{
					ID:              addStageAndProduct.Activity.ID,
					Name:            addStageAndProduct.Activity.Name,
					ProductName:     addStageAndProduct.ProductName,
					TracifiedItemId: addStageAndProduct.ProductID,
				})
			}
		}
	}

	if !activityExist || len(dbData.Activity) == 0 {
		activity = append(activity, model.Activity{
			ID:              addStageAndProduct.Activity.ID,
			Name:            addStageAndProduct.Activity.Name,
			ProductName:     addStageAndProduct.ProductName,
			TracifiedItemId: addStageAndProduct.ProductID,
		})
	}
	fmt.Println(products)
	fmt.Println(activity)
	dbData.Activity = activity
	dbData.Products = products
	object.UpdateLiquidityPool(dbData)

	if err != nil {
		logrus.Info("Can not find the Equation")
		w.WriteHeader(http.StatusNoContent)
		result := "Can not find the Equation"
		json.NewEncoder(w).Encode(result)
		return
	}

	log.Println("Equation result")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dbData)
	return
}
