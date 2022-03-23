package logz_test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/ibrt/golang-errors/errorz"
	"github.com/ibrt/golang-fixtures/fixturez"
	"github.com/ibrt/golang-inject-clock/clockz"
	"github.com/ibrt/golang-inject-clock/clockz/testclockz"
	"github.com/stretchr/testify/require"

	"github.com/ibrt/golang-inject-logs/logz"
)

var (
	_ sentry.Transport = &testTransport{}
)

type testTransport struct {
	clientOptions *sentry.ClientOptions
	events        []*sentry.Event
	isFlushed     bool
}

// Flush implements the sentry.Transport interface.
func (t *testTransport) Flush(_ time.Duration) bool {
	t.events = nil
	t.isFlushed = true
	return true
}

// Configure implements the sentry.Transport interface.
func (t *testTransport) Configure(clientOptions sentry.ClientOptions) {
	t.clientOptions = &clientOptions
}

// SendEvent implements the sentry.Transport interface.
func (t *testTransport) SendEvent(event *sentry.Event) {
	t.events = append(t.events, event)
	t.isFlushed = false
}

func setupLogs(ctx context.Context) (context.Context, func(), *testTransport) {
	transport := &testTransport{}

	cfg := &logz.Config{
		SentryLevel:            logz.Debug,
		OutputLevel:            logz.Debug,
		OutputFormat:           logz.JSON,
		SentryDSN:              "",
		SentrySampleRate:       1,
		SentryTracesSampleRate: 1,
		ReleaseTimeoutSeconds:  5,
		Environment:            "environment",
		Release:                "release",
		ServerName:             "serverName",
		SentryTransport:        transport,
		BeforeSend: func(event *sentry.Event, _ *sentry.EventHint) *sentry.Event {
			return event
		},
	}

	ctx = logz.NewConfigSingletonInjector(cfg)(ctx)
	injector, releaser := logz.Initializer(ctx)
	ctx = injector(ctx)

	return ctx, releaser, transport
}

type ModuleSuite struct {
	*fixturez.DefaultConfigMixin
	Clock *testclockz.MockHelper
}

func TestModule(t *testing.T) {
	fixturez.RunSuite(t, &ModuleSuite{})
}

func (s *ModuleSuite) TestDebug(ctx context.Context, t *testing.T) {
	c := fixturez.CaptureOutput()
	defer c.Close()
	ctx, releaser, transport := setupLogs(ctx)
	defer releaser()

	logz.Get(ctx).SetUser(&logz.User{ID: "user-id"})
	logz.Get(ctx).AddMetadata("gk", "gv")

	logz.Get(ctx).Debug("message: %v", logz.A("value"), logz.M("k", "v"))
	buf, err := json.Marshal(map[string]interface{}{
		"level": "debug",
		"time":  clockz.Get(ctx).Now(),
		"msg":   "message: value",
		"k":     "v",
		"gk":    "gv",
		"uid":   "user-id",
	})
	fixturez.RequireNoError(t, err)
	require.JSONEq(t, string(buf), string(c.GetErr()))

	require.Len(t, transport.events, 1)
	event := transport.events[0]

	require.NotEmpty(t, event.Contexts)
	event.Contexts = nil

	require.NotEmpty(t, event.EventID)
	event.EventID = ""

	require.NotEmpty(t, event.Sdk)
	event.Sdk = sentry.SdkInfo{}

	require.Len(t, event.Threads, 1)
	require.NotNil(t, event.Threads[0].Stacktrace)
	require.Equal(t, "(*ModuleSuite).TestDebug", event.Threads[0].Stacktrace.Frames[len(event.Threads[0].Stacktrace.Frames)-1].Function)
	event.Threads = nil

	require.Equal(t, &sentry.Event{
		Tags:        make(map[string]string),
		Environment: "environment",
		Release:     "release",
		ServerName:  "serverName",
		Platform:    "go",
		Extra:       map[string]interface{}{"k": "v", "gk": "gv"},
		Level:       sentry.LevelDebug,
		Message:     "message: value",
		Timestamp:   clockz.Get(ctx).Now(),
		User: sentry.User{
			ID: "user-id",
		},
	}, event)
}

func (s *ModuleSuite) TestInfo(ctx context.Context, t *testing.T) {
	c := fixturez.CaptureOutput()
	defer c.Close()
	ctx, releaser, transport := setupLogs(ctx)
	defer releaser()

	logz.Get(ctx).SetUser(&logz.User{ID: "user-id"})
	logz.Get(ctx).AddMetadata("gk", "gv")

	logz.Get(ctx).Info("message: %v", logz.A("value"), logz.M("k", "v"))
	buf, err := json.Marshal(map[string]interface{}{
		"level": "info",
		"time":  clockz.Get(ctx).Now(),
		"msg":   "message: value",
		"k":     "v",
		"gk":    "gv",
		"uid":   "user-id",
	})
	fixturez.RequireNoError(t, err)
	require.JSONEq(t, string(buf), string(c.GetErr()))

	require.Len(t, transport.events, 1)
	event := transport.events[0]

	require.NotEmpty(t, event.Contexts)
	event.Contexts = nil

	require.NotEmpty(t, event.EventID)
	event.EventID = ""

	require.NotEmpty(t, event.Sdk)
	event.Sdk = sentry.SdkInfo{}

	require.Len(t, event.Threads, 1)
	require.NotNil(t, event.Threads[0].Stacktrace)
	require.Equal(t, "(*ModuleSuite).TestInfo", event.Threads[0].Stacktrace.Frames[len(event.Threads[0].Stacktrace.Frames)-1].Function)
	event.Threads = nil

	require.Equal(t, &sentry.Event{
		Tags:        make(map[string]string),
		Environment: "environment",
		Release:     "release",
		ServerName:  "serverName",
		Platform:    "go",
		Extra:       map[string]interface{}{"k": "v", "gk": "gv"},
		Level:       sentry.LevelInfo,
		Message:     "message: value",
		Timestamp:   clockz.Get(ctx).Now(),
		User: sentry.User{
			ID: "user-id",
		},
	}, event)
}

func (s *ModuleSuite) TestWarning(ctx context.Context, t *testing.T) {
	c := fixturez.CaptureOutput()
	defer c.Close()
	ctx, releaser, transport := setupLogs(ctx)
	defer releaser()

	logz.Get(ctx).SetUser(&logz.User{ID: "user-id"})
	logz.Get(ctx).AddMetadata("gk", "gv")

	logz.Get(ctx).Warning(errorz.Errorf("message: %v", errorz.A("value"), errorz.M("k", "v")))
	buf, err := json.Marshal(map[string]interface{}{
		"level": "warning",
		"time":  clockz.Get(ctx).Now(),
		"msg":   "message: value",
		"k":     "v",
		"gk":    "gv",
		"uid":   "user-id",
	})
	fixturez.RequireNoError(t, err)
	require.JSONEq(t, string(buf), string(c.GetErr()))

	require.Len(t, transport.events, 1)
	event := transport.events[0]

	require.NotEmpty(t, event.Contexts)
	event.Contexts = nil

	require.NotEmpty(t, event.EventID)
	event.EventID = ""

	require.NotEmpty(t, event.Sdk)
	event.Sdk = sentry.SdkInfo{}

	require.Len(t, event.Exception, 1)
	require.NotNil(t, event.Exception[0].Stacktrace)
	require.Equal(t, "(*ModuleSuite).TestWarning", event.Exception[0].Stacktrace.Frames[len(event.Exception[0].Stacktrace.Frames)-1].Function)
	event.Exception[0].Stacktrace = nil

	require.Equal(t, &sentry.Event{
		Tags:        make(map[string]string),
		Environment: "environment",
		Release:     "release",
		ServerName:  "serverName",
		Platform:    "go",
		Extra:       map[string]interface{}{"k": "v", "gk": "gv"},
		Level:       sentry.LevelWarning,
		Timestamp:   clockz.Get(ctx).Now(),
		User: sentry.User{
			ID: "user-id",
		},
		Exception: []sentry.Exception{
			{
				Type:  "*errors.errorString",
				Value: "message: value",
			},
		},
	}, event)
}

func (s *ModuleSuite) TestError(ctx context.Context, t *testing.T) {
	c := fixturez.CaptureOutput()
	defer c.Close()
	ctx, releaser, transport := setupLogs(ctx)
	defer releaser()

	logz.Get(ctx).SetUser(&logz.User{ID: "user-id"})
	logz.Get(ctx).AddMetadata("gk", "gv")

	logz.Get(ctx).Error(errorz.Errorf("message: %v", errorz.A("value"), errorz.M("k", "v")))
	buf, err := json.Marshal(map[string]interface{}{
		"level": "error",
		"time":  clockz.Get(ctx).Now(),
		"msg":   "message: value",
		"k":     "v",
		"gk":    "gv",
		"uid":   "user-id",
	})
	fixturez.RequireNoError(t, err)
	require.JSONEq(t, string(buf), string(c.GetErr()))

	require.Len(t, transport.events, 1)
	event := transport.events[0]

	require.NotEmpty(t, event.Contexts)
	event.Contexts = nil

	require.NotEmpty(t, event.EventID)
	event.EventID = ""

	require.NotEmpty(t, event.Sdk)
	event.Sdk = sentry.SdkInfo{}

	require.Len(t, event.Exception, 1)
	require.NotNil(t, event.Exception[0].Stacktrace)
	require.Equal(t, "(*ModuleSuite).TestError", event.Exception[0].Stacktrace.Frames[len(event.Exception[0].Stacktrace.Frames)-1].Function)
	event.Exception[0].Stacktrace = nil

	require.Equal(t, &sentry.Event{
		Tags:        make(map[string]string),
		Environment: "environment",
		Release:     "release",
		ServerName:  "serverName",
		Platform:    "go",
		Extra:       map[string]interface{}{"k": "v", "gk": "gv"},
		Level:       sentry.LevelError,
		Timestamp:   clockz.Get(ctx).Now(),
		User: sentry.User{
			ID: "user-id",
		},
		Exception: []sentry.Exception{
			{
				Type:  "*errors.errorString",
				Value: "message: value",
			},
		},
	}, event)
}

func (s *ModuleSuite) TestTracing(ctx context.Context, t *testing.T) {
	//c := fixturez.CaptureOutput()
	//defer c.Close()
	ctx, releaser, transport := setupLogs(ctx)
	defer releaser()
	startTime := clockz.Get(ctx).Now()

	func() {
		testReq := httptest.NewRequest("GET", "/path", nil)
		testReq.RemoteAddr = "1.2.3.4:4321"
		ctx, releaseTransaction := logz.Get(ctx).TraceHTTPRequestServer(testReq, []byte(`body`))
		defer releaseTransaction()

		logz.Get(ctx).SetUser(&logz.User{
			ID: "user-id",
			Metadata: logz.Metadata{
				"ku": "vu",
			},
		})

		logz.Get(ctx).AddMetadata("kt", "vt")
		s.Clock.Mock.Add(time.Second)

		ctx, releaseSpan := logz.Get(ctx).TraceSpan("test", "Test Span.")
		defer releaseSpan()

		logz.Get(ctx).AddMetadata("ks", "vs")
		s.Clock.Mock.Add(time.Second)
	}()

	require.Len(t, transport.events, 1)
	event := transport.events[0]

	require.NotEmpty(t, event.Contexts)
	event.Contexts = nil

	require.NotEmpty(t, event.EventID)
	event.EventID = ""

	require.NotEmpty(t, event.Sdk)
	event.Sdk = sentry.SdkInfo{}

	require.Len(t, event.Spans, 1)
	span := event.Spans[0]
	require.NotEmpty(t, span.TraceID)
	require.NotEmpty(t, span.SpanID)
	require.NotEmpty(t, span.ParentSpanID)
	require.Equal(t, startTime.Add(time.Second), span.StartTime)
	require.Equal(t, startTime.Add(2*time.Second), span.EndTime)
	require.Equal(t, "test", span.Op)
	require.Equal(t, "Test Span.", span.Description)
	require.Equal(t, map[string]interface{}{"ks": "vs"}, span.Data)
	event.Spans = nil

	require.Equal(t, &sentry.Event{
		Environment: "environment",
		Release:     "release",
		ServerName:  "serverName",
		Type:        "transaction",
		StartTime:   startTime,
		Timestamp:   startTime.Add(2 * time.Second),
		Transaction: "GET /path",
		Platform:    "go",
		Extra:       map[string]interface{}{"ku": "vu", "kt": "vt"},
		Level:       sentry.LevelInfo,
		Request: &sentry.Request{
			URL:         "http://example.com/path",
			Method:      "GET",
			Data:        "body",
			QueryString: "",
			Cookies:     "",
			Headers: map[string]string{
				"Host": "example.com",
			},
			Env: map[string]string{
				"REMOTE_ADDR": "1.2.3.4",
				"REMOTE_PORT": "4321",
			},
		},
		User: sentry.User{
			ID: "user-id",
		},
	}, event)
}

func (s *ModuleSuite) TestNoopLogs(_ context.Context, t *testing.T) {
	c := fixturez.CaptureOutput()
	defer c.Close()

	fixturez.RequireNotPanics(t, func() {
		noopLogs := logz.Get(context.Background())
		noopLogs.Debug("message")
		noopLogs.Info("message")
		noopLogs.Warning(errorz.Errorf("error"))
		noopLogs.Error(errorz.Errorf("error"))
		_, releaser := noopLogs.TraceHTTPRequestServer(nil, nil)
		releaser()
		_, releaser = noopLogs.TraceSpan("op", "desc")
		releaser()
		noopLogs.SetUser(nil)
		noopLogs.AddMetadata("k", "v")
	})

	require.Empty(t, c.GetErr())
	require.Empty(t, c.GetOut())
}

func (s *ModuleSuite) TestTextOutput(ctx context.Context, t *testing.T) {
	c := fixturez.CaptureOutput()
	defer c.Close()

	cfg := &logz.Config{
		SentryLevel:            logz.Debug,
		OutputLevel:            logz.Debug,
		OutputFormat:           logz.Text,
		SentryDSN:              "",
		SentrySampleRate:       1,
		SentryTracesSampleRate: 1,
		ReleaseTimeoutSeconds:  5,
		Environment:            "environment",
		Release:                "release",
		ServerName:             "serverName",
	}

	ctx = logz.NewConfigSingletonInjector(cfg)(ctx)
	injector, releaser := logz.Initializer(ctx)
	defer releaser()
	ctx = injector(ctx)

	logz.Get(ctx).Debug("message: %v", logz.A("value"), logz.M("k", "v"))
	out := string(c.GetErr())
	require.Contains(t, out, "DEBU")
	require.Contains(t, out, "message: value")
}

/*


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

	logz.Get(ctx).SetUser(&logz.User{
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
		ReleaseTimeoutSeconds:  5,
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


}*/
