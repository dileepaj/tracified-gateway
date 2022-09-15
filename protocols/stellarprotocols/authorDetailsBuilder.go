package stellarprotocols

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

func BuildAuthorManageData(expertKey string) (txnbuild.ManageData, error) {

	authorKey := ""
	authorValue := ""

	//check if the string is 128 characters
	if len(expertKey) > 128 {
		logrus.Error("Expert public key is greater than 128 character limit")
		return txnbuild.ManageData{}, errors.New("Expert public key is greater than 128 character limit")
	} else {
		//check if the expert key is greater than 64 character limit
		if len(expertKey) > 64 {
			//divide the expert key to 2 parts with each of 64 bytes
			authorKey = expertKey[0:64]
			authorValue = expertKey[64:]

		} else if len(expertKey) < 64 || len(expertKey) == 64 {
			//add to key field directly
			authorKey = expertKey
			authorValue = fmt.Sprintf("%s", strings.Repeat("0", 64))
		}

	}

	logrus.Info("Author detials key ", authorKey)
	logrus.Info("Author details value ", authorValue)

	authorBuilder := txnbuild.ManageData{
		Name:  authorKey,
		Value: []byte(authorValue),
	}

	return authorBuilder, nil
}
