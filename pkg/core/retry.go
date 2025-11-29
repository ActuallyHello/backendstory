package core

import (
	"log/slog"
	"time"
)

func Retry(operation string, fn func() (any, error), options ...RetryOption) (any, error) {
	retryOptions := RetryOptions{
		maxRetries: 10,
		delay:      1 * time.Second,
	}

	for _, option := range options {
		option(&retryOptions)
	}

	var err error
	var res any
	isSuccess := false

	for i := 0; i < retryOptions.maxRetries; i++ {
		res, err = fn()
		if err != nil {
			slog.Warn("Retry operation!", "operation", operation, "attempt", i, "err", err)
			<-time.After(retryOptions.delay)
			continue
		}
		isSuccess = true
		break
	}

	if !isSuccess {
		slog.Error("Retry operations failed!", "operation", operation, "err", err)
		return nil, err
	}
	return res, nil
}

type RetryOptions struct {
	maxRetries int
	delay      time.Duration
}

type RetryOption func(*RetryOptions)

func SetMaxRetriesOpt(maxRetries int) RetryOption {
	return func(ro *RetryOptions) {
		ro.maxRetries = maxRetries
	}
}

func SetMaxDelayOpt(delay time.Duration) RetryOption {
	return func(ro *RetryOptions) {
		ro.delay = delay
	}
}
