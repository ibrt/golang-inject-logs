package logz

import (
	"encoding/hex"
	"regexp"

	"github.com/getsentry/sentry-go"
)

const (
	sentryTraceHeader         = "sentry-trace"
	sentryMaxRequestBodyBytes = 10 * 1024
	logsRequestExtraKey       = "golang-inject-logs-request"
)

var (
	sentryTraceRegexp = regexp.MustCompile(`^([[:xdigit:]]{32})-([[:xdigit:]]{16})(?:-([01]))?$`)
)

func traceBeforeSend(event *sentry.Event) *sentry.Event {
	if event == nil {
		return nil
	}

	if req, ok := event.Extra[logsRequestExtraKey].(*sentry.Request); ok {
		if len(req.Data) > sentryMaxRequestBodyBytes {
			req.Data = ""
		}
		event.Request = req

		delete(event.Extra, logsRequestExtraKey)
	}
	return event
}

func newTraceSpanOption(headers map[string]string) sentry.SpanOption {
	return func(span *sentry.Span) {
		if headers != nil {
			if sentryTrace := headers[sentryTraceHeader]; sentryTrace != "" {
				if traceID, parentSpanID, sampled, ok := parseSentryTraceHeader(sentryTrace); ok {
					span.TraceID = traceID
					span.ParentSpanID = parentSpanID
					span.Sampled = sampled
				}
			}
		}
	}
}

func parseSentryTraceHeader(value string) (sentry.TraceID, sentry.SpanID, sentry.Sampled, bool) {
	var traceID sentry.TraceID
	var parentSpanID sentry.SpanID

	sampled := sentry.SampledFalse
	m := sentryTraceRegexp.FindSubmatch([]byte(value))

	if m == nil {
		return traceID, parentSpanID, sampled, false
	}

	_, _ = hex.Decode(traceID[:], m[1])
	_, _ = hex.Decode(parentSpanID[:], m[2])

	if len(m[3]) != 0 {
		switch m[3][0] {
		case '0':
			sampled = sentry.SampledFalse
		case '1':
			sampled = sentry.SampledTrue
		}
	}

	return traceID, parentSpanID, sampled, true
}
