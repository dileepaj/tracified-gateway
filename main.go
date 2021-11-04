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

	c.AddFunc("@every 1m", func() {
		services.CheckTempOrphan()
	})
	c.Start()

	router := routes.NewRouter()
	fmt.Println("Gateway Started @port " + port + " with " + envName + " environment")
	http.ListenAndServe(port, handlers.CORS(originsOk, headersOk, methodsOk)(router))

}
