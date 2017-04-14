package printer

import (
	"context"

	"github.com/go-kit/kit/endpoint"
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
		return statusResponse{err}, nil
	}
}
