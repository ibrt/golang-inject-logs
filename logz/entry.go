package logz

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/ibrt/golang-inject-clock/clockz"
)

var (
	_ Option = OptionFunc(nil)
	_ Option = &Args{}
	_ Option = &Metadata{}
)

// Option describes an option which can be applied to an entry.
type Option interface {
	Apply(e *entry)
}

// OptionFunc describes an option which can be applied to an entry.
type OptionFunc func(e *entry)

// Apply implements the Option interface.
func (f OptionFunc) Apply(e *entry) {
	f(e)
}

// Args describes a list of args used for formatting an entry message.
type Args []interface{}

// Apply implements the Option interface.
func (a Args) Apply(_ *entry) {
	// intentionally empty
}

// A is a shorthand builder for args.
func A(a ...interface{}) Args {
	return a
}

// Metadata describes metadata which can be attached to an entry.
type Metadata map[string]interface{}

// Apply implements the Option interface.
func (m Metadata) Apply(e *entry) {
	for k, v := range m {
		e.metadata[k] = v
	}
}

// M is a shorthand for providing metadata to an entry.
func M(k string, v interface{}) OptionFunc {
	return func(e *entry) {
		e.metadata[k] = v
	}
}

type entry struct {
	level     Level
	timestamp time.Time
	callers   []uintptr
	message   string
	metadata  Metadata
}

func (e *entry) toSentryEvent() *sentry.Event {
	event := sentry.NewEvent()
	event.Level = sentry.Level(e.level)
	event.Timestamp = e.timestamp
	event.Message = e.message
	event.Extra = e.metadata

	event.Threads = []sentry.Thread{{
		Stacktrace: callersToSentryStacktrace(e.callers),
		Crashed:    false,
		Current:    true,
	}}

	return event
}

func newEntry(ctx context.Context, level Level, skipCallers int, format string, options ...Option) *entry {
	callers := make([]uintptr, 1024)
	callers = callers[:runtime.Callers(2+skipCallers, callers[:])]

	var mergedArgs []interface{}
	for _, option := range options {
		if args, ok := option.(Args); ok {
			mergedArgs = append(mergedArgs, args...)
		}
	}

	e := &entry{
		level:     level,
		timestamp: clockz.Get(ctx).Now(),
		callers:   callers,
		message:   fmt.Sprintf(format, mergedArgs...),
		metadata:  Metadata{},
	}

	for _, o := range options {
		o.Apply(e)
	}

	return e
}
