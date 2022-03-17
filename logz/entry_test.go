package logz

import (
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/require"
)

func TestEntry(t *testing.T) {
	e := newEntry(Info, 0, "message: %v", A("value"), M("k1", "v1"), Metadata{"k2": "v2"})
	event := e.toSentryEvent()
	require.Equal(t, sentry.LevelInfo, event.Level)
	require.Equal(t, "message: value", event.Message)
	require.Equal(t, map[string]interface{}{"k1": "v1", "k2": "v2"}, event.Extra)
	require.Len(t, event.Threads, 1)
	require.Equal(t, "TestEntry", event.Threads[0].Stacktrace.Frames[len(event.Threads[0].Stacktrace.Frames)-1].Function)

	e = newEntry(Warning, 0, "message: %v", A("value"), M("k1", "v1"), Metadata{"k2": "v2"})
	event = e.toSentryEvent()
	require.Equal(t, sentry.LevelWarning, event.Level)
	require.Equal(t, "message: value", event.Message)
	require.Equal(t, map[string]interface{}{"k1": "v1", "k2": "v2"}, event.Extra)
	require.Len(t, event.Threads, 1)
	require.Equal(t, "TestEntry", event.Threads[0].Stacktrace.Frames[len(event.Threads[0].Stacktrace.Frames)-1].Function)
}
