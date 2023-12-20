package service_test

import (
	"bytes"
	"context"
	"os"
	"os/signal"
	"ports-service/internal/domain"
	"ports-service/internal/service"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

var invalidJson = `sfgsfg`

var portsSample = `{
	"VNDAD": {
		"name": "Da Nang",
		"city": "Da Nang",
		"country": "Viet Nam",
		"alias": [],
		"regions": [],
		"coordinates": [55.0272904, 24.9857145],
		"province": "Da Nang",
		"timezone": "Asia/Saigon",
		"unlocs": [
		  "VNDAD"
		],
		"code": "55204"
	  },
	  "VNHPH": {
		"name": "Haiphong",
		"city": "Haiphong",
		"country": "Viet Nam",
		"alias": [],
		"regions": [],
		"coordinates": [52.6126027, 24.1915137],
		"province": "Haiphong",
		"timezone": "Asia/Saigon",
		"unlocs": [
		  "VNHPH"
		],
		"code": "55201"
	  },
	  "VNNHA": {
		"name": "Nha Trang",
		"city": "Nha Trang",
		"country": "Viet Nam",
		"alias": [],
		"regions": [],
		"coordinates": [
		  1009.1967488,
		  12.2387911
		],
		"province": "Khanh Hoa Province",
		"timezone": "Asia/Saigon",
		"unlocs": [
		  "VNNHA"
		],
		"code": "55208"
	  }
}`

type mockPortRepository struct {
	ports map[string]domain.Port
}

func (m *mockPortRepository) GetPort(_ context.Context, key string) (*domain.Port, error) {
	port, exists := m.ports[key]
	if !exists {
		return nil, nil
	}
	return &port, nil
}

func (m *mockPortRepository) UpsertPort(_ context.Context, port domain.Port) error {
	m.ports[port.Unloc] = port
	return nil
}

func TestPortService_Success(t *testing.T) {

	repo := &mockPortRepository{
		ports: make(map[string]domain.Port),
	}

	ctx := context.Background()

	stream := service.NewPortService(ctx, nil, repo)

	reader := bytes.NewReader([]byte(portsSample))
	go func(t *testing.T) {
		for data := range stream.Watch() {
			assert.NoError(t, data.Error)
		}
	}(t)

	stream.Start(reader)
	port, err := repo.GetPort(ctx, "VNDAD")
	assert.NoError(t, err)
	assert.NotNil(t, port)
	assert.Equal(t, "Da Nang", port.Name)
	assert.Equal(t, "Da Nang", port.City)
}

func TestPortService_Invalid(t *testing.T) {
	repo := &mockPortRepository{
		ports: make(map[string]domain.Port),
	}

	ctx := context.Background()

	stream := service.NewPortService(ctx, nil, repo)

	reader := bytes.NewReader([]byte(portsSample))
	go func(t *testing.T) {
		for data := range stream.Watch() {
			assert.NoError(t, data.Error)
		}
	}(t)

	stream.Start(reader)
	port, err := repo.GetPort(ctx, "VNNHA")
	assert.NoError(t, err)
	assert.Nil(t, port)
}

func TestPortService_Terminate(t *testing.T) {
	repo := &mockPortRepository{
		ports: make(map[string]domain.Port),
	}

	ctx := context.Background()

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)

	terminate <- syscall.SIGTERM

	stream := service.NewPortService(ctx, terminate, repo)

	reader := bytes.NewReader([]byte(portsSample))
	go func(t *testing.T) {
		for data := range stream.Watch() {
			assert.NoError(t, data.Error)
		}
	}(t)

	stream.Start(reader)
	port, err := repo.GetPort(ctx, "VNDAD")
	assert.NoError(t, err)
	assert.Nil(t, port)
}
