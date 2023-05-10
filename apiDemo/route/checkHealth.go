package route

import (
	"net/http"

	"github.com/dileepaj/tracified-gateway/apiDemo/controller"
	"github.com/dileepaj/tracified-gateway/apiDemo/middleware"
	"github.com/gorilla/mux"
)

var required = []int{1, 2}

func healthRoutes(r *mux.Router) {
	r.HandleFunc("/health", middleware.HeaderReader(controller.HealthCheck)).Methods(http.MethodGet)
	r.HandleFunc("/health/auth", middleware.HeaderReader(middleware.Authentication(required, controller.HealthCheck))).Methods(http.MethodGet)
}
