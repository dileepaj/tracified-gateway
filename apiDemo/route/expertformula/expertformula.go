package expertformula

import (
	"net/http"

	"github.com/dileepaj/tracified-gateway/apiDemo/controller"
	"github.com/dileepaj/tracified-gateway/apiDemo/middleware"
	"github.com/gorilla/mux"
)

var required = []int{1, 2}

func ExpertFormulaRoutes(r *mux.Router) {
	r.HandleFunc("/socialimapact/expertformula", middleware.HeaderReader(middleware.Authentication(required, controller.SocialImpactExpertFormula))).Methods(http.MethodGet)
}
