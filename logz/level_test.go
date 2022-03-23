package logz

import (
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/ibrt/golang-fixtures/fixturez"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestLevel(t *testing.T) {
	require.Equal(t, sentry.LevelDebug, Debug.toSentry())
	require.Equal(t, sentry.LevelInfo, Info.toSentry())
	require.Equal(t, sentry.LevelWarning, Warning.toSentry())
	require.Equal(t, sentry.LevelError, Error.toSentry())
	require.Equal(t, logrus.DebugLevel, Debug.toLogrus())
	require.Equal(t, logrus.InfoLevel, Info.toLogrus())
	require.Equal(t, logrus.WarnLevel, Warning.toLogrus())
	require.Equal(t, logrus.ErrorLevel, Error.toLogrus())
	require.Equal(t, Debug, levelFromSentry(sentry.LevelDebug))
	require.Equal(t, Info, levelFromSentry(sentry.LevelInfo))
	require.Equal(t, Warning, levelFromSentry(sentry.LevelWarning))
	require.Equal(t, Error, levelFromSentry(sentry.LevelError))
	require.Equal(t, Error, levelFromSentry(sentry.LevelFatal))

	fixturez.RequirePanicsWith(t, "unknown level: unknown", func() {
		Level("unknown").toSentry()
	})

	fixturez.RequirePanicsWith(t, "unknown level: unknown", func() {
		Level("unknown").toLogrus()
	})

	fixturez.RequirePanicsWith(t, "unknown level: unknown", func() {
		levelFromSentry("unknown")
	})
}
