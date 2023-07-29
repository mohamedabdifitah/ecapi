package api

import (
	"github.com/mohamedabdifitah/ecapi/controller"
)

func RouterDefinition() {
	CustomerRouteDefinition()
}
func CustomerRouteDefinition() {
	CustomerRouter.GET("/all", controller.GetAllCustomers)

}
