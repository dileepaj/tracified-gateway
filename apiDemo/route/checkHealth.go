package route

import (
	"net/http"

	"github.com/dileepaj/tracified-gateway/apiDemo/controller"
	"github.com/dileepaj/tracified-gateway/apiDemo/middleware"
	"github.com/gorilla/mux"
)


func healthRoutes(r *mux.Router){
	r.HandleFunc("/health", middleware.JSONMiddleware(controller.HealthCheck)).Methods(http.MethodGet)
}