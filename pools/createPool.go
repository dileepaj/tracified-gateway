package pools

import (
	"errors"
)

var (
	depositorPK = "GBPI4RF4IUOXTQ7XHFHPOPBMJ6KSMHT7SCHBFDQ7MDBTZGITVFPQWGYZ"
	depositorSK = "SBILQPUR3BXXZN2O6RFCJO2RMUJ2JXHUSZQN5C6DW67I632SDU4EISFN"
)

type BuildPool struct {
	Coin1               string
	DepositeAmountCoin1 string
	Coin2               string
	DepositeAmountCoin2 string
	Ratio               int
	poolDepositorPK     string
	PoolId              string
	ProductId           string
	EquationId          string
	TenantId            string
}

func CreatePoolsUsingJson(pools []BuildPool) {
	for i := 0; i < len(pools); i++ {
		CreatePool(pools[i])
	}
}

func CreatePool(buildPool BuildPool) (string, error) {
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

	return depostHash2, nil
}
