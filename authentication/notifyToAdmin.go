package authentication

import (
	"github.com/dileepaj/tracified-gateway/configs"
	"github.com/dileepaj/tracified-gateway/services"
	"github.com/sirupsen/logrus"
)

func NotifyToAdmin(expertId, expertPk string) {
	message := `<center><h1 style='color: brown;'>Security Risk</h1></center><p>Dear Admins,</p><p> 
	This email is auto-generated to notify that the user with the following details exceeded the weekly limit of 20 requests for defining formulas.
	<p><b>ExpertID:</b> ` + expertId + `</p><p>
	<b>Public Key:</b> ` + expertPk + `</p><br><p>Thank you</p><hr><center><img src="https://tracified-platform-images.s3.ap-south-1.amazonaws.com/tracified-logo+(1).png" style="width:20em"></center>`

	for _, email := range configs.NotificationEmails {
		err := services.SendingEmail(message, "RISK", email)
		if err != nil {
			logrus.Error("Email sending issue ", err)
		}
	}
}
