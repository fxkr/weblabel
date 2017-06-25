package printer

import (
	"context"
	"encoding/json"
	"image/png"
	"net/http"
	"strings"

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

	printHandler := kithttp.NewServer(
		makePrintEndpoint(s),
		decodePrintRequest,
		encodeResponse,
		opts...,
	)

	printImageHandler := kithttp.NewServer(
		makePrintImageEndpoint(s),
		decodePrintImageRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()
	r.Handle("/api/v1/printer/status", statusHandler).Methods("GET")
	r.Handle("/api/v1/printer/print", printHandler).Methods("POST")
	r.Handle("/api/v1/printer/image", printImageHandler).Methods("POST")

	return r
}

func decodeStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return statusRequest{}, nil
}

func decodePrintRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body printRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, errors.Wrap(ErrBadRequest, err.Error())
	}

	return body, nil
}

func decodePrintImageRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body printImageRequest

	r.ParseMultipartForm(2 << 20)

	err := json.NewDecoder(strings.NewReader(r.FormValue("data"))).Decode(&body)
	if err != nil {
		return nil, errors.Wrap(ErrBadRequest, err.Error())
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		return nil, errors.Wrap(ErrBadRequest, err.Error())
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, errors.Wrap(ErrBadRequest, err.Error())
	}

	body.Image = img

	return body, nil
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
