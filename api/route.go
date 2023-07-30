package api

import (
	"github.com/mohamedabdifitah/ecapi/controller"
)

func RouterDefinition() {
	CustomerRouteDefinition()
}
func CustomerRouteDefinition() {
	CustomerRouter.GET("/all", controller.GetAllCustomers)
	CustomerRouter.GET("/get/:id", controller.GetCustomer)
	CustomerRouter.POST("/singup/email", controller.SingUpCustomerWithEmail)
}
