package expertformula

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

/*
des - build the value manage data for the semantic constant
return txnbuild.Manage data object

	value - 64 bytes
	future use - 64 bytes
*/

func (expertFormula ExpertFormula) BuildSemanticValue(value any) (txnbuild.ManageData, error) {

	semanticConstantValue := fmt.Sprintf("%g", value.(float64))
	// check value is 64 character
	if len(semanticConstantValue) > 64 {
		logrus.Error("Value is greater than 64 character limit(semanticValueBuilder.go) ")
		return txnbuild.ManageData{}, errors.New("value is greater than 64 character limit ")
	} else {
		if len(semanticConstantValue) < 64 {
			// add 0s to the rest of the name
			remain := 64 - len(semanticConstantValue)
			setReaminder := strings.Repeat("0", remain)
			semanticConstantValue = setReaminder + semanticConstantValue
		}
	}
	decodedStrFetureUsed, err := hex.DecodeString(fmt.Sprintf("%0128d", 0))
	if err != nil {
		return txnbuild.ManageData{}, err
	}
	strFetureUsed := string(decodedStrFetureUsed)

	nameString := semanticConstantValue
	valueString := strFetureUsed

	logrus.Info("Semantic Key : ", nameString)
	logrus.Info("Semantic Value : ", valueString)

	if len(nameString) > 64 || len(valueString) > 64 {
		logrus.Error("Key or Value on semantic value builder exceeded 64 character limit(semanticValueBuilder.go) ")
		return txnbuild.ManageData{}, errors.New("key or Value on semantic value builder exceeded 64 character limit ")
	}

	semanticValueBuilder := txnbuild.ManageData{
		Name:  nameString,
		Value: []byte(valueString),
	}

	return semanticValueBuilder, nil
}
