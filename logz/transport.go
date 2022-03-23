package logz

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

type logsTransport struct {
	logrusLogger *logrus.Logger
	transport    sentry.Transport
}

func newLogsTransport(logrusLogger *logrus.Logger, transport sentry.Transport) *logsTransport {
	if transport == nil {
		transport = sentry.NewHTTPTransport()
	}

	return &logsTransport{
		logrusLogger: logrusLogger,
		transport:    transport,
	}
}

// Flush implements the sentry.Transport interface.
func (t *logsTransport) Flush(timeout time.Duration) bool {
	return t.transport.Flush(timeout)
}

// Configure implements the sentry.Transport interface.
func (t *logsTransport) Configure(options sentry.ClientOptions) {
	t.transport.Configure(options)
}

// SendEvent implements the sentry.Transport interface.
func (t *logsTransport) SendEvent(event *sentry.Event) {
	event = traceBeforeSend(event)

	logrusEntry := t.logrusLogger.
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
	t.transport.SendEvent(event)
}
