package pools

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/model"
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
					Coin1:               portion[i].FieldAndCoin[j].CoinName,
					DepositeAmountCoin1: strconv.Itoa(coinDetails1),
					Coin2:               portion[i].FieldAndCoin[j+1].CoinName,
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
func RemoveDivisionAndOperator(equationJson model.CreatePool) (model.CreatePool, error) {
	// equation is a equation-Json
	portion := equationJson.EquationSubPortion

	if len(portion) > 0 {
		for i := 0; i < len(portion); i++ {
			if len(portion[i].FieldAndCoin) > 0 {
				userInputCount := 0
				for j := 0; j < len(portion[i].FieldAndCoin); j++ {
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
					return model.CreatePool{}, errors.New("Equation's sub portion can not contain two user input")
				}
			} else {
				return model.CreatePool{}, errors.New("Equation-JSON's Sub portions are empty")
			}
		}
		// reomve the oprator from equation
		for i := 0; i < len(portion); i++ {
			if len(portion[i].FieldAndCoin) > 0 {
				for j := 0; j < len(portion[i].FieldAndCoin); j++ {
					if portion[i].FieldAndCoin[j].VariableType == "operator" {
						portion[i].FieldAndCoin = append(portion[i].FieldAndCoin[:j], portion[i].FieldAndCoin[j+1:]...)
					}
					if j == len(portion[i].FieldAndCoin)-1 {
						portion[i].FieldAndCoin[j].CoinName = equationJson.MetricCoin.CoinName
					}
				}
			} else {
				return model.CreatePool{}, errors.New("Equation-JSON's Sub portions are empty")
			}
		}
		// Find element in a slice and move it to first position
		for i := 0; i < len(portion); i++ {
			if len(portion[i].FieldAndCoin) > 0 {
				portion[i].FieldAndCoin = rearrangedArray(portion[i].FieldAndCoin, "userInput")
			} else {
				return model.CreatePool{}, errors.New("Equation-JSON's Sub portions are empty")
			}
		}
	} else {
		return model.CreatePool{}, errors.New("Equation-JSON's Sub portions are empty")
	}
	return equationJson, nil
}

// rearrangedArray ,  Find element(coinsAndFiled) in a slice by variable type and move it to first position
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

func BuilPathPaymentJson(coinConvertObject model.BatchCoinConvert, batchAccountPK string, batchAccountSK string) ([]model.BuildPathPayment, error) {
	var buildPathPayments []model.BuildPathPayment
	if coinConvertObject.UserInputs != nil && len(coinConvertObject.UserInputs) > 0 {
		for _, inputCoin := range coinConvertObject.UserInputs {
			buildPathPayment := model.BuildPathPayment{
				SendingCoin:        model.Coin{Id: inputCoin.Id, CoinName: inputCoin.CoinName, Amount: inputCoin.Value},
				ReceivingCoin:      model.Coin{Id: coinConvertObject.MetrixCoin.Id, CoinName: coinConvertObject.MetrixCoin.CointName, FieldName: coinConvertObject.MetrixCoin.FieldName},
				BatchAccountPK:     batchAccountPK,
				BatchAccountSK:     batchAccountSK,
				CoinIssuerAccontPK: coinIsuserPK,
				PoolId:             "",
			}

			buildPathPayments = append(buildPathPayments, buildPathPayment)

		}
	} else {
		return buildPathPayments, errors.New("user Input coin values are empty")
	}
	return buildPathPayments, nil
}
