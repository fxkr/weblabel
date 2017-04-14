package printer

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

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

	r := mux.NewRouter()
	r.Handle("/printer/v1/status", statusHandler).Methods("GET")
	r.Handle("/printer/v1/print", printHandler).Methods("POST")

	return r
}

func decodeStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return statusRequest{}, nil
}

func decodePrintRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body printRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	switch err {
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
