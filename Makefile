test: ## Run tests locally with race detector and test coverage
	go test ./... -race -cover

lint: ## Perform linting. Packages goimports and linter should be manually installed.
	go vet ./...
	goimports -w `find . -name '*.go'`
	golangci-lint run

docker-up: ## Spin up application and the requirements
	docker-compose -f ./docker-compose.yml up --force-recreate

docker-down: ## Stop the application and the requirements
	docker-compose -f ./docker-compose.yml down

docker-build-no-cache: ## Build the docker image of the application without caching 
	docker build --no-cache -t ports-service -f ./Dockerfile.multi . 

docker-build: ## Build the docker image of the application
	docker build -t ports-service -f ./Dockerfile.multi .