package logz

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"github.com/ibrt/golang-errors/errorz"
	"github.com/ibrt/golang-inject/injectz"
	"github.com/sirupsen/logrus"
)

type contextKey int

const (
	logsConfigContextKey contextKey = iota
	logsContextKey
)

var (
	_ Logs        = &logsImpl{}
	_ ContextLogs = &contextLogsImpl{}

	validate = validator.New()
)

// Level describes a logs level.
type Level string

func (l Level) toLogrus() logrus.Level {
	switch l {
	case Debug:
		return logrus.DebugLevel
	case Info:
		return logrus.InfoLevel
	case Warning:
		return logrus.WarnLevel
	case Error:
		return logrus.ErrorLevel
	default:
		panic(errorz.Errorf("unknown level: %v", errorz.A(l)))
	}
}

func (l Level) toSentry() sentry.Level {
	switch l {
	case Debug:
		return sentry.LevelDebug
	case Info:
		return sentry.LevelInfo
	case Warning:
		return sentry.LevelWarning
	case Error:
		return sentry.LevelError
	default:
		panic(errorz.Errorf("unknown level: %v", errorz.A(l)))
	}
}

func levelFromSentry(l sentry.Level) Level {
	switch l {
	case sentry.LevelFatal, sentry.LevelError:
		return Error
	case sentry.LevelWarning:
		return Warning
	case sentry.LevelInfo:
		return Info
	case sentry.LevelDebug:
		return Debug
	default:
		panic(errorz.Errorf("unknown level: %v", errorz.A(l)))
	}
}

// Known levels.
const (
	Debug   Level = "debug"
	Info    Level = "info"
	Warning Level = "warning"
	Error   Level = "error"
)

// OutputFormat describes the format for output logs.
type OutputFormat string

// Known formats.
const (
	Text OutputFormat = "text"
	JSON OutputFormat = "json"
)

// LogsConfig describes the configuration for the logz module.
type LogsConfig struct {
	SentryLevel      Level         `json:"sentryLevel" validate:"required,oneof=debug info warning error"`
	OutputLevel      Level         `json:"outputLevel" validate:"required,oneof=debug info warning error"`
	OutputFormat     OutputFormat  `json:"format" validate:"required,oneof=text json"`
	SentryDSN        string        `json:"sentryDsn"`
	SentrySampleRate float64       `json:"sampleRate" validate:"required"`
	ReleaseTimeout   time.Duration `json:"releaseTimeout"`
	Environment      string        `json:"environment"`
	Release          string        `json:"release"`
	ServerName       string        `json:"serverName"`
}

// NewConfigSingletonInjector always inject the given LogsConfig.
func NewConfigSingletonInjector(cfg *LogsConfig) injectz.Injector {
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
		newEntry(Debug, skipCallers+1, format, options...).toSentryEvent())
}

// Info logs an info message.
func (l *logsImpl) Info(ctx context.Context, skipCallers int, format string, options ...Option) {
	sentry.GetHubFromContext(ctx).CaptureEvent(
		newEntry(Info, skipCallers+1, format, options...).toSentryEvent())
}

// Warning logs a warning.
func (l *logsImpl) Warning(ctx context.Context, err error) {
	sentry.GetHubFromContext(ctx).CaptureEvent(
		errorToSentryEvent(errorz.Wrap(err, errorz.Skip()), Warning))
}

// Error logs an error.
func (l *logsImpl) Error(ctx context.Context, err error) {
	sentry.GetHubFromContext(ctx).CaptureEvent(
		errorToSentryEvent(errorz.Wrap(err, errorz.Skip()), Error))
}

// TraceHTTPRequestServer starts tracing an inbound HTTP request.
func (l *logsImpl) TraceHTTPRequestServer(ctx context.Context, req *http.Request, reqBody []byte) (context.Context, func()) {
	sentryHub := sentry.GetHubFromContext(ctx).Clone()
	ctx = sentry.SetHubOnContext(ctx, sentryHub)

	span := sentry.StartSpan(ctx, "http.server",
		sentry.TransactionName(fmt.Sprintf("%s %s", req.Method, req.URL.Path)),
		sentry.ContinueFromRequest(req))

	sentryHub.Scope().SetRequest(req)
	if len(reqBody) > 0 {
		sentryHub.Scope().SetRequestBody(reqBody)
	}

	return ctx, func() {
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
	sentry.GetHubFromContext(ctx).Scope().SetExtra(k, v)
}

// ContextLogs describes a Logs with a cached context.
type ContextLogs interface {
	Debug(format string, options ...Option)
	Info(format string, options ...Option)
	Warning(err error)
	Error(err error)
	TraceHTTPRequestServer(req *http.Request, reqBody []byte) (context.Context, func())
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
	l.logs.Warning(l.ctx, errorz.Wrap(err, errorz.Skip()))
}

// Error logs an error.
func (l *contextLogsImpl) Error(err error) {
	l.logs.Error(l.ctx, errorz.Wrap(err, errorz.Skip()))
}

// TraceHTTPRequestServer starts tracing an inbound HTTP request.
func (l *contextLogsImpl) TraceHTTPRequestServer(req *http.Request, reqBody []byte) (context.Context, func()) {
	return l.logs.TraceHTTPRequestServer(l.ctx, req, reqBody)
}

// SetUser sets the user in the current scope.
func (l *contextLogsImpl) SetUser(user *User) {
	l.logs.SetUser(l.ctx, user)
}

// AddMetadata adds the given metadata to the current scope.
func (l *contextLogsImpl) AddMetadata(k string, v interface{}) {
	l.logs.AddMetadata(l.ctx, k, v)
}

// Initializer is a Logs initializer which provides a default implementation using Sentry.
// It is possible to inject a different Logs implementation using NewSingletonInjector and a custom Initializer.
func Initializer(ctx context.Context) (injectz.Injector, injectz.Releaser) {
	cfg := ctx.Value(logsConfigContextKey).(*LogsConfig)
	errorz.MaybeMustWrap(validate.Struct(cfg))

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
		TracesSampleRate: cfg.SentrySampleRate,
		ServerName:       cfg.ServerName,
		Release:          cfg.Release,
		Environment:      cfg.Environment,
		BeforeSend: func(event *sentry.Event, _ *sentry.EventHint) *sentry.Event {
			logrusEntry := logrusLogger.
				WithTime(event.Timestamp).
				WithFields(event.Extra)

			if event.User.ID != "" {
				logrusEntry = logrusEntry.WithField("user-id", event.User.ID)
			}

			message := event.Message
			if len(event.Exception) > 0 {
				message = event.Exception[0].Value
			}

			logrusEntry.Log(levelFromSentry(event.Level).toLogrus(), message)
			return event
		},
	})
	errorz.MaybeMustWrap(err)
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
			sentryHub.Flush(cfg.ReleaseTimeout)
		}
}

// NewSingletonInjector always injects the given Logs.
func NewSingletonInjector(l Logs) injectz.Injector {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logsContextKey, l)
	}
}

// Get extracts the Logs from context and wraps it as ContextLogs. Panics if not found.
func Get(ctx context.Context) ContextLogs {
	return &contextLogsImpl{
		ctx:  ctx,
		logs: ctx.Value(logsContextKey).(Logs),
	}
}

// MaybeGet is like Get but returns nil if not found.
func MaybeGet(ctx context.Context) ContextLogs {
	logs, ok := ctx.Value(logsContextKey).(Logs)
	if !ok {
		return nil
	}

	return &contextLogsImpl{
		ctx:  ctx,
		logs: logs,
	}
}
