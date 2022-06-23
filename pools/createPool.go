package pools

import (
	"errors"
	"fmt"
)

var (
	depositorPK = "GCLDZY7I6TXTJ5LY4B6RJELSEGPPUYVXUAC7CRW546EJ2TD77WHPYVED"
	depositorSK = "SCGHVQQKUMWLCDV2FVDAMCFFTMT2NALYLWA7SAPMA67QCME5W6CABPSI"
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
	fmt.Println( depostHash2,err)
	return depostHash2, nil
}
