package services

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/sirupsen/logrus"
	gomail "gopkg.in/gomail.v2"
)

const (
	lowercase = "abcdefghijklmnopqrstuvwxyz"
	uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers   = "0123456789"
)

func createPassword(length int, charSets []string) string {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	// Create a string to hold the password.
	password := ""

	// Loop until the password is the desired length.
	for len(password) < length {
		// Pick a random character set.
		charSet := charSets[seededRand.Intn(len(charSets))]

		// Pick a random character from the chosen character set.
		c := charSet[rand.Intn(len(charSet))]

		// Append the character to the password.
		password += string(c)
	}

	return password
}
func GeneratePassword() (string, error) {
	passwordLength, pwdLenerr := strconv.Atoi(commons.GoDotEnvVariable("TRUST_NETWORK_PASSWORDRESET_LENGTH"))
	if pwdLenerr != nil {
		logrus.Println("error getting password from env: ", pwdLenerr.Error())
		return "", pwdLenerr
	}
	charSets := []string{lowercase, uppercase, numbers}
	newpassword := createPassword(passwordLength, charSets)
	return newpassword, nil

}

func SendEmail(newpassword string, senderEmail string) error {
	hostEmail := commons.GoDotEnvVariable("HOST_EMAIL")
	hostEmailPort, emailportErr := strconv.Atoi(commons.GoDotEnvVariable("EMAIL_PORT"))
	if emailportErr != nil {
		logrus.Println("error getting email port from env: ", emailportErr.Error())
		return emailportErr
	}
	senderEmailAddress := commons.GoDotEnvVariable("SENDER_EMAILADRESS")
	senderEmailKey := commons.GoDotEnvVariable("SENDER_EMAILADRESS_APPPWD")

	emailToSend := getEmail(newpassword)

	msg := gomail.NewMessage()
	msg.SetHeader("From", senderEmailAddress)
	msg.SetHeader("To", senderEmail)
	msg.SetHeader("Subject", "Blockchain meetup password reset")
	msg.SetBody("text/html", emailToSend)
	n := gomail.NewDialer(hostEmail, hostEmailPort, senderEmailAddress, senderEmailKey)
	if emailSendErr := n.DialAndSend(msg); emailSendErr != nil {
		logrus.Println("Error sending email: ", emailSendErr.Error())
		return emailSendErr
	}
	return nil
}
func EncodeTrustnetworkResetPassword(userPassword string) ([]byte, error) {
	logrus.Println("password to save : ", userPassword)
	pwdbyte := commons.Encrypt(userPassword)
	logrus.Println("Encrypted pwd : ", pwdbyte)

	decrypt := commons.Decrypt(pwdbyte)
	logrus.Println("decrypted: ", decrypt)
	return pwdbyte, nil
}
func DecodeTrustnetworkResetPassword(userPassword string) (string, error) {
	return "", nil
}

func getEmail(password string) string {
	var emailTemplate = `<table border="0" cellpadding="0" cellspacing="0" width="100%">
    <!-- LOGO -->
    <tr>
        <td bgcolor="#021d28" align="center">
            <table border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 600px;">
                <tr>
                    <td align="center" valign="top" style="padding: 40px 10px 40px 10px;"> </td>
                </tr>
            </table>
        </td>
    </tr>
    <tr>
        <td bgcolor="#021d28" align="center" style="padding: 0px 10px 0px 10px;">
            <table border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 600px;">
                <tr>
                   <td style="padding: 40px 20px 20px; border-radius: 4px 4px 0px 0px; color: #111111; font-family: Lato, Helvetica, Arial, sans-serif; font-size: 40px; font-weight: 400; letter-spacing: 4px; line-height: 48px; height: 31px;" align="center" valign="top" bgcolor="#ffffff">
				   <!-- BC MEETUP LOGO -->
				   <img style="width: 3em;" src="https://tracified-profile-images.s3.ap-south-1.amazonaws.com/RURI+1.png">
				   </td>
                </tr>
            </table>
        </td>
    </tr>
    <tr>
        <td bgcolor="#f4f4f4" align="center" style="padding: 0px 10px 0px 10px;">
            <table border="0" cellpadding="0" cellspacing="0" width="100%" style="max-width: 600px;">
              <tr>
      <td style="padding: 20px 30px 15px 20px; color: #666666; font-family: Lato, Helvetica, Arial, sans-serif; font-size: 16px; font-weight: 400; line-height: 25px; width: 100%;" align="left" bgcolor="#ffffff">
      <p style="margin: 0;">Hi<br />Please use the follwing password to complete your password reset process</p>
      </td>
      <tr>
                    <td bgcolor="#ffffff" align="left">
                        <table width="100%" border="0" cellspacing="0" cellpadding="0">
                            <tr>
                                <td bgcolor="#ffffff" align="center" style="padding: 2px 30px 10px 30px;">
                                    <table border="0" cellspacing="0" cellpadding="0">
                                        <tr>
                                            <td align="center" style="border-radius: 3px;" bgcolor="#00466a"><div style="font-size: 20px; font-family: Helvetica, Arial, sans-serif; color: #62FFA3; text-decoration: none; color: #62FFA3; text-decoration: none; padding: 15px 25px; border-radius: 2px; border: 1px solid #00466a; display: inline-block; text-align:center">` + password + `</a></td>
                                        </tr>
                                    </table>
                                </td>
                            </tr>
                        </table>
                    </td>
                <tr>
                    
                </tr>
        <td bgcolor="#ffffff" align="left" style="padding: 0px 30px 10px 30px; color: #666666; font-family: 'Lato', Helvetica, Arial, sans-serif; font-size: 16px; font-weight: 400; line-height: 25px;">
          <p style="margin: 0;"> <strong>Note - </strong>Please note that the One Time Password is valid for a period of one month only.
                    </td>						
      </tr>
                <tr>
        <td bgcolor="#ffffff" align="left" style="padding: 0px 30px 10px 30px; color: #666666; font-family: 'Lato', Helvetica, Arial, sans-serif; font-size: 16px; font-weight: 400; line-height: 25px;">
          <p style="margin: 0;">Enjoy your NFT !
                    </td>						
      </tr>
                <tr>
                    <td bgcolor="#ffffff" align="left" style="padding: 0px 30px 40px 30px; border-radius: 0px 0px 4px 4px; color: #666666; font-family: 'Lato', Helvetica, Arial, sans-serif; font-size: 16px; font-weight: 400; line-height: 25px;">
                        <p style="margin: 0;">Cheers,<br>Team RURI</p>
                        <hr style="background-color: #D9D9D9; ">
                        <p style="color: #878787;"><center>Powered by</center></p>
                        <center><img src="https://tracified-platform-images.s3.ap-south-1.amazonaws.com/Tracified-NFT-v2.png" style="width:20em"></center>
                    </td>
                </tr>
            </table>
        </td>
    </tr>        
  </table>`
	return emailTemplate
}
