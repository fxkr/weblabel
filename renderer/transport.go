package renderer

import (
	"bytes"
	"context"
	"encoding/json"
	"image"
	"image/png"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
)

func MakeHandler(ctx context.Context, s Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	statusHandler := kithttp.NewServer(
		makeStatusEndpoint(s),
		decodeStatusRequest,
		encodeResponse,
		opts...,
	)

	renderHandler := kithttp.NewServer(
		makeRenderEndpoint(s),
		decodeRenderRequest,
		encodeImageResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/api/v1/renderer/status", statusHandler).Methods("GET")
	r.Handle("/api/v1/renderer/render", renderHandler).Methods("POST")

	return r
}

func decodeStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return statusRequest{}, nil
}

func decodeRenderRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body renderRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, errors.Wrap(ErrBadRequest, err.Error())
	}

	return body, nil
}

func encodeImageResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	img, ok := response.(image.Image)
	if !ok {
		return encodeResponse(ctx, w, response) // fallback
	}

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, img); err != nil {
		return errors.Wrap(err, "Failed to encode PNG image")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

	if _, err := w.Write(buffer.Bytes()); err != nil {
		return errors.Wrap(err, "Failed to write image to client")
	}

	return nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, errors.WithStack(e.error()), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return errors.WithStack(json.NewEncoder(w).Encode(response))
}

type errorer interface {
	error() error
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	switch errors.Cause(err) {
	case ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
	case ErrBadRequest:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
