package printer

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) Status(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Status",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Status(ctx)
}
