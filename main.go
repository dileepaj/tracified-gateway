package main

import (
	// "fmt"

	"fmt"
	// "log"
	"github.com/tracified-gateway/api/routes"
	"net/http"
	// "github.com/robfig/cron"
	"github.com/gorilla/handlers"

)

func main() {
	
	headersOk := handlers.AllowedHeaders([]string{"Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// c := cron.New()
	// // c.AddFunc("0 30 * * * *", func() { fmt.Println("Every hour on the half hour") })
	// // c.AddFunc("@hourly",      func() { fmt.Println("Every hour") })
	// c.AddFunc("@every 30s", func() { fmt.Println("Every thirty") })
	// c.Start()
	router := routes.NewRouter()	
	fmt.Println("Gateway Started @port :8030...!")
	http.ListenAndServe(":8030", handlers.CORS(originsOk, headersOk, methodsOk)(router))

}
