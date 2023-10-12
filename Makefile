# setting developemnt flag 
export APP_ENV = development
$ENV:APP_ENV = "development"
dev:
	watcher
test:
	go test ./...