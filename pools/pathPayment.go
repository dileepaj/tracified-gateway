package pools

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	sdk "github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

// CoinConvert convert the coin (do a path payment operation by sponsering)
func CoinConvert(pathPayment model.BuildPathPayment) (model.BuildPathPayment, error) {
	if pathPayment.SendingCoin.GeneratedName == "" || pathPayment.SendingCoin.Amount == "" || pathPayment.ReceivingCoin.CoinName == "" || pathPayment.CoinIssuerAccontPK == "" {
		log.Error("CoinConvert() method's parameters have a empty values")
		return model.BuildPathPayment{}, errors.New("metric coin or input coins can not be empty")
	}
	convertion, err := GetConvertedCoinAmount(pathPayment.SendingCoin.GeneratedName, pathPayment.SendingCoin.Amount, pathPayment.ReceivingCoin.GeneratedName, pathPayment.CoinIssuerAccontPK)
	if err != nil {
		logrus.Error(err)
		return model.BuildPathPayment{}, err
	}
	_, err0 := CreateCoinSponsering(pathPayment.SendingCoin.GeneratedName, pathPayment.BatchAccountPK, pathPayment.BatchAccountSK)
	if err0 != nil {
		logrus.Error(err0)
		return model.BuildPathPayment{}, err0
	}
	_, err1 := CreateCoinSponsering(pathPayment.ReceivingCoin.GeneratedName, pathPayment.BatchAccountPK, pathPayment.BatchAccountSK)
	if err1 != nil {
		logrus.Error(err1)
		return model.BuildPathPayment{}, err1
	}
	_, err2 := IssueCoin(pathPayment.SendingCoin.GeneratedName, pathPayment.BatchAccountPK, pathPayment.SendingCoin.Amount)
	if err2 != nil {
		logrus.Error(err2)
		return model.BuildPathPayment{}, err2
	}

	sponserAccount, err := client.AccountDetail(sdk.AccountRequest{AccountID: sponsorPK})
	if err != nil {
		logrus.Error(err)
		return model.BuildPathPayment{}, err
	}

	// batch account and tradring account are same
	traderSign, err := keypair.ParseFull(pathPayment.BatchAccountSK)
	if err != nil {
		logrus.Error(err)
		return model.BuildPathPayment{}, err
	}

	sponserAccountSign, err := keypair.ParseFull(sponsorSK)
	if err != nil {
		logrus.Error(err)
		return model.BuildPathPayment{}, err
	}
	sendAsset := txnbuild.CreditAsset{pathPayment.SendingCoin.GeneratedName, pathPayment.CoinIssuerAccontPK}
	if err != nil {
		logrus.Error(err)
		return model.BuildPathPayment{}, err
	}
	destAsset := txnbuild.CreditAsset{pathPayment.ReceivingCoin.GeneratedName, pathPayment.CoinIssuerAccontPK}
	if err != nil {
		logrus.Error(err)
		return model.BuildPathPayment{}, err
	}

	// intermediateAssertArray coin convertion path as a array(this path take from stellar endpoint)
	var intermediateAssertArray []txnbuild.Asset
	for _, pathCoin := range convertion.IntermediateCoin {
		intermediateAsset := txnbuild.CreditAsset{pathCoin.CoinName, pathCoin.Issuer}
		intermediateAssertArray = append(intermediateAssertArray, intermediateAsset)
	}

	sponsoringPathPayment := []txnbuild.Operation{
		&txnbuild.BeginSponsoringFutureReserves{
			SponsoredID:   pathPayment.BatchAccountPK,
			SourceAccount: sponsorPK,
		},

		&txnbuild.PathPaymentStrictSend{
			SendAsset:     sendAsset,
			SendAmount:    pathPayment.SendingCoin.Amount,
			Destination:   pathPayment.BatchAccountPK,
			DestAsset:     destAsset,
			DestMin:       convertion.Destination.Amount,
			Path:          intermediateAssertArray,
			SourceAccount: pathPayment.BatchAccountPK,
		},
		&txnbuild.EndSponsoringFutureReserves{
			SourceAccount: pathPayment.BatchAccountPK,
		},
	}
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sponserAccount,
			IncrementSequenceNum: true,
			Operations:           sponsoringPathPayment,
			BaseFee:              txnbuild.MinBaseFee,
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
		},
	)
	if err != nil {
		logrus.Error(err)
		return model.BuildPathPayment{}, err
	}

	signedTx, err := tx.Sign(network.TestNetworkPassphrase, traderSign, sponserAccountSign)
	if err != nil {
		logrus.Error(err)
		return model.BuildPathPayment{}, err
	}

	response, err := client.SubmitTransaction(signedTx)
	if err != nil {
		logrus.Error(err)
		return model.BuildPathPayment{}, err
	} else {
		logrus.Info("CoinConverted ", response.Hash)
		pathPayment.ReceivingCoin.Amount = convertion.Destination.Amount
		pathPayment.Hash = response.Hash
		return pathPayment, nil
	}
}

// GetConvertedCoinAmount,  get distination recived coin ammount after converting the coin and get the coin convesion path (using stella call)
func GetConvertedCoinAmount(from string, fromAmount string, to string, assetIssuer string) (model.DestinationCoin, error) {
	var destinationAssert model.DestinationCoin
	url := commons.GetHorizonClient().HorizonURL + "paths/strict-send?source_asset_type=credit_alphanum12&source_asset_code=" + from + "&source_asset_issuer=" + assetIssuer + "&source_amount=" + fromAmount + "&destination_assets=" + to + "%3A" + assetIssuer
	result, err := http.Get(url)
	if err != nil {
		log.Error("Unable to reach Stellar network", url)
		return destinationAssert, err
	}
	if result.StatusCode != 200 {
		return destinationAssert, errors.New(result.Status + " The request you sent to pool assert convertion was invalid in some way" + " " + url)
	}
	defer result.Body.Close()
	coinconvertionInfo, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Error(err)
		return destinationAssert, err
	}
	var raw map[string]interface{}
	var raw1 []interface{}
	var raw2 []interface{}

	json.Unmarshal(coinconvertionInfo, &raw)

	out1, _ := json.Marshal(raw["_embedded"])
	json.Unmarshal(out1, &raw)

	out2, _ := json.Marshal(raw["records"])
	json.Unmarshal(out2, &raw1)
	if len(raw1) <= 0 {
		return destinationAssert, errors.New("pool not Created")
	}
	record := raw1[0].(map[string]interface{})
	// retrive the distination recived coin ammount
	destinationAmount := fmt.Sprintf("%v", record["destination_amount"])

	// set the distination recived coin ammount
	destinationAssert.Destination.Amount = destinationAmount
	destinationAssert.Destination.CoinName = to

	out3, _ := json.Marshal(record["path"])
	json.Unmarshal(out3, &raw2)

	// retrive the coin converion paths and push it to array
	for i := range raw2 {
		path := raw2[i].(map[string]interface{})
		pathAssert := model.CoinPath{
			Type:     fmt.Sprintf("%v", path["asset_type"]),
			CoinName: fmt.Sprintf("%v", path["asset_code"]),
			Issuer:   fmt.Sprintf("%v", path["asset_issuer"]),
		}
		destinationAssert.IntermediateCoin = append(destinationAssert.IntermediateCoin, pathAssert)
	}
	if destinationAmount == "" {
		log.Error("Destination amount is empty" + url)
		return destinationAssert, errors.New("Destination amount is empty")
	}
	return destinationAssert, nil
}

func PathPaymentHandle(newBatchConvertCoinObj model.BatchCoinConvert) (string, error) {
	var batchAccountPK string
	var batchAccountSK string
	var coinConversions []model.BuildPathPayment
	// check if there is an account in the DB for the batchID and get the account
	object := dao.Connection{}
	data, _ := object.GetBatchSpecificAccount(newBatchConvertCoinObj.FormulaTypeName, newBatchConvertCoinObj.EquationID,
		newBatchConvertCoinObj.ProductName, newBatchConvertCoinObj.TenantID).Then(func(data interface{}) interface{} {
		return data
	}).Await()

	if data == nil {
		// add account to the DB
		batchAccount := model.CoinAccount{
			FormulaType:     newBatchConvertCoinObj.FormulaType,
			FormulaTypeID:   newBatchConvertCoinObj.FormulaTypeID,
			FormulaTypeName: newBatchConvertCoinObj.FormulaTypeName,
			TenantID:        newBatchConvertCoinObj.TenantID,
			ProductName:     newBatchConvertCoinObj.ProductName,
			EquationID:      newBatchConvertCoinObj.EquationID,
			StageID:         newBatchConvertCoinObj.StageId,
			MetricCoin:      newBatchConvertCoinObj.MetricCoin,
		}
		// if not create the sponsering account
		batchPK, batchSK, err := CreateSponseredAccount(batchAccount)
		batchAccountPK = batchPK
		batchAccountSK = batchSK

		if err != nil {
			logrus.Error("Can not Create Batch Account " + err.Error())
			return "", err
		}

	} else {

		decryptedPK := (data.(model.CoinAccount)).CoinAccountPK
		decryptedSK := (data.(model.CoinAccount)).CoinAccountSK

		// decrypt account details
		// decryptedPK := commons.Decrypt([]byte(encryptedPK))
		// decryptedSK := commons.Decrypt([]byte(encryptedSK))

		// if there is an account go to path payments directly
		batchAccountPK = decryptedPK 
		batchAccountSK = commons.Decrypt([]byte(decryptedSK))

		logrus.Info("account PK  ", batchAccountPK)
		//logrus.Info("account SK  ", batchAccountPK)

		if batchAccountPK == "" || batchAccountSK == "" {
			logrus.Error("Can not Create Batch Account")
			return "", errors.New("Can not Create Batch Account")
		}

	}

	// CoinConvertionJson return CoinConvertionJson that used to do a coin convert via pools
	pathpayments, err := CoinConvertionJson(newBatchConvertCoinObj, batchAccountPK, batchAccountSK)
	if err != nil {
		logrus.Error("Can not create Path Payment Json ", err)
		return "", err
	}

	for _, pathPayment := range pathpayments {
		coinConversion, err := CoinConvert(pathPayment)
		if err != nil {
			logrus.Error("Coin converion issue ", err)
		} else {
			coinConversions = append(coinConversions, coinConversion)
		}
	}
	if len(coinConversions) <= 0 {
		logrus.Info("Can not convert any Coin ", coinConversions)
		return "", errors.New("an not convert any Coin")
	}
	// build response with all coin details
	buildCoinConvertionResponse := model.BuildPathPaymentJSon{
		CoinConertions: coinConversions,
		ProductIdName:  newBatchConvertCoinObj.ProductName,
		EquationId:     newBatchConvertCoinObj.EquationID,
		TenantId:       newBatchConvertCoinObj.TenantID,
	}
	err1 := object.InsertCoinConversionDetails(buildCoinConvertionResponse)
	if err1 != nil {
		logrus.Error("Error when inserting coin conversion details to DB " + err.Error())
		return "", nil
	} else {
		logrus.Info("Coin conversion details added to the DB")
		out, err := json.Marshal(buildCoinConvertionResponse)
		if err != nil {
			logrus.Info("Path payment to the DB JSON Marshal Error")
		}
		return string(out), nil
	}
}
