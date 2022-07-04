package pools

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

var (
	depositorPK = "GCFECJ5KFYMS4RWCKTPG7PHBIYBKI64R5RMFVE2VDGGFKVRZEWHZ5DZC"
	depositorSK = "SBDZYQ2DYC3KBCQJPXYWYFINPVCY2N5ARIVM5E5EHPZ2SOJD2IGUSSU7"
)

// CreatePoolsUsingJson , loop the pool creation-Josn and call poolcreation method
func CreatePoolsUsingJson(pools []model.BuildPool) ([]model.BuildPool, error) {
	var createdPools []model.BuildPool
	for i := 0; i < len(pools); i++ {
		pool, err := CreatePool(pools[i])
		if err == nil {
			createdPools = append(createdPools, pool)
		} else {
			logrus.Error(err)
		}
	}
	return createdPools, nil
}

// CreatePool carate the pool by handling all blackchain calls
func CreatePool(buildPool model.BuildPool) (model.BuildPool, error) {
	var depostHash2 string
	CreateCoin(buildPool.Coin1, depositorPK, depositorSK)
	CreateCoin(buildPool.Coin2, depositorPK, depositorSK)
	IssueCoin(buildPool.Coin1, depositorPK, buildPool.DepositeAmountCoin1)
	IssueCoin(buildPool.Coin2, depositorPK, buildPool.DepositeAmountCoin2)
	poolId, err := GeneratePoolId(buildPool.Coin1, buildPool.Coin2)
	if !err {
		return model.BuildPool{}, errors.New("Can not create poold Id")
	}
	// EstablishPoolTrustline , this method create trust line(change trusst opration) with pool
	_, err1 := EstablishPoolTrustline(buildPool.Coin1, buildPool.Coin2, depositorPK, depositorSK)
	if err1 != nil {
		return model.BuildPool{}, err1
	}
	// Checked the Asset's names are in lexicographic order if not change the DepositeAmountCoin1 and asset DepositeAmountCoin2
	if buildPool.Coin1 < buildPool.Coin2 {
		depostHash, err := DepositeToPool(poolId, depositorPK, depositorSK, buildPool.DepositeAmountCoin1, buildPool.DepositeAmountCoin2)
		if err != nil {
			return model.BuildPool{}, err
		}
		depostHash2 = depostHash
	} else {
		depostHash, err := DepositeToPool(poolId, depositorPK, depositorSK, buildPool.DepositeAmountCoin2, buildPool.DepositeAmountCoin1)
		if err != nil {
			return model.BuildPool{}, err
		}
		depostHash2 = depostHash
	}
	buildPool.PoolDepositeHash = depostHash2
	logrus.Info(depostHash2, err)
	return buildPool, nil
}
