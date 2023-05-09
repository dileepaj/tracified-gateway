package routes

import (
	"github.com/dileepaj/tracified-gateway/apiDemo/model"
	"github.com/dileepaj/tracified-gateway/apiDemo/routes/proof"
)

var ApplicationRoutes model.Routers

func init() {
	routes := []model.Routers{
		healthRoutes,
		proof.ProofRoutes,
	}

	for _, r := range routes {
		ApplicationRoutes = append(ApplicationRoutes, r...)
	}
}
