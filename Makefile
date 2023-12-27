# Makefile
# setting developemnt flag 
export APP_ENV = development
$ENV:APP_ENV = "development"
dev:
	go run .
test:
	go test ./...