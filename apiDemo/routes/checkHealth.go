package routes

import (
	"github.com/dileepaj/tracified-gateway/apiDemo/controller"
	"github.com/dileepaj/tracified-gateway/apiDemo/model"
)

// This routes use to check the API status
var healthRoutes = model.Routers{
	model.Router{
		Name:    "Connection test API",
		Method:  "GET",
		Path:    "/transaction/genesis",
		Handler: controller.HealthCheck,
	},
}
