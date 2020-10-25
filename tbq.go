package tbqueue

//go:generate mockgen -source=tbq.go -destination=./source_mock_test.go -package=tbqueue

import (
	"context"
	"fmt"
	"time"

	"github.com/fpapadopou/tbq/config"
	"github.com/fpapadopou/tbq/redis"
)

// TBQ holds the connection to the actual source and is responsible for sending/receiving messages to/from the source.
type TBQ struct {
	s    Source
	proc ProcessorFunc
}

// Source defines the interface that a queue implementation must satisfy in order
// to be used by TBQ.
type Source interface {
	Send(context.Context, interface{}, int64) error
	Receive(context.Context) (interface{}, error)
}

// ProcessorFunc specifies how each received message
type ProcessorFunc func(item interface{}) error

// Publish sends items to the queue along with the time at which they should be processed.
func (q *TBQ) Publish(ctx context.Context, item interface{}, processTime time.Time) error {

	return nil
}

// Consume starts a consumer job for the current queue.
func (q *TBQ) Consume(ctx context.Context) error {

	return nil
}

// New initializes and returns a new instance of TBQ.
func New(proc ProcessorFunc) (TBQ, error) {

	c, err := config.New()
	if err != nil {
		return TBQ{}, fmt.Errorf("failed to load config: %v", err)
	}

	tbq := TBQ{
		proc: proc,
	}

	switch c.SourceType {
	case "redis":
		s, err := redis.New(c)
		if err != nil {
			return TBQ{}, fmt.Errorf("failed to create source (redis): %v", err)
		}
		tbq.s = s
	default:
		return TBQ{}, fmt.Errorf("failed to load config: %v", err)
	}

	return tbq, nil
}
