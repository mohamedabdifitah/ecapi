package controller

import (
	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/gin-gonic/gin"
	"github.com/mohamedabdifitah/ecapi/service"
)

func Search(c *gin.Context) {
	var text string = c.Query("text")
	queries := []search.IndexedQuery{
		search.NewIndexedQuery("merchant", opt.Query(text)),
		search.NewIndexedQuery("menu", opt.Query(text)),
	}
	res, err := service.MultipleSearchDocument(queries)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	c.JSON(200, res)
}
