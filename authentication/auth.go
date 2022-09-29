package authentication

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"strings"
	"time"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

type AuthLayer struct {
	FormulaId string
	ExpertPK  string
	CiperText string
	Plaintext []model.FormulaItemRequest
}

func (authObject AuthLayer) ValidateExpertRequest() (error, int) {
	err1, code1 := authObject.isExceedRequestLimit()
	// if err1 != nil {
	// 	return err1, code1
	// } else {
	// 	err2, code2 := authObject.isSignatureValid()
	// 	if err2 != nil {
	// 		return err2, code2
	// 	}
	return err1, code1
	//}
}

func (authObject AuthLayer) isExceedRequestLimit() (error, int) {
	t := time.Now()
	time.Local = time.UTC
	convertedFromTime := time.Date(t.Year(), t.Month(), t.Day(), 0o0, 0o0, 0o0, 0o0, t.UTC().Location())
	apiReq := model.API_ThrottlerRequest{
		RequestEntityType: "Test",
		RequestEntity:     "PK",
		FormulaID:         authObject.FormulaId,
		AllowedAmount:     5,
		FromTime:          convertedFromTime,
		ToTime:            convertedFromTime.AddDate(0, 0, +1),
	}
	err, errCode, _ := APIThrottler(apiReq)
	if err != nil {
		return err, errCode
	}
	return nil, 200
}

func (authObject AuthLayer) isSignatureValid() (error, int) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		logrus.Error(err)
	}
	publicKey := privateKey.PublicKey
	secretMessage := "This is super secret message!"
	signature := CreateSignature(secretMessage, *privateKey)

	// Export the keys to pem string
	privateKeyPemString, _ := ExportRsaPrivateKeyAsPemStr(privateKey)
	publicKeyPemString, _ := ExportRsaPublicKeyAsPemStr(&privateKey.PublicKey)

	// Import the keys from pem string
	priv_parsed, _ := ParseRsaPrivateKeyFromPemStr(privateKeyPemString)
	pub_parsed, _ := ParseRsaPublicKeyFromPemStr(publicKeyPemString)

	// Export the newly imported keys
	priv_parsed_pem, _ := ExportRsaPrivateKeyAsPemStr(priv_parsed)
	pub_parsed_pem, _ := ExportRsaPublicKeyAsPemStr(pub_parsed)

	startRemovedPK := strings.ReplaceAll(pub_parsed_pem, "-----BEGIN RSA PUBLIC KEY-----\n", "")
	startRemovedSK := strings.ReplaceAll(priv_parsed_pem, "-----BEGIN RSA PRIVATE KEY-----\n", "")
	endRemovedPK := strings.ReplaceAll(startRemovedPK, "-----END RSA PRIVATE KEY-----\n", "")
	endRemovedSK := strings.ReplaceAll(startRemovedSK, "-----END RSA PUBLIC KEY-----\n", "")

	logrus.Info("Private key ", endRemovedSK)
	logrus.Info("Public key ", endRemovedPK)

	isSignatureValied := VerifySignature(secretMessage, signature, publicKey)
	logrus.Info("isSignatureValied ", isSignatureValied)
	if !isSignatureValied {
		return errors.New("Digital signature verification issue"), 403
	}
	return nil, 200
}
