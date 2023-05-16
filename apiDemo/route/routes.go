package route

import (
	proofroute "github.com/dileepaj/tracified-gateway/apiDemo/route/proof"
	"github.com/dileepaj/tracified-gateway/apiDemo/route/expertformula"
	"github.com/gorilla/mux"
)

func DefineRoutes(r *mux.Router) {
	healthRoutes(r)
	proofroute.PogRoutes(r)
	expertformula.ExpertFormulaRoutes(r)
}
