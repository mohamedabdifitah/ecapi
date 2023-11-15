package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetReviewById(id primitive.ObjectID) (*Review, error) {
	var review *Review
	query := bson.M{"_id": id}
	result := ReviewCollection.FindOne(Ctx, query)
	err := result.Decode(&review)
	if err != nil {
		return nil, err
	}
	return review, nil
}

// get all the reviews
func GetReviews(query bson.D) ([]*Review, error) {
	var reviews []*Review
	cursor, err := ReviewCollection.Find(Ctx, query)
	if err != nil {
		return nil, err
	}
	for cursor.Next(Ctx) {
		var review *Review
		err := cursor.Decode(&review)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	cursor.Close(Ctx)
	return reviews, nil
}
func (r *Review) Create() (*mongo.InsertOneResult, error) {
	r.Metadata.CreatedAt = time.Now().UTC()
	res, err := ReviewCollection.InsertOne(Ctx, &r)
	if err != nil {
		return nil, err
	}
	// go service.ProduceMessage("", "review_create", "", r)
	return res, nil
}

// Update review
func ReviewUpdate(filter bson.M, update bson.D) (*mongo.UpdateResult, error) {
	result, err := ReviewCollection.UpdateOne(Ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// delte review
func DeleteReview(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": id}
	result, err := ReviewCollection.DeleteOne(Ctx, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Calculate all score and return five 5-star rating score i.e 2.9 ,3.7 etc.
func CalculateRate(scores []int) int {
	var ScoreTotal int = scores[4]*5 + scores[3]*4 + scores[2]*3 + scores[1]*2 + scores[0]*1
	var ResponseTotal int = scores[4] + scores[3] + scores[2] + scores[1] + scores[0]
	return ScoreTotal / ResponseTotal
}
func ReviewMenu(id primitive.ObjectID) {
	_ = bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "rate",
					Value: bson.M{},
				},
			},
		},
		{
			Key: "$inc",
			Value: bson.D{
				{Key: "rate.participants", Value: 1},
			},
		},
	}

}
