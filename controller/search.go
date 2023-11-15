package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/meilisearch/meilisearch-go"
	"github.com/mohamedabdifitah/ecapi/service"
)

var menu_searchable_attributes []string = []string{"title", "description", "category", "reciepe", "price"}
var merchant_searchable_attributes []string = []string{"business_name", "generes", "description", "address"}
var merchant_retivable_attributes []string = []string{"id", "business_name", "generes", "description", "location", "address", "closed", "likes"}
var menu_retivable_attributes []string = []string{"id", "title", "description", "category", "reciepe", "merchant_external_id", "price", "service_availablity"}

func Search(c *gin.Context) {
	var query string = c.Query("query")
	var index string = c.Query("index")
	var searchable_attributes []string
	var retivable_attributes []string
	switch index {
	case "merchant":
		searchable_attributes = merchant_searchable_attributes
		retivable_attributes = merchant_retivable_attributes
	case "menu":
		searchable_attributes = menu_searchable_attributes
		retivable_attributes = menu_retivable_attributes
	}
	var body meilisearch.SearchRequest = meilisearch.SearchRequest{
		AttributesToSearchOn: searchable_attributes,
		AttributesToRetrieve: retivable_attributes,
	}
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
	var body struct {
		Menu     string `json:"menu"`
		Merchant string `json:"merchant"`
	}
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(200, err)
		return
	}
	query := []meilisearch.SearchRequest{
		{
			IndexUID:             "menu",
			Query:                body.Menu,
			AttributesToRetrieve: menu_retivable_attributes,
			AttributesToSearchOn: menu_searchable_attributes,
		},
		{

			IndexUID:             "merchant",
			Query:                body.Merchant,
			AttributesToRetrieve: merchant_retivable_attributes,
			AttributesToSearchOn: merchant_searchable_attributes,
		},
	}
	res, err := service.MultiSearch(query)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	c.JSON(200, res)
}
