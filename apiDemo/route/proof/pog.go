package proof

import (
	"net/http"

	"github.com/dileepaj/tracified-gateway/apiDemo/controller"
	"github.com/dileepaj/tracified-gateway/apiDemo/middleware"
	"github.com/gorilla/mux"
)

func PogRoutes(r *mux.Router) {
	a:=[]int{1, 2}
	r.HandleFunc("/transaction/genesis", middleware.HeaderReader(middleware.Authentication(a, controller.SubmitGenesis))).Methods(http.MethodGet)
}
