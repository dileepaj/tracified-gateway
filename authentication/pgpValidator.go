package authentication

import (
	"encoding/base64"
	"errors"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

func PGPValidator(sha256hash, signature, originalMsg string) (error, bool) {
	object := dao.Connection{}

	// get the public key from the DB using its hash
	publicKeyDet, errWhenGettingPublicKey := object.GetRSAPublicKeyBySHA256PK(sha256hash).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errWhenGettingPublicKey != nil {
		logrus.Error("Error when getting the public key from gateway datastore " + errWhenGettingPublicKey.Error())
		return errors.New("Error when getting the public key from gateway datastore " + errWhenGettingPublicKey.Error()), false
	}
	if publicKeyDet == nil {
		logrus.Error("Public key does not exist in the gateway datastore for the hash ", sha256hash)
		return errors.New("Public key does not exist in the gateway datastore for the hash " + sha256hash), false
	}
	publicKeyData := publicKeyDet.(model.RSAPublickey)
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKeyData.PgpPublickkey)
	if err != nil {
		logrus.Error("Open PGP, public key decoding error  ", err)
		return errors.New("Open PGP, public key decoding error  " + err.Error()), false
	}
	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		logrus.Error("Open PGP, signature decoding error  ", err)
		return errors.New("Open PGP, signature decoding error  " + err.Error()), false
	}
	verifiedPlainText, err := helper.VerifyCleartextMessageArmored(string(decodedPublicKey), string(decodedSignature), crypto.GetUnixTime())
	if err != nil {
		return err, false
	}
	if verifiedPlainText == originalMsg {
		return nil, true
	} else {
		return errors.New("signature txt is not matching"), false
	}
}
