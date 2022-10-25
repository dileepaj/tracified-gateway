package authentication

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/configs"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
)

type AuthLayer struct {
	FormulaId    string
	ExpertPK     string
	ExpertUserID string
	CiperText    string
	Plaintext    []model.FormulaItemRequest
}

func (authObject AuthLayer) ValidateExpertRequest() (error, int) {
	if configs.ValidateAgaintTrustNetworkExpertFormulaEnable {
		// validation againt trust network
		errInTrustNetworkValidation := ValidateAgainstTrustNetwork(authObject.ExpertPK)
		if errInTrustNetworkValidation != nil {
			logrus.Error("Expert is not in the trust network ", errInTrustNetworkValidation)
			return errInTrustNetworkValidation, http.StatusBadRequest
		}
	}

	err1, code1 := authObject.isExceedRequestLimitPerDay()
	if err1 != nil {
		return err1, code1
	} else {
		err1, code1 := authObject.isExceedRequestLimitPerWeek()
		if err1 != nil {
			return err1, code1
		}
		return err1, code1
	}
}

func (authObject AuthLayer) isExceedRequestLimitPerDay() (error, int) {
	t := time.Now()
	time.Local = time.UTC
	convertedFromTime := time.Date(t.Year(), t.Month(), t.Day(), 0o0, 0o0, 0o0, 0o0, t.UTC().Location())
	allowReqPerDay, err := strconv.Atoi(commons.GoDotEnvVariable("ALLOWREQUESTPERDAY"))
	if err != nil {
		logrus.Error("Issue when converting string to int  ", err)
	}
	apiReq := model.API_ThrottlerRequest{
		RequestEntityType: "Test",
		RequestEntity:     "PK",
		FormulaID:         authObject.FormulaId,
		AllowedAmount:     allowReqPerDay,
		FromTime:          convertedFromTime,
		ToTime:            convertedFromTime.AddDate(0, 0, +1),
	}
	err, errCode, _ := APIThrottler(apiReq)
	if err != nil {
		return err, errCode
	}
	return nil, 200
}

func (authObject AuthLayer) isExceedRequestLimitPerWeek() (error, int) {
	t := time.Now()
	time.Local = time.UTC
	convertedFromTime := time.Date(t.Year(), t.Month(), t.Day(), 0o0, 0o0, 0o0, 0o0, t.UTC().Location())
	allowReqPerWeek, err := strconv.Atoi(commons.GoDotEnvVariable("ALLOWREQUESTPERWEEK"))
	if err != nil {
		logrus.Error("Issue when converting string to int  ", err)
	}
	apiReq := model.API_ThrottlerRequest{
		RequestEntityType: "Test",
		RequestEntity:     "PK",
		FormulaID:         authObject.FormulaId,
		AllowedAmount:     allowReqPerWeek,
		FromTime:          convertedFromTime.AddDate(0, 0, -6),
		ToTime:            convertedFromTime.AddDate(0, 0, +1),
	}
	err, errCode, _ := APIThrottler(apiReq)
	if err != nil {
		if errCode != 429 {
			return err, errCode
		}
		NotifyToAdmin(authObject.ExpertPK, authObject.ExpertPK)
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
	// // Export the keys to pem string
	// privateKeyPemString, _ := ExportRsaPrivateKeyAsPemStr(privateKey)
	// publicKeyPemString, _ := ExportRsaPublicKeyAsPemStr(&privateKey.PublicKey)
	// // Import the keys from pem string
	// priv_parsed, _ := ParseRsaPrivateKeyFromPemStr(privateKeyPemString)
	// pub_parsed, _ := ParseRsaPublicKeyFromPemStr(publicKeyPemString)
	// // Export the newly imported keys
	// priv_parsed_pem, _ := ExportRsaPrivateKeyAsPemStr(priv_parsed)
	// pub_parsed_pem, _ := ExportRsaPublicKeyAsPemStr(pub_parsed)
	// startRemovedPK := strings.ReplaceAll(pub_parsed_pem, "-----BEGIN RSA PUBLIC KEY-----\n", "")
	// startRemovedSK := strings.ReplaceAll(priv_parsed_pem, "-----BEGIN RSA PRIVATE KEY-----\n", "")
	// endRemovedPK := strings.ReplaceAll(startRemovedPK, "-----END RSA PRIVATE KEY-----\n", "")
	// endRemovedSK := strings.ReplaceAll(startRemovedSK, "-----END RSA PUBLIC KEY-----\n", "")
	// logrus.Info("Private key ", endRemovedSK)
	// logrus.Info("Public key ", endRemovedPK)
	isSignatureValied := VerifySignature(secretMessage, signature, publicKey)
	logrus.Info("isSignatureValied ", isSignatureValied)
	if !isSignatureValied {
		return errors.New("Digital signature verification issue"), 403
	}
	return nil, 200
}
