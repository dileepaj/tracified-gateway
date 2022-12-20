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
	<tr bgcolor="#2f3440">
	  <td style="height: 70px"></td>
	  <td style="height: 70px"></td>
	  <td style="height: 70px"></td>
	</tr>
	<tr bgcolor="#2f3440">
	  <td style="height: 190px"></td>
	  <td bgcolor="#2f3440" style="height: 190px; width: 60%" align="center">
		<table border="0" cellpadding="0" cellspacing="0" width="100%">
		  <tr>
			<td
			  style="
				border-radius: 20px 20px 0px 0px;
				height: 190px;
				box-shadow: 0px 4px 31px 3px rgba(0, 0, 0, 0.15);
			  "
			  align="center"
			  bgcolor="#ffffff"
			>
			  <img
				style="width: 20em; max-width: 80%"
				src="https://tracified-platform-images.s3.ap-south-1.amazonaws.com/Tracified-Blockchain-Meetup-Logo-v2-02+1.png"
			  />
			</td>
		  </tr>
		</table>
	  </td>
	  <td style="height: 190px"></td>
	</tr>
	<tr bgcolor="#f0f0f0">
	  <td></td>
	  <td
		bgcolor="#f0f0f0"
		style="width: 60%; padding-bottom: 70px"
		align="center"
	  >
		<table border="0" cellpadding="0" cellspacing="0" width="100%">
		  <tr>
			<td
			  style="
				padding: 20px 40px 0px 40px;
				box-shadow: 0px 25px 31px 3px rgba(0, 0, 0, 0.15);
			  "
			  align="center"
			  bgcolor="#ffffff"
			>
			  <p
				style="
				  @import url('https://fonts.googleapis.com/css2?family=Inter:wght@100;200;300;400;500;600;700;800;900&display=swap');
				  word-wrap: break-word;
				  font-family: 'Inter', sans-serif;
				  font-style: normal;
				  font-weight: 400;
				  font-size: 20px;
				  line-height: 180.02%;
				  text-align: left;
				"
			  >
				Hi,<br /><br />
				We have received a request to change the password for your Trust
				Network account. <br />
				Please use this code to reset the password: <br />
				` + password + `<br /><br />
				If you did not request a password reset, you can safely ignore
				this email.
				<br /><br />
				Cheers,<br />
				<strong>Tracified Team</strong>
			  </p>
			</td>
		  </tr>
		  <tr>
			<td
			  style="
				border-radius: 0px 0px 20px 20px;
				padding: 0px 40px 20px 40px;
				box-shadow: 0px 25px 31px 3px rgba(0, 0, 0, 0.15);
			  "
			  align="center"
			  bgcolor="#ffffff"
			>
			  <div
				style="
				  background: #d9d9d9;
				  height: 1px;
				  width: 100%;
				  margin: 35px 0px 35px 0px;
				"
			  ></div>
			  <img
				src="https://tracified-platform-images.s3.ap-south-1.amazonaws.com/tracified-logo+(1).png"
				style="width: 210px; height: 44px; text-align: center"
			  />
			  <div
				style="
				  width: 100%;
				  margin-top: 20px;
				  margin-bottom: 20px;
				  text-align: center;
				"
			  >
				<p style="display: inline; padding: 50px 12px 4px 12px">
				  <a
					href="https://www.facebook.com/tracified/?ref=page_internal"
					target="_blank"
					><img
					  style="filter: opacity(0.3) drop-shadow(0 0 0 #acacac)"
					  src="https://tracified-platform-images.s3.ap-south-1.amazonaws.com/NFT_Market/Facebook.png"
				  /></a>
				</p>
				<p style="display: inline; padding: 50px 12px 4px 12px">
				  <a href="https://twitter.com/Tracified1" target="_blank"
					><img
					  style="filter: opacity(0.3) drop-shadow(0 0 0 #acacac)"
					  src="https://tracified-platform-images.s3.ap-south-1.amazonaws.com/NFT_Market/Twitter.png"
				  /></a>
				</p>
				<p style="display: inline; padding: 50px 12px 4px 12px">
				  <a
					href="https://www.instagram.com/tracified_official/"
					target="_blank"
					><img
					  style="filter: opacity(0.3) drop-shadow(0 0 0 #acacac)"
					  src="https://tracified-platform-images.s3.ap-south-1.amazonaws.com/NFT_Market/Instagram.png"
				  /></a>
				</p>
				<p style="display: inline; padding: 50px 12px 4px 12px">
				  <a
					href="https://www.instagram.com/tracified_official/"
					target="_blank"
					><img
					  style="filter: opacity(0.3) drop-shadow(0 0 0 #acacac)"
					  src="https://tracified-platform-images.s3.ap-south-1.amazonaws.com/NFT_Market/LinkedIn.png"
				  /></a>
				</p>
				<p style="display: inline; padding: 50px 12px 4px 12px">
				  <a href="#"
					><img
					  style="filter: opacity(0.3) drop-shadow(0 0 0 #acacac)"
					  src="https://tracified-platform-images.s3.ap-south-1.amazonaws.com/NFT_Market/YouTube.png"
				  /></a>
				</p>
			  </div>
			</td>
		  </tr>
		</table>
	  </td>
	  <td></td>
	</tr>
	<tr bgcolor="#f0f0f0">
	  <td style="height: 70px"></td>
	  <td style="height: 70px"></td>
	  <td style="height: 70px"></td>
	</tr>
  </table>
  `
	return emailTemplate
}
