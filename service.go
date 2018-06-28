package main

import (
	"context"
	"errors"
	"strings"
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
