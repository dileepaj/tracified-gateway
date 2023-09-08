package notificationhandler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	gomail "gopkg.in/gomail.v2"
)

// CheckStellarAccountBalance checks the balance of a Stellar account associated with the given public key.
// If the balance falls below a certain threshold, it sends low balance notifications to a list of email addresses.
// It returns an error if there is any issue during the process.
func CheckStellarAccountBalance(publicKey string) error {
	// Initialize a Stellar Horizon client
	client := commons.GetHorizonClient()

	// Parse the provided public key
	_, err := keypair.Parse(publicKey)
	if err != nil {
		// Log a fatal error and terminate if the public key is invalid
		logrus.Error("Invalid Stellar public key:", err)
		return err
	}

	// Create an account request to fetch account details from Horizon
	accountRequest := horizonclient.AccountRequest{
		AccountID: publicKey,
	}

	// Fetch the account details for the given public key
	account, err := client.AccountDetail(accountRequest)
	if err != nil {
		// Log a fatal error and terminate if there is an issue fetching account details
		logrus.Error("Error fetching account details:", err)
	}

	// Find and display the XLM balance and selling liabilities
	var xlmBalance horizon.Balance
	var sellingLiabilities string
	for _, balance := range account.Balances {
		if balance.Asset.Type == "native" {
			xlmBalance = balance
			sellingLiabilities = balance.SellingLiabilities
			break
		}
	}

	// Parse the XLM balance and selling liabilities as float64
	floatBalance, err := strconv.ParseFloat(xlmBalance.Balance, 64)
	if err != nil {
		logrus.Error("Error:", err)
		return err
	}

	f, err := strconv.ParseFloat(sellingLiabilities, 64)
	if err != nil {
		logrus.Error("Error:", err)
		return err
	}

	// Calculate the minimum account balance required
	minimumBalance, err := GetAccountminBalance(account.SubentryCount, account.NumSponsoring, account.NumSponsored, f)

	// Retrieve the buffer amount from environment variables
	bufferAmount := commons.GoDotEnvVariable("STELLAR_KEY_LOW_BALANCE_BUFFER_AMOUNT")
	bufferAmountPass, bufferConvertError := strconv.ParseFloat(bufferAmount, 64)
	if bufferConvertError != nil {
		logrus.Error("Error:", err)
		return err
	}

	// Check if the account balance is below the minimum balance plus buffer
	if (minimumBalance + bufferAmountPass) >= floatBalance {
		logrus.Printf("Low Balance detected for %s sending emails", publicKey)

		// Generate an email message template
		message := GetEmailMessageTemplate(publicKey, xlmBalance.Balance, minimumBalance)

		// Get a list of email addresses to send notifications
		GetEmailList := GetEmailList()

		// Send low balance notifications to each email address in the list
		for _, email := range GetEmailList {

			SendLowBalanceEmail(message, "Blockchain account low balance", email)
		}
	}

	// Return nil (no error) if the process completes successfully
	return nil
}

// GetAccountminBalance calculates the minimum account balance required for an account
// based on the provided parameters such as subEntryCount, numberOfSponsoring, numOfSponsored, and sellingLiabilities.
// It returns the calculated minimum balance and any error encountered during the calculation.
func GetAccountminBalance(subEntryCount int32, numberOfSponsoring uint32, numOfSponsored uint32, sellingLiabilities float64) (float64, error) {
	// Initialize the minimumBalance variable to 0.0
	minimumBalance := 0.0

	// Retrieve the base reserve value from environment variables
	baseReserve := commons.GoDotEnvVariable("STELLAR_BASE_RESERVE")

	// Convert the baseReserve from string to float64
	basefee, err := strconv.ParseFloat(baseReserve, 64)
	if err != nil {
		// Log an error message and return 0.0 if there is an error converting the base reserve
		logrus.Error("Error:", err)
		return minimumBalance, err
	}

	// Define the constant value for the Stellar base reserve
	var stellarBaseReserve = basefee

	// Calculate the mandatory minimum balance as twice the base reserve
	mandatoryMinimumBalance := stellarBaseReserve * 2

	// Calculate the minimum balance based on the provided parameters
	minimumBalance = (mandatoryMinimumBalance+float64(subEntryCount)+float64(numberOfSponsoring)-float64(numOfSponsored))*stellarBaseReserve + float64(sellingLiabilities)

	// Print the calculated minimum balance for debugging purposes
	fmt.Printf("Minimum Balance: %f\n", minimumBalance)

	// Return the calculated minimum balance and no error
	return minimumBalance, nil
}

// SendLowBalanceEmail sends an email with the specified message, subject, and recipient.
// It returns an error if there is any issue with sending the email.
func SendLowBalanceEmail(message, subject, to string) error {
	// Create a new email message
	msg := gomail.NewMessage()

	// Set the "From" header of the email
	msg.SetHeader("From", commons.GoDotEnvVariable("MAIL_SENDER"))

	// Set the "To" header of the email
	msg.SetHeader("To", to)

	// Set the "Subject" header of the email
	msg.SetHeader("Subject", subject)

	// Set the body of the email as HTML with the provided message
	msg.SetBody("text/html", message)

	// Convert the MAIL_PORT from a string to an integer
	port, err := strconv.Atoi(commons.GoDotEnvVariable("MAIL_PORT"))
	if err != nil {
		// Log an error message if there is an issue converting the string to an int
		logrus.Error("Issue when converting string to int  ", err)
		return err
	}

	// Create a new mail dialer with the mail host, port, sender, and sender's app key
	n := gomail.NewDialer(commons.GoDotEnvVariable("MAIL_HOST"), port, commons.GoDotEnvVariable("MAIL_SENDER"), commons.GoDotEnvVariable("MAIL_SENDER_APP_KEY"))

	// Dial and send the email message
	err1 := n.DialAndSend(msg)
	if err1 != nil {
		// Log an error message if there is an issue sending the email
		logrus.Error("Email sending issue ", err1)
		return err1
	}

	// Return any error that occurred during the email sending process
	return err1
}

// GetEmailList retrieves a list of email addresses from a configuration variable
// and returns them as a slice of strings.
func GetEmailList() []string {
	// Retrieve the email list as a string from the configuration variable
	list := commons.GoDotEnvVariable("NOTIFIER_EMAILS_FOR_LOW_BALANCE_WARNING")

	// Remove square brackets from the string
	input := strings.Trim(list, "[]")

	// Split the string into a slice of email addresses using commas as the delimiter
	emailList := strings.Split(input, ",")

	// Trim leading and trailing whitespace from each email address
	for i, email := range emailList {
		emailList[i] = strings.TrimSpace(email)
	}

	// Return the list of email addresses as a slice of strings
	return emailList
}

// GetEmailMessageTemplate generates an HTML email message template with placeholders for
// the public key, account balance, and minimum balance. It returns the HTML template as a string.
// The placeholders will be filled with actual values when sending an email.
func GetEmailMessageTemplate(publicKey string, balance string, minimumBalance float64) string {
	env := commons.GoDotEnvVariable("ENVIRONMENT")
	// Define an HTML email template with placeholders for the provided values
	template := `
	<html>
	<head>
		<style>
			body {
				font-family: Arial, sans-serif;
			}
			.container {
				padding: 20px;
				background-color: #f7f7f7;
				border-radius: 10px;
			}
			.strong {
				font-weight: bold;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<p>` + env + ` Tracified account: <span class="strong">` + publicKey + `</span> is low in balance (<span class="strong">` + balance + ` XLM</span>). Please fund it. The minimum balance for the account to be active is: <span class="strong">` + strconv.FormatFloat(minimumBalance, 'f', -1, 64) + ` XLM</span>.</p>
		</div>
	</body>
	</html>
	`

	// Return the generated HTML template as a string
	return template
}
