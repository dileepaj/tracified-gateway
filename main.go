package main

import (

	"fmt"
	"github.com/dileepaj/tracified-gateway/api/routes"
	"net/http"
	"github.com/robfig/cron"
	"github.com/gorilla/handlers"
	"github.com/dileepaj/tracified-gateway/services"

)

func main() {
	
	headersOk := handlers.AllowedHeaders([]string{"Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	c := cron.New()
	c.AddFunc("@every 30s", func() { 
		services.CheckCOCStatus()
	})
	c.Start()

	router := routes.NewRouter()	
	fmt.Println("Gateway Started @port :8030...!")
	http.ListenAndServe(":8030", handlers.CORS(originsOk, headersOk, methodsOk)(router))

}
