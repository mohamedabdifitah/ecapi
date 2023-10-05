package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/meilisearch/meilisearch-go"
	"github.com/mohamedabdifitah/ecapi/service"
)

func Search(c *gin.Context) {
	var text string = c.Query("text")
	fmt.Println(text)
	hints, err := service.Melli.Index("data").Search(text, &meilisearch.SearchRequest{})
	if err != nil {
		c.String(400, err.Error())
		return
	}
	c.JSON(200, hints)
}
