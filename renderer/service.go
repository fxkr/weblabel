package renderer

import (
	"context"
	"image"

	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

type Service interface {

	// Returns nil if the service is operational.
	Status(ctx context.Context) error

	// Render an image based on parameters.
	Render(ctx context.Context, doc Document) (image.Image, error)
}

type service struct {
	renderer Renderer
	logger   log.Logger
}

type Document struct {
	Text string `json:"text"`
}

func NewService(renderer Renderer, logger log.Logger) Service {
	return &service{
		renderer,
		logger,
	}
}

func (s *service) Status(ctx context.Context) error {
	return nil
}

func (s *service) Render(ctx context.Context, doc Document) (image.Image, error) {
	img, err := s.renderer.Render(doc.Text)
	err = errors.Wrap(err, "Failed to render image")
	return img, err
}
