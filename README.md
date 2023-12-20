## A service for ports

A Makefile is provided to establish the specified requirements of the project. 
To avoid docker build and other commands together with the testing of the application, The user can use make commands 
to test and run the service.

### commands
``make test `` -> test with -race and -cover flags.

-race flag should always be used. Specifically when goroutines access the same variable concurrently
-cover flag returns the code coverage of the tests created

``make lint`` Runs format on the code, checks the imports and the does the linting. 

> **Warning**
> goimports and golangci-lint should be manually installed

``make docker-up`` Spins up the docker containers for the service. Here a redis in memory database is required and the builded docker image for the port service

``make docker-down`` Stops the application and the dependencies 

``make docker-build-no-cache`` Build the application using the flag --no-cache in case needed.

``make docker-build`` Builds the application using cached layers of the image.


