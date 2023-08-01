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
	CustomerRouter.PUT("/update/:id", controller.UpdateCustomer)
	CustomerRouter.PATCH("/change/password/:id", controller.ChangeCustomerPassword)
	CustomerRouter.DELETE("/delete/:id", controller.DeleteCustomer)
	CustomerRouter.POST("/singup/email", controller.SingUpCustomerWithEmail)
}
