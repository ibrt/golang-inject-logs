package logz

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/ibrt/golang-errors/errorz"
	"github.com/stretchr/testify/require"
)

type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}

type testErrorWithCause struct {
	message string
	err     error
}

func (e *testErrorWithCause) Error() string {
	return e.message
}

func (e *testErrorWithCause) Cause() error {
	return e.err
}

type testErrorWithUnwrap struct {
	message string
	err     error
}

func (e *testErrorWithUnwrap) Error() string {
	return e.message
}

func (e *testErrorWithUnwrap) Unwrap() error {
	return e.err
}

func TestErrorToSentryEvent(t *testing.T) {
	event := errorToSentryEvent(fmt.Errorf("test error"), Error)
	require.Equal(t, sentry.LevelError, event.Level)
	require.Empty(t, event.Extra)
	require.Len(t, event.Exception, 1)
	require.Equal(t, "*errors.errorString", event.Exception[0].Type)
	require.Equal(t, "test error", event.Exception[0].Value)
	require.Equal(t, "TestErrorToSentryEvent", event.Exception[0].Stacktrace.Frames[len(event.Exception[0].Stacktrace.Frames)-1].Function)

	event = errorToSentryEvent(errorz.Errorf("test error"), Warning)
	require.Equal(t, sentry.LevelWarning, event.Level)
	require.Empty(t, event.Extra)
	require.Len(t, event.Exception, 1)
	require.Equal(t, "*errors.errorString", event.Exception[0].Type)
	require.Equal(t, "test error", event.Exception[0].Value)
	require.Equal(t, "TestErrorToSentryEvent", event.Exception[0].Stacktrace.Frames[len(event.Exception[0].Stacktrace.Frames)-1].Function)

	event = errorToSentryEvent(errorz.Errorf("test error", errorz.ID("test-id"), errorz.Status(http.StatusBadRequest), errorz.M("k", "v")), Warning)
	require.Equal(t, sentry.LevelWarning, event.Level)
	require.Equal(t, map[string]interface{}{
		"status": http.StatusBadRequest,
		"k":      "v",
	}, event.Extra)
	require.Len(t, event.Exception, 1)
	require.Equal(t, "test-id", event.Exception[0].Type)
	require.Equal(t, "test error", event.Exception[0].Value)
	require.Equal(t, "TestErrorToSentryEvent", event.Exception[0].Stacktrace.Frames[len(event.Exception[0].Stacktrace.Frames)-1].Function)

	tErr := &testError{message: "test error"}
	event = errorToSentryEvent(errorz.Wrap(tErr), Warning)
	require.Equal(t, sentry.LevelWarning, event.Level)
	require.Equal(t, map[string]interface{}{"unwrapped[0][*logz.testError]": tErr}, event.Extra)
	require.Len(t, event.Exception, 1)
	require.Equal(t, "*logz.testError", event.Exception[0].Type)
	require.Equal(t, "test error", event.Exception[0].Value)
	require.Equal(t, "TestErrorToSentryEvent", event.Exception[0].Stacktrace.Frames[len(event.Exception[0].Stacktrace.Frames)-1].Function)

	tErr = &testError{message: "inner test error"}
	tErrC := &testErrorWithCause{message: "outer test error", err: tErr}
	event = errorToSentryEvent(errorz.Wrap(tErrC), Warning)
	require.Equal(t, sentry.LevelWarning, event.Level)
	require.Equal(t, map[string]interface{}{
		"unwrapped[0][*logz.testErrorWithCause]": tErrC,
		"unwrapped[1][*logz.testError]":          tErr,
	}, event.Extra)
	require.Len(t, event.Exception, 2)
	require.Equal(t, "*logz.testErrorWithCause", event.Exception[0].Type)
	require.Equal(t, "outer test error", event.Exception[0].Value)
	require.Equal(t, "TestErrorToSentryEvent", event.Exception[0].Stacktrace.Frames[len(event.Exception[0].Stacktrace.Frames)-1].Function)
	require.Equal(t, "*logz.testError", event.Exception[1].Type)
	require.Equal(t, "inner test error", event.Exception[1].Value)
	require.Nil(t, event.Exception[1].Stacktrace)

	tErr = &testError{message: "inner test error"}
	tErrU := &testErrorWithUnwrap{message: "outer test error", err: tErr}
	event = errorToSentryEvent(errorz.Wrap(tErrU), Warning)
	require.Equal(t, sentry.LevelWarning, event.Level)
	require.Equal(t, map[string]interface{}{
		"unwrapped[0][*logz.testErrorWithUnwrap]": tErrU,
		"unwrapped[1][*logz.testError]":           tErr,
	}, event.Extra)
	require.Len(t, event.Exception, 2)
	require.Equal(t, "*logz.testErrorWithUnwrap", event.Exception[0].Type)
	require.Equal(t, "outer test error", event.Exception[0].Value)
	require.Equal(t, "TestErrorToSentryEvent", event.Exception[0].Stacktrace.Frames[len(event.Exception[0].Stacktrace.Frames)-1].Function)
	require.Equal(t, "*logz.testError", event.Exception[1].Type)
	require.Equal(t, "inner test error", event.Exception[1].Value)
	require.Nil(t, event.Exception[1].Stacktrace)
}
