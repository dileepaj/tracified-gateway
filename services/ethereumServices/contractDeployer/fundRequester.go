package contractdeployer

import (
	"errors"
	"strconv"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/configs"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/sirupsen/logrus"
	gomail "gopkg.in/gomail.v2"
)

//Request funds from the admins via email
//1 Ethereum, 2 Polygon
func RequestFunds(blockchainType int) error {
	var mainnetUrl string
	var testnetUrl string
	var message string
	var subject string
	logger := utilities.NewCustomLogger()
	if blockchainType == 1 {
		subject = `Gateway Ethereum account should be funded`
		mainnetUrl = `https://etherscan.io/address/` + commons.GoDotEnvVariable("ETHEREUMPUBKEY")
		testnetUrl = `https://sepolia.etherscan.io/address/` + commons.GoDotEnvVariable("ETHEREUMPUBKEY")
		message = `<center><h1 style='color: brown;'>Gateway Ethereum account should be funded</h1></center><p>Dear Admins,</p><p> 
		This email is auto-generated to notify that the following gateway Ethereum account is low on Eths, please fund the account.
		<p><b>Public key:</b> ` + commons.GoDotEnvVariable("ETHEREUMPUBKEY") + `</p>` + `<p><a href="` + mainnetUrl + `">View Account on Mainnet</a></p>` + `<p><a href="` + testnetUrl + `">View Account on Testnet</a></p><br><br><p>Thank you</p>`
	} else if blockchainType == 2 {
		subject = `Gateway Polygon account should be funded`
		mainnetUrl = `https://polygonscan.com/address/` + commons.GoDotEnvVariable("ETHEREUMPUBKEY")
		testnetUrl = `https://mumbai.polygonscan.com/address/` + commons.GoDotEnvVariable("ETHEREUMPUBKEY")
		message = `<center><h1 style='color: brown;'>Gateway Polygon account should be funded</h1></center><p>Dear Admins,</p><p> 
		This email is auto-generated to notify that the following gateway Polygon account is low on Matics, please fund the account.
		<p><b>Public key:</b> ` + commons.GoDotEnvVariable("ETHEREUMPUBKEY") + `</p>` + `<p><a href="` + mainnetUrl + `">View Account on Mainnet</a></p>` + `<p><a href="` + testnetUrl + `">View Account on Testnet</a></p><br><br><p>Thank you</p>`
	} else {
		logger.LogWriter("Invalid blockchain type", constants.ERROR)
		return errors.New("Invalid blockchain type")
	}

	for _, email := range configs.EthereumNotificationEmails {
		msg := gomail.NewMessage()
		msg.SetHeader("From", commons.GoDotEnvVariable("EMAILADRESSFORNOTIFICATIONSENDER"))
		msg.SetHeader("To", email)
		msg.SetHeader("Subject", subject)
		msg.SetBody("text/html", message)
		port, errWhenConvertingToStr := strconv.Atoi(commons.GoDotEnvVariable("GOMAILPORT"))
		if errWhenConvertingToStr != nil {
			logrus.Error("Issue when converting string to int, ERROR : " + errWhenConvertingToStr.Error())
			return errors.New("Issue when converting string to int, ERROR : " + errWhenConvertingToStr.Error())
		}
		n := gomail.NewDialer(commons.GoDotEnvVariable("GMAILHOST"), port, commons.GoDotEnvVariable("EMAILADRESSFORNOTIFICATIONSENDER"), commons.GoDotEnvVariable("SENDER_EMAILADRESS_APPPWD"))
		errWhenDialAndSending := n.DialAndSend(msg)
		if errWhenDialAndSending != nil {
			logrus.Error("Email sending issue, ERROR : " + errWhenDialAndSending.Error())
			return errors.New("Email sending issue, ERROR : " + errWhenDialAndSending.Error())
		}

	}
	return nil
}
