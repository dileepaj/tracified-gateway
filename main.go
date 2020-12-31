package main

import (
	"fmt"
	"github.com/astaxie/beego/core/config"
	"github.com/dileepaj/tracified-gateway/commons"
	"log"

	// "log"
	"net/http"
	"os"

	"github.com/dileepaj/tracified-gateway/api/routes"
	"github.com/dileepaj/tracified-gateway/services"
	"github.com/gorilla/handlers"

	// "github.com/joho/godotenv"
	"github.com/robfig/cron"
)

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":8000"
}

// func getEnvironment() {
// 	err := godotenv.Load()
// 	  if err != nil {
// 	    log.Fatal("Error loading .env file")
// 	  }
// }

func main() {

	//Read conf/app.conf file
	conf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		log.Fatalf("failed to parse config file err: %s", err.Error())
	}
	commons.ConstructConnectionPool(conf)

	// getEnvironment()
	port := getPort()
	headersOk := handlers.AllowedHeaders([]string{"Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	c := cron.New()
	c.AddFunc("@every 30m", func() {
		services.CheckCOCStatus()
	})

	c.AddFunc("@every 1m", func() {
		services.CheckTempOrphan()
	})
	c.Start()

	router := routes.NewRouter()
	fmt.Println("Gateway Started @port " + port)
	http.ListenAndServe(port, handlers.CORS(originsOk, headersOk, methodsOk)(router))

}
