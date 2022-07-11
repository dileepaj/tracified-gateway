package pools

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/support/log"
)

var (
	depositorPK = "GBBW4XYGCVOYMMIVABFZH7ZPZDKWWOPYRS5W23L2JYFH3DEPUU2PAZQ7"
	depositorSK = "SD7D26XD37Y746MUQCAOEHGTPV6F3RVTEUPIJUILXV4O3EG3WORPLVO4"
)

// CreatePoolsUsingJson , loop the pool creation-Josn and call poolcreation method
func CreatePoolsUsingJson(pools []model.BuildPool) ([]model.BuildPool, error, bool) {
	var createdPools []model.BuildPool
	var isCreated bool
	for i := 0; i < len(pools); i++ {
		pool, err, created := CreatePool(pools[i])
		if err == nil {
			createdPools = append(createdPools, pool)
			isCreated = created
		} else {
			logrus.Error(err)
		}
	}
	return createdPools, nil, isCreated
}

// CreatePool carate the pool by handling all blackchain calls
func CreatePool(buildPool model.BuildPool) (model.BuildPool, error, bool) {
	var depostHash2 string

	poolId, err := GeneratePoolId(buildPool.Coin1, buildPool.Coin2)
	if !err {
		return model.BuildPool{}, errors.New("Can not create poold Id"), false
	}
	poolIdString := hex.EncodeToString(poolId[:])

	//calling the stellar endpoint to see if the pool exists 
	url1 := commons.GetHorizonClient().HorizonURL + "liquidity_pools/" + poolIdString
	result, err0 := http.Get(url1)
	if err0 != nil{
		log.Error("Unable to reach Stellar network",url1)
		return model.BuildPool{}, err0, false
	}
	if result.StatusCode == 404{
		logrus.Info("Pool " + poolIdString + " is not created")

		//deposite to the pool
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
	}
	if result.StatusCode == 200{
		poolInfo, err := ioutil.ReadAll(result.Body)
		if err != nil{
			log.Error(err)
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

		//checking if the pool has been deposited
		if (amount1 == "0.0000000" && amount2 == "0.0000000"){
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
		}else{
			logrus.Info("Pool " + poolIdString + " is deposited")
			return model.BuildPool{}, nil, false
		}
	}
	// CreateCoin(buildPool.Coin1, depositorPK, depositorSK)
	// CreateCoin(buildPool.Coin2, depositorPK, depositorSK)
	// IssueCoin(buildPool.Coin1, depositorPK, buildPool.DepositeAmountCoin1)
	// IssueCoin(buildPool.Coin2, depositorPK, buildPool.DepositeAmountCoin2)
	
	// // EstablishPoolTrustline , this method create trust line(change trusst opration) with pool
	// _, err1 := EstablishPoolTrustline(buildPool.Coin1, buildPool.Coin2, depositorPK, depositorSK)
	// if err1 != nil {
	// 	return model.BuildPool{}, err1
	// }
	// // Checked the Asset's names are in lexicographic order if not change the DepositeAmountCoin1 and asset DepositeAmountCoin2
	// if buildPool.Coin1 < buildPool.Coin2 {
	// 	depostHash, err := DepositeToPool(poolId, depositorPK, depositorSK, buildPool.DepositeAmountCoin1, buildPool.DepositeAmountCoin2)
	// 	if err != nil {
	// 		return model.BuildPool{}, err
	// 	}
	// 	depostHash2 = depostHash
	// } else {
	// 	depostHash, err := DepositeToPool(poolId, depositorPK, depositorSK, buildPool.DepositeAmountCoin2, buildPool.DepositeAmountCoin1)
	// 	if err != nil {
	// 		return model.BuildPool{}, err
	// 	}
	// 	depostHash2 = depostHash
	// }
	// buildPool.PoolDepositeHash = depostHash2
	// buildPool.PoolId = poolIdString
	// logrus.Info(depostHash2, err)
	// return buildPool, nil

	return model.BuildPool{}, nil, false
}
