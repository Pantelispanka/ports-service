package redis_test

import (
	"ports-service/internal/infrustructure/repository/redis"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPortRepo_InvalidRedisUrl(t *testing.T) {
	_, err := redis.NewPortRepo("wrong_schema")
	assert.Error(t, err)
	assert.Equal(t,
		"redis: invalid URL scheme: ", err.Error())
}
