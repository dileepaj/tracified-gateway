package proof

import (
	"github.com/dileepaj/tracified-gateway/apiDemo/controller"
	"github.com/dileepaj/tracified-gateway/apiDemo/model"
)

var ProofRoutes = model.Routers{
	model.Router{
		Name:    "Submit Genesis XDR",
		Method:  "POST",
		Path:    "/api/health",
		Handler: controller.SubmitGenesis,
	},
}
