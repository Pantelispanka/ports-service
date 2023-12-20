test: ## Run tests locally with race detector and test coverage
	go test ./... -race -cover

lint: ## Perform linting. Packages goimports and linter should be manually installed.
	go vet ./...
	goimports -w `find . -name '*.go'`
	golangci-lint run