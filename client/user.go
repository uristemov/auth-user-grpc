package client

import (
	"context"
	"github.com/uristemov/auth-user-grpc/models"
)

type Config struct {
	Address  string
	Protocol string
	Insecure bool
}

type Client interface {
	Connect() error
	Close() error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}
