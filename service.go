package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
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

/*
type Midleware func(endpoint.Endpoint) endpoint.Endpoint

func loggingMidleware(logger log.Logger) Midleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint { // return a Middleware
		return func(ctx context.Context, request interface{}) (interface{}, error) { // Return an endpoint
			logger.Log("msj", "calling endpoint")
			defer logger.Log("msj", "called endpoint")
			return next(ctx, request) // return an endpoint (call the original one)
		}
	}
}*/

type loggingMiddleware struct {
	logger log.Logger
	next   StringService
}

func (mw loggingMiddleware) Uppercase(ctx context.Context, s string) (output string, err error) {
	defer func(t time.Time) {
		mw.logger.Log(
			"method", "uppercase",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(t),
		)
	}(time.Now())
	return mw.next.Uppercase(ctx, s)
}

func (mw loggingMiddleware) Count(ctx context.Context, s string) (n int) {
	defer func(t time.Time) {
		mw.logger.Log(
			"method", "count",
			"input", s,
			"n", n,
			"took", time.Since(t),
		)
	}(time.Now())
	n = mw.next.Count(ctx, s)
	return
}

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	cuntResult     metrics.Histogram
	next           StringService
}

func (im instrumentingMiddleware) Uppercase(ctx context.Context, s string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "uppercase", "error", fmt.Sprint(err != nil)}
		im.requestCount.With(lvs...).Add(1)
		im.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output, err = im.next.Uppercase(ctx, s)
	return
}

func (im instrumentingMiddleware) Count(ctx context.Context, s string) (n int) {
	defer func(begin time.Time) {
		lvs := []string{"method", "uppercase", "error", "false"}
		im.requestCount.With(lvs...).Add(1)
		im.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		im.cuntResult.Observe(float64(n))
	}(time.Now())
	n = im.next.Count(ctx, s)
	return
}
