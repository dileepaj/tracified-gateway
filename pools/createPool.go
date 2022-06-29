package pools

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

var (
	depositorPK = "GCJGTNEUW27ON4X2LIWBKEY525XQGZFSBIN56CP324VNVYD42BAD3PR2"
	depositorSK = "SANLXA5YA3ECV3WFIJ6WPU5KPE45H46JFQA4CGKMJFXTDJZ35TCTNEIO"
)

func CreatePoolsUsingJson(pools []model.BuildPool) ([]string, error) {
	var hashes []string
	for i := 0; i < len(pools); i++ {
		hash, err := CreatePool(pools[i])
		if err == nil {
			hashes = append(hashes, hash)
		} else {
			logrus.Error(err)
		}
	}
	return hashes, nil
}

func CreatePool(buildPool model.BuildPool) (string, error) {
	var depostHash2 string
	CreateCoin(buildPool.Coin1, depositorPK, depositorSK)
	CreateCoin(buildPool.Coin2, depositorPK, depositorSK)
	IssueCoin(buildPool.Coin1, depositorPK, buildPool.DepositeAmountCoin1)
	IssueCoin(buildPool.Coin2, depositorPK, buildPool.DepositeAmountCoin2)
	poolId, err := GeneratePoolId(buildPool.Coin1, buildPool.Coin2)
	if !err {
		return "", errors.New("Can not create poold Id")
	}
	poolCreationHash, err1 := EstablishPoolTrustline(buildPool.Coin1, buildPool.Coin2, depositorPK, depositorSK)
	if err1 != nil {
		return poolCreationHash, err1
	}
	//Checked the Asset's names are in lexicographic order if not change the DepositeAmountCoin1 and asset DepositeAmountCoin2
	if buildPool.Coin1 < buildPool.Coin2 {
		depostHash, err := DepositeToPool(poolId, depositorPK, depositorSK, buildPool.DepositeAmountCoin1, buildPool.DepositeAmountCoin2)
		if err != nil {
			return "", err
		}
		depostHash2 = depostHash
	} else {
		depostHash, err := DepositeToPool(poolId, depositorPK, depositorSK, buildPool.DepositeAmountCoin2, buildPool.DepositeAmountCoin1)
		if err != nil {
			return "", err
		}
		depostHash2 = depostHash
	}

	logrus.Info(depostHash2, err)
	return depostHash2, nil
}
