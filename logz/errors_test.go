package logz

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/ibrt/golang-errors/errorz"
	"github.com/ibrt/golang-fixtures/fixturez"
	"github.com/ibrt/golang-inject-clock/clockz"
	"github.com/ibrt/golang-inject-clock/clockz/testclockz"
	"github.com/stretchr/testify/require"
)

type TestError struct {
	message string
}

func (e *TestError) Error() string {
	return e.message
}

type TestErrorWithCause struct {
	message string
	err     error
}

func (e *TestErrorWithCause) Error() string {
	return e.message
}

func (e *TestErrorWithCause) Cause() error {
	return e.err
}

type TestErrorWithUnwrap struct {
	message string
	err     error
}

func (e *TestErrorWithUnwrap) Error() string {
	return e.message
}

func (e *TestErrorWithUnwrap) Unwrap() error {
	return e.err
}

type ErrorsSuite struct {
	*fixturez.DefaultConfigMixin
	Clock *testclockz.MockHelper
}

func TestErrors(t *testing.T) {
	fixturez.RunSuite(t, &ErrorsSuite{})
}

func (s *ErrorsSuite) TestErrorToSentryEvent(ctx context.Context, t *testing.T) {
	event := errorToSentryEvent(ctx, fmt.Errorf("test error"), Error)
	require.Len(t, event.Exception, 1)
	require.NotNil(t, event.Exception[0].Stacktrace)
	event.Exception[0].Stacktrace = nil
	require.Equal(t, &sentry.Event{
		Contexts:  make(map[string]interface{}),
		Tags:      make(map[string]string),
		Modules:   make(map[string]string),
		Extra:     map[string]interface{}{},
		Level:     sentry.LevelError,
		Timestamp: clockz.Get(ctx).Now(),
		Exception: []sentry.Exception{
			{
				Type:  "*errors.errorString",
				Value: "test error",
			},
		},
	}, event)

	event = errorToSentryEvent(ctx, errorz.Errorf("test error"), Warning)
	require.Len(t, event.Exception, 1)
	require.NotNil(t, event.Exception[0].Stacktrace)
	event.Exception[0].Stacktrace = nil
	require.Equal(t, &sentry.Event{
		Contexts:  make(map[string]interface{}),
		Tags:      make(map[string]string),
		Modules:   make(map[string]string),
		Extra:     map[string]interface{}{},
		Level:     sentry.LevelWarning,
		Timestamp: clockz.Get(ctx).Now(),
		Exception: []sentry.Exception{
			{
				Type:  "*errors.errorString",
				Value: "test error",
			},
		},
	}, event)

	event = errorToSentryEvent(ctx, errorz.Errorf("test error", errorz.ID("test-id"), errorz.Status(http.StatusBadRequest), errorz.M("k", "v")), Warning)
	require.Len(t, event.Exception, 1)
	require.NotNil(t, event.Exception[0].Stacktrace)
	event.Exception[0].Stacktrace = nil
	require.Equal(t, &sentry.Event{
		Contexts: make(map[string]interface{}),
		Tags:     make(map[string]string),
		Modules:  make(map[string]string),
		Extra: map[string]interface{}{
			"status": http.StatusBadRequest,
			"k":      "v",
		},
		Level:     sentry.LevelWarning,
		Timestamp: clockz.Get(ctx).Now(),
		Exception: []sentry.Exception{
			{
				Type:  "test-id",
				Value: "test error",
			},
		},
	}, event)

	tErr := &TestError{message: "test error"}
	event = errorToSentryEvent(ctx, errorz.Wrap(tErr), Warning)
	require.Len(t, event.Exception, 1)
	require.NotNil(t, event.Exception[0].Stacktrace)
	event.Exception[0].Stacktrace = nil
	require.Equal(t, &sentry.Event{
		Contexts:  make(map[string]interface{}),
		Tags:      make(map[string]string),
		Modules:   make(map[string]string),
		Extra:     map[string]interface{}{"unwrapped[0][*logz.TestError]": tErr},
		Level:     sentry.LevelWarning,
		Timestamp: clockz.Get(ctx).Now(),
		Exception: []sentry.Exception{
			{
				Type:  "*logz.TestError",
				Value: "test error",
			},
		},
	}, event)

	tErr = &TestError{message: "inner test error"}
	tErrC := &TestErrorWithCause{message: "outer test error", err: tErr}
	event = errorToSentryEvent(ctx, errorz.Wrap(tErrC), Warning)
	require.Len(t, event.Exception, 2)
	require.NotNil(t, event.Exception[0].Stacktrace)
	event.Exception[0].Stacktrace = nil
	require.Equal(t, &sentry.Event{
		Contexts: make(map[string]interface{}),
		Tags:     make(map[string]string),
		Modules:  make(map[string]string),
		Extra: map[string]interface{}{
			"unwrapped[0][*logz.TestErrorWithCause]": tErrC,
			"unwrapped[1][*logz.TestError]":          tErr,
		},
		Level:     sentry.LevelWarning,
		Timestamp: clockz.Get(ctx).Now(),
		Exception: []sentry.Exception{
			{
				Type:  "*logz.TestErrorWithCause",
				Value: "outer test error",
			},
			{
				Type:  "*logz.TestError",
				Value: "inner test error",
			},
		},
	}, event)

	tErr = &TestError{message: "inner test error"}
	tErrU := &TestErrorWithUnwrap{message: "outer test error", err: tErr}
	event = errorToSentryEvent(ctx, errorz.Wrap(tErrU), Warning)
	require.Len(t, event.Exception, 2)
	require.NotNil(t, event.Exception[0].Stacktrace)
	event.Exception[0].Stacktrace = nil
	require.Equal(t, &sentry.Event{
		Contexts: make(map[string]interface{}),
		Tags:     make(map[string]string),
		Modules:  make(map[string]string),
		Extra: map[string]interface{}{
			"unwrapped[0][*logz.TestErrorWithUnwrap]": tErrU,
			"unwrapped[1][*logz.TestError]":           tErr,
		},
		Level:     sentry.LevelWarning,
		Timestamp: clockz.Get(ctx).Now(),
		Exception: []sentry.Exception{
			{
				Type:  "*logz.TestErrorWithUnwrap",
				Value: "outer test error",
			},
			{
				Type:  "*logz.TestError",
				Value: "inner test error",
			},
		},
	}, event)
}
