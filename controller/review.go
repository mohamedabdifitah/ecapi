package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedabdifitah/ecapi/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserReview(c *gin.Context) {
	id := c.Param("id")
	review := db.Review{
		From: id,
	}
	reviews, err := review.GetByUser()
	if err != nil {
		c.JSON(500, err)
	}
	c.JSON(200, reviews)
}
func GetReviewToMe(c *gin.Context) {
	Type := c.Param("type")
	ExternalId := c.Param("eid")
	review := db.Review{
		ExternalId: ExternalId,
		Type:       Type,
	}
	reviews, err := review.GetReviewsToMe()
	if err != nil {
		c.JSON(500, err)
	}
	c.JSON(200, reviews)

}
func GetAllReview(c *gin.Context) {
	review := db.Review{}
	reviews, err := review.GetAll()
	if err != nil {
		c.JSON(500, err)
	}
	c.JSON(200, reviews)
}
func CreateReview(c *gin.Context) {
	var body *ReviewBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	review := &db.Review{
		Rate:       body.Rate,
		From:       body.From,
		Message:    body.Message,
		Type:       body.Type,
		ExternalId: body.ExternalId,
	}
	res, err := review.Create()
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(201, res)
}
func GetReviewById(c *gin.Context) {
	var id string = c.Param("id")
	objecId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, "invalid id")
		return
	}
	review := &db.Review{
		Id: objecId,
	}
	err = review.GetById()
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			c.String(200, "review not found")
			return
		}
		c.String(200, err.Error())
		return
	}
	c.JSON(200, review)
}
func UpdateReview(c *gin.Context) {
	var body *ReviewBody
	var id string = c.Param("id")
	objecId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, "invalid id")
		return
	}
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	review := &db.Review{
		Id:      objecId,
		Rate:    body.Rate,
		Message: body.Message,
	}
	res, err := review.Update()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if res.MatchedCount > 0 && res.ModifiedCount == 0 {

		c.String(400, "review is same as original")
		return
	}
	if res.MatchedCount == 0 && res.ModifiedCount == 0 {
		c.String(200, "review not found")
		return
	}
	c.String(200, "review updated successfully")
}
func DeleteReview(c *gin.Context) {
	var id string = c.Param("id")
	objecId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.String(400, "Invalid id")
		return
	}
	review := db.Review{
		Id: objecId,
	}
	res, err := review.Delete()
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, res)
}
