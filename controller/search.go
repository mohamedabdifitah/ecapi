package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/meilisearch/meilisearch-go"
	"github.com/mohamedabdifitah/ecapi/service"
)

func Search(c *gin.Context) {
	var query string = c.Query("query")
	var index string = c.Query("index")
	var body meilisearch.SearchRequest
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(400, err.Error())
		c.Abort()
		return
	}
	res, err := service.Search(index, query, body)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	c.JSON(200, &res)
}
func MultiSearch(c *gin.Context) {
	var queries []meilisearch.SearchRequest
	err := c.ShouldBindJSON(&queries)
	if err != nil {
		c.JSON(200, err)
		return
	}
	// query := []meilisearch.SearchRequest{
	// 	{
	// 		IndexUID: "movies",
	// 		Query:    "pooh",
	// 		Limit:    5,
	// 	},
	// 	{
	// 		IndexUID: "movies",
	// 		Query:    "nemo",
	// 		Limit:    5,
	// 	},
	// 	{
	// 		IndexUID: "movie_ratings",
	// 		Query:    "us",
	// 	},
	// }
	res, err := service.MultiSearch(queries)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	c.JSON(200, res)
}
