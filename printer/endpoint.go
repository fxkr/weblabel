package printer

import (
	"context"
	"image"

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

type printRequest struct {
	Text string `json:"text"`
}

type printResponse struct {
	Err error `json:"error,omitempty"`
}

func (r printResponse) error() error { return r.Err }

func makePrintEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(printRequest)
		err := s.Print(ctx, req)
		err = errors.WithStack(err)
		return printResponse{err}, nil
	}
}

type printImageRequest struct {
	Image image.Image `json:"-"`
}

type printImageResponse struct {
	Err error `json:"error,omitempty"`
}

func makePrintImageEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(printImageRequest)
		err := s.PrintImage(ctx, req.Image)
		err = errors.WithStack(err)
		return printImageResponse{err}, nil
	}
}
