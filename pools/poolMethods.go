package pools

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
	"strconv"

	//"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
	sdk "github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/support/log"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

// TODO create a common method to submit all pool-related blockchain transactions

var (
	coinIsuserPK = "GD5WHVKSMPQNFDJOPOMMFVVTLRT5GZXBVSVUZZAWMW6KPS4WQABSG4GH"
	coinIsuserSK = "SA6P6HYPLYMWJUYI3XTMHKJ3CU5EYQKZVSURHC64F5RNHLXFWPXCMQYP"
)
var poolCoin []txnbuild.Asset

var client = sdk.DefaultTestNetClient

// IssueCoin ==> issue coin to given specific account
func IssueCoin(coinName string, coinReceiverPK string, amount string) (string, error) {
	issuerAccount, err := client.AccountDetail(sdk.AccountRequest{AccountID: coinIsuserPK})
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	issuer, err := keypair.ParseFull(coinIsuserSK)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	coin, err := txnbuild.CreditAsset{Code: coinName, Issuer: coinIsuserPK}.ToAsset()
	if err != nil {
		logrus.Error(err)
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
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	resp, err := client.SubmitTransaction(signedTx)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	if err != nil {
		logrus.Error(err)
		return "", err
	} else {
		logrus.Info("IssueCoin ", resp.Hash)
		return resp.Hash, nil
	}
}

// CreateCoin ==> create the coin (create trust line)
func CreateCoin(coinName string, coinReceiverPK string, coinReciverSK string) (string, error) {
	// validate weather the asset is issued by the issuer previously
	assetIssued := assetIssued(coinName)
	logrus.Info(assetIssued)

	// validate weather there is a trustline for the relevent assset
	trustLineCreated := trustlineCreated(coinName, coinReceiverPK)

	// if asset is not issued and there is no DB records, then complete the transaction
	if trustLineCreated == false {
		// Load the corresponding account for both A and C.
		coinReceiverAccount, err := client.AccountDetail(sdk.AccountRequest{AccountID: coinReceiverPK})
		if err != nil {
			logrus.Error(err)
			return "", err
		}
		coinReceiver, err := keypair.ParseFull(coinReciverSK)
		if err != nil {
			logrus.Error(err)
			return "", err
		}
		coin, err := txnbuild.CreditAsset{Code: coinName, Issuer: coinIsuserPK}.ToChangeTrustAsset()
		if err != nil {
			logrus.Error(err)
			return "", err
		}
		// First, the receiving (distribution) account must trust the asset from the issuer.
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
		if err != nil {
			logrus.Error(err)
			return "", err
		}
		resp, err := client.SubmitTransaction(signedTx)
		if err != nil {
			logrus.Error(err)
			return "", err
		}
		if err != nil {
			logrus.Error(err)
			return "", err
		} else {
			// add trustline to DB
			InsertTrustline(coinName, coinReceiverPK)
			logrus.Info("CreateCoin ", resp.Hash)
			return resp.Hash, nil
		}
	} else {
		logrus.Info("Trustline for the relevent assset alredy created")
		return "", nil
	}
}

// CreateCoin ==> create the coin (create trust line)
func CreateCoinSponsering(coinName string, coinReceiverPK string, coinReciverSK string) (string, error) {
	// validate weather there is a trustline for the relevent assset
	trustLineCreated := trustlineCreated(coinName, coinReceiverPK)

	// if asset is not issued and there is no DB records, then complete the transaction
	if trustLineCreated == false {
		sponserAccount, err := client.AccountDetail(sdk.AccountRequest{AccountID: sponsorPK})
		if err != nil {
			logrus.Error(err)
			return "", err
		}
		coinReceiverSign, err := keypair.ParseFull(coinReciverSK)
		if err != nil {
			logrus.Error(err)
			return "", err
		}
		coin, err := txnbuild.CreditAsset{Code: coinName, Issuer: coinIsuserPK}.ToChangeTrustAsset()
		if err != nil {
			logrus.Error(err)
			return "", err
		}
		sponserSign, err := keypair.ParseFull(sponsorSK)
		if err != nil {
			logrus.Error(err)
			return "", err
		}

		sponseringChangeTrust:= []txnbuild.Operation{
			&txnbuild.BeginSponsoringFutureReserves{
				SponsoredID:   coinReceiverPK,
				SourceAccount: sponsorPK,
			},
	
			&txnbuild.ChangeTrust{
				Line:          coin,
				Limit:         "",
				SourceAccount: coinReceiverPK,
			},
			&txnbuild.EndSponsoringFutureReserves{
				SourceAccount: coinReceiverPK,
			},
		}

		// First, the receiving (distribution) account must trust the asset from the issuer.
		tx, err := txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        &sponserAccount,
				IncrementSequenceNum: true,
				Operations:           sponseringChangeTrust,
				BaseFee:              txnbuild.MinBaseFee,
				Memo:                 nil,
				Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
			},
		)
		if err != nil {
			logrus.Error(err)
			return "", err
		}
		signedTx, err := tx.Sign(network.TestNetworkPassphrase, coinReceiverSign,sponserSign)
		if err != nil {
			logrus.Error(err)
			return "", err
		}
		resp, err := client.SubmitTransaction(signedTx)
		if err != nil {
			logrus.Error(err)
			return "", err
		}
		if err != nil {
			logrus.Error(err)
			return "", err
		} else {
			// add trustline to DB
			InsertTrustline(coinName, coinReceiverPK)
			logrus.Info("CreateCoin ", resp.Hash)
			return resp.Hash, nil
		}
	} else {
		logrus.Info("Trustline for the relevent assset alredy created")
		return "", nil
	}
}

// func orderAsset(a string, aVlaue int64, b string, bValue int64) []txnbuild.Asset {
// 	poolCoin = []txnbuild.Asset{}
// 	coinA, err1 := txnbuild.CreditAsset{Code: a, Issuer: coinIsuserPK}.ToAsset()
// 	check(err1)
// 	coinB, err2 := txnbuild.CreditAsset{Code: b, Issuer: coinIsuserPK}.ToAsset()
// 	check(err2)

// GeneratePoolId return the poolid specific to account
func GeneratePoolId(a string, b string) (txnbuild.LiquidityPoolId, bool) {
	coinA, err1 := txnbuild.CreditAsset{Code: a, Issuer: coinIsuserPK}.ToAsset()
	if err1 != nil {
		logrus.Error(err1)
		return txnbuild.LiquidityPoolId{}, false
	}
	coinB, err2 := txnbuild.CreditAsset{Code: b, Issuer: coinIsuserPK}.ToAsset()
	if err2 != nil {
		logrus.Error(err2)
		return txnbuild.LiquidityPoolId{}, false
	}
	poolId, err := txnbuild.LiquidityPoolShareChangeTrustAsset{LiquidityPoolParameters: txnbuild.LiquidityPoolParameters{
		AssetA: coinA,
		AssetB: coinB,
		Fee:    txnbuild.LiquidityPoolFeeV18,
	}}.GetLiquidityPoolID()
	if !err {
		poolId, err = txnbuild.LiquidityPoolShareChangeTrustAsset{LiquidityPoolParameters: txnbuild.LiquidityPoolParameters{
			AssetA: coinB,
			AssetB: coinA,
			Fee:    txnbuild.LiquidityPoolFeeV18,
		}}.GetLiquidityPoolID()
	}
	logrus.Info("PoolId ", poolId)
	return poolId, err
}

// EstablishPoolTrustline return the blockchain hash that use to create trsutline with pool
func EstablishPoolTrustline(a string, b string, coinReceiverPK string, coinReciverSK string) (string, error) {
	poolCoin = []txnbuild.Asset{}
	// assets array
	coins := []string{a, b}

	// reodered the Asset's names in lexicographic order (if not fail the operation)
	sort.Strings(coins)
	coinA, err1 := txnbuild.CreditAsset{Code: coins[0], Issuer: coinIsuserPK}.ToAsset()
	if err1 != nil {
		logrus.Error(err1)
		return "", err1
	}
	coinB, err2 := txnbuild.CreditAsset{Code: coins[1], Issuer: coinIsuserPK}.ToAsset()
	if err2 != nil {
		logrus.Error(err2)
		return "", err2
	}
	poolShareAsset := txnbuild.LiquidityPoolShareChangeTrustAsset{LiquidityPoolParameters: txnbuild.LiquidityPoolParameters{
		AssetA: coinA,
		AssetB: coinB,
		Fee:    txnbuild.LiquidityPoolFeeV18,
	}}

	distributorAccount, err := client.AccountDetail(sdk.AccountRequest{AccountID: coinReceiverPK})
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	distributor, err := keypair.ParseFull(coinReciverSK)
	if err != nil {
		logrus.Error(err)
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
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	resp, err := client.SubmitTransaction(signedTx)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	check(err)
	if err != nil {
		return "", err
	} else {
		log.Info("EstablishPoolTrustline ", resp.Hash)
		return resp.Hash, nil
	}
}

// DepositeToPool, deposite the coin to pool
// TODO return thepool Id
func DepositeToPool(poolId txnbuild.LiquidityPoolId, coinReceiverPK string, coinReciverSK string, maxReserveA string, maxReserveB string) (string, error) {
	distributorAccount, err := client.AccountDetail(sdk.AccountRequest{AccountID: coinReceiverPK})
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	distributor, err := keypair.ParseFull(coinReciverSK)
	if err != nil {
		logrus.Error(err)
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
		logrus.Error(err)
		return "", err
	}
	signedTx, err := tx.Sign(network.TestNetworkPassphrase, distributor)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	resp, err := client.SubmitTransaction(signedTx)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	if err != nil {
		logrus.Error(err)
		return "", err
	} else {
		log.Info("DepositeToPool ", resp.Hash)
		return resp.Hash, nil
	}
}

// check error checker logs the errors
func check(err error) {
	if err != nil {
		logrus.Error(err)
	}
}

//roundFloat return the round float number
func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// check if the issuer has issued the assets
func assetIssued(coinName string) bool {
	result, err := http.Get("https://horizon-testnet.stellar.org/assets?asset_code=" + coinName + "&asset_issuer=" + coinIsuserPK)
	if err != nil {
		log.Error("Error while loading assets for " + coinIsuserPK + err.Error())
	}

	assetsInfo, err1 := ioutil.ReadAll(result.Body)
	if err1 != nil {
		log.Error("Error while reading the respone " + err.Error())
	}

	var raw map[string]interface{}
	var raw1 []interface{}
	json.Unmarshal(assetsInfo, &raw)

	out1, _ := json.Marshal(raw["_embedded"])
	json.Unmarshal(out1, &raw)

	out2, _ := json.Marshal(raw["records"])
	json.Unmarshal(out2, &raw1)

	// checking if the raw is empty
	if len(raw1) == 0 {
		log.Info("Asset is not issued yet")
		return false
	} else {
		return true
	}
}

// Insert trustline to DB
func InsertTrustline(coinName string, coinReceiverPK string) {
	trustlineHistory := model.TrustlineHistory{
		CoinIssuer:   coinIsuserPK,
		CoinReceiver: coinReceiverPK,
		Asset:        coinName,
	}

	object := dao.Connection{}
	err := object.InsertTrustlineHistory(trustlineHistory)
	if err != nil {
		log.Error("Error when inserting trustline to DB " + err.Error())
	} else {
		log.Info("Trustline added to the DB")
	}
}

// check if a trustline is already created for a particular asset
func trustlineCreated(coinName string, coinReceiverPK string) bool {
	object := dao.Connection{}

	data, _ := object.GetTrustline(coinName, coinIsuserPK, coinReceiverPK).Then(func(data interface{}) interface{} {
		return data
	}).Await()

	if data == nil {
		logrus.Info("No trustlines created")
		return false
	} else {
		logrus.Info("Trustline already created")
		return true
	}
}
