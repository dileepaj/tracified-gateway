package constants

import "github.com/dileepaj/tracified-gateway/commons"

// PublicKey The Stellar Testnet PublicKey for Gateway's Account, Should be replaced when stellar resets the testnet
// var PublicKey ="GDUTBNH7VVFUZTX4WCMW34WO6TOV6VH4QSUSZ5YOXF4TTRUWMQAZ3HSF"

// Public network Keys
// var SecretKey = "SCJW6N76UWARARCTHGVT5DYL26L4S77KEUXBHKUJBF526R2QOI4HC7IU"

// Public network Keys and test network key
var PublicKey = commons.GoDotEnvVariable("PublicKey")

// SecretKey The Stellar Testnet SecretKey for Gateway's Account, Should be replaced when stellar resets the testnet
var SecretKey = commons.GoDotEnvVariable("SecretKey")
