package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

var ErrEmpty = errors.New("Empty string")

type StringService interface {
	Uppercase(context.Context, string) (string, error)
	Count(context.Context, string) int
}

type stringService struct {
}

func (stringService) Uppercase(ctx context.Context, s string) (string, error) {

	if s == "" {
		return "", ErrEmpty
	}

	return strings.ToUpper(s), nil
}

func (stringService) Count(ctx context.Context, s string) int {
	return len(s)
}

type uppercaseRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}

type countRequest struct {
	S string `json:"s"`
}

type countResponse struct {
	V int `json:"v"`
}

func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(uppercaseRequest)
		v, err := svc.Uppercase(ctx, req.S)
		if err != nil {
			return uppercaseResponse{v, err.Error()}, nil
		}
		return uppercaseResponse{v, ""}, nil
	}
}

func makeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		v := svc.Count(ctx, req.S)
		return countResponse{v}, nil
	}

}

func main() {

	svc := stringService{}
	uppercaseHandler := httptransport.NewServer(makeUppercaseEndpoint(svc), decodeUppercaseRequest, encodeResponse)
	countHandler := httptransport.NewServer(makeCountEndpoint(svc), decodeCountRequest, encodeResponse)
	http.Handle("/uuppeercase", uppercaseHandler)
	http.Handle("/count", countHandler)
}

func decodeUppercaseRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {

		log.Println(err)
		return nil, err
	}
	return request, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
