package controller

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/mohamedabdifitah/ecapi/db"
// )

// type Recommendations struct {
// 	title string        `json:"title"`
// 	Items []db.Merchant `json:"items"`
// }

// // order again => depent
// func ListCategoriesRecomendations(c *gin.Context) {
// 	var categories []string = []string{"fries", "restaurant", "chicken", "meat", "drink", "vegetable"}
// 	c.JSON(200, categories)
// 	// fries , chicken , resturant , drinks , grocessory , etc.
// }
// func PopularFilter(c *gin.Context) {
// 	var filters []string = []string{"popular", "near by", "new", "rating", "recent viewed", "order again", "favorites", "discount"}
// 	c.JSON(200, filters)
// }
// func FilterMerchant(c *gin.Context) {

// }
// func PopularMerchant(c *gin.Context) {
// 	// []Recommendations
// 	// like top 10 most popular restaurant
// 	// top 10 most ordered restaurant
// 	// top 10 news restaurant
// 	// top 10 most loved restaurant
// }
// func GetRecomended(c *gin.Context) {
// 	location := c.Query("location")
// 	menu := c.Query("menu")
// 	Popular := c.Query("popular")

// 	// filter by category
// 	// cache response
// }
// func SearchBy(c *gin.Context) {
// 	name := c.Query("name")
// 	description := c.Query("description")
// 	menu := c.Query("menu")
// }
