package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *Review) GetById() error {
	query := bson.M{"_id": r.Id}
	result := ReviewCollection.FindOne(Ctx, query)
	err := result.Decode(&r)
	return err
}

// get all the reviews user review
func (r *Review) GetByUser() ([]*Review, error) {
	var reviews []*Review
	filter := bson.D{
		{Key: "from", Value: r.From},
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
func (r *Review) GetReviewsToMe() ([]*Review, error) {
	var reviews []*Review
	filter := bson.D{
		{Key: "type", Value: r.Type},
		{Key: "external_id", Value: r.ExternalId},
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
func (r *Review) GetAll() ([]*Review, error) {
	var reviews []*Review
	cursor, err := ReviewCollection.Find(Ctx, bson.D{})
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
	return res, nil
}

// Update review
func (r *Review) Update() (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": r.Id}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "rate", Value: r.Rate},
		{Key: "message", Value: r.Message},
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
	return result, nil
}
