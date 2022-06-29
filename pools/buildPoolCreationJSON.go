package pools

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/model"
)

// var poolJson []model.BuildPool
// pool1 := model.BuildPool{
// 	Coin1:               "M",
// 	DepositeAmountCoin1: "10000",
// 	Coin2:               "N",
// 	DepositeAmountCoin2: "30000",
// 	Ratio:               2,
// }
// pool2 := model.BuildPool{
// 	Coin1:               "N",
// 	DepositeAmountCoin1: "10000",
// 	Coin2:               "O",
// 	DepositeAmountCoin2: "70000",
// 	Ratio:               2,
// }
// pool3 := model.BuildPool{
// 	Coin1:               "O",
// 	DepositeAmountCoin1: "10000",
// 	Coin2:               "P",
// 	DepositeAmountCoin2: "90000",
// 	Ratio:               2,
// }

// BuildPoolCreationJSON return the pool creation json by Restructing the admin equationJson
func BuildPoolCreationJSON(equationJson model.CreatePool) ([]model.BuildPool, error) {
	var poolJson []model.BuildPool
	// var a []int
	portion := equationJson.EquationSubPortion
	for i := 0; i < len(portion); i++ {
		for j := 0; j < len(portion[i].FieldAndCoin); j++ {
			if j+1 < len(portion[i].FieldAndCoin) {
				//
				ratio := portion[i].FieldAndCoin[j+1].Value
				// a = append(a, j, j+1)
				removeDecimal := strings.Split(ratio, ".")

				multiplicationZeros := len(removeDecimal[0]) - 1
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
				for j := 0; j < len(portion[i].FieldAndCoin); j++ {
					if portion[i].FieldAndCoin[j].VariableType == "operator" && portion[i].FieldAndCoin[j].Value == "/" {
						portion[i].FieldAndCoin[j].Value = "*"
						if float, err := strconv.ParseFloat(portion[i].FieldAndCoin[j+1].Value, 32); err == nil {
							var value float64 = 1 / float
							portion[i].FieldAndCoin[j+1].Value = fmt.Sprintf("%f", value)
						}
					}
				}
			} else {
				return model.CreatePool{}, errors.New("Equation-JSON's Sub portions are empty")
			}
		}

		// reomve the oprator from equation
		for i := 0; i < len(portion); i++ {
			for j := 0; j < len(portion[i].FieldAndCoin); j++ {
				if portion[i].FieldAndCoin[j].VariableType == "operator" {
					portion[i].FieldAndCoin = append(portion[i].FieldAndCoin[:j], portion[i].FieldAndCoin[j+1:]...)
				}
			}
		}
	} else {
		return model.CreatePool{}, errors.New("Equation-JSON's Sub portions are empty")
	}
	return equationJson, nil
}
