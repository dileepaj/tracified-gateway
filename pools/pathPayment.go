package pools

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	log "github.com/sirupsen/logrus"
	sdk "github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func CoinConvert(pathPayment model.BuildPathPayment) (string, error) {
	destinationAmount, err0 := GetConvertedCoinAmount(pathPayment.SendingCoin.CoinName, pathPayment.SendingCoin.Amount, pathPayment.ReceivingCoin.CoinName, pathPayment.CoinIssuerAccontPK)
	if err0 != nil {
		return "", err0
	}
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

	sendAsset := txnbuild.CreditAsset{pathPayment.SendingCoin.CoinName, pathPayment.CoinIssuerAccontPK}
	check(err)
	destAsset := txnbuild.CreditAsset{pathPayment.ReceivingCoin.CoinName, pathPayment.CoinIssuerAccontPK}
	check(err)

	var intermediateAssertArray []txnbuild.Asset
	for i := 0; i < len(pathPayment.IntermediateCoins); i++ {
		intermediateAsset := txnbuild.CreditAsset{pathPayment.IntermediateCoins[i].CoinName, pathPayment.CoinIssuerAccontPK}
		check(err)
		intermediateAssertArray = append(intermediateAssertArray, intermediateAsset)
	}

	op := txnbuild.PathPaymentStrictSend{
		SendAsset:     sendAsset,
		SendAmount:    pathPayment.SendingCoin.Amount,
		Destination:   pathPayment.BatchAccountPK,
		DestAsset:     destAsset,
		DestMin:       destinationAmount,
		Path:          intermediateAssertArray,
		SourceAccount: traderAccount.AccountID,
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
	if err != nil {
		return "", err
	}

	signedTx, err := tx.Sign(network.TestNetworkPassphrase, trader)
	if err != nil {
		return "", err
	}

	resp, err := client.SubmitTransaction(signedTx)
	if err != nil {
		return "", err
	} else {
		return resp.Hash, nil
	}
}

// get distination recived coin ammount after converting the coin
func GetConvertedCoinAmount(from string, fromAmount string, to string, assetIssuer string) (string, error) {
	result, err := http.Get(commons.GetHorizonClient().HorizonURL + "paths/strict-send?source_asset_type=credit_alphanum4&source_asset_code=" + from + "&source_asset_issuer=" + assetIssuer + "&source_amount=" + fromAmount + "&destination_assets=" + to + "%3A" + assetIssuer)
	if err != nil {
		log.Error("Unable to reach Stellar network in result1")
		return "", err
	}
	if result.StatusCode != 200 {
		return "", errors.New(result.Status)
	}
	defer result.Body.Close()
	coinconvertionInfo, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return "", err
	}
	var raw map[string]interface{}
	var raw1 []interface{}
	json.Unmarshal(coinconvertionInfo, &raw)

	out1, _ := json.Marshal(raw["_embedded"])
	json.Unmarshal(out1, &raw)

	out2, _ := json.Marshal(raw["records"])
	json.Unmarshal(out2, &raw1)

	record := raw1[0].(map[string]interface{})
	destinationAmount := fmt.Sprintf("%v", record["destination_amount"])

	if destinationAmount == "" {
		return destinationAmount, errors.New("Destination amount is empty")
	}
	return destinationAmount, nil
}
