package main

import (
	"fmt"
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

func main() {
	port := getPort()
	headersOk := handlers.AllowedHeaders([]string{"Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	c := cron.New()
	c.AddFunc("@every 30m", func() {
		services.CheckCOCStatus()
	})
	c.Start()

	router := routes.NewRouter()
	fmt.Println("Gateway Started @port " + port)
	http.ListenAndServe(port, handlers.CORS(originsOk, headersOk, methodsOk)(router))

}
