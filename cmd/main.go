package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"ports-service/internal/infrustructure/repository/redis"
	portservice "ports-service/internal/service"
	"syscall"
)

// The main function
func main() {
	log.Println("Starting ...")
	terminateCh := make(chan os.Signal, 1)
	signal.Notify(terminateCh, syscall.SIGINT, syscall.SIGTERM)

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		fmt.Println("Please set REDIS_URL variable")
		os.Exit(1)
	}

	filePath := os.Getenv("FILE_PATH")
	if filePath == "" {
		fmt.Println("Please set FILE_PATH variable")
		os.Exit(1)
	}

	portRepo, err := redis.NewPortRepo(redisURL)

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Could not create port repo. Please check redis URL")
		os.Exit(1)
	}

	ctx := context.Background()
	stream := portservice.NewPortService(ctx, terminateCh, portRepo)
	count := make(chan int, 1)
	errorCounter := make(chan int, 1)
	go func(counter chan int, errorCounter chan int) {
		i := 0
		j := 0
		for data := range stream.Watch() {
			if data.Error != nil {
				j += 1
				select {
				case errorCounter <- j:

				default:
					select {
					case <-errorCounter:
						errorCounter <- j
					default:
					}
				}
				log.Println(data.Error)
			}
			i += 1
			select {
			case counter <- i: // channel was empty - ok

			default: // channel if full - we have to delete a value from it with some precautions to not get locked in our own channel
				select {
				case <-counter: // read stale value and put a fresh one
					counter <- i
				default: // consumer have read it - so skip and not get locked
				}
			}

		}
	}(count, errorCounter)

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	stream.Start(file)

	counted := <-count
	errorsCounted := <-errorCounter
	fmt.Printf("Upsert %d ports \n", counted)
	fmt.Printf("Error in %d ports \n", errorsCounted)
	log.Println("Program ended...")
}
