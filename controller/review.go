package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohamedabdifitah/ecapi/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserReviews(c *gin.Context) {
	id := c.Param("id")
	filter := bson.D{
		{Key: "from", Value: id},
	}
	reviews, err := db.GetReviews(filter)
	if err != nil {
		c.JSON(400, err)
		return
	}
	c.JSON(200, reviews)
}
func GetReviewToInstace(Type string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ExternalId := c.Param("id")
		filter := bson.D{
			{Key: "Type", Value: Type},
			{Key: "external_id", Value: ExternalId},
		}
		reviews, err := db.GetReviews(filter)
		if err != nil {
			c.JSON(400, err)
			c.Abort()
			return
		}
		c.JSON(200, reviews)
		c.Abort()
	}
}

func GetAllReview(c *gin.Context) {
	reviews, err := db.GetReviews(bson.D{})
	if err != nil {
		c.JSON(500, err)
		return
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
		From:       body.From,
		Rate:       body.Rate,
		Message:    body.Message,
		ExternalId: body.ExternalId,
		OrderId:    body.OrderId,
		Options:    body.Options,
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
	review, err := db.GetReviewById(objecId)
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
	filter := bson.M{"_id": objecId}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "message", Value: body.Message},
		{Key: "rate", Value: body.Rate},
		{Key: "options", Value: body.Options},
		{Key: "metadata.updated_at", Value: time.Now().UTC()},
	}}}
	res, err := db.ReviewUpdate(filter, update)
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
	res, err := db.DeleteReview(objecId)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, res)
}
