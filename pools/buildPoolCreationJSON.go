package pools

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
)

// BuildPoolCreationJSON return the pool creation json by Restructing the admin equationJson
func BuildPoolCreationJSON(equationJson model.CreatePool) ([]model.BuildPool, error) {
	var poolJson []model.BuildPool
	portion := equationJson.EquationSubPortion
	for i := 0; i < len(portion); i++ {
		// loop the subportion and create a pair array(poolJson) n and n+1
		for j := 0; j < len(portion[i].FieldAndCoin); j++ {
			if j+1 < len(portion[i].FieldAndCoin) {
				//
				ratio := portion[i].FieldAndCoin[j+1].Value

				// ration string split by .
				removeDecimal := strings.Split(ratio, ".")
				// counting the character in a ratio string value and calculating the 10th multiplication factor
				multiplicationZeros := len(removeDecimal[0]) - 1

				// multiplicationFactor==> coin multiplication factor
				// Note maximum 1000000000 coin can be issued by the issuer
				var multiplicationFactor int
				if multiplicationZeros == 0 {
					multiplicationFactor = 100000000
				} else {
					multiplicationFactor = int(100000000 / math.Pow(10, float64(multiplicationZeros)))
				}
				value, _ := strconv.ParseFloat(ratio, 64)
				coinDetails1 := 1 * multiplicationFactor

				coinDetails2 := int(value * float64(multiplicationFactor))
				pool1 := model.BuildPool{
					Coin1:               portion[i].FieldAndCoin[j].GeneratedName,
					DepositeAmountCoin1: strconv.Itoa(coinDetails1),
					Coin2:               portion[i].FieldAndCoin[j+1].GeneratedName,
					DepositeAmountCoin2: strconv.Itoa(coinDetails2),
					Ratio:               ratio,
				}
				poolJson = append(poolJson, pool1)
			}
		}
	}
	return poolJson, nil
}

// RemoveDivisionAndRestructed return restructed equato==ion json by
// turning the divisor (the second fraction) upside down (switching its numerator with its denominator)
// and changing the division symbol to a multiplication symbol at the same time
func RemoveDivisionAndOperator(equationJson model.CreatePool) (model.CreatePool, []model.CoinMap, error) {
	// equation is a equation-Json
	var generatedCoinName string
	portion := equationJson.EquationSubPortion
	var coinMap []model.CoinMap
	if len(portion) > 0 {
		for i := 0; i < len(portion); i++ {
			if len(portion[i].FieldAndCoin) > 2 {
				userInputCount := 0
				for j := 0; j < len(portion[i].FieldAndCoin); j++ {
					if portion[i].FieldAndCoin[j].VariableType != "operator" || portion[i].FieldAndCoin[j].CoinName != "" {
						// timestamp := makeTimestamp()
						// strTimestamp := strconv.Itoa(int(timestamp))
						// generatedCoinName := portion[i].FieldAndCoin[j].CoinName[0:1] +
						// 	equationJson.EquationID[0:3] + equationJson.ProductName[0:1] + equationJson.TenantID[0:4] + strTimestamp[10:13]
						generatedSID, err := shortid.Generate()
						if err != nil {
							logrus.Error("Cannot generate ShortID", err.Error())
						}

						logrus.Info("Generated string ", generatedSID)

						replacerOne := strings.Replace(generatedSID, "-", randomChar(), 12)
						replacerTwo := strings.Replace(replacerOne, "_", randomChar(), 12)

						logrus.Info("Replaced string ", replacerTwo)

						//check if the length is less than 12
						if(len(replacerTwo) < 12){
							remainingChars := 12 - len(replacerTwo)
							if(remainingChars == 3){
								generatedCoinName = portion[i].FieldAndCoin[j].CoinName[0:1] + equationJson.TenantID[0:1] + equationJson.EquationID[0:1] + replacerTwo
							} else if(remainingChars == 2){
								generatedCoinName = portion[i].FieldAndCoin[j].CoinName[0:1] + equationJson.TenantID[0:1] + replacerTwo
							} else if(remainingChars == 1){
								generatedCoinName = portion[i].FieldAndCoin[j].CoinName[0:1] + replacerTwo
							} else{
								generatedCoinName = replacerTwo
							}
						} else{
							generatedCoinName = replacerTwo
						}
	
						logrus.Info("Generated and modified Coin name  01 : ", generatedCoinName)
						portion[i].FieldAndCoin[j].GeneratedName = generatedCoinName
						
					}

					// count the userIput type variable in a  sub portion
					if portion[i].FieldAndCoin[j].VariableType == "userInput" {
						userInputCount++
					}
					if portion[i].FieldAndCoin[j].VariableType == "operator" && portion[i].FieldAndCoin[j].Value == "/" {
						portion[i].FieldAndCoin[j].Value = "*"
						if float, err := strconv.ParseFloat(portion[i].FieldAndCoin[j+1].Value, 32); err == nil {
							var value float64 = 1 / float
							portion[i].FieldAndCoin[j+1].Value = fmt.Sprintf("%f", value)
						}
					}
				}
				// checked whether sub-portion contain only one user input
				if userInputCount != 1 {
					return model.CreatePool{}, []model.CoinMap{}, errors.New("Equation's sub portion can not contain two user input")
				}
			} else {
				return model.CreatePool{}, []model.CoinMap{}, errors.New("Equation-JSON's Sub portions are empty or contaion a one element")
			}
		}
		// reomve the oprator from equation
		generatedSID, err := shortid.Generate()
		if err != nil {
			logrus.Error("Cannot generate ShortID", err.Error())
		}

		logrus.Info("Genereated string ", generatedSID)

		replacerOne := strings.Replace(generatedSID, "-", randomChar(), 12)
		replacerTwo := strings.Replace(replacerOne, "_", randomChar(), 12)

		logrus.Info("Replaced string ", replacerTwo)

		//check if the length is less than 12
		if(len(replacerTwo) < 12){
			remainingChars := 12 - len(replacerTwo)
			if(remainingChars == 3){
				generatedCoinName = equationJson.MetricCoin.CoinName[0:1] + equationJson.TenantID[0:1] + equationJson.EquationID[0:1] + replacerTwo
			} else if(remainingChars == 2){
				generatedCoinName = equationJson.MetricCoin.CoinName[0:1] + equationJson.TenantID[0:1] +  replacerTwo
			} else if(remainingChars == 1){
				generatedCoinName = equationJson.MetricCoin.CoinName[0:1] +  replacerTwo
			} else{
				generatedCoinName = replacerTwo
			}
		} else{
			generatedCoinName = replacerTwo
		}
		
		logrus.Info("Generated and modified Coin name  02 : ", generatedCoinName)

		for i := 0; i < len(portion); i++ {
			if len(portion[i].FieldAndCoin) > 0 {
				for j := 0; j < len(portion[i].FieldAndCoin); j++ {
					if portion[i].FieldAndCoin[j].VariableType == "operator" {
						portion[i].FieldAndCoin = append(portion[i].FieldAndCoin[:j], portion[i].FieldAndCoin[j+1:]...)
					}
					if j == len(portion[i].FieldAndCoin)-1 {
						// timestamp := makeTimestamp()
						// strTimestamp := strconv.Itoa(int(timestamp))
						// generatedCoinName := equationJson.MetricCoin.CoinName[0:1] +
						// 	equationJson.EquationID[0:3] + equationJson.ProductName[0:1] + equationJson.TenantID[0:4] + strTimestamp[10:13]
						
						portion[i].FieldAndCoin[j].CoinName = equationJson.MetricCoin.CoinName
						portion[i].FieldAndCoin[j].GeneratedName = generatedCoinName
						
						coinMap1 := model.CoinMap{
							CoinName:      equationJson.MetricCoin.CoinName,
							GeneratedName: generatedCoinName,
						}
						coinMap = append(coinMap, coinMap1)
					} else {
						coinMap1 := model.CoinMap{
							CoinName:      portion[i].FieldAndCoin[j].CoinName,
							GeneratedName: portion[i].FieldAndCoin[j].GeneratedName,
						}
						coinMap = append(coinMap, coinMap1)
					}
				}
			} else {
				return model.CreatePool{}, []model.CoinMap{}, errors.New("Equation-JSON's Sub portions are empty")
			}
		}
		// Find element in a slice and move it to first position
		for i := 0; i < len(portion); i++ {
			if len(portion[i].FieldAndCoin) > 0 {
				portion[i].FieldAndCoin = rearrangedArray(portion[i].FieldAndCoin, "userInput")
			} else {
				return model.CreatePool{}, []model.CoinMap{}, errors.New("Equation-JSON's Sub portions are empty")
			}
		}
	} else {
		return model.CreatePool{}, []model.CoinMap{}, errors.New("Equation-JSON's Sub portions are empty")
	}
	return equationJson, coinMap, nil
}

// rearrangedArray return the orded array ==> Find element(coinsAndFiled) in a slice by variable type and move it to first position
func rearrangedArray(poolJson []model.FieldAndCoin, find string) []model.FieldAndCoin {
	if len((poolJson)) == 0 || (poolJson)[0].VariableType == find {
		return poolJson
	}
	if (poolJson)[len(poolJson)-1].VariableType == find {
		(poolJson) = append([]model.FieldAndCoin{poolJson[len(poolJson)-1]}, (poolJson)[:len(poolJson)-1]...)
		return poolJson
	}
	for p, x := range poolJson {
		if x.VariableType == find {
			(poolJson) = append([]model.FieldAndCoin{x}, append((poolJson)[:p], (poolJson)[p+1:]...)...)
		}
	}
	return poolJson
}

// CoinConvertionJson ==>this method  rerstructed the coinConvertion request body recived by backend
// metric Coin is used as a received coin because all sub-portions of an equation finally should be the same units
func CoinConvertionJson(coinConvertObject model.BatchCoinConvert, batchAccountPK string, batchAccountSK string) ([]model.BuildPathPayment, error) {
	object := dao.Connection{}
	data, err := object.GetPoolFromDB(coinConvertObject.FormulaType ,coinConvertObject.EquationID, coinConvertObject.ProductName, coinConvertObject.TenantID).Then(func(data interface{}) interface{} {
		logrus.Info("Pool data taken from DB")
		return data
	}).Await()
	if err != nil {
		logrus.Info("Cannot take pool data from DB")
		return []model.BuildPathPayment{}, err
	}
	if data == nil {
		return []model.BuildPathPayment{}, errors.New("Can not find Pool")
	}
	//find  coin name assign it to generated coin name 
	for i := 0; i < len(coinConvertObject.UserInputs); i++ {
		for k := 0; k < len(data.(model.BuildPoolResponse).CoinMap); k++ {
			fmt.Println(data.(model.BuildPoolResponse).CoinMap[k].CoinName, data.(model.BuildPoolResponse).CoinMap[k].GeneratedName)
			if coinConvertObject.UserInputs[i].CoinName == data.(model.BuildPoolResponse).CoinMap[k].CoinName {
				coinConvertObject.UserInputs[i].GeneratedName = data.(model.BuildPoolResponse).CoinMap[k].GeneratedName
			}
			if coinConvertObject.MetricCoin.CoinName == data.(model.BuildPoolResponse).CoinMap[k].CoinName {
				coinConvertObject.MetricCoin.GeneratedName = data.(model.BuildPoolResponse).CoinMap[k].GeneratedName
			}
		}
	}

	var buildPathPayments []model.BuildPathPayment
	if coinConvertObject.UserInputs != nil && len(coinConvertObject.UserInputs) > 0 {
		for _, inputCoin := range coinConvertObject.UserInputs {
			fmt.Println(inputCoin.CoinName)

			buildPathPayment := model.BuildPathPayment{
				SendingCoin: model.Coin{
					Id: inputCoin.Id, CoinName: inputCoin.CoinName,
					Amount: inputCoin.Value, GeneratedName: inputCoin.GeneratedName,
				},
				ReceivingCoin: model.Coin{
					Id: coinConvertObject.MetricCoin.Id, CoinName: coinConvertObject.MetricCoin.CoinName,
					FieldName: coinConvertObject.MetricCoin.FieldName, GeneratedName: coinConvertObject.MetricCoin.GeneratedName,
				},
				BatchAccountPK:     batchAccountPK,
				BatchAccountSK:     batchAccountSK,
				CoinIssuerAccontPK: coinIsuserPK,
				PoolId:             "",
			}
			buildPathPayments = append(buildPathPayments, buildPathPayment)
		}
	} else {
		return []model.BuildPathPayment{}, errors.New("user Input coin values are empty")
	}
	return buildPathPayments, nil
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}


func randomChar() string{
	rand.Seed(time.Now().UnixNano())
	charset := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	c := charset[rand.Intn(len(charset))]
	logrus.Info("Generated character " ,string(c))
	return string(c)
}
