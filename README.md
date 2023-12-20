## A service for ports

  

A Makefile is provided to establish the specified requirements of the project.

To avoid docker build and other commands together with the testing of the application, The user can use make commands

to test and run the service.

  

### commands

``make test `` -> test with -race and -cover flags.

  

-race flag should always be used. Specifically when goroutines access the same variable concurrently

-cover flag returns the code coverage of the tests created

  

``make lint`` Runs format on the code, checks the imports and the does the linting.

  

>  **Warning**

> goimports and golangci-lint should be manually installed

  

``make docker-up`` Spins up the docker containers for the service. Here a redis in memory database is required and the builded docker image for the port service

  

``make docker-down`` Stops the application and the dependencies

  

``make docker-build-no-cache`` Build the application using the flag --no-cache in case needed.

  

``make docker-build`` Builds the application using cached layers of the image.

  
  

### Run the app

  

```make docker-build```

Build the application

  

```make docker-up```

Spins up the application and the dependencies


## Docs

### main
main function is found in cmd folder.

### Service Domain

In the domain the model required to parse the json and handle it in go is found. Along with a validator and the corresponding unit tests.

### Infra

Here the database is found at corresponding repositories. We used redis as described in requirements of the project.
Here mainly integration tests have been implemented and small unit tests to increase the code coverage of the tests.

### Service

The service that handles everything. The dependencies are injected through a repo interface. The main dependency is the redis repo that upserts the entity.  For the service test we mock the dependencies and test the functionality intended and described by the requirements of the project. The json is parsed as a stream and not all together to ensure that the application won't break. A stream is created with the dependencies that handled everything and is followed in main. 

If an entity doesn't validate then the service will continue to parse the json and write the rest on the database while it will print out a log of the failed entity.
Before each parsing checks the signal to see if a terminate signal has been received.



>  **Warning**

> To run locally exposed env variables should be present on the system for golang to retreive them. 
The REDIS_URL and the FILE_PATH are critical for the application to run. 