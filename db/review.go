package db

import (
	"fmt"
	"time"

	"github.com/mohamedabdifitah/ecapi/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *Review) GetById() error {
	query := bson.M{"_id": r.Id}
	result := ReviewCollection.FindOne(Ctx, query)
	err := result.Decode(&r)
	return err
}

// get all the reviews user made for order , merchant or driver
func GetReviewsByUser(id string) ([]*Review, error) {
	var reviews []*Review
	filter := bson.D{
		{Key: "from", Value: id},
	}
	cursor, err := ReviewCollection.Find(Ctx, filter)
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

// reviews to particular instance like
// all review of merchant 1
func GetReviewsToInstance(Type string, id string) ([]*Review, error) {
	var reviews []*Review
	filter := bson.D{
		{Key: Type, Value: id},
	}
	cursor, err := ReviewCollection.Find(Ctx, filter)
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

// get all the reviews
func GetAllReviews(query bson.D) ([]*Review, error) {
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
	// rate := math.Round(r.DriverReview.Rate)
	objectid, err := primitive.ObjectIDFromHex(r.MerchantReview.ExternalId)
	if err != nil {
		return nil, fmt.Errorf("merchant id is not valid: %v", err)
	}
	merchant := Merchant{
		Id: objectid,
	}

	err = merchant.GetById()
	if err != nil {
		return nil, fmt.Errorf("merchant id is not found")
	}
	response, err := UpdateMerchant(bson.M{"_id": merchant.Id}, bson.D{{Key: "$set", Value: bson.D{
		{
			Key:   "rate.stats.$[]",
			Value: 1,
		},
	},
	}})
	if err != nil {
		return nil, err
	}
	fmt.Println(response)
	return res, nil
}

// Update review
func (r *Review) Update() (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": r.Id}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "metadata.updated_at", Value: time.Now().UTC()},
	}}}
	result, err := ReviewCollection.UpdateOne(Ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// delte review
func (r *Review) Delete() (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": r.Id}
	result, err := ReviewCollection.DeleteOne(Ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = service.PublishTopic("review", r); err != nil {
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
