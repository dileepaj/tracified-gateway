package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dileepaj/tracified-gateway/adminDAO"
	"github.com/dileepaj/tracified-gateway/api/routes"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/pools"
	"github.com/dileepaj/tracified-gateway/services"
	"github.com/gorilla/handlers"
	"github.com/robfig/cron"
)

func getPort() string {
	p := os.Getenv("GATEWAY_PORT")
	if p != "" {
		return ":" + p
	}
	return ":8000"
}

func main() {
	// godotenv package
	envName := commons.GoDotEnvVariable("BRANCH_NAME")

	// getEnvironment()
	port := getPort()
	headersOk := handlers.AllowedHeaders([]string{"Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	commons.ConstructConnectionPool()
	adminDAO.ConstructAdminConnectionPool()

	c := cron.New()
	c.AddFunc("@every 30m", func() {
		services.CheckCOCStatus()
	})

	c.AddFunc("@every 12h", func() {
		services.CheckTestimonialStatus()
		services.CheckOrganizationStatus()
	})

	c.AddFunc("@every 1m", func() {
		services.CheckTempOrphan()
	})
	c.Start()
	router := routes.NewRouter()

	var poolJson []pools.BuildPool
	pool1 := pools.BuildPool{
		Coin1:               "N",
		DepositeAmountCoin1: "10000",
		Coin2:               "O",
		DepositeAmountCoin2: "70000",
		Ratio:               2,
	}
	pool2 := pools.BuildPool{
		Coin1:               "P",
		DepositeAmountCoin1: "10000",
		Coin2:               "Q",
		DepositeAmountCoin2: "60000",
		Ratio:               2,
	}
	poolJson = append(poolJson, pool1,pool2)
	pools.CreatePoolsUsingJson(poolJson)

	// var buildPathPayment pools.BuildPathPayment
	// buildPathPayment = pools.BuildPathPayment{
	// 	SendingCoin: pools.Coin{
	// 		CoinName: "E",
	// 		Amount:   "100",
	// 	},
	// 	ReceivingCoin: pools.Coin{
	// 		CoinName: "F",
	// 		Amount:   "",
	// 	},
	// 	BatchAccountPK:     "GD6ZW4L3Y5E3JEW4TLSYGV3PC7TBYN6AXIGVW54J6HRH2J3HDZMBA62C",
	// 	BatchAccountSK:     "SD6G4TWP5PTCKIO4XOHCZE5IDJDNOEIVTOCRC6YRG5B3IO42SRMWYKU4",
	// 	CoinIssuerAccontPK: "GDBXHHHG7CKIODJIUPU46W52RUDUMJ3PJQOSWF24R3VGVRLPVHWNT5DI",
	// 	PoolId:             "",
	// 	ProductId:          "",
	// 	EquationId:         "",
	// 	TenantId:           "",
	// }
	//amount, err:=pools.CoinConvert(buildPathPayment)

	//amount, err := pools.GetConvertedCoinAmount("E", "100", "", "F", "GBSLTJX4NRMTPTQ2EJZJ24U44K7ZWY3CPGBZTV623PLTLIFXWK3T4CD6")

	fmt.Println("Gateway Started @port " + port + " with " + envName + " environment")
	http.ListenAndServe(port, handlers.CORS(originsOk, headersOk, methodsOk)(router))
}
