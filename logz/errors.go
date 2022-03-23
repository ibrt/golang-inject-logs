package logz

import (
	"context"
	"fmt"
	"reflect"
	"runtime"

	"github.com/getsentry/sentry-go"
	"github.com/ibrt/golang-errors/errorz"
	"github.com/ibrt/golang-inject-clock/clockz"
)

const (
	maxErrorDepth = 10
	statusKey     = "status"
	unwrappedKey  = "unwrapped[%v][%v]"
)

func errorToSentryEvent(ctx context.Context, err error, level Level) *sentry.Event {
	err = errorz.Wrap(err, errorz.SkipPackage()) // ensure error is wrapped to start
	event := sentry.NewEvent()
	event.Timestamp = clockz.Get(ctx).Now()
	event.Level = level.toSentry()

	if status := errorz.GetStatus(err); status != 0 {
		event.Extra[statusKey] = status.Int()
	}

	for k, v := range errorz.GetMetadata(err) {
		event.Extra[k] = v
	}

	for i := 0; i < maxErrorDepth && err != nil; i++ {
		uErr := errorz.Unwrap(err)
		uType := reflect.TypeOf(uErr).String()

		event.Exception = append(event.Exception, sentry.Exception{
			Type: func() string {
				if id := errorz.GetID(err); id != "" {
					return id.String()
				}
				return uType
			}(),
			Value:      err.Error(),
			Stacktrace: extractSentryStacktrace(err),
		})

		err = uErr

		switch uType {
		case "*errors.errorString":
		default:
			event.Extra[fmt.Sprintf(unwrappedKey, i, reflect.TypeOf(err).String())] = err
		}

		switch pErr := err.(type) {
		case interface{ Unwrap() error }:
			err = pErr.Unwrap()
		case interface{ Cause() error }:
			err = pErr.Cause()
		default:
			err = nil
		}
	}

	return event
}

func extractSentryStacktrace(err error) *sentry.Stacktrace {
	if errorz.Unwrap(err) == err {
		return sentry.ExtractStacktrace(err)
	}
	return callersToSentryStacktrace(errorz.GetCallers(err))
}

func callersToSentryStacktrace(callers []uintptr) *sentry.Stacktrace {
	s := &sentry.Stacktrace{
		Frames: make([]sentry.Frame, 0),
	}

	callersFrames := runtime.CallersFrames(callers)

	for {
		callerFrame, more := callersFrames.Next()
		s.Frames = append([]sentry.Frame{sentry.NewFrame(callerFrame)}, s.Frames...)

		if !more {
			break
		}
	}

	return s
}
