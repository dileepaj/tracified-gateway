package authentication

import (
	"github.com/dileepaj/tracified-gateway/configs"
	"github.com/dileepaj/tracified-gateway/services"
	"github.com/sirupsen/logrus"
)

func NotifyToAdmin() {
	message := `<center><h1 style='color: brown;'>Security Risk</h1></center><p>Dear Sir,Madam,</p><p>Exceed the weekly request limit.</p><hr><center><img src="https://tracified-platform-images.s3.ap-south-1.amazonaws.com/tracified-logo+(1).png" style="width:20em"></center>`

	for _, email := range configs.NotificationEmails {
		err := services.SendingEmail(message, "RISK", email)
		if err != nil {
			logrus.Error("Email sending issue ", err)
		}
	}
}