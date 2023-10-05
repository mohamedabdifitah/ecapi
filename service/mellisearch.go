package service

import (
	"github.com/meilisearch/meilisearch-go"
)

var Melli *meilisearch.Client

func InitMelliClient() {
	Melli = meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   "http://localhost:7700",
		APIKey: "aX0s_rDilOB0l1yl_AnwaQYcauOzEBeQ24WVb9Lyrrg",
	})
}
