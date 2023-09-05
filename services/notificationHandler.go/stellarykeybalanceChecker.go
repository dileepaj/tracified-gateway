package notificationhandler

import (
	"fmt"
	"log"
	"strconv"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/configs"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	gomail "gopkg.in/gomail.v2"
)

func CheckStellarAccountBalance(publicKey string) error {
	client := commons.GetHorizonClient()
	_, err := keypair.Parse(publicKey)
	if err != nil {
		log.Fatal("Invalid Stellar public key:", err)
	}

	// Fetch the account details for the given public key
	accountRequest := horizonclient.AccountRequest{
		AccountID: publicKey,
	}

	account, err := client.AccountDetail(accountRequest)
	if err != nil {
		log.Fatal("Error fetching account details:", err)
	}

	// Find and display the XLM balance
	var xlmBalance horizon.Balance
	var sellingLiabilities string
	for _, balance := range account.Balances {
		if balance.Asset.Type == "native" {
			xlmBalance = balance
			sellingLiabilities = balance.SellingLiabilities
			fmt.Println("seeling liabilities: ", balance.SellingLiabilities)
			break
		}
	}

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
	minimumBalance, err := GetAccountminBalance(account.SubentryCount, account.NumSponsoring, account.NumSponsored, f)
	bufferAmount := commons.GoDotEnvVariable("STELLAR_KEY_LOW_BALANCE_BUFFER_AMOUNT")
	bufferAmountPass, bufferConvertError := strconv.ParseFloat(bufferAmount, 64)
	if bufferConvertError != nil {
		logrus.Error("Error:", err)
		return err
	}
	if (minimumBalance + bufferAmountPass) >= floatBalance {
		message := `
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
				<p>Tracified account: <span class="strong">` + publicKey + `</span> is low in balance (<span class="strong">` + xlmBalance.Balance + ` XLM</span>). Please fund it. The minimum balance for the account to be active is: <span class="strong">` + strconv.FormatFloat(minimumBalance, 'f', -1, 64) + ` XLM</span>.</p>
			</div>
		</body>
		</html>
		`
		for _, email := range configs.NftEmailList {
			SendLowBalanceEmail(message, "Blockchain account low balance", email)
		}

	}
	return nil
}

func GetAccountminBalance(subEntryCount int32, numberOfSponsoring uint32, numOfSponsored uint32, sellingLiabilities float64) (float64, error) {
	minimumBalance := 0.0
	baseReserve := commons.GoDotEnvVariable("STELLAR_BASE_RESERVE")
	basefee, err := strconv.ParseFloat(baseReserve, 64)
	if err != nil {
		logrus.Error("Error:", err)
		return minimumBalance, err
	}
	// Define the base reserve constant
	var stellarBaseReserve = basefee

	// Start with twice the base reserve as the minimum balance
	mandatoryMinimumBalance := stellarBaseReserve * 2

	// Calculate the minimum balance
	minimumBalance = (mandatoryMinimumBalance+float64(subEntryCount)+float64(numberOfSponsoring)-float64(numOfSponsored))*stellarBaseReserve + float64(sellingLiabilities)
	fmt.Printf("Minimum Balance: %f\n", minimumBalance)
	return minimumBalance, nil
}

func SendLowBalanceEmail(message, subject, to string) error {
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
