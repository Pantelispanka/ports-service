package redis

import (
	"context"
	"encoding/json"
	"ports-service/internal/domain"

	"github.com/redis/go-redis/v9"
)

// The Port repo
type PortRepository struct {
	redisClient *redis.Client
}

// Create and ingect the Prt repositpry repo in other packages
func NewPortRepo(redisUrl string) (*PortRepository, error) {
	options, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)

	return &PortRepository{
		redisClient: client,
	}, nil
}

// Upsert a port
func (pr *PortRepository) UpsertPort(ctx context.Context, p domain.Port) error {
	key := p.Unloc
	value, err := json.Marshal(p)
	if err != nil {
		return err
	}
	err = pr.redisClient.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// Get a port by it's key, here the unloc
func (pr *PortRepository) GetPort(ctx context.Context, key string) (*domain.Port, error) {
	p, err := pr.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var port domain.Port
	err = json.Unmarshal(p, &port)
	if err != nil {
		return nil, err
	}
	return &port, nil
}
