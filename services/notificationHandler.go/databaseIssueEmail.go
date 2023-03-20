package notificationhandler

import (
	"errors"
	"strconv"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/configs"
	"github.com/sirupsen/logrus"
	gomail "gopkg.in/gomail.v2"
)

func InformDBConnectionIssue(information string, errorMessage string) error {
	subject := "Gateway Database Error"
	message := `<center><h1 style='color: brown;'>An Error Occurred in the Gateway Database</h1></center><p>Dear Admins,</p><p> 
	This email is auto-generated to notify that an error occurred in the gateway database connection when trying to <b> ` + information + `</b>.
	<p><b>Error:</b> ` + errorMessage + `.</p> please check and get the necessary actions.`
	for _, email := range configs.DatabaseNotificationEmails {
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
			logrus.Error("Email sending issue for database collection issue, ERROR : " + errWhenDialAndSending.Error())
			return errors.New("Email sending issue, ERROR : " + errWhenDialAndSending.Error())
		}

	}

	return nil
}