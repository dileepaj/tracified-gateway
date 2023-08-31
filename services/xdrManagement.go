package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	fosponsoring "github.com/dileepaj/tracified-gateway/nft/stellar/FOSponsoring"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var accountStatus string

func SubmitFOData(Response model.TransactionData) (string, error) {
	if commons.GoDotEnvVariable("FONEW_FLAG") == "TRUE" {
		if Response.XDR != "" && Response.FOUser != "" && Response.AccountIssuer != "" {
			resp, err := http.Get(commons.GetHorizonClient().HorizonURL + "accounts/" + Response.FOUser)
			if err != nil {
				logrus.Error("Error making HTTP request:", err)
				return "", nil
			}
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusNotFound {
				accountStatus = "0"
			} else {
				accountStatus = "1"
			}
			if accountStatus == "0" {
				logrus.Error("Account of FO User is inactive")
			}

			if accountStatus == "1" {
				result1, err := http.Get(commons.GetHorizonClient().HorizonURL + "accounts/" + Response.AccountIssuer)
				body, err := ioutil.ReadAll(result1.Body)
				if err != nil {
					log.Error("Error while read response " + err.Error())
				}
				var balances model.BalanceResponse

				err = json.Unmarshal(body, &balances)
				if err != nil {
					log.Error("Error while json.Unmarshal(body, &balance) " + err.Error())
				}

				balance := balances.Balances[0].Balance
				floatValue, err := strconv.ParseFloat(balance, 64)
				if floatValue < 10 {
					hash, err := fosponsoring.FundAccount(Response.AccountIssuer)
					if err != nil {
						log.Error("Error while funding issuer " + err.Error())
					}
					logrus.Info("funded and hash is : ", hash)
					var TransactionPayload = model.TransactionData{
						FOUser:        Response.FOUser,
						AccountIssuer: Response.AccountIssuer,
						XDR:           Response.XDR,
					}
					xdr, err := fosponsoring.BuildSignedSponsoredXDR(TransactionPayload)
					if err != nil {
						log.Error(err)
					} else {
						logrus.Info("xdr base64 been passed to frontend : ", xdr)
						return xdr, nil
					}
				} else {
					var TransactionPayload = model.TransactionData{
						FOUser:        Response.FOUser,
						AccountIssuer: Response.AccountIssuer,
						XDR:           Response.XDR,
					}
					xdr, err := fosponsoring.BuildSignedSponsoredXDR(TransactionPayload)
					if err != nil {
						log.Error(err)
					} else {
						logrus.Info("xdr base64 been passed to frontend : ", xdr)
						return xdr, nil
					}
				}
			} else {
				return "", nil
			}
		}
		return "", nil
	}
	return "", nil
}
