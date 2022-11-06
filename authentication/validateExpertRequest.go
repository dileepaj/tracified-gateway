package authentication

import (
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
	Signature    string
	Plaintext    string
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
	// PGP validator
	if configs.DigitalSIgnatureValidationEnabled {
		errWhenValidatingDigitalSignature, isValidated := PGPValidator(authObject.ExpertPK, authObject.Signature, authObject.Plaintext)
		if errWhenValidatingDigitalSignature != nil {
			logrus.Error("Digital signature validation issue ", errWhenValidatingDigitalSignature)
			return errWhenValidatingDigitalSignature, http.StatusUnauthorized
		}
		if !isValidated {
			logrus.Error("Signature validation failed, incorrect credentials")
			return errors.New("Signature validation failed, incorrect credentials"), http.StatusUnauthorized
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
		RequestEntityType: "EXPERT",
		RequestEntity:     authObject.ExpertPK,
		FormulaID:         authObject.FormulaId,
		AllowedAmount:     allowReqPerDay,
		FromTime:          convertedFromTime,
		ToTime:            convertedFromTime.AddDate(0, 0, +1),
	}
	err, errCode, count := APIThrottler(apiReq, true)
	logrus.Info("  REQUESTPERDAY  ", count)
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
		RequestEntityType: "EXPERT",
		RequestEntity:     authObject.ExpertPK,
		FormulaID:         authObject.FormulaId,
		AllowedAmount:     allowReqPerWeek,
		FromTime:          convertedFromTime.AddDate(0, 0, -6),
		ToTime:            convertedFromTime.AddDate(0, 0, +1),
	}
	err, errCode, count := APIThrottler(apiReq, false)
	logrus.Info("  REQUESTPERWEEK  ", count)
	if err != nil {
		if errCode != 429 {
			return err, errCode
		}
		NotifyToAdmin(authObject.ExpertPK, authObject.ExpertPK)
	}
	return nil, 200
}
