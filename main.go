package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dileepaj/tracified-gateway/adminDAO"
	"github.com/dileepaj/tracified-gateway/api/routes"
	"github.com/dileepaj/tracified-gateway/commons"
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

	// var buildPathPayment model.BuildPathPayment
	// buildPathPayment = model.BuildPathPayment{
	// 	SendingCoin: model.Coin{
	// 		CoinName: "NEW1",
	// 		Amount:   "200",
	// 	},
	// 	ReceivingCoin: model.Coin{
	// 		CoinName: "NEW3",
	// 		Amount:   "",
	// 	},
	// 	IntermediateCoins: []model.Coin{
	// 		{
	// 			CoinName: "NEW2",
	// 			Amount:   "",
	// 		},
	// 	},
	// 	BatchAccountPK:     "GCBZ7J5434MIU3AYKCI2FPMLBV5LQBKIZYG2C5QMVEWOTIT2XM2AVWSG",
	// 	BatchAccountSK:     "SA4C7PM67PYJQ2SMRRXDUIX5EUMV725JGDXZXMLKG2VPLW4UYHJLUVSI",
	// 	CoinIssuerAccontPK: "GBRCIPHDMVGMQUUCP2DWHB55RMZOVL6JPE4KCH2AS2MODVHL6NHC642R",
	// 	PoolId:             "",
	// 	ProductId:          "",
	// 	EquationId:         "",
	// 	TenantId:           "",
	// }
	// hash, err := pools.CoinConvert(buildPathPayment)
	// fmt.Println("dsadasdsa----------------          ", hash, err)
	// amount, err := pools.GetConvertedCoinAmount("BTC", "100", "USDT", "GBRCIPHDMVGMQUUCP2DWHB55RMZOVL6JPE4KCH2AS2MODVHL6NHC642R")
	// fmt.Println("final--",amount, err)

	fmt.Println("Gateway Started @port " + port + " with " + envName + " environment")
	http.ListenAndServe(port, handlers.CORS(originsOk, headersOk, methodsOk)(router))
}
