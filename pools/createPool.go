package pools

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

var (
	depositorPK = "GDUXXB3FHCHZJJEHJ3ZRVBW4LCCQDQCH7P5KHL2S5EIEZ6DOC2AWXWCE"
	depositorSK = "SDCFURSX7IP4YYHQ4BB6MNCS7IQCC3S7IXM5HHEBXHT256GK5XZEFOCC"
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
	depostHash2, err2 := DepositeToPool(poolId, depositorPK, depositorSK, buildPool.DepositeAmountCoin1, buildPool.DepositeAmountCoin2)
	if err2 != nil {
		return "", err2
	}
	logrus.Info(depostHash2, err)
	return depostHash2, nil
}
