package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mohamedabdifitah/ecapi/controller"
	"github.com/mohamedabdifitah/ecapi/middleware"
)

type GroupRoute struct {
	prefix string
	routes []Route
}
type Route struct {
	method   string
	path     string
	handlers []gin.HandlerFunc
}

// TODO: implement validation of body , params and quaery and headers for preventing errors smooth expereicne for api users.
func (r *GroupRoute) register() {
	for _, route := range r.routes {
		switch route.method {
		case "GET":
			server.GET(r.prefix+route.path, route.handlers...)
			continue
		case "POST":
			server.POST(r.prefix+route.path, route.handlers...)
			continue
		case "PUT":
			server.PUT(r.prefix+route.path, route.handlers...)
			continue
		case "DELETE":
			server.DELETE(r.prefix+route.path, route.handlers...)
			continue
		case "OPTIONS":
			server.OPTIONS(r.prefix+route.path, route.handlers...)
			continue
		case "PATCH":
			server.PATCH(r.prefix+route.path, route.handlers...)
			continue
		case "HEAD":
			server.HEAD(r.prefix+route.path, route.handlers...)
			continue
		default:
			server.Any(r.prefix+route.path, route.handlers...)
			continue
		}

	}
}

var Routes []GroupRoute = []GroupRoute{
	{
		prefix: "/customer",
		routes: []Route{
			{
				method:   "GET",
				path:     "/all",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"admin"}), controller.GetAllCustomers},
			},
			{
				method:   "GET",
				path:     "/get/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{}), controller.GetCustomer},
			},
			{
				method:   "GET",
				path:     "/signin/google",
				handlers: []gin.HandlerFunc{controller.SiginWithGoogle},
			},
			{
				method:   "PUT",
				path:     "/update/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"admin", "customer"}), controller.UpdateCustomer},
			},
			{
				method:   "PATCH",
				path:     "/change/password",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"customer"}), controller.ChangeCustomerPassword},
			},
			{
				method:   "PATCH",
				path:     "/change/email",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"customer"}), controller.ChangeCustomerEmail},
			},
			{
				method:   "DELETE",
				path:     "/delete",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"customer"}), controller.DeleteCustomer},
			},
			{
				method:   "POST",
				path:     "/signup/email",
				handlers: []gin.HandlerFunc{controller.SignUpCustomerWithEmail},
			},
			{
				method:   "POST",
				path:     "/signin/email",
				handlers: []gin.HandlerFunc{controller.CustomerEmailLogin},
			},
			{
				method:   "PATCH",
				path:     "/change/device",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"customer"}), controller.ChangeCustomerDevice},
			},
			{
				method:   "PATCH",
				path:     "/change/webhooks",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"customer"}), controller.ChangeCustomerWebhooks},
			},
			{
				method:   "PATCH",
				path:     "/change/profile",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"customer"}), controller.ChangeCustomerDevice},
			},
		},
	},
	{
		prefix: "/review",
		routes: []Route{
			{
				method:   "GET",
				path:     "/get/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{}), controller.GetReviewById},
			},
			{
				method:   "GET",
				path:     "/user/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{}), controller.GetUserReviews},
			},
			{
				method:   "GET",
				path:     "/merchant/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{}), controller.GetReviewToInstace("REVIEW_MERCHANT")},
			},
			{
				method:   "GET",
				path:     "/driver/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{}), controller.GetReviewToInstace("REVIEW_DRIVER")},
			},
			{
				method:   "GET",
				path:     "/menu/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{}), controller.GetReviewToInstace("REVIEW_MENU")},
			},
			{
				method:   "GET",
				path:     "/all",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"admin"}), controller.GetAllReview},
			},
			{
				method:   "PUT",
				path:     "/update/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"customer"}), controller.UpdateReview},
			},
			{
				method:   "POST",
				path:     "/create",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"customer"}), controller.CreateReview},
			},
			{
				method:   "DELETE",
				path:     "/delete/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"customer", "admin"}), controller.DeleteReview},
			},
		},
	},
	{
		prefix: "/menu",
		routes: []Route{
			{
				method:   "GET",
				path:     "/get/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{}), controller.GetMenu},
			},
			{
				method:   "GET",
				path:     "/all",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{}), controller.GetMenus},
			},
			{
				method:   "GET",
				path:     "/merchant/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{}), controller.GetMenuFromMerchant},
			},
			{
				method:   "POST",
				path:     "/create",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"merchant"}), controller.CreateMenu},
			},
			{
				method:   "DELETE",
				path:     "/delete/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"merchant"}), controller.DeleteMenu},
			},
			{
				method:   "PUT",
				path:     "/update/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"merchant"}), controller.UpdateMenu},
			},
			{
				method:   "PATCH",
				path:     "/images/add/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"merchant"}), controller.AddMenuImages},
			},
		},
	},
	{
		prefix: "/order",
		routes: []Route{
			{
				method:   "GET",
				path:     "/all",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"admin"}), controller.GetAllOrders},
			},
			{
				method:   "GET",
				path:     "/get/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{}), controller.GetOrderByid},
			},
			{
				method:   "GET",
				path:     "/customer/all/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"customer", "admin"}), controller.GetOrderByCustomer},
			},
			{
				method:   "POST",
				path:     "/place",
				handlers: []gin.HandlerFunc{controller.PlaceOrder},
			},
			{
				method:   "GET",
				path:     "/location",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{}), controller.GetOrderByLocation},
			},
			{
				method:   "PATCH",
				path:     "/driver/accept/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"driver"}), controller.AccpetOrderByDriver},
			},
			{
				method:   "PATCH",
				path:     "/assign/:oid/:did",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"merchant", "admin"}), controller.AssignOrderToDriver},
			},
			{
				method:   "PATCH",
				path:     "/drop/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"driver", "admin"}), controller.DropOrderByDriver},
			},
			{
				method:   "PATCH",
				path:     "/merchant/decline/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"merchant"}), controller.RejectOrderByMerchant},
			},
			{
				method:   "PATCH",
				path:     "/merchant/accept/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"merchant"}), controller.MerchantOrderAccept},
			},
			{
				method:   "PATCH",
				path:     "/customer/cancel",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"customer", "admin"}), controller.CancelOrder},
			},
			{
				method:   "PATCH",
				path:     "/stage/prepare/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"merchant"}), controller.ChangeOrderStatus("preparing")},
			},
			{
				method:   "PATCH",
				path:     "/stage/pickuped/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"merchant"}), controller.ChangeOrderStatus("pickuped")},
			},
			{
				method:   "PATCH",
				path:     "/stage/ready/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"merchant"}), controller.ChangeOrderStatus("ready")},
			},
			{
				method:   "PATCH",
				path:     "/stage/delivered/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"driver"}), controller.OrderIsDelivered},
			},
			{
				method:   "GET",
				path:     "/merchant/all/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"merchant", "admin"}), controller.GetOrderByMerchant},
			},
			{
				method:   "GET",
				path:     "/driver/all/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"driver", "admin"}), controller.GetOrderByDriver},
			},
		},
	},
	{
		prefix: "/merchant",
		routes: []Route{
			{
				method:   "GET",
				path:     "/all",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{}), controller.GetAllMerchants},
			},
			{
				method:   "GET",
				path:     "/get/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{}), controller.GetMerchant},
			},
			{
				method:   "GET",
				path:     "/location",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{}), controller.GetMerchantByLocation},
			},
			{
				method:   "PATCH",
				path:     "/change/password/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"merchant"}), controller.ChangeMerchantPassword},
			},
			{
				method:   "PATCH",
				path:     "/change/phone/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"merchant"}), controller.ChangeMerchantPhone},
			},
			{
				method:   "PATCH",
				path:     "/change/device/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"merchant"}), controller.ChangeMerchantDevice},
			},
			{
				method:   "PATCH",
				path:     "/change/profile/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"merchant"}), controller.ChangeMerchantProfile},
			},
			{
				method:   "PATCH",
				path:     "/change/webhooks",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"merchant"}), controller.ChangeMerchantWebhooks},
			},
			{
				method:   "DELETE",
				path:     "/delete/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"merchant", "admin"}), controller.DeleteMerchant},
			},
			{
				method:   "PUT",
				path:     "/update/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"admin", "merchant"}), controller.UpdateMerchant},
			},
			{
				method:   "POST",
				path:     "/signup",
				handlers: []gin.HandlerFunc{controller.SignUpMerchant},
			},
			{
				method:   "POST",
				path:     "/signin/phone",
				handlers: []gin.HandlerFunc{controller.MerchantPhoneLogin},
			},
			{
				method:   "GET",
				path:     "/filter",
				handlers: []gin.HandlerFunc{controller.FilterMerchants},
			},
		},
	},
	{
		prefix: "/driver",
		routes: []Route{
			{
				method:   "GET",
				path:     "/all",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"merchant", "admin"}), controller.GetAllDrivers},
			},
			{
				method:   "GET",
				path:     "/get/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{}), controller.GetDriver},
			},
			{
				method:   "GET",
				path:     "/list",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{}), controller.GetListDrivers},
			},
			{
				method:   "GET",
				path:     "/location",
				handlers: []gin.HandlerFunc{controller.GetDrivesByLocation},
			},
			{
				method:   "PATCH",
				path:     "/change/password/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"driver"}), controller.ChangeDriverPassword},
			},
			{
				method:   "PATCH",
				path:     "/change/phone/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"driver"}), controller.ChangeDriverPhone},
			},
			{
				method:   "PATCH",
				path:     "/change/device/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"driver"}), controller.ChangeDriverDevice},
			},
			{
				method:   "PATCH",
				path:     "/change/location",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"driver"}), controller.ChangeDriverLocation},
			},
			{
				method:   "PATCH",
				path:     "/change/webhook",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"driver"}), controller.ChangeDriverWebhooks},
			},
			{
				method:   "PATCH",
				path:     "/change/profile/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"driver"}), controller.ChangeDriverProfile},
			},
			{
				method:   "PATCH",
				path:     "/change/email/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"driver"}), controller.ChangeDriverEmail},
			},
			{
				method:   "DELETE",
				path:     "/delete/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"driver", "admin"}), controller.DeleteDriver},
			},
			{
				method:   "PUT",
				path:     "/update/:id",
				handlers: []gin.HandlerFunc{middleware.AuthorizeRolesMiddleware([]string{"admin", "driver"}), controller.UpdateDriver},
			},
			{
				method:   "POST",
				path:     "/signup",
				handlers: []gin.HandlerFunc{controller.SignUpDriverWithEmail},
			},
			{
				method:   "POST",
				path:     "/signin",
				handlers: []gin.HandlerFunc{controller.DriverEmailLogin},
			},
		},
	},
	{
		prefix: "/",
		routes: []Route{
			{
				method:   "GET",
				path:     "search",
				handlers: []gin.HandlerFunc{controller.Search},
			},
			{
				method:   "GET",
				path:     "multisearch",
				handlers: []gin.HandlerFunc{controller.MultiSearch},
			},
			{
				method:   "GET",
				path:     "callback/google",
				handlers: []gin.HandlerFunc{controller.GoogleCallBack},
			},
		},
	},
	{
		prefix: "/upload",
		routes: []Route{
			{
				method:   "POST",
				path:     "/file",
				handlers: []gin.HandlerFunc{controller.UploadFile},
			},
		},
	},
	{
		prefix: "/verify",
		routes: []Route{
			{
				method:   "POST",
				path:     "/email",
				handlers: []gin.HandlerFunc{controller.VerifyOtpEmail},
			},
			{
				method:   "POST",
				path:     "/phone",
				handlers: []gin.HandlerFunc{controller.VerifyOtpPhone},
			},
		},
	},
}
