package main

import (
	"net/http"
	"os"

	"github.com/dileepaj/tracified-gateway/api/routes"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/configs"
	"github.com/dileepaj/tracified-gateway/services"
	notificationhandler "github.com/dileepaj/tracified-gateway/services/notificationHandler.go"
	"github.com/dileepaj/tracified-gateway/services/rabbitmq"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/handlers"
	"github.com/robfig/cron/v3"
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
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Token"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// Register services, this will be used to schedule workers
	// Added here to avoid circular import, refactor if there is a better way
	configs.QueueBackLinks.Method = services.SubmitBacklinksDataToStellar
	configs.QueueTransaction.Method = services.SubmitUserDataToStellar

	// Initialize the cron scheduler with Delay a job's execution if the previous run hasn't completed yet
	c := cron.New(
		cron.WithChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
			cron.Recover(cron.DefaultLogger),
		),
	)

	c.AddFunc("@every 30m", func() {
		notificationhandler.CheckStellarAccountBalance(commons.GoDotEnvVariable("NFTSTELLARISSUERPUBLICKEYK"))
		notificationhandler.CheckStellarAccountBalance(commons.GoDotEnvVariable("SPONSORERPK"))
	})

	c.AddFunc("@every 12h", func() {
		services.CheckTestimonialStatus()
		services.CheckOrganizationStatus()
	})

	c.AddFunc("@every 30s", func() {
		services.QueueScheduleWorkers()
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
	// initial log file when server starts
	utilities.CreateLogFile()
	// create logger
	logger := utilities.NewCustomLogger()
	logger.LogWriter("Gateway Started @port "+port+" with "+envName+" environment", 1)

	http.ListenAndServe(port, handlers.CORS(originsOk, headersOk, methodsOk)(router))
}
