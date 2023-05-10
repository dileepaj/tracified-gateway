package proof

import (
	"net/http"

	"github.com/dileepaj/tracified-gateway/apiDemo/controller"
	"github.com/dileepaj/tracified-gateway/apiDemo/middleware"
	"github.com/gorilla/mux"
)

var required = []int{1, 2}

func PogRoutes(r *mux.Router) {
	r.HandleFunc("/transaction/genesis", middleware.HeaderReader(middleware.Authentication(required, controller.SubmitGenesis))).Methods(http.MethodGet)
}
