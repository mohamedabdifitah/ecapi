package service

import (
	"fmt"
	"os"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	// "github.com/meilisearch/meilisearch-go/"
)

type Record struct {
	ObjectID string `json:"objectID"`
	Name     string `json:"name"`
}

var client *search.Client

func InitAlgolia() {
	// Connect and authenticate with your Algolia app
	client = search.NewClient(os.Getenv("ALGO_APPID"), os.Getenv("ALGO_API_KEY"))
	fmt.Println("Connected to Algolia")
}

func CreateDocument(name string, record interface{}) error {
	// Create a new index and add a record
	index := client.InitIndex(name)
	resSave, err := index.SaveObject(record)
	if err != nil {
		return err
	}
	err = resSave.Wait()
	if err != nil {
		return err
	}
	return nil

}
func MultipleSearchDocument(queries []search.IndexedQuery, options ...interface{}) (search.MultipleQueriesRes, error) {
	// Search the index and print the results
	results, err := client.MultipleQueries(queries, "", options...)
	if err != nil {
		return search.MultipleQueriesRes{}, err
	}
	return results, nil
}
func SearchDocument(name string, query string, options ...interface{}) (*search.QueryRes, error) {
	index := client.InitIndex(name)
	results, err := index.Search(query, options...)
	if err != nil {
		return nil, err
	}
	return &results, nil
}
