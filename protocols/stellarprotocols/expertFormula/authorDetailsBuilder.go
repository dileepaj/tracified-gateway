package expertformula

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

/*
BuildAuthorManageData
des-Build the author definition manage data
*/
func (expertFormula ExpertFormula) BuildAuthorManageData(expertKey string) (txnbuild.ManageData, txnbuild.ManageData, error) {
	authorKey1 := ""
	authorValue1 := ""
	authorKey2 := ""
	authorValue2 := ""

	// check if the string is 128 characters
	if len(expertKey) != 256 {
		logrus.Error("Expert public key should be equal to 256 bytes")
		return txnbuild.ManageData{}, txnbuild.ManageData{}, errors.New("Expert public key should be equal to 256 bytes")
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
		return txnbuild.ManageData{}, txnbuild.ManageData{}, errors.New("Length issue on key or value fields on the author details building")
	}
	if len(authorKey2) > 64 || len(authorValue2) > 64 {
		logrus.Error("Key string length : ", len(authorKey2))
		logrus.Error("Value string length : ", len(authorValue2))
		return txnbuild.ManageData{}, txnbuild.ManageData{}, errors.New("Length issue on key or value fields on the author details building")
	}

	return authorBuilder1, authorBuilder2, nil
}

func (expertFormula ExpertFormula) BuildPublisherManageData(publicKeyHash string) (txnbuild.ManageData, error) {
	// check if the string is 128 characters
	if len(publicKeyHash) != 64 {
		logrus.Error("Expert public key is equal to 64 character limit, It is a sha256")
		return txnbuild.ManageData{}, errors.New("Expert public key is equal to 64 character limit, It is a sha256 value")
	}
	decodedStrFetureUsed, err := hex.DecodeString(fmt.Sprintf("%0128d", 0))
	if err != nil {
		return txnbuild.ManageData{}, err
	}
	authorBuilder := txnbuild.ManageData{
		Name:  publicKeyHash,
		Value: decodedStrFetureUsed,
	}

	// check the lengths of the key and value
	if len(publicKeyHash) != 64 || len(decodedStrFetureUsed) != 64 {
		logrus.Error("Key string length : ", len(publicKeyHash))
		logrus.Error("Value string length : ", len(decodedStrFetureUsed))
		return txnbuild.ManageData{}, errors.New("Length issue on key or value fields on the author details building")
	}
	logrus.Info("Author detials key ", publicKeyHash)
	logrus.Info("Author details value ", decodedStrFetureUsed)
	return authorBuilder, nil
}
