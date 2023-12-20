package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"ports-service/internal/domain"
)

// The port repository interface. This should implement the two functions documented here for dependency injection.
// For example we can pass multiple implementations if needed. Here we work with an in memory database. This
// coud be any other database like Mongo or an SQL db.
type PortRepository interface {
	UpsertPort(ctx context.Context, p domain.Port) error
	GetPort(ctx context.Context, key string) (*domain.Port, error)
}

type PortService struct {
	portRepo PortRepository
}

// Create new service and pass the dependencies
func NewPortService(pr PortRepository) *PortService {
	return &PortService{
		portRepo: pr,
	}
}

// Entry represents each stream. If the stream fails, an error will be present.
type Entry struct {
	Error error
	Port  domain.Port
}

// Stream helps transmit each streams withing a channel.
type Stream struct {
	portRepo PortRepository
	stream   chan Entry
	context  context.Context
	shutdown chan os.Signal
}

// NewJSONStream returns a new `Stream` type.
func NewJSONStream(ctx context.Context, shutdown chan os.Signal, pr PortRepository) Stream {
	return Stream{
		portRepo: pr,
		stream:   make(chan Entry),
		context:  ctx,
		shutdown: shutdown,
	}
}

// Watch watches JSON streams. Each stream entry will either have an error or a
// Port object. Client code does not need to explicitly exit after catching an
// error as the `Start` method will close the channel automatically.
func (s Stream) Watch() <-chan Entry {
	return s.stream
}

func (s Stream) Start(reader io.Reader) {
	// Stop streaming channel as soon as nothing left to read in the file.
	defer close(s.stream)

	decoder := json.NewDecoder(reader)

	// Read opening delimiter. `[` or `{`
	if _, err := decoder.Token(); err != nil {
		s.stream <- Entry{Error: fmt.Errorf("decode opening delimiter: %w", err)}
		return
	}

	// Read file content as long as there is something.
	for decoder.More() {
		// Check for termination signal
		select {
		case <-s.shutdown:
			return // Gracefully terminate
		default:
			// Continue processing
		}
		t, err := decoder.Token()
		if err != nil {
			s.stream <- Entry{Error: fmt.Errorf("decode error: %w", err)}
			return
		}
		key := t.(string)

		var port domain.Port
		if err := decoder.Decode(&port); err != nil {
			s.stream <- Entry{Error: fmt.Errorf("decode error: %w", err)}
			return
		}

		port.Unloc = key
		s.stream <- Entry{Port: port}

		err = port.Validate()
		if err != nil {
			fmt.Printf("WARNING: failed to validate port with code: %s", port.Code)
			fmt.Printf("WARNING: %s", err.Error())
			return
		}
		err = s.portRepo.UpsertPort(s.context, port)
		if err != nil {
			fmt.Printf("WARNING: failed to upsert port: %v", err)
			continue
		}
	}

	// Read closing delimiter. `]` or `}`
	if _, err := decoder.Token(); err != nil {
		s.stream <- Entry{Error: fmt.Errorf("decode closing delimiter: %w", err)}
		return
	}

}
