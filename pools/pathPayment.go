package pools

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dileepaj/tracified-gateway/commons"
	log "github.com/sirupsen/logrus"
	sdk "github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

type Coin struct {
	CoinName string
	Amount  string
}

type BuildPathPayment struct {
	SendingCoin        Coin
	IntermediateCoins  []Coin
	ReceivingCoin      Coin
	BatchAccountPK     string
	BatchAccountSK     string
	CoinIssuerAccontPK string
	PoolId             string
	ProductId          string
	EquationId         string
	TenantId           string
}

func CoinConvert(pathPayment BuildPathPayment) (string, error) {
	_, err := CreateCoin(pathPayment.SendingCoin.CoinName, pathPayment.BatchAccountPK, pathPayment.BatchAccountSK)
	if err != nil {
		return "", err
	}
	_, err1 := CreateCoin(pathPayment.ReceivingCoin.CoinName, pathPayment.BatchAccountPK, pathPayment.BatchAccountSK)
	if err1 != nil {
		return "", err1
	}
	_, err2 := IssueCoin(pathPayment.SendingCoin.CoinName, pathPayment.BatchAccountPK, pathPayment.SendingCoin.Amount)
	if err2 != nil {
		return "", err2
	}

	traderAccount, err := client.AccountDetail(sdk.AccountRequest{AccountID: pathPayment.BatchAccountPK})
	if err != nil {
		return "", err
	}
	trader, err := keypair.ParseFull(pathPayment.BatchAccountSK)
	if err != nil {
		return "", err
	}

	sendAsset, err := txnbuild.CreditAsset{pathPayment.SendingCoin.CoinName, pathPayment.CoinIssuerAccontPK}.ToAsset()
	destAsset, err := txnbuild.CreditAsset{pathPayment.ReceivingCoin.CoinName, pathPayment.CoinIssuerAccontPK}.ToAsset()

	op := txnbuild.PathPaymentStrictSend{
		SendAsset:   sendAsset,
		SendAmount:  pathPayment.SendingCoin.Amount,
		Destination: pathPayment.BatchAccountPK,
		DestAsset:   destAsset,
		DestMin:     pathPayment.ReceivingCoin.Amount,
		Path:        []txnbuild.Asset{},
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &traderAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&op},
			BaseFee:              txnbuild.MinBaseFee,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
		},
	)
	check(err)

	signedTx, err := tx.Sign(network.TestNetworkPassphrase, trader)
	check(err)

	resp, err := client.SubmitTransaction(signedTx)
	check(err)
	if err != nil {
		return "", err
	} else {
		return resp.Hash, nil
	}
}

func CheckTrustline() {
}

func GetConvertedCoinAmount(from string, fromAmount string, intermediate string, to string, assetIssuer string) (string, error) {
	result, err := http.Get(commons.GetHorizonClient().HorizonURL + "paths/strict-send?source_asset_type=credit_alphanum4&source_asset_code=" + from + "&source_asset_issuer=" + assetIssuer + "&source_amount=" + fromAmount + "&destination_assets=" + to + "%3A" + assetIssuer)
	fmt.Println(commons.GetHorizonClient().HorizonURL + "paths/strict-send?source_asset_type=credit_alphanum4&source_asset_code=" + from + "&source_asset_issuer=" + assetIssuer + "&source_amount=" + fromAmount + "&destination_assets=" + to + "%3A" + assetIssuer)
	if err != nil {
		log.Error("Unable to reach Stellar network in result1")
		return "", err
	}
	if result.StatusCode != 200 {
		return "", errors.New(result.Status)
	}
	b, err := ioutil.ReadAll(result.Body)

	fmt.Println(string(b))
	var r1 []records
	json.Unmarshal([]byte(string(b)), &r1)

	// The object stored in the "birds" key is also stored as
	// a map[string]interface{} type, and its type is asserted from
	// the interface{} type
	fmt.Println(r1)

	return "destination_amount", nil
}

type a struct {
	_embedded embedded
}
type embedded struct {
	records []records
}

type records struct {
	source_asset_type        string
	source_asset_code        string
	source_asset_issuer      string
	source_amount            string
	destination_asset_type   string
	destination_asset_code   string
	destination_asset_issuer string
	destination_amount       string
	path                     []string
}

// issue
// Public Key	GAMO5NXHTOBD3IFKXVTMBCK2SIF7U6NXXUQKKQ6BYGFMZTDNQGOGF6D4
// Secret Key	SC4KP4JBJZQIYCHMVNCWMYMQN5PSI54SPUXKJMCWZQM6SOXBHLMXKSNM

// dis
// Public Key	GCQ6FGXZWQBXRKZOOZDOXRRHEAXL5UWXB5B45E2KNTPZC75F24K57VCA
// Secret Key	SCRZSFE5QFCRIFUGTJ5TGFHD7GTYAAYOR2ZZWVQV2JFHR2PNB45DNY55
