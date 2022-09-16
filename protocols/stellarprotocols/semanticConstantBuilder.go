package stellarprotocols

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

func BuildSemanticConstantManageData(description string) (txnbuild.ManageData, error) {

	valueType := "2"
	semanticConstValueID := ""
	semanticConstDataType := ""
	semanticConstDescription := ""
	semanticConstValue := ""

	//check if the description have 40 characters
	if len(description) > 40 {
		logrus.Error("Description 40 character limit exceeded")
		return txnbuild.ManageData{}, errors.New("Description 40 character limit exceeded")

	} else {
		if len(description) < 40 {
			//add 0s to the rest of the name
			remain := 40 - len(description)
			setRemainder := fmt.Sprintf("%s", strings.Repeat("0", remain))
			semanticConstDescription = description + setRemainder
		} else if len(description) == 40 {
			semanticConstDescription = description
		}
	}

	keyString := valueType + semanticConstValueID + semanticConstDataType + semanticConstDescription
	valueString := semanticConstValue

	logrus.Info("Semantic constant key: ", keyString)
	logrus.Info("Semantic constant value: ", valueString)

	// Building the manage data operation
	semanticConstManageData := txnbuild.ManageData{
		Name:  keyString,
		Value: []byte(valueString),
	}

	return semanticConstManageData, nil
}
