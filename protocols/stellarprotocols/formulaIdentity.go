package stellarprotocols

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

func BuildFormulaIdentity(expertId int, formulaName string, formulaDecription string) (txnbuild.ManageData, error) {

	formName := ""
	authorID := ""
	formDescription := ""

	//check if the formula name is 15 characters
	if len(formulaName) > 15 {
		logrus.Error("Formula name 15 character limit exceeded")
		return txnbuild.ManageData{}, errors.New("Formula name 15 character limit exceeded")
	} else {
		if len(formulaName) < 15 {
			//add 0s to the rest of the name
			remain := 15 - len(formulaName)
			setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain))
			formName = formulaName + setReaminder
		} else if len(formulaName) == 15 {
			formName = formulaName
		}
	}

	//convert the formula id to string
	expIdString := strconv.Itoa(expertId)

	//check if the expert ID have 8 characters
	if len(expIdString) > 8 {
		logrus.Error("Expert ID 8 character limit exceeded")
		return txnbuild.ManageData{}, errors.New("Expert ID 8 character limit exceeded")
	} else {
		if len(expIdString) < 8 {
			//add 0s to the start
			remain := 8 - len(expIdString)
			setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain))
			authorID = setReaminder + expIdString
		} else if len(expIdString) == 8 {
			authorID = expIdString
		}
	}

	//check if the formula description have 64 characters
	if len(formulaDecription) > 64 {
		logrus.Error("Formula description 64 character limit exceeded")
		return txnbuild.ManageData{}, errors.New("Formula description 64 character limit exceeded")

	} else {
		if len(formulaDecription) < 64 {
			//add 0s to the rest of the name
			remain := 64 - len(formulaDecription)
			setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain))
			formDescription = formulaDecription + setReaminder
		} else if len(formulaDecription) == 64 {
			formDescription = formulaDecription
		}
	}

	keyString := formName + authorID
	valueString := formDescription

	logrus.Info("Formula identity key ", keyString)
	logrus.Info("Formula identity value ", valueString)

	//building the manage data operation
	formulaIdentityBuilder := txnbuild.ManageData{
		Name:  keyString,
		Value: []byte(valueString),
	}

	return formulaIdentityBuilder, nil

}
