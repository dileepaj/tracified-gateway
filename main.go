package main

import (
	"net/http"
	"os"

	"github.com/dileepaj/tracified-gateway/adminDAO"
	"github.com/dileepaj/tracified-gateway/api/routes"
	"github.com/dileepaj/tracified-gateway/apiDemo/route"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/services"
	"github.com/dileepaj/tracified-gateway/services/rabbitmq"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	envName := commons.GoDotEnvVariable("ENVIRONMENT")

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

	c.AddFunc("@every 5m", func() {
		services.CheckContractStatus()
	})
	c.Start()
	router := routes.NewRouter()
	// rabbit mq server
	go rabbitmq.ReceiverRmq()
	go rabbitmq.ReleaseLock()
	// serve swagger documentation
	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	router.Handle("/docs", sh)
	router.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	//initial log file when server starts
	utilities.CreateLogFile()
	//create logger
	logger := utilities.NewCustomLogger()
	logger.LogWriter("Gateway Started @port "+port+" with "+envName+" environment", 1)

	//http.ListenAndServe(port, handlers.CORS(originsOk, headersOk, methodsOk)(router))

	// API-Demo core re-structure
	r := mux.NewRouter()
	// Define public routes
	route.DefineRoutes(r)
	// Start the server
	http.ListenAndServe(":1776", handlers.CORS(originsOk, headersOk, methodsOk)(r))
}
