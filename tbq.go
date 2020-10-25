package tbqueue

//go:generate mockgen -source=tbq.go -destination=./source_mock_test.go -package=tbqueue

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fpapadopou/tbq/config"
	"github.com/fpapadopou/tbq/redis"
)

const (
	consumerInterval  = 500
	consumerRetries   = 3
	consumerRetryWait = 1000
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
type ProcessorFunc func(ctx context.Context, item interface{}) error

// Publish sends items to the queue along with the time at which they should be processed.
func (q *TBQ) Publish(ctx context.Context, item interface{}, processTime time.Time) error {

	timestamp := processTime.UTC().Unix()

	err := q.s.Send(ctx, item, timestamp)
	if err != nil {
		log.Printf("tbqueue.Publish failed: %v", err)
		return fmt.Errorf("tbqueue.Publish: failed to publish item  ")
	}

	return nil
}

// Consume starts a consumer job for the current queue.
func (q *TBQ) Consume(ctx context.Context) error {

	var err error

	for i := 1; i <= consumerRetries; i++ {
		err = q.doConsume(ctx)
		if err == nil {
			return nil
		}

		if i <= consumerRetries {
			log.Printf("tbqueue.Consume: consume attempt %d/%d, retrying in %d milliseconds", i, consumerRetries, consumerRetryWait)
			time.Sleep(consumerRetryWait * time.Millisecond)
		}
	}

	return err
}

func (q *TBQ) doConsume(ctx context.Context) error {
	// Start a ticker consuming and processing one item every consumeInterval milliseconds.
	ticker := time.NewTicker(consumerInterval * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// TODO: Need to wait for all processorFunc instances to finish?
			log.Print("tbqueue.doConsume: stopping consumer..")
			return nil
		case <-ticker.C:
			// Consume next message and start processing it.
			item, err := q.s.Receive(ctx)
			if err != nil {
				return fmt.Errorf("tbqueue.Consume: error while receiving message: %v", err)
			}
			err = q.proc(ctx, item)
			if err != nil {
				return fmt.Errorf("tbqueue.Consume: error while processing message: %v", err)
			}
		}
	}
}

// New initializes and returns a new instance of TBQ using the specified processor function for the consumer.
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
			return TBQ{}, fmt.Errorf("tbqueue.New: failed to create source (redis): %v", err)
		}
		tbq.s = s
	default:
		return TBQ{}, fmt.Errorf("tbqueue.New: unrecognized source type: %s", c.SourceType)
	}

	return tbq, nil
}
