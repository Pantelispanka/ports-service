package redis_test

import (
	"context"
	"fmt"
	"ports-service/internal/domain"
	"ports-service/internal/infrustructure/repository/redis"

	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	// "github.com/ory/dockertest"
	// "github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/assert"
)

var (
	redisContainer testcontainers.Container
	redisRepo      *redis.PortRepository
)

func setupContainer(t *testing.T) (func(), error) {
	ctx := context.Background()

	// Define the Redis container configuration
	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	// Create and start the Redis container
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start Redis container: %w", err)
	}

	// Get the Redis container's host and port
	redisHost, err := redisC.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Redis container host: %w", err)
	}
	redisPort, err := redisC.MappedPort(ctx, "6379")
	if err != nil {
		return nil, fmt.Errorf("failed to get Redis container port: %w", err)
	}

	// Create the RedisPortRepository with the container's host and port
	redisURL := fmt.Sprintf("redis://%s:%s/0", redisHost, redisPort.Port())
	repo, err := redis.NewPortRepo(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis repository: %w", err)
	}

	// Set the global variables for the container and repository
	redisContainer = redisC
	redisRepo = repo

	// Create a cleanup function to terminate the container after the tests
	cleanup := func() {
		err := redisContainer.Terminate(ctx)
		if err != nil {
			t.Errorf("failed to terminate Redis container: %v", err)
		}
	}

	return cleanup, nil
}

func TestRedisRepo(t *testing.T) {
	clean, err := setupContainer(t)
	defer clean()
	ctx := context.Background()
	port := domain.Port{
		Unloc: "AEAUH",
		City:  "Abu Dhabi",
		Code:  "52001",
	}
	err = redisRepo.UpsertPort(ctx, port)
	if err != nil {
		fmt.Printf(err.Error())
	}
	assert.NoError(t, err)

	p, err := redisRepo.GetPort(ctx, "AEAUH")
	fmt.Println(p.City)

}
