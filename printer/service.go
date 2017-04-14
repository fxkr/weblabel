package printer

import (
	"context"

	"github.com/go-kit/kit/log"
)

type Service interface {

	// Returns nil if the service is operational.
	Status(ctx context.Context) error

	// Sends a document to the printer.
	Print(ctx context.Context, req printRequest) error
}

type service struct {
	printer Printer
	logger  log.Logger
}

func NewService(printer Printer, logger log.Logger) Service {
	return &service{
		printer,
		logger,
	}
}

func (s *service) Status(ctx context.Context) error {
	return nil
}

func (s *service) Print(ctx context.Context, req printRequest) error {
	return s.printer.Text(req.Text)
}
