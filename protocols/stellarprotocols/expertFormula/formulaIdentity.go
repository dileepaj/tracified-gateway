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
*/
func (expertFormula ExpertFormula)BuildFormulaIdentity(expertId uint64, formulaName string, formulaDecription string) (txnbuild.ManageData, error) {
	formName := ""
	authorID := stellarprotocols.UInt64ToByteString(expertId)
	formDescription := ""
	//check if the formula name is 15 characters
	if len(formulaName) > 15 {
		logrus.Error("Formula name 15 character limit exceeded")
		return txnbuild.ManageData{}, errors.New("Formula name 15 character limit exceeded")
	} else {
		if len(formulaName) == 15 {
			formName = formulaName
		} else if len(formulaName) < 15 {
			formName = formulaName + "/"
		}
	}
	//checking the formula name length if not 15 append 0s to the end
	if len(formName) < 15 {
		remain := 15 - len(formName)
		setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		formName = formName + setReaminder
	}
	//checking if the expert ID is having 8 characters
	if len(authorID) < 8 {
		remain := 8 - len(authorID)
		setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		authorID = authorID + setReaminder
	}
	//check if the formula description have 64 characters
	if len(formulaDecription) > 64 {
		logrus.Error("Formula description 64 character limit exceeded")
		return txnbuild.ManageData{}, errors.New("Formula description 64 character limit exceeded")
	} else {
		if len(formulaDecription) == 64 {
			formDescription = formulaDecription
		} else if len(formulaDecription) < 64 {
			formDescription = formulaDecription + "/"
		}
	}
	//checking if the formula description has the 64 bytes
	if len(formDescription) < 64 {
		remain := 64 - len(formDescription)
		setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		formDescription = formDescription + setReaminder
	}
	// define a 41 zeros string for future use
	decodedStrFutureUse, err := hex.DecodeString(fmt.Sprintf("%082d", 0))
	if err != nil {
		return txnbuild.ManageData{}, err
	}
	strFutureUse := string(decodedStrFutureUse)

	valueString := formName + authorID + strFutureUse 
	keyString := formDescription
	
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
		return txnbuild.ManageData{}, errors.New("Length issue on key or value fields on the formula identity building")
	}
	return formulaIdentityBuilder, nil
}
