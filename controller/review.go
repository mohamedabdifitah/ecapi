package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohamedabdifitah/ecapi/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserReviews(c *gin.Context) {
	id := c.Param("id")
	reviews, err := db.GetReviewsByUser(id)
	if err != nil {
		c.JSON(400, err)
	}
	c.JSON(200, reviews)
}
func GetReviewMerchant(c *gin.Context) {
	ExternalId := c.Param("id")
	// bson.D{{"merchant_review.external_id", "64ea387416182c259943067b"}}
	reviews, err := db.GetReviewsToInstance("merchant_review.external_id", ExternalId)
	if err != nil {
		c.JSON(500, err)
	}
	c.JSON(200, reviews)

}
func GetReviewDriver(c *gin.Context) {
	ExternalId := c.Param("id")
	reviews, err := db.GetReviewsToInstance("driver_review.external_id", ExternalId)
	if err != nil {
		c.JSON(500, err)
	}
	c.JSON(200, reviews)

}
func GetAllReview(c *gin.Context) {
	reviews, err := db.GetAllReviews(bson.D{})
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
		From: body.From,
		MerchantReview: db.ReviewColl{
			Rate:       body.MerchantReview.Rate,
			Message:    body.MerchantReview.ExternalId,
			ExternalId: body.MerchantReview.ExternalId,
		},
		OrderId: body.OrderId,
		DriverReview: db.ReviewColl{
			Rate:       body.DriverReview.Rate,
			Message:    body.MerchantReview.ExternalId,
			ExternalId: body.MerchantReview.ExternalId,
		},
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
		Id: objecId,
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
