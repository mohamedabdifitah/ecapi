package api

import (
	"github.com/mohamedabdifitah/ecapi/controller"
	"github.com/mohamedabdifitah/ecapi/middleware"
)

func RouterDefinition() {

	CustomerRouteDefinition()
}
func CustomerRouteDefinition() {
	CustomerRouter.GET("/all", middleware.AuthorizeRolesMiddleware([]string{"admin", "manager"}), controller.GetAllCustomers)
	CustomerRouter.GET("/get/:id", middleware.AuthorizeRolesMiddleware([]string{}), controller.GetCustomer)
	CustomerRouter.PUT("/update/:id", controller.UpdateCustomer)
	CustomerRouter.PATCH("/change/password/:id", controller.ChangeCustomerPassword)
	CustomerRouter.PATCH("/change/email/:id", controller.ChangeCustomerEmail)
	CustomerRouter.DELETE("/delete/:id", controller.DeleteCustomer)
	CustomerRouter.POST("/singup/email", controller.SingUpCustomerWithEmail)
	CustomerRouter.POST("/signin/email", controller.CustomerEmailLogin)
}
