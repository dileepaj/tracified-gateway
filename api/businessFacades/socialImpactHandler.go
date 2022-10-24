package businessFacades

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/authentication"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols"
	"github.com/dileepaj/tracified-gateway/validations"
	"github.com/sirupsen/logrus"
)

/*
BuildSocialImpactFormula
des-This handler the expert formula bulding,
	In side this handler
	 * validate the JSON request body
	 * changethe formual elemet type according to the protocol
	 * call the SocialImpactExpertFormula method to build the formula
*/

func BuildSocialImpactExpertFormula(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var formulaJSON model.FormulaBuildingRequest
	err := json.NewDecoder(r.Body).Decode(&formulaJSON)
	if err != nil {
		commons.JSONErrorReturn(w, r, err.Error(), http.StatusBadRequest, "Error while decoding the body ")
		return
	}
	errInJsonValidation := validations.ValidateFormulaBuilder(formulaJSON)
	if errInJsonValidation != nil {
		commons.JSONErrorReturn(w, r, errInJsonValidation.Error(), http.StatusBadRequest, "Request body failed the validation check :")
		return
	} else {
		tstMsgEnc := `-----BEGIN PGP MESSAGE----- Version: Keybase OpenPGP v2.0.76 Comment: https://keybase.io/crypto wYwDN3XXL/8v5ksBA/4vRQSu7g/mZovxvzT+XcebWYv8cgtWvK7g4lP7IMopseGl seg6fO3dbfsWiTc7QfVBw8N5IjNSWjrN821etoPDLfujho3aV1FLHtDkFXnOv2Fs VpdXH0p16uSr57hivA1AX4gqrHDDIsFoD2evI01LehrDmb9w1UJ7nX5cMj8et9LA NgHpCi7103et/c+sBJQgQQtEH0YZPi2ddRXwg9aeBZhtpPUHMyWej5oG2crgz/kD tZr3GDkGs05gTcp6Eo5tZ60uI8MQZVkgTbXpFCgUOAwptNSRfvb/ag96kGiTztuh qzPWf9IdW+vEVfy1I/J2vFx2hWI0Q+BzsPixynOIXBfIgcLuDSo/yMc/JhEvwtpV TYvIB04xue2yO3a7cndUVHTXWd14P5tycURUBmTsKFz4gkq1znPm75ug7S9/WARr zf0M+10omwaZMLPKzPDvp5UbCB+XsNTN5xLyzdwY+IetCY0hrBSejzdEVFMXCRm4 mbF3Ct6Vvw== =Ssk2 -----END PGP MESSAGE-----`

		//PGP validation
		errWhenValidatingPGP, isPGPValidationPassed := authentication.PGPValidator("124", tstMsgEnc, "test message")
		if errWhenValidatingPGP != nil {
			logrus.Error("Error when validating PGP keypair ", errWhenValidatingPGP.Error())
			commons.JSONErrorReturn(w, r, errWhenValidatingPGP.Error(), http.StatusBadRequest, "Error when validating PGP keypair ")
			return
		}

		if !isPGPValidationPassed {
			logrus.Error("PGP key is not in the trust network")
			commons.JSONErrorReturn(w, r, "", http.StatusBadRequest, "PGP key is not in the trust network")
			return
		}

		// validation againt trust network
		errInTrustNetworkValidation := authentication.ValidateAgainstTrustNetwork(formulaJSON.User.Publickey)
		if errInTrustNetworkValidation != nil {
			commons.JSONErrorReturn(w, r, errInTrustNetworkValidation.Error(), http.StatusBadRequest, "Expert is not in the trust network ")
			return
		}
		authLayer := authentication.AuthLayer{
			FormulaId:    formulaJSON.MetricExpertFormula.ID,
			ExpertPK:     formulaJSON.User.Publickey,
			ExpertUserID: formulaJSON.User.ID,
			CiperText:    formulaJSON.MetricExpertFormula.CiperText,
			Plaintext:    formulaJSON.MetricExpertFormula.Formula,
		}
		err, errCode := authLayer.ValidateExpertRequest()
		if err != nil {
			commons.JSONErrorReturn(w, r, err.Error(), errCode, "Authentication Issue ")
			return
		}

		formulaArray := formulaJSON.MetricExpertFormula.Formula
		fieldCount := 0
		for i, element := range formulaJSON.MetricExpertFormula.Formula {
			if element.Type == "DATA" {
				formulaArray[i].Type = "VARIABLE"
			} else if element.Type == "CONSTANT" && element.MetricReferenceId != "" {
				formulaArray[i].Type = "REFERREDCONSTANT"
			} else if element.Type == "CONSTANT" && element.MetricReferenceId == "" {
				formulaArray[i].Type = "SEMANTICCONSTANT"
			}
			if element.Type != "OPERATOR" {
				fieldCount++
			}
		}
		formulaJSON.MetricExpertFormula.Formula = formulaArray
		// build the abstract struct and call the SocialImpactExpertFormula
		socialImpactBuilder := protocols.AbstractSocialImpact{
			Blockchain:  formulaJSON.MetricExpertFormula.Blockchain,
			FormulaJSON: formulaJSON,
			FieldCount:  fieldCount,
		}
		socialImpactBuilder.SocialImpactExpertFormula(w, r)
	}
}

// BindMetric method : binds the metric with mutiple formulas
func BindMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	permissionStatus := authentication.HasPermission(r.Header.Get("Authorization"))
	if !permissionStatus.Status || !permissionStatus.IsSubscriptionPaid {
		commons.JSONErrorReturn(w, r, "", http.StatusUnauthorized, "Status Unauthorized")
		return
	} else {
		var metricBindJSON model.MetricDataBindingRequest
		err := json.NewDecoder(r.Body).Decode(&metricBindJSON)
		if err != nil {
			commons.JSONErrorReturn(w, r, err.Error(), http.StatusBadRequest, "Error while decoding the body ")
			return
		}
		errInJsonValidationInMetricBind := validations.ValidateMetricDataBindingRequest(metricBindJSON)
		if errInJsonValidationInMetricBind != nil {
			commons.JSONErrorReturn(w, r, errInJsonValidationInMetricBind.Error(), http.StatusBadRequest, "Request body is invalid, Error :")
			return
		} else {
			metricBuilder := protocols.AbstractSocialImpactMetricBinding{
				Blockchain:     metricBindJSON.Metric.Blockchain,
				MetricBindJSON: metricBindJSON,
			}
			metricBuilder.SocialImpactMetricBinding(w, r)
		}
	}
}
