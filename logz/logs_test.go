package logz_test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/ibrt/golang-errors/errorz"
	"github.com/ibrt/golang-fixtures/fixturez"
	"github.com/stretchr/testify/require"

	"github.com/ibrt/golang-inject-logs/logz"
)

func TestModule_Debug(t *testing.T) {
	c := fixturez.CaptureOutput()
	defer c.Close()

	var recEvent *sentry.Event

	ctx := getCfg(func(event *sentry.Event, _ *sentry.EventHint) *sentry.Event {
		recEvent = event
		return event
	})

	injector, releaser := logz.Initializer(ctx)
	defer releaser()

	logs := logz.Get(injector(ctx))
	require.NotNil(t, logs)

	logs.Debug("message: %v", logz.A("value"), logz.M("k", "v"))
	v := make(map[string]interface{})
	fixturez.RequireNoError(t, json.Unmarshal(c.GetErr(), &v))
	require.Equal(t, "v", v["k"])
	require.Equal(t, "debug", v["level"])
	require.Equal(t, "message: value", v["msg"])
	require.NotNil(t, v["time"])

	require.NotNil(t, recEvent)
	require.Equal(t, sentry.LevelDebug, recEvent.Level)
	require.Equal(t, "message: value", recEvent.Message)
	require.Equal(t, map[string]interface{}{"k": "v"}, recEvent.Extra)
	require.Len(t, recEvent.Threads, 1)
	require.Equal(t, "TestModule_Debug", recEvent.Threads[0].Stacktrace.Frames[len(recEvent.Threads[0].Stacktrace.Frames)-1].Function)
}

func TestModule_Info(t *testing.T) {
	c := fixturez.CaptureOutput()
	defer c.Close()

	var recEvent *sentry.Event

	ctx := getCfg(func(event *sentry.Event, _ *sentry.EventHint) *sentry.Event {
		recEvent = event
		return event
	})

	injector, releaser := logz.Initializer(ctx)
	defer releaser()

	logs := logz.MaybeGet(injector(ctx))
	require.NotNil(t, logs)

	logs.Info("message: %v", logz.A("value"), logz.M("k", "v"))
	v := make(map[string]interface{})
	fixturez.RequireNoError(t, json.Unmarshal(c.GetErr(), &v))
	require.Equal(t, "v", v["k"])
	require.Equal(t, "info", v["level"])
	require.Equal(t, "message: value", v["msg"])
	require.NotNil(t, v["time"])

	require.NotNil(t, recEvent)
	require.Equal(t, sentry.LevelInfo, recEvent.Level)
	require.Equal(t, "message: value", recEvent.Message)
	require.Equal(t, map[string]interface{}{"k": "v"}, recEvent.Extra)
	require.Len(t, recEvent.Threads, 1)
	require.Equal(t, "TestModule_Info", recEvent.Threads[0].Stacktrace.Frames[len(recEvent.Threads[0].Stacktrace.Frames)-1].Function)
}

func TestModule_Warning(t *testing.T) {
	c := fixturez.CaptureOutput()
	defer c.Close()

	var recEvent *sentry.Event

	ctx := getCfg(func(event *sentry.Event, _ *sentry.EventHint) *sentry.Event {
		recEvent = event
		return event
	})

	injector, releaser := logz.Initializer(ctx)
	defer releaser()

	logs := logz.Get(injector(ctx))
	require.NotNil(t, logs)

	logs.Warning(errorz.Errorf("message: %v", errorz.A("value"), errorz.M("k", "v")))
	v := make(map[string]interface{})
	fixturez.RequireNoError(t, json.Unmarshal(c.GetErr(), &v))
	require.Equal(t, "v", v["k"])
	require.Equal(t, "warning", v["level"])
	require.Equal(t, "message: value", v["msg"])
	require.NotNil(t, v["time"])

	require.NotNil(t, recEvent)
	require.Equal(t, sentry.LevelWarning, recEvent.Level)
	require.Equal(t, map[string]interface{}{"k": "v"}, recEvent.Extra)
	require.Len(t, recEvent.Exception, 1)
	require.Equal(t, "*errors.errorString", recEvent.Exception[0].Type)
	require.Equal(t, "message: value", recEvent.Exception[0].Value)
	require.Equal(t, "TestModule_Warning", recEvent.Exception[0].Stacktrace.Frames[len(recEvent.Exception[0].Stacktrace.Frames)-1].Function)
}

func TestModule_Error(t *testing.T) {
	c := fixturez.CaptureOutput()
	defer c.Close()

	var recEvent *sentry.Event

	ctx := getCfg(func(event *sentry.Event, _ *sentry.EventHint) *sentry.Event {
		recEvent = event
		return event
	})

	injector, releaser := logz.Initializer(ctx)
	defer releaser()

	logs := logz.Get(injector(ctx))
	require.NotNil(t, logs)

	logs.Error(errorz.Errorf("message: %v", errorz.A("value"), errorz.M("k", "v")))
	v := make(map[string]interface{})
	fixturez.RequireNoError(t, json.Unmarshal(c.GetErr(), &v))
	require.Equal(t, "v", v["k"])
	require.Equal(t, "error", v["level"])
	require.Equal(t, "message: value", v["msg"])
	require.NotNil(t, v["time"])

	require.NotNil(t, recEvent)
	require.Equal(t, sentry.LevelError, recEvent.Level)
	require.Equal(t, map[string]interface{}{"k": "v"}, recEvent.Extra)
	require.Len(t, recEvent.Exception, 1)
	require.Equal(t, "*errors.errorString", recEvent.Exception[0].Type)
	require.Equal(t, "message: value", recEvent.Exception[0].Value)
	require.Equal(t, "TestModule_Error", recEvent.Exception[0].Stacktrace.Frames[len(recEvent.Exception[0].Stacktrace.Frames)-1].Function)
}

func TestModule_Tracing(t *testing.T) {
	c := fixturez.CaptureOutput()
	defer c.Close()

	var recEvent *sentry.Event

	ctx := getCfg(func(event *sentry.Event, _ *sentry.EventHint) *sentry.Event {
		recEvent = event
		return event
	})

	injector, releaser := logz.Initializer(ctx)
	defer releaser()
	ctx = injector(ctx)

	testReq := httptest.NewRequest("GET", "/path", nil)
	ctx, releaseTransaction := logz.Get(ctx).TraceHTTPRequestServer(testReq, []byte(`body`))

	logz.MaybeGet(ctx).SetUser(&logz.User{
		ID: "some-user-id",
		Metadata: logz.Metadata{
			"ku": "vu",
		},
	})

	logz.Get(ctx).AddMetadata("k2", "v2")

	ctx, releaseSpan := logz.Get(ctx).TraceSpan("test", "Test Span.")

	logz.Get(ctx).AddMetadata("sk1", "sv1")
	logz.Get(ctx).Error(errorz.Errorf("message: %v", errorz.A("value"), errorz.M("k1", "v1")))

	releaseSpan()
	releaseTransaction()

	v := make(map[string]interface{})
	fixturez.RequireNoError(t, json.Unmarshal(c.GetErr(), &v))
	require.Equal(t, "some-user-id", v["user-id"])
	require.Equal(t, "v1", v["k1"])
	require.Equal(t, "v2", v["k2"])
	require.Equal(t, "vu", v["ku"])
	require.Equal(t, "error", v["level"])
	require.Equal(t, "message: value", v["msg"])
	require.NotNil(t, v["time"])

	require.NotNil(t, recEvent)
	require.Equal(t, sentry.LevelError, recEvent.Level)
	require.Equal(t, map[string]interface{}{"k1": "v1", "k2": "v2", "ku": "vu"}, recEvent.Extra)
	require.Equal(t, "some-user-id", recEvent.User.ID)
	require.Len(t, recEvent.Exception, 1)
	require.Equal(t, "*errors.errorString", recEvent.Exception[0].Type)
	require.Equal(t, "message: value", recEvent.Exception[0].Value)
	require.Equal(t, "TestModule_Tracing", recEvent.Exception[0].Stacktrace.Frames[len(recEvent.Exception[0].Stacktrace.Frames)-1].Function)
}

func TestModule_Text(t *testing.T) {
	c := fixturez.CaptureOutput()
	defer c.Close()

	cfg := &logz.Config{
		SentryLevel:            logz.Debug,
		OutputLevel:            logz.Debug,
		OutputFormat:           logz.Text,
		SentryDSN:              "",
		SentrySampleRate:       1,
		SentryTracesSampleRate: 1,
		ReleaseTimeout:         1,
		Environment:            "environment",
		Release:                "release",
		ServerName:             "serverName",
	}

	ctx := logz.NewConfigSingletonInjector(cfg)(context.Background())
	injector, releaser := logz.Initializer(ctx)
	defer releaser()

	logs := logz.Get(injector(ctx))
	require.NotNil(t, logs)

	logs.Debug("message: %v", logz.A("value"), logz.M("k", "v"))
	out := string(c.GetErr())
	require.Contains(t, out, "DEBU")
	require.Contains(t, out, "message: value")
}

func TestGetters(t *testing.T) {
	require.Panics(t, func() {
		logz.Get(context.Background())
	})
	require.Panics(t, func() {
		logz.Initializer(context.Background())
	})
	fixturez.RequirePanicsWith(t, "Key: 'Config.SentryLevel' Error:Field validation for 'SentryLevel' failed on the 'required' tag\nKey: 'Config.OutputLevel' Error:Field validation for 'OutputLevel' failed on the 'required' tag\nKey: 'Config.OutputFormat' Error:Field validation for 'OutputFormat' failed on the 'required' tag\nKey: 'Config.SentrySampleRate' Error:Field validation for 'SentrySampleRate' failed on the 'required' tag\nKey: 'Config.SentryTracesSampleRate' Error:Field validation for 'SentryTracesSampleRate' failed on the 'required' tag", func() {
		ctx := logz.NewConfigSingletonInjector(&logz.Config{})(context.Background())
		logz.Initializer(ctx)
	})
	fixturez.RequireNotPanics(t, func() {
		require.Nil(t, logz.MaybeGet(context.Background()))
	})
}

func getCfg(beforeSendFunc logz.BeforeSendFunc) context.Context {
	cfg := &logz.Config{
		SentryLevel:            logz.Debug,
		OutputLevel:            logz.Debug,
		OutputFormat:           logz.JSON,
		SentryDSN:              "",
		SentrySampleRate:       1,
		SentryTracesSampleRate: 1,
		ReleaseTimeout:         1,
		Environment:            "environment",
		Release:                "release",
		ServerName:             "serverName",
		BeforeSend:             beforeSendFunc,
	}

	return logz.NewConfigSingletonInjector(cfg)(context.Background())
}
