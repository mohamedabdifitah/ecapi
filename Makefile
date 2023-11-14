# Makefile
# setting developemnt flag 
export APP_ENV = development
$ENV:APP_ENV = "development"
db:
	cd C:\Melli && meilisearch --master-key hAHhoSkUf1KM2STEW0X5wAo755pMnr6DiRqYrQ1d3H8
dev:
	nodemon --exec go run . --signal SIGTERM
test:
	go test ./...