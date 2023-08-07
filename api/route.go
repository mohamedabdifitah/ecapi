package api

import (
	"github.com/mohamedabdifitah/ecapi/controller"
	"github.com/mohamedabdifitah/ecapi/middleware"
)

func RouterDefinition() {
	CustomerRouteDefinition()
	MerchantRouteDefinition()
	DriverRouteDefinition()
	MenuRouterDefinition()
}
func CustomerRouteDefinition() {
	CustomerRouter.GET("/all", middleware.AuthorizeRolesMiddleware([]string{"admin", "manager"}), controller.GetAllCustomers)
	CustomerRouter.GET("/get/:id", middleware.AuthorizeRolesMiddleware([]string{}), controller.GetCustomer)
	CustomerRouter.PUT("/update/:id", controller.UpdateCustomer)
	CustomerRouter.PATCH("/change/password/:id", controller.ChangeCustomerPassword)
	CustomerRouter.PATCH("/change/email/:id", controller.ChangeCustomerEmail)
	CustomerRouter.DELETE("/delete/:id", controller.DeleteCustomer)
	CustomerRouter.POST("/signup/email", controller.SignUpCustomerWithEmail)
	CustomerRouter.POST("/signin/email", controller.CustomerEmailLogin)
}
func MerchantRouteDefinition() {
	MerchantRouter.GET("/all", controller.GetAllMerchants)
	MerchantRouter.GET("/get/:id", controller.GetMerchant)
	MerchantRouter.PUT("/update/:id", controller.UpdateMerchant)
	MerchantRouter.PATCH("/change/password/:id", controller.ChangeMerchantPassword)
	MerchantRouter.PATCH("/change/phone/:id", controller.ChangeMerchantPhone)
	MerchantRouter.DELETE("/delete/:id", controller.DeleteMerchant)
	MerchantRouter.POST("/signup/phone", controller.SignUpMerchantWithPhone)
	MerchantRouter.POST("/signin/phone", controller.MerchantPhoneLogin)
}
func DriverRouteDefinition() {
	DriverRouter.GET("/all", controller.GetAllDrivers)
	DriverRouter.GET("/get/:id", controller.GetDriver)
	DriverRouter.PUT("/update/:id", controller.UpdateDriver)
	DriverRouter.PATCH("/change/password/:id", controller.ChangeDriverPassword)
	DriverRouter.PATCH("/change/phone/:id", controller.ChangeDriverPhone)
	DriverRouter.PATCH("/change/email/:id", controller.ChangeDriverEmail)
	DriverRouter.DELETE("/delete/:id", controller.DeleteDriver)
	DriverRouter.POST("/signup", controller.SignUpDriverWithPhone)
	DriverRouter.POST("/signin/phone", controller.DriverPhoneLogin)
}
func MenuRouterDefinition() {
	MenuRouter.GET("/get/:id", controller.GetMenu)
	MenuRouter.GET("/all", controller.GetMenus)
	MenuRouter.POST("/create", controller.CreateMenu)
	MenuRouter.PUT("/update/:id", controller.UpdateMenu)
}
