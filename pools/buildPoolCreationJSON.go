package pools

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

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
				object := dao.Connection{}
				data, _ := object.GetPool(pool1.Coin1, pool1.Coin2).Then(func(data interface{}) interface{} {
					return data
				}).Await()
				if data != nil {
					logrus.Error("pool already deposited " + " Coin2 " + pool1.Coin1 + " Coin2 " + pool1.Coin2)
					return []model.BuildPool{}, errors.New("pool is already deposited " + " Coin2 " + pool1.Coin1 + " Coin2 " + pool1.Coin2)
				}

				pool2 := model.Pool{
					EquationId:          equationJson.EquationID,
					ProductId:           equationJson.ProductID,
					TenantId:            equationJson.TenantID,
					FormulatType:        equationJson.FormulaType,
					Coin1:               portion[i].FieldAndCoin[j].GeneratedName,
					DepositeAmountCoin1: strconv.Itoa(coinDetails1),
					Coin2:               portion[i].FieldAndCoin[j+1].GeneratedName,
					DepositeAmountCoin2: strconv.Itoa(coinDetails2),
					Ratio:               ratio,
				}
				err := object.InsertPoool(pool2)
				if err != nil {
					logrus.Error("Pool did not add to DB ", err)
					return nil, err
				}
				logrus.Info("pool deposite" + " Coin2 " + pool1.Coin1 + " Coin2 " + pool1.Coin2)
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
	portion := equationJson.EquationSubPortion
	var coinMap []model.CoinMap
	if len(portion) > 0 {
		for i := 0; i < len(portion); i++ {
			if len(portion[i].FieldAndCoin) > 2 {
				userInputCount := 0
				for j := 0; j < len(portion[i].FieldAndCoin); j++ {
					if portion[i].FieldAndCoin[j].Value=="(" || portion[i].FieldAndCoin[j].Value==")" {
						return model.CreatePool{}, []model.CoinMap{}, errors.New("Equation can not caontaion  ( , ) oprators")
					}
					// Check if the coin name's character equalto 4
					if portion[i].FieldAndCoin[j].VariableType != "OPERATOR" && len(portion[i].FieldAndCoin[j].CoinName) != 4 {
						return model.CreatePool{}, []model.CoinMap{}, errors.New("Coin name character limit should be 4")
					}
					if portion[i].FieldAndCoin[j].VariableType != "OPERATOR" || portion[i].FieldAndCoin[j].CoinName != "" {
						// timestamp := makeTimestamp()
						generatedSID, err := GenerateCoinName(equationJson.TenantID, portion[i].FieldAndCoin[j].CoinName, portion[i].FieldAndCoin[j].Description,
							portion[i].FieldAndCoin[j].VariableType, equationJson.EquationID, portion[i].FieldAndCoin[j].FieldName)
						if err != nil {
							logrus.Error("Cannot generate ShortID", err.Error())
							return model.CreatePool{}, []model.CoinMap{}, err
						}
						logrus.Info("Generated and modified Coin name  01 : ", generatedSID)
						portion[i].FieldAndCoin[j].GeneratedName = generatedSID
					}

					// count the userIput type variable in a  sub portion
					if portion[i].FieldAndCoin[j].VariableType == "USERINPUT" {
						userInputCount++
					}
					if portion[i].FieldAndCoin[j].VariableType == "OPERATOR" && portion[i].FieldAndCoin[j].Value == "/" {
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

		for i := 0; i < len(portion); i++ {
			if len(portion[i].FieldAndCoin) > 0 {
				for j := 0; j < len(portion[i].FieldAndCoin); j++ {
					if portion[i].FieldAndCoin[j].VariableType == "OPERATOR" {
						portion[i].FieldAndCoin = append(portion[i].FieldAndCoin[:j], portion[i].FieldAndCoin[j+1:]...)
					}
					if j == len(portion[i].FieldAndCoin)-1 {
						generatedSID, err := GenerateCoinName(equationJson.TenantID, equationJson.MetricCoin.CoinName,
							equationJson.MetricCoin.Description, "result", equationJson.EquationID, "")
						if err != nil {
							logrus.Error("Can not generate Coin name", err.Error())
							return model.CreatePool{}, []model.CoinMap{}, err
						}
						logrus.Info("Generated string ", generatedSID)
						portion[i].FieldAndCoin[j].CoinName = strings.ToUpper(equationJson.MetricCoin.CoinName)
						portion[i].FieldAndCoin[j].GeneratedName = generatedSID
						equationJson.MetricCoin.GeneratedName = generatedSID

						coinMap1 := model.CoinMap{
							CoinName:      strings.ToUpper(equationJson.MetricCoin.CoinName),
							GeneratedName: generatedSID,
						}
						coinMap = append(coinMap, coinMap1)
					} else {
						coinMap1 := model.CoinMap{
							CoinName:      strings.ToUpper(portion[i].FieldAndCoin[j].CoinName),
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
				portion[i].FieldAndCoin = rearrangedArray(portion[i].FieldAndCoin, "USERINPUT")
			} else {
				return model.CreatePool{}, []model.CoinMap{}, errors.New("Equation-JSON's Sub portions are empty")
			}
		}
	} else {
		return model.CreatePool{}, []model.CoinMap{}, errors.New("Equation-JSON's Sub portions are empty")
	}
	equationJson.MetricCoin.CoinName = strings.ToUpper(equationJson.MetricCoin.CoinName)
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
	data, _ := object.GetLiquidityPool(coinConvertObject.EquationID, coinConvertObject.ProductName,
		coinConvertObject.TenantID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if data == nil {
		return []model.BuildPathPayment{}, errors.New("Can not find Pool")
	}
	// find  coin name assign it to generated coin name
	for i := 0; i < len(coinConvertObject.UserInputs); i++ {
		for k := 0; k < len(data.(model.BuildPoolResponse).CoinMap); k++ {
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
			buildPathPayment := model.BuildPathPayment{
				SendingCoin: model.Coin{
					Id: inputCoin.Id, CoinName: inputCoin.CoinName,
					Amount: inputCoin.Value, GeneratedName: inputCoin.GeneratedName,
				},
				ReceivingCoin: model.Coin{
					CoinName:  coinConvertObject.MetricCoin.CoinName,
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

// GenerateCoinName return a generated coin name
// Coin name structure
//! generated coin include 12 characters
//! start	 characters 	0-1 "TF" -->  2 characters Tracifed
//! 		 characters 	2-4 ""   --> first 3 character of tenant Id  -upper case
//! 		 characters 	5-7 ""   --> first 3 character of user insertd coin name -upper case
//! 		 characters 	8-12 ""   -->  4 character  straing from 0001 (if coin description not equal ,auto increment count)
func GenerateCoinName(tenantID, coinName, description, coinType, fieldName, equationId string) (string, error) {
	var generatedCoinName string
	object := dao.Connection{}
	data, _ := object.GetCoinName(strings.ToUpper(coinName)).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	coinNameObj := model.CoinName{
		TenantID:    tenantID,
		EquationID:  equationId,
		FieldName:   fieldName,
		Type:        coinType,
		CoinName:    strings.ToUpper(coinName),
		Description: description,
		Count:       "00001",
	}
	if data != nil {
		if strings.Replace(strings.ToLower(data.(model.CoinName).Description), " ", "", -1) == strings.Replace(strings.ToLower(description), " ", "", -1) {
			generatedCoinName = data.(model.CoinName).GeneratedCoinName
		} else {
			// string to int
			i, err := strconv.Atoi(data.(model.CoinName).Count)
			if err != nil {
				return "", err
			}
			i++
			count := strconv.Itoa(i)
			if len(strconv.Itoa(i)) < 5 {
				zero := `%0` + `5` + `d`
				count = fmt.Sprintf(zero, i)
			}
			generatedCoinName = "TFD" + strings.ToUpper(coinName[0:4]) + count
			coinNameObj.GeneratedCoinName = generatedCoinName
			err1 := object.InsertCoinName(coinNameObj)
			if err1 != nil {
				return "", err1
			}
		}
	} else {
		generatedCoinName = "TFD" + strings.ToUpper(coinName[0:4]) + "00001"
		coinNameObj.GeneratedCoinName = generatedCoinName
		err := object.InsertCoinName(coinNameObj)
		if err != nil {
			return "", err
		}
	}
	logrus.Info("coinname: ", coinName, "generatedCoinName  ", generatedCoinName)
	if len(generatedCoinName) != 12 {
		logrus.Error("coinname: ", coinName, "generatedCoinName: ", generatedCoinName)
		return "", errors.New("length issue in generated coin name ")
	}
	return generatedCoinName, nil
}
