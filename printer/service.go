package printer

import (
	"context"

	"github.com/go-kit/kit/log"
)

type Service interface {

	// Returns nil if the service is operational.
	Status(ctx context.Context) error
}

type service struct {
	logger log.Logger
}

func NewService(logger log.Logger) Service {
	return &service{
		logger,
	}
}

func (s *service) Status(ctx context.Context) error {
	return nil
}
