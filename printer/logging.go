package printer

import (
	"context"
	"fmt"
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
			"err", fmt.Sprintf("%+v", err),
		)
	}(time.Now())
	return s.Service.Status(ctx)
}

func (s *loggingService) Print(ctx context.Context, req printRequest) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "Print",
			"took", time.Since(begin),
			"err", fmt.Sprintf("%+v", err),
		)
	}(time.Now())
	return s.Service.Print(ctx, req)
}
