package proof

import (
	"net/http"

	"github.com/dileepaj/tracified-gateway/apiDemo/controller"
	"github.com/dileepaj/tracified-gateway/apiDemo/middleware"
	"github.com/gorilla/mux"
)

func PogRoutes(r *mux.Router) {
	r.HandleFunc("/pog", middleware.JSONMiddleware(middleware.AuthenticationMiddleware(controller.HealthCheck))).Methods(http.MethodGet)
}
