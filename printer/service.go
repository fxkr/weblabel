package printer

import (
	"context"
	"image"

	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"

	"github.com/fxkr/weblabel/renderer"
)

type Service interface {

	// Returns nil if the service is operational.
	Status(ctx context.Context) error

	// Sends a document to the printer.
	Print(ctx context.Context, req printRequest) error

	// Directly sends raster graphics to the printer.
	PrintImage(ctx context.Context, img image.Image) error
}

type service struct {
	printer  Printer
	renderer renderer.Service
	logger   log.Logger
}

func NewService(printer Printer, renderer renderer.Service, logger log.Logger) Service {
	return &service{
		printer,
		renderer,
		logger,
	}
}

func (s *service) Status(ctx context.Context) error {
	return nil
}

func (s *service) Print(ctx context.Context, req printRequest) error {

	img, err := s.renderer.Render(ctx, renderer.Document{Text: req.Text})
	if err != nil {
		return errors.Wrap(err, "Failed to render label to file")
	}

	return s.printer.Image(img)
}

func (s *service) PrintImage(ctx context.Context, img image.Image) error {
	return s.printer.Image(img)
}
