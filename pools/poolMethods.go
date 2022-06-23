package pools

import (
	"math"
	"strconv"

	"github.com/sirupsen/logrus"
	sdk "github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

var (
	coinIseerPK = "GBSLTJX4NRMTPTQ2EJZJ24U44K7ZWY3CPGBZTV623PLTLIFXWK3T4CD6"
	coinIsserSK = "SAC6WJ46JO5JG6D5GWPCK6JLY52KIIJT4RZXNAMHUUB23QSUBIUXIMSZ"
)
var poolCoin []txnbuild.Asset

// coinIseerPK = "GDBXHHHG7CKIODJIUPU46W52RUDUMJ3PJQOSWF24R3VGVRLPVHWNT5DI"
// coinIsserSK = "SDY2GF4NBSR6WDTAOWUCCHGCQTIDNOVGZ5KH2XHOOF4FPDQANGEJVCDR"

// depositorCoin
//Public Key	GCLDZY7I6TXTJ5LY4B6RJELSEGPPUYVXUAC7CRW546EJ2TD77WHPYVED
//Secret Key	SCGHVQQKUMWLCDV2FVDAMCFFTMT2NALYLWA7SAPMA67QCME5W6CABPSI

// trader
// Public Key	GD6ZW4L3Y5E3JEW4TLSYGV3PC7TBYN6AXIGVW54J6HRH2J3HDZMBA62C
// Secret Key	SD6G4TWP5PTCKIO4XOHCZE5IDJDNOEIVTOCRC6YRG5B3IO42SRMWYKU4

var client = sdk.DefaultTestNetClient

func IssueCoin(coinName string, coinReceiverPK string, amount string) (string, error) {
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
			Operations:           []txnbuild.Operation{&txnbuild.Payment{Destination: coinReceiverPK, Asset: coin, Amount: amount}},
			BaseFee:              txnbuild.MinBaseFee,
			Memo:                 nil,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
		},
	)

	signedTx, err := tx.Sign(network.TestNetworkPassphrase, issuer)
	check(err)
	resp, err := client.SubmitTransaction(signedTx)
	check(err)
	if err != nil {
		return "", err
	} else {
		return resp.Hash, nil
	}
}

func CreateCoin(coinName string, coinReceiverPK string, coinReciverSK string) (string, error) {
	// Load the corresponding account for both A and C.
	coinReceiverAccount, err := client.AccountDetail(sdk.AccountRequest{AccountID: coinReceiverPK})
	if err != nil {
		return "", err
	}
	coinReceiver, err := keypair.ParseFull(coinReciverSK)
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
			SourceAccount:        &coinReceiverAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&txnbuild.ChangeTrust{Line: coin, Limit: "", SourceAccount: ""}},
			BaseFee:              txnbuild.MinBaseFee,
			Memo:                 nil,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
		},
	)
	signedTx, err := tx.Sign(network.TestNetworkPassphrase, coinReceiver)
	check(err)
	resp, err := client.SubmitTransaction(signedTx)
	check(err)
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
	coinA, err1 := txnbuild.CreditAsset{Code: a, Issuer: coinIseerPK}.ToAsset()
	check(err1)
	coinB, err2 := txnbuild.CreditAsset{Code: b, Issuer: coinIseerPK}.ToAsset()
	check(err2)
	poolId, err := txnbuild.LiquidityPoolShareChangeTrustAsset{LiquidityPoolParameters: txnbuild.LiquidityPoolParameters{
		AssetA: coinA,
		AssetB: coinB,
		Fee:    txnbuild.LiquidityPoolFeeV18,
	}}.GetLiquidityPoolID()
	return poolId, err
}

func EstablishPoolTrustline(a string, b string, coinReceiverPK string, coinReciverSK string) (string, error) {
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
		AssetA: coinA,
		AssetB: coinB,
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
	check(err)
	resp, err := client.SubmitTransaction(signedTx)
	check(err)
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
	reserveA, err := strconv.Atoi(maxReserveA)
	reserveB, err := strconv.Atoi(maxReserveB)
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
					N: xdr.Int32(reserveA),
					D: xdr.Int32(reserveB),
				},
				MaxPrice: xdr.Price{
					N: xdr.Int32(reserveA),
					D: xdr.Int32(reserveB),
				},
			}},
			BaseFee:       txnbuild.MinBaseFee,
			Memo:          nil,
			Preconditions: txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
		},
	)
	if err != nil {
		return "", err
	}
	signedTx, err := tx.Sign(network.TestNetworkPassphrase, distributor)
	check(err)
	resp, err := client.SubmitTransaction(signedTx)
	check(err)
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
