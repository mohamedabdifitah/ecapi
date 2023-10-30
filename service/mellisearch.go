package service

import (
	"os"

	"github.com/meilisearch/meilisearch-go"
)

var Melli *meilisearch.Client

func InitMelliClient() {
	Melli = meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   "http://localhost:7700",
		APIKey: os.Getenv("MELLI_API_KEY"),
	})
}
func AddDocument(index string, data interface{}) error {
	_, err := Melli.Index(index).AddDocuments(data)
	if err != nil {
		return err
	}
	return nil
}
func Search(index string, query string, SearchRequest meilisearch.SearchRequest) (*meilisearch.SearchResponse, error) {
	res, err := Melli.Index(index).Search(query, &SearchRequest)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// example
//
//	[]meilisearch.SearchRequest{
//		{
//			IndexUID: "movies",
//			Query:    "pooh",
//			Limit:    5,
//		},
//		{
//			IndexUID: "movies",
//			Query:    "nemo",
//			Limit:    5,
//		},
//		{
//			IndexUID: "movie_ratings",
//			Query:    "us",
//		},
//	},
func MultiSearch(queries []meilisearch.SearchRequest) (*meilisearch.MultiSearchResponse, error) {
	res, err := Melli.MultiSearch(&meilisearch.MultiSearchRequest{
		Queries: queries,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
