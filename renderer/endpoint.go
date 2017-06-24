package renderer

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
)

type statusRequest struct {
}

type statusResponse struct {
	Err error `json:"error,omitempty"`
}

func (r statusResponse) error() error { return r.Err }

func makeStatusEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(statusRequest)
		err := s.Status(ctx)
		err = errors.WithStack(err)
		return statusResponse{err}, nil
	}
}

type renderRequest struct {
	Document Document `json:"document"`
}

type errorResponse struct {
	Err error `json:"error,omitempty"`
}

func (r errorResponse) error() error { return r.Err }

func makeRenderEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(renderRequest)
		img, err := s.Render(ctx, req.Document)
		err = errors.WithStack(err)
		if err != nil {
			return errorResponse{err}, nil
		}
		return img, nil
	}
}
