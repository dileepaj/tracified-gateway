package services

import (
	"strconv"

	gomail "gopkg.in/gomail.v2"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/sirupsen/logrus"
)

func SendingEmail(message, subject, to string) error {
	// message:=`<center><h1 style='color: brown;'>Security Risk</h1></center><p>Dear Sir,Madam,</p><p>Exceed the weekly request limit.</p><hr><center><img src="https://tracified-platform-images.s3.ap-south-1.amazonaws.com/tracified-logo+(1).png" style="width:20em"></center>`
	msg := gomail.NewMessage()
	msg.SetHeader("From", commons.GoDotEnvVariable("GOMAILSENDER"))
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", message)
	port, err := strconv.Atoi(commons.GoDotEnvVariable("GOMAILPORT"))
	if err != nil {
		logrus.Error("Issue when converting string to int  ", err)
		return err
	}
	n := gomail.NewDialer(commons.GoDotEnvVariable("GMAILHOST"), port, commons.GoDotEnvVariable("GOMAILSENDER"), commons.GoDotEnvVariable("GOMAILSENDERPW"))
	err1 := n.DialAndSend(msg)
	if err1 != nil {
		logrus.Error("Email sending issue ", err1)
		return err1
	}
	return err1
}
