# setting developemnt flag 
export APP_ENV = development
$ENV:APP_ENV = "development"
db:
	meilisearch --master-key 4eh9oD6BRbTftjnsjBCd0SOO3jGmO_-x6ZAQo6Mbr3c
dev:
	go run .
test:
	go test ./...