package db

import "go.mongodb.org/mongo-driver/bson"

func GetCategories(query bson.D) ([]*Category, error) {
	var categories []*Category
	cursor, err := ReviewCollection.Find(Ctx, query)
	if err != nil {
		return nil, err
	}
	for cursor.Next(Ctx) {
		var category *Category
		err := cursor.Decode(&category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	cursor.Close(Ctx)
	return categories, nil
}
func GetCategoryMenues(category string) ([]*Menu, error) {
	query := bson.D{
		{
			Key: "category", Value: category,
		},
	}
	menues, err := GetMenues(query)
	if err != nil {
		return nil, err
	}
	return menues, nil
}
