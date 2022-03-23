package logz

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/ibrt/golang-errors/errorz"
	"github.com/ibrt/golang-inject-clock/clockz"
	"github.com/ibrt/golang-inject/injectz"
	"github.com/ibrt/golang-validation/vz"
	"github.com/sirupsen/logrus"
)

type contextKey int

const (
	logsConfigContextKey contextKey = iota
	logsContextKey
	logsSpanContextKey
)

var (
	_ Logs        = &logsImpl{}
	_ Logs        = &noopLogsImpl{}
	_ ContextLogs = &contextLogsImpl{}

	noopLogs = &noopLogsImpl{}
)

// OutputFormat describes the format for output logs.
type OutputFormat string

// Known formats.
const (
	Text OutputFormat = "text"
	JSON OutputFormat = "json"
)

// BeforeSendFunc describes a function called before sending out an event.
type BeforeSendFunc func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event

// Config describes the configuration for Logs.
type Config struct {
	SentryLevel            Level            `json:"sentryLevel" validate:"required,oneof=debug info warning error"`
	OutputLevel            Level            `json:"outputLevel" validate:"required,oneof=debug info warning error"`
	OutputFormat           OutputFormat     `json:"format" validate:"required,oneof=text json"`
	SentryDSN              string           `json:"sentryDsn"`
	SentrySampleRate       float64          `json:"sentrySampleRate" validate:"required"`
	SentryTracesSampleRate float64          `json:"sentryTracesSampleRate" validate:"required"`
	SentryTransport        sentry.Transport `json:"-"`
	ReleaseTimeoutSeconds  int              `json:"releaseTimeoutSeconds"`
	Environment            string           `json:"environment"`
	Release                string           `json:"release"`
	ServerName             string           `json:"serverName"`
	BeforeSend             BeforeSendFunc   `json:"-"`
}

// Validate implements the vz.Validator interface.
func (c *Config) Validate() error {
	return errorz.MaybeWrap(vz.ValidateStruct(c), errorz.Skip())
}

// NewConfigSingletonInjector always inject the given *Config.
func NewConfigSingletonInjector(cfg *Config) injectz.Injector {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logsConfigContextKey, cfg)
	}
}

// Logs describes the logs module.
type Logs interface {
	Debug(ctx context.Context, skipCallers int, format string, options ...Option)
	Info(ctx context.Context, skipCallers int, format string, options ...Option)
	Warning(ctx context.Context, err error)
	Error(ctx context.Context, err error)
	TraceHTTPRequestServer(ctx context.Context, req *http.Request, reqBody []byte) (context.Context, func())
	TraceSpan(ctx context.Context, op, desc string) (context.Context, func())
	SetUser(ctx context.Context, user *User)
	AddMetadata(ctx context.Context, k string, v interface{})
}

type logsImpl struct {
	logrusLogger *logrus.Logger
	sentryHub    *sentry.Hub
}

// Debug logs a debug message.
func (l *logsImpl) Debug(ctx context.Context, skipCallers int, format string, options ...Option) {
	sentry.GetHubFromContext(ctx).CaptureEvent(
		newEntry(ctx, Debug, skipCallers+1, format, options...).toSentryEvent())
}

// Info logs an info message.
func (l *logsImpl) Info(ctx context.Context, skipCallers int, format string, options ...Option) {
	sentry.GetHubFromContext(ctx).CaptureEvent(
		newEntry(ctx, Info, skipCallers+1, format, options...).toSentryEvent())
}

// Warning logs a warning.
func (l *logsImpl) Warning(ctx context.Context, err error) {
	sentry.GetHubFromContext(ctx).CaptureEvent(
		errorToSentryEvent(ctx, errorz.Wrap(err, errorz.SkipPackage()), Warning))
}

// Error logs an error.
func (l *logsImpl) Error(ctx context.Context, err error) {
	sentry.GetHubFromContext(ctx).CaptureEvent(
		errorToSentryEvent(ctx, errorz.Wrap(err, errorz.SkipPackage()), Error))
}

// TraceHTTPRequestServer starts tracing an inbound HTTP request.
func (l *logsImpl) TraceHTTPRequestServer(ctx context.Context, req *http.Request, reqBody []byte) (context.Context, func()) {
	sentryHub := sentry.GetHubFromContext(ctx).Clone()
	ctx = sentry.SetHubOnContext(ctx, sentryHub)

	span := sentry.StartSpan(ctx, "http.server",
		sentry.TransactionName(fmt.Sprintf("%s %s", req.Method, req.URL.Path)),
		sentry.ContinueFromRequest(req))
	span.StartTime = clockz.Get(ctx).Now()
	ctx = span.Context()

	sentryHub.Scope().SetRequest(req)
	if len(reqBody) > 0 {
		sentryHub.Scope().SetRequestBody(reqBody)
	}

	return ctx, func() {
		span.EndTime = clockz.Get(ctx).Now()
		span.Finish()
	}
}

// TraceSpan starts tracing a span, for example an outgoing HTTP request or database query.
func (l *logsImpl) TraceSpan(ctx context.Context, op, desc string) (context.Context, func()) {
	span := sentry.StartSpan(ctx, op)
	span.StartTime = clockz.Get(ctx).Now()
	span.Data = make(map[string]interface{})
	span.Description = desc

	ctx = span.Context()
	ctx = context.WithValue(ctx, logsSpanContextKey, span)

	return ctx, func() {
		span.EndTime = clockz.Get(ctx).Now()
		span.Finish()
	}
}

// SetUser sets the user in the current scope.
func (l *logsImpl) SetUser(ctx context.Context, user *User) {
	if user != nil {
		scope := sentry.GetHubFromContext(ctx).Scope()
		scope.SetUser(sentry.User{
			ID:    user.ID,
			Email: user.Email,
		})
		scope.SetExtras(user.Metadata)
	}
}

// AddMetadata adds the given metadata to the current scope.
func (l *logsImpl) AddMetadata(ctx context.Context, k string, v interface{}) {
	if span, ok := ctx.Value(logsSpanContextKey).(*sentry.Span); ok {
		span.Data[k] = v
		return
	}

	sentry.GetHubFromContext(ctx).Scope().SetExtra(k, v)
}

type noopLogsImpl struct {
}

// Debug logs a debug message.
func (l *noopLogsImpl) Debug(_ context.Context, _ int, _ string, _ ...Option) {
	// nothing to do here
}

// Info logs an info message.
func (l *noopLogsImpl) Info(_ context.Context, _ int, _ string, _ ...Option) {
	// nothing to do here
}

// Warning logs a warning.
func (l *noopLogsImpl) Warning(_ context.Context, _ error) {
	// nothing to do here
}

// Error logs an error.
func (l *noopLogsImpl) Error(_ context.Context, _ error) {
	// nothing to do here
}

// TraceHTTPRequestServer starts tracing an inbound HTTP request.
func (l *noopLogsImpl) TraceHTTPRequestServer(ctx context.Context, _ *http.Request, _ []byte) (context.Context, func()) {
	return ctx, func() {}
}

// TraceSpan starts tracing a span, for example an outgoing HTTP request or database query.
func (l *noopLogsImpl) TraceSpan(ctx context.Context, _, _ string) (context.Context, func()) {
	return ctx, func() {}
}

// SetUser sets the user in the current scope.
func (l *noopLogsImpl) SetUser(_ context.Context, _ *User) {
	// nothing to do here
}

// AddMetadata adds the given metadata to the current scope.
func (l *noopLogsImpl) AddMetadata(_ context.Context, _ string, _ interface{}) {
	// nothing to do here
}

// ContextLogs describes a Logs with a cached context.
type ContextLogs interface {
	Debug(format string, options ...Option)
	Info(format string, options ...Option)
	Warning(err error)
	Error(err error)
	TraceHTTPRequestServer(req *http.Request, reqBody []byte) (context.Context, func())
	TraceSpan(op, desc string) (context.Context, func())
	SetUser(user *User)
	AddMetadata(k string, v interface{})
}

type contextLogsImpl struct {
	ctx  context.Context
	logs Logs
}

// Debug logs a debug message.
func (l *contextLogsImpl) Debug(format string, options ...Option) {
	l.logs.Debug(l.ctx, 1, format, options...)
}

// Info logs an info message.
func (l *contextLogsImpl) Info(format string, options ...Option) {
	l.logs.Info(l.ctx, 1, format, options...)
}

// Warning logs a warning.
func (l *contextLogsImpl) Warning(err error) {
	l.logs.Warning(l.ctx, errorz.Wrap(err, errorz.SkipPackage()))
}

// Error logs an error.
func (l *contextLogsImpl) Error(err error) {
	l.logs.Error(l.ctx, errorz.Wrap(err, errorz.SkipPackage()))
}

// TraceHTTPRequestServer starts tracing an inbound HTTP request.
func (l *contextLogsImpl) TraceHTTPRequestServer(req *http.Request, reqBody []byte) (context.Context, func()) {
	return l.logs.TraceHTTPRequestServer(l.ctx, req, reqBody)
}

// TraceSpan starts tracing a span, for example an outgoing HTTP request or database query.
func (l *contextLogsImpl) TraceSpan(op, desc string) (context.Context, func()) {
	return l.logs.TraceSpan(l.ctx, op, desc)
}

// SetUser sets the user in the current scope.
func (l *contextLogsImpl) SetUser(user *User) {
	l.logs.SetUser(l.ctx, user)
}

// AddMetadata adds the given metadata to the current scope.
func (l *contextLogsImpl) AddMetadata(k string, v interface{}) {
	l.logs.AddMetadata(l.ctx, k, v)
}

// Initializer is a Logs initializer which provides a default implementation using Logrus and Sentry.
func Initializer(ctx context.Context) (injectz.Injector, injectz.Releaser) {
	cfg := ctx.Value(logsConfigContextKey).(*Config)
	errorz.MaybeMustWrap(cfg.Validate(), errorz.SkipPackage())

	logrusLogger := logrus.New()
	logrusLogger.SetLevel(cfg.OutputLevel.toLogrus())

	if cfg.OutputFormat == JSON {
		logrusLogger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	} else {
		logrusLogger.SetFormatter(&logrus.TextFormatter{
			ForceColors: true,
		})
	}

	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:              cfg.SentryDSN,
		SampleRate:       cfg.SentrySampleRate,
		TracesSampleRate: cfg.SentryTracesSampleRate,
		ServerName:       cfg.ServerName,
		Release:          cfg.Release,
		Environment:      cfg.Environment,
		Transport:        cfg.SentryTransport,
		BeforeSend: func(event *sentry.Event, eventHint *sentry.EventHint) *sentry.Event {
			if cfg.BeforeSend != nil {
				event = cfg.BeforeSend(event, eventHint)
			}

			logrusEntry := logrusLogger.
				WithTime(event.Timestamp).
				WithFields(event.Extra)

			if event.User.ID != "" {
				logrusEntry = logrusEntry.WithField("uid", event.User.ID)
			}

			message := event.Message
			if len(event.Exception) > 0 {
				message = event.Exception[0].Value
			}

			logrusEntry.Log(levelFromSentry(event.Level).toLogrus(), message)
			return event
		},
	})
	errorz.MaybeMustWrap(err, errorz.SkipPackage())
	sentryHub := sentry.NewHub(client, sentry.NewScope())

	return injectz.NewInjectors(
			NewSingletonInjector(&logsImpl{
				logrusLogger: logrusLogger,
				sentryHub:    sentryHub,
			}),
			func(ctx context.Context) context.Context {
				return sentry.SetHubOnContext(ctx, sentryHub)
			}),
		func() {
			sentryHub.Flush(time.Duration(cfg.ReleaseTimeoutSeconds) * time.Second)
		}
}

// NewSingletonInjector always injects the given Logs.
func NewSingletonInjector(l Logs) injectz.Injector {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logsContextKey, l)
	}
}

// Get extracts the Logs from context and wraps it as ContextLogs, returns a no-op Logs if not found.
func Get(ctx context.Context) ContextLogs {
	logs, ok := ctx.Value(logsContextKey).(Logs)
	if !ok {
		logs = noopLogs
	}

	return &contextLogsImpl{
		ctx:  ctx,
		logs: logs,
	}
}
