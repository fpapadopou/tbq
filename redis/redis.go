package redis

import (
	"context"
	config "github.com/fpapadopou/tbq/config"
)

// Redis holds all the functionality that's needed in order to interact with a Redis connection.
type Redis struct{}

// Send implements the respective method of the Source interface.
func (r *Redis) Send(ctx context.Context, item interface{}) error {

	return nil
}

// Receive implements the respective method of the Source interface.
func (r *Redis) Receive(ctx context.Context) (interface{}, error) {

	return nil, nil
}

// New creates a new Redis object using the provided configuration.
func New(conf config.Config) (*Redis, error) {

	return &Redis{}, nil
}
