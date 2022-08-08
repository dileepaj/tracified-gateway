package pools

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

var (
	depositorPK = commons.GoDotEnvVariable("DEPOSITEPK")
	depositorSK = commons.GoDotEnvVariable("DEPOSITERSeed")
)

// CreatePoolsUsingJson , loop the pool creation-Josn and call poolcreation method
func CreatePoolsUsingJson(pools []model.BuildPool) ([]model.BuildPool, error, bool) {
	object := dao.Connection{}
	var createdPools []model.BuildPool
	var isCreated bool
	for i := 0; i < len(pools); i++ {
		data, _ := object.GetCreatedPool(pools[i].Coin1, pools[i].Coin2).Then(func(data interface{}) interface{} {
			return data
		}).Await()
		if data == nil {
			pool, err, created := CreatePool(pools[i])
			logrus.Error("Pool does not added to Creat Pool")
			logrus.Info("Pool added to Creat Pool")
			if err == nil {
				createdPools = append(createdPools, pool)
				isCreated = created
			} else {
				logrus.Error(err)
			}
		} else {
			pools[i].PoolDepositeHash = data.(model.BuildPool).PoolDepositeHash
			pools[i].PoolId = data.(model.BuildPool).PoolId
			err1 := object.InsertCreatedPoool(pools[i])
			createdPools = append(createdPools, pools[i])
			isCreated = true
			if err1 != nil {
				logrus.Error("Pool did not add to DB ", err1)
			}
		}
		logrus.Error("pool already deposited " + " Coin1 " + pools[i].Coin1 + " Coin2 " + pools[i].Coin2)
	}
	return createdPools, nil, isCreated
}

// CreatePool carate the pool by handling all blackchain calls
func CreatePool(buildPool model.BuildPool) (model.BuildPool, error, bool) {
	object := dao.Connection{}
	var depostHash2 string
	logrus.Info("Coin 1 ", buildPool.Coin1)
	logrus.Info("Coin 2 ", buildPool.Coin2)
	poolId, err := GeneratePoolId(buildPool.Coin1, buildPool.Coin2)
	if !err {
		return model.BuildPool{}, errors.New("Can not create poold Id"), false
	}
	poolIdString := hex.EncodeToString(poolId[:])

	// calling the stellar endpoint to see if the pool exists
	url1 := commons.GetHorizonClient().HorizonURL + "liquidity_pools/" + poolIdString
	result, err0 := http.Get(url1)
	if err0 != nil {
		logrus.Error("Unable to reach Stellar network", url1)
		return model.BuildPool{}, err0, false
	}
	if result.StatusCode == 404 {
		logrus.Info("Pool " + poolIdString + " is not created")

		// deposite to the pool
		CreateCoin(buildPool.Coin1, depositorPK, depositorSK)
		CreateCoin(buildPool.Coin2, depositorPK, depositorSK)
		IssueCoin(buildPool.Coin1, depositorPK, buildPool.DepositeAmountCoin1)
		IssueCoin(buildPool.Coin2, depositorPK, buildPool.DepositeAmountCoin2)

		// EstablishPoolTrustline , this method create trust line(change trusst opration) with pool
		_, err1 := EstablishPoolTrustline(buildPool.Coin1, buildPool.Coin2, depositorPK, depositorSK)
		if err1 != nil {
			return model.BuildPool{}, err1, false
		}
		// Checked the Asset's names are in lexicographic order if not change the DepositeAmountCoin1 and asset DepositeAmountCoin2
		if buildPool.Coin1 < buildPool.Coin2 {
			depostHash, err := DepositeToPool(poolId, depositorPK, depositorSK, buildPool.DepositeAmountCoin1, buildPool.DepositeAmountCoin2)
			if err != nil {
				return model.BuildPool{}, err, false
			}
			depostHash2 = depostHash
		} else {
			depostHash, err := DepositeToPool(poolId, depositorPK, depositorSK, buildPool.DepositeAmountCoin2, buildPool.DepositeAmountCoin1)
			if err != nil {
				return model.BuildPool{}, err, false
			}
			depostHash2 = depostHash
		}
		buildPool.PoolDepositeHash = depostHash2
		buildPool.PoolId = poolIdString
		logrus.Info(depostHash2, err)
		if buildPool.PoolDepositeHash != "" {
			err1 := object.InsertCreatedPoool(buildPool)
			if err1 != nil {
				logrus.Error("Pool did not add to DB ", err1)
			}
		}
		return buildPool, nil, true
	}
	if result.StatusCode == 200 {
		poolInfo, err := ioutil.ReadAll(result.Body)
		if err != nil {
			logrus.Error(err)
		}

		var raw map[string]interface{}
		var raw1 []interface{}

		json.Unmarshal(poolInfo, &raw)

		out, _ := json.Marshal(raw["reserves"])
		json.Unmarshal(out, &raw1)

		asset1 := raw1[0].(map[string]interface{})
		asset2 := raw1[1].(map[string]interface{})

		amount1 := fmt.Sprintf("%v", asset1["amount"])
		amount2 := fmt.Sprintf("%v", asset2["amount"])

		// checking if the pool has been deposited
		if amount1 == "0.0000000" && amount2 == "0.0000000" {
			logrus.Info("Pool " + poolIdString + " is not deposited")

			// deposite to the pool
			CreateCoin(buildPool.Coin1, depositorPK, depositorSK)
			CreateCoin(buildPool.Coin2, depositorPK, depositorSK)
			IssueCoin(buildPool.Coin1, depositorPK, buildPool.DepositeAmountCoin1)
			IssueCoin(buildPool.Coin2, depositorPK, buildPool.DepositeAmountCoin2)

			// EstablishPoolTrustline , this method create trust line(change trusst opration) with pool
			_, err1 := EstablishPoolTrustline(buildPool.Coin1, buildPool.Coin2, depositorPK, depositorSK)
			if err1 != nil {
				return model.BuildPool{}, err1, false
			}
			// Checked the Asset's names are in lexicographic order if not change the DepositeAmountCoin1 and asset DepositeAmountCoin2
			if buildPool.Coin1 < buildPool.Coin2 {
				depostHash, err := DepositeToPool(poolId, depositorPK, depositorSK, buildPool.DepositeAmountCoin1, buildPool.DepositeAmountCoin2)
				if err != nil {
					return model.BuildPool{}, err, false
				}
				depostHash2 = depostHash
			} else {
				depostHash, err := DepositeToPool(poolId, depositorPK, depositorSK, buildPool.DepositeAmountCoin2, buildPool.DepositeAmountCoin1)
				if err != nil {
					return model.BuildPool{}, err, false
				}
				depostHash2 = depostHash
			}
			buildPool.PoolDepositeHash = depostHash2
			buildPool.PoolId = poolIdString
			logrus.Info(depostHash2, err)
			return buildPool, nil, true
		} else {
			logrus.Info("Pool " + poolIdString + " is deposited")
			return model.BuildPool{}, nil, false
		}
	}
	return model.BuildPool{}, nil, false
}

func PoolCreateHandle(equationJsonObj model.CreatePool, coinMap []model.CoinMap, poolCreationJSON []model.BuildPool) (string, error) {
	// logrus.Info("PoolCreationJSON ", poolCreationJSON)
	object := dao.Connection{}
	data, _ := object.GetLiquidityPool(equationJsonObj.EquationID, equationJsonObj.TenantID, equationJsonObj.FormulaType).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if data != nil {
		logrus.Error("GetLiquidityPool did not empty, Pool already created")
		return "", errors.New("Pool already created")
	}
	// create the pools
	cratedPools, err, isPoolCreated := CreatePoolsUsingJson(poolCreationJSON)
	if err != nil {
		logrus.Error("CreatePoolsUsingJson ", cratedPools, err)
		return "", err
	}
	// logrus.Info("Data", equationJsonObj.Activity, equationJsonObj.EquatinStringFormate, equationJsonObj.SimpleifedEquation)
	response := model.BuildPoolResponse{
		EquationId:           equationJsonObj.EquationID,
		Products:             equationJsonObj.Products,
		TenantId:             equationJsonObj.TenantID,
		Activity:             equationJsonObj.Activity,
		MetricCoin:           equationJsonObj.MetricCoin,
		EquationStringFormat: equationJsonObj.EquatinStringFormate,
		SimplifiedEquation:   equationJsonObj.SimpleifedEquation,
		FormulaType:          equationJsonObj.FormulaType,
		CoinMap:              coinMap,
		BuildPools:           cratedPools,
		EquationSubPortion:   equationJsonObj.EquationSubPortion,
		CreatedAt:            equationJsonObj.CreatedAt,
		UpdatedAt:            equationJsonObj.UpdatedAt,
		Active:               equationJsonObj.Active,
	}
	// check if the pool is created
	if isPoolCreated {
		logrus.Println("New pools are created")
		// insert the pool to the DB
		object := dao.Connection{}
		err1 := object.InsertLiquidityPool(response)
		if err1 != nil {
			logrus.Println("Error when inserting pool to DB " + err.Error())
			return "", errors.New("Error when inserting pool to DB " + err.Error())
		} else {
			logrus.Println("Pool added to the DB")
			out, err := json.Marshal(response)
			if err != nil {
				logrus.Error("Pool added to the DB JSON Marshal Error")
			}
			return string(out), nil
		}
	} else {
		logrus.Println("Pool is already created and deposited")
		return "", errors.New("Pool is already created and deposited")
	}
}
