package expertformula

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/dileepaj/tracified-gateway/configs"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

/*
BuildAuthorManageData
des-Build the author definition manage data
	public key - 64 bytes
	future use - 64 bytes
*/
func (expertFormula ExpertFormula) BuildAuthorManageData(expertKey string) (txnbuild.ManageData, txnbuild.ManageData, error) {
	authorKey1 := ""
	authorValue1 := ""
	authorKey2 := ""
	authorValue2 := ""

	// check if the string is 128 characters
	if len(expertKey) != 256 {
		logrus.Error("Expert public key should be equal to 256 bytes")
		return txnbuild.ManageData{}, txnbuild.ManageData{}, errors.New("expert public key should be equal to 256 bytes")
	} else {
		// check if the expert key is greater than 64 character limit
		if len(expertKey) > 64 {
			// divide the expert key to 2 parts with each of 64 bytes
			authorKey1 = expertKey[0:64]
			authorValue1 = expertKey[64:128]
			authorKey2 = expertKey[128:192]
			authorValue2 = expertKey[192:]
		}
	}
	authorBuilder1 := txnbuild.ManageData{
		Name:  authorKey1,
		Value: []byte(authorValue1),
	}
	authorBuilder2 := txnbuild.ManageData{
		Name:  authorKey2,
		Value: []byte(authorValue2),
	}

	// check the lengths of the key and value
	if len(authorKey1) > 64 || len(authorValue1) > 64 {
		logrus.Error("Key string length : ", len(authorKey1))
		logrus.Error("Value string length : ", len(authorValue1))
		return txnbuild.ManageData{}, txnbuild.ManageData{}, errors.New("length issue on key or value fields on the author details building")
	}
	if len(authorKey2) > 64 || len(authorValue2) > 64 {
		logrus.Error("Key string length : ", len(authorKey2))
		logrus.Error("Value string length : ", len(authorValue2))
		return txnbuild.ManageData{}, txnbuild.ManageData{}, errors.New("length issue on key or value fields on the author details building")
	}

	return authorBuilder1, authorBuilder2, nil
}

func (expertFormula ExpertFormula) BuildPublicManageData(publicKeyHash string) (txnbuild.ManageData, error) {
	// check if the string is 64 characters
	if configs.PGPkeyEnable && len(publicKeyHash) != 64 {
		logrus.Error("Expert public key should be equal to 64 character limit, It is a sha256(authorDetailsBuilder.go)")
		return txnbuild.ManageData{}, errors.New("expert public key should be equal to 64 character limit, It is a sha256 value")
	}
	if !configs.PGPkeyEnable && len(publicKeyHash) > 64 {
		logrus.Error("Expert public key should be less than 64 character limit, It is a stellar public key(authorDetailsBuilder.go)")
		return txnbuild.ManageData{}, errors.New("expert public key should be less than 64 character limit, It is a stellar public key")
	}
	decodedStrFutureUse, err := hex.DecodeString(fmt.Sprintf("%0128d", 0))
	if err != nil {
		return txnbuild.ManageData{}, err
	}
	authorBuilder := txnbuild.ManageData{
		Name:  publicKeyHash,
		Value: decodedStrFutureUse,
	}

	// check the lengths of the key and value
	if len(publicKeyHash) > 64 || len(decodedStrFutureUse) != 64 {
		logrus.Error("Key string length : ", len(publicKeyHash))
		logrus.Error("Value string length : ", len(decodedStrFutureUse))
		return txnbuild.ManageData{}, errors.New("length issue on key or value fields on the author details building")
	}
	logrus.Info("Author detials key ", publicKeyHash)
	logrus.Info("Author details value ", decodedStrFutureUse)
	return authorBuilder, nil
}
