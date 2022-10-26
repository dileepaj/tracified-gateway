package expertformula

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

/*
BuildFormulaIdentity
des-Build the formula idenitiy manage data
	formula name - 64 (full key)
	author id -  8 (value)
	future use - 56 (value)
*/
func (expertFormula ExpertFormula) BuildFormulaIdentity(expertId uint64, formulaName string) (txnbuild.ManageData, error) {
	authorID := stellarprotocols.UInt64ToByteString(expertId)
	formName := ""

	//checking if the expert ID is having 8 characters
	if len(authorID) < 8 {
		remain := 8 - len(authorID)
		setReaminder := strings.Repeat("0", remain)
		authorID = authorID + setReaminder
	}
	//check if the formula name have 64 characters
	if len(formulaName) > 64 {
		logrus.Error("Formula name 64 character limit exceeded(formulaIdentity.go)")
		return txnbuild.ManageData{}, errors.New("formula name 64 character limit exceeded")
	} else {
		if len(formulaName) == 64 {
			formName = formulaName
		} else if len(formulaName) < 64 {
			formName = formulaName + "/"
		}
	}
	//checking if the formula name has the 64 bytes
	if len(formName) < 64 {
		remain := 64 - len(formName)
		setReaminder := strings.Repeat("0", remain)
		formName = formName + setReaminder
	}
	// define a 41 zeros string for future use
	decodedStrFutureUse, err := hex.DecodeString(fmt.Sprintf("%0112d", 0))
	if err != nil {
		return txnbuild.ManageData{}, err
	}
	strFutureUse := string(decodedStrFutureUse)

	valueString := authorID + strFutureUse
	keyString := formName

	logrus.Info("Formula identity key ", keyString)
	logrus.Info("Formula identity value ", valueString)
	//building the manage data operation
	formulaIdentityBuilder := txnbuild.ManageData{
		Name:  keyString,
		Value: []byte(valueString),
	}
	//check the lengths of the key and value
	if len(keyString) > 64 || len(valueString) > 64 {
		logrus.Error("Key string length : ", len(keyString))
		logrus.Error("Value string length : ", len(valueString))
		return txnbuild.ManageData{}, errors.New("length issue on key or value fields on the formula identity building")
	}
	return formulaIdentityBuilder, nil
}
