# setting developemnt flag 
export APP_ENV = development
$ENV:APP_ENV = "development"
run: 
	docker-compose up
dev:
	watcher
test:
	go test ./...