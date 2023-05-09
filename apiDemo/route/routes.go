package route

import (
	proofroute "github.com/dileepaj/tracified-gateway/apiDemo/route/proof"
	"github.com/gorilla/mux"
)

func DefineRoutes(r *mux.Router) {
	healthRoutes(r)
	proofroute.PogRoutes(r)
}
