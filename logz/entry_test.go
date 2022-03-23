package logz

import (
	"context"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/ibrt/golang-fixtures/fixturez"
	"github.com/ibrt/golang-inject-clock/clockz"
	"github.com/ibrt/golang-inject-clock/clockz/testclockz"
	"github.com/stretchr/testify/require"
)

type EntrySuite struct {
	*fixturez.DefaultConfigMixin
	Clock *testclockz.MockHelper
}

func TestEntry(t *testing.T) {
	fixturez.RunSuite(t, &EntrySuite{})
}

func (s *EntrySuite) TestEntry(ctx context.Context, t *testing.T) {
	e := newEntry(ctx, Debug, 0, "message: %v", A("value"), M("k1", "v1"), Metadata{"k2": "v2"})
	event := e.toSentryEvent()
	require.Len(t, event.Threads, 1)
	event.Threads = nil
	require.Equal(t, &sentry.Event{
		Contexts:  make(map[string]interface{}),
		Tags:      make(map[string]string),
		Modules:   make(map[string]string),
		Extra:     map[string]interface{}{"k1": "v1", "k2": "v2"},
		Level:     sentry.LevelDebug,
		Message:   "message: value",
		Timestamp: clockz.Get(ctx).Now(),
	}, event)

	e = newEntry(ctx, Warning, 0, "message: %v", A("value"), M("k1", "v1"), Metadata{"k2": "v2"})
	event = e.toSentryEvent()
	require.Len(t, event.Threads, 1)
	event.Threads = nil
	require.Equal(t, &sentry.Event{
		Contexts:  make(map[string]interface{}),
		Tags:      make(map[string]string),
		Modules:   make(map[string]string),
		Extra:     map[string]interface{}{"k1": "v1", "k2": "v2"},
		Level:     sentry.LevelWarning,
		Message:   "message: value",
		Timestamp: clockz.Get(ctx).Now(),
	}, event)
}
