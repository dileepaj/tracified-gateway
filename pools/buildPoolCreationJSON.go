package pools

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
					Coin1Name:           portion[i].FieldAndCoin[j].CoinName[0:4],
					Coin1FullName:       portion[i].FieldAndCoin[j].FullCoinName,
					DepositeAmountCoin1: strconv.Itoa(coinDetails1),
					Coin2Name:           portion[i].FieldAndCoin[j+1].CoinName[0:4],
					Coin2FullName:       portion[i].FieldAndCoin[j+1].FullCoinName,
					Coin2:               portion[i].FieldAndCoin[j+1].GeneratedName,
					DepositeAmountCoin2: strconv.Itoa(coinDetails2),
					Ratio:               ratio,
					EquationId:          equationJson.EquationID,
					TenantId:            equationJson.TenantID,
					FormulatType:        equationJson.FormulaType,
					Activity:            equationJson.Activity,
					MetricCoin:          equationJson.MetricCoin,
				}
				object := dao.Connection{}
				data, _ := object.GetPool(pool1.Coin1, pool1.Coin2).Then(func(data interface{}) interface{} {
					return data
				}).Await()
				if data != nil {
					logrus.Error("pool already deposited " + " Coin1 " + pool1.Coin1 + " Coin2 " + pool1.Coin2)
					// return []model.BuildPool{}, errors.New("pool is already deposited " + " Coin2 " + pool1.Coin1 + " Coin2 " + pool1.Coin2)
				}

				pool2 := model.Pool{
					EquationId:          equationJson.EquationID,
					Products:            equationJson.Products,
					TenantId:            equationJson.TenantID,
					FormulatType:        equationJson.FormulaType,
					Coin1:               portion[i].FieldAndCoin[j].GeneratedName,
					Coin1Name:           portion[i].FieldAndCoin[j].CoinName,
					Coin1FullName:       portion[i].FieldAndCoin[j].FullCoinName,
					DepositeAmountCoin1: strconv.Itoa(coinDetails1),
					Coin2:               portion[i].FieldAndCoin[j+1].GeneratedName,
					Coin2FullName:       portion[i].FieldAndCoin[j+1].FullCoinName,
					Coin2Name:           portion[i].FieldAndCoin[j+1].CoinName,
					DepositeAmountCoin2: strconv.Itoa(coinDetails2),
					Ratio:               ratio,
				}
				err := object.InsertPoool(pool2)
				if err != nil {
					logrus.Error("Pool did not add to DB ", err)
					return nil, err
				}
				logrus.Info("pool deposite  " + " Coin2 " + pool1.Coin1 + " Coin2 " + pool1.Coin2)
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
					if portion[i].FieldAndCoin[j].Value == "(" || portion[i].FieldAndCoin[j].Value == ")" {
						return model.CreatePool{}, []model.CoinMap{}, errors.New("Equation can not caontaion  ( , ) oprators")
					}
					// Check if the coin name's character equalto 4
					if portion[i].FieldAndCoin[j].VariableType != "OPERATOR" && len(portion[i].FieldAndCoin[j].CoinName) > 4 {
						return model.CreatePool{}, []model.CoinMap{}, errors.New("Coin name character limit should be 4")
					}
					if portion[i].FieldAndCoin[j].VariableType != "OPERATOR" || portion[i].FieldAndCoin[j].CoinName != "" {
						// timestamp := makeTimestamp()
						if j != len(portion[i].FieldAndCoin)-1 {
							generatedSID, err := GenerateCoinName(equationJson.TenantID, portion[i].FieldAndCoin[j].CoinName, portion[i].FieldAndCoin[j].Description,
								portion[i].FieldAndCoin[j].VariableType, equationJson.EquationID, equationJson.MetricCoin.ID, portion[i].FieldAndCoin[j].FullCoinName, portion[i].FieldAndCoin[j].ID)
							if err != nil {
								logrus.Error("Cannot generate ShortID", err.Error())
								return model.CreatePool{}, []model.CoinMap{}, err
							}
							// logrus.Info("Generated and modified Coin name  01 : ", generatedSID)
							portion[i].FieldAndCoin[j].GeneratedName = generatedSID
						}
					}

					// count the userIput type variable in a  sub portion
					if portion[i].FieldAndCoin[j].VariableType == "DATA" {
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
							equationJson.MetricCoin.Description, "result", equationJson.EquationID, equationJson.MetricCoin.ID,
							equationJson.MetricCoin.FullCoinName, equationJson.MetricCoin.ID)
						if err != nil {
							logrus.Error("Can not generate Coin name ", err.Error())
							return model.CreatePool{}, []model.CoinMap{}, err
						}
						logrus.Info("Generated string ", generatedSID)
						portion[i].FieldAndCoin[j].CoinName = strings.ToUpper(equationJson.MetricCoin.CoinName)
						portion[i].FieldAndCoin[j].GeneratedName = generatedSID
						equationJson.MetricCoin.GeneratedName = generatedSID

						coinMap1 := model.CoinMap{
							ID:            equationJson.MetricCoin.ID,
							CoinName:      strings.ToUpper(equationJson.MetricCoin.CoinName),
							FullCoinName:  equationJson.MetricCoin.FullCoinName,
							GeneratedName: generatedSID,
							Description:   equationJson.MetricCoin.Description,
						}
						coinMap = append(coinMap, coinMap1)
					} else {
						coinMap1 := model.CoinMap{
							ID:            portion[i].FieldAndCoin[j].ID,
							CoinName:      strings.ToUpper(portion[i].FieldAndCoin[j].CoinName),
							FullCoinName:  portion[i].FieldAndCoin[j].FullCoinName,
							GeneratedName: portion[i].FieldAndCoin[j].GeneratedName,
							Description:   portion[i].FieldAndCoin[j].Description,
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
				portion[i].FieldAndCoin = rearrangedArray(portion[i].FieldAndCoin, "DATA")
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
func CoinConvertionJson(batchAccountPK, batchAccountSK, formulaId, activityId string,coinConvertObject model.CoinConvertBody) ([]model.BuildPathPayment, error) {
	object := dao.Connection{}
	data, _ := object.GetLiquidityPoolByProductAndActivity(formulaId,
		coinConvertObject.TenantID, coinConvertObject.Type, activityId, coinConvertObject.Event.Details.StageID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if data == nil {
		return []model.BuildPathPayment{}, errors.New("Can not find Pool")
	}
	// find  coin name assign it to generated coin name
	for i := 0; i < len(coinConvertObject.Inputs); i++ {
		for k := 0; k < len(data.(model.BuildPoolResponse).CoinMap); k++ {
			if coinConvertObject.Inputs[i].CoinName == data.(model.BuildPoolResponse).CoinMap[k].FullCoinName {
				coinConvertObject.Inputs[i].GeneratedName = data.(model.BuildPoolResponse).CoinMap[k].GeneratedName
				coinConvertObject.Inputs[i].ID = data.(model.BuildPoolResponse).CoinMap[k].ID
				coinConvertObject.Inputs[i].Description = data.(model.BuildPoolResponse).CoinMap[k].Description
			}
			if coinConvertObject.Metric.Name == data.(model.BuildPoolResponse).CoinMap[k].FullCoinName {
				coinConvertObject.Metric.GeneratedName = data.(model.BuildPoolResponse).CoinMap[k].GeneratedName
				coinConvertObject.Metric.ID = data.(model.BuildPoolResponse).CoinMap[k].ID
				coinConvertObject.Metric.Description = data.(model.BuildPoolResponse).CoinMap[k].Description
			}
		}
	}

	var buildPathPayments []model.BuildPathPayment
	if coinConvertObject.Inputs != nil && len(coinConvertObject.Inputs) > 0 {
		for _, inputCoin := range coinConvertObject.Inputs {
			buildPathPayment := model.BuildPathPayment{
				SendingCoin: model.Coin{
					ID:              inputCoin.ID,
					FullCoinName:    inputCoin.CoinName,
					CoinName:        strings.ToUpper(inputCoin.CoinName[0:4]),
					GeneratedName:   inputCoin.GeneratedName,
					Description:     inputCoin.Description,
					Amount:          fmt.Sprintf("%f", inputCoin.Input),
					RescaledAmmount: fmt.Sprintf("%f", inputCoin.Input/100),
				},
				ReceivingCoin: model.Coin{
					ID:            coinConvertObject.Metric.ID,
					FullCoinName:  coinConvertObject.Metric.Description,
					CoinName:      CreateCoinnameUsingValue(strings.ToUpper(coinConvertObject.Metric.Description)),
					GeneratedName: coinConvertObject.Metric.GeneratedName,
					Description:   coinConvertObject.Metric.Description,
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
func GenerateCoinName(tenantID, coinName, description, coinType, equationId, metricId, fullCoinName, coinID string) (string, error) {
	var generatedCoinName string
	object := dao.Connection{}
	data, _ := object.GetCoinName(strings.ToUpper(coinName)).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	coinNameObj := model.CoinName{
		TenantID:          tenantID,
		EquationID:        equationId,
		Type:              coinType,
		CoinID:            coinID,
		CoinName:          strings.ToUpper(coinName),
		GeneratedCoinName: generatedCoinName,
		FullCoinName:      fullCoinName,
		Description:       description,
		Count:             "000001",
		MetricID:          metricId,
		Timestamp:         primitive.NewDateTimeFromTime(time.Now()),
	}
	if data != nil {
		// if coin name are equal check the discription
		if strings.Replace(strings.ToLower(data.(model.CoinName).Description), " ", "", -1) != strings.Replace(strings.ToLower(description), " ", "", -1) ||
			strings.Replace(strings.ToLower(data.(model.CoinName).FullCoinName), " ", "", -1) != strings.Replace(strings.ToLower(fullCoinName), " ", "", -1) {
			// if description and coin anme does not equal create new coin name by increting count by 1
			// string to int
			i, err := strconv.Atoi(data.(model.CoinName).Count)
			if err != nil {
				return "", err
			}
			i++
			if i > 999999 {
				return "", errors.New("exceeded the duplication coin limit  " + data.(model.CoinName).FullCoinName + "  " + data.(model.CoinName).CoinName)
			}
			count := strconv.Itoa(i)

			if len(strconv.Itoa(i)) < 5 {
				zero := `%0` + `6` + `d`
				count = fmt.Sprintf(zero, i)
			}
			generatedCoinName = "TF" + strings.ToUpper(coinName[0:4]) + count
			coinNameObj.GeneratedCoinName = generatedCoinName
			coinNameObj.Count = count
			err1 := object.InsertCoinName(coinNameObj)
			if err1 != nil {
				return "", err1
			}
		} else {
			generatedCoinName = data.(model.CoinName).GeneratedCoinName
			coinNameObj := model.CoinName{
				TenantID:          tenantID,
				EquationID:        equationId,
				Type:              data.(model.CoinName).Type,
				CoinID:            coinID,
				CoinName:          data.(model.CoinName).CoinName,
				GeneratedCoinName: data.(model.CoinName).GeneratedCoinName,
				FullCoinName:      data.(model.CoinName).FullCoinName,
				Description:       data.(model.CoinName).Description,
				Count:             data.(model.CoinName).Count,
				MetricID:          metricId,
				Timestamp:         data.(model.CoinName).Timestamp,
			}
			err1 := object.UpdateCoinName(coinNameObj)
			if err1 != nil {
				return "", err1
			}
		}
	} else {
		generatedCoinName = "TF" + strings.ToUpper(coinName[0:4]) + "000001"
		coinNameObj.GeneratedCoinName = generatedCoinName
		err := object.InsertCoinName(coinNameObj)
		if err != nil {
			return "", err
		}
	}
	logrus.Info("coinname: ", coinName, "  generatedCoinName  ", generatedCoinName)
	if len(generatedCoinName) != 12 {
		logrus.Error("coinname: ", coinName, "generatedCoinName: ", generatedCoinName)
		return "", errors.New("length issue in generated coin name ")
	}
	return generatedCoinName, nil
}

// CreateCoinnameUsingValue ==> convert CONSTANT value to 4 character coin name
func CreateCoinnameUsingValue(value string) string {
	// replace . with "D"
	coinName := strings.Replace(value, ".", "D", 1)
	// replace 0 with "Z"
	if len(coinName) < 4 {
		var buffer bytes.Buffer
		count := 4 - len(coinName)
		for i := 1; i <= count; i++ {
			buffer.WriteString("Z")
		}
		return coinName + buffer.String()
	}
	// return first 4 charater
	return coinName[0:4]
}

func CreateCoinnameUsingDescription(str string) string {
	// fruitsString := "apple banana orange pear"
	nonAlphanumericRegex := regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
	wordsString := nonAlphanumericRegex.ReplaceAllString(str, "")
	words := strings.Split(wordsString, " ")
	words2 := []string{}
	coinName := []string{}
	// coinName:=[]string{}
	for _, element := range words {
		word := strings.Replace(element, " ", "", -1)
		word = strings.Replace(word, "	", "", -1)
		if word != "" {
			words2 = append(words2, word)
		}
	}
	if len(words2) >= 4 {
		for _, element := range words2 {
			coinName = append(coinName, element[0:1])
		}
	} else if len(words2) == 3 {
		for i, element := range words2 {
			if i == 0 {
				coinName = append(coinName, element[0:2])
			} else {
				coinName = append(coinName, element[0:1])
			}
		}
	} else if len(words2) == 2 {
		for i, element := range words2 {
			if i == 0 {
				coinName = append(coinName, element[0:2])
				coinName = append(coinName, element[len(element)-1:])
			} else {
				coinName = append(coinName, element[0:1])
			}
		}
	} else {
		coinName = append(coinName, words2[0][0:3])
		coinName = append(coinName, words2[0][len(words2[0])-1:])

	}
	return strings.ToUpper(strings.Join(coinName, "")[0:4])
}
