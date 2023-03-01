package contractdeployer

import (
	"errors"
	"strconv"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/configs"
	"github.com/sirupsen/logrus"
	gomail "gopkg.in/gomail.v2"
)

func RequestFunds() error {
	url := `https://etherscan.io/address/` + commons.GoDotEnvVariable("ETHEREUMPUBKEY")
	message := `<center><h1 style='color: brown;'>Gateway Ethereum account should be funded</h1></center><p>Dear Admins,</p><p> 
This email is auto-generated to notify that the following gateway Ethereum account is low on Eths, please fund the account.
<p><b>Public key:</b> ` + commons.GoDotEnvVariable("ETHEREUMPUBKEY") + `</p>` + `<p><a href="` + url + `">View Account</p><br><br><p>Thank you</p>`

	subject := `Gateway Ethereum account should be funded`

	for _, email := range configs.EthereumNotificationEmails {
		msg := gomail.NewMessage()
		msg.SetHeader("From", commons.GoDotEnvVariable("sender_emailadress"))
		msg.SetHeader("To", email)
		msg.SetHeader("Subject", subject)
		msg.SetBody("text/html", message)
		port, errWhenConvertingToStr := strconv.Atoi(commons.GoDotEnvVariable("GOMAILPORT"))
		if errWhenConvertingToStr != nil {
			logrus.Error("Issue when converting string to int, ERROR : " + errWhenConvertingToStr.Error())
			return errors.New("Issue when converting string to int, ERROR : " + errWhenConvertingToStr.Error())
		}
		n := gomail.NewDialer(commons.GoDotEnvVariable("GMAILHOST"), port, commons.GoDotEnvVariable("sender_emailadress"), commons.GoDotEnvVariable("SENDER_EMAILADRESS_APPPWD"))
		errWhenDialAndSending := n.DialAndSend(msg)
		if errWhenDialAndSending != nil {
			logrus.Error("Email sending issue, ERROR : " + errWhenDialAndSending.Error())
			return errors.New("Email sending issue, ERROR : " + errWhenDialAndSending.Error())
		}

	}
	return nil
}
