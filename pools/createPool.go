package pools

import (
	"math"

	"github.com/sirupsen/logrus"
	sdk "github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

var (
	coinIseerPK = "GDBXHHHG7CKIODJIUPU46W52RUDUMJ3PJQOSWF24R3VGVRLPVHWNT5DI"
	coinIsserSK = "SDY2GF4NBSR6WDTAOWUCCHGCQTIDNOVGZ5KH2XHOOF4FPDQANGEJVCDR"
)
var poolCoin []txnbuild.Asset

// reciverCoin
// Public Key	GADOJTCR6QHI5IZDQEGEZMCHR7SO4LG6R4PYPLMD76HOHKRGXADEWSEL
// Secret Key	SDSEJCPISEOGUGQR6CSX4XECR2N4J75VXHPWXSNWREWG5MZU7YWTWSKX
var client = sdk.DefaultTestNetClient

func IssueCoin(coinReceiverPK string, coinName string) (string, error) {
	issuerAccount, err := client.AccountDetail(sdk.AccountRequest{AccountID: coinIseerPK})
	if err != nil {
		return "", err
	}
	issuer, err := keypair.ParseFull(coinIsserSK)
	if err != nil {
		return "", err
	}
	coin, err := txnbuild.CreditAsset{Code: coinName, Issuer: coinIseerPK}.ToAsset()
	if err != nil {
		return "", err
	}

	// Second, the issuing account actually sends a payment using the asset
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &issuerAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&txnbuild.Payment{Destination: coinReceiverPK, Asset: coin, Amount: "1000"}},
			BaseFee:              txnbuild.MinBaseFee,
			Memo:                 nil,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
		},
	)

	signedTx, err := tx.Sign(network.TestNetworkPassphrase, issuer)
	resp, err := client.SubmitTransaction(signedTx)

	if err != nil {
		return "", err
	} else {
		return resp.Hash, nil
	}
}

func CreateCoin(coinName string, coinReceiverPK string, coinReciverSK string) (string, error) {
	// Load the corresponding account for both A and C.
	distributorAccount, err := client.AccountDetail(sdk.AccountRequest{AccountID: coinReceiverPK})
	if err != nil {
		return "", err
	}
	distributor, err := keypair.ParseFull(coinReciverSK)
	if err != nil {
		return "", err
	}
	coin, err := txnbuild.CreditAsset{Code: coinName, Issuer: coinIseerPK}.ToChangeTrustAsset()
	if err != nil {
		return "", err
	} // First, the receiving (distribution) account must trust the asset from the
	// issuer.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &distributorAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&txnbuild.ChangeTrust{Line: coin, Limit: "", SourceAccount: ""}},
			BaseFee:              txnbuild.MinBaseFee,
			Memo:                 nil,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
		},
	)
	signedTx, err := tx.Sign(network.TestNetworkPassphrase, distributor)
	resp, err := client.SubmitTransaction(signedTx)

	if err != nil {
		return "", err
	} else {
		return resp.Hash, nil
	}
}

func orderAsset(a string, aVlaue int64, b string, bValue int64) []txnbuild.Asset {
	poolCoin = []txnbuild.Asset{}
	coinA, err1 := txnbuild.CreditAsset{Code: a, Issuer: coinIseerPK}.ToAsset()
	check(err1)
	coinB, err2 := txnbuild.CreditAsset{Code: b, Issuer: coinIseerPK}.ToAsset()
	check(err2)

	if aVlaue > bValue {
		poolCoin = append(poolCoin, coinA, coinB)
		return poolCoin
	} else {
		poolCoin = append(poolCoin, coinB, coinA)
		return poolCoin
	}
}

func GeneratePoolId(a string, b string) (txnbuild.LiquidityPoolId, bool) {
	// orderAsset("test1", 1000, "t2", 2000)
	coinA, err1 := txnbuild.CreditAsset{Code: a, Issuer: coinIseerPK}.ToAsset()
	check(err1)
	coinB, err2 := txnbuild.CreditAsset{Code: b, Issuer: coinIseerPK}.ToAsset()
	check(err2)
	poolId, err := txnbuild.LiquidityPoolShareChangeTrustAsset{LiquidityPoolParameters: txnbuild.LiquidityPoolParameters{
		AssetA: coinB,
		AssetB: coinA,
		Fee:    txnbuild.LiquidityPoolFeeV18,
	}}.GetLiquidityPoolID()
	// fmt.Println(poolId, err)
	return poolId, err
}

func establishPoolTrustline(a string, b string, coinReceiverPK string, coinReciverSK string) (string, error) {
	poolCoin = []txnbuild.Asset{}
	coinA, err1 := txnbuild.CreditAsset{Code: a, Issuer: coinIseerPK}.ToAsset()
	if err1 != nil {
		return "", err1
	}
	coinB, err2 := txnbuild.CreditAsset{Code: b, Issuer: coinIseerPK}.ToAsset()
	if err2 != nil {
		return "", err2
	}
	poolShareAsset := txnbuild.LiquidityPoolShareChangeTrustAsset{LiquidityPoolParameters: txnbuild.LiquidityPoolParameters{
		AssetA: coinB,
		AssetB: coinA,
		Fee:    txnbuild.LiquidityPoolFeeV18,
	}}

	distributorAccount, err := client.AccountDetail(sdk.AccountRequest{AccountID: coinReceiverPK})
	if err != nil {
		return "", err
	}
	distributor, err := keypair.ParseFull(coinReciverSK)
	if err != nil {
		return "", err
	}
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &distributorAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&txnbuild.ChangeTrust{Line: poolShareAsset, Limit: "", SourceAccount: ""}},
			BaseFee:              txnbuild.MinBaseFee,
			Memo:                 nil,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
		},
	)

	signedTx, err := tx.Sign(network.TestNetworkPassphrase, distributor)
	resp, err := client.SubmitTransaction(signedTx)

	if err != nil {
		return "", err
	} else {
		return resp.Hash, nil
	}
}

func DepositeToPool(poolId txnbuild.LiquidityPoolId, coinReceiverPK string, coinReciverSK string, maxReserveA string, maxReserveB string) (string, error) {
	distributorAccount, err := client.AccountDetail(sdk.AccountRequest{AccountID: coinReceiverPK})
	if err != nil {
		return "", err
	}
	distributor, err := keypair.ParseFull(coinReciverSK)
	if err != nil {
		return "", err
	}
	exactPrice := 2000.00 /1000.00
	minPrice := exactPrice - (exactPrice * 0.10)
	maxPrice := exactPrice + (exactPrice * 0.10)

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &distributorAccount,
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{&txnbuild.LiquidityPoolDeposit{
				SourceAccount:   distributorAccount.AccountID,
				LiquidityPoolID: poolId,
				MaxAmountA:      maxReserveA,
				MaxAmountB:      maxReserveB,
				MinPrice: xdr.Price{
					N: xdr.Int32(minPrice),
					D: 1,
				},
				MaxPrice: xdr.Price{
					N: xdr.Int32(maxPrice),
					D: 1,
				},
			}},
			BaseFee:       txnbuild.MinBaseFee,
			Memo:          nil,
			Preconditions: txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
		},
	)

	signedTx, err := tx.Sign(network.TestNetworkPassphrase, distributor)
	resp, err := client.SubmitTransaction(signedTx)

	if err != nil {
		return "", err
	} else {
		return resp.Hash, nil
	}
}

func check(err error) {
	if err != nil {
		logrus.Error(err)
	}
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
