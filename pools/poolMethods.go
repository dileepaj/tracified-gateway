package pools

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
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

var (
	coinIseerPK = "GBRCIPHDMVGMQUUCP2DWHB55RMZOVL6JPE4KCH2AS2MODVHL6NHC642R"
	coinIsserSK = "SDM7UK2J3DYGQYU5DNES5PP5BKTQTOZ3U6OHXX367OBWBEYTUULMDIF6"
)
var poolCoin []txnbuild.Asset

// coinIseerPK = "GBRCIPHDMVGMQUUCP2DWHB55RMZOVL6JPE4KCH2AS2MODVHL6NHC642R"
// coinIsserSK = "SDY2GF4NBSR6WDTAOWUCCHGCQTIDNOVGZ5KH2XHOOF4FPDQANGEJVCDR"

// depositorCoin
//Public Key	GBPI4RF4IUOXTQ7XHFHPOPBMJ6KSMHT7SCHBFDQ7MDBTZGITVFPQWGYZ
//Secret Key	SBILQPUR3BXXZN2O6RFCJO2RMUJ2JXHUSZQN5C6DW67I632SDU4EISFN

// trader
// Public Key	GCBZ7J5434MIU3AYKCI2FPMLBV5LQBKIZYG2C5QMVEWOTIT2XM2AVWSG
// Secret Key	SA4C7PM67PYJQ2SMRRXDUIX5EUMV725JGDXZXMLKG2VPLW4UYHJLUVSI

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

	//validate weather the asset is issued by the issuer previously
	//assetIssued := assetIssued(coinName)
	//fmt.Println(assetIssued)

	//validate weather there is a trustline for the relevent assset 
	trustLineCreated := trustlineCreated(coinName, coinReceiverPK)
	//fmt.Println(trustLineCreated)

	//if asset is not issued and there is no DB records, then complete the transaction
	if(trustLineCreated == false){
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
			//add trustline to DB
			InsertTrustline(coinName, coinReceiverPK)
			return resp.Hash, nil
		}
	} else{
		return "", nil
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

//check if the issuer has issued the assets
func assetIssued(coinName string) bool{
	//fmt.Println(coinName)
	result, err := http.Get("https://horizon-testnet.stellar.org/assets?asset_code=" + coinName +  "&asset_issuer=" + coinIseerPK)
	if err != nil {
		log.Error("Error while loading assets for " + coinIseerPK + err.Error())
	}

	assetsInfo, err1 := ioutil.ReadAll(result.Body)
	if err1 != nil{
		log.Error("Error while reading the respone " +  err.Error())
	}

	var raw map[string]interface{}
	var raw1 []interface{}
	json.Unmarshal(assetsInfo, &raw)

	out1, _ := json.Marshal(raw["_embedded"])
	json.Unmarshal(out1, &raw)

	out2, _ := json.Marshal(raw["records"])
	json.Unmarshal(out2, &raw1)

	//checking if the raw is empty
	if len(raw1) == 0 {
		log.Info("Asset is not issued yet")
		return false
	} else{
		return true
	}

}

//Insert trustline to DB
func InsertTrustline(coinName string, coinReceiverPK string){
	trustlineHistory := model.TrustlineHistory{
		CoinIssuer: coinIseerPK,
		CoinReceiver: coinReceiverPK,
		Asset: coinName,
	}

	object := dao.Connection{}
	err := object.InsertTrustlineHistory(trustlineHistory)
	if err != nil{
		log.Error("Error when inserting trustline to DB " + err.Error())
	} else{
		log.Info("Trustline added to the DB")
	}
}

//check if a trustline is already created for a particular asset
func trustlineCreated(coinName string, coinReceiverPK string) bool{
	object := dao.Connection{}

	data,_ := object.GetTrustline(coinName, coinIseerPK, coinReceiverPK).Then(func(data interface{}) interface{}{
		return data
	}).Await()

	if data == nil{
		//fmt.Println("No trustlines created")
		return false
	} else{
		//fmt.Println("Trustline already created")
		return true
	}

}
