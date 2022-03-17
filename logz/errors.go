package logz

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/getsentry/sentry-go"
	"github.com/ibrt/golang-errors/errorz"
)

const (
	maxErrorDepth = 10
	statusKey     = "status"
	unwrappedKey  = "unwrapped[%v][%v]"
)

func errorToSentryEvent(err error, level Level) *sentry.Event {
	err = errorz.Wrap(err, errorz.Skip()) // ensure error is wrapped to start
	event := sentry.NewEvent()
	event.Level = level.toSentry()

	if status := errorz.GetStatus(err); status != 0 {
		event.Extra[statusKey] = status
	}

	for k, v := range errorz.GetMetadata(err) {
		event.Extra[k] = v
	}

	for i := 0; i < maxErrorDepth && err != nil; i++ {
		event.Exception = append(event.Exception, sentry.Exception{
			Type: func() string {
				if id := errorz.GetID(err); id != "" {
					return id.String()
				}
				return reflect.TypeOf(err).String()
			}(),
			Value:      err.Error(),
			Stacktrace: extractSentryStacktrace(err),
		})

		err = errorz.Unwrap(err)

		// TODO(ibrt): Skip common errors.
		event.Extra[fmt.Sprintf(unwrappedKey, i, reflect.TypeOf(err).String())] = err

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
