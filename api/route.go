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
	ReviewRouterDefinition()
	OrderRouterDefinition()
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
	CustomerRouter.PATCH("/device/change/:id", controller.ChangeCustomerDevice)
	CustomerRouter.PATCH("/change/profile/:id", controller.ChangeCustomerProfile)

}
func MerchantRouteDefinition() {
	MerchantRouter.GET("/all", controller.GetAllMerchants)
	MerchantRouter.GET("/get/:id", controller.GetMerchant)
	MerchantRouter.GET("/location", controller.GetMerchantByLocation)
	MerchantRouter.PUT("/update/:id", controller.UpdateMerchant)
	MerchantRouter.PATCH("/change/password/:id", controller.ChangeMerchantPassword)
	MerchantRouter.PATCH("/change/phone/:id", controller.ChangeMerchantPhone)
	MerchantRouter.DELETE("/delete/:id", controller.DeleteMerchant)
	MerchantRouter.POST("/signup/phone", controller.SignUpMerchantWithPhone)
	MerchantRouter.POST("/signin/phone", controller.MerchantPhoneLogin)
	MerchantRouter.PATCH("/device/change/:id", controller.ChangeMerchantDevice)
	MerchantRouter.PATCH("/change/profile/:id", controller.ChangeMerchantProfile)
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
	DriverRouter.PATCH("/device/change/:id", controller.ChangeDriverDevice)
	DriverRouter.PATCH("/change/profile/:id", controller.ChangeDriverProfile)
}
func MenuRouterDefinition() {
	MenuRouter.GET("/get/:id", controller.GetMenu)
	MenuRouter.GET("/all", controller.GetMenus)
	MenuRouter.GET("/merchant/:id", controller.GetMenuFromMerchant)
	MenuRouter.POST("/create", controller.CreateMenu)
	MenuRouter.DELETE("/delete/:id", controller.DeleteMenu)
	MenuRouter.PUT("/update/:id", controller.UpdateMenu)
	MenuRouter.PUT("/image/add/:id", controller.PutImageMenues)
}

func ReviewRouterDefinition() {
	ReviewRouter.GET("/get/:id", controller.GetReviewById)
	ReviewRouter.GET("/user/:id", controller.GetUserReview)
	ReviewRouter.GET("/to/:type/:eid", controller.GetReviewToMe)
	ReviewRouter.GET("/all", controller.GetAllReview)
	ReviewRouter.PUT("/update/:id", controller.UpdateReview)
	ReviewRouter.POST("/create", controller.CreateReview)
	ReviewRouter.DELETE("/delete/:id", controller.DeleteReview)

}
func OrderRouterDefinition() {
	OrderRouter.GET("/all", controller.GetAllOrders)
	OrderRouter.GET("/get/:id", controller.GetOrderByid)
	OrderRouter.POST("/place", controller.PlaceOrder)
	OrderRouter.GET("/location")       // polygons , longitude and latitude are
	OrderRouter.POST("/driver/accept") // driver accepts to deliver the request of order
	OrderRouter.POST("/change/driver")
	OrderRouter.POST("/change/merchant")
	OrderRouter.POST("/merchant/decline")
	OrderRouter.POST("/merchant/accept")
	OrderRouter.PATCH("/cancel")
	OrderRouter.GET("/merchant/get/all/:id")
	OrderRouter.GET("/driver/get/all/:id")
}
