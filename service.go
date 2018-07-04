package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// StringService _
type StringService interface {
	Uppercase(context.Context, string) (string, error)
	Count(context.Context, string) int
}

type stringService struct{}

func (stringService) Uppercase(ctx context.Context, str string) (string, error) {
	if str == "" {
		return "", errors.New("Empty string")
	}
	return strings.ToUpper(str), nil
}

func (stringService) Count(ctx context.Context, str string) int {
	return len(str)
}

// Model expected client request for uppercase resource
type uppercaseRequest struct {
	Str string `json:"str"`
}

// Model server response for uppercase resource
type uppercaseResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"` // Errors don't JSON-Marshal, so we use a string
}

// Model expected client request for count resource
type countRequest struct {
	Str string `json:"str"`
}

// Model server response for count resource
type countReponse struct {
	Count int `json:"count"`
}

// Parse request body into a uppercaseRequest object
func decodeUppercaseRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req uppercaseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// Parse request body into a countRequest object
func decodeCountRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req countRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// We can use one encode procedure for both calls, since we only need to encode the provided response into a json object
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
